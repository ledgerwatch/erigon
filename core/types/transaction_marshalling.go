package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/common/hexutil"
	"github.com/protolambda/ztyp/view"
	"github.com/valyala/fastjson"
)

// txJSON is the JSON representation of transactions.
type txJSON struct {
	Type hexutil.Uint64 `json:"type"`

	// Common transaction fields:
	Nonce    *hexutil.Uint64 `json:"nonce"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	FeeCap   *hexutil.Big    `json:"maxFeePerGas"`
	Tip      *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Gas      *hexutil.Uint64 `json:"gas"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"input"`
	V        *hexutil.Big    `json:"v"`
	R        *hexutil.Big    `json:"r"`
	S        *hexutil.Big    `json:"s"`
	To       *common.Address `json:"to"`

	// Access list transaction fields:
	ChainID    *hexutil.Big `json:"chainId,omitempty"`
	AccessList *AccessList  `json:"accessList,omitempty"`

	// Blob transaction fields:
	MaxFeePerDataGas    *hexutil.Big  `json:"maxFeePerDataGas,omitempty"`
	BlobVersionedHashes []common.Hash `json:"blobVersionedHashes,omitempty"`
	Blobs               Blobs         `json:"blobs,omitempty"`
	BlobKzgs            BlobKzgs      `json:"blobKzgs,omitempty"`
	KzgAggregatedProof  KZGProof      `json:"kzgAggregatedProof,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

func (tx LegacyTx) MarshalJSON() ([]byte, error) {
	var enc txJSON
	// These are set for all tx types.
	enc.Hash = tx.Hash()
	enc.Type = hexutil.Uint64(tx.Type())
	enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
	enc.Gas = (*hexutil.Uint64)(&tx.Gas)
	enc.GasPrice = (*hexutil.Big)(tx.GasPrice.ToBig())
	enc.Value = (*hexutil.Big)(tx.Value.ToBig())
	enc.Data = (*hexutil.Bytes)(&tx.Data)
	enc.To = tx.To
	enc.V = (*hexutil.Big)(tx.V.ToBig())
	enc.R = (*hexutil.Big)(tx.R.ToBig())
	enc.S = (*hexutil.Big)(tx.S.ToBig())
	return json.Marshal(&enc)
}

func (tx AccessListTx) MarshalJSON() ([]byte, error) {
	var enc txJSON
	// These are set for all tx types.
	enc.Hash = tx.Hash()
	enc.Type = hexutil.Uint64(tx.Type())
	enc.ChainID = (*hexutil.Big)(tx.ChainID.ToBig())
	enc.AccessList = &tx.AccessList
	enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
	enc.Gas = (*hexutil.Uint64)(&tx.Gas)
	enc.GasPrice = (*hexutil.Big)(tx.GasPrice.ToBig())
	enc.Value = (*hexutil.Big)(tx.Value.ToBig())
	enc.Data = (*hexutil.Bytes)(&tx.Data)
	enc.To = tx.To
	enc.V = (*hexutil.Big)(tx.V.ToBig())
	enc.R = (*hexutil.Big)(tx.R.ToBig())
	enc.S = (*hexutil.Big)(tx.S.ToBig())
	return json.Marshal(&enc)
}

func (tx DynamicFeeTransaction) MarshalJSON() ([]byte, error) {
	var enc txJSON
	// These are set for all tx types.
	enc.Hash = tx.Hash()
	enc.Type = hexutil.Uint64(tx.Type())
	enc.ChainID = (*hexutil.Big)(tx.ChainID.ToBig())
	enc.AccessList = &tx.AccessList
	enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
	enc.Gas = (*hexutil.Uint64)(&tx.Gas)
	enc.FeeCap = (*hexutil.Big)(tx.FeeCap.ToBig())
	enc.Tip = (*hexutil.Big)(tx.Tip.ToBig())
	enc.Value = (*hexutil.Big)(tx.Value.ToBig())
	enc.Data = (*hexutil.Bytes)(&tx.Data)
	enc.To = tx.To
	enc.V = (*hexutil.Big)(tx.V.ToBig())
	enc.R = (*hexutil.Big)(tx.R.ToBig())
	enc.S = (*hexutil.Big)(tx.S.ToBig())
	return json.Marshal(&enc)
}

func (tx SignedBlobTx) MarshalJSON() ([]byte, error) {
	var enc txJSON
	enc.ChainID = (*hexutil.Big)(u256ToBig(&tx.Message.ChainID))
	enc.AccessList = (*AccessList)(&tx.Message.AccessList)
	enc.Nonce = (*hexutil.Uint64)(&tx.Message.Nonce)
	enc.Gas = (*hexutil.Uint64)(&tx.Message.Gas)
	enc.FeeCap = (*hexutil.Big)(u256ToBig(&tx.Message.GasFeeCap)) // MaxFeePerGas
	enc.Tip = (*hexutil.Big)(u256ToBig(&tx.Message.GasTipCap))    // MaxPriorityFeePerGas
	enc.Value = (*hexutil.Big)(u256ToBig(&tx.Message.Value))
	enc.Data = (*hexutil.Bytes)(&tx.Message.Data)
	enc.To = tx.GetTo()
	v, r, s := tx.RawSignatureValues()
	enc.V = (*hexutil.Big)(v.ToBig())
	enc.R = (*hexutil.Big)(r.ToBig())
	enc.S = (*hexutil.Big)(s.ToBig())
	enc.MaxFeePerDataGas = (*hexutil.Big)(u256ToBig(&tx.Message.MaxFeePerDataGas))
	enc.BlobVersionedHashes = tx.Message.BlobVersionedHashes

	// tx.WrapData is temp solution, this will require more detailed research + possible refactoring
	if tx.WrapData != nil {
		enc.Blobs = tx.WrapData.blobs()
		enc.BlobKzgs = tx.WrapData.kzgs()
		enc.KzgAggregatedProof = tx.WrapData.aggregatedProof()
	}
	return json.Marshal(&enc)
}

func UnmarshalTransactionFromJSON(input []byte) (Transaction, error) {
	var p fastjson.Parser
	v, err := p.ParseBytes(input)
	if err != nil {
		return nil, fmt.Errorf("parse transaction json: %w", err)
	}
	// check the type
	txTypeHex := v.GetStringBytes("type")
	var txType hexutil.Uint64 = LegacyTxType
	if txTypeHex != nil {
		if err = txType.UnmarshalText(txTypeHex); err != nil {
			return nil, err
		}
	}
	switch byte(txType) {
	case LegacyTxType:
		tx := &LegacyTx{}
		if err = tx.UnmarshalJSON(input); err != nil {
			return nil, err
		}
		return tx, nil
	case AccessListTxType:
		tx := &AccessListTx{}
		if err = tx.UnmarshalJSON(input); err != nil {
			return nil, err
		}
		return tx, nil
	case DynamicFeeTxType:
		tx := &DynamicFeeTransaction{}
		if err = tx.UnmarshalJSON(input); err != nil {
			return nil, err
		}
		return tx, nil
	case BlobTxType:
		tx := &SignedBlobTx{}
		if err = tx.UnmarshalJSON(input); err != nil {
			return nil, err
		}
		return tx, err
	default:
		return nil, fmt.Errorf("unknown transaction type: %v", txType)
	}
}

func (tx *LegacyTx) UnmarshalJSON(input []byte) error {
	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.To != nil {
		tx.To = dec.To
	}
	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' in transaction")
	}
	tx.Nonce = uint64(*dec.Nonce)
	if dec.GasPrice == nil {
		return errors.New("missing required field 'gasPrice' in transaction")
	}
	var overflow bool
	tx.GasPrice, overflow = uint256.FromBig(dec.GasPrice.ToInt())
	if overflow {
		return errors.New("'gasPrice' in transaction does not fit in 256 bits")
	}
	if dec.Gas == nil {
		return errors.New("missing required field 'gas' in transaction")
	}
	tx.Gas = uint64(*dec.Gas)
	if dec.Value == nil {
		return errors.New("missing required field 'value' in transaction")
	}
	tx.Value, overflow = uint256.FromBig(dec.Value.ToInt())
	if overflow {
		return errors.New("'value' in transaction does not fit in 256 bits")
	}
	if dec.Data == nil {
		return errors.New("missing required field 'input' in transaction")
	}
	tx.Data = *dec.Data
	if dec.V == nil {
		return errors.New("missing required field 'v' in transaction")
	}
	overflow = tx.V.SetFromBig(dec.V.ToInt())
	if overflow {
		return fmt.Errorf("dec.V higher than 2^256-1")
	}
	if dec.R == nil {
		return errors.New("missing required field 'r' in transaction")
	}
	overflow = tx.R.SetFromBig(dec.R.ToInt())
	if overflow {
		return fmt.Errorf("dec.R higher than 2^256-1")
	}
	if dec.S == nil {
		return errors.New("missing required field 's' in transaction")
	}
	overflow = tx.S.SetFromBig(dec.S.ToInt())
	if overflow {
		return fmt.Errorf("dec.S higher than 2^256-1")
	}
	if overflow {
		return errors.New("'s' in transaction does not fit in 256 bits")
	}
	withSignature := !tx.V.IsZero() || !tx.R.IsZero() || !tx.S.IsZero()
	if withSignature {
		if err := sanityCheckSignature(&tx.V, &tx.R, &tx.S, true); err != nil {
			return err
		}
	}
	return nil
}

