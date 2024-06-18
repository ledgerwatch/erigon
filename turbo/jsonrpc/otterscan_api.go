package jsonrpc

import (
	"context"
	"fmt"
	"math/big"

	hexutil2 "github.com/ledgerwatch/erigon-lib/common/hexutil"
	"github.com/ledgerwatch/erigon-lib/kv/order"
	"github.com/ledgerwatch/erigon-lib/kv/rawdbv3"
	"github.com/ledgerwatch/log/v3"
	"golang.org/x/sync/errgroup"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon-lib/chain"
	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/hexutility"
	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/consensus"
	"github.com/ledgerwatch/erigon/core"
	"github.com/ledgerwatch/erigon/core/rawdb"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/core/vm"
	"github.com/ledgerwatch/erigon/eth/ethutils"
	"github.com/ledgerwatch/erigon/rpc"
	"github.com/ledgerwatch/erigon/turbo/adapter/ethapi"
	"github.com/ledgerwatch/erigon/turbo/rpchelper"
	"github.com/ledgerwatch/erigon/turbo/transactions"
)

// API_LEVEL Must be incremented every time new additions are made
const API_LEVEL = 8

type TransactionsWithReceipts struct {
	Txs       []*RPCTransaction        `json:"txs"`
	Receipts  []map[string]interface{} `json:"receipts"`
	FirstPage bool                     `json:"firstPage"`
	LastPage  bool                     `json:"lastPage"`
}

type OtterscanAPI interface {
	GetApiLevel() uint8
	GetInternalOperations(ctx context.Context, hash common.Hash) ([]*InternalOperation, error)
	SearchTransactionsBefore(ctx context.Context, addr common.Address, blockNum uint64, pageSize uint16) (*TransactionsWithReceipts, error)
	SearchTransactionsAfter(ctx context.Context, addr common.Address, blockNum uint64, pageSize uint16) (*TransactionsWithReceipts, error)
	GetBlockDetails(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error)
	GetBlockDetailsByHash(ctx context.Context, hash common.Hash) (map[string]interface{}, error)
	GetBlockTransactions(ctx context.Context, number rpc.BlockNumber, pageNumber uint8, pageSize uint8) (map[string]interface{}, error)
	HasCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (bool, error)
	TraceTransaction(ctx context.Context, hash common.Hash) ([]*TraceEntry, error)
	GetTransactionError(ctx context.Context, hash common.Hash) (hexutility.Bytes, error)
	GetTransactionBySenderAndNonce(ctx context.Context, addr common.Address, nonce uint64) (*common.Hash, error)
	GetContractCreator(ctx context.Context, addr common.Address) (*ContractCreatorData, error)
	Debug(ctx context.Context) error
	Debug2(ctx context.Context) error
	Debug3(ctx context.Context) error
	Debug4(ctx context.Context) error
}

// TEST ADDR with 10 creation/self-destructs on sepolia
var addr = common.HexToAddress("0xC4a96da3483Ccd935944a0eeCdB20aE42476C296")

