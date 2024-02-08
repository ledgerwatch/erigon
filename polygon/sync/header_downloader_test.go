package sync

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ledgerwatch/log/v3"
	"github.com/stretchr/testify/require"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/polygon/heimdall"
	"github.com/ledgerwatch/erigon/polygon/p2p"
	"github.com/ledgerwatch/erigon/turbo/testlog"
)

func newHeaderDownloaderTest(t *testing.T) *headerDownloaderTest {
	return newHeaderDownloaderTestWithOpts(t, headerDownloaderTestOpts{})
}

func newHeaderDownloaderTestWithOpts(t *testing.T, opts headerDownloaderTestOpts) *headerDownloaderTest {
	ctrl := gomock.NewController(t)
	checkpointStore := NewMockCheckpointStore(ctrl)
	milestoneStore := NewMockMilestoneStore(ctrl)
	heimdall := heimdall.NewMockHeimdall(ctrl)
	peerManager := p2p.NewMockPeerManager(ctrl)
	peerManager.EXPECT().MaxPeers().Return(100).Times(1)
	logger := testlog.Logger(t, log.LvlDebug)
	headerVerifier := opts.getOrCreateDefaultHeaderVerifier()
	headerDownloader := NewHeaderDownloader(logger, peerManager, heimdall, headerVerifier)
	return &headerDownloaderTest{
		heimdall:         heimdall,
		peerManager:      peerManager,
		headerDownloader: headerDownloader,
		milestoneStore:   milestoneStore,
		checkpointStore:  checkpointStore,
	}
}

type headerDownloaderTestOpts struct {
	headerVerifier AccumulatedHeadersVerifier
}

func (opts headerDownloaderTestOpts) getOrCreateDefaultHeaderVerifier() AccumulatedHeadersVerifier {
	if opts.headerVerifier == nil {
		return func(_ heimdall.Waypoint, _ []*types.Header) error {
			return nil
		}
	}

	return opts.headerVerifier
}

type headerDownloaderTest struct {
	heimdall         *heimdall.MockHeimdall
	peerManager      *p2p.MockPeerManager
	milestoneStore   *MockMilestoneStore
	checkpointStore  *MockCheckpointStore
	headerDownloader *HeaderDownloader
}

func (hdt headerDownloaderTest) fakePeers(count int, blockNums ...uint64) p2p.PeersSyncProgress {
	peers := make(p2p.PeersSyncProgress, count)
	for i := range peers {
		var blockNum uint64
		if i < len(blockNums) {
			blockNum = blockNums[i]
		} else {
			blockNum = math.MaxUint64
		}

		pidBytes := make([]byte, 64)
		binary.LittleEndian.PutUint64(pidBytes, uint64(i))
		peers[i] = &p2p.PeerSyncProgress{
			Id:              p2p.PeerId(pidBytes),
			MaxSeenBlockNum: blockNum,
		}
	}

	return peers
}

func (hdt headerDownloaderTest) fakeCheckpoints(count int) heimdall.Waypoints {
	checkpoints := make(heimdall.Waypoints, count)
	for i := range checkpoints {
		num := i + 1
		checkpoints[i] = &heimdall.Checkpoint{
			Fields: heimdall.WaypointFields{
				StartBlock: big.NewInt(int64(num)),
				EndBlock:   big.NewInt(int64(num)),
				RootHash:   common.BytesToHash([]byte(fmt.Sprintf("0x%d", num))),
			},
		}
	}

	return checkpoints
}

func (hdt headerDownloaderTest) fakeMilestones(count int) heimdall.Waypoints {
	milestones := make(heimdall.Waypoints, count)
	for i := range milestones {
		num := i + 1
		milestones[i] = &heimdall.Milestone{
			Fields: heimdall.WaypointFields{
				StartBlock: big.NewInt(int64(num)),
				EndBlock:   big.NewInt(int64(num)),
				RootHash:   common.BytesToHash([]byte(fmt.Sprintf("0x%d", num))),
			},
		}
	}

	return milestones
}

type downloadHeadersMock func(context.Context, *big.Int, *big.Int, string) ([]*types.Header, error)

func (hdt headerDownloaderTest) defaultDownloadHeadersMock() downloadHeadersMock {
	return func(ctx context.Context, start *big.Int, end *big.Int, peerID string) ([]*types.Header, error) {
		res := make([]*types.Header, new(big.Int).Sub(end, start).Uint64()+1)
		for i := new(big.Int).Set(start); i.Cmp(end) < 1; i.Add(i, new(big.Int).SetUint64(1)) {
			res[new(big.Int).Sub(i, start).Uint64()] = &types.Header{Number: new(big.Int).Set(i)}
		}
		return res, nil
	}
}

func (hdt headerDownloaderTest) defaultWriteHeadersMock(capture *[]*types.Header) func([]*types.Header) error {
	return func(headers []*types.Header) error {
		*capture = append(*capture, headers...)
		return nil
	}
}