func (tx *AccessListTx) UnmarshalJSON(input []byte) error {
	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	// Access list is optional for now.
	if dec.AccessList != nil {
		tx.AccessList = *dec.AccessList
	}
	if dec.ChainID == nil {
		return errors.New("missing required field 'chainId' in transaction")
	}
	var overflow bool
	tx.ChainID, overflow = uint256.FromBig(dec.ChainID.ToInt())
	if overflow {
		return errors.New("'chainId' in transaction does not fit in 256 bits")
	}
	if dec.To != nil {
		tx.To = dec.To
	}
	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' in transaction")
	}
	tx.Nonce = uint64(*dec.Nonce)
	if dec.GasPrice == nil {
		return errors.New("missing required field 'gasPrice' in transaction")
	}
	tx.GasPrice, overflow = uint256.FromBig(dec.GasPrice.ToInt())
	if overflow {
		return errors.New("'gasPrice' in transaction does not fit in 256 bits")
	}
	if dec.Gas == nil {
		return errors.New("missing required field 'gas' in transaction")
	}
	tx.Gas = uint64(*dec.Gas)
	if dec.Value == nil {
		return errors.New("missing required field 'value' in transaction")
	}
	tx.Value, overflow = uint256.FromBig(dec.Value.ToInt())
	if overflow {
		return errors.New("'value' in transaction does not fit in 256 bits")
	}
	if dec.Data == nil {
		return errors.New("missing required field 'input' in transaction")
	}
	tx.Data = *dec.Data
	if dec.V == nil {
		return errors.New("missing required field 'v' in transaction")
	}
	overflow = tx.V.SetFromBig(dec.V.ToInt())
	if overflow {
		return fmt.Errorf("dec.V higher than 2^256-1")
	}
	if dec.R == nil {
		return errors.New("missing required field 'r' in transaction")
	}
	overflow = tx.R.SetFromBig(dec.R.ToInt())
	if overflow {
		return fmt.Errorf("dec.R higher than 2^256-1")
	}
	if dec.S == nil {
		return errors.New("missing required field 's' in transaction")
	}
	overflow = tx.S.SetFromBig(dec.S.ToInt())
	if overflow {
		return fmt.Errorf("dec.S higher than 2^256-1")
	}
	withSignature := !tx.V.IsZero() || !tx.R.IsZero() || !tx.S.IsZero()
	if withSignature {
		if err := sanityCheckSignature(&tx.V, &tx.R, &tx.S, false); err != nil {
			return err
		}
	}
	return nil
}

