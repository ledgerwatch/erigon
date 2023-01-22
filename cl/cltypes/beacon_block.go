package cltypes

import (
	"bytes"
	"fmt"

	libcommon "github.com/ledgerwatch/erigon-lib/common"

	"github.com/ledgerwatch/erigon/cl/clparams"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz_utils"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
	"github.com/ledgerwatch/erigon/cl/utils"
	"github.com/ledgerwatch/erigon/ethdb/cbor"
)

/*
 * Block body for Consensus Layer to be stored internally (payload and attestations are stored separatedly).
 */
type BeaconBlockForStorage struct {
	// Non-body fields
	Signature     [96]byte
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    libcommon.Hash
	StateRoot     libcommon.Hash
	// Body fields
	RandaoReveal      [96]byte
	Eth1Data          *Eth1Data
	Graffiti          []byte
	ProposerSlashings []*ProposerSlashing
	AttesterSlashings []*AttesterSlashing
	Deposits          []*Deposit
	VoluntaryExits    []*SignedVoluntaryExit
	SyncAggregate     *SyncAggregate
	// Metadatas
	Eth1Number    uint64
	Eth1BlockHash libcommon.Hash
	Eth2BlockRoot libcommon.Hash
	// Version type
	Version uint8
}

const (
	MaxAttesterSlashings = 2
	MaxProposerSlashings = 16
	MaxAttestations      = 128
	MaxDeposits          = 16
	MaxVoluntaryExits    = 16
)

func getBeaconBlockMinimumSize(v clparams.StateVersion) (size uint32) {
	switch v {
	case clparams.BellatrixVersion:
		size = 384
	case clparams.AltairVersion:
		size = 380
	case clparams.Phase0Version:
		size = 220
	default:
		panic("unimplemented version")
	}
	return
}

type SignedBeaconBlock struct {
	Signature [96]byte
	Block     *BeaconBlock
}

type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    libcommon.Hash
	StateRoot     libcommon.Hash
	Body          *BeaconBody
}

type BeaconBody struct {
	// A byte array used for randomness in the beacon chain
	RandaoReveal [96]byte
	// Data related to the Ethereum 1.0 chain
	Eth1Data *Eth1Data
	// A byte array used to customize validators' behavior
	Graffiti []byte
	// A list of slashing events for validators who included invalid blocks in the chain
	ProposerSlashings []*ProposerSlashing
	// A list of slashing events for validators who included invalid attestations in the chain
	AttesterSlashings []*AttesterSlashing
	// A list of attestations included in the block
	Attestations []*Attestation
	// A list of deposits made to the Ethereum 1.0 chain
	Deposits []*Deposit
	// A list of validators who have voluntarily exited the beacon chain
	VoluntaryExits []*SignedVoluntaryExit
	// A summary of the current state of the beacon chain
	SyncAggregate *SyncAggregate
	// Data related to crosslink records and executing operations on the Ethereum 2.0 chain
	ExecutionPayload *Eth1Block
	// The version of the beacon chain
	Version clparams.StateVersion
}

// Getters

// Version returns beacon block version.
func (b *SignedBeaconBlock) Version() clparams.StateVersion {
	return b.Block.Body.Version
}

// Version returns beacon block version.
func (b *BeaconBlock) Version() clparams.StateVersion {
	return b.Body.Version
}