func (api *OtterscanAPIImpl) Debug(ctx context.Context) error {
	tx, err := api.db.BeginRo(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ttx := tx.(kv.TemporalTx)

	startBlock := uint64(0)
	minTxNum, err := rawdbv3.TxNums.Min(ttx, startBlock)
	if err != nil {
		return err
	}
	maxTxNum, err := rawdbv3.TxNums.Max(ttx, ^uint64(0))
	if err != nil {
		return err
	}
	log.Info("Range", "startBlock", startBlock, "minTxNum", minTxNum, "maxTxNum", maxTxNum, "diff", maxTxNum-minTxNum+1)

	return nil
}

func (api *OtterscanAPIImpl) Debug2(ctx context.Context) error {
	tx, err := api.db.BeginRo(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ttx := tx.(kv.TemporalTx)

	startBlock := uint64(0)
	minTxNum, err := rawdbv3.TxNums.Min(ttx, startBlock)
	if err != nil {
		return err
	}
	maxTxNum, err := rawdbv3.TxNums.Max(ttx, ^uint64(0))
	if err != nil {
		return err
	}
	log.Info("Range", "startBlock", startBlock, "minTxNum", minTxNum, "maxTxNum", maxTxNum, "diff", maxTxNum-minTxNum+1)

	addr := common.HexToAddress("0x35b7f08dc14d16c47c0ae202c901fd64b739966c")

	it2, err := ttx.IndexRange(kv.CodeHistoryIdx, addr.Bytes(), int(minTxNum), int(maxTxNum+1), order.Asc, kv.Unlim)
	if err != nil {
		return err
	}
	log.Info("CHANGES", "addr", addr)
	it3 := rawdbv3.TxNums2BlockNums(ttx, it2, order.Asc)
	for it3.HasNext() {
		_, blockNum, txIdx, isFinalTxn, _, err := it3.Next()
		if err != nil {
			return err
		}
		log.Info("TX", "blockNum", blockNum, "txIdx", txIdx, "isFinal", isFinalTxn)
	}

	return nil
}

func (api *OtterscanAPIImpl) Debug3(ctx context.Context) error {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ttx := tx.(kv.TemporalTx)

	testKey := []byte("test")

	// startBlock := uint64(0)
	// txNum, err := rawdbv3.TxNums.Min(ttx, startBlock)
	// if err != nil {
	// 	return err
	// }
	// k := dbutils.EncodeBlockNumber(txNum)
	ts, err := ttx.IndexRange(kv.TblTestIIIdx, testKey, -1, -1, order.Asc, kv.Unlim)
	if err != nil {
		return err
	}
	it := rawdbv3.TxNums2BlockNums(tx, ts, order.Asc)
	log.Info("XXXXX")
	count := 0
	for it.HasNext() && count < 100 {
		count++
		txNum, blockNum, txIndex, _, _, err := it.Next()
		if err != nil {
			return err
		}
		log.Info("AAAAA", "txNum", txNum, "blockNum", blockNum, "txIndex", txIndex)
	}
	log.Info("YYYYY")

	return nil
}

func (api *OtterscanAPIImpl) Debug4(ctx context.Context) error {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := api.printDupTable(tx, kv.TblTestIIKeys, 100); err != nil {
		return err
	}
	if err := api.printDupTable(tx, kv.TblTestIIIdx, 100); err != nil {
		return err
	}

	return nil
}

func (api *OtterscanAPIImpl) printDupTable(tx kv.Tx, table string, limit int) error {
	c, err := tx.CursorDupSort(table)
	if err != nil {
		return err
	}
	defer c.Close()

	log.Info("TABLE", "t", table)
	count := 0
	k, v, err := c.First()
	if err != nil {
		return err
	}
	for k != nil && (limit != -1 && count < limit) {
		count++
		log.Info("W", "k", hexutility.Encode(k), "v", hexutility.Encode(v), "c", count)
		k, v, err = c.NextDup()
		if err != nil {
			return err
		}
		if k == nil {
			k, v, err = c.NextNoDup()
			if err != nil {
				return err
			}
			log.Info("D")
		}
	}

	return nil
}

type OtterscanAPIImpl struct {
	*BaseAPI
	db          kv.RoDB
	maxPageSize uint64
}

func NewOtterscanAPI(base *BaseAPI, db kv.RoDB, maxPageSize uint64) *OtterscanAPIImpl {
	return &OtterscanAPIImpl{
		BaseAPI:     base,
		db:          db,
		maxPageSize: maxPageSize,
	}
}

func (api *OtterscanAPIImpl) GetApiLevel() uint8 {
	return API_LEVEL
}

// TODO: dedup from eth_txs.go#GetTransactionByHash
func (api *OtterscanAPIImpl) getTransactionByHash(ctx context.Context, tx kv.Tx, hash common.Hash) (types.Transaction, *types.Block, common.Hash, uint64, uint64, error) {
	// https://infura.io/docs/ethereum/json-rpc/eth-getTransactionByHash
	blockNum, ok, err := api.txnLookup(ctx, tx, hash)
	if err != nil {
		return nil, nil, common.Hash{}, 0, 0, err
	}
	if !ok {
		return nil, nil, common.Hash{}, 0, 0, nil
	}

	block, err := api.blockByNumberWithSenders(ctx, tx, blockNum)
	if err != nil {
		return nil, nil, common.Hash{}, 0, 0, err
	}
	if block == nil {
		return nil, nil, common.Hash{}, 0, 0, nil
	}
	blockHash := block.Hash()
	var txnIndex uint64
	var txn types.Transaction
	for i, transaction := range block.Transactions() {
		if transaction.Hash() == hash {
			txn = transaction
			txnIndex = uint64(i)
			break
		}
	}

	// Add GasPrice for the DynamicFeeTransaction
	// var baseFee *big.Int
	// if chainConfig.IsLondon(blockNum) && blockHash != (common.Hash{}) {
	// 	baseFee = block.BaseFee()
	// }

	// if no transaction was found then we return nil
	if txn == nil {
		return nil, nil, common.Hash{}, 0, 0, nil
	}
	return txn, block, blockHash, blockNum, txnIndex, nil
}

func (api *OtterscanAPIImpl) runTracer(ctx context.Context, tx kv.Tx, hash common.Hash, tracer vm.EVMLogger) (*core.ExecutionResult, error) {
	txn, block, _, _, txIndex, err := api.getTransactionByHash(ctx, tx, hash)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, fmt.Errorf("transaction %#x not found", hash)
	}

	chainConfig, err := api.chainConfig(ctx, tx)
	if err != nil {
		return nil, err
	}
	engine := api.engine()

	msg, blockCtx, txCtx, ibs, _, err := transactions.ComputeTxEnv(ctx, engine, block, chainConfig, api._blockReader, tx, int(txIndex))
	if err != nil {
		return nil, err
	}

	var vmConfig vm.Config
	if tracer == nil {
		vmConfig = vm.Config{}
	} else {
		vmConfig = vm.Config{Debug: true, Tracer: tracer}
	}
	vmenv := vm.NewEVM(blockCtx, txCtx, ibs, chainConfig, vmConfig)

	result, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()).AddBlobGas(msg.BlobGas()), true, false /* gasBailout */)
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %v", err)
	}

	return result, nil
}

