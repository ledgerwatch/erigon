package stagedsync

import (
	"context"
	"fmt"

	"github.com/ledgerwatch/turbo-geth/eth/stagedsync/stages"
	"github.com/ledgerwatch/turbo-geth/turbo/snapshotsync"

	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
)

type FinishCfg struct {
	db     ethdb.RwKV
	tmpDir string
}

func StageFinishCfg(db ethdb.RwKV, tmpDir string) FinishCfg {
	return FinishCfg{
		db:     db,
		tmpDir: tmpDir,
	}
}

func FinishForward(s *StageState, tx ethdb.RwTx, cfg FinishCfg, btClient *snapshotsync.Client, snBuilder *snapshotsync.SnapshotMigrator) error {
	useExternalTx := tx != nil
	if !useExternalTx {
		var err error
		tx, err = cfg.db.BeginRw(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}

	var executionAt uint64
	var err error
	if executionAt, err = s.ExecutionAt(tx); err != nil {
		return err
	}
	if executionAt <= s.BlockNumber {
		s.Done()
		return nil
	}

	err = MigrateSnapshot(executionAt, tx, cfg.db, btClient, snBuilder)
	if err != nil {
		return err
	}
	err = s.DoneAndUpdate(tx, executionAt)
	if err != nil {
		return err
	}
	if !useExternalTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func UnwindFinish(u *UnwindState, s *StageState, tx ethdb.RwTx, cfg FinishCfg) error {
	useExternalTx := tx != nil
	if !useExternalTx {
		var err error
		tx, err = cfg.db.BeginRw(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}
	err := u.Done(tx)
	if err != nil {
		return err
	}
	if !useExternalTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func NotifyNewHeaders2(finishStageBeforeSync, unwindTo uint64, notifier ChainEventNotifier, db ethdb.Database) error {
	finishAt, err := stages.GetStageProgress(db, stages.Finish) // because later stages can be disabled
	if err != nil {
		return err
	}
	notifyFrom := finishStageBeforeSync + 1
	if unwindTo < finishStageBeforeSync {
		notifyFrom = unwindTo + 1
	}
	return NotifyNewHeaders(notifyFrom, finishAt, notifier, db)
}
func NotifyNewHeaders(from, to uint64, notifier ChainEventNotifier, db ethdb.Database) error {
	if notifier == nil {
		log.Warn("rpc notifier is not set, rpc daemon won't be updated about headers")
		return nil
	}
	for i := from; i <= to; i++ {
		header := rawdb.ReadHeaderByNumber(db, i)
		if header == nil {
			return fmt.Errorf("could not find canonical header for number: %d", i)
		}
		notifier.OnNewHeader(header)
	}

	return nil
}

func MigrateSnapshot(to uint64, tx ethdb.RwTx, db ethdb.RwKV, btClient *snapshotsync.Client, mg *snapshotsync.SnapshotMigrator) error {
	if mg == nil {
		return nil
	}

	snBlock := snapshotsync.CalculateEpoch(to, snapshotsync.EpochSize)
	return mg.Migrate(db, tx, snBlock, btClient)
}