func (b *BeaconBody) EncodeSSZ(dst []byte) ([]byte, error) {
	buf := dst
	var err error
	//start := len(buf)
	offset := getBeaconBlockMinimumSize(b.Version)
	// Write "easy" fields
	buf = append(buf, b.RandaoReveal[:]...)
	if buf, err = b.Eth1Data.EncodeSSZ(buf); err != nil {
		return nil, err
	}
	if len(b.Graffiti) != 32 {
		return nil, fmt.Errorf("bad graffiti length")
	}
	buf = append(buf, b.Graffiti...)
	// Write offsets for proposer slashings
	buf = append(buf, ssz_utils.OffsetSSZ(offset)...)
	offset += uint32(len(b.ProposerSlashings)) * 416
	// Attester slashings offset
	buf = append(buf, ssz_utils.OffsetSSZ(offset)...)
	for _, slashing := range b.AttesterSlashings {
		offset += uint32(slashing.EncodingSizeSSZ()) + 4
	}
	// Attestation offset
	buf = append(buf, ssz_utils.OffsetSSZ(offset)...)
	for _, attestation := range b.Attestations {
		offset += uint32(attestation.EncodingSizeSSZ()) + 4
	}
	// Deposits offset
	buf = append(buf, ssz_utils.OffsetSSZ(offset)...)
	offset += uint32(len(b.Deposits)) * 1240
	// Voluntary Exit offset
	buf = append(buf, ssz_utils.OffsetSSZ(offset)...)
	offset += uint32(len(b.VoluntaryExits)) * 112
	// Encode Sync Aggregate
	if b.Version >= clparams.AltairVersion {
		buf = b.SyncAggregate.EncodeSSZ(buf)
	}
	if b.Version >= clparams.BellatrixVersion {
		buf = append(buf, ssz_utils.OffsetSSZ(offset)...)
	}
	// Now start encoding the rest of the fields.
	if len(b.AttesterSlashings) > MaxAttesterSlashings {
		return nil, fmt.Errorf("Encode(SSZ): too many attester slashings")
	}
	if len(b.ProposerSlashings) > MaxAttesterSlashings {
		return nil, fmt.Errorf("Encode(SSZ): too many proposer slashings")
	}
	if len(b.Attestations) > MaxAttestations {
		return nil, fmt.Errorf("Encode(SSZ): too many attestations")
	}
	if len(b.Deposits) > MaxDeposits {
		return nil, fmt.Errorf("Encode(SSZ): too many attestations")
	}
	if len(b.VoluntaryExits) > MaxVoluntaryExits {
		return nil, fmt.Errorf("Encode(SSZ): too many attestations")
	}
	// Write proposer slashings
	for _, proposerSlashing := range b.ProposerSlashings {
		if buf, err = proposerSlashing.EncodeSSZ(buf); err != nil {
			return nil, err
		}
	}
	// Write attester slashings as a dynamic list.
	subOffset := len(b.AttesterSlashings) * 4
	for _, attesterSlashing := range b.AttesterSlashings {
		buf = append(buf, ssz_utils.OffsetSSZ(uint32(subOffset))...)
		subOffset += attesterSlashing.EncodingSizeSSZ()
	}

	for _, attesterSlashing := range b.AttesterSlashings {
		buf, err = attesterSlashing.EncodeSSZ(buf)
		if err != nil {
			return nil, err
		}
	}
	// Attestation
	subOffset = len(b.Attestations) * 4
	for _, attestation := range b.Attestations {
		buf = append(buf, ssz_utils.OffsetSSZ(uint32(subOffset))...)
		subOffset += attestation.EncodingSizeSSZ()
	}
	for _, attestation := range b.Attestations {
		buf, err = attestation.EncodeSSZ(buf)
		if err != nil {
			return nil, err
		}
	}

	for _, deposit := range b.Deposits {
		buf = deposit.EncodeSSZ(buf)
	}

	for _, exit := range b.VoluntaryExits {
		buf = exit.EncodeSSZ(buf)
	}

	if b.Version >= clparams.BellatrixVersion {
		buf, err = b.ExecutionPayload.EncodeSSZ(buf)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (b *BeaconBody) EncodingSizeSSZ() (size int) {
	size = int(getBeaconBlockMinimumSize(b.Version))

	size += len(b.ProposerSlashings) * 416

	for _, slashing := range b.AttesterSlashings {
		size += 4
		size += slashing.EncodingSizeSSZ()
	}

	for _, attestation := range b.Attestations {
		size += 4
		size += attestation.EncodingSizeSSZ()
	}

	size += len(b.Deposits) * 1240
	size += len(b.VoluntaryExits) * 112

	if b.Version >= clparams.BellatrixVersion {
		if b.ExecutionPayload == nil {
			b.ExecutionPayload = new(Eth1Block)
		}
		size += b.ExecutionPayload.EncodingSizeSSZ()
	}

	return
}

func (b *BeaconBody) DecodeSSZ(buf []byte, version clparams.StateVersion) error {
	b.Version = version
	var err error

	if len(buf) < b.EncodingSizeSSZ() {
		return ssz_utils.ErrLowBufferSize
	}

	// Start wildly decoding this thing
	copy(b.RandaoReveal[:], buf)
	// Decode ethereum 1 data.
	b.Eth1Data = new(Eth1Data)
	if err := b.Eth1Data.DecodeSSZ(buf[96:168]); err != nil {
		return err
	}
	// Decode graffiti.
	b.Graffiti = libcommon.Copy(buf[168:200])

	// Decode offsets
	offSetProposerSlashings := ssz_utils.DecodeOffset(buf[200:])
	offsetAttesterSlashings := ssz_utils.DecodeOffset(buf[204:])
	offsetAttestations := ssz_utils.DecodeOffset(buf[208:])
	offsetDeposits := ssz_utils.DecodeOffset(buf[212:])
	offsetExits := ssz_utils.DecodeOffset(buf[216:])
	// Decode sync aggregate if we are past altair.
	if version >= clparams.AltairVersion {
		if len(buf) < 380 {
			return ssz_utils.ErrLowBufferSize
		}
		b.SyncAggregate = new(SyncAggregate)
		if err := b.SyncAggregate.DecodeSSZ(buf[220:380]); err != nil {
			return err
		}
	}

	// Execution Payload offset if past bellatrix.
	var offsetExecution uint32
	if version >= clparams.BellatrixVersion {
		offsetExecution = ssz_utils.DecodeOffset(buf[380:])
	}
	// Decode Proposer slashings
	proposerSlashingLength := 416
	b.ProposerSlashings, err = ssz_utils.DecodeStaticList[*ProposerSlashing](buf, offSetProposerSlashings, offsetAttesterSlashings, uint32(proposerSlashingLength), MaxProposerSlashings)
	if err != nil {
		return err
	}
	// Decode attester slashings
	b.AttesterSlashings, err = ssz_utils.DecodeDynamicList[*AttesterSlashing](buf, offsetAttesterSlashings, offsetAttestations, uint32(MaxAttesterSlashings))
	if err != nil {
		return err
	}
	// Decode attestations
	b.Attestations, err = ssz_utils.DecodeDynamicList[*Attestation](buf, offsetAttestations, offsetDeposits, uint32(MaxAttestations))
	if err != nil {
		return err
	}
	// Decode deposits
	depositsLength := 1240
	b.Deposits, err = ssz_utils.DecodeStaticList[*Deposit](buf, offsetDeposits, offsetExits, uint32(depositsLength), MaxDeposits)
	if err != nil {
		return err
	}
	// Decode exits
	exitLength := 112
	endExitBuffer := len(buf)
	if b.Version >= clparams.BellatrixVersion {
		endExitBuffer = int(offsetExecution)
	}
	b.VoluntaryExits, err = ssz_utils.DecodeStaticList[*SignedVoluntaryExit](buf, offsetExits, uint32(endExitBuffer), uint32(exitLength), MaxVoluntaryExits)
	if err != nil {
		return err
	}
	if b.Version >= clparams.BellatrixVersion {
		b.ExecutionPayload = new(Eth1Block)
		if err := b.ExecutionPayload.DecodeSSZ(buf[offsetExecution:], b.Version); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconBody) HashSSZ() ([32]byte, error) {
	leaves := make([][32]byte, 0, 16)
	// Signature leaf
	randaoLeaf, err := merkle_tree.SignatureRoot(b.RandaoReveal)
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, randaoLeaf)
	// Eth1Data Leaf
	dataLeaf, err := b.Eth1Data.HashSSZ()
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, dataLeaf)
	// Graffiti leaf
	var graffitiLeaf [32]byte
	copy(graffitiLeaf[:], b.Graffiti)
	leaves = append(leaves, graffitiLeaf)
	// Proposer slashings leaf
	proposerLeaf, err := merkle_tree.ListObjectSSZRoot(b.ProposerSlashings, MaxProposerSlashings)
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, proposerLeaf)
	// Attester slashings leaf
	attesterLeaf, err := merkle_tree.ListObjectSSZRoot(b.AttesterSlashings, MaxAttesterSlashings)
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, attesterLeaf)
	// Attestations leaf
	attestationLeaf, err := merkle_tree.ListObjectSSZRoot(b.Attestations, MaxAttestations)
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, attestationLeaf)
	// Deposits leaf
	depositLeaf, err := merkle_tree.ListObjectSSZRoot(b.Deposits, MaxDeposits)
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, depositLeaf)
	// Voluntary exits leaf
	exitLeaf, err := merkle_tree.ListObjectSSZRoot(b.VoluntaryExits, MaxVoluntaryExits)
	if err != nil {
		return [32]byte{}, err
	}
	leaves = append(leaves, exitLeaf)
	// Sync aggreate leaf
	if b.Version >= clparams.AltairVersion {
		aggLeaf, err := b.SyncAggregate.HashSSZ()
		if err != nil {
			return [32]byte{}, err
		}
		leaves = append(leaves, aggLeaf)
	}
	if b.Version >= clparams.BellatrixVersion {
		payloadLeaf, err := b.ExecutionPayload.HashSSZ()
		if err != nil {
			return [32]byte{}, err
		}
		leaves = append(leaves, payloadLeaf)
	}
	if b.Version == clparams.Phase0Version {
		return merkle_tree.ArraysRoot(leaves, 8)
	}
	return merkle_tree.ArraysRoot(leaves, 16)
}

