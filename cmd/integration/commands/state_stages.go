package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/cmd/utils"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/changeset"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"github.com/ledgerwatch/turbo-geth/eth/stagedsync"
	"github.com/ledgerwatch/turbo-geth/eth/stagedsync/stages"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/spf13/cobra"
)

var stateStags = &cobra.Command{
	Use: "state_stages",
	Short: `Move all StateStages (which happen after senders) forward. 
			Stops at StageSenders progress or at "--block".
			Each iteration test will move forward "--unwind_every" blocks, then unwind "--unwind" blocks.
			Use reset_state command to re-run this test.
			When finish all cycles, does comparison to "--reference_chaindata" if flag provided.
		`,
	Example: "go run ./cmd/integration state_stages --chaindata=... --verbosity=3 --unwind=100 --unwind_every=100000 --block=2000000",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := utils.RootContext()
		if err := syncBySmallSteps(ctx, chaindata); err != nil {
			log.Error("Error", "err", err)
			return err
		}

		if referenceChaindata != "" {
			if err := compareStates(ctx, chaindata, referenceChaindata); err != nil {
				log.Error(err.Error())
				return err
			}

		}
		return nil
	},
}

func init() {
	withChaindata(stateStags)
	withReferenceChaindata(stateStags)
	withUnwind(stateStags)
	withUnwindEvery(stateStags)
	withBlock(stateStags)

	rootCmd.AddCommand(stateStags)
}

func syncBySmallSteps(ctx context.Context, chaindata string) error {
	core.UsePlainStateExecution = true
	db := ethdb.MustOpen(chaindata)
	defer db.Close()

	ch := ctx.Done()

	expectedAccountChanges := make(map[uint64][]byte)
	expectedStorageChanges := make(map[uint64][]byte)
	changeSetHook := func(blockNum uint64, csw *state.ChangeSetWriter) {
		accountChanges, err := csw.GetAccountChanges()
		if err != nil {
			panic(err)
		}
		expectedAccountChanges[blockNum], err = changeset.EncodeAccountsPlain(accountChanges)
		if err != nil {
			panic(err)
		}

		storageChanges, err := csw.GetStorageChanges()
		if err != nil {
			panic(err)
		}
		if storageChanges.Len() > 0 {
			expectedStorageChanges[blockNum], err = changeset.EncodeStoragePlain(storageChanges)
			if err != nil {
				panic(err)
			}
		}
	}

	var tx ethdb.DbWithPendingMutations = ethdb.NewTxDbWithoutTransaction(db)
	defer tx.Rollback()

	bc, st, progress := newSync(ch, tx, changeSetHook)
	defer bc.Stop()

	st.BeforeStageRun(stages.Execution, func() error {
		if hasTx, ok := tx.(ethdb.HasTx); ok && hasTx.Tx() != nil {
			return nil
		}
		log.Debug("cycle: begin transaction")
		var errTx error
		tx, errTx = tx.Begin()
		return errTx
	})
	st.BeforeStageRun(stages.TxPool, func() error {
		log.Debug("cycle: commit transaction")
		var errTx error
		_, errTx = tx.Commit()
		return errTx
	})
	st.BeforeUnwind(func() error {
		if hasTx, ok := tx.(ethdb.HasTx); ok && hasTx.Tx() != nil {
			return nil
		}
		log.Debug("cycle unwind: begin transaction")
		var errTx error
		tx, errTx = tx.Begin()
		return errTx
	})
	st.AfterUnwind(func() error {
		log.Debug("cycle unwind: commit transaction")
		_, errCommit := tx.Commit()
		return errCommit
	})

	st.DisableStages(stages.Headers, stages.BlockHashes, stages.Bodies, stages.Senders)

	senderStageProgress := progress(stages.Senders).BlockNumber

	var stopAt = senderStageProgress
	if block > 0 && block < senderStageProgress {
		stopAt = block
	}

	for progress(stages.Execution).BlockNumber < stopAt {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// All stages forward to `execStage + unwindEvery` block
		execToBlock := progress(stages.Execution).BlockNumber + unwindEvery
		if execToBlock > stopAt {
			execToBlock = stopAt + 1
			unwind = 0
		}

		// set block limit of execute stage
		st.MockExecFunc(stages.Execution, func(stageState *stagedsync.StageState, unwinder stagedsync.Unwinder) error {
			if err := stagedsync.SpawnExecuteBlocksStage(stageState, tx, bc.Config(), bc, bc.GetVMConfig(), execToBlock, ch, false, changeSetHook); err != nil {
				return fmt.Errorf("spawnExecuteBlocksStage: %w", err)
			}
			return nil
		})

		if err := st.Run(db, tx); err != nil {
			return err
		}

		for blockN := range expectedAccountChanges {
			if err := checkChangeSet(db, blockN, expectedAccountChanges[blockN], expectedStorageChanges[blockN]); err != nil {
				return err
			}
			delete(expectedAccountChanges, blockN)
			delete(expectedStorageChanges, blockN)
		}

		// Unwind all stages to `execStage - unwind` block
		if unwind == 0 {
			continue
		}

		execStage := progress(stages.Execution)
		to := execStage.BlockNumber - unwind

		if err := st.UnwindTo(to, tx); err != nil {
			return err
		}
	}

	return nil
}

