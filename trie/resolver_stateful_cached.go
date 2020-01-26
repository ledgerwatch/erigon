package trie

import (
	"bytes"
	"fmt"
	"math/big"
	"runtime/debug"
	"strings"

	"github.com/ledgerwatch/bolt"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/hexutil"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

type ResolverStatefulCached struct {
	*ResolverStateful
}

func NewResolverStatefulCached(topLevels int, requests []*ResolveRequest, hookFunction hookFunction) *ResolverStatefulCached {
	return &ResolverStatefulCached{
		NewResolverStateful(topLevels, requests, hookFunction),
	}
}

func hexIncrement(in []byte) ([]byte, error) {
	digit, err := hexutil.DecodeBig(string(in))
	if err != nil {
		return nil, err
	}
	digit.Add(digit, big.NewInt(1))
	out := hexutil.EncodeBig(digit)
	if len(out) != len(in) {
		return nil, nil
	}
	return []byte(out), nil
}

// keyIsBefore - kind of bytes.Compare, but nil is the last key. And return
func keyIsBefore(k1, k2 []byte) (bool, []byte) {
	if k1 == nil {
		return false, k2
	}

	if k2 == nil {
		return true, k1
	}

	switch bytes.Compare(k1, k2) {
	case -1, 0:
		return true, k1
	default:
		return false, k2
	}
}

func (tr *ResolverStatefulCached) RebuildTrie(
	db ethdb.Database,
	blockNr uint64,
	accounts bool,
	historical bool) error {
	startkeys, fixedbits := tr.PrepareResolveParams()
	if db == nil {
		var b strings.Builder
		fmt.Fprintf(&b, "ResolveWithDb(db=nil), accounts: %t\n", accounts)
		for i, sk := range startkeys {
			fmt.Fprintf(&b, "sk %x, bits: %d\n", sk, fixedbits[i])
		}
		return fmt.Errorf("unexpected resolution: %s at %s", b.String(), debug.Stack())
	}

	typed, ok := db.(*ethdb.BoltDatabase)
	if !ok {
		panic("only Bolt supported yet")
	}
	boltDb := typed.GetDb()

	var err error
	if accounts {
		if historical {
			panic("historical not supported yet")
			//err = db.MultiWalkAsOf(dbutils.AccountsBucket, dbutils.AccountsHistoryBucket, startkeys, fixedbits, blockNr+1, tr.WalkerAccounts)
		} else {
			//err = db.MultiWalk(dbutils.AccountsBucket, startkeys, fixedbits, tr.WalkerAccounts)
			err = tr.MultiWalk2(boltDb, dbutils.AccountsBucket, startkeys, fixedbits, tr.WalkerAccounts)
		}
	} else {
		if historical {
			panic("historical not supported yet")
			//err = db.MultiWalkAsOf(dbutils.StorageBucket, dbutils.StorageHistoryBucket, startkeys, fixedbits, blockNr+1, tr.WalkerStorage)
		} else {
			//err = db.MultiWalk(dbutils.StorageBucket, startkeys, fixedbits, tr.WalkerStorage)
			err = tr.MultiWalk2(boltDb, dbutils.AccountsBucket, startkeys, fixedbits, tr.WalkerAccounts)
		}
	}
	if err != nil {
		return err
	}
	return tr.finaliseRoot()
}

//func (tr *ResolverStatefulCached) Walker2(isAccount bool, keyIdx int, k []byte, v []byte, useCache bool) error {
//	// Algo here. Make hashgen
//	return tr.Walker(isAccount bool, keyIdx, k, v)
//}

func (tr *ResolverStatefulCached) WalkerAccounts(keyIdx int, k []byte, v []byte, useCache bool) error {
	return tr.Walker(true, keyIdx, k, v)
}

func (tr *ResolverStatefulCached) WalkerStorage(keyIdx int, k []byte, v []byte, useCache bool) error {
	return tr.Walker(false, keyIdx, k, v)
}

func (tr *ResolverStatefulCached) MultiWalk2(db *bolt.DB, bucket []byte, startkeys [][]byte, fixedbits []uint, walker func(keyIdx int, k []byte, v []byte, useCache bool) error) error {
	if len(startkeys) == 0 {
		return nil
	}
	rangeIdx := 0 // What is the current range we are extracting
	fixedbytes, mask := ethdb.Bytesmask(fixedbits[rangeIdx])
	startkey := startkeys[rangeIdx]
	err := db.View(func(tx *bolt.Tx) error {
		cacheBucket := tx.Bucket(dbutils.IntermediateTrieHashesBucket)
		if cacheBucket == nil {
			return nil
		}
		cache := cacheBucket.Cursor()

		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		c := b.Cursor()

		k, v := c.Seek(startkey)
		cacheK, cacheV := cache.Seek(startkey)
		_ = cacheK
		_ = cacheV

		fmt.Printf("Walk Start: \n\t%x, %x, \n\t%x, %x\n", k, v, cacheK, cacheV)

		for k != nil || cacheK != nil {
			useCache, minKey := keyIsBefore(cacheK, k)
			_ = minKey
			// Adjust rangeIdx if needed
			cmp := int(-1)
			for fixedbytes > 0 && cmp != 0 {
				useCache, minKey = keyIsBefore(cacheK, k)

				cmp = bytes.Compare(minKey[:fixedbytes-1], startkey[:fixedbytes-1])
				switch cmp {
				case 0:
					k1 := minKey[fixedbytes-1] & mask
					k2 := startkey[fixedbytes-1] & mask
					if k1 < k2 {
						cmp = -1
					} else if k1 > k2 {
						cmp = 1
					}
				case -1:
					k, v = c.SeekTo(startkey)
					cacheK, cacheV = cache.SeekTo(startkey)
					if k == nil {
						return nil
					}
				default:
					rangeIdx++
					if rangeIdx == len(startkeys) {
						return nil
					}
					fixedbytes, mask = ethdb.Bytesmask(fixedbits[rangeIdx])
					startkey = startkeys[rangeIdx]
				}
			}

			if useCache {
				if len(cacheV) > 0 {
					if err := walker(rangeIdx, cacheK, cacheV, useCache); err != nil {
						return err
					}
				}

				next, err := hexIncrement(cacheK)
				if err != nil {
					return err
				}
				k, v = c.SeekTo(next)
				cacheK, cacheV = cache.SeekTo(next)
				if k == nil {
					return nil
				}
				continue
			}

			if len(v) > 0 {
				if err := walker(rangeIdx, k, v, useCache); err != nil {
					return err
				}
			}
			k, v = c.Next()
		}
		return nil
	})
	return err
}
