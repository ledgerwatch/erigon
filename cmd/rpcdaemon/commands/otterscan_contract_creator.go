package commands

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/cmp"
	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon-lib/kv/bitmapdb"
	"github.com/ledgerwatch/erigon-lib/kv/rawdbv3"
	"github.com/ledgerwatch/erigon-lib/kv/temporal/historyv2"
	"github.com/ledgerwatch/erigon/core/state/temporal"
	"github.com/ledgerwatch/erigon/eth/stagedsync/stages"
	"github.com/ledgerwatch/log/v3"

	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/erigon/core/types/accounts"
)

type ContractCreatorData struct {
	Tx      libcommon.Hash    `json:"hash"`
	Creator libcommon.Address `json:"creator"`
}

func (api *OtterscanAPIImpl) GetContractCreator(ctx context.Context, addr libcommon.Address) (*ContractCreatorData, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	reader := state.NewPlainStateReader(tx)
	plainStateAcc, err := reader.ReadAccountData(addr)
	if err != nil {
		return nil, err
	}

	// No state == non existent
	if plainStateAcc == nil {
		return nil, nil
	}

	// EOA?
	if plainStateAcc.IsEmptyCodeHash() {
		return nil, nil
	}

	chainConfig, err := api.chainConfig(tx)
	if err != nil {
		return nil, err
	}

	var acc accounts.Account
	if api.historyV3(tx) {
		ttx := tx.(kv.TemporalTx)
		headNumber, err := stages.GetStageProgress(tx, stages.Execution)
		if err != nil {
			return nil, err
		}
		lastTxNum, _ := rawdbv3.TxNums.Max(tx, headNumber)

		// The sort.Search function finds the first block where the incarnation has
		// changed to the desired one, so we get the previous block from the bitmap;
		// however if the creationTxnID block is already the first one from the bitmap, it means
		// the block we want is the max block from the previous shard.
		var creationTxnID uint64
		var searchErr error

		_ = sort.Search(int(lastTxNum), func(i int) bool {
			v, ok, err := ttx.DomainGet(temporal.AccountsDomain, addr[:], nil, uint64(i))
			if err != nil {
				log.Error("Unexpected error, couldn't find changeset", "txNum", i, "addr", addr)
				panic(err)
			}
			if !ok {
				return false
			}
			if len(v) == 0 {
				creationTxnID = cmp.Max(creationTxnID, uint64(i))
				return false
			}

			if err := acc.DecodeForStorage(v); err != nil {
				searchErr = err
				return false
			}
			if acc.Incarnation < plainStateAcc.Incarnation {
				creationTxnID = cmp.Max(creationTxnID, uint64(i))
				return false
			}
			return true
		})
		if searchErr != nil {
			return nil, searchErr
		}

		// Trace block, find tx and contract creator
		tracer := NewCreateTracer(ctx, addr)
		_, bn, _ := rawdbv3.TxNums.FindBlockNum(tx, creationTxnID)
		minTxNum, _ := rawdbv3.TxNums.Min(tx, bn)
		txIndex := creationTxnID - minTxNum - 1 /* system-contract */
		if err := api.genericTracer(tx, ctx, bn, creationTxnID, txIndex, chainConfig, tracer); err != nil {
			return nil, err
		}
		return &ContractCreatorData{
			Tx:      tracer.Tx.Hash(),
			Creator: tracer.Creator,
		}, nil
	}

	// Contract; search for creation tx; navigate forward on AccountsHistory/ChangeSets
	//
	// We search shards in forward order on purpose because popular contracts may have
	// dozens of states changes due to ETH deposits/withdraw after contract creation,
	// so it is optimal to search from the beginning even if the contract has multiple
	// incarnations.
	accHistory, err := tx.Cursor(kv.AccountsHistory)
	if err != nil {
		return nil, err
	}
	defer accHistory.Close()

	accCS, err := tx.CursorDupSort(kv.AccountChangeSet)
	if err != nil {
		return nil, err
	}
	defer accCS.Close()

	// Locate shard that contains the block where incarnation changed
	acs := historyv2.Mapper[kv.AccountChangeSet]
	k, v, err := accHistory.Seek(acs.IndexChunkKey(addr.Bytes(), 0))
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(k, addr.Bytes()) {
		log.Error("Couldn't find any shard for account history", "addr", addr)
		return nil, fmt.Errorf("could't find any shard for account history addr=%v", addr)
	}

	bm := bitmapdb.NewBitmap64()
	defer bitmapdb.ReturnToPool64(bm)
	prevShardMaxBl := uint64(0)
	for {
		_, err := bm.ReadFrom(bytes.NewReader(v))
		if err != nil {
			return nil, err
		}

		// Shortcut precheck
		st, err := acs.Find(accCS, bm.Maximum(), addr.Bytes())
		if err != nil {
			return nil, err
		}
		if st == nil {
			log.Error("Unexpected error, couldn't find changeset", "block", bm.Maximum(), "addr", addr)
			return nil, fmt.Errorf("unexpected error, couldn't find changeset block=%v addr=%v", bm.Maximum(), addr)
		}

		// Found the shard where the incarnation change happens; ignore all
		// next shards
		if err := acc.DecodeForStorage(st); err != nil {
			return nil, err
		}
		if acc.Incarnation >= plainStateAcc.Incarnation {
			break
		}
		prevShardMaxBl = bm.Maximum()

		k, v, err = accHistory.Next()
		if err != nil {
			return nil, err
		}

		// No more shards; it means the max bl from previous shard
		// contains the incarnation change
		if !bytes.HasPrefix(k, addr.Bytes()) {
			break
		}
	}

	// Binary search block number inside shard; get first block where desired
	// incarnation appears
	blocks := bm.ToArray()
	var searchErr error
	r := sort.Search(len(blocks), func(i int) bool {
		bl := blocks[i]
		st, err := acs.Find(accCS, bl, addr.Bytes())
		if err != nil {
			searchErr = err
			return false
		}
		if st == nil {
			log.Error("Unexpected error, couldn't find changeset", "block", bl, "addr", addr)
			return false
		}

		if err := acc.DecodeForStorage(st); err != nil {
			searchErr = err
			return false
		}
		if acc.Incarnation < plainStateAcc.Incarnation {
			return false
		}
		return true
	})

	if searchErr != nil {
		return nil, searchErr
	}

	// The sort.Search function finds the first block where the incarnation has
	// changed to the desired one, so we get the previous block from the bitmap;
	// however if the found block is already the first one from the bitmap, it means
	// the block we want is the max block from the previous shard.
	blockFound := prevShardMaxBl
	if r > 0 {
		blockFound = blocks[r-1]
	}
	// Trace block, find tx and contract creator
	tracer := NewCreateTracer(ctx, addr)
	if err := api.genericTracer(tx, ctx, blockFound, 0, 0, chainConfig, tracer); err != nil {
		return nil, err
	}

	return &ContractCreatorData{
		Tx:      tracer.Tx.Hash(),
		Creator: tracer.Creator,
	}, nil
}
