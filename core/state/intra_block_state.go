// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package state provides a caching layer atop the Ethereum state trie.
package state

import (
	"fmt"
	"sort"

	"github.com/holiman/uint256"

	"github.com/ledgerwatch/erigon-lib/chain"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	types2 "github.com/ledgerwatch/erigon-lib/types"
	"github.com/ledgerwatch/erigon/common/u256"
	"github.com/ledgerwatch/erigon/core/blockstm"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/core/types/accounts"
	"github.com/ledgerwatch/erigon/crypto"
	"github.com/ledgerwatch/erigon/turbo/trie"
)

type revision struct {
	id           int
	journalIndex int
}

// SystemAddress - sender address for internal state updates.
var SystemAddress = libcommon.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe")

// BalanceIncrease represents the increase of balance of an account that did not require
// reading the account first
type BalanceIncrease struct {
	increase    uint256.Int
	transferred bool // Set to true when the corresponding stateObject is created and balance increase is transferred to the stateObject
	count       int  // Number of increases - this needs tracking for proper reversion
}

// IntraBlockState is responsible for caching and managing state changes
// that occur during block's execution.
// NOT THREAD SAFE!
type IntraBlockState struct {
	stateReader StateReader

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects      map[libcommon.Address]*stateObject
	stateObjectsDirty map[libcommon.Address]struct{}

	nilAccounts map[libcommon.Address]struct{} // Remember non-existent account to avoid reading them again

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by IntraBlockState.Commit.
	savedErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash libcommon.Hash
	txIndex      int
	logs         map[libcommon.Hash][]*types.Log
	logSize      uint

	// Per-transaction access list
	accessList *accessList

	// Transient storage
	transientStorage transientStorage

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int
	trace          bool
	balanceInc     map[libcommon.Address]*BalanceIncrease // Map of balance increases (without first reading the account)

	// Block-stm related fields
	mvHashmap   *blockstm.MVHashMap
	incarnation int
	readMap     map[blockstm.Key]blockstm.ReadDescriptor
	writeMap    map[blockstm.Key]blockstm.WriteDescriptor
	dep         int
}

// Create a new state from a given trie
func New(stateReader StateReader) *IntraBlockState {
	return &IntraBlockState{
		stateReader:       stateReader,
		stateObjects:      map[libcommon.Address]*stateObject{},
		stateObjectsDirty: map[libcommon.Address]struct{}{},
		nilAccounts:       map[libcommon.Address]struct{}{},
		logs:              map[libcommon.Hash][]*types.Log{},
		journal:           newJournal(),
		accessList:        newAccessList(),
		transientStorage:  newTransientStorage(),
		balanceInc:        map[libcommon.Address]*BalanceIncrease{},
		readMap:           make(map[blockstm.Key]blockstm.ReadDescriptor),
		writeMap:          make(map[blockstm.Key]blockstm.WriteDescriptor),
	}
}

func (ibs *IntraBlockState) SetStateReader(stateReader StateReader) {
	ibs.stateReader = stateReader
}

func (ibs *IntraBlockState) GetStateReader() StateReader {
	return ibs.stateReader
}

func NewWithMVHashmap(stateReader StateReader, mvhm *blockstm.MVHashMap) *IntraBlockState {
	ibs := New(stateReader)
	ibs.mvHashmap = mvhm
	ibs.dep = -1
	return ibs
}

func (ibs *IntraBlockState) SetMVHashmap(mvhm *blockstm.MVHashMap) {
	ibs.mvHashmap = mvhm
	ibs.dep = -1
}

func (ibs *IntraBlockState) GetMVHashmap() *blockstm.MVHashMap {
	return ibs.mvHashmap
}

func (ibs *IntraBlockState) MVWriteList() []blockstm.WriteDescriptor {
	return ibs.MVFullWriteList()
}

func (ibs *IntraBlockState) MVFullWriteList() []blockstm.WriteDescriptor {
	writes := make([]blockstm.WriteDescriptor, 0, len(ibs.writeMap))

	for _, v := range ibs.writeMap {
		writes = append(writes, v)
	}

	return writes
}

func (ibs *IntraBlockState) MVReadMap() map[blockstm.Key]blockstm.ReadDescriptor {
	return ibs.readMap
}

func (ibs *IntraBlockState) MVReadList() []blockstm.ReadDescriptor {
	reads := make([]blockstm.ReadDescriptor, 0, len(ibs.readMap))

	for _, v := range ibs.MVReadMap() {
		reads = append(reads, v)
	}

	return reads
}

func (ibs *IntraBlockState) ensureReadMap() {
	if ibs.readMap == nil {
		ibs.readMap = make(map[blockstm.Key]blockstm.ReadDescriptor)
	}
}

