// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package ethconfig

import (
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/consensus/ethash"
	"github.com/ledgerwatch/erigon/core"
	"github.com/ledgerwatch/erigon/eth/gasprice"
	"github.com/ledgerwatch/erigon/ethdb"
	"github.com/ledgerwatch/erigon/params"
)

// MarshalTOML marshals as TOML.
func (c Config) MarshalTOML() (interface{}, error) {
	type Config struct {
		Genesis          *core.Genesis `toml:",omitempty"`
		NetworkID        uint64
		EthDiscoveryURLs []string
		Whitelist        map[uint64]common.Hash `toml:"-"`
		StorageMode      string
		OnlyAnnounce     bool
		Miner            params.MiningConfig
		Ethash           ethash.Config
		TxPool           core.TxPoolConfig
		GPO              gasprice.Config
		DocRoot          string                         `toml:"-"`
		RPCGasCap        uint64                         `toml:",omitempty"`
		RPCTxFeeCap      float64                        `toml:",omitempty"`
		Checkpoint       *params.TrustedCheckpoint      `toml:",omitempty"`
		CheckpointOracle *params.CheckpointOracleConfig `toml:",omitempty"`
	}
	var enc Config
	enc.Genesis = c.Genesis
	enc.NetworkID = c.NetworkID
	enc.EthDiscoveryURLs = c.EthDiscoveryURLs
	enc.Whitelist = c.Whitelist
	enc.StorageMode = c.Prune.ToString()
	enc.Miner = c.Miner
	enc.Ethash = c.Ethash
	enc.TxPool = c.TxPool
	enc.GPO = c.GPO
	enc.DocRoot = c.DocRoot
	enc.RPCGasCap = c.RPCGasCap
	enc.RPCTxFeeCap = c.RPCTxFeeCap
	enc.Checkpoint = c.Checkpoint
	enc.CheckpointOracle = c.CheckpointOracle
	return &enc, nil
}

// UnmarshalTOML unmarshals from TOML.
func (c *Config) UnmarshalTOML(unmarshal func(interface{}) error) error {
	type Config struct {
		Genesis            *core.Genesis `toml:",omitempty"`
		NetworkID          *uint64
		EthDiscoveryURLs   []string
		Whitelist          map[uint64]common.Hash `toml:"-"`
		Mode               *string
		Experiments        *[]string
		OnlyAnnounce       *bool
		SkipBcVersionCheck *bool `toml:"-"`
		DatabaseHandles    *int  `toml:"-"`
		DatabaseFreezer    *string
		Miner              *params.MiningConfig
		Ethash             *ethash.Config
		TxPool             *core.TxPoolConfig
		GPO                *gasprice.Config
		DocRoot            *string                        `toml:"-"`
		RPCGasCap          *uint64                        `toml:",omitempty"`
		RPCTxFeeCap        *float64                       `toml:",omitempty"`
		Checkpoint         *params.TrustedCheckpoint      `toml:",omitempty"`
		CheckpointOracle   *params.CheckpointOracleConfig `toml:",omitempty"`
	}
	var dec Config
	if err := unmarshal(&dec); err != nil {
		return err
	}
	if dec.Genesis != nil {
		c.Genesis = dec.Genesis
	}
	if dec.NetworkID != nil {
		c.NetworkID = *dec.NetworkID
	}
	if dec.EthDiscoveryURLs != nil {
		c.EthDiscoveryURLs = dec.EthDiscoveryURLs
	}
	if dec.Whitelist != nil {
		c.Whitelist = dec.Whitelist
	}
	if dec.Mode != nil {
		mode, err := ethdb.PruneFromString(*dec.Mode, *dec.Experiments)
		if err != nil {
			return err
		}
		c.Prune = mode
	}
	if dec.Miner != nil {
		c.Miner = *dec.Miner
	}
	if dec.Ethash != nil {
		c.Ethash = *dec.Ethash
	}
	if dec.TxPool != nil {
		c.TxPool = *dec.TxPool
	}
	if dec.GPO != nil {
		c.GPO = *dec.GPO
	}
	if dec.DocRoot != nil {
		c.DocRoot = *dec.DocRoot
	}
	if dec.RPCGasCap != nil {
		c.RPCGasCap = *dec.RPCGasCap
	}
	if dec.RPCTxFeeCap != nil {
		c.RPCTxFeeCap = *dec.RPCTxFeeCap
	}
	if dec.Checkpoint != nil {
		c.Checkpoint = dec.Checkpoint
	}
	if dec.CheckpointOracle != nil {
		c.CheckpointOracle = dec.CheckpointOracle
	}
	return nil
}
