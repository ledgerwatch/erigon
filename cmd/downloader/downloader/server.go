package downloader

import (
	"context"
	"errors"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/ledgerwatch/erigon-lib/gointerfaces"
	proto_downloader "github.com/ledgerwatch/erigon-lib/gointerfaces/downloader"
	prototypes "github.com/ledgerwatch/erigon-lib/gointerfaces/types"
	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/turbo/snapshotsync/snapshothashes"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrNotSupportedNetworkID = errors.New("not supported network id")
	ErrNotSupportedSnapshot  = errors.New("not supported snapshot for this network id")
)
var (
	_ proto_downloader.DownloaderServer = &SNDownloaderServer{}
)

func NewServer(db kv.RwDB, client *Client) (*SNDownloaderServer, error) {
	sn := &SNDownloaderServer{
		db: db,
		t:  client,
	}
	return sn, nil
}

func (s *SNDownloaderServer) Load(ctx context.Context) error {
	if err := BuildTorrentFilesIfNeed(ctx, s.t.snapshotsDir); err != nil {
		return err
	}
	preverifiedHashes := snapshothashes.Goerli // TODO: remove hard-coded hashes from downloader
	if err := AddTorrentFiles(ctx, s.t.snapshotsDir, s.t.Cli, preverifiedHashes); err != nil {
		return err
	}
	return nil
}

type SNDownloaderServer struct {
	proto_downloader.UnimplementedDownloaderServer
	t  *Client
	db kv.RwDB
}

func (s *SNDownloaderServer) Download(ctx context.Context, request *proto_downloader.DownloadRequest) (*emptypb.Empty, error) {
	infoHashes := Proto2InfoHashes(request.TorrentHashes)

	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()
	if err := ResolveAbsentTorrents(ctx, s.t.Cli, infoHashes); err != nil {
		return nil, err
	}
	for _, t := range s.t.Cli.Torrents() {
		t.AllowDataDownload()
		t.AllowDataUpload()
	}
	return &emptypb.Empty{}, nil
}

func (s *SNDownloaderServer) Snapshots(ctx context.Context, request *proto_downloader.SnapshotsRequest) (*proto_downloader.SnapshotsInfoReply, error) {
	torrents := s.t.Cli.Torrents()
	infoItems := make([]*proto_downloader.SnapshotsInfo, len(torrents))
	for i, t := range torrents {
		readiness := int32(0)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.GotInfo():
			readiness = int32(100 * (float64(t.BytesCompleted()) / float64(t.Info().TotalLength())))
		default:
		}

		infoItems[i] = &proto_downloader.SnapshotsInfo{
			Readiness: readiness,
			Path:      t.Name(),
		}
	}
	return &proto_downloader.SnapshotsInfoReply{Info: infoItems}, nil
}

func (s *SNDownloaderServer) Stats(ctx context.Context) map[string]torrent.TorrentStats {
	stats := map[string]torrent.TorrentStats{}
	torrents := s.t.Cli.Torrents()
	for _, t := range torrents {
		stats[t.Name()] = t.Stats()
	}
	return stats
}

func Proto2InfoHashes(in []*prototypes.H160) []metainfo.Hash {
	infoHashes := make([]metainfo.Hash, len(in))
	i := 0
	for _, h := range in {
		infoHashes[i] = gointerfaces.ConvertH160toAddress(h)
		i++
	}
	return infoHashes
}

func InfoHashes2Proto(in []metainfo.Hash) []*prototypes.H160 {
	infoHashes := make([]*prototypes.H160, len(in))
	i := 0
	for _, h := range in {
		infoHashes[i] = gointerfaces.ConvertAddressToH160(h)
		i++
	}
	return infoHashes
}