func (ibs *IntraBlockState) ensureWriteMap() {
	if ibs.writeMap == nil {
		ibs.writeMap = make(map[blockstm.Key]blockstm.WriteDescriptor)
	}
}

func (ibs *IntraBlockState) ClearReadMap() {
	ibs.readMap = make(map[blockstm.Key]blockstm.ReadDescriptor)
}

func (ibs *IntraBlockState) ClearWriteMap() {
	ibs.writeMap = make(map[blockstm.Key]blockstm.WriteDescriptor)
}

func (ibs *IntraBlockState) HadInvalidRead() bool {
	return ibs.dep >= 0
}

func (ibs *IntraBlockState) DepTxIndex() int {
	return ibs.dep
}

func (ibs *IntraBlockState) SetSTMIncarnation(inc int) {
	ibs.incarnation = inc
}

func MVRead[T any](s *IntraBlockState, k blockstm.Key, defaultV T, readStorage func(sdb *IntraBlockState) T) (v T) {
	if s.mvHashmap == nil {
		return readStorage(s)
	}

	s.ensureReadMap()

	if s.writeMap != nil {
		if _, ok := s.writeMap[k]; ok {
			return readStorage(s)
		}
	}

	if !k.IsAddress() {
		// If we are reading subpath from a deleted account, return default value instead of reading from MVHashmap
		addr := k.GetAddress()
		stateObject := s.getStateObject(addr)
		if stateObject == nil || stateObject.deleted {
			readStorage(s)
			return defaultV
		}
	}

	res := s.mvHashmap.Read(k, s.txIndex)

	var rd blockstm.ReadDescriptor

	rd.V = blockstm.Version{
		TxnIndex:    res.DepIdx(),
		Incarnation: res.Incarnation(),
	}

	rd.Path = k

	switch res.Status() {
	case blockstm.MVReadResultDone:
		{
			v = readStorage(res.Value().(*IntraBlockState))
			rd.Kind = blockstm.ReadKindMap
		}
	case blockstm.MVReadResultDependency:
		{
			s.dep = res.DepIdx()
			panic("Found denpendency")
		}
	case blockstm.MVReadResultNone:
		{
			v = readStorage(s)
			rd.Kind = blockstm.ReadKindStorage
		}
	default:
		return defaultV
	}

	// TODO: I assume we don't want to overwrite an existing read because this could - for example - change a storage
	//  read to map if the same value is read multiple times.
	if _, ok := s.readMap[k]; !ok {
		s.readMap[k] = rd
	}

	return
}

func (s *IntraBlockState) Version() blockstm.Version {
	return blockstm.Version{
		TxnIndex:    s.txIndex,
		Incarnation: s.incarnation,
	}
}

func MVWrite(s *IntraBlockState, k blockstm.Key) {
	if s.mvHashmap != nil {
		s.ensureWriteMap()

		s.writeMap[k] = blockstm.WriteDescriptor{
			Path: k,
			V:    s.Version(),
			Val:  s,
		}
	}
}

func MVWritten(s *IntraBlockState, k blockstm.Key) bool {
	if s.mvHashmap == nil || s.writeMap == nil {
		return false
	}

	_, ok := s.writeMap[k]

	return ok
}

// mvRecordWritten checks whether a state object is already present in the current MV writeMap.
// If yes, it returns the object directly.
// If not, it clones the object and inserts it into the writeMap before returning it.
func (s *IntraBlockState) mvRecordWritten(object *stateObject) *stateObject {
	if s.mvHashmap == nil {
		return object
	}

	addrKey := blockstm.NewAddressKey(object.Address())

	if MVWritten(s, addrKey) {
		return object
	}

	// Deepcopy is needed to ensure that objects are not written by multiple transactions at the same time, because
	// the input state object can come from a different transaction.
	s.setStateObject(object.address, object.deepCopy(s))
	MVWrite(s, addrKey)

	return s.stateObjects[object.Address()]
}

// Apply entries in the write set to MVHashMap. Note that this function does not clear the write set.
func (s *IntraBlockState) FlushMVWriteSet() {
	if s.mvHashmap != nil && s.writeMap != nil {
		s.mvHashmap.FlushMVWriteSet(s.MVFullWriteList())
	}
}

// Apply entries in a given write set to StateDB. Note that this function does not change MVHashMap nor write set
// of the current StateDB.
func (sw *IntraBlockState) ApplyMVWriteSet(writes []blockstm.WriteDescriptor) {
	for i := range writes {
		path := writes[i].Path
		sr := writes[i].Val.(*IntraBlockState)

		addr := path.GetAddress()

		if sr.getStateObject(addr) != nil {
			sw.SetIncarnation(addr, sr.GetIncarnation(addr))

			if path.IsState() {
				stateKey := path.GetStateKey()
				var state uint256.Int
				sr.GetState(addr, &stateKey, &state)
				sw.SetState(addr, &stateKey, state)
			} else if path.IsAddress() {
				continue
			} else {
				switch path.GetSubpath() {
				case BalancePath:
					sw.SetBalance(addr, sr.GetBalance(addr))
				case NoncePath:
					sw.SetNonce(addr, sr.GetNonce(addr))
				case CodePath:
					sw.SetCode(addr, sr.GetCode(addr))
				case SuicidePath:
					stateObject := sr.getStateObject(addr)
					if stateObject != nil && stateObject.deleted {
						sw.Selfdestruct(addr)
					}
				default:
					panic(fmt.Errorf("unknown key type: %d", path.GetSubpath()))
				}
			}
		}
	}
}

