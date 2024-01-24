package heimdall

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ledgerwatch/log/v3"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/polygon/bor/borcfg"
)

// Heimdall is a wrapper of Heimdall HTTP API
//
//go:generate mockgen -destination=./heimdall_mock.go -package=sync . Heimdall
type Heimdall interface {
	LastCheckpointId(ctx context.Context, io CheckpointIO) (CheckpointId, bool, error)
	LastMilestoneId(ctx context.Context, io MilestoneIO) (MilestoneId, bool, error)
	LastSpanId(ctx context.Context, io SpanIO) (SpanId, bool, error)

	FetchCheckpoints(ctx context.Context, io CheckpointIO, start CheckpointId, end CheckpointId) (Waypoints, error)
	FetchMilestones(ctx context.Context, io MilestoneIO, start MilestoneId, end MilestoneId) (Waypoints, error)
	FetchSpans(ctx context.Context, io SpanIO, start SpanId, end SpanId) ([]*Span, error)

	FetchCheckpointsFromBlock(ctx context.Context, io CheckpointIO, startBlock uint64) (Waypoints, error)
	FetchMilestonesFromBlock(ctx context.Context, io MilestoneIO, startBlock uint64) (Waypoints, error)
	FetchSpansFromBlock(ctx context.Context, io SpanIO, startBlock uint64) ([]*Span, error)

	OnCheckpointEvent(ctx context.Context, io CheckpointIO, callback func(*Checkpoint)) error
	OnMilestoneEvent(ctx context.Context, io MilestoneIO, callback func(*Milestone)) error
	OnSpanEvent(ctx context.Context, io SpanIO, callback func(*Span)) error
}

// ErrIncompleteMilestoneRange happens when FetchMilestones is called with an old start block because old milestones are evicted
var ErrIncompleteMilestoneRange = errors.New("milestone range doesn't contain the start block")

type heimdallImpl struct {
	client          HeimdallClient
	pollDelay       time.Duration
	lastBlockNumber chan uint64
	cfg             *borcfg.BorConfig
	logger          log.Logger
}

func NewHeimdall(client HeimdallClient, logger log.Logger) Heimdall {
	h := heimdallImpl{
		client:    client,
		pollDelay: time.Second,
		logger:    logger,
	}
	return &h
}

func (h *heimdallImpl) LastCheckpointId(ctx context.Context, io CheckpointIO) (CheckpointId, bool, error) {
	// todo get this from io if its likely not changed (need timeout)

	count, err := h.client.FetchCheckpointCount(ctx)

	if err != nil {
		return 0, false, err
	}

	return CheckpointId(count), true, nil
}

func (h *heimdallImpl) FetchCheckpointsFromBlock(ctx context.Context, io CheckpointIO, startBlock uint64) (Waypoints, error) {
	count, _, err := h.LastCheckpointId(ctx, io)

	if err != nil {
		return nil, err
	}

	var checkpoints []Waypoint

	for i := count; i >= 1; i-- {
		c, err := h.FetchCheckpoints(ctx, io, i, i)
		if err != nil {
			return nil, err
		}

		cmpResult := c[0].CmpRange(startBlock)
		// the start block is past the last checkpoint
		if cmpResult > 0 {
			return nil, nil
		}

		checkpoints = append(checkpoints, c...)

		// the checkpoint contains the start block
		if cmpResult == 0 {
			break
		}
	}

	common.SliceReverse(checkpoints)
	return checkpoints, nil
}

func (h *heimdallImpl) FetchCheckpoints(ctx context.Context, io CheckpointIO, start CheckpointId, end CheckpointId) (Waypoints, error) {
	return nil, fmt.Errorf("TODO")
}

func (h *heimdallImpl) LastMilestoneId(ctx context.Context, io MilestoneIO) (MilestoneId, bool, error) {
	// todo get this from io if its likely not changed (need timeout)

	count, err := h.client.FetchMilestoneCount(ctx)

	if err != nil {
		return 0, false, err
	}

	return MilestoneId(count), true, nil
}