func (b *BeaconBlock) EncodeSSZ(buf []byte) (dst []byte, err error) {
	dst = buf
	// Encode base params
	dst = append(dst, ssz_utils.Uint64SSZ(b.Slot)...)
	dst = append(dst, ssz_utils.Uint64SSZ(b.ProposerIndex)...)
	dst = append(dst, b.ParentRoot[:]...)
	dst = append(dst, b.StateRoot[:]...)
	// Encode body
	dst = append(dst, ssz_utils.OffsetSSZ(84)...)
	if dst, err = b.Body.EncodeSSZ(dst); err != nil {
		return
	}

	return
}

func (b *BeaconBlock) EncodingSizeSSZ() int {
	if b.Body == nil {
		b.Body = new(BeaconBody)
	}
	return 80 + b.Body.EncodingSizeSSZ()
}

func (b *BeaconBlock) DecodeSSZ(buf []byte, version clparams.StateVersion) error {
	if len(buf) < b.EncodingSizeSSZ() {
		return ssz_utils.ErrLowBufferSize
	}
	b.Slot = ssz_utils.UnmarshalUint64SSZ(buf)
	b.ProposerIndex = ssz_utils.UnmarshalUint64SSZ(buf[8:])
	copy(b.ParentRoot[:], buf[16:])
	copy(b.StateRoot[:], buf[48:])
	b.Body = new(BeaconBody)
	return b.Body.DecodeSSZ(buf[84:], version)
}

