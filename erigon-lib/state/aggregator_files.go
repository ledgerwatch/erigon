/*
   Copyright 2022 The Erigon contributors

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

package state

import (
	"math/bits"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/length"
	"github.com/ledgerwatch/erigon-lib/kv"
)

type filesItemI struct {
	fi []*filesItem
	i  int
}

type SelectedStaticFilesV3 struct {
	d     [kv.DomainLen][]*filesItem
	dHist [kv.DomainLen][]*filesItem
	dIdx  [kv.DomainLen][]*filesItem
	dI    [kv.DomainLen]int
	ii    [kv.StandaloneIdxLen]*filesItemI
}

func (sf SelectedStaticFilesV3) Close() {
	clist := make([][]*filesItem, 0, uint16(kv.DomainLen)+kv.StandaloneIdxLen)
	for id := range sf.d {
		clist = append(clist, sf.d[id], sf.dIdx[id], sf.dHist[id])
	}

	for _, i := range sf.ii {
		clist = append(clist, i.fi)
	}
	for _, group := range clist {
		for _, item := range group {
			if item != nil {
				if item.decompressor != nil {
					item.decompressor.Close()
				}
				if item.index != nil {
					item.index.Close()
				}
			}
		}
	}
}

func (ac *AggregatorRoTx) staticFilesInRange(r RangesV3) (sf SelectedStaticFilesV3, err error) {
	for id := range ac.d {
		if r.d[id].any() {
			sf.d[id], sf.dIdx[id], sf.dHist[id], sf.dI[id] = ac.d[id].staticFilesInRange(r.d[id])
		}
	}
	for id, rng := range r.ranges {
		if rng != nil && rng.needMerge {
			fi, i := ac.iis[id].staticFilesInRange(rng.from, rng.to)
			sf.ii[id] = &filesItemI{fi, i}
		}
	}
	return sf, err
}

type MergedFilesV3 struct {
	d     [kv.DomainLen]*filesItem
	dHist [kv.DomainLen]*filesItem
	dIdx  [kv.DomainLen]*filesItem
	iis   [kv.StandaloneIdxLen]*filesItem
}

func (mf MergedFilesV3) FrozenList() (frozen []string) {
	for id, d := range mf.d {
		if d == nil {
			continue
		}
		frozen = append(frozen, d.decompressor.FileName())

		if mf.dHist[id] != nil && mf.dHist[id].frozen {
			frozen = append(frozen, mf.dHist[id].decompressor.FileName())
		}
		if mf.dIdx[id] != nil && mf.dIdx[id].frozen {
			frozen = append(frozen, mf.dIdx[id].decompressor.FileName())
		}
	}

	for _, ii := range mf.iis {
		if ii != nil && ii.frozen {
			frozen = append(frozen, ii.decompressor.FileName())
		}
	}
	return frozen
}
func (mf MergedFilesV3) Close() {
	clist := make([]*filesItem, 0, kv.DomainLen+4)
	for id := range mf.d {
		clist = append(clist, mf.d[id], mf.dHist[id], mf.dIdx[id])
	}
	clist = append(clist, mf.iis[:]...)

	for _, item := range clist {
		if item != nil {
			if item.decompressor != nil {
				item.decompressor.Close()
			}
			if item.index != nil {
				item.index.Close()
			}
		}
	}
}

type MergedFiles struct {
	d     [kv.DomainLen]*filesItem
	dHist [kv.DomainLen]*filesItem
	dIdx  [kv.DomainLen]*filesItem
}

func (mf MergedFiles) FillV3(m *MergedFilesV3) MergedFiles {
	for id := range m.d {
		mf.d[id], mf.dHist[id], mf.dIdx[id] = m.d[id], m.dHist[id], m.dIdx[id]
	}
	return mf
}

func (mf MergedFiles) Close() {
	for id := range mf.d {
		for _, item := range []*filesItem{mf.d[id], mf.dHist[id], mf.dIdx[id]} {
			if item != nil {
				if item.decompressor != nil {
					item.decompressor.Close()
				}
				if item.decompressor != nil {
					item.index.Close()
				}
				if item.bindex != nil {
					item.bindex.Close()
				}
			}
		}
	}
}

func DecodeAccountBytes(enc []byte) (nonce uint64, balance *uint256.Int, hash []byte) {
	balance = new(uint256.Int)

	if len(enc) > 0 {
		pos := 0
		nonceBytes := int(enc[pos])
		pos++
		if nonceBytes > 0 {
			nonce = bytesToUint64(enc[pos : pos+nonceBytes])
			pos += nonceBytes
		}
		balanceBytes := int(enc[pos])
		pos++
		if balanceBytes > 0 {
			balance.SetBytes(enc[pos : pos+balanceBytes])
			pos += balanceBytes
		}
		codeHashBytes := int(enc[pos])
		pos++
		if codeHashBytes > 0 {
			codeHash := make([]byte, length.Hash)
			copy(codeHash, enc[pos:pos+codeHashBytes])
		}
	}
	return
}

func EncodeAccountBytes(nonce uint64, balance *uint256.Int, hash []byte, incarnation uint64) []byte {
	l := 1
	if nonce > 0 {
		l += common.BitLenToByteLen(bits.Len64(nonce))
	}
	l++
	if !balance.IsZero() {
		l += balance.ByteLen()
	}
	l++
	if len(hash) == length.Hash {
		l += 32
	}
	l++
	if incarnation > 0 {
		l += common.BitLenToByteLen(bits.Len64(incarnation))
	}
	value := make([]byte, l)
	pos := 0

	if nonce == 0 {
		value[pos] = 0
		pos++
	} else {
		nonceBytes := common.BitLenToByteLen(bits.Len64(nonce))
		value[pos] = byte(nonceBytes)
		var nonce = nonce
		for i := nonceBytes; i > 0; i-- {
			value[pos+i] = byte(nonce)
			nonce >>= 8
		}
		pos += nonceBytes + 1
	}
	if balance.IsZero() {
		value[pos] = 0
		pos++
	} else {
		balanceBytes := balance.ByteLen()
		value[pos] = byte(balanceBytes)
		pos++
		balance.WriteToSlice(value[pos : pos+balanceBytes])
		pos += balanceBytes
	}
	if len(hash) == 0 {
		value[pos] = 0
		pos++
	} else {
		value[pos] = 32
		pos++
		copy(value[pos:pos+32], hash)
		pos += 32
	}
	if incarnation == 0 {
		value[pos] = 0
	} else {
		incBytes := common.BitLenToByteLen(bits.Len64(incarnation))
		value[pos] = byte(incBytes)
		var inc = incarnation
		for i := incBytes; i > 0; i-- {
			value[pos+i] = byte(inc)
			inc >>= 8
		}
	}
	return value
}

func bytesToUint64(buf []byte) (x uint64) {
	for i, b := range buf {
		x = x<<8 + uint64(b)
		if i == 7 {
			return
		}
	}
	return
}
