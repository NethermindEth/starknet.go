package types

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
)

func TestAddressStr(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		input    *felt.Felt
		expected string
	}{
		{
			name:     "Zero value",
			input:    &felt.Zero,
			expected: "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:     "Small value",
			input:    new(felt.Felt).SetUint64(123),
			expected: "0x000000000000000000000000000000000000000000000000000000000000007b",
		},
		{
			name:     "Medium value",
			input:    func() *felt.Felt { f, _ := new(felt.Felt).SetString("0xabc123"); return f }(),
			expected: "0x0000000000000000000000000000000000000000000000000000000000abc123",
		},
		{
			name:     "Address-like value",
			input:    func() *felt.Felt { f, _ := new(felt.Felt).SetString("0x06bb9425718d801fd06f144abb82eced725f0e81db61d2f9f4c9a26ece46a829"); return f }(),
			expected: "0x06bb9425718d801fd06f144abb82eced725f0e81db61d2f9f4c9a26ece46a829",
		},
		{
			name:     "Edge case - already 66 characters",
			input:    func() *felt.Felt { f, _ := new(felt.Felt).SetString("0x00000000000000000000000000000000000000000000000000000000000000ff"); return f }(),
			expected: "0x00000000000000000000000000000000000000000000000000000000000000ff",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			feltWrapper := NewFelt(tc.input)
			formatted := feltWrapper.AddressStr()
			assert.Equal(t, tc.expected, formatted, "The formatted address should match the expected value")
			assert.Equal(t, 66, len(formatted), "The formatted address should be exactly 66 characters long")
		})
	}
}
