package commands

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/common/hexutil"
	"github.com/ledgerwatch/erigon/core/state"
)

// ParityAPI the interface for the parity_ RPC commands
type ParityAPI interface {
	ListStorageKeys(ctx context.Context, account common.Address, quantity int, offset *hexutil.Bytes) ([]hexutil.Bytes, error)
}

// ParityAPIImpl data structure to store things needed for parity_ commands
type ParityAPIImpl struct {
	db kv.RoDB
}

// NewParityAPIImpl returns ParityAPIImpl instance
func NewParityAPIImpl(db kv.RoDB) *ParityAPIImpl {
	return &ParityAPIImpl{
		db: db,
	}
}

// ListStorageKeys implements parity_listStorageKeys. Returns all storage keys of the given address
func (api *ParityAPIImpl) ListStorageKeys(ctx context.Context, account common.Address, quantity int, offset *hexutil.Bytes) ([]hexutil.Bytes, error) {
	tx, txErr := api.db.BeginRo(ctx)
	if txErr != nil {
		return nil, fmt.Errorf("listStorageKeys cannot open tx: %w", txErr)
	}
	defer tx.Rollback()
	a, err := state.NewPlainStateReader(tx).ReadAccountData(account)
	if err != nil {
		return nil, err
	} else if a == nil {
		return nil, fmt.Errorf("acc not found")
	}
	c, err := tx.Cursor(kv.PlainState)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	keys := make([]hexutil.Bytes, 0)
	var (
		k []byte
	)

	if offset != nil {
		k, _, err = c.SeekExact(*offset)
	} else {
		k, _, err = c.Seek(account.Bytes())
	}
	if err != nil {
		return nil, err
	}
	for ; k != nil && err == nil && len(keys) != quantity; k, _, err = c.Next() {
		if err != nil {
			return nil, err
		}
		if !bytes.HasPrefix(k, account.Bytes()) {
			break
		}
		if len(k) <= common.AddressLength {
			continue
		}
		keys = append(keys, k[common.AddressLength:])
	}
	return keys, nil
}