func (api *OtterscanAPIImpl) GetInternalOperations(ctx context.Context, hash common.Hash) ([]*InternalOperation, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	tracer := NewOperationsTracer(ctx)
	if _, err := api.runTracer(ctx, tx, hash, tracer); err != nil {
		return nil, err
	}

	return tracer.Results, nil
}

// Search transactions that touch a certain address.
//
// It searches back a certain block (excluding); the results are sorted descending.
//
// The pageSize indicates how many txs may be returned. If there are less txs than pageSize,
// they are just returned. But it may return a little more than pageSize if there are more txs
// than the necessary to fill pageSize in the last found block, i.e., let's say you want pageSize == 25,
// you already found 24 txs, the next block contains 4 matches, then this function will return 28 txs.
func (api *OtterscanAPIImpl) SearchTransactionsBefore(ctx context.Context, addr common.Address, blockNum uint64, pageSize uint16) (*TransactionsWithReceipts, error) {
	if uint64(pageSize) > api.maxPageSize {
		return nil, fmt.Errorf("max allowed page size: %v", api.maxPageSize)
	}

	dbtx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer dbtx.Rollback()

	return api.searchTransactionsBeforeV3(dbtx.(kv.TemporalTx), ctx, addr, blockNum, pageSize)
}

// Search transactions that touch a certain address.
//
// It searches forward a certain block (excluding); the results are sorted descending.
//
// The pageSize indicates how many txs may be returned. If there are less txs than pageSize,
// they are just returned. But it may return a little more than pageSize if there are more txs
// than the necessary to fill pageSize in the last found block, i.e., let's say you want pageSize == 25,
// you already found 24 txs, the next block contains 4 matches, then this function will return 28 txs.
func (api *OtterscanAPIImpl) SearchTransactionsAfter(ctx context.Context, addr common.Address, blockNum uint64, pageSize uint16) (*TransactionsWithReceipts, error) {
	if uint64(pageSize) > api.maxPageSize {
		return nil, fmt.Errorf("max allowed page size: %v", api.maxPageSize)
	}

	dbtx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer dbtx.Rollback()

	return api.searchTransactionsAfterV3(dbtx.(kv.TemporalTx), ctx, addr, blockNum, pageSize)
}

