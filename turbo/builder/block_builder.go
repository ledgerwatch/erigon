package builder

import (
	"sync"
	"sync/atomic"

	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/core"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/log/v3"
)

type BlockBuilderFunc func(tx kv.Tx, param *core.BlockBuilderParameters, interrupt *int32) (*types.Block, error)

// BlockBuilder wraps a goroutine that builds Proof-of-Stake payloads (PoS "mining")
type BlockBuilder struct {
	interrupt   int32
	syncCond    *sync.Cond
	block       *types.Block
	err         error
	emptyHeader *types.Header
}

// Creates a new BlockBuilder.
// BlockBuilder is responsible for rolling back the transaction eventually.
func NewBlockBuilder(tx kv.Tx, build BlockBuilderFunc, param *core.BlockBuilderParameters, emptyHeader *types.Header) *BlockBuilder {
	b := new(BlockBuilder)
	b.syncCond = sync.NewCond(new(sync.Mutex))
	b.emptyHeader = emptyHeader

	go func() {
		block, err := build(tx, param, &b.interrupt)
		tx.Rollback()

		b.syncCond.L.Lock()
		defer b.syncCond.L.Unlock()
		b.block = block
		b.err = err
		b.syncCond.Broadcast()
	}()

	return b
}

func (b *BlockBuilder) Stop() *types.Block {
	atomic.StoreInt32(&b.interrupt, 1)

	b.syncCond.L.Lock()
	defer b.syncCond.L.Unlock()
	for b.block == nil && b.err == nil {
		b.syncCond.Wait()
	}

	if b.err != nil {
		log.Error("BlockBuilder", "err", b.err)
		return types.NewBlock(b.emptyHeader, nil, nil, nil)
	}

	return b.block
}
