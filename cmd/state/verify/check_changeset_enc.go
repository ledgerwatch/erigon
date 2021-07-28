package verify

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ledgerwatch/erigon/common/changeset"
	"github.com/ledgerwatch/erigon/ethdb/kv"
	"github.com/ledgerwatch/erigon/ethdb/olddb"
	"golang.org/x/sync/errgroup"
)

type Walker interface {
	Walk(f func(k, v []byte) error) error
	Find(k []byte) ([]byte, error)
}

func CheckEnc(chaindata string) error {
	db := olddb.MustOpen(chaindata)
	defer db.Close()
	var (
		currentSize uint64
		newSize     uint64
	)

	//set test methods
	chainDataStorageDecoder := changeset.Mapper[kv.StorageChangeSet].Decode
	testStorageEncoder := changeset.Mapper[kv.StorageChangeSet].Encode
	testStorageDecoder := changeset.Mapper[kv.StorageChangeSet].Decode

	startTime := time.Now()
	ch := make(chan struct {
		k []byte
		v []byte
	})
	stop := make(chan struct{})
	//run workers
	g, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < runtime.NumCPU()-1; i++ {
		g.Go(func() error {
			for {
				select {
				case v := <-ch:
					blockNum, kk, vv := chainDataStorageDecoder(v.k, v.v)
					cs := changeset.NewStorageChangeSet()
					_ = cs.Add(v.k, v.v)
					atomic.AddUint64(&currentSize, uint64(len(v.v)))
					innerErr := testStorageEncoder(blockNum, cs, func(k, v []byte) error {
						atomic.AddUint64(&newSize, uint64(len(v)))
						_, a, b := testStorageDecoder(k, v)
						if !bytes.Equal(kk, a) {
							return fmt.Errorf("incorrect order. block: %d", blockNum)
						}
						if !bytes.Equal(vv, b) {
							return fmt.Errorf("incorrect value. block: %d, key: %x", blockNum, a)
						}
						return nil
					})
					if innerErr != nil {
						return innerErr
					}

				case <-ctx.Done():
					return nil
				case <-stop:
					return nil
				}
			}
		})
	}

	g.Go(func() error {
		var i uint64
		defer func() {
			close(stop)
		}()
		return db.View(context.Background(), func(tx kv.Tx) error {
			return tx.ForEach(kv.StorageChangeSet, []byte{}, func(k, v []byte) error {
				if i%100_000 == 0 {
					blockNum := binary.BigEndian.Uint64(k)
					fmt.Printf("Processed %dK, block number %d, current %d, new %d, time %s\n",
						i/1000,
						blockNum,
						atomic.LoadUint64(&currentSize),
						atomic.LoadUint64(&newSize),
						time.Since(startTime))
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:

				}

				i++
				ch <- struct {
					k []byte
					v []byte
				}{k: k, v: v}
				return nil
			})
		})
	})
	err := g.Wait()
	if err != nil {
		return err
	}

	fmt.Println("-- Final size --")
	fmt.Println("Current:", currentSize)
	fmt.Println("New:", newSize)

	return nil
}