func (sdb *IntraBlockState) SetTrace(trace bool) {
	sdb.trace = trace
}

// setErrorUnsafe sets error but should be called in medhods that already have locks
func (sdb *IntraBlockState) setErrorUnsafe(err error) {
	if sdb.savedErr == nil {
		sdb.savedErr = err
	}
}

func (sdb *IntraBlockState) Error() error {
	return sdb.savedErr
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying state trie to avoid reloading data for the next operations.
func (sdb *IntraBlockState) Reset() {
	//if len(sdb.nilAccounts) == 0 || len(sdb.stateObjects) == 0 || len(sdb.stateObjectsDirty) == 0 || len(sdb.balanceInc) == 0 {
	//	log.Warn("zero", "len(sdb.nilAccounts)", len(sdb.nilAccounts),
	//		"len(sdb.stateObjects)", len(sdb.stateObjects),
	//		"len(sdb.stateObjectsDirty)", len(sdb.stateObjectsDirty),
	//		"len(sdb.balanceInc)", len(sdb.balanceInc))
	//}
	sdb.nilAccounts = make(map[libcommon.Address]struct{})
	sdb.stateObjects = make(map[libcommon.Address]*stateObject)
	sdb.stateObjectsDirty = make(map[libcommon.Address]struct{})
	sdb.logs = make(map[libcommon.Hash][]*types.Log)
	sdb.balanceInc = make(map[libcommon.Address]*BalanceIncrease)
	sdb.thash = libcommon.Hash{}
	sdb.bhash = libcommon.Hash{}
	sdb.txIndex = 0
	sdb.logSize = 0
	//sdb.clearJournalAndRefund()
	//sdb.accessList = newAccessList() // this reset by .Prepare() method
	//sdb.nilAccounts = make(map[libcommon.Address]struct{})
	//sdb.stateObjects = make(map[libcommon.Address]*stateObject)
	//sdb.stateObjectsDirty = make(map[libcommon.Address]struct{})
	//sdb.thash = libcommon.Hash{}
	//sdb.bhash = libcommon.Hash{}
	//sdb.txIndex = 0
	//sdb.logs = make(map[libcommon.Hash][]*types.Log)
	//sdb.logSize = 0
	//sdb.clearJournalAndRefund()
	//sdb.accessList = newAccessList()
	//sdb.balanceInc = make(map[libcommon.Address]*BalanceIncrease)

	sdb.readMap = nil
	sdb.writeMap = nil
	sdb.dep = -1
}

func (sdb *IntraBlockState) AddLog(log2 *types.Log) {
	sdb.journal.append(addLogChange{txhash: sdb.thash})
	log2.TxHash = sdb.thash
	log2.BlockHash = sdb.bhash
	log2.TxIndex = uint(sdb.txIndex)
	log2.Index = sdb.logSize
	sdb.logs[sdb.thash] = append(sdb.logs[sdb.thash], log2)
	sdb.logSize++
}

func (sdb *IntraBlockState) GetLogs(hash libcommon.Hash) []*types.Log {
	return sdb.logs[hash]
}

func (sdb *IntraBlockState) Logs() []*types.Log {
	var logs []*types.Log
	for _, lgs := range sdb.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

// AddRefund adds gas to the refund counter
func (sdb *IntraBlockState) AddRefund(gas uint64) {
	sdb.journal.append(refundChange{prev: sdb.refund})
	sdb.refund += gas
}

// SubRefund removes gas from the refund counter.
// This method will panic if the refund counter goes below zero
func (sdb *IntraBlockState) SubRefund(gas uint64) {
	sdb.journal.append(refundChange{prev: sdb.refund})
	if gas > sdb.refund {
		sdb.setErrorUnsafe(fmt.Errorf("refund counter below zero"))
	}
	sdb.refund -= gas
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (sdb *IntraBlockState) Exist(addr libcommon.Address) bool {
	s := sdb.getStateObject(addr)
	return s != nil && !s.deleted
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (sdb *IntraBlockState) Empty(addr libcommon.Address) bool {
	so := sdb.getStateObject(addr)
	return so == nil || so.deleted || so.empty()
}

const BalancePath = 1
const NoncePath = 2
const CodePath = 3
const SuicidePath = 4

// GetBalance retrieves the balance from the given address or 0 if object not found
// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) GetBalance(addr libcommon.Address) *uint256.Int {
	return MVRead(sdb, blockstm.NewSubpathKey(addr, BalancePath), u256.Num0, func(s *IntraBlockState) *uint256.Int {
		stateObject := s.getStateObject(addr)
		if stateObject != nil && !stateObject.deleted {
			return stateObject.Balance()
		}
		return u256.Num0
	})
}

// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) GetNonce(addr libcommon.Address) uint64 {
	return MVRead(sdb, blockstm.NewSubpathKey(addr, NoncePath), 0, func(s *IntraBlockState) uint64 {
		stateObject := s.getStateObject(addr)
		if stateObject != nil && !stateObject.deleted {
			return stateObject.Nonce()
		}

		return 0
	})
}

// TxIndex returns the current transaction index set by Prepare.
func (sdb *IntraBlockState) TxIndex() int {
	return sdb.txIndex
}

// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) GetCode(addr libcommon.Address) []byte {
	return MVRead(sdb, blockstm.NewSubpathKey(addr, CodePath), nil, func(s *IntraBlockState) []byte {
		stateObject := s.getStateObject(addr)
		if stateObject != nil && !stateObject.deleted {
			if s.trace {
				fmt.Printf("GetCode %x, returned %d\n", addr, len(stateObject.Code()))
			}
			return stateObject.Code()
		}
		if s.trace {
			fmt.Printf("GetCode %x, returned nil\n", addr)
		}
		return nil
	})
}

// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) GetCodeSize(addr libcommon.Address) int {
	return MVRead(sdb, blockstm.NewSubpathKey(addr, CodePath), 0, func(s *IntraBlockState) int {
		stateObject := s.getStateObject(addr)
		if stateObject == nil || stateObject.deleted {
			return 0
		}
		if stateObject.code != nil {
			return len(stateObject.code)
		}
		l, err := s.stateReader.ReadAccountCodeSize(addr, stateObject.data.Incarnation, stateObject.data.CodeHash)
		if err != nil {
			s.setErrorUnsafe(err)
		}
		return l
	})
}

