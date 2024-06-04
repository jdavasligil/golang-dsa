package bitset_test

import (
	"testing"

	"github.com/jdavasligil/golang-dsa/bitset"
)

func TestBitset8(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		tests := []struct {
			b        bitset.BitSet8
			idx      uint8
			expected bitset.BitSet8
		}{
			{b: 0b00000000, idx: 0, expected: 0b00000001},
			{b: 0b00000001, idx: 0, expected: 0b00000001},
			{b: 0b00000000, idx: 255, expected: 0b00000000},
			{b: 0b00000000, idx: 7, expected: 0b10000000},
		}

		for i, test := range tests {
			test.b.Set(test.idx)
			if test.b != test.expected {
				t.Errorf("Test %d failed. Expected %v - Got %v", i, test.expected, test.b)
			}
		}
	})
	t.Run("Unset", func(t *testing.T) {
		tests := []struct {
			b        bitset.BitSet8
			idx      uint8
			expected bitset.BitSet8
		}{
			{b: 0b01000000, idx: 0, expected: 0b01000000},
			{b: 0b00000001, idx: 0, expected: 0b00000000},
			{b: 0b00000000, idx: 255, expected: 0b00000000},
			{b: 0b10000001, idx: 7, expected: 0b00000001},
		}

		for i, test := range tests {
			test.b.Unset(test.idx)
			if test.b != test.expected {
				t.Errorf("Test %d failed. Expected %v - Got %v", i, test.expected, test.b)
			}
		}
	})
	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			b        bitset.BitSet8
			idx      uint8
			expected bool
		}{
			{b: 0b01000000, idx: 0, expected: false},
			{b: 0b00000001, idx: 0, expected: true},
			{b: 0b00000000, idx: 255, expected: false},
			{b: 0b10000001, idx: 7, expected: true},
		}

		for i, test := range tests {
			if test.b.Get(test.idx) != test.expected {
				t.Errorf("Test %d failed. Expected %v - Got %v", i, test.expected, test.b)
			}
		}
	})
}
