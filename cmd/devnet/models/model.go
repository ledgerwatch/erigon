package models

import (
	"fmt"
	"github.com/ledgerwatch/erigon/cmd/rpctest/rpctest"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/common/hexutil"
	"github.com/ledgerwatch/erigon/crypto"
	"github.com/ledgerwatch/erigon/p2p"
)

type (
	// TransactionType is the type of transaction attempted to be made, can be regular or contract
	TransactionType string
	// BlockNumber represents the block number type
	BlockNumber string

	// RPCMethod is the type for rpc methods used
	RPCMethod string
	// SubMethod is the type for sub methods used in subscriptions
	SubMethod string
)

const (
	// BuildDirArg is the build directory for the devnet executable
	BuildDirArg = "./build/bin/devnet"
	// DataDirArg is the datadir flag
	DataDirArg = "--datadir"
	// ChainArg is the chain flag
	ChainArg = "--chain"
	// DevPeriodArg is the dev.period flag
	DevPeriodArg = "--dev.period"
	// ConsoleVerbosityArg is the log.console.verbosity flag
	ConsoleVerbosityArg = "--log.console.verbosity"
	// LogDirArg is the log.dir.path flag
	LogDirArg = "--log.dir.path"
	// Mine is the mine flag
	Mine = "--mine"
	// NoDiscover is the nodiscover flag
	NoDiscover = "--nodiscover"
	// PrivateApiAddrArg is the private.api.addr flag
	PrivateApiAddrArg = "--private.api.addr"
	// StaticPeersArg is the staticpeers flag
	StaticPeersArg = "--staticpeers"
	// HttpApiArg is the http.api flag
	HttpApiArg = "--http.api"
	// WSArg is the --ws flag for rpcdaemon
	WSArg = "--ws"

	// DataDirParam is the datadir parameter
	DataDirParam = "./dev"
	// ChainParam is the chain parameter
	ChainParam = "dev"
	// DevPeriodParam is the dev.period parameter
	DevPeriodParam = "0"
	// ConsoleVerbosityParam is the verbosity parameter for the console logs
	ConsoleVerbosityParam = "0"
	// LogDirParam is the log directory parameter for logging to disk
	LogDirParam = "./cmd/devnet/debug_logs"
	// PrivateApiParamMine is the private.api.addr parameter for the mining node
	PrivateApiParamMine = "localhost:9090"
	// PrivateApiParamNoMine is the private.api.addr parameter for the non-mining node
	PrivateApiParamNoMine = "localhost:9091"
	// HttpApiParam is the http.api default parameter for rpcdaemon
	HttpApiParam = "admin,eth,erigon,web3,net,debug,trace,txpool,parity"
	//// WSParam is the ws default parameter for rpcdaemon
	//WSParam = ""

	// ErigonUrl is the default url for rpc connections
	ErigonUrl = "http://localhost:8545"
	// Localhost is the default localhost endpoint for web socket attachments
	Localhost = "127.0.0.1:8545"

	// ReqId is the request id for each request
	ReqId = 0
	// MaxNumberOfBlockChecks is the max number of blocks to look for a transaction in
	MaxNumberOfBlockChecks = 1

	// Latest is the parameter for the latest block
	Latest BlockNumber = "latest"
	// Earliest is the parameter for the earliest block
	Earliest BlockNumber = "earliest"
	// Pending is the parameter for the pending block
	Pending BlockNumber = "pending"

	// hexPrivateKey is the hex value for the private key
	hexPrivateKey = "26e86e45f6fc45ec6e2ecd128cec80fa1d1505e5507dcd2ae58c3130a7a97b48"
	// DevAddress is the developer address for sending
	DevAddress = "0x67b1d87101671b127f5f8714789C7192f7ad340e"

	// NonContractTx is the transaction type for sending ether
	NonContractTx TransactionType = "non-contract"
	// ContractTx is the transaction type for sending ether
	ContractTx TransactionType = "contract"

	// ETHGetTransactionCount represents the eth_getTransactionCount method
	ETHGetTransactionCount RPCMethod = "eth_getTransactionCount"
	// ETHGetBalance represents the eth_getBalance method
	ETHGetBalance RPCMethod = "eth_getBalance"
	// ETHSendRawTransaction represents the eth_sendRawTransaction method
	ETHSendRawTransaction RPCMethod = "eth_sendRawTransaction"
	// ETHGetBlockByNumber represents the eth_getBlockByNumber method
	ETHGetBlockByNumber RPCMethod = "eth_getBlockByNumber"
	// AdminNodeInfo represents the admin_nodeInfo method
	AdminNodeInfo RPCMethod = "admin_nodeInfo"
	// TxpoolContent represents the txpool_content method
	TxpoolContent RPCMethod = "txpool_content"

	// ETHNewHeads represents the eth_newHeads sub method
	ETHNewHeads SubMethod = "eth_newHeads"
)

var (
	// DevSignedPrivateKey is the signed private key for signing transactions
	DevSignedPrivateKey, _ = crypto.HexToECDSA(hexPrivateKey)
)

// Responses for the rpc calls
type (
	// TxpoolContextResponse is the response for calls made to txpool_content
	TxpoolContextResponse map[string]interface{}

	// AdminNodeInfoResponse is the response for calls made to admin_nodeInfo
	AdminNodeInfoResponse struct {
		rpctest.CommonResponse
		Result p2p.NodeInfo `json:"result"`
	}

	//content := map[string]map[string]map[string]*RPCTransaction{
	//	"pending": make(map[string]map[string]*RPCTransaction),
	//	"baseFee": make(map[string]map[string]*RPCTransaction),
	//	"queued":  make(map[string]map[string]*RPCTransaction),
	//}
)

// Block represents a block
type Block struct {
	Number       *hexutil.Big
	Transactions []common.Hash
}

// ParameterFromArgument merges the argument and parameter and returns a flag input string
func ParameterFromArgument(arg, param string) (string, error) {
	if arg == "" {
		return "", ErrInvalidArgument
	}
	return fmt.Sprintf("%s=%s", arg, param), nil
}
