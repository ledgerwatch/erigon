// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
//nolint:errcheck,prealloc
package core

import (
	"fmt"
	"math/big"
	"testing"

	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/datadir"
	"github.com/ledgerwatch/log/v3"
	"golang.org/x/crypto/sha3"

	"github.com/ledgerwatch/erigon/common/u256"
	"github.com/ledgerwatch/erigon/consensus/ethash"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/crypto"
	"github.com/ledgerwatch/erigon/params"
	"github.com/ledgerwatch/erigon/rlp"
)

func getBlock(tb testing.TB, transactions int, uncles int, dataSize int, tmpDir string, logger log.Logger) *types.Block {
	_, db, _ := temporaltest.NewTestDB(tb, datadir.New(tmpDir))
	var (
		aa = libcommon.HexToAddress("0x000000000000000000000000000000000000aaaa")
		// Generate a canonical chain to act as the main dataset
		engine = ethash.NewFaker()
		// A sender who makes transactions, has some funds
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &types.Genesis{
			Config: params.TestChainConfig,
			Alloc:  types.GenesisAlloc{address: {Balance: funds}},
		}
		genesis = MustCommitGenesis(gspec, db, tmpDir, logger)
	)

	// We need to generate as many blocks +1 as uncles
	chain, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, uncles+1, func(n int, b *BlockGen) {
		if n == uncles {
			// Add transactions and stuff on the last block
			for i := 0; i < transactions; i++ {
				tx, _ := types.SignTx(types.NewTransaction(uint64(i), aa,
					u256.Num0, 50000, u256.Num1, make([]byte, dataSize)), *types.LatestSignerForChainID(nil), key)
				b.AddTx(tx)
			}
			for i := 0; i < uncles; i++ {
				b.AddUncle(&types.Header{ParentHash: b.PrevBlock(n - 1 - i).Hash(), Number: big.NewInt(int64(n - i))})
			}
		}
	})
	block := chain.TopBlock
	return block
}

// TestRlpIterator tests that individual transactions can be picked out
// from blocks without full unmarshalling/marshalling
func TestRlpIterator(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		txs      int
		uncles   int
		datasize int
	}{
		{0, 0, 0},
		{0, 2, 0},
		{10, 0, 0},
		{10, 2, 0},
		{10, 2, 50},
	} {
		testRlpIterator(t, tt.txs, tt.uncles, tt.datasize)
	}
}

func testRlpIterator(t *testing.T, txs, uncles, datasize int) {
	desc := fmt.Sprintf("%d txs [%d datasize] and %d uncles", txs, datasize, uncles)
	bodyRlp, _ := rlp.EncodeToBytes(getBlock(t, txs, uncles, datasize, "", log.Root()).Body())
	it, err := rlp.NewListIterator(bodyRlp)
	if err != nil {
		t.Fatal(err)
	}
	// Check that txs exist
	if !it.Next() {
		t.Fatal("expected two elems, got zero")
	}
	txdata := it.Value()
	// Check that uncles exist
	if !it.Next() {
		t.Fatal("expected three elems, got two")
	}
	// No more after that
	if it.Next() {
		t.Fatal("expected only three elems, got more")
	}
	txIt, err := rlp.NewListIterator(txdata)
	if err != nil {
		t.Fatal(err)
	}
	var gotHashes []libcommon.Hash
	var expHashes []libcommon.Hash
	for txIt.Next() {
		gotHashes = append(gotHashes, crypto.Keccak256Hash(txIt.Value()))
	}

	var expBody types.Body
	err = rlp.DecodeBytes(bodyRlp, &expBody)
	if err != nil {
		t.Fatal(err)
	}
	for _, tx := range expBody.Transactions {
		expHashes = append(expHashes, tx.Hash())
	}
	if gotLen, expLen := len(gotHashes), len(expHashes); gotLen != expLen {
		t.Fatalf("testcase %v: length wrong, got %d exp %d", desc, gotLen, expLen)
	}
	// also sanity check against input
	if gotLen := len(gotHashes); gotLen != txs {
		t.Fatalf("testcase %v: length wrong, got %d exp %d", desc, gotLen, txs)
	}
	for i, got := range gotHashes {
		if exp := expHashes[i]; got != exp {
			t.Errorf("testcase %v: hash wrong, got %x, exp %x", desc, got, exp)
		}
	}
}

// BenchmarkHashing compares the speeds of hashing a rlp raw data directly
// without the unmarshalling/marshalling step
func BenchmarkHashing(b *testing.B) {
	// Make a pretty fat block
	var (
		bodyRlp  []byte
		blockRlp []byte
	)
	{
		block := getBlock(b, 200, 2, 50, "", log.Root())
		bodyRlp, _ = rlp.EncodeToBytes(block.Body())
		blockRlp, _ = rlp.EncodeToBytes(block)
	}
	var got libcommon.Hash
	var hasher = sha3.NewLegacyKeccak256()
	b.Run("iteratorhashing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var hash libcommon.Hash
			it, err := rlp.NewListIterator(bodyRlp)
			if err != nil {
				b.Fatal(err)
			}
			it.Next()
			txs := it.Value()
			txIt, err := rlp.NewListIterator(txs)
			if err != nil {
				b.Fatal(err)
			}
			for txIt.Next() {
				hasher.Reset()
				hasher.Write(txIt.Value())
				hasher.Sum(hash[:0])
				got = hash
			}
		}
	})
	var exp libcommon.Hash
	b.Run("fullbodyhashing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var body types.Body
			rlp.DecodeBytes(bodyRlp, &body)
			for _, tx := range body.Transactions {
				exp = tx.Hash()
			}
		}
	})
	b.Run("fullblockhashing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var block types.Block
			rlp.DecodeBytes(blockRlp, &block)
			for _, tx := range block.Transactions() {
				tx.Hash()
			}
		}
	})
	if got != exp {
		b.Fatalf("hash wrong, got %x exp %x", got, exp)
	}
}
