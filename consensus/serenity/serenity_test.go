package serenity

import (
	"math/big"
	"testing"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/params"
)

type readerMock struct{}

func (r readerMock) Config() *params.ChainConfig {
	return nil
}

func (r readerMock) CurrentHeader() *types.Header {
	return nil
}

func (r readerMock) GetHeader(common.Hash, uint64) *types.Header {
	return nil
}

func (r readerMock) GetHeaderByNumber(uint64) *types.Header {
	return nil
}

func (r readerMock) GetHeaderByHash(common.Hash) *types.Header {
	return nil
}

// The thing only that changes beetwen normal ethash checks other than POW, is difficulty
// and nonce so we are gonna test those
func TestVerifyHeaderDifficulty(t *testing.T) {
	header := &types.Header{
		Difficulty: big.NewInt(1),
	}

	parent := &types.Header{}

	serenity := New()

	err := serenity.verifyHeader(readerMock{}, header, parent)
	if err != errInvalidDifficulty {
		if err != nil {
			t.Fatalf("Serenity should not accept non-zero difficulty, got %s", err.Error())
		} else {
			t.Fatalf("Serenity should not accept non-zero difficulty")
		}
	}
}

func TestVerifyHeaderNonce(t *testing.T) {
	header := &types.Header{
		Nonce:      types.BlockNonce{1, 0, 0, 0, 0, 0, 0, 0},
		Difficulty: big.NewInt(0),
	}

	parent := &types.Header{}

	serenity := New()

	err := serenity.verifyHeader(readerMock{}, header, parent)
	if err != errInvalidNonce {
		if err != nil {
			t.Fatalf("Serenity should not accept non-zero difficulty, got %s", err.Error())
		} else {
			t.Fatalf("Serenity should not accept non-zero difficulty")
		}
	}
}
