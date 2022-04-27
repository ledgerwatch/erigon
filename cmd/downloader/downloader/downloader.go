package downloader

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	common2 "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/dir"
	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/cmd/downloader/downloader/torrentcfg"
	"github.com/ledgerwatch/log/v3"
)

const ASSERT = false

type Protocols struct {
	TorrentClient *torrent.Client
	DB            kv.RwDB
	cfg           *torrentcfg.Cfg

	statsLock   *sync.RWMutex
	stats       AggStats
	snapshotDir *dir.Rw
}

func New(cfg *torrentcfg.Cfg, snapshotDir *dir.Rw) (*Protocols, error) {
	peerID, err := readPeerID(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("get peer id: %w", err)
	}
	cfg.PeerID = string(peerID)
	torrentClient, err := torrent.NewClient(cfg.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("fail to start torrent client: %w", err)
	}
	if len(peerID) == 0 {
		if err = savePeerID(cfg.DB, torrentClient.PeerID()); err != nil {
			return nil, fmt.Errorf("save peer id: %w", err)
		}
	}

	return &Protocols{
		cfg:           cfg,
		TorrentClient: torrentClient,
		DB:            cfg.DB,
		statsLock:     &sync.RWMutex{},
		snapshotDir:   snapshotDir,
	}, nil
}

func savePeerID(db kv.RwDB, peerID torrent.PeerID) error {
	return db.Update(context.Background(), func(tx kv.RwTx) error {
		return tx.Put(kv.BittorrentInfo, []byte(kv.BittorrentPeerID), peerID[:])
	})
}

func readPeerID(db kv.RoDB) (peerID []byte, err error) {
	if err = db.View(context.Background(), func(tx kv.Tx) error {
		peerIDFromDB, err := tx.GetOne(kv.BittorrentInfo, []byte(kv.BittorrentPeerID))
		if err != nil {
			return fmt.Errorf("get peer id: %w", err)
		}
		peerID = common2.Copy(peerIDFromDB)
		return nil
	}); err != nil {
		return nil, err
	}
	return peerID, nil
}