func (sdb *IntraBlockState) GetCodeHash(addr libcommon.Address) libcommon.Hash {
	return MVRead(sdb, blockstm.NewSubpathKey(addr, CodePath), libcommon.Hash{}, func(s *IntraBlockState) libcommon.Hash {
		stateObject := s.getStateObject(addr)
		if stateObject == nil || stateObject.deleted {
			return libcommon.Hash{}
		}
		return libcommon.BytesToHash(stateObject.CodeHash())
	})
}

// GetState retrieves a value from the given account's storage trie.
// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) GetState(addr libcommon.Address, key *libcommon.Hash, value *uint256.Int) {
	MVRead(sdb, blockstm.NewStateKey(addr, *key), nil, func(s *IntraBlockState) *uint256.Int {
		stateObject := s.getStateObject(addr)
		if stateObject != nil && !stateObject.deleted {
			stateObject.GetState(key, value)
		} else {
			value.Clear()
		}

		return value
	})
}

// GetCommittedState retrieves a value from the given account's committed storage trie.
// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) GetCommittedState(addr libcommon.Address, key *libcommon.Hash, value *uint256.Int) {
	MVRead(sdb, blockstm.NewStateKey(addr, *key), nil, func(s *IntraBlockState) *uint256.Int {
		stateObject := s.getStateObject(addr)
		if stateObject != nil && !stateObject.deleted {
			stateObject.GetCommittedState(key, value)
		} else {
			value.Clear()
		}

		return value
	})
}

func (sdb *IntraBlockState) HasSelfdestructed(addr libcommon.Address) bool {
	return MVRead(sdb, blockstm.NewSubpathKey(addr, SuicidePath), false, func(s *IntraBlockState) bool {
		stateObject := s.getStateObject(addr)
		if stateObject == nil {
			return false
		}
		if stateObject.deleted {
			return false
		}
		if stateObject.created {
			return false
		}
		return stateObject.selfdestructed
	})
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) AddBalance(addr libcommon.Address, amount *uint256.Int) {
	if sdb.trace {
		fmt.Printf("AddBalance %x, %d\n", addr, amount)
	}

	if sdb.mvHashmap != nil {
		sdb.GetBalance(addr)
	}

	// If this account has not been read, add to the balance increment map
	_, needAccount := sdb.stateObjects[addr]
	if (!needAccount && addr == ripemd && amount.IsZero()) || sdb.mvHashmap != nil {
		needAccount = true
	}
	if !needAccount {
		sdb.journal.append(balanceIncrease{
			account:  &addr,
			increase: *amount,
		})
		bi, ok := sdb.balanceInc[addr]
		if !ok {
			bi = &BalanceIncrease{}
			sdb.balanceInc[addr] = bi
		}
		bi.increase.Add(&bi.increase, amount)
		bi.count++
		MVWrite(sdb, blockstm.NewSubpathKey(addr, BalancePath))
		return
	}

	stateObject := sdb.GetOrNewStateObject(addr)
	stateObject = sdb.mvRecordWritten(stateObject)
	stateObject.AddBalance(amount)
	MVWrite(sdb, blockstm.NewSubpathKey(addr, BalancePath))
}

