package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon-lib/aggregator"
	"github.com/ledgerwatch/erigon-lib/kv"
	kv2 "github.com/ledgerwatch/erigon-lib/kv/mdbx"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/core"
	"github.com/ledgerwatch/erigon/core/rawdb"
	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/core/types/accounts"
	"github.com/ledgerwatch/erigon/core/vm"
	"github.com/ledgerwatch/erigon/params"
	"github.com/ledgerwatch/log/v3"
	"github.com/spf13/cobra"
)

func init() {
	withBlock(erigon2Cmd)
	withDatadir(erigon2Cmd)
	rootCmd.AddCommand(erigon2Cmd)
}

var erigon2Cmd = &cobra.Command{
	Use:   "erigon2",
	Short: "Exerimental command to re-execute blocks from beginning using erigon2 state representation",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New()
		return Erigon2(genesis, logger, block, datadir)
	},
}

func Erigon2(genesis *core.Genesis, logger log.Logger, blockNum uint64, datadir string) error {
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		interruptCh <- true
	}()
	historyDb, err := kv2.NewMDBX(logger).Path(path.Join(datadir, "chaindata")).Open()
	if err != nil {
		return err
	}
	defer historyDb.Close()
	ctx := context.Background()
	historyTx, err1 := historyDb.BeginRo(ctx)
	if err1 != nil {
		return err1
	}
	defer historyTx.Rollback()
	stateDbPath := path.Join(datadir, "statedb")
	if _, err = os.Stat(stateDbPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else if err = os.RemoveAll(stateDbPath); err != nil {
		return err
	}
	db, err2 := kv2.NewMDBX(logger).Path(stateDbPath).Open()
	if err2 != nil {
		return err2
	}
	defer db.Close()
	aggPath := path.Join(datadir, "aggregator")
	if _, err = os.Stat(aggPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else if err = os.RemoveAll(aggPath); err != nil {
		return err
	}
	if err = os.Mkdir(aggPath, os.ModePerm); err != nil {
		return err
	}
	agg, err3 := aggregator.NewAggregator(aggPath, 90000, 4096)
	if err3 != nil {
		return fmt.Errorf("create aggregator: %w", err3)
	}
	chainConfig := genesis.Config
	vmConfig := vm.Config{}

	noOpWriter := state.NewNoopWriter()
	interrupt := false
	block := uint64(0)
	var rwTx kv.RwTx
	defer func() {
		rwTx.Rollback()
	}()
	var tx kv.Tx
	defer func() {
		tx.Rollback()
	}()
	if rwTx, err = db.BeginRw(ctx); err != nil {
		return err
	}
	_, genesisIbs, err4 := genesis.ToBlock()
	if err4 != nil {
		return err4
	}
	genesisW, err5 := agg.MakeStateWriter(rwTx, 0)
	if err5 != nil {
		return err5
	}
	if err = genesisIbs.CommitBlock(params.Rules{}, &WriterWrapper{w: genesisW}); err != nil {
		return fmt.Errorf("cannot write state: %w", err)
	}
	if err = rwTx.Commit(); err != nil {
		return err
	}
	for !interrupt {
		block++
		if block >= blockNum {
			break
		}
		blockHash, err := rawdb.ReadCanonicalHash(historyTx, block)
		if err != nil {
			return err
		}
		var b *types.Block
		b, _, err = rawdb.ReadBlockWithSenders(historyTx, blockHash, block)
		if err != nil {
			return err
		}
		if b == nil {
			break
		}
		if tx, err = db.BeginRo(ctx); err != nil {
			return err
		}
		if rwTx, err = db.BeginRw(ctx); err != nil {
			return err
		}
		r := agg.MakeStateReader(tx, block)
		var w *aggregator.Writer
		if w, err = agg.MakeStateWriter(rwTx, block); err != nil {
			return err
		}
		intraBlockState := state.New(&ReaderWrapper{r: r})
		getHeader := func(hash common.Hash, number uint64) *types.Header { return rawdb.ReadHeader(historyTx, hash, number) }
		if _, err = runBlock(intraBlockState, noOpWriter, &WriterWrapper{w: w}, chainConfig, getHeader, nil, b, vmConfig); err != nil {
			return fmt.Errorf("block %d: %w", block, err)
		}
		if block%1000 == 0 {
			log.Info("Processed", "blocks", block)
		}
		// Check for interrupts
		select {
		case interrupt = <-interruptCh:
			fmt.Println("interrupted, please wait for cleanup...")
		default:
		}
		tx.Rollback()
		if err = w.Finish(); err != nil {
			return err
		}
		if err = rwTx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// Implements StateReader and StateWriter
type ReaderWrapper struct {
	r *aggregator.Reader
}

type WriterWrapper struct {
	w *aggregator.Writer
}

func (rw *ReaderWrapper) ReadAccountData(address common.Address) (*accounts.Account, error) {
	enc, err := rw.r.ReadAccountData(address.Bytes())
	if err != nil {
		return nil, err
	}
	if len(enc) == 0 {
		return nil, nil
	}
	var a accounts.Account
	if err = a.DecodeForStorage(enc); err != nil {
		return nil, err
	}
	return &a, nil
}

func (rw *ReaderWrapper) ReadAccountStorage(address common.Address, incarnation uint64, key *common.Hash) ([]byte, error) {
	enc, err := rw.r.ReadAccountStorage(address.Bytes(), key.Bytes())
	if err != nil {
		return nil, err
	}
	if enc == nil {
		return nil, nil
	}
	return enc.Bytes(), nil
}

func (rw *ReaderWrapper) ReadAccountCode(address common.Address, incarnation uint64, codeHash common.Hash) ([]byte, error) {
	enc, err := rw.r.ReadAccountCode(address.Bytes())
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (rw *ReaderWrapper) ReadAccountCodeSize(address common.Address, incarnation uint64, codeHash common.Hash) (int, error) {
	enc, err := rw.r.ReadAccountCode(address.Bytes())
	if err != nil {
		return 0, err
	}
	return len(enc), nil
}

func (rw *ReaderWrapper) ReadAccountIncarnation(address common.Address) (uint64, error) {
	return 0, nil
}

func (ww *WriterWrapper) UpdateAccountData(address common.Address, original, account *accounts.Account) error {
	value := make([]byte, account.EncodingLengthForStorage())
	account.EncodeForStorage(value)
	if err := ww.w.UpdateAccountData(address.Bytes(), value); err != nil {
		return err
	}
	return nil
}

func (ww *WriterWrapper) UpdateAccountCode(address common.Address, incarnation uint64, codeHash common.Hash, code []byte) error {
	if err := ww.w.UpdateAccountCode(address.Bytes(), code); err != nil {
		return err
	}
	return nil
}

func (ww *WriterWrapper) DeleteAccount(address common.Address, original *accounts.Account) error {
	if original == nil || !original.Initialised {
		return nil
	}
	if err := ww.w.DeleteAccount(address.Bytes()); err != nil {
		return err
	}
	return nil
}

func (ww *WriterWrapper) WriteAccountStorage(address common.Address, incarnation uint64, key *common.Hash, original, value *uint256.Int) error {
	if err := ww.w.WriteAccountStorage(address.Bytes(), key.Bytes(), value); err != nil {
		return err
	}
	return nil
}

func (ww *WriterWrapper) CreateContract(address common.Address) error {
	return nil
}