func (tx *DynamicFeeTransaction) UnmarshalJSON(input []byte) error {
	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	// Access list is optional for now.
	if dec.AccessList != nil {
		tx.AccessList = *dec.AccessList
	}
	if dec.ChainID == nil {
		return errors.New("missing required field 'chainId' in transaction")
	}
	var overflow bool
	tx.ChainID, overflow = uint256.FromBig(dec.ChainID.ToInt())
	if overflow {
		return errors.New("'chainId' in transaction does not fit in 256 bits")
	}
	if dec.To != nil {
		tx.To = dec.To
	}
	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' in transaction")
	}
	tx.Nonce = uint64(*dec.Nonce)
	if dec.GasPrice == nil {
		return errors.New("missing required field 'gasPrice' in transaction")
	}
	tx.Tip, overflow = uint256.FromBig(dec.Tip.ToInt())
	if overflow {
		return errors.New("'tip' in transaction does not fit in 256 bits")
	}
	tx.FeeCap, overflow = uint256.FromBig(dec.FeeCap.ToInt())
	if overflow {
		return errors.New("'feeCap' in transaction does not fit in 256 bits")
	}
	if dec.Gas == nil {
		return errors.New("missing required field 'gas' in transaction")
	}
	tx.Gas = uint64(*dec.Gas)
	if dec.Value == nil {
		return errors.New("missing required field 'value' in transaction")
	}
	tx.Value, overflow = uint256.FromBig(dec.Value.ToInt())
	if overflow {
		return errors.New("'value' in transaction does not fit in 256 bits")
	}
	if dec.Data == nil {
		return errors.New("missing required field 'input' in transaction")
	}
	tx.Data = *dec.Data
	if dec.V == nil {
		return errors.New("missing required field 'v' in transaction")
	}
	overflow = tx.V.SetFromBig(dec.V.ToInt())
	if overflow {
		return fmt.Errorf("dec.V higher than 2^256-1")
	}
	if dec.R == nil {
		return errors.New("missing required field 'r' in transaction")
	}
	overflow = tx.R.SetFromBig(dec.R.ToInt())
	if overflow {
		return fmt.Errorf("dec.R higher than 2^256-1")
	}
	if dec.S == nil {
		return errors.New("missing required field 's' in transaction")
	}
	overflow = tx.S.SetFromBig(dec.S.ToInt())
	if overflow {
		return fmt.Errorf("dec.S higher than 2^256-1")
	}
	if overflow {
		return errors.New("'s' in transaction does not fit in 256 bits")
	}
	withSignature := !tx.V.IsZero() || !tx.R.IsZero() || !tx.S.IsZero()
	if withSignature {
		if err := sanityCheckSignature(&tx.V, &tx.R, &tx.S, false); err != nil {
			return err
		}
	}
	return nil
}