// SubBalance subtracts amount from the account associated with addr.
// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) SubBalance(addr libcommon.Address, amount *uint256.Int) {
	if sdb.trace {
		fmt.Printf("SubBalance %x, %d\n", addr, amount)
	}

	if sdb.mvHashmap != nil {
		// ensure a read balance operation is recorded in mvHashmap
		sdb.GetBalance(addr)
	}

	MVWrite(sdb, blockstm.NewSubpathKey(addr, BalancePath))

	stateObject := sdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject = sdb.mvRecordWritten(stateObject)
		stateObject.SubBalance(amount)
	}
}

// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) SetBalance(addr libcommon.Address, amount *uint256.Int) {
	stateObject := sdb.GetOrNewStateObject(addr)

	MVWrite(sdb, blockstm.NewSubpathKey(addr, BalancePath))

	if stateObject != nil {
		stateObject = sdb.mvRecordWritten(stateObject)
		stateObject.SetBalance(amount)
	}
}

// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) SetNonce(addr libcommon.Address, nonce uint64) {
	stateObject := sdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject = sdb.mvRecordWritten(stateObject)
		stateObject.SetNonce(nonce)
		MVWrite(sdb, blockstm.NewSubpathKey(addr, NoncePath))
	}
}

// DESCRIBED: docs/programmers_guide/guide.md#code-hash
// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) SetCode(addr libcommon.Address, code []byte) {
	stateObject := sdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject = sdb.mvRecordWritten(stateObject)
		stateObject.SetCode(crypto.Keccak256Hash(code), code)
		MVWrite(sdb, blockstm.NewSubpathKey(addr, CodePath))
	}
}

// DESCRIBED: docs/programmers_guide/guide.md#address---identifier-of-an-account
func (sdb *IntraBlockState) SetState(addr libcommon.Address, key *libcommon.Hash, value uint256.Int) {
	stateObject := sdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject = sdb.mvRecordWritten(stateObject)
		stateObject.SetState(key, value)
		MVWrite(sdb, blockstm.NewStateKey(addr, *key))
	}
}

// SetStorage replaces the entire storage for the specified account with given
// storage. This function should only be used for debugging.
func (sdb *IntraBlockState) SetStorage(addr libcommon.Address, storage Storage) {
	stateObject := sdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetStorage(storage)
	}
}

func (sdb *IntraBlockState) SetBlockSTMIncarnation(incarnation int) {
	sdb.incarnation = incarnation
}

// SetIncarnation sets incarnation for account if account exists
func (sdb *IntraBlockState) SetIncarnation(addr libcommon.Address, incarnation uint64) {
	stateObject := sdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.setIncarnation(incarnation)
	}
}

func (sdb *IntraBlockState) GetIncarnation(addr libcommon.Address) uint64 {
	stateObject := sdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.data.Incarnation
	}
	return 0
}

// Selfdestruct marks the given account as suicided.
// This clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (sdb *IntraBlockState) Selfdestruct(addr libcommon.Address) bool {
	stateObject := sdb.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		return false
	}

	stateObject = sdb.mvRecordWritten(stateObject)
	sdb.journal.append(selfdestructChange{
		account:     &addr,
		prev:        stateObject.selfdestructed,
		prevbalance: *stateObject.Balance(),
	})
	stateObject.markSelfdestructed()
	stateObject.created = false
	stateObject.data.Balance.Clear()

	MVWrite(sdb, blockstm.NewSubpathKey(addr, SuicidePath))
	MVWrite(sdb, blockstm.NewSubpathKey(addr, BalancePath))

	return true
}

// SetTransientState sets transient storage for a given account. It
// adds the change to the journal so that it can be rolled back
// to its previous value if there is a revert.
func (sdb *IntraBlockState) SetTransientState(addr libcommon.Address, key libcommon.Hash, value uint256.Int) {
	prev := sdb.GetTransientState(addr, key)
	if prev == value {
		return
	}

	sdb.journal.append(transientStorageChange{
		account:  &addr,
		key:      key,
		prevalue: prev,
	})
	sdb.setTransientState(addr, key, value)
}

// setTransientState is a lower level setter for transient storage. It
// is called during a revert to prevent modifications to the journal.
func (sdb *IntraBlockState) setTransientState(addr libcommon.Address, key libcommon.Hash, value uint256.Int) {
	sdb.transientStorage.Set(addr, key, value)
}

