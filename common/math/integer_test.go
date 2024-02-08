// Copyright 2017 The go-ethereum Authors
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

package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type operation byte

const (
	sub operation = iota
	add
	mul
)

func TestHexOrDecimal64(t *testing.T) {
	tests := []struct {
		input string
		num   uint64
		ok    bool
	}{
		{"", 0, true},
		{"0", 0, true},
		{"0x0", 0, true},
		{"12345678", 12345678, true},
		{"0x12345678", 0x12345678, true},
		{"0X12345678", 0x12345678, true},
		// Tests for leading zero behaviour:
		{"0123456789", 123456789, true}, // note: not octal
		{"0x00", 0, true},
		{"0x012345678abc", 0x12345678abc, true},
		// Invalid syntax:
		{"abcdef", 0, false},
		{"0xgg", 0, false},
		// Doesn't fit into 64 bits:
		{"18446744073709551617", 0, false},
	}
	for _, test := range tests {
		var num HexOrDecimal64
		err := num.UnmarshalText([]byte(test.input))
		if (err == nil) != test.ok {
			t.Errorf("ParseUint64(%q) -> (err == nil) = %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if err == nil && uint64(num) != test.num {
			t.Errorf("ParseUint64(%q) -> %d, want %d", test.input, num, test.num)
		}
	}
}

func TestMustParseUint64(t *testing.T) {
	if v := MustParseUint64("12345"); v != 12345 {
		t.Errorf(`MustParseUint64("12345") = %d, want 12345`, v)
	}
}

func TestMustParseUint64Panic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustParseBig should've panicked")
		}
	}()
	MustParseUint64("ggg")
}

func TestAbsoluteDifference(t *testing.T) {
	x1 := uint64(99)
	x2 := uint64(45)
	assert.Equal(t, AbsoluteDifference(x1, x2), x1-x2)
	assert.Equal(t, AbsoluteDifference(x2, x1), x1-x2)
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		number   uint64
		expected bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{13, true},
		{25, false},
		{29, true},
		{7919, true}, // testing with a larger prime
		{7920, false},
		{9223372036854775807, true}, // largest 64-bit prime number
	}

	for _, test := range tests {
		if result := IsPrime(test.number); result != test.expected {
			t.Errorf("IsPrime(%d) = %v, want %v", test.number, result, test.expected)
		}
	}
}
