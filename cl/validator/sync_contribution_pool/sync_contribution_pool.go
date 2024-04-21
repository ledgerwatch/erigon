package sync_contribution_pool

import (
	"errors"
	"sync"

	"github.com/Giulio2002/bls"
	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/cl/clparams"
	"github.com/ledgerwatch/erigon/cl/cltypes"
	"github.com/ledgerwatch/erigon/cl/cltypes/solid"
	"github.com/ledgerwatch/erigon/cl/phase1/core/state"
	"github.com/ledgerwatch/erigon/cl/utils"
)

type syncContributionKey struct {
	slot              uint64
	subcommitteeIndex uint64
	beaconBlockRoot   common.Hash
}

type syncContributionPoolImpl struct {
	// syncContributionPool is a map of sync contributions, indexed by slot, subcommittee index and block root.
	syncContributionPool map[syncContributionKey]*cltypes.Contribution
	beaconCfg            *clparams.BeaconChainConfig

	mu sync.Mutex
}

var ErrIsSuperset = errors.New("sync contribution is a superset of existing attestation")

func NewSyncContributionPool() SyncContributionPool {
	return &syncContributionPoolImpl{
		syncContributionPool: make(map[syncContributionKey]*cltypes.Contribution),
	}
}

func getSyncCommitteeFromState(s *state.CachingBeaconState) *solid.SyncCommittee {
	cfg := s.BeaconConfig()
	if cfg.SyncCommitteePeriod(s.Slot()) == cfg.SyncCommitteePeriod(s.Slot()+1) {
		return s.CurrentSyncCommittee()
	}
	return s.NextSyncCommittee()

}

func (s *syncContributionPoolImpl) AddSyncContribution(headState *state.CachingBeaconState, contribution *cltypes.Contribution) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := syncContributionKey{
		slot:              contribution.Slot,
		subcommitteeIndex: contribution.SubcommitteeIndex,
		beaconBlockRoot:   contribution.BeaconBlockRoot,
	}
	baseContribution := &cltypes.Contribution{
		Slot:              contribution.Slot,
		SubcommitteeIndex: contribution.SubcommitteeIndex,
		BeaconBlockRoot:   contribution.BeaconBlockRoot,
		AggregationBits:   make([]byte, cltypes.SyncCommitteeAggregationBitsSize),
		Signature:         bls.InfiniteSignature,
	}

	if val, ok := s.syncContributionPool[key]; ok {
		baseContribution = val.Copy()
	}
	// Time to aggregate the giga aggregatable.
	if utils.IsSupersetBitlist(baseContribution.AggregationBits, contribution.AggregationBits) {
		return ErrIsSuperset // Skip it if it is just a superset.
	}
	// Aggregate the bits.
	utils.MergeBitlists(baseContribution.AggregationBits, contribution.AggregationBits)
	// Aggregate the signature.
	aggregatedSignature, err := bls.AggregateSignatures([][]byte{
		baseContribution.Signature[:],
		contribution.Signature[:],
	})
	if err != nil {
		return err
	}
	copy(baseContribution.Signature[:], aggregatedSignature)

	// Make a copy.
	s.syncContributionPool[key] = baseContribution.Copy()
	s.cleanupOldContributions(headState)
	return nil
}

func (s *syncContributionPoolImpl) GetSyncContribution(slot, subcommitteeIndex uint64, beaconBlockRoot common.Hash) *cltypes.Contribution {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := syncContributionKey{
		slot:              slot,
		subcommitteeIndex: subcommitteeIndex,
		beaconBlockRoot:   beaconBlockRoot,
	}

	contribution, ok := s.syncContributionPool[key]
	// Return a copies.
	if !ok {
		// if we dont have it return an empty contribution (no aggregation bits).
		return &cltypes.Contribution{
			Slot:              slot,
			SubcommitteeIndex: subcommitteeIndex,
			BeaconBlockRoot:   beaconBlockRoot,
			AggregationBits:   make([]byte, cltypes.SyncCommitteeAggregationBitsSize),
			Signature:         bls.InfiniteSignature,
		}
	}
	return contribution.Copy()
}

func (s *syncContributionPoolImpl) cleanupOldContributions(headState *state.CachingBeaconState) {

	for key := range s.syncContributionPool {
		if headState.Slot() != key.slot {
			delete(s.syncContributionPool, key)
		}
	}
}