// GetTransientState gets transient storage for a given account.
func (sdb *IntraBlockState) GetTransientState(addr libcommon.Address, key libcommon.Hash) uint256.Int {
	return sdb.transientStorage.Get(addr, key)
}

func (sdb *IntraBlockState) getStateObject(addr libcommon.Address) *stateObject {
	return MVRead(sdb, blockstm.NewAddressKey(addr), nil, func(s *IntraBlockState) *stateObject {
		// Prefer 'live' objects.
		if obj := s.stateObjects[addr]; obj != nil {
			return obj
		}

		// Load the object from the database.
		if _, ok := s.nilAccounts[addr]; ok {
			if bi, ok := s.balanceInc[addr]; ok && !bi.transferred && sdb.mvHashmap == nil {
				return s.createObject(addr, nil)
			}
			return nil
		}
		account, err := s.stateReader.ReadAccountData(addr)
		if err != nil {
			s.setErrorUnsafe(err)
			return nil
		}
		if account == nil {
			s.nilAccounts[addr] = struct{}{}
			if bi, ok := s.balanceInc[addr]; ok && !bi.transferred && sdb.mvHashmap == nil {
				return s.createObject(addr, nil)
			}
			return nil
		}

		// Insert into the live set.
		obj := newObject(s, addr, account, account)
		s.setStateObject(addr, obj)
		return obj
	})
}

func (sdb *IntraBlockState) setStateObject(addr libcommon.Address, object *stateObject) {
	if bi, ok := sdb.balanceInc[addr]; ok && !bi.transferred && sdb.mvHashmap == nil {
		object.data.Balance.Add(&object.data.Balance, &bi.increase)
		bi.transferred = true
		sdb.journal.append(balanceIncreaseTransfer{bi: bi})
	}
	sdb.stateObjects[addr] = object
}

// Retrieve a state object or create a new state object if nil.
func (sdb *IntraBlockState) GetOrNewStateObject(addr libcommon.Address) *stateObject {
	stateObject := sdb.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject = sdb.createObject(addr, stateObject /* previous */)
	}
	return stateObject
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten.
func (sdb *IntraBlockState) createObject(addr libcommon.Address, previous *stateObject) (newobj *stateObject) {
	account := new(accounts.Account)
	var original *accounts.Account
	if previous == nil {
		original = &accounts.Account{}
	} else {
		original = &previous.original
	}
	account.Root.SetBytes(trie.EmptyRoot[:]) // old storage should be ignored
	newobj = newObject(sdb, addr, account, original)
	newobj.setNonce(0) // sets the object to dirty

	MVWrite(sdb, blockstm.NewAddressKey(addr))

	if previous == nil {
		sdb.journal.append(createObjectChange{account: &addr})
	} else {
		sdb.journal.append(resetObjectChange{account: &addr, prev: previous})
	}
	sdb.setStateObject(addr, newobj)
	return newobj
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//  1. sends funds to sha(account ++ (nonce + 1))
//  2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (sdb *IntraBlockState) CreateAccount(addr libcommon.Address, contractCreation bool) {
	var prevInc uint64
	previous := sdb.getStateObject(addr)
	if contractCreation {
		if previous != nil && previous.selfdestructed {
			prevInc = previous.data.Incarnation
		} else {
			if inc, err := sdb.stateReader.ReadAccountIncarnation(addr); err == nil {
				prevInc = inc
			} else {
				sdb.savedErr = err
			}
		}
	}

	newObj := sdb.createObject(addr, previous)
	if previous != nil && !previous.selfdestructed {
		newObj.data.Balance.Set(&previous.data.Balance)
		newObj.data.Initialised = true
	}
	newObj.data.Initialised = true

	MVWrite(sdb, blockstm.NewSubpathKey(addr, BalancePath))

	if contractCreation {
		newObj.created = true
		newObj.data.Incarnation = prevInc + 1
	} else {
		newObj.selfdestructed = false
	}
}

