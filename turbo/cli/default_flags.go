package cli

import (
	"github.com/ledgerwatch/erigon/cmd/utils"

	"github.com/urfave/cli"
)

// DefaultFlags contains all flags that are used and supported by Erigon binary.
var DefaultFlags = []cli.Flag{
	utils.DataDirFlag,
	utils.EthashDatasetDirFlag,
	utils.SnapshotFlag,
	utils.TxPoolDisableFlag,
	utils.TxPoolLocalsFlag,
	utils.TxPoolNoLocalsFlag,
	utils.TxPoolPriceLimitFlag,
	utils.TxPoolPriceBumpFlag,
	utils.TxPoolAccountSlotsFlag,
	utils.TxPoolGlobalSlotsFlag,
	utils.TxPoolGlobalBaseFeeSlotsFlag,
	utils.TxPoolAccountQueueFlag,
	utils.TxPoolGlobalQueueFlag,
	utils.TxPoolLifetimeFlag,
	utils.TxPoolTraceSendersFlag,
	PruneFlag,
	PruneHistoryFlag,
	PruneReceiptFlag,
	PruneTxIndexFlag,
	PruneCallTracesFlag,
	PruneHistoryBeforeFlag,
	PruneReceiptBeforeFlag,
	PruneTxIndexBeforeFlag,
	PruneCallTracesBeforeFlag,
	BatchSizeFlag,
	BlockDownloaderWindowFlag,
	DatabaseVerbosityFlag,
	PrivateApiAddr,
	PrivateApiRateLimit,
	EtlBufferSizeFlag,
	TLSFlag,
	TLSCertFlag,
	TLSKeyFlag,
	TLSCACertFlag,
	StateStreamDisableFlag,
	SyncLoopThrottleFlag,
	BadBlockFlag,

	utils.HTTPEnabledFlag,
	utils.HTTPListenAddrFlag,
	utils.HTTPPortFlag,
	utils.AuthRpcAddr,
	utils.AuthRpcPort,
	utils.JWTSecretPath,
	utils.HttpCompressionFlag,
	utils.HTTPCORSDomainFlag,
	utils.HTTPVirtualHostsFlag,
	utils.AuthRpcVirtualHostsFlag,
	utils.HTTPApiFlag,
	utils.WSEnabledFlag,
	utils.WsCompressionFlag,
	utils.HTTPTraceFlag,
	utils.StateCacheFlag,
	utils.RpcBatchConcurrencyFlag,
	utils.RpcStreamingDisableFlag,
	utils.DBReadConcurrencyFlag,
	utils.RpcAccessListFlag,
	utils.RpcTraceCompatFlag,
	utils.RpcGasCapFlag,
	utils.MemoryOverlayFlag,
	utils.TxpoolApiAddrFlag,
	utils.TraceMaxtracesFlag,
	HTTPReadTimeoutFlag,
	HTTPWriteTimeoutFlag,
	HTTPIdleTimeoutFlag,
	AuthRpcReadTimeoutFlag,
	AuthRpcWriteTimeoutFlag,
	AuthRpcIdleTimeoutFlag,
	EvmCallTimeoutFlag,

	utils.SnapKeepBlocksFlag,
	utils.SnapStopFlag,
	utils.DbPageSizeFlag,
	utils.TorrentPortFlag,
	utils.TorrentMaxPeersFlag,
	utils.TorrentConnsPerFileFlag,
	utils.TorrentDownloadSlotsFlag,
	utils.TorrentUploadRateFlag,
	utils.TorrentDownloadRateFlag,
	utils.TorrentVerbosityFlag,
	utils.ListenPortFlag,
	utils.P2pProtocolVersionFlag,
	utils.NATFlag,
	utils.NoDiscoverFlag,
	utils.DiscoveryV5Flag,
	utils.NetrestrictFlag,
	utils.NodeKeyFileFlag,
	utils.NodeKeyHexFlag,
	utils.DNSDiscoveryFlag,
	utils.BootnodesFlag,
	utils.StaticPeersFlag,
	utils.TrustedPeersFlag,
	utils.MaxPeersFlag,
	utils.ChainFlag,
	utils.DeveloperPeriodFlag,
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
	utils.HistoryV2Flag,
	utils.IdentityFlag,
	utils.CliqueSnapshotCheckpointIntervalFlag,
	utils.CliqueSnapshotInmemorySnapshotsFlag,
	utils.CliqueSnapshotInmemorySignaturesFlag,
	utils.CliqueDataDirFlag,
	utils.EnabledIssuance,
	utils.MiningEnabledFlag,
	utils.ProposingDisableFlag,
	utils.MinerNotifyFlag,
	utils.MinerGasLimitFlag,
	utils.MinerEtherbaseFlag,
	utils.MinerExtraDataFlag,
	utils.MinerNoVerfiyFlag,
	utils.MinerSigningKeyFileFlag,
	utils.SentryAddrFlag,
	utils.SentryLogPeerInfoFlag,
	utils.DownloaderAddrFlag,
	utils.NoDownloaderFlag,
	utils.DownloaderVerifyFlag,
	HealthCheckFlag,
	utils.HeimdallURLFlag,
	utils.WithoutHeimdallFlag,
	utils.EthStatsURLFlag,
	utils.OverrideTerminalTotalDifficulty,
	utils.OverrideMergeNetsplitBlock,

	utils.ConfigFlag,
}