// AddSyncCommitteeMessage aggregates a sync committee message to a contribution to the pool.
func (s *syncContributionPoolImpl) AddSyncCommitteeMessage(headState *state.CachingBeaconState, subCommittee uint64, message *cltypes.SyncCommitteeMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := headState.BeaconConfig()

	key := syncContributionKey{
		slot:              message.Slot,
		subcommitteeIndex: subCommittee,
		beaconBlockRoot:   message.BeaconBlockRoot,
	}

	// We retrieve a base contribution
	contribution, ok := s.syncContributionPool[key]
	if !ok {
		contribution = &cltypes.Contribution{
			Slot:              message.Slot,
			SubcommitteeIndex: subCommittee,
			BeaconBlockRoot:   message.BeaconBlockRoot,
			AggregationBits:   make([]byte, cltypes.SyncCommitteeAggregationBitsSize),
			Signature:         bls.InfiniteSignature,
		}
	}
	// We use the a copy of this contribution
	contribution = contribution.Copy() // make a copy
	// First we find the aggregation bits to which this validator needs to be turned on.
	publicKey, err := headState.ValidatorPublicKey(int(message.ValidatorIndex))
	if err != nil {
		return err
	}

	committee := getSyncCommitteeFromState(headState).GetCommittee()
	subCommitteeSize := cfg.SyncCommitteeSize / cfg.SyncCommitteeSubnetCount
	startSubCommittee := subCommittee * subCommitteeSize
	for i := startSubCommittee; i < startSubCommittee+subCommitteeSize; i++ {
		if committee[i] == publicKey { // turn on this bit
			utils.FlipBitOn(contribution.AggregationBits, int(i-startSubCommittee))
		}
	}
	// Compute the aggregated signature.
	aggregatedSignature, err := bls.AggregateSignatures([][]byte{
		contribution.Signature[:],
		message.Signature[:],
	})
	if err != nil {
		return err
	}
	copy(contribution.Signature[:], aggregatedSignature)
	s.syncContributionPool[key] = contribution
	s.cleanupOldContributions(headState)
	return nil
}

// GetSyncAggregate computes and returns the sync aggregate for the sync messages pointing to a given beacon block root.
/*
def process_sync_committee_contributions(block: BeaconBlock,
	contributions: Set[SyncCommitteeContribution]) -> None:
	sync_aggregate = SyncAggregate()
	signatures = []
	sync_subcommittee_size = SYNC_COMMITTEE_SIZE // SYNC_COMMITTEE_SUBNET_COUNT

	for contribution in contributions:
		subcommittee_index = contribution.subcommittee_index
		for index, participated in enumerate(contribution.aggregation_bits):
			if participated:
				participant_index = sync_subcommittee_size * subcommittee_index + index
				sync_aggregate.sync_committee_bits[participant_index] = True
		signatures.append(contribution.signature)

	sync_aggregate.sync_committee_signature = bls.Aggregate(signatures)
	block.body.sync_aggregate = sync_aggregate
*/
func (s *syncContributionPoolImpl) GetSyncAggregate(slot uint64, beaconBlockRoot common.Hash) (*cltypes.SyncAggregate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// find all contributions for the given beacon block root.
	contributions := []*cltypes.Contribution{}
	for key, contribution := range s.syncContributionPool {
		if key.beaconBlockRoot == beaconBlockRoot && slot == key.slot {
			contributions = append(contributions, contribution)
		}
	}
	if len(contributions) == 0 {
		return &cltypes.SyncAggregate{ // return an empty aggregate.
			SyncCommiteeSignature: bls.InfiniteSignature,
		}, nil
	}
	aggregate := &cltypes.SyncAggregate{}
	signatures := [][]byte{}
	syncSubCommitteeIndex := cltypes.SyncCommitteeSize / s.beaconCfg.SyncCommitteeSubnetCount
	// triple for-loop for the win.
	for _, contribution := range contributions {
		for i, _ := range contribution.AggregationBits {
			for j := 0; j < 8; j++ {
				bitIndex := i*8 + j
				partecipated := utils.IsBitOn(contribution.AggregationBits, bitIndex)
				if partecipated {
					participantIndex := syncSubCommitteeIndex*contribution.SubcommitteeIndex + uint64(bitIndex)
					utils.FlipBitOn(aggregate.SyncCommiteeBits[:], int(participantIndex))
				}
			}
		}
		signatures = append(signatures, contribution.Signature[:])
	}
	// Aggregate the signatures.
	aggregateSignature, err := bls.AggregateSignatures(signatures)
	if err != nil {
		return &cltypes.SyncAggregate{ // return an empty aggregate.
			SyncCommiteeSignature: bls.InfiniteSignature,
		}, err
	}
	copy(aggregate.SyncCommiteeSignature[:], aggregateSignature)
	return aggregate, nil
}
