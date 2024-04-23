package p2p

import (
	"context"
	"errors"
	"fmt"

	"github.com/ledgerwatch/log/v3"

	"github.com/ledgerwatch/erigon/core/types"
)

func NewPenalizingFetcher(logger log.Logger, fetcher Fetcher, peerPenalizer PeerPenalizer) Fetcher {
	return newPenalizingFetcher(logger, fetcher, peerPenalizer)
}

func newPenalizingFetcher(logger log.Logger, fetcher Fetcher, peerPenalizer PeerPenalizer) *penalizingFetcher {
	return &penalizingFetcher{
		Fetcher:       fetcher,
		logger:        logger,
		peerPenalizer: peerPenalizer,
	}
}

type penalizingFetcher struct {
	Fetcher
	logger        log.Logger
	peerPenalizer PeerPenalizer
}

func (pf *penalizingFetcher) FetchHeaders(ctx context.Context, start uint64, end uint64, peerId *PeerId) ([]*types.Header, int, error) {
	headers, size, err := pf.Fetcher.FetchHeaders(ctx, start, end, peerId)
	if err != nil {
		return nil, 0, pf.maybePenalize(ctx, peerId, err, &ErrTooManyHeaders{}, &ErrNonSequentialHeaderNumbers{})
	}

	return headers, size, nil
}

func (pf *penalizingFetcher) FetchBodies(ctx context.Context, headers []*types.Header, peerId *PeerId) ([]*types.Body, int, error) {
	bodies, size, err := pf.Fetcher.FetchBodies(ctx, headers, peerId)
	if err != nil {
		return nil, 0, pf.maybePenalize(ctx, peerId, err, &ErrTooManyBodies{})
	}

	return bodies, size, nil
}

func (pf *penalizingFetcher) maybePenalize(ctx context.Context, peerId *PeerId, err error, penalizeErrs ...error) error {
	var shouldPenalize bool
	for _, penalizeErr := range penalizeErrs {
		if errors.Is(err, penalizeErr) {
			shouldPenalize = true
			break
		}
	}

	if shouldPenalize {
		pf.logger.Debug(
			"[p2p.penalizing.fetcher] penalizing peer - penalize-able fetcher issue",
			"peerId", peerId,
			"err", err,
		)

		if penalizeErr := pf.peerPenalizer.Penalize(ctx, peerId); penalizeErr != nil {
			err = fmt.Errorf("%w: %w", penalizeErr, err)
		}
	}

	return err
}
