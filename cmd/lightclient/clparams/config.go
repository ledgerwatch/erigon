/*
   Copyright 2022 Erigon-Lightclient contributors
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package clparams

import (
	"math"
	"time"

	"github.com/ledgerwatch/erigon/cmd/lightclient/bytesutils"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/crypto"
	"github.com/ledgerwatch/erigon/p2p/discover"
	"github.com/ledgerwatch/erigon/p2p/enode"
)

type NetworkType int

const (
	MainnetNetwork NetworkType = 0
	GoerliNetwork  NetworkType = 1
	SepoliaNetwork NetworkType = 2
)

var (
	MainnetBootstrapNodes = []string{
		// Teku team's bootnode
		"enr:-KG4QOtcP9X1FbIMOe17QNMKqDxCpm14jcX5tiOE4_TyMrFqbmhPZHK_ZPG2Gxb1GE2xdtodOfx9-cgvNtxnRyHEmC0ghGV0aDKQ9aX9QgAAAAD__________4JpZIJ2NIJpcIQDE8KdiXNlY3AyNTZrMaEDhpehBDbZjM_L9ek699Y7vhUJ-eAdMyQW_Fil522Y0fODdGNwgiMog3VkcIIjKA",
		"enr:-KG4QL-eqFoHy0cI31THvtZjpYUu_Jdw_MO7skQRJxY1g5HTN1A0epPCU6vi0gLGUgrzpU-ygeMSS8ewVxDpKfYmxMMGhGV0aDKQtTA_KgAAAAD__________4JpZIJ2NIJpcIQ2_DUbiXNlY3AyNTZrMaED8GJ2vzUqgL6-KD1xalo1CsmY4X1HaDnyl6Y_WayCo9GDdGNwgiMog3VkcIIjKA",
		// Prylab team's bootnodes
		"enr:-Ku4QImhMc1z8yCiNJ1TyUxdcfNucje3BGwEHzodEZUan8PherEo4sF7pPHPSIB1NNuSg5fZy7qFsjmUKs2ea1Whi0EBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQOVphkDqal4QzPMksc5wnpuC3gvSC8AfbFOnZY_On34wIN1ZHCCIyg",
		"enr:-Ku4QP2xDnEtUXIjzJ_DhlCRN9SN99RYQPJL92TMlSv7U5C1YnYLjwOQHgZIUXw6c-BvRg2Yc2QsZxxoS_pPRVe0yK8Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQMeFF5GrS7UZpAH2Ly84aLK-TyvH-dRo0JM1i8yygH50YN1ZHCCJxA",
		"enr:-Ku4QPp9z1W4tAO8Ber_NQierYaOStqhDqQdOPY3bB3jDgkjcbk6YrEnVYIiCBbTxuar3CzS528d2iE7TdJsrL-dEKoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQMw5fqqkw2hHC4F5HZZDPsNmPdB1Gi8JPQK7pRc9XHh-oN1ZHCCKvg",
		// Lighthouse team's bootnodes
		"enr:-Jq4QItoFUuug_n_qbYbU0OY04-np2wT8rUCauOOXNi0H3BWbDj-zbfZb7otA7jZ6flbBpx1LNZK2TDebZ9dEKx84LYBhGV0aDKQtTA_KgEAAAD__________4JpZIJ2NIJpcISsaa0ZiXNlY3AyNTZrMaEDHAD2JKYevx89W0CcFJFiskdcEzkH_Wdv9iW42qLK79ODdWRwgiMo",
		"enr:-Jq4QN_YBsUOqQsty1OGvYv48PMaiEt1AzGD1NkYQHaxZoTyVGqMYXg0K9c0LPNWC9pkXmggApp8nygYLsQwScwAgfgBhGV0aDKQtTA_KgEAAAD__________4JpZIJ2NIJpcISLosQxiXNlY3AyNTZrMaEDBJj7_dLFACaxBfaI8KZTh_SSJUjhyAyfshimvSqo22WDdWRwgiMo",
		// EF bootnodes
		"enr:-Ku4QHqVeJ8PPICcWk1vSn_XcSkjOkNiTg6Fmii5j6vUQgvzMc9L1goFnLKgXqBJspJjIsB91LTOleFmyWWrFVATGngBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhAMRHkWJc2VjcDI1NmsxoQKLVXFOhp2uX6jeT0DvvDpPcU8FWMjQdR4wMuORMhpX24N1ZHCCIyg",
		"enr:-Ku4QG-2_Md3sZIAUebGYT6g0SMskIml77l6yR-M_JXc-UdNHCmHQeOiMLbylPejyJsdAPsTHJyjJB2sYGDLe0dn8uYBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhBLY-NyJc2VjcDI1NmsxoQORcM6e19T1T9gi7jxEZjk_sjVLGFscUNqAY9obgZaxbIN1ZHCCIyg",
		"enr:-Ku4QPn5eVhcoF1opaFEvg1b6JNFD2rqVkHQ8HApOKK61OIcIXD127bKWgAtbwI7pnxx6cDyk_nI88TrZKQaGMZj0q0Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhDayLMaJc2VjcDI1NmsxoQK2sBOLGcUb4AwuYzFuAVCaNHA-dy24UuEKkeFNgCVCsIN1ZHCCIyg",
		"enr:-Ku4QEWzdnVtXc2Q0ZVigfCGggOVB2Vc1ZCPEc6j21NIFLODSJbvNaef1g4PxhPwl_3kax86YPheFUSLXPRs98vvYsoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhDZBrP2Jc2VjcDI1NmsxoQM6jr8Rb1ktLEsVcKAPa08wCsKUmvoQ8khiOl_SLozf9IN1ZHCCIyg",
		// Nimbus bootnodes
		"enr:-LK4QA8FfhaAjlb_BXsXxSfiysR7R52Nhi9JBt4F8SPssu8hdE1BXQQEtVDC3qStCW60LSO7hEsVHv5zm8_6Vnjhcn0Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhAN4aBKJc2VjcDI1NmsxoQJerDhsJ-KxZ8sHySMOCmTO6sHM3iCFQ6VMvLTe948MyYN0Y3CCI4yDdWRwgiOM",
		"enr:-LK4QKWrXTpV9T78hNG6s8AM6IO4XH9kFT91uZtFg1GcsJ6dKovDOr1jtAAFPnS2lvNltkOGA9k29BUN7lFh_sjuc9QBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhANAdd-Jc2VjcDI1NmsxoQLQa6ai7y9PMN5hpLe5HmiJSlYzMuzP7ZhwRiwHvqNXdoN0Y3CCI4yDdWRwgiOM",
	}
	GoerliBootstrapNodes = []string{
		"enr:-Ku4QFmUkNp0g9bsLX2PfVeIyT-9WO-PZlrqZBNtEyofOOfLMScDjaTzGxIb1Ns9Wo5Pm_8nlq-SZwcQfTH2cgO-s88Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpDkvpOTAAAQIP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQLV_jMOIxKbjHFKgrkFvwDvpexo6Nd58TK5k7ss4Vt0IoN1ZHCCG1g",
		"enr:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo",
		"enr:-Ly4QFPk-cTMxZ3jWTafiNblEZkQIXGF2aVzCIGW0uHp6KaEAvBMoctE8S7YU0qZtuS7By0AA4YMfKoN9ls_GJRccVpFh2F0dG5ldHOI__________-EZXRoMpCC9KcrAgAQIIS2AQAAAAAAgmlkgnY0gmlwhKh3joWJc2VjcDI1NmsxoQKrxz8M1IHwJqRIpDqdVW_U1PeixMW5SfnBD-8idYIQrIhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA",
		"enr:-L64QJmwSDtaHVgGiqIxJWUtxWg6uLCipsms6j-8BdsOJfTWAs7CLF9HJnVqFE728O-JYUDCxzKvRdeMqBSauHVCMdaCAVWHYXR0bmV0c4j__________4RldGgykIL0pysCABAghLYBAAAAAACCaWSCdjSCaXCEQWxOdolzZWNwMjU2azGhA7Qmod9fK86WidPOzLsn5_8QyzL7ZcJ1Reca7RnD54vuiHN5bmNuZXRzD4N0Y3CCIyiDdWRwgiMo",
		"enr:-KG4QCIzJZTY_fs_2vqWEatJL9RrtnPwDCv-jRBuO5FQ2qBrfJubWOWazri6s9HsyZdu-fRUfEzkebhf1nvO42_FVzwDhGV0aDKQed8EKAAAECD__________4JpZIJ2NIJpcISHtbYziXNlY3AyNTZrMaED4m9AqVs6F32rSCGsjtYcsyfQE2K8nDiGmocUY_iq-TSDdGNwgiMog3VkcIIjKA",
	}

	SepoliaBootstrapNodes = []string{
		// EF boot nodes
		"enr:-Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk",
		"enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA",
		// Teku boot node
		"enr:-Ly4QFoZTWR8ulxGVsWydTNGdwEESueIdj-wB6UmmjUcm-AOPxnQi7wprzwcdo7-1jBW_JxELlUKJdJES8TDsbl1EdNlh2F0dG5ldHOI__78_v2bsV-EZXRoMpA2-lATkAAAcf__________gmlkgnY0gmlwhBLYJjGJc2VjcDI1NmsxoQI0gujXac9rMAb48NtMqtSTyHIeNYlpjkbYpWJw46PmYYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA",
	}
)

type NetworkConfig struct {
	GossipMaxSize                   uint64        `json:"gossip_max_size"`                    // The maximum allowed size of uncompressed gossip messages.
	MaxRequestBlocks                uint64        `json:"max_request_blocks"`                 // Maximum number of blocks in a single request
	MinEpochsForBlockRequests       uint64        `json:"min_epochs_for_block_requests"`      // The minimum epoch range over which a node must serve blocks
	MaxChunkSize                    uint64        `json:"max_chunk_size"`                     // The maximum allowed size of uncompressed req/resp chunked responses.
	AttestationSubnetCount          uint64        `json:"attestation_subnet_count"`           // The number of attestation subnets used in the gossipsub protocol.
	TtfbTimeout                     time.Duration `json:"ttfbt_timeout"`                      // The maximum time to wait for first byte of request response (time-to-first-byte).
	RespTimeout                     time.Duration `json:"resp_timeout"`                       // The maximum time for complete response transfer.
	AttestationPropagationSlotRange uint64        `json:"attestation_propagation_slot_range"` // The maximum number of slots during which an attestation can be propagated.
	MaximumGossipClockDisparity     time.Duration `json:"maximum_gossip_clock_disparity"`     // The maximum milliseconds of clock disparity assumed between honest nodes.
	MessageDomainInvalidSnappy      [4]byte       `json:"message_domain_invalid_snappy"`      // 4-byte domain for gossip message-id isolation of invalid snappy messages
	MessageDomainValidSnappy        [4]byte       `json:"message_domain_valid_snappy"`        // 4-byte domain for gossip message-id isolation of valid snappy messages

	// DiscoveryV5 Config
	Eth2key                    string // ETH2Key is the ENR key of the Ethereum consensus object in an enr.
	AttSubnetKey               string // AttSubnetKey is the ENR key of the subnet bitfield in the enr.
	SyncCommsSubnetKey         string // SyncCommsSubnetKey is the ENR key of the sync committee subnet bitfield in the enr.
	MinimumPeersInSubnetSearch uint64 // PeersInSubnetSearch is the required amount of peers that we need to be able to lookup in a subnet search.

	ContractDeploymentBlock uint64 // the eth1 block in which the deposit contract is deployed.
	bootNodes               []string
}

type GenesisConfig struct {
	GenesisValidatorRoot common.Hash // Merkle Root at Genesis
	GenesisTime          uint64      // Unix time at Genesis
}

var NetworkConfigs map[NetworkType]NetworkConfig = map[NetworkType]NetworkConfig{
	MainnetNetwork: {
		GossipMaxSize:                   1 << 20, // 1 MiB
		MaxChunkSize:                    1 << 20, // 1 MiB
		AttestationSubnetCount:          64,
		AttestationPropagationSlotRange: 32,
		MaxRequestBlocks:                1 << 10, // 1024
		TtfbTimeout:                     5 * time.Second,
		RespTimeout:                     10 * time.Second,
		MaximumGossipClockDisparity:     500 * time.Millisecond,
		MessageDomainInvalidSnappy:      [4]byte{00, 00, 00, 00},
		MessageDomainValidSnappy:        [4]byte{01, 00, 00, 00},
		Eth2key:                         "eth2",
		AttSubnetKey:                    "attnets",
		SyncCommsSubnetKey:              "syncnets",
		MinimumPeersInSubnetSearch:      20,
		ContractDeploymentBlock:         11184524,
		bootNodes:                       MainnetBootstrapNodes,
	},

	SepoliaNetwork: {
		GossipMaxSize:                   1 << 20, // 1 MiB
		MaxChunkSize:                    1 << 20, // 1 MiB
		AttestationSubnetCount:          64,
		AttestationPropagationSlotRange: 32,
		MaxRequestBlocks:                1 << 10, // 1024
		TtfbTimeout:                     5 * time.Second,
		RespTimeout:                     10 * time.Second,
		MaximumGossipClockDisparity:     500 * time.Millisecond,
		MessageDomainInvalidSnappy:      [4]byte{00, 00, 00, 00},
		MessageDomainValidSnappy:        [4]byte{01, 00, 00, 00},
		Eth2key:                         "eth2",
		AttSubnetKey:                    "attnets",
		SyncCommsSubnetKey:              "syncnets",
		MinimumPeersInSubnetSearch:      20,
		ContractDeploymentBlock:         1273020,
		bootNodes:                       SepoliaBootstrapNodes,
	},

	GoerliNetwork: {
		GossipMaxSize:                   1 << 20, // 1 MiB
		MaxChunkSize:                    1 << 20, // 1 MiB
		AttestationSubnetCount:          64,
		AttestationPropagationSlotRange: 32,
		MaxRequestBlocks:                1 << 10, // 1024
		TtfbTimeout:                     5 * time.Second,
		RespTimeout:                     10 * time.Second,
		MaximumGossipClockDisparity:     500 * time.Millisecond,
		MessageDomainInvalidSnappy:      [4]byte{00, 00, 00, 00},
		MessageDomainValidSnappy:        [4]byte{01, 00, 00, 00},
		Eth2key:                         "eth2",
		AttSubnetKey:                    "attnets",
		SyncCommsSubnetKey:              "syncnets",
		MinimumPeersInSubnetSearch:      20,
		ContractDeploymentBlock:         4367322,
		bootNodes:                       GoerliBootstrapNodes,
	},
}

var GenesisConfigs map[NetworkType]GenesisConfig = map[NetworkType]GenesisConfig{
	MainnetNetwork: {
		GenesisValidatorRoot: common.HexToHash("4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95"),
		GenesisTime:          1606824023,
	},
	SepoliaNetwork: {
		GenesisValidatorRoot: common.HexToHash("d8ea171f3c94aea21ebc42a1ed61052acf3f9209c00e4efbaaddac09ed9b8078"),
		GenesisTime:          1655733600,
	},
	GoerliNetwork: {
		GenesisValidatorRoot: common.HexToHash("043db0d9a83813551ee2f33450d23797757d430911a9320530ad8a0eabc43efb"),
		GenesisTime:          1616508000,
	},
}

// BeaconChainConfig contains constant configs for node to participate in beacon chain.
type BeaconChainConfig struct {
	// Constants (non-configurable)
	GenesisSlot              uint64 `yaml:"GENESIS_SLOT"`                // GenesisSlot represents the first canonical slot number of the beacon chain.
	GenesisEpoch             uint64 `yaml:"GENESIS_EPOCH"`               // GenesisEpoch represents the first canonical epoch number of the beacon chain.
	FarFutureEpoch           uint64 `yaml:"FAR_FUTURE_EPOCH"`            // FarFutureEpoch represents a epoch extremely far away in the future used as the default penalization epoch for validators.
	FarFutureSlot            uint64 `yaml:"FAR_FUTURE_SLOT"`             // FarFutureSlot represents a slot extremely far away in the future.
	BaseRewardsPerEpoch      uint64 `yaml:"BASE_REWARDS_PER_EPOCH"`      // BaseRewardsPerEpoch is used to calculate the per epoch rewards.
	DepositContractTreeDepth uint64 `yaml:"DEPOSIT_CONTRACT_TREE_DEPTH"` // DepositContractTreeDepth depth of the Merkle trie of deposits in the validator deposit contract on the PoW chain.
	JustificationBitsLength  uint64 `yaml:"JUSTIFICATION_BITS_LENGTH"`   // JustificationBitsLength defines number of epochs to track when implementing k-finality in Casper FFG.

	// Misc constants.
	PresetBase                     string `yaml:"PRESET_BASE" spec:"true"`                        // PresetBase represents the underlying spec preset this config is based on.
	ConfigName                     string `yaml:"CONFIG_NAME" spec:"true"`                        // ConfigName for allowing an easy human-readable way of knowing what chain is being used.
	TargetCommitteeSize            uint64 `yaml:"TARGET_COMMITTEE_SIZE" spec:"true"`              // TargetCommitteeSize is the number of validators in a committee when the chain is healthy.
	MaxValidatorsPerCommittee      uint64 `yaml:"MAX_VALIDATORS_PER_COMMITTEE" spec:"true"`       // MaxValidatorsPerCommittee defines the upper bound of the size of a committee.
	MaxCommitteesPerSlot           uint64 `yaml:"MAX_COMMITTEES_PER_SLOT" spec:"true"`            // MaxCommitteesPerSlot defines the max amount of committee in a single slot.
	MinPerEpochChurnLimit          uint64 `yaml:"MIN_PER_EPOCH_CHURN_LIMIT" spec:"true"`          // MinPerEpochChurnLimit is the minimum amount of churn allotted for validator rotations.
	ChurnLimitQuotient             uint64 `yaml:"CHURN_LIMIT_QUOTIENT" spec:"true"`               // ChurnLimitQuotient is used to determine the limit of how many validators can rotate per epoch.
	ShuffleRoundCount              uint64 `yaml:"SHUFFLE_ROUND_COUNT" spec:"true"`                // ShuffleRoundCount is used for retrieving the permuted index.
	MinGenesisActiveValidatorCount uint64 `yaml:"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT" spec:"true"` // MinGenesisActiveValidatorCount defines how many validator deposits needed to kick off beacon chain.
	MinGenesisTime                 uint64 `yaml:"MIN_GENESIS_TIME" spec:"true"`                   // MinGenesisTime is the time that needed to pass before kicking off beacon chain.
	TargetAggregatorsPerCommittee  uint64 `yaml:"TARGET_AGGREGATORS_PER_COMMITTEE" spec:"true"`   // TargetAggregatorsPerCommittee defines the number of aggregators inside one committee.
	HysteresisQuotient             uint64 `yaml:"HYSTERESIS_QUOTIENT" spec:"true"`                // HysteresisQuotient defines the hysteresis quotient for effective balance calculations.
	HysteresisDownwardMultiplier   uint64 `yaml:"HYSTERESIS_DOWNWARD_MULTIPLIER" spec:"true"`     // HysteresisDownwardMultiplier defines the hysteresis downward multiplier for effective balance calculations.
	HysteresisUpwardMultiplier     uint64 `yaml:"HYSTERESIS_UPWARD_MULTIPLIER" spec:"true"`       // HysteresisUpwardMultiplier defines the hysteresis upward multiplier for effective balance calculations.

	// Gwei value constants.
	MinDepositAmount          uint64 `yaml:"MIN_DEPOSIT_AMOUNT" spec:"true"`          // MinDepositAmount is the minimum amount of Gwei a validator can send to the deposit contract at once (lower amounts will be reverted).
	MaxEffectiveBalance       uint64 `yaml:"MAX_EFFECTIVE_BALANCE" spec:"true"`       // MaxEffectiveBalance is the maximal amount of Gwei that is effective for staking.
	EjectionBalance           uint64 `yaml:"EJECTION_BALANCE" spec:"true"`            // EjectionBalance is the minimal GWei a validator needs to have before ejected.
	EffectiveBalanceIncrement uint64 `yaml:"EFFECTIVE_BALANCE_INCREMENT" spec:"true"` // EffectiveBalanceIncrement is used for converting the high balance into the low balance for validators.

	// Initial value constants.
	BLSWithdrawalPrefixByte         byte     `yaml:"BLS_WITHDRAWAL_PREFIX" spec:"true"`          // BLSWithdrawalPrefixByte is used for BLS withdrawal and it's the first byte.
	ETH1AddressWithdrawalPrefixByte byte     `yaml:"ETH1_ADDRESS_WITHDRAWAL_PREFIX" spec:"true"` // ETH1AddressWithdrawalPrefixByte is used for withdrawals and it's the first byte.
	ZeroHash                        [32]byte // ZeroHash is used to represent a zeroed out 32 byte array.

	// Time parameters constants.
	GenesisDelay                              uint64 `yaml:"GENESIS_DELAY" spec:"true"`                   // GenesisDelay is the minimum number of seconds to delay starting the Ethereum Beacon Chain genesis. Must be at least 1 second.
	MinAttestationInclusionDelay              uint64 `yaml:"MIN_ATTESTATION_INCLUSION_DELAY" spec:"true"` // MinAttestationInclusionDelay defines how many slots validator has to wait to include attestation for beacon block.
	SecondsPerSlot                            uint64 `yaml:"SECONDS_PER_SLOT" spec:"true"`                // SecondsPerSlot is how many seconds are in a single slot.
	SlotsPerEpoch                             uint64 `yaml:"SLOTS_PER_EPOCH" spec:"true"`                 // SlotsPerEpoch is the number of slots in an epoch.
	SqrRootSlotsPerEpoch                      uint64 // SqrRootSlotsPerEpoch is a hard coded value where we take the square root of `SlotsPerEpoch` and round down.
	MinSeedLookahead                          uint64 `yaml:"MIN_SEED_LOOKAHEAD" spec:"true"`                  // MinSeedLookahead is the duration of randao look ahead seed.
	MaxSeedLookahead                          uint64 `yaml:"MAX_SEED_LOOKAHEAD" spec:"true"`                  // MaxSeedLookahead is the duration a validator has to wait for entry and exit in epoch.
	EpochsPerEth1VotingPeriod                 uint64 `yaml:"EPOCHS_PER_ETH1_VOTING_PERIOD" spec:"true"`       // EpochsPerEth1VotingPeriod defines how often the merkle root of deposit receipts get updated in beacon node on per epoch basis.
	SlotsPerHistoricalRoot                    uint64 `yaml:"SLOTS_PER_HISTORICAL_ROOT" spec:"true"`           // SlotsPerHistoricalRoot defines how often the historical root is saved.
	MinValidatorWithdrawabilityDelay          uint64 `yaml:"MIN_VALIDATOR_WITHDRAWABILITY_DELAY" spec:"true"` // MinValidatorWithdrawabilityDelay is the shortest amount of time a validator has to wait to withdraw.
	ShardCommitteePeriod                      uint64 `yaml:"SHARD_COMMITTEE_PERIOD" spec:"true"`              // ShardCommitteePeriod is the minimum amount of epochs a validator must participate before exiting.
	MinEpochsToInactivityPenalty              uint64 `yaml:"MIN_EPOCHS_TO_INACTIVITY_PENALTY" spec:"true"`    // MinEpochsToInactivityPenalty defines the minimum amount of epochs since finality to begin penalizing inactivity.
	Eth1FollowDistance                        uint64 `yaml:"ETH1_FOLLOW_DISTANCE" spec:"true"`                // Eth1FollowDistance is the number of eth1.0 blocks to wait before considering a new deposit for voting. This only applies after the chain as been started.
	SafeSlotsToUpdateJustified                uint64 `yaml:"SAFE_SLOTS_TO_UPDATE_JUSTIFIED" spec:"true"`      // SafeSlotsToUpdateJustified is the minimal slots needed to update justified check point.
	DeprecatedSafeSlotsToImportOptimistically uint64 `yaml:"SAFE_SLOTS_TO_IMPORT_OPTIMISTICALLY" spec:"true"` // SafeSlotsToImportOptimistically is the minimal number of slots to wait before importing optimistically a pre-merge block
	SecondsPerETH1Block                       uint64 `yaml:"SECONDS_PER_ETH1_BLOCK" spec:"true"`              // SecondsPerETH1Block is the approximate time for a single eth1 block to be produced.

	// Fork choice algorithm constants.
	ProposerScoreBoost uint64 `yaml:"PROPOSER_SCORE_BOOST" spec:"true"` // ProposerScoreBoost defines a value that is a % of the committee weight for fork-choice boosting.
	IntervalsPerSlot   uint64 `yaml:"INTERVALS_PER_SLOT" spec:"true"`   // IntervalsPerSlot defines the number of fork choice intervals in a slot defined in the fork choice spec.

	// Ethereum PoW parameters.
	DepositChainID         uint64 `yaml:"DEPOSIT_CHAIN_ID" spec:"true"`         // DepositChainID of the eth1 network. This used for replay protection.
	DepositNetworkID       uint64 `yaml:"DEPOSIT_NETWORK_ID" spec:"true"`       // DepositNetworkID of the eth1 network. This used for replay protection.
	DepositContractAddress string `yaml:"DEPOSIT_CONTRACT_ADDRESS" spec:"true"` // DepositContractAddress is the address of the deposit contract.

	// Validator parameters.
	RandomSubnetsPerValidator         uint64 `yaml:"RANDOM_SUBNETS_PER_VALIDATOR" spec:"true"`          // RandomSubnetsPerValidator specifies the amount of subnets a validator has to be subscribed to at one time.
	EpochsPerRandomSubnetSubscription uint64 `yaml:"EPOCHS_PER_RANDOM_SUBNET_SUBSCRIPTION" spec:"true"` // EpochsPerRandomSubnetSubscription specifies the minimum duration a validator is connected to their subnet.

	// State list lengths
	EpochsPerHistoricalVector uint64 `yaml:"EPOCHS_PER_HISTORICAL_VECTOR" spec:"true"` // EpochsPerHistoricalVector defines max length in epoch to store old historical stats in beacon state.
	EpochsPerSlashingsVector  uint64 `yaml:"EPOCHS_PER_SLASHINGS_VECTOR" spec:"true"`  // EpochsPerSlashingsVector defines max length in epoch to store old stats to recompute slashing witness.
	HistoricalRootsLimit      uint64 `yaml:"HISTORICAL_ROOTS_LIMIT" spec:"true"`       // HistoricalRootsLimit defines max historical roots that can be saved in state before roll over.
	ValidatorRegistryLimit    uint64 `yaml:"VALIDATOR_REGISTRY_LIMIT" spec:"true"`     // ValidatorRegistryLimit defines the upper bound of validators can participate in eth2.

	// Reward and penalty quotients constants.
	BaseRewardFactor               uint64 `yaml:"BASE_REWARD_FACTOR" spec:"true"`               // BaseRewardFactor is used to calculate validator per-slot interest rate.
	WhistleBlowerRewardQuotient    uint64 `yaml:"WHISTLEBLOWER_REWARD_QUOTIENT" spec:"true"`    // WhistleBlowerRewardQuotient is used to calculate whistle blower reward.
	ProposerRewardQuotient         uint64 `yaml:"PROPOSER_REWARD_QUOTIENT" spec:"true"`         // ProposerRewardQuotient is used to calculate the reward for proposers.
	InactivityPenaltyQuotient      uint64 `yaml:"INACTIVITY_PENALTY_QUOTIENT" spec:"true"`      // InactivityPenaltyQuotient is used to calculate the penalty for a validator that is offline.
	MinSlashingPenaltyQuotient     uint64 `yaml:"MIN_SLASHING_PENALTY_QUOTIENT" spec:"true"`    // MinSlashingPenaltyQuotient is used to calculate the minimum penalty to prevent DoS attacks.
	ProportionalSlashingMultiplier uint64 `yaml:"PROPORTIONAL_SLASHING_MULTIPLIER" spec:"true"` // ProportionalSlashingMultiplier is used as a multiplier on slashed penalties.

	// Max operations per block constants.
	MaxProposerSlashings uint64 `yaml:"MAX_PROPOSER_SLASHINGS" spec:"true"` // MaxProposerSlashings defines the maximum number of slashings of proposers possible in a block.
	MaxAttesterSlashings uint64 `yaml:"MAX_ATTESTER_SLASHINGS" spec:"true"` // MaxAttesterSlashings defines the maximum number of casper FFG slashings possible in a block.
	MaxAttestations      uint64 `yaml:"MAX_ATTESTATIONS" spec:"true"`       // MaxAttestations defines the maximum allowed attestations in a beacon block.
	MaxDeposits          uint64 `yaml:"MAX_DEPOSITS" spec:"true"`           // MaxDeposits defines the maximum number of validator deposits in a block.
	MaxVoluntaryExits    uint64 `yaml:"MAX_VOLUNTARY_EXITS" spec:"true"`    // MaxVoluntaryExits defines the maximum number of validator exits in a block.

	// BLS domain values.
	DomainBeaconProposer              [4]byte `yaml:"DOMAIN_BEACON_PROPOSER" spec:"true"`                // DomainBeaconProposer defines the BLS signature domain for beacon proposal verification.
	DomainRandao                      [4]byte `yaml:"DOMAIN_RANDAO" spec:"true"`                         // DomainRandao defines the BLS signature domain for randao verification.
	DomainBeaconAttester              [4]byte `yaml:"DOMAIN_BEACON_ATTESTER" spec:"true"`                // DomainBeaconAttester defines the BLS signature domain for attestation verification.
	DomainDeposit                     [4]byte `yaml:"DOMAIN_DEPOSIT" spec:"true"`                        // DomainDeposit defines the BLS signature domain for deposit verification.
	DomainVoluntaryExit               [4]byte `yaml:"DOMAIN_VOLUNTARY_EXIT" spec:"true"`                 // DomainVoluntaryExit defines the BLS signature domain for exit verification.
	DomainSelectionProof              [4]byte `yaml:"DOMAIN_SELECTION_PROOF" spec:"true"`                // DomainSelectionProof defines the BLS signature domain for selection proof.
	DomainAggregateAndProof           [4]byte `yaml:"DOMAIN_AGGREGATE_AND_PROOF" spec:"true"`            // DomainAggregateAndProof defines the BLS signature domain for aggregate and proof.
	DomainSyncCommittee               [4]byte `yaml:"DOMAIN_SYNC_COMMITTEE" spec:"true"`                 // DomainVoluntaryExit defines the BLS signature domain for sync committee.
	DomainSyncCommitteeSelectionProof [4]byte `yaml:"DOMAIN_SYNC_COMMITTEE_SELECTION_PROOF" spec:"true"` // DomainSelectionProof defines the BLS signature domain for sync committee selection proof.
	DomainContributionAndProof        [4]byte `yaml:"DOMAIN_CONTRIBUTION_AND_PROOF" spec:"true"`         // DomainAggregateAndProof defines the BLS signature domain for contribution and proof.
	DomainApplicationMask             [4]byte `yaml:"DOMAIN_APPLICATION_MASK" spec:"true"`               // DomainApplicationMask defines the BLS signature domain for application mask.
	DomainApplicationBuilder          [4]byte // DomainApplicationBuilder defines the BLS signature domain for application builder.
	DomainBLSToExecutionChange        [4]byte // DomainBLSToExecutionChange defines the BLS signature domain to change withdrawal addresses to ETH1 prefix

	// Prysm constants.
	GweiPerEth                     uint64        // GweiPerEth is the amount of gwei corresponding to 1 eth.
	BLSSecretKeyLength             int           // BLSSecretKeyLength defines the expected length of BLS secret keys in bytes.
	BLSPubkeyLength                int           // BLSPubkeyLength defines the expected length of BLS public keys in bytes.
	DefaultBufferSize              int           // DefaultBufferSize for channels across the Prysm repository.
	ValidatorPrivkeyFileName       string        // ValidatorPrivKeyFileName specifies the string name of a validator private key file.
	WithdrawalPrivkeyFileName      string        // WithdrawalPrivKeyFileName specifies the string name of a withdrawal private key file.
	RPCSyncCheck                   time.Duration // Number of seconds to query the sync service, to find out if the node is synced or not.
	EmptySignature                 [96]byte      // EmptySignature is used to represent a zeroed out BLS Signature.
	DefaultPageSize                int           // DefaultPageSize defines the default page size for RPC server request.
	MaxPeersToSync                 int           // MaxPeersToSync describes the limit for number of peers in round robin sync.
	SlotsPerArchivedPoint          uint64        // SlotsPerArchivedPoint defines the number of slots per one archived point.
	GenesisCountdownInterval       time.Duration // How often to log the countdown until the genesis time is reached.
	BeaconStateFieldCount          int           // BeaconStateFieldCount defines how many fields are in beacon state.
	BeaconStateAltairFieldCount    int           // BeaconStateAltairFieldCount defines how many fields are in beacon state hard fork 1.
	BeaconStateBellatrixFieldCount int           // BeaconStateBellatrixFieldCount defines how many fields are in beacon state post upgrade to the Bellatrix.

	// Slasher constants.
	WeakSubjectivityPeriod    uint64 // WeakSubjectivityPeriod defines the time period expressed in number of epochs were proof of stake network should validate block headers and attestations for slashable events.
	PruneSlasherStoragePeriod uint64 // PruneSlasherStoragePeriod defines the time period expressed in number of epochs were proof of stake network should prune attestation and block header store.

	// Slashing protection constants.
	SlashingProtectionPruningEpochs uint64 // SlashingProtectionPruningEpochs defines a period after which all prior epochs are pruned in the validator database.

	// Fork-related values.
	GenesisForkVersion   []byte `yaml:"GENESIS_FORK_VERSION" spec:"true"`   // GenesisForkVersion is used to track fork version between state transitions.
	AltairForkVersion    []byte `yaml:"ALTAIR_FORK_VERSION" spec:"true"`    // AltairForkVersion is used to represent the fork version for altair.
	AltairForkEpoch      uint64 `yaml:"ALTAIR_FORK_EPOCH" spec:"true"`      // AltairForkEpoch is used to represent the assigned fork epoch for altair.
	BellatrixForkVersion []byte `yaml:"BELLATRIX_FORK_VERSION" spec:"true"` // BellatrixForkVersion is used to represent the fork version for bellatrix.
	BellatrixForkEpoch   uint64 `yaml:"BELLATRIX_FORK_EPOCH" spec:"true"`   // BellatrixForkEpoch is used to represent the assigned fork epoch for bellatrix.
	ShardingForkVersion  []byte `yaml:"SHARDING_FORK_VERSION" spec:"true"`  // ShardingForkVersion is used to represent the fork version for sharding.
	ShardingForkEpoch    uint64 `yaml:"SHARDING_FORK_EPOCH" spec:"true"`    // ShardingForkEpoch is used to represent the assigned fork epoch for sharding.
	CapellaForkVersion   []byte `yaml:"CAPELLA_FORK_VERSION" spec:"true"`   // CapellaForkVersion is used to represent the fork version for capella.
	CapellaForkEpoch     uint64 `yaml:"CAPELLA_FORK_EPOCH" spec:"true"`     // CapellaForkEpoch is used to represent the assigned fork epoch for capella.

	ForkVersionSchedule map[[VersionLength]byte]uint64 // Schedule of fork epochs by version.
	ForkVersionNames    map[[VersionLength]byte]string // Human-readable names of fork versions.

	// Weak subjectivity values.
	SafetyDecay uint64 // SafetyDecay is defined as the loss in the 1/3 consensus safety margin of the casper FFG mechanism.

	// New values introduced in Altair hard fork 1.
	// Participation flag indices.
	TimelySourceFlagIndex uint8 `yaml:"TIMELY_SOURCE_FLAG_INDEX" spec:"true"` // TimelySourceFlagIndex is the source flag position of the participation bits.
	TimelyTargetFlagIndex uint8 `yaml:"TIMELY_TARGET_FLAG_INDEX" spec:"true"` // TimelyTargetFlagIndex is the target flag position of the participation bits.
	TimelyHeadFlagIndex   uint8 `yaml:"TIMELY_HEAD_FLAG_INDEX" spec:"true"`   // TimelyHeadFlagIndex is the head flag position of the participation bits.

	// Incentivization weights.
	TimelySourceWeight uint64 `yaml:"TIMELY_SOURCE_WEIGHT" spec:"true"` // TimelySourceWeight is the factor of how much source rewards receives.
	TimelyTargetWeight uint64 `yaml:"TIMELY_TARGET_WEIGHT" spec:"true"` // TimelyTargetWeight is the factor of how much target rewards receives.
	TimelyHeadWeight   uint64 `yaml:"TIMELY_HEAD_WEIGHT" spec:"true"`   // TimelyHeadWeight is the factor of how much head rewards receives.
	SyncRewardWeight   uint64 `yaml:"SYNC_REWARD_WEIGHT" spec:"true"`   // SyncRewardWeight is the factor of how much sync committee rewards receives.
	WeightDenominator  uint64 `yaml:"WEIGHT_DENOMINATOR" spec:"true"`   // WeightDenominator accounts for total rewards denomination.
	ProposerWeight     uint64 `yaml:"PROPOSER_WEIGHT" spec:"true"`      // ProposerWeight is the factor of how much proposer rewards receives.

	// Validator related.
	TargetAggregatorsPerSyncSubcommittee uint64 `yaml:"TARGET_AGGREGATORS_PER_SYNC_SUBCOMMITTEE" spec:"true"` // TargetAggregatorsPerSyncSubcommittee for aggregating in sync committee.
	SyncCommitteeSubnetCount             uint64 `yaml:"SYNC_COMMITTEE_SUBNET_COUNT" spec:"true"`              // SyncCommitteeSubnetCount for sync committee subnet count.

	// Misc.
	SyncCommitteeSize            uint64 `yaml:"SYNC_COMMITTEE_SIZE" spec:"true"`              // SyncCommitteeSize for light client sync committee size.
	InactivityScoreBias          uint64 `yaml:"INACTIVITY_SCORE_BIAS" spec:"true"`            // InactivityScoreBias for calculating score bias penalties during inactivity
	InactivityScoreRecoveryRate  uint64 `yaml:"INACTIVITY_SCORE_RECOVERY_RATE" spec:"true"`   // InactivityScoreRecoveryRate for recovering score bias penalties during inactivity.
	EpochsPerSyncCommitteePeriod uint64 `yaml:"EPOCHS_PER_SYNC_COMMITTEE_PERIOD" spec:"true"` // EpochsPerSyncCommitteePeriod defines how many epochs per sync committee period.

	// Updated penalty values. This moves penalty parameters toward their final, maximum security values.
	// Note: We do not override previous configuration values but instead creates new values and replaces usage throughout.
	InactivityPenaltyQuotientAltair         uint64 `yaml:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR" spec:"true"`         // InactivityPenaltyQuotientAltair for penalties during inactivity post Altair hard fork.
	MinSlashingPenaltyQuotientAltair        uint64 `yaml:"MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR" spec:"true"`       // MinSlashingPenaltyQuotientAltair for slashing penalties post Altair hard fork.
	ProportionalSlashingMultiplierAltair    uint64 `yaml:"PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR" spec:"true"`    // ProportionalSlashingMultiplierAltair for slashing penalties' multiplier post Alair hard fork.
	MinSlashingPenaltyQuotientBellatrix     uint64 `yaml:"MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX" spec:"true"`    // MinSlashingPenaltyQuotientBellatrix for slashing penalties post Bellatrix hard fork.
	ProportionalSlashingMultiplierBellatrix uint64 `yaml:"PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX" spec:"true"` // ProportionalSlashingMultiplierBellatrix for slashing penalties' multiplier post Bellatrix hard fork.
	InactivityPenaltyQuotientBellatrix      uint64 `yaml:"INACTIVITY_PENALTY_QUOTIENT_BELLATRIX" spec:"true"`      // InactivityPenaltyQuotientBellatrix for penalties during inactivity post Bellatrix hard fork.

	// Light client
	MinSyncCommitteeParticipants uint64 `yaml:"MIN_SYNC_COMMITTEE_PARTICIPANTS" spec:"true"` // MinSyncCommitteeParticipants defines the minimum amount of sync committee participants for which the light client acknowledges the signature.

	// Bellatrix
	TerminalBlockHash                common.Hash    `yaml:"TERMINAL_BLOCK_HASH" spec:"true"`                  // TerminalBlockHash of beacon chain.
	TerminalBlockHashActivationEpoch uint64         `yaml:"TERMINAL_BLOCK_HASH_ACTIVATION_EPOCH" spec:"true"` // TerminalBlockHashActivationEpoch of beacon chain.
	TerminalTotalDifficulty          string         `yaml:"TERMINAL_TOTAL_DIFFICULTY" spec:"true"`            // TerminalTotalDifficulty is part of the experimental Bellatrix spec. This value is type is currently TBD.
	DefaultFeeRecipient              common.Address // DefaultFeeRecipient where the transaction fee goes to.
	EthBurnAddressHex                string         // EthBurnAddressHex is the constant eth address written in hex format to burn fees in that network. the default is 0x0
	DefaultBuilderGasLimit           uint64         // DefaultBuilderGasLimit is the default used to set the gaslimit for the Builder APIs, typically at around 30M wei.

	// Mev-boost circuit breaker
	MaxBuilderConsecutiveMissedSlots uint64 // MaxBuilderConsecutiveMissedSlots defines the number of consecutive skip slot to fallback from using relay/builder to local execution engine for block construction.
	MaxBuilderEpochMissedSlots       uint64 // MaxBuilderEpochMissedSlots is defines the number of total skip slot (per epoch rolling windows) to fallback from using relay/builder to local execution engine for block construction.
}

// InitializeForkSchedule initializes the schedules forks baked into the config.
func (b *BeaconChainConfig) InitializeForkSchedule() {
	b.ForkVersionSchedule = configForkSchedule(b)
	b.ForkVersionNames = configForkNames(b)
}

func toBytes4(in []byte) (ret [4]byte) {
	copy(ret[:], in)
	return
}

func configForkSchedule(b *BeaconChainConfig) map[[VersionLength]byte]uint64 {
	fvs := map[[VersionLength]byte]uint64{}
	fvs[toBytes4(b.GenesisForkVersion)] = b.GenesisEpoch
	fvs[toBytes4(b.AltairForkVersion)] = b.AltairForkEpoch
	fvs[toBytes4(b.BellatrixForkVersion)] = b.BellatrixForkEpoch
	return fvs
}

func configForkNames(b *BeaconChainConfig) map[[VersionLength]byte]string {
	fvn := map[[VersionLength]byte]string{}
	fvn[toBytes4(b.GenesisForkVersion)] = "phase0"
	fvn[toBytes4(b.AltairForkVersion)] = "altair"
	fvn[toBytes4(b.BellatrixForkVersion)] = "bellatrix"
	return fvn
}

var mainnetBeaconConfig BeaconChainConfig = BeaconChainConfig{
	// Constants (Non-configurable)
	FarFutureEpoch:           math.MaxUint64,
	FarFutureSlot:            math.MaxUint64,
	BaseRewardsPerEpoch:      4,
	DepositContractTreeDepth: 32,
	GenesisDelay:             604800, // 1 week.

	// Misc constant.
	TargetCommitteeSize:            128,
	MaxValidatorsPerCommittee:      2048,
	MaxCommitteesPerSlot:           64,
	MinPerEpochChurnLimit:          4,
	ChurnLimitQuotient:             1 << 16,
	ShuffleRoundCount:              90,
	MinGenesisActiveValidatorCount: 16384,
	MinGenesisTime:                 1606824000, // Dec 1, 2020, 12pm UTC.
	TargetAggregatorsPerCommittee:  16,
	HysteresisQuotient:             4,
	HysteresisDownwardMultiplier:   1,
	HysteresisUpwardMultiplier:     5,

	// Gwei value constants.
	MinDepositAmount:          1 * 1e9,
	MaxEffectiveBalance:       32 * 1e9,
	EjectionBalance:           16 * 1e9,
	EffectiveBalanceIncrement: 1 * 1e9,

	// Initial value constants.
	BLSWithdrawalPrefixByte:         byte(0),
	ETH1AddressWithdrawalPrefixByte: byte(1),
	ZeroHash:                        [32]byte{},

	// Time parameter constants.
	MinAttestationInclusionDelay:     1,
	SecondsPerSlot:                   12,
	SlotsPerEpoch:                    32,
	SqrRootSlotsPerEpoch:             5,
	MinSeedLookahead:                 1,
	MaxSeedLookahead:                 4,
	EpochsPerEth1VotingPeriod:        64,
	SlotsPerHistoricalRoot:           8192,
	MinValidatorWithdrawabilityDelay: 256,
	ShardCommitteePeriod:             256,
	MinEpochsToInactivityPenalty:     4,
	Eth1FollowDistance:               2048,
	SafeSlotsToUpdateJustified:       8,

	// Fork choice algorithm constants.
	ProposerScoreBoost: 40,
	IntervalsPerSlot:   3,

	// Ethereum PoW parameters.
	DepositChainID:         1, // Chain ID of eth1 mainnet.
	DepositNetworkID:       1, // Network ID of eth1 mainnet.
	DepositContractAddress: "0x00000000219ab540356cBB839Cbe05303d7705Fa",

	// Validator params.
	RandomSubnetsPerValidator:         1 << 0,
	EpochsPerRandomSubnetSubscription: 1 << 8,

	// While eth1 mainnet block times are closer to 13s, we must conform with other clients in
	// order to vote on the correct eth1 blocks.
	//
	// Additional context: https://github.com/ethereum/consensus-specs/issues/2132
	// Bug prompting this change: https://github.com/prysmaticlabs/prysm/issues/7856
	// Future optimization: https://github.com/prysmaticlabs/prysm/issues/7739
	SecondsPerETH1Block: 14,

	// State list length constants.
	EpochsPerHistoricalVector: 65536,
	EpochsPerSlashingsVector:  8192,
	HistoricalRootsLimit:      16777216,
	ValidatorRegistryLimit:    1099511627776,

	// Reward and penalty quotients constants.
	BaseRewardFactor:               64,
	WhistleBlowerRewardQuotient:    512,
	ProposerRewardQuotient:         8,
	InactivityPenaltyQuotient:      67108864,
	MinSlashingPenaltyQuotient:     128,
	ProportionalSlashingMultiplier: 1,

	// Max operations per block constants.
	MaxProposerSlashings: 16,
	MaxAttesterSlashings: 2,
	MaxAttestations:      128,
	MaxDeposits:          16,
	MaxVoluntaryExits:    16,

	// BLS domain values.
	DomainBeaconProposer:              bytesutils.Uint32ToBytes4(0x00000000),
	DomainBeaconAttester:              bytesutils.Uint32ToBytes4(0x01000000),
	DomainRandao:                      bytesutils.Uint32ToBytes4(0x02000000),
	DomainDeposit:                     bytesutils.Uint32ToBytes4(0x03000000),
	DomainVoluntaryExit:               bytesutils.Uint32ToBytes4(0x04000000),
	DomainSelectionProof:              bytesutils.Uint32ToBytes4(0x05000000),
	DomainAggregateAndProof:           bytesutils.Uint32ToBytes4(0x06000000),
	DomainSyncCommittee:               bytesutils.Uint32ToBytes4(0x07000000),
	DomainSyncCommitteeSelectionProof: bytesutils.Uint32ToBytes4(0x08000000),
	DomainContributionAndProof:        bytesutils.Uint32ToBytes4(0x09000000),
	DomainApplicationMask:             bytesutils.Uint32ToBytes4(0x00000001),
	DomainApplicationBuilder:          bytesutils.Uint32ToBytes4(0x00000001),
	DomainBLSToExecutionChange:        bytesutils.Uint32ToBytes4(0x0A000000),

	// Prysm constants.
	GweiPerEth:                     1000000000,
	BLSSecretKeyLength:             32,
	BLSPubkeyLength:                48,
	DefaultBufferSize:              10000,
	WithdrawalPrivkeyFileName:      "/shardwithdrawalkey",
	ValidatorPrivkeyFileName:       "/validatorprivatekey",
	RPCSyncCheck:                   1,
	EmptySignature:                 [96]byte{},
	DefaultPageSize:                250,
	MaxPeersToSync:                 15,
	SlotsPerArchivedPoint:          2048,
	GenesisCountdownInterval:       time.Minute,
	ConfigName:                     "mainnet",
	PresetBase:                     "mainnet",
	BeaconStateFieldCount:          21,
	BeaconStateAltairFieldCount:    24,
	BeaconStateBellatrixFieldCount: 25,

	// Slasher related values.
	WeakSubjectivityPeriod:          54000,
	PruneSlasherStoragePeriod:       10,
	SlashingProtectionPruningEpochs: 512,

	// Weak subjectivity values.
	SafetyDecay: 10,

	// Fork related values.
	GenesisEpoch:         0,
	GenesisForkVersion:   []byte{0, 0, 0, 0},
	AltairForkVersion:    []byte{1, 0, 0, 0},
	AltairForkEpoch:      74240,
	BellatrixForkVersion: []byte{2, 0, 0, 0},
	BellatrixForkEpoch:   144869,
	CapellaForkVersion:   []byte{3, 0, 0, 0},
	CapellaForkEpoch:     math.MaxUint64,
	ShardingForkVersion:  []byte{4, 0, 0, 0},
	ShardingForkEpoch:    math.MaxUint64,

	// New values introduced in Altair hard fork 1.
	// Participation flag indices.
	TimelySourceFlagIndex: 0,
	TimelyTargetFlagIndex: 1,
	TimelyHeadFlagIndex:   2,

	// Incentivization weight values.
	TimelySourceWeight: 14,
	TimelyTargetWeight: 26,
	TimelyHeadWeight:   14,
	SyncRewardWeight:   2,
	ProposerWeight:     8,
	WeightDenominator:  64,

	// Validator related values.
	TargetAggregatorsPerSyncSubcommittee: 16,
	SyncCommitteeSubnetCount:             4,

	// Misc values.
	SyncCommitteeSize:            512,
	InactivityScoreBias:          4,
	InactivityScoreRecoveryRate:  16,
	EpochsPerSyncCommitteePeriod: 256,

	// Updated penalty values.
	InactivityPenaltyQuotientAltair:         3 * 1 << 24, //50331648
	MinSlashingPenaltyQuotientAltair:        64,
	ProportionalSlashingMultiplierAltair:    2,
	MinSlashingPenaltyQuotientBellatrix:     32,
	ProportionalSlashingMultiplierBellatrix: 3,
	InactivityPenaltyQuotientBellatrix:      1 << 24,

	// Light client
	MinSyncCommitteeParticipants: 1,

	// Bellatrix
	TerminalBlockHashActivationEpoch: 18446744073709551615,
	TerminalBlockHash:                [32]byte{},
	TerminalTotalDifficulty:          "58750000000000000000000", // Estimated: Sept 15, 2022
	EthBurnAddressHex:                "0x0000000000000000000000000000000000000000",
	DefaultBuilderGasLimit:           uint64(30000000),

	// Mevboost circuit breaker
	MaxBuilderConsecutiveMissedSlots: 3,
	MaxBuilderEpochMissedSlots:       8,
}

func mainnetConfig() BeaconChainConfig {
	cfg := mainnetBeaconConfig
	cfg.InitializeForkSchedule()
	return cfg
}

func sepoliaConfig() BeaconChainConfig {
	cfg := mainnetBeaconConfig
	cfg.MinGenesisTime = 1655647200
	cfg.GenesisDelay = 86400
	cfg.MinGenesisActiveValidatorCount = 1300
	cfg.ConfigName = "sepolia"
	cfg.GenesisForkVersion = []byte{0x90, 0x00, 0x00, 0x69}
	cfg.SecondsPerETH1Block = 14
	cfg.DepositChainID = SepoliaEth1ChainId
	cfg.DepositNetworkID = SepoliaEth1ChainId
	cfg.AltairForkEpoch = 50
	cfg.AltairForkVersion = []byte{0x90, 0x00, 0x00, 0x70}
	cfg.BellatrixForkEpoch = 100
	cfg.BellatrixForkVersion = []byte{0x90, 0x00, 0x00, 0x71}
	cfg.TerminalTotalDifficulty = "17000000000000000"
	cfg.DepositContractAddress = "0x7f02C3E3c98b133055B8B348B2Ac625669Ed295D"
	cfg.InitializeForkSchedule()
	return cfg
}

func goerliConfig() BeaconChainConfig {
	cfg := mainnetBeaconConfig
	cfg.MinGenesisTime = 1614588812
	cfg.GenesisDelay = 1919188
	cfg.ConfigName = "prater"
	cfg.GenesisForkVersion = []byte{0x00, 0x00, 0x10, 0x20}
	cfg.SecondsPerETH1Block = 14
	cfg.DepositChainID = GoerliEth1ChainId
	cfg.DepositNetworkID = GoerliEth1ChainId
	cfg.AltairForkEpoch = 36660
	cfg.AltairForkVersion = []byte{0x1, 0x0, 0x10, 0x20}
	cfg.CapellaForkVersion = []byte{0x3, 0x0, 0x10, 0x20}
	cfg.ShardingForkVersion = []byte{0x4, 0x0, 0x10, 0x20}
	cfg.BellatrixForkEpoch = 112260
	cfg.BellatrixForkVersion = []byte{0x2, 0x0, 0x10, 0x20}
	cfg.TerminalTotalDifficulty = "10790000"
	cfg.DepositContractAddress = "0xff50ed3d0ec03aC01D4C79aAd74928BFF48a7b2b"
	cfg.InitializeForkSchedule()
	return cfg
}

// Beacon configs
var BeaconConfigs map[NetworkType]BeaconChainConfig = map[NetworkType]BeaconChainConfig{
	MainnetNetwork: mainnetConfig(),
	SepoliaNetwork: sepoliaConfig(),
	GoerliNetwork:  goerliConfig(),
}

func GetConfigsByNetwork(net NetworkType) (*discover.Config, GenesisConfig, NetworkConfig, BeaconChainConfig, error) {
	networkConfig := NetworkConfigs[net]
	bootnodes := networkConfig.bootNodes
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, GenesisConfig{}, NetworkConfig{}, BeaconChainConfig{}, err
	}
	enodes := []*enode.Node{}
	for _, addr := range bootnodes {
		enode, err := enode.Parse(enode.ValidSchemes, addr)
		if err != nil {
			return nil, GenesisConfig{}, NetworkConfig{}, BeaconChainConfig{}, err
		}
		enodes = append(enodes, enode)
	}
	return &discover.Config{
		PrivateKey: privateKey,
		Bootnodes:  enodes,
	}, GenesisConfigs[net], networkConfig, BeaconConfigs[net], nil
}
