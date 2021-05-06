// Copyright 2014 The go-ethereum Authors
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

package types

import (
	"io"
	"sync"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/hexutil"
	"github.com/ledgerwatch/turbo-geth/rlp"
)

// go:generate gencodec -type Log -field-override logMarshaling -out gen_log_json.go

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required" codec:"1"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required" codec:"2"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required" codec:"3"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber" codec:"-"`
	// hash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required" codec:"-"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" codec:"-"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash" codec:"-"`
	// index of the log in the block
	Index uint `json:"logIndex" codec:"-"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed" codec:"-"`
}

type Logs []*Log

type logMarshaling struct {
	Data        hexutil.Bytes
	BlockNumber hexutil.Uint64
	TxIndex     hexutil.Uint
	Index       hexutil.Uint
}

type rlpLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}

// rlpStorageLog is the storage encoding of a log.
type rlpStorageLog rlpLog

// legacyRlpStorageLog is the previous storage encoding of a log including some redundant fields.
type legacyRlpStorageLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
	//BlockNumber uint64
	//TxHash      common.Hash
	//TxIndex     uint
	//BlockHash   common.Hash
	//Index       uint
}

// EncodeRLP implements rlp.Encoder.
func (l *Log) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpLog{Address: l.Address, Topics: l.Topics, Data: l.Data})
}

// DecodeRLP implements rlp.Decoder.
func (l *Log) DecodeRLP(s *rlp.Stream) error {
	var dec rlpLog
	err := s.Decode(&dec)
	if err == nil {
		l.Address, l.Topics, l.Data = dec.Address, dec.Topics, dec.Data
	}
	return err
}

// LogForStorage is a wrapper around a Log that flattens and parses the entire content of
// a log including non-consensus fields.
type LogForStorage Log

// EncodeRLP implements rlp.Encoder.
func (l *LogForStorage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpStorageLog{
		Address: l.Address,
		Topics:  l.Topics,
		Data:    l.Data,
		//BlockNumber: l.BlockNumber,
		//TxHash:      l.TxHash,
		//TxIndex:     l.TxIndex,
		//BlockHash:   l.BlockHash,
		//Index:       l.Index,
	})
}

// DecodeRLP implements rlp.Decoder.
//
// Note some redundant fields(e.g. block number, tx hash etc) will be assembled later.
func (l *LogForStorage) DecodeRLP(s *rlp.Stream) error {
	blob, err := s.Raw()
	if err != nil {
		return err
	}
	var dec rlpStorageLog
	err = rlp.DecodeBytes(blob, &dec)
	if err == nil {
		*l = LogForStorage{
			Address: dec.Address,
			Topics:  dec.Topics,
			Data:    dec.Data,
		}
	} else {
		// Try to decode log with previous definition.
		var dec legacyRlpStorageLog
		err = rlp.DecodeBytes(blob, &dec)
		if err == nil {
			*l = LogForStorage{
				Address: dec.Address,
				Topics:  dec.Topics,
				Data:    dec.Data,
			}
		}
	}
	return err
}

// LogsPubSub - it's safe to use this class as non-pointer, do double-unsubscribe
type LogsPubSub struct {
	sync.Mutex
	id    uint
	last  Logs
	chans map[uint]chan Logs
}

func (s *LogsPubSub) Sub() (ch chan Logs, unsubscribe func()) {
	s.Lock()
	defer s.Unlock()
	if s.chans == nil {
		s.chans = make(map[uint]chan Logs)
	}
	s.id++
	id := s.id
	ch = make(chan Logs, 1)
	s.chans[id] = ch
	return ch, func() { s.unsubscribe(id) }
}

func (s *LogsPubSub) Pub(v Logs) {
	s.Lock()
	defer s.Unlock()
	s.last = v
	for _, ch := range s.chans {
		ch <- v
	}
}

func (s *LogsPubSub) Last() Logs {
	s.Lock()
	defer s.Unlock()
	return s.last
}

func (s *LogsPubSub) unsubscribe(id uint) {
	s.Lock()
	defer s.Unlock()
	ch, ok := s.chans[id]
	if !ok { // double-unsubscribe support
		return
	}
	close(ch)
	delete(s.chans, id)
}