func TestHeaderDownloadUsingMilestones(t *testing.T) {
	test := newHeaderDownloaderTest(t)
	test.heimdall.EXPECT().
		FetchMilestonesFromBlock(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(test.fakeMilestones(4), nil).
		Times(1)
	test.peerManager.EXPECT().
		PeersSyncProgress().
		Return(test.fakePeers(8)).
		Times(1)
	test.peerManager.EXPECT().
		DownloadHeaders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(test.defaultDownloadHeadersMock()).
		Times(4)
	var persistedHeaders []*types.Header
	test.milestoneStore.EXPECT().
		PutHeaders(gomock.Any()).
		DoAndReturn(test.defaultWriteHeadersMock(&persistedHeaders)).
		Times(1)

	err := test.headerDownloader.DownloadUsingMilestones(context.Background(), test.milestoneStore, 1)
	require.NoError(t, err)
	require.Len(t, persistedHeaders, 4)
	// check headers are written in order
	require.Equal(t, uint64(1), persistedHeaders[0].Number.Uint64())
	require.Equal(t, uint64(2), persistedHeaders[1].Number.Uint64())
	require.Equal(t, uint64(3), persistedHeaders[2].Number.Uint64())
	require.Equal(t, uint64(4), persistedHeaders[3].Number.Uint64())
}

func TestHeaderDownloadUsingCheckpoints(t *testing.T) {
	test := newHeaderDownloaderTest(t)
	test.heimdall.EXPECT().
		FetchCheckpointsFromBlock(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(test.fakeCheckpoints(8), nil).
		Times(1)
	test.peerManager.EXPECT().
		PeersSyncProgress().
		Return(test.fakePeers(2)).
		Times(4)
	test.peerManager.EXPECT().
		DownloadHeaders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(test.defaultDownloadHeadersMock()).
		Times(8)
	var persistedHeaders []*types.Header
	test.checkpointStore.EXPECT().
		PutHeaders(gomock.Any()).
		DoAndReturn(test.defaultWriteHeadersMock(&persistedHeaders)).
		Times(4)

	err := test.headerDownloader.DownloadUsingCheckpoints(context.Background(), test.checkpointStore, 1)
	require.NoError(t, err)
	require.Len(t, persistedHeaders, 8)
	// check headers are written in order
	require.Equal(t, uint64(1), persistedHeaders[0].Number.Uint64())
	require.Equal(t, uint64(2), persistedHeaders[1].Number.Uint64())
	require.Equal(t, uint64(3), persistedHeaders[2].Number.Uint64())
	require.Equal(t, uint64(4), persistedHeaders[3].Number.Uint64())
	require.Equal(t, uint64(5), persistedHeaders[4].Number.Uint64())
	require.Equal(t, uint64(6), persistedHeaders[5].Number.Uint64())
	require.Equal(t, uint64(7), persistedHeaders[6].Number.Uint64())
	require.Equal(t, uint64(8), persistedHeaders[7].Number.Uint64())
}

func TestHeaderDownloadWhenInvalidStateThenPenalizePeerAndReDownload(t *testing.T) {
	var firstTimeInvalidReturned bool
	firstTimeInvalidReturnedPtr := &firstTimeInvalidReturned
	test := newHeaderDownloaderTestWithOpts(t, headerDownloaderTestOpts{
		headerVerifier: func(waypoint heimdall.Waypoint, headers []*types.Header) error {
			if waypoint.StartBlock().Cmp(new(big.Int).SetUint64(2)) == 0 && !*firstTimeInvalidReturnedPtr {
				*firstTimeInvalidReturnedPtr = true
				return errors.New("invalid checkpoint")
			}
			return nil
		},
	})
	test.heimdall.EXPECT().
		FetchCheckpointsFromBlock(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(test.fakeCheckpoints(6), nil).
		Times(1)
	test.peerManager.EXPECT().
		PeersSyncProgress().
		Return(test.fakePeers(3)).
		Times(3)
	test.peerManager.EXPECT().
		DownloadHeaders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(test.defaultDownloadHeadersMock()).
		// request 1,2,3 in parallel
		// -> 2 fails
		// requests 2,3,4 in parallel
		// 3 is cached
		// requests 5,6 in parallel
		// in total 6 requests + 1 request for re-requesting checkpoint 2
		// total = 7 (note this also tests caching works)
		Times(7)
	test.peerManager.EXPECT().
		Penalize(gomock.Eq("peer2")).
		Times(1)
	var persistedHeadersFirstTime, persistedHeadersRemaining []*types.Header
	gomock.InOrder(
		test.checkpointStore.EXPECT().
			PutHeaders(gomock.Any()).
			DoAndReturn(test.defaultWriteHeadersMock(&persistedHeadersFirstTime)).
			Times(1),
		test.checkpointStore.EXPECT().
			PutHeaders(gomock.Any()).
			DoAndReturn(test.defaultWriteHeadersMock(&persistedHeadersRemaining)).
			Times(2),
	)

	err := test.headerDownloader.DownloadUsingCheckpoints(context.Background(), test.checkpointStore, 1)
	require.NoError(t, err)
	require.Len(t, persistedHeadersFirstTime, 1)
	require.Len(t, persistedHeadersRemaining, 5)
}

func TestHeaderDownloadWhenZeroPeersTriesAgain(t *testing.T) {
	test := newHeaderDownloaderTest(t)
	test.heimdall.EXPECT().
		FetchCheckpointsFromBlock(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(test.fakeCheckpoints(8), nil).
		Times(1)
	test.peerManager.EXPECT().
		DownloadHeaders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(test.defaultDownloadHeadersMock()).
		Times(8)
	var persistedHeaders []*types.Header
	test.checkpointStore.EXPECT().
		PutHeaders(gomock.Any()).
		DoAndReturn(test.defaultWriteHeadersMock(&persistedHeaders)).
		Times(4)
	gomock.InOrder(
		// first, no peers at all
		test.peerManager.EXPECT().
			PeersSyncProgress().
			Return(nil).
			Times(1),
		// second, 2 peers but not synced enough for us to use
		test.peerManager.EXPECT().
			PeersSyncProgress().
			Return(test.fakePeers(2, 0, 0)).
			Times(1),
		// then, 2 fully synced peers that we can use
		test.peerManager.EXPECT().
			PeersSyncProgress().
			Return(test.fakePeers(2)).
			Times(4),
	)

	err := test.headerDownloader.DownloadUsingCheckpoints(context.Background(), test.checkpointStore, 1)
	require.NoError(t, err)
	require.Len(t, persistedHeaders, 8)
}