func checkChangeSet(db *ethdb.ObjectDatabase, blockNum uint64, expectedAccountChanges []byte, expectedStorageChanges []byte) error {
	dbAccountChanges, err := db.GetChangeSetByBlock(false /* storage */, blockNum)
	if err != nil {
		return err
	}

	if !bytes.Equal(dbAccountChanges, expectedAccountChanges) {
		fmt.Printf("Unexpected account changes in block %d\nIn the database: ======================\n", blockNum)
		if err = changeset.AccountChangeSetPlainBytes(dbAccountChanges).Walk(func(k, v []byte) error {
			fmt.Printf("0x%x: %x\n", k, v)
			return nil
		}); err != nil {
			return err
		}
		fmt.Printf("Expected: ==========================\n")
		if err = changeset.AccountChangeSetPlainBytes(expectedAccountChanges).Walk(func(k, v []byte) error {
			fmt.Printf("0x%x %x\n", k, v)
			return nil
		}); err != nil {
			return err
		}
		return fmt.Errorf("check change set failed")
	}

	dbStorageChanges, err := db.GetChangeSetByBlock(true /* storage */, blockNum)
	if err != nil {
		return err
	}
	equal := true
	if !bytes.Equal(dbStorageChanges, expectedStorageChanges) {
		var addrs [][]byte
		var keys [][]byte
		var vals [][]byte
		if err = changeset.StorageChangeSetPlainBytes(dbStorageChanges).Walk(func(k, v []byte) error {
			addrs = append(addrs, common.CopyBytes(k[:common.AddressLength]))
			keys = append(keys, common.CopyBytes(k[common.AddressLength+common.IncarnationLength:]))
			vals = append(vals, common.CopyBytes(v))
			return nil
		}); err != nil {
			return err
		}
		i := 0
		if err = changeset.StorageChangeSetPlainBytes(expectedStorageChanges).Walk(func(k, v []byte) error {
			if !equal {
				return nil
			}
			if i >= len(addrs) {
				equal = false
				return nil
			}
			if !bytes.Equal(k[:common.AddressLength], addrs[i]) {
				equal = false
				return nil
			}
			if !bytes.Equal(k[common.AddressLength+common.IncarnationLength:], keys[i]) {
				equal = false
				return nil
			}
			if !bytes.Equal(v, vals[i]) {
				equal = false
				return nil
			}
			i++
			return nil
		}); err != nil {
			return err
		}
	}
	if !equal {
		fmt.Printf("Unexpected storage changes in block %d\nIn the database: ======================\n", blockNum)
		if err = changeset.StorageChangeSetPlainBytes(dbStorageChanges).Walk(func(k, v []byte) error {
			fmt.Printf("0x%x: [%x]\n", k, v)
			return nil
		}); err != nil {
			return err
		}
		fmt.Printf("Expected: ==========================\n")
		if err = changeset.StorageChangeSetPlainBytes(expectedStorageChanges).Walk(func(k, v []byte) error {
			fmt.Printf("0x%x: [%x]\n", k, v)
			return nil
		}); err != nil {
			return err
		}
		return fmt.Errorf("check change set failed")
	}
	return nil
}