func (h *heimdallImpl) FetchMilestonesFromBlock(ctx context.Context, io MilestoneIO, startBlock uint64) (Waypoints, error) {
	last, _, err := h.LastMilestoneId(ctx, io)

	if err != nil {
		return nil, err
	}

	var milestones Waypoints

	for i := last; i >= 1; i-- {
		m, err := h.client.FetchMilestone(ctx, int64(i))
		if err != nil {
			if errors.Is(err, ErrNotInMilestoneList) {
				common.SliceReverse(milestones)
				return milestones, ErrIncompleteMilestoneRange
			}
			return nil, err
		}

		cmpResult := m.CmpRange(startBlock)
		// the start block is past the last milestone
		if cmpResult > 0 {
			return nil, nil
		}

		milestones = append(milestones, m)

		// the checkpoint contains the start block
		if cmpResult == 0 {
			break
		}
	}

	common.SliceReverse(milestones)
	return milestones, nil
}

func (h *heimdallImpl) FetchMilestones(ctx context.Context, io MilestoneIO, start MilestoneId, end MilestoneId) (Waypoints, error) {
	return nil, fmt.Errorf("TODO")
}

func (h *heimdallImpl) LastSpanId(ctx context.Context, io SpanIO) (SpanId, bool, error) {
	return 0, false, fmt.Errorf("TODO")
}

func (h *heimdallImpl) FetchSpansFromBlock(ctx context.Context, io SpanIO, startBlock uint64) ([]*Span, error) {
	return nil, fmt.Errorf("TODO")
}

func (h *heimdallImpl) FetchSpans(ctx context.Context, io SpanIO, start SpanId, end SpanId) ([]*Span, error) {
	var spans []*Span

	lastSpanId, exists, err := io.LastSpanId(ctx)

	if err != nil {
		return nil, err
	}

	if exists && start <= lastSpanId {
		if end <= lastSpanId {
			lastSpanId = end
		}

		for id := start; id <= lastSpanId; id++ {
			span, err := io.ReadSpan(ctx, id)

			if err != nil {
				return nil, err
			}

			spans = append(spans, span)
		}

		start = lastSpanId + 1
	}

	for id := start; id <= end; id++ {
		span, err := h.client.Span(ctx, uint64(id))

		if err != nil {
			return nil, err
		}

		err = io.WriteSpan(ctx, span)

		if err != nil {
			return nil, err
		}

		spans = append(spans, span)
	}

	return spans, nil
}

func (h *heimdallImpl) OnSpanEvent(ctx context.Context, io SpanIO, callback func(*Span)) error {
	go func() {
		var blockNumber uint64

		for {
			select {
			case <-ctx.Done():
				return
			case lastBlockNumber := <-h.lastBlockNumber:
				if lastBlockNumber > blockNumber {
					blockNumber = lastBlockNumber

					requiredSpanId := SpanIdAt(blockNumber)

					if IsBlockInLastSprintOfSpan(blockNumber, h.cfg) {
						requiredSpanId++

						if spans, err := h.FetchSpans(ctx, io, requiredSpanId, requiredSpanId); err == nil && len(spans) > 0 {
							go callback(spans[0])
						}
					}
				}
			}
		}
	}()

	return nil
}

func (h *heimdallImpl) OnCheckpointEvent(ctx context.Context, io CheckpointIO, callback func(*Checkpoint)) error {
	return fmt.Errorf("TODO")
}

func (h *heimdallImpl) OnMilestoneEvent(ctx context.Context, io MilestoneIO, callback func(*Milestone)) error {
	currentCount, err := h.client.FetchMilestoneCount(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			count, err := h.client.FetchMilestoneCount(ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					h.logger.Error("heimdallImpl.OnMilestoneEvent FetchMilestoneCount error", "err", err)
				}
				break
			}

			if count <= currentCount {
				pollDelayTimer := time.NewTimer(h.pollDelay)
				select {
				case <-ctx.Done():
					return
				case <-pollDelayTimer.C:
				}
			} else {
				currentCount = count
				m, err := h.client.FetchMilestone(ctx, count)
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						h.logger.Error("heimdallImpl.OnMilestoneEvent FetchMilestone error", "err", err)
					}
					break
				}

				go callback(m)
			}
		}
	}()

	return nil
}