// Snapshot returns an identifier for the current revision of the state.
func (sdb *IntraBlockState) Snapshot() int {
	id := sdb.nextRevisionID
	sdb.nextRevisionID++
	sdb.validRevisions = append(sdb.validRevisions, revision{id, sdb.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (sdb *IntraBlockState) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(sdb.validRevisions), func(i int) bool {
		return sdb.validRevisions[i].id >= revid
	})
	if idx == len(sdb.validRevisions) || sdb.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := sdb.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	sdb.journal.revert(sdb, snapshot)
	sdb.validRevisions = sdb.validRevisions[:idx]
}

// GetRefund returns the current value of the refund counter.
func (sdb *IntraBlockState) GetRefund() uint64 {
	return sdb.refund
}

func updateAccount(EIP161Enabled bool, isAura bool, stateWriter StateWriter, addr libcommon.Address, stateObject *stateObject, isDirty bool) error {
	emptyRemoval := EIP161Enabled && stateObject.empty() && (!isAura || addr != SystemAddress)
	if stateObject.selfdestructed || (isDirty && emptyRemoval) {
		if err := stateWriter.DeleteAccount(addr, &stateObject.original); err != nil {
			return err
		}
		stateObject.deleted = true
	}
	if isDirty && (stateObject.created || !stateObject.selfdestructed) && !emptyRemoval {
		stateObject.deleted = false
		// Write any contract code associated with the state object
		if stateObject.code != nil && stateObject.dirtyCode {
			if err := stateWriter.UpdateAccountCode(addr, stateObject.data.Incarnation, stateObject.data.CodeHash, stateObject.code); err != nil {
				return err
			}
		}
		if stateObject.created {
			if err := stateWriter.CreateContract(addr); err != nil {
				return err
			}
		}
		if err := stateObject.updateTrie(stateWriter); err != nil {
			return err
		}
		if err := stateWriter.UpdateAccountData(addr, &stateObject.original, &stateObject.data); err != nil {
			return err
		}
	}
	return nil
}

func printAccount(EIP161Enabled bool, addr libcommon.Address, stateObject *stateObject, isDirty bool) {
	emptyRemoval := EIP161Enabled && stateObject.empty()
	if stateObject.selfdestructed || (isDirty && emptyRemoval) {
		fmt.Printf("delete: %x\n", addr)
	}
	if isDirty && (stateObject.created || !stateObject.selfdestructed) && !emptyRemoval {
		// Write any contract code associated with the state object
		if stateObject.code != nil && stateObject.dirtyCode {
			fmt.Printf("UpdateCode: %x,%x\n", addr, stateObject.CodeHash())
		}
		if stateObject.created {
			fmt.Printf("CreateContract: %x\n", addr)
		}
		stateObject.printTrie()
		if stateObject.data.Balance.IsUint64() {
			fmt.Printf("UpdateAccountData: %x, balance=%d, nonce=%d\n", addr, stateObject.data.Balance.Uint64(), stateObject.data.Nonce)
		} else {
			div := uint256.NewInt(1_000_000_000)
			fmt.Printf("UpdateAccountData: %x, balance=%d*%d, nonce=%d\n", addr, uint256.NewInt(0).Div(&stateObject.data.Balance, div).Uint64(), div.Uint64(), stateObject.data.Nonce)
		}
	}
}

// FinalizeTx should be called after every transaction.
func (sdb *IntraBlockState) FinalizeTx(chainRules *chain.Rules, stateWriter StateWriter) error {
	for addr, bi := range sdb.balanceInc {
		if !bi.transferred {
			sdb.getStateObject(addr)
		}
	}
	for addr := range sdb.journal.dirties {
		so, exist := sdb.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `sdb.journal.dirties` but not in `sdb.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}

		if err := updateAccount(chainRules.IsSpuriousDragon, chainRules.IsAura, stateWriter, addr, so, true); err != nil {
			return err
		}

		sdb.stateObjectsDirty[addr] = struct{}{}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	sdb.clearJournalAndRefund()
	return nil
}

func (sdb *IntraBlockState) SoftFinalise() {
	for addr := range sdb.journal.dirties {
		_, exist := sdb.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `sdb.journal.dirties` but not in `sdb.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}
		sdb.stateObjectsDirty[addr] = struct{}{}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	sdb.clearJournalAndRefund()
}

// CommitBlock finalizes the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (sdb *IntraBlockState) CommitBlock(chainRules *chain.Rules, stateWriter StateWriter) error {
	for addr, bi := range sdb.balanceInc {
		if !bi.transferred {
			sdb.getStateObject(addr)
		}
	}
	return sdb.MakeWriteSet(chainRules, stateWriter)
}

func (sdb *IntraBlockState) BalanceIncreaseSet() map[libcommon.Address]uint256.Int {
	s := make(map[libcommon.Address]uint256.Int, len(sdb.balanceInc))
	for addr, bi := range sdb.balanceInc {
		if !bi.transferred {
			s[addr] = bi.increase
		}
	}
	return s
}