func (cli *Protocols) Start(ctx context.Context, silent bool) error {
	if err := CreateTorrentFilesAndAdd(ctx, cli.snapshotDir, cli.TorrentClient); err != nil {
		return fmt.Errorf("CreateTorrentFilesAndAdd: %w", err)
	}

	go func() {
		for {
			torrents := cli.TorrentClient.Torrents()
			for _, t := range torrents {
				<-t.GotInfo()
				if t.Complete.Bool() {
					continue
				}
				t.AllowDataDownload()
				select {
				case <-t.Complete.On():
					fmt.Printf("on: %s, %t\n", t.Name(), t.Complete.Bool())
					//case <-t.Complete.Off():
					//	fmt.Printf("off: %s\n", t.Name())
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		var m runtime.MemStats
		logEvery := time.NewTicker(20 * time.Second)
		defer logEvery.Stop()

		interval := 5 * time.Second
		statEvery := time.NewTicker(interval)
		defer statEvery.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-statEvery.C:
				cli.ReCalcStats(interval)

			case <-logEvery.C:
				if silent {
					continue
				}

				torrents := cli.TorrentClient.Torrents()
				stats := cli.Stats()

				fmt.Printf("alex111: %d, %d\n ", stats.MetadataReady, len(torrents))
				if stats.MetadataReady < len(torrents) {
					log.Info(fmt.Sprintf("[torrent] Waiting for torrents metadata: %d/%d", stats.MetadataReady, len(torrents)))
					continue
				}

				runtime.ReadMemStats(&m)
				if stats.Completed {
					log.Info("[torrent] Seeding",
						"download", common2.ByteCount(uint64(stats.DownloadRate))+"/s",
						"upload", common2.ByteCount(uint64(stats.UploadRate))+"/s",
						"unique_peers", stats.PeersCount,
						"files", stats.TorrentsCount,
						"alloc", common2.ByteCount(m.Alloc), "sys", common2.ByteCount(m.Sys))
					continue
				}

				log.Info("[torrent] Downloading",
					"Progress", fmt.Sprintf("%.2f%%", stats.Progress),
					"download", common2.ByteCount(uint64(stats.DownloadRate))+"/s",
					"upload", common2.ByteCount(uint64(stats.UploadRate))+"/s",
					"unique_peers", stats.PeersCount,
					"files", stats.TorrentsCount,
					"alloc", common2.ByteCount(m.Alloc), "sys", common2.ByteCount(m.Sys))
				if stats.PeersCount == 0 {
					ips := cli.TorrentClient.BadPeerIPs()
					if len(ips) > 0 {
						log.Info("[torrent] Stats", "banned", ips)
					}
				}
			}
		}
	}()
	return nil
}

func (cli *Protocols) ReCalcStats(interval time.Duration) {
	cli.statsLock.Lock()
	defer cli.statsLock.Unlock()
	prevStats, stats := cli.stats, cli.stats

	var aggBytesCompleted, aggLen int64
	peers := map[torrent.PeerID]*torrent.PeerConn{}
	torrents := cli.TorrentClient.Torrents()
	connStats := cli.TorrentClient.ConnStats()

	stats.BytesRead += uint64(connStats.BytesReadUsefulIntendedData.Int64())
	stats.BytesWritten += uint64(connStats.BytesWrittenData.Int64())

	stats.BytesTotal, stats.BytesCompleted, stats.PeerConnections, stats.MetadataReady = 0, 0, 0, 0
	for _, t := range torrents {
		select {
		case <-t.GotInfo():
			stats.MetadataReady++
		default:
			return // not all torrents are resolved yet
		}
		aggBytesCompleted += t.BytesCompleted()
		aggLen += t.Length()

		for _, peer := range t.PeerConns() {
			peers[peer.PeerID] = peer
		}

		stats.Completed = stats.Completed && t.Complete.Bool()
		stats.BytesCompleted += uint64(t.BytesCompleted())
		stats.BytesTotal += uint64(t.Info().TotalLength())
		stats.PeerConnections += uint64(len(t.PeerConns()))
	}

	stats.DownloadRate += (stats.BytesRead - prevStats.BytesRead) / uint64(interval.Seconds())
	stats.UploadRate += (stats.BytesWritten - prevStats.BytesWritten) / uint64(interval.Seconds())

	stats.Progress = float32(float64(100) * (float64(aggBytesCompleted) / float64(aggLen)))
	if stats.Progress == 100 && !stats.Completed {
		stats.Progress = 99.99
	}

	stats.PeersCount = int32(len(peers))
	stats.TorrentsCount = len(torrents)

	cli.stats = stats
}

func (cli *Protocols) Stats() AggStats {
	cli.statsLock.RLock()
	defer cli.statsLock.RUnlock()
	return cli.stats
}

func (cli *Protocols) Close() {
	for _, tr := range cli.TorrentClient.Torrents() {
		tr.Drop()
	}
	cli.TorrentClient.Close()
	cli.DB.Close()
	if cli.cfg.CompletionCloser != nil {
		cli.cfg.CompletionCloser.Close() //nolint
	}
}

func (cli *Protocols) PeerID() []byte {
	peerID := cli.TorrentClient.PeerID()
	return peerID[:]
}

func (cli *Protocols) StopSeeding(hash metainfo.Hash) error {
	t, ok := cli.TorrentClient.Torrent(hash)
	if !ok {
		return nil
	}
	ch := t.Closed()
	t.Drop()
	<-ch
	return nil
}

type AggStats struct {
	Completed                  bool
	BytesCompleted, BytesTotal uint64
	PeerConnections            uint64

	DownloadRate uint64
	UploadRate   uint64
	PeersCount   int32

	Progress                     float32
	TorrentsCount, MetadataReady int

	BytesRead    uint64
	BytesWritten uint64
}

func AddTorrentFile(ctx context.Context, torrentFilePath string, torrentClient *torrent.Client) (mi *metainfo.MetaInfo, err error) {
	mi, err = metainfo.LoadFromFile(torrentFilePath)
	if err != nil {
		return nil, err
	}
	mi.AnnounceList = Trackers

	t := time.Now()
	torrent, err := torrentClient.AddTorrent(mi)
	if err != nil {
		return mi, err
	}
	torrent.DisallowDataDownload()
	took := time.Since(t)
	if took > 3*time.Second {
		log.Info("[torrent] Check validity", "file", torrentFilePath, "took", took)
	}
	return mi, nil
}

// AddTorrentFiles - adding .torrent files to torrentClient (and checking their hashes), if .torrent file
// added first time - pieces verification process will start (disk IO heavy) - Progress
// kept in `piece completion storage` (surviving reboot). Once it done - no disk IO needed again.
// Don't need call torrent.VerifyData manually
func AddTorrentFiles(ctx context.Context, snapshotsDir *dir.Rw, torrentClient *torrent.Client) error {
	files, err := AllTorrentPaths(snapshotsDir.Path)
	if err != nil {
		return err
	}
	for _, torrentFilePath := range files {
		if _, err := AddTorrentFile(ctx, torrentFilePath, torrentClient); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

	}

	return nil
}

// ResolveAbsentTorrents - add hard-coded hashes (if client doesn't have) as magnet links and download everything
func ResolveAbsentTorrents(ctx context.Context, cli *Protocols, preverifiedHashes []metainfo.Hash, snapshotDir *dir.Rw) error {
	mi := &metainfo.MetaInfo{AnnounceList: Trackers}
	torrentClient := cli.TorrentClient
	go func() {
		for i := range preverifiedHashes {
			if _, ok := torrentClient.Torrent(preverifiedHashes[i]); ok {
				continue
			}
			magnet := mi.Magnet(&preverifiedHashes[i], nil)
			t, err := torrentClient.AddMagnet(magnet.String())
			if err != nil {
				_ = err
			}
			t.DisallowDataDownload()
			t.AllowDataUpload()
		}
	}()

	logEvery := time.NewTicker(10 * time.Second)
	defer logEvery.Stop()
	gotInfo := 0
	torrents := torrentClient.Torrents()

	for _, t := range torrents {

	LogLoop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-logEvery.C:
				if gotInfo < len(torrents) {
					log.Info(fmt.Sprintf("[torrent] Waiting for torrents metadata: %d/%d", gotInfo, len(torrents)))
					continue
				}
			case <-t.GotInfo():
				gotInfo++
				mi := t.Metainfo()
				if err := CreateTorrentFileIfNotExists(snapshotDir, t.Info(), &mi); err != nil {
					return err
				}
				break LogLoop
			}
		}
	}

	return nil
}

func VerifyDtaFiles(ctx context.Context, snapshotDir string) error {
	logEvery := time.NewTicker(5 * time.Second)
	defer logEvery.Stop()
	files, err := AllTorrentPaths(snapshotDir)
	if err != nil {
		return err
	}
	totalPieces := 0
	for _, f := range files {
		metaInfo, err := metainfo.LoadFromFile(f)
		if err != nil {
			return err
		}
		info, err := metaInfo.UnmarshalInfo()
		if err != nil {
			return err
		}
		totalPieces += info.NumPieces()
	}

	j := 0
	for _, f := range files {
		metaInfo, err := metainfo.LoadFromFile(f)
		if err != nil {
			return err
		}
		info, err := metaInfo.UnmarshalInfo()
		if err != nil {
			return err
		}

		err = verifyTorrent(&info, snapshotDir, func(i int, good bool) error {
			j++
			if !good {
				log.Error("[torrent] Verify hash mismatch", "at piece", i, "file", f)
				return fmt.Errorf("invalid file")
			}
			select {
			case <-logEvery.C:
				log.Info("[torrent] Verify", "Progress", fmt.Sprintf("%.2f%%", 100*float64(j)/float64(totalPieces)))
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	log.Info("[torrent] Verify succeed")
	return nil
}
