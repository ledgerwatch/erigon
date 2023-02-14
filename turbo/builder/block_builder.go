package builder

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/ledgerwatch/log/v3"

	"github.com/ledgerwatch/erigon/core"
	"github.com/ledgerwatch/erigon/core/types"
)

type BlockBuilderFunc func(param *core.BlockBuilderParameters, interrupt *int32) (*types.BlockWithReceipts, error)

// BlockBuilder wraps a goroutine that builds Proof-of-Stake payloads (PoS "mining")
type BlockBuilder struct {
	interrupt int32
	syncCond  *sync.Cond
	block     *types.BlockWithReceipts
	err       error
}

func NewBlockBuilder(build BlockBuilderFunc, param *core.BlockBuilderParameters) *BlockBuilder {
	builder := new(BlockBuilder)
	builder.syncCond = sync.NewCond(new(sync.Mutex))

	go func() {
		log.Info("Building block...")
		t := time.Now()
		block, err := build(param, &builder.interrupt)
		if err != nil {
			log.Warn("Failed to build a block", "err", err)
		} else {
			b := block.Block
			log.Info("Built block", "hash", b.Hash(), "height", b.NumberU64(), "txs", len(b.Transactions()), "gas used %", 100*float64(b.GasUsed())/float64(b.GasLimit()), "time", time.Since(t))
		}

		builder.syncCond.L.Lock()
		defer builder.syncCond.L.Unlock()
		builder.block = block
		builder.err = err
		builder.syncCond.Broadcast()
	}()

	return builder
}

func (b *BlockBuilder) Stop() (*types.BlockWithReceipts, error) {
	atomic.StoreInt32(&b.interrupt, 1)

	b.syncCond.L.Lock()
	defer b.syncCond.L.Unlock()
	for b.block == nil && b.err == nil {
		b.syncCond.Wait()
	}

	return b.block, b.err
}

func (b *BlockBuilder) Block() *types.Block {
	b.syncCond.L.Lock()
	defer b.syncCond.L.Unlock()

	if b.block == nil {
		return nil
	}
	return b.block.Block
}
