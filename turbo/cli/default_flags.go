package cli

import (
	"github.com/ledgerwatch/turbo-geth/cmd/utils"

	"github.com/urfave/cli"
)

var DefaultFlags = []cli.Flag{
	utils.DataDirFlag,
	utils.KeyStoreDirFlag,
	utils.EthashDatasetDirFlag,
	utils.TxPoolLocalsFlag,
	utils.TxPoolNoLocalsFlag,
	utils.TxPoolJournalFlag,
	utils.TxPoolRejournalFlag,
	utils.TxPoolPriceLimitFlag,
	utils.TxPoolPriceBumpFlag,
	utils.TxPoolAccountSlotsFlag,
	utils.TxPoolGlobalSlotsFlag,
	utils.TxPoolAccountQueueFlag,
	utils.TxPoolGlobalQueueFlag,
	utils.TxPoolLifetimeFlag,
	utils.TxLookupLimitFlag,
	utils.StorageModeFlag,
	utils.HddFlag,
	utils.DatabaseFlag,
	utils.LMDBMapSizeFlag,
	utils.PrivateApiAddr,
	utils.ListenPortFlag,
	utils.NATFlag,
	utils.NoDiscoverFlag,
	utils.DiscoveryV5Flag,
	utils.NetrestrictFlag,
	utils.NodeKeyFileFlag,
	utils.NodeKeyHexFlag,
	utils.DNSDiscoveryFlag,
	utils.RopstenFlag,
	utils.RinkebyFlag,
	utils.GoerliFlag,
	utils.YoloV1Flag,
	utils.VMEnableDebugFlag,
	utils.NetworkIdFlag,
	utils.FakePoWFlag,
	utils.GpoBlocksFlag,
	utils.GpoPercentileFlag,
	utils.EWASMInterpreterFlag,
	utils.EVMInterpreterFlag,
	utils.IPCDisabledFlag,
	utils.IPCPathFlag,
	utils.InsecureUnlockAllowedFlag,
	utils.MetricsEnabledFlag,
	utils.MetricsEnabledExpensiveFlag,
	utils.MetricsHTTPFlag,
	utils.MetricsPortFlag,
	utils.IdentityFlag,
}
