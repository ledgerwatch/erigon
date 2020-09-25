package commands

import (
	"context"
	"fmt"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/rpc"
)

// CallParam a parameter for a trace_callMany routine
type CallParam struct {
	_ types.Transaction
	_ []string
}

// CallParams array of callMany structs
type CallParams []CallParam

// Call Implements trace_call
func (api *TraceAPIImpl) Call(ctx context.Context, call CallParam, blockNr rpc.BlockNumber) ([]interface{}, error) {
	var stub []interface{}
	return stub, fmt.Errorf("function trace_call not implemented")
}

// CallMany Implements trace_call
func (api *TraceAPIImpl) CallMany(ctx context.Context, calls CallParams) ([]interface{}, error) {
	var stub []interface{}
	return stub, fmt.Errorf("function trace_callMany not implemented")
}

// RawTransaction Implements trace_rawtransaction
func (api *TraceAPIImpl) RawTransaction(ctx context.Context, txHash common.Hash, traceTypes []string) ([]interface{}, error) {
	var stub []interface{}
	return stub, fmt.Errorf("function trace_rawTransaction not implemented")
}

// ReplayBlockTransactions Implements trace_replayBlockTransactions
func (api *TraceAPIImpl) ReplayBlockTransactions(ctx context.Context, blockNr rpc.BlockNumber, traceTypes []string) ([]interface{}, error) {
	var stub []interface{}
	return stub, fmt.Errorf("function trace_replayBlockTransactions not implemented")
}

// ReplayTransaction Implements trace_replaytransactions
func (api *TraceAPIImpl) ReplayTransaction(ctx context.Context, txHash common.Hash, traceTypes []string) ([]interface{}, error) {
	var stub []interface{}
	return stub, fmt.Errorf("function trace_replayTransaction not implemented")
}
