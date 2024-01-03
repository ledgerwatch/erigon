package borcfg

import (
	"math/big"
	"sort"
	"strconv"

	"github.com/ledgerwatch/erigon-lib/common"
)

// BorConfig is the consensus engine configs for Matic bor based sealing.
type BorConfig struct {
	Period                map[string]uint64 `json:"period"`                // Number of seconds between blocks to enforce
	ProducerDelay         map[string]uint64 `json:"producerDelay"`         // Number of seconds delay between two producer interval
	Sprint                map[string]uint64 `json:"sprint"`                // Epoch length to proposer
	BackupMultiplier      map[string]uint64 `json:"backupMultiplier"`      // Backup multiplier to determine the wiggle time
	ValidatorContract     string            `json:"validatorContract"`     // Validator set contract
	StateReceiverContract string            `json:"stateReceiverContract"` // State receiver contract

	OverrideStateSyncRecords map[string]int         `json:"overrideStateSyncRecords"` // override state records count
	BlockAlloc               map[string]interface{} `json:"blockAlloc"`

	JaipurBlock                *big.Int          `json:"jaipurBlock"`                // Jaipur switch block (nil = no fork, 0 = already on jaipur)
	DelhiBlock                 *big.Int          `json:"delhiBlock"`                 // Delhi switch block (nil = no fork, 0 = already on delhi)
	IndoreBlock                *big.Int          `json:"indoreBlock"`                // Indore switch block (nil = no fork, 0 = already on indore)
	AgraBlock                  *big.Int          `json:"agraBlock"`                  // Agra switch block (nil = no fork, 0 = already in agra)
	StateSyncConfirmationDelay map[string]uint64 `json:"stateSyncConfirmationDelay"` // StateSync Confirmation Delay, in seconds, to calculate `to`

	ParallelUniverseBlock *big.Int `json:"parallelUniverseBlock"` // TODO: update all occurrence, change name and finalize number (hardfork for block-stm related changes)

	sprints sprints
}

// String implements the stringer interface, returning the consensus engine details.
func (c *BorConfig) String() string {
	return "bor"
}

func (c *BorConfig) CalculateProducerDelay(number uint64) uint64 {
	return borKeyValueConfigHelper(c.ProducerDelay, number)
}

func (c *BorConfig) CalculateSprint(number uint64) uint64 {
	if c.sprints == nil {
		c.sprints = asSprints(c.Sprint)
	}

	for i := 0; i < len(c.sprints)-1; i++ {
		if number >= c.sprints[i].from && number < c.sprints[i+1].from {
			return c.sprints[i].size
		}
	}

	return c.sprints[len(c.sprints)-1].size
}

func (c *BorConfig) CalculateSprintCount(from, to uint64) int {
	switch {
	case from > to:
		return 0
	case from < to:
		to--
	}

	if c.sprints == nil {
		c.sprints = asSprints(c.Sprint)
	}

	count := uint64(0)
	startCalc := from

	zeroth := func(boundary uint64, size uint64) uint64 {
		if boundary%size == 0 {
			return 1
		}

		return 0
	}

	for i := 0; i < len(c.sprints)-1; i++ {
		if startCalc >= c.sprints[i].from && startCalc < c.sprints[i+1].from {
			if to >= c.sprints[i].from && to < c.sprints[i+1].from {
				if startCalc == to {
					return int(count + zeroth(startCalc, c.sprints[i].size))
				}
				return int(count + zeroth(startCalc, c.sprints[i].size) + (to-startCalc)/c.sprints[i].size)
			} else {
				endCalc := c.sprints[i+1].from - 1
				count += zeroth(startCalc, c.sprints[i].size) + (endCalc-startCalc)/c.sprints[i].size
				startCalc = endCalc + 1
			}
		}
	}

	if startCalc == to {
		return int(count + zeroth(startCalc, c.sprints[len(c.sprints)-1].size))
	}

	return int(count + zeroth(startCalc, c.sprints[len(c.sprints)-1].size) + (to-startCalc)/c.sprints[len(c.sprints)-1].size)
}

func (c *BorConfig) CalculateBackupMultiplier(number uint64) uint64 {
	return borKeyValueConfigHelper(c.BackupMultiplier, number)
}

func (c *BorConfig) CalculatePeriod(number uint64) uint64 {
	return borKeyValueConfigHelper(c.Period, number)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s *big.Int, head uint64) bool {
	if s == nil {
		return false
	}
	return s.Uint64() <= head
}

func (c *BorConfig) IsJaipur(number uint64) bool {
	return isForked(c.JaipurBlock, number)
}

func (c *BorConfig) IsDelhi(number uint64) bool {
	return isForked(c.DelhiBlock, number)
}

func (c *BorConfig) IsIndore(number uint64) bool {
	return isForked(c.IndoreBlock, number)
}

// IsAgra returns whether num is either equal to the Agra fork block or greater.
// The Agra hard fork is based on the Shanghai hard fork, but it doesn't include withdrawals.
// Also Agra is activated based on the block number rather than the timestamp.
// Refer to https://forum.polygon.technology/t/pip-28-agra-hardfork
func (c *BorConfig) IsAgra(num uint64) bool {
	return isForked(c.AgraBlock, num)
}

func (c *BorConfig) GetAgraBlock() *big.Int {
	return c.AgraBlock
}

// TODO: modify this function once the block number is finalized
func (c *BorConfig) IsParallelUniverse(number uint64) bool {
	if c.ParallelUniverseBlock != nil {
		if c.ParallelUniverseBlock.Cmp(big.NewInt(0)) == 0 {
			return false
		}
	}

	return isForked(c.ParallelUniverseBlock, number)
}

func (c *BorConfig) CalculateStateSyncDelay(number uint64) uint64 {
	return borKeyValueConfigHelper(c.StateSyncConfirmationDelay, number)
}

func borKeyValueConfigHelper[T uint64 | common.Address](field map[string]T, number uint64) T {
	fieldUint := make(map[uint64]T)
	for k, v := range field {
		keyUint, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			panic(err)
		}
		fieldUint[keyUint] = v
	}

	keys := common.SortedKeys(fieldUint)

	for i := 0; i < len(keys)-1; i++ {
		if number >= keys[i] && number < keys[i+1] {
			return fieldUint[keys[i]]
		}
	}

	return fieldUint[keys[len(keys)-1]]
}

type sprint struct {
	from, size uint64
}

type sprints []sprint

func (s sprints) Len() int {
	return len(s)
}

func (s sprints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sprints) Less(i, j int) bool {
	return s[i].from < s[j].from
}

func asSprints(configSprints map[string]uint64) sprints {
	sprints := make(sprints, len(configSprints))

	i := 0
	for key, value := range configSprints {
		sprints[i].from, _ = strconv.ParseUint(key, 10, 64)
		sprints[i].size = value
		i++
	}

	sort.Sort(sprints)

	return sprints
}
