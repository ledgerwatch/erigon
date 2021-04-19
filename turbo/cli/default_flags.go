package cli

import (
	"github.com/ledgerwatch/turbo-geth/cmd/utils"

	"github.com/urfave/cli"
)

// DefaultFlags contains all flags that are used and supported by turbo-geth binary.
var DefaultFlags = []cli.Flag{
	utils.DataDirFlag,
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
	StorageModeFlag,
	SnapshotModeFlag,
	SeedSnapshotsFlag,
	ExternalSnapshotDownloaderAddrFlag,
	//CacheSizeFlag,
	BatchSizeFlag,
	DatabaseFlag,
	PrivateApiAddr,
	EtlBufferSizeFlag,
	LMDBMapSizeFlag,
	TLSFlag,
	TLSCertFlag,
	TLSKeyFlag,
	TLSCACertFlag,
	utils.ListenPortFlag,
	utils.NATFlag,
	utils.NoDiscoverFlag,
	utils.DiscoveryV5Flag,
	utils.NetrestrictFlag,
	utils.NodeKeyFileFlag,
	utils.NodeKeyHexFlag,
	utils.DNSDiscoveryFlag,
	utils.ChainFlag,
	utils.VMEnableDebugFlag,
	utils.NetworkIdFlag,
	utils.FakePoWFlag,
	utils.GpoBlocksFlag,
	utils.GpoPercentileFlag,
	utils.InsecureUnlockAllowedFlag,
	utils.MetricsEnabledFlag,
	utils.MetricsEnabledExpensiveFlag,
	utils.MetricsHTTPFlag,
	utils.MetricsPortFlag,
	utils.IdentityFlag,
	utils.CliqueSnapshotCheckpointIntervalFlag,
	utils.CliqueSnapshotInmemorySnapshotsFlag,
	utils.CliqueSnapshotInmemorySignaturesFlag,
	utils.CliqueDataDirFlag,
	SilkwormFlag,
	utils.MiningEnabledFlag,
	utils.MinerNotifyFlag,
	utils.MinerGasTargetFlag,
	utils.MinerGasLimitFlag,
	utils.MinerEtherbaseFlag,
	utils.MinerExtraDataFlag,
	utils.MinerNoVerfiyFlag,
}