func (tx *SignedBlobTx) UnmarshalJSON(input []byte) error {
	var dec txJSON
	var overflow bool
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.AccessList != nil {
		tx.Message.AccessList = AccessListView(*dec.AccessList)
	}

	if dec.ChainID == nil {
		return errors.New("missing required field 'chainId' in transaction")
	}
	tx.Message.ChainID.SetFromBig((*big.Int)(dec.ChainID))

	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' in transaction")
	}
	tx.Message.Nonce = view.Uint64View(*dec.Nonce)

	if dec.Tip == nil {
		return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
	}
	overflow = tx.Message.GasTipCap.SetFromBig((*big.Int)(dec.Tip))
	if overflow {
		return errors.New("'tip' in transaction does not fit in 256 bits")
	}

	if dec.FeeCap == nil {
		return errors.New("missing required field 'maxFeePerGas' for txdata")
	}
	overflow = tx.Message.GasFeeCap.SetFromBig((*big.Int)(dec.FeeCap))
	if overflow {
		return errors.New("'feeCap' in transaction does not fit in 256 bits")
	}

	if dec.Gas == nil {
		return errors.New("missing required field 'gas' for txdata")
	}
	tx.Message.Gas = view.Uint64View(*dec.Gas)

	if dec.Value == nil {
		return errors.New("missing required field 'value' in transaction")
	}
	overflow = tx.Message.Value.SetFromBig((*big.Int)(dec.Value))
	if overflow {
		return errors.New("'value' in transaction does not fit in 256 bits")
	}

	if dec.Data == nil {
		return errors.New("missing required field 'input' in transaction")
	}
	tx.Message.Data = TxDataView(*dec.Data)

	if dec.V == nil {
		return errors.New("missing required field 'v' in transaction")
	}
	tx.Signature.V = view.Uint8View((*big.Int)(dec.V).Uint64())

	if dec.R == nil {
		return errors.New("missing required field 'r' in transaction")
	}
	tx.Signature.R.SetFromBig((*big.Int)(dec.R))

	if dec.S == nil {
		return errors.New("missing required field 's' in transaction")
	}

	tx.Signature.S.SetFromBig((*big.Int)(dec.S))
	withSignature := (*big.Int)(dec.V).Sign() != 0 || (*big.Int)(dec.R).Sign() != 0 || (*big.Int)(dec.S).Sign() != 0
	if withSignature {
		if err := sanityCheckSignature(uint256.NewInt(uint64(tx.Signature.V)), (*uint256.Int)(&tx.Signature.R), (*uint256.Int)(&tx.Signature.S), false); err != nil {
			return err
		}
	}
	tx.Message.MaxFeePerDataGas.SetFromBig((*big.Int)(dec.MaxFeePerDataGas))
	if dec.MaxFeePerDataGas == nil {
		return errors.New("missing required field 'maxFeePerDataGas' for txdata")
	}
	tx.Message.BlobVersionedHashes = dec.BlobVersionedHashes
	// A BlobTx may not contain data
	if len(dec.Blobs) != 0 || len(dec.BlobKzgs) != 0 {
		tx.WrapData = &BlobTxWrapData{
			BlobKzgs:           dec.BlobKzgs,
			Blobs:              dec.Blobs,
			KzgAggregatedProof: dec.KzgAggregatedProof,
		}
		// Verify that versioned hashes match kzgs, and kzgs match blobs.
		if err := tx.WrapData.validateBlobTransactionWrapper(tx); err != nil {
			return fmt.Errorf("blob wrapping data is invalid: %v", err)
		}
	}

	return nil
}