func (api *OtterscanAPIImpl) traceBlocks(ctx context.Context, addr common.Address, chainConfig *chain.Config, pageSize, resultCount uint16, callFromToProvider BlockProvider) ([]*TransactionsWithReceipts, bool, error) {
	// Estimate the common case of user address having at most 1 interaction/block and
	// trace N := remaining page matches as number of blocks to trace concurrently.
	// TODO: this is not optimimal for big contract addresses; implement some better heuristics.
	estBlocksToTrace := pageSize - resultCount
	results := make([]*TransactionsWithReceipts, estBlocksToTrace)
	totalBlocksTraced := 0
	hasMore := true

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(1024) // we don't want limit much here, but protecting from infinity attack
	for i := 0; i < int(estBlocksToTrace); i++ {
		i := i // we will pass it to goroutine

		var nextBlock uint64
		var err error
		nextBlock, hasMore, err = callFromToProvider()
		if err != nil {
			return nil, false, err
		}
		// TODO: nextBlock == 0 seems redundant with hasMore == false
		if !hasMore && nextBlock == 0 {
			break
		}

		totalBlocksTraced++

		eg.Go(func() error {
			// don't return error from searchTraceBlock - to avoid 1 block fail impact to other blocks
			// if return error - `errgroup` will interrupt all other goroutines
			// but passing `ctx` - then user still can cancel request
			api.searchTraceBlock(ctx, addr, chainConfig, i, nextBlock, results)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, false, err
	}

	return results[:totalBlocksTraced], hasMore, nil
}

func delegateGetBlockByNumber(tx kv.Tx, b *types.Block, number rpc.BlockNumber, inclTx bool) (map[string]interface{}, error) {
	td, err := rawdb.ReadTd(tx, b.Hash(), b.NumberU64())
	if err != nil {
		return nil, err
	}
	additionalFields := make(map[string]interface{})
	response, err := ethapi.RPCMarshalBlock(b, inclTx, inclTx, additionalFields)
	if !inclTx {
		delete(response, "transactions") // workaround for https://github.com/ledgerwatch/erigon/issues/4989#issuecomment-1218415666
	}
	response["totalDifficulty"] = (*hexutil2.Big)(td)
	response["transactionCount"] = b.Transactions().Len()

	if err == nil && number == rpc.PendingBlockNumber {
		// Pending blocks need to nil out a few fields
		for _, field := range []string{"hash", "nonce", "miner"} {
			response[field] = nil
		}
	}

	// Explicitly drop unwanted fields
	response["logsBloom"] = nil
	return response, err
}

// TODO: temporary workaround due to API breakage from watch_the_burn
type internalIssuance struct {
	BlockReward string `json:"blockReward,omitempty"`
	UncleReward string `json:"uncleReward,omitempty"`
	Issuance    string `json:"issuance,omitempty"`
}

func delegateIssuance(tx kv.Tx, block *types.Block, chainConfig *chain.Config, engine consensus.EngineReader) (internalIssuance, error) {
	// TODO: aura seems to be already broken in the original version of this RPC method
	rewards, err := engine.CalculateRewards(chainConfig, block.HeaderNoCopy(), block.Uncles(), func(contract common.Address, data []byte) ([]byte, error) {
		return nil, nil
	})
	if err != nil {
		return internalIssuance{}, err
	}

	blockReward := uint256.NewInt(0)
	uncleReward := uint256.NewInt(0)
	for _, r := range rewards {
		if r.Kind == consensus.RewardAuthor {
			blockReward.Add(blockReward, &r.Amount)
		}
		if r.Kind == consensus.RewardUncle {
			uncleReward.Add(uncleReward, &r.Amount)
		}
	}

	var ret internalIssuance
	ret.BlockReward = hexutil2.EncodeBig(blockReward.ToBig())
	ret.UncleReward = hexutil2.EncodeBig(uncleReward.ToBig())

	blockReward.Add(blockReward, uncleReward)
	ret.Issuance = hexutil2.EncodeBig(blockReward.ToBig())
	return ret, nil
}

func delegateBlockFees(ctx context.Context, tx kv.Tx, block *types.Block, senders []common.Address, chainConfig *chain.Config, receipts types.Receipts) (*big.Int, error) {
	fee := big.NewInt(0)
	gasUsed := big.NewInt(0)

	totalFees := big.NewInt(0)
	for _, receipt := range receipts {
		txn := block.Transactions()[receipt.TransactionIndex]
		effectiveGasPrice := uint64(0)
		if !chainConfig.IsLondon(block.NumberU64()) {
			effectiveGasPrice = txn.GetPrice().Uint64()
		} else {
			baseFee, _ := uint256.FromBig(block.BaseFee())
			gasPrice := new(big.Int).Add(block.BaseFee(), txn.GetEffectiveGasTip(baseFee).ToBig())
			effectiveGasPrice = gasPrice.Uint64()
		}

		fee.SetUint64(effectiveGasPrice)
		gasUsed.SetUint64(receipt.GasUsed)
		fee.Mul(fee, gasUsed)

		totalFees.Add(totalFees, fee)
	}

	return totalFees, nil
}

func (api *OtterscanAPIImpl) getBlockWithSenders(ctx context.Context, number rpc.BlockNumber, tx kv.Tx) (*types.Block, []common.Address, error) {
	if number == rpc.PendingBlockNumber {
		return api.pendingBlock(), nil, nil
	}

	n, hash, _, err := rpchelper.GetBlockNumber(rpc.BlockNumberOrHashWithNumber(number), tx, api.filters)
	if err != nil {
		return nil, nil, err
	}

	block, err := api.blockWithSenders(ctx, tx, hash, n)
	if err != nil {
		return nil, nil, err
	}
	if block == nil {
		return nil, nil, nil
	}
	return block, block.Body().SendersFromTxs(), nil
}

func (api *OtterscanAPIImpl) GetBlockTransactions(ctx context.Context, number rpc.BlockNumber, pageNumber uint8, pageSize uint8) (map[string]interface{}, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b, senders, err := api.getBlockWithSenders(ctx, number, tx)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}

	chainConfig, err := api.chainConfig(ctx, tx)
	if err != nil {
		return nil, err
	}

	getBlockRes, err := delegateGetBlockByNumber(tx, b, number, true)
	if err != nil {
		return nil, err
	}

	// Receipts
	receipts, err := api.getReceipts(ctx, tx, b, senders)
	if err != nil {
		return nil, fmt.Errorf("getReceipts error: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(receipts))
	for _, receipt := range receipts {
		txn := b.Transactions()[receipt.TransactionIndex]
		marshalledRcpt := ethutils.MarshalReceipt(receipt, txn, chainConfig, b.HeaderNoCopy(), txn.Hash(), true)
		marshalledRcpt["logs"] = nil
		marshalledRcpt["logsBloom"] = nil
		result = append(result, marshalledRcpt)
	}

	// Pruned block attrs
	prunedBlock := map[string]interface{}{}
	for _, k := range []string{"timestamp", "miner", "baseFeePerGas"} {
		prunedBlock[k] = getBlockRes[k]
	}

	// Crop tx input to 4bytes
	var txs = getBlockRes["transactions"].([]interface{})
	for _, rawTx := range txs {
		rpcTx := rawTx.(*ethapi.RPCTransaction)
		if len(rpcTx.Input) >= 4 {
			rpcTx.Input = rpcTx.Input[:4]
		}
	}

	// Crop page
	pageEnd := b.Transactions().Len() - int(pageNumber)*int(pageSize)
	pageStart := pageEnd - int(pageSize)
	if pageEnd < 0 {
		pageEnd = 0
	}
	if pageStart < 0 {
		pageStart = 0
	}

	response := map[string]interface{}{}
	getBlockRes["transactions"] = getBlockRes["transactions"].([]interface{})[pageStart:pageEnd]
	response["fullblock"] = getBlockRes
	response["receipts"] = result[pageStart:pageEnd]
	return response, nil
}