func (b *BeaconBlock) HashSSZ() ([32]byte, error) {
	blockRoot, err := b.Body.HashSSZ()
	if err != nil {
		return [32]byte{}, err
	}
	return merkle_tree.ArraysRoot([][32]byte{
		merkle_tree.Uint64Root(b.Slot),
		merkle_tree.Uint64Root(b.ProposerIndex),
		b.ParentRoot,
		b.StateRoot,
		blockRoot,
	}, 8)
}

func (b *SignedBeaconBlock) EncodeSSZ(buf []byte) ([]byte, error) {
	dst := buf
	var err error
	dst = append(dst, ssz_utils.OffsetSSZ(100)...)
	dst = append(dst, b.Signature[:]...)
	dst, err = b.Block.EncodeSSZ(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (b *SignedBeaconBlock) EncodingSizeSSZ() int {
	if b.Block == nil {
		b.Block = new(BeaconBlock)
	}
	return 100 + b.Block.EncodingSizeSSZ()
}

func (b *SignedBeaconBlock) DecodeSSZ(buf []byte) error {
	return b.DecodeSSZWithVersion(buf, int(clparams.BellatrixVersion))
}

func (b *SignedBeaconBlock) DecodeSSZWithVersion(buf []byte, s int) error {
	if len(buf) < b.EncodingSizeSSZ() {
		return ssz_utils.ErrLowBufferSize
	}
	copy(b.Signature[:], buf[4:100])
	return b.Block.DecodeSSZ(buf[100:], clparams.StateVersion(s))
}

func (b *SignedBeaconBlock) HashSSZ() ([32]byte, error) {
	blockRoot, err := b.Block.HashSSZ()
	if err != nil {
		return [32]byte{}, err
	}
	signatureRoot, err := merkle_tree.SignatureRoot(b.Signature)
	if err != nil {
		return [32]byte{}, err
	}
	return merkle_tree.ArraysRoot([][32]byte{blockRoot, signatureRoot}, 2)
}

// EncodeForStorage encodes beacon block in snappy compressed CBOR format.
func (b *SignedBeaconBlock) EncodeForStorage() ([]byte, error) {
	var (
		blockRoot libcommon.Hash
		err       error
	)
	if blockRoot, err = b.Block.HashSSZ(); err != nil {
		return nil, err
	}
	storageObject := &BeaconBlockForStorage{
		Signature:         b.Signature,
		Slot:              b.Block.Slot,
		ProposerIndex:     b.Block.ProposerIndex,
		ParentRoot:        b.Block.ParentRoot,
		StateRoot:         b.Block.StateRoot,
		RandaoReveal:      b.Block.Body.RandaoReveal,
		Eth1Data:          b.Block.Body.Eth1Data,
		Graffiti:          b.Block.Body.Graffiti,
		ProposerSlashings: b.Block.Body.ProposerSlashings,
		AttesterSlashings: b.Block.Body.AttesterSlashings,
		Deposits:          b.Block.Body.Deposits,
		VoluntaryExits:    b.Block.Body.VoluntaryExits,
		SyncAggregate:     b.Block.Body.SyncAggregate,
		Version:           uint8(b.Version()),
		Eth2BlockRoot:     blockRoot,
	}

	if b.Version() >= clparams.BellatrixVersion {
		eth1Block := b.Block.Body.ExecutionPayload
		storageObject.Eth1Number = eth1Block.NumberU64()
		storageObject.Eth1BlockHash = eth1Block.Header.BlockHashCL
	}
	var buffer bytes.Buffer
	if err := cbor.Marshal(&buffer, storageObject); err != nil {
		return nil, err
	}
	return utils.CompressSnappy(buffer.Bytes()), nil
}

// DecodeBeaconBlockForStorage decodes beacon block in snappy compressed CBOR format.
func DecodeBeaconBlockForStorage(buf []byte) (block *SignedBeaconBlock, eth1Number uint64, eth1Hash libcommon.Hash, eth2Hash libcommon.Hash, err error) {
	decompressedBuf, err := utils.DecompressSnappy(buf)
	if err != nil {
		return nil, 0, libcommon.Hash{}, libcommon.Hash{}, err
	}
	storageObject := &BeaconBlockForStorage{}
	var buffer bytes.Buffer
	if _, err := buffer.Write(decompressedBuf); err != nil {
		return nil, 0, libcommon.Hash{}, libcommon.Hash{}, err
	}
	if err := cbor.Unmarshal(storageObject, &buffer); err != nil {
		return nil, 0, libcommon.Hash{}, libcommon.Hash{}, err
	}

	return &SignedBeaconBlock{
		Signature: storageObject.Signature,
		Block: &BeaconBlock{
			Slot:          storageObject.Slot,
			ProposerIndex: storageObject.ProposerIndex,
			ParentRoot:    storageObject.ParentRoot,
			StateRoot:     storageObject.StateRoot,
			Body: &BeaconBody{
				RandaoReveal:      storageObject.RandaoReveal,
				Eth1Data:          storageObject.Eth1Data,
				Graffiti:          storageObject.Graffiti,
				ProposerSlashings: storageObject.ProposerSlashings,
				AttesterSlashings: storageObject.AttesterSlashings,
				Deposits:          storageObject.Deposits,
				VoluntaryExits:    storageObject.VoluntaryExits,
				SyncAggregate:     storageObject.SyncAggregate,
				Version:           clparams.StateVersion(storageObject.Version),
			},
		},
	}, storageObject.Eth1Number, storageObject.Eth1BlockHash, storageObject.Eth2BlockRoot, nil
}