func (sdb *IntraBlockState) MakeWriteSet(chainRules *chain.Rules, stateWriter StateWriter) error {
	for addr := range sdb.journal.dirties {
		sdb.stateObjectsDirty[addr] = struct{}{}
	}
	for addr, stateObject := range sdb.stateObjects {
		_, isDirty := sdb.stateObjectsDirty[addr]
		if err := updateAccount(chainRules.IsSpuriousDragon, chainRules.IsAura, stateWriter, addr, stateObject, isDirty); err != nil {
			return err
		}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	sdb.clearJournalAndRefund()
	return nil
}

func (sdb *IntraBlockState) Print(chainRules chain.Rules) {
	for addr, stateObject := range sdb.stateObjects {
		_, isDirty := sdb.stateObjectsDirty[addr]
		_, isDirty2 := sdb.journal.dirties[addr]

		printAccount(chainRules.IsSpuriousDragon, addr, stateObject, isDirty || isDirty2)
	}
}

// SetTxContext sets the current transaction hash and index and block hash which are
// used when the EVM emits new state logs. It should be invoked before
// transaction execution.
func (sdb *IntraBlockState) SetTxContext(thash, bhash libcommon.Hash, ti int) {
	sdb.thash = thash
	sdb.bhash = bhash
	sdb.txIndex = ti
}

// no not lock
func (sdb *IntraBlockState) clearJournalAndRefund() {
	sdb.journal = newJournal()
	sdb.validRevisions = sdb.validRevisions[:0]
	sdb.refund = 0
}

// Prepare handles the preparatory steps for executing a state transition.
// This method must be invoked before state transition.
//
// Berlin fork:
// - Add sender to access list (EIP-2929)
// - Add destination to access list (EIP-2929)
// - Add precompiles to access list (EIP-2929)
// - Add the contents of the optional tx access list (EIP-2930)
//
// Shanghai fork:
// - Add coinbase to access list (EIP-3651)
//
// Cancun fork:
// - Reset transient storage (EIP-1153)
func (sdb *IntraBlockState) Prepare(rules *chain.Rules, sender, coinbase libcommon.Address, dst *libcommon.Address,
	precompiles []libcommon.Address, list types2.AccessList,
) {
	if rules.IsBerlin {
		// Clear out any leftover from previous executions
		al := newAccessList()
		sdb.accessList = al

		al.AddAddress(sender)
		if dst != nil {
			al.AddAddress(*dst)
			// If it's a create-tx, the destination will be added inside evm.create
		}
		for _, addr := range precompiles {
			al.AddAddress(addr)
		}
		for _, el := range list {
			al.AddAddress(el.Address)
			for _, key := range el.StorageKeys {
				al.AddSlot(el.Address, key)
			}
		}
		if rules.IsShanghai { // EIP-3651: warm coinbase
			al.AddAddress(coinbase)
		}
	}
	// Reset transient storage at the beginning of transaction execution
	sdb.transientStorage = newTransientStorage()
}

// AddAddressToAccessList adds the given address to the access list
func (sdb *IntraBlockState) AddAddressToAccessList(addr libcommon.Address) {
	if sdb.accessList.AddAddress(addr) {
		sdb.journal.append(accessListAddAccountChange{&addr})
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (sdb *IntraBlockState) AddSlotToAccessList(addr libcommon.Address, slot libcommon.Hash) {
	addrMod, slotMod := sdb.accessList.AddSlot(addr, slot)
	if addrMod {
		// In practice, this should not happen, since there is no way to enter the
		// scope of 'address' without having the 'address' become already added
		// to the access list (via call-variant, create, etc).
		// Better safe than sorry, though
		sdb.journal.append(accessListAddAccountChange{&addr})
	}
	if slotMod {
		sdb.journal.append(accessListAddSlotChange{
			address: &addr,
			slot:    &slot,
		})
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (sdb *IntraBlockState) AddressInAccessList(addr libcommon.Address) bool {
	return sdb.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (sdb *IntraBlockState) SlotInAccessList(addr libcommon.Address, slot libcommon.Hash) (addressPresent bool, slotPresent bool) {
	return sdb.accessList.Contains(addr, slot)
}

// Copy intra block state
func (sdb *IntraBlockState) Copy() *IntraBlockState {
	state := New(sdb.stateReader)
	state.stateObjects = make(map[libcommon.Address]*stateObject, len(sdb.stateObjectsDirty))
	state.stateObjectsDirty = make(map[libcommon.Address]struct{}, len(sdb.stateObjectsDirty))

	for addr := range sdb.journal.dirties {
		if object, exist := sdb.stateObjects[addr]; exist {
			state.stateObjects[addr] = object.deepCopy(state)

			state.stateObjectsDirty[addr] = struct{}{} // Mark the copy dirty to force internal (code/state) commits
		}
	}

	state.validRevisions = append(state.validRevisions, sdb.validRevisions...)
	state.refund = sdb.refund

	for addr := range sdb.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = sdb.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsDirty[addr] = struct{}{}
	}

	for hash, logs := range sdb.logs {
		cpy := make([]*types.Log, len(logs))
		for i, l := range logs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		state.logs[hash] = cpy
	}

	state.accessList = sdb.accessList.Copy()

	state.thash = sdb.thash
	state.bhash = sdb.bhash
	state.txIndex = sdb.txIndex

	if sdb.mvHashmap != nil {
		state.mvHashmap = sdb.mvHashmap
	}

	return state
}
