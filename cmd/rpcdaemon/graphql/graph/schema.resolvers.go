package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ledgerwatch/erigon/cmd/rpcdaemon/graphql/graph/model"
	"github.com/ledgerwatch/erigon/common/hexutil"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/rpc"
)

// SendRawTransaction is the resolver for the sendRawTransaction field.
func (r *mutationResolver) SendRawTransaction(ctx context.Context, data string) (string, error) {
	panic(fmt.Errorf("not implemented: SendRawTransaction - sendRawTransaction"))
}

// Block is the resolver for the block field.
func (r *queryResolver) Block(ctx context.Context, number *string, hash *string) (*model.Block, error) {
	var blockNumber rpc.BlockNumber

	if number != nil {
		// Block number is not null, test for a positive long integer
		bNum, err := strconv.ParseUint(*number, 10, 64)
		if err == nil {
			// Positive integer, go ahead
			blockNumber = rpc.BlockNumber(bNum)
		} else {
			bNum, err := hexutil.DecodeUint64(*number)
			if err == nil {
				// Hexadecimal, 0x prefixed
				blockNumber = rpc.BlockNumber(bNum)
			} else {
				var err error
				return nil, err
			}
		}
	} else {
		if hash != nil {
			blockHash, _ := hexutil.DecodeBig(*hash)
			fmt.Println("TODO/GraphQL/Implement me, get Block by hash=", blockHash)
			hash = nil
		}
	}

	if number == nil && hash == nil {
		// If neither number or hash is specified (nil), we should deliver "latest" block
		// blockNumber = rpc.LatestExecutedBlockNumber
		blockNumber = rpc.LatestBlockNumber
	}

	res, err := r.GraphQLAPI.GetBlockDetails(ctx, blockNumber)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	block := &model.Block{}
	absBlk := res["block"]

	if absBlk != nil {
		blk := absBlk.(map[string]interface{})

		block.Difficulty = *convertDataToStringP(blk, "difficulty")
		block.ExtraData = *convertDataToStringP(blk, "extraData")
		block.GasLimit = uint64(*convertDataToUint64P(blk, "gasLimit"))
		block.GasUsed = *convertDataToUint64P(blk, "gasUsed")
		block.Hash = *convertDataToStringP(blk, "hash")
		block.Miner = &model.Account{}
		address := convertDataToStringP(blk, "miner")
		if address != nil {
			block.Miner.Address = strings.ToLower(*address)
		}
		mixHash := convertDataToStringP(blk, "mixHash")
		if mixHash != nil {
			block.MixHash = *mixHash
		}
		blockNonce := convertDataToStringP(blk, "nonce")
		if blockNonce != nil {
			block.Nonce = *blockNonce
		}
		block.Number = *convertDataToUint64P(blk, "number")
		block.Ommers = []*model.Block{}
		block.Parent = &model.Block{}
		block.Parent.Hash = *convertDataToStringP(blk, "parentHash")
		block.ReceiptsRoot = *convertDataToStringP(blk, "receiptsRoot")
		block.StateRoot = *convertDataToStringP(blk, "stateRoot")
		block.Timestamp = *convertDataToUint64P(blk, "timestamp") // int in the schema but Geth displays in HEX !!!
		block.TransactionCount = convertDataToIntP(blk, "transactionCount")
		block.TransactionsRoot = *convertDataToStringP(blk, "transactionsRoot")
		block.TotalDifficulty = *convertDataToStringP(blk, "totalDifficulty")
		block.Transactions = []*model.Transaction{}

		block.LogsBloom = "0x" + *convertDataToStringP(blk, "logsBloom")
		block.OmmerHash = "" // OmmerHash:     gointerfaces.ConvertHashToH256(header.UncleHash),

		/*
			Missing Block fields to fill :
			- ommerHash
		*/

		absRcp := res["receipts"]
		rcp := absRcp.([]map[string]interface{})
		for _, transReceipt := range rcp {
			trans := &model.Transaction{}
			trans.CumulativeGasUsed = convertDataToUint64P(transReceipt, "cumulativeGasUsed")
			trans.InputData = *convertDataToStringP(transReceipt, "data")
			trans.EffectiveGasPrice = convertDataToStringP(transReceipt, "effectiveGasPrice")
			trans.GasPrice = *convertDataToStringP(transReceipt, "gasPrice")
			trans.GasUsed = convertDataToUint64P(transReceipt, "gasUsed")
			trans.Hash = *convertDataToStringP(transReceipt, "transactionHash")
			trans.Index = convertDataToIntP(transReceipt, "transactionIndex")
			trans.Nonce = *convertDataToUint64P(transReceipt, "nonce")
			trans.Status = convertDataToUint64P(transReceipt, "status")
			trans.Type = convertDataToIntP(transReceipt, "type")
			trans.Value = *convertDataToStringP(transReceipt, "value")

			trans.Logs = make([]*model.Log, 0)
			for _, rlog := range transReceipt["logs"].(types.Logs) {
				tlog := model.Log{
					Index: int(rlog.Index),
					Data:  "0x" + hex.EncodeToString(rlog.Data),
				}
				tlog.Account = &model.Account{}
				tlog.Account.Address = rlog.Address.String()

				for _, rtopic := range rlog.Topics {
					tlog.Topics = append(tlog.Topics, rtopic.String())
				}

				trans.Logs = append(trans.Logs, &tlog)
			}

			trans.From = &model.Account{}
			trans.From.Address = strings.ToLower(*convertDataToStringP(transReceipt, "from"))

			trans.To = &model.Account{}
			address := convertDataToStringP(transReceipt, "to")
			// To address could be nil in case of contract creation
			if address != nil {
				trans.To.Address = strings.ToLower(*convertDataToStringP(transReceipt, "to"))
			}

			block.Transactions = append(block.Transactions, trans)
		}
	}

	return block, ctx.Err()
}

// Blocks is the resolver for the blocks field.
func (r *queryResolver) Blocks(ctx context.Context, from *uint64, to *uint64) ([]*model.Block, error) {
	panic(fmt.Errorf("not implemented: Blocks - blocks"))
}

// Pending is the resolver for the pending field.
func (r *queryResolver) Pending(ctx context.Context) (*model.Pending, error) {
	panic(fmt.Errorf("not implemented: Pending - pending"))
}

// Transaction is the resolver for the transaction field.
func (r *queryResolver) Transaction(ctx context.Context, hash string) (*model.Transaction, error) {
	panic(fmt.Errorf("not implemented: Transaction - transaction"))
}

// Logs is the resolver for the logs field.
func (r *queryResolver) Logs(ctx context.Context, filter model.FilterCriteria) ([]*model.Log, error) {
	panic(fmt.Errorf("not implemented: Logs - logs"))
}

// GasPrice is the resolver for the gasPrice field.
func (r *queryResolver) GasPrice(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented: GasPrice - gasPrice"))
}

// MaxPriorityFeePerGas is the resolver for the maxPriorityFeePerGas field.
func (r *queryResolver) MaxPriorityFeePerGas(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented: MaxPriorityFeePerGas - maxPriorityFeePerGas"))
}

// Syncing is the resolver for the syncing field.
func (r *queryResolver) Syncing(ctx context.Context) (*model.SyncState, error) {
	panic(fmt.Errorf("not implemented: Syncing - syncing"))
}

// ChainID is the resolver for the chainID field.
func (r *queryResolver) ChainID(ctx context.Context) (string, error) {
	chainID, err := r.GraphQLAPI.GetChainID(ctx)

	return "0x" + strconv.FormatUint(chainID.Uint64(), 16), err
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
