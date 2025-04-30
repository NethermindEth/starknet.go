package utils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToByteArrFelt(t *testing.T) {
	var tests = []struct {
		in  string
		out []string
	}{
		{
			in:  "hello",
			out: []string{"0x0", "0x68656c6c6f", "0x5"},
		},
		{
			in:  "Long string, more than 31 characters.",
			out: []string{"0x1", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x63746572732e", "0x6"},
		},
		{
			in:  "Blockchain secure digital asset",
			out: []string{"0x1", "0x426c6f636b636861696e20736563757265206469676974616c206173736574", "0x0", "0x0"},
		},
		{
			in:  "Decentralized applications offer transparency and user control",
			out: []string{"0x2", "0x446563656e7472616c697a6564206170706c69636174696f6e73206f666665", "0x72207472616e73706172656e637920616e64207573657220636f6e74726f6c", "0x0", "0x0"},
		},
		{
			in:  "12345",
			out: []string{"0x0", "0x3132333435", "0x5"},
		},
		{
			in:  "1234567890123456789012345678901",
			out: []string{"0x1", "0x31323334353637383930313233343536373839303132333435363738393031", "0x0", "0x0"},
		},
		{
			in:  "12345678901234567890123456789012",
			out: []string{"0x1", "0x31323334353637383930313233343536373839303132333435363738393031", "0x32", "0x1"},
		},
	}

	for _, tc := range tests {
		res, err := StringToByteArrFelt(tc.in)
		require.NoError(t, err, "error returned from StringToByteArrFelt")
		require.Len(t, res, len(tc.out), "invalid conversion: array sizes do not match")

		out, err := HexArrToFelt(tc.out)
		require.NoError(t, err, "error returned from HexArrToFelt")
		require.Exactly(t, out, res, "invalid conversion: values do not match")
	}
}

func TestByteArrFeltToString(t *testing.T) {
	var tests = []struct {
		in  []string
		out string
	}{
		{
			in:  []string{"0x0", "0x68656c6c6f", "0x5"},
			out: "hello",
		},
		{
			in:  []string{"0x1", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x63746572732e", "0x6"},
			out: "Long string, more than 31 characters.",
		},
		{
			in:  []string{"0x1", "0x426c6f636b636861696e20736563757265206469676974616c206173736574", "0x0", "0x0"},
			out: "Blockchain secure digital asset",
		},
		{
			in:  []string{"0x2", "0x446563656e7472616c697a6564206170706c69636174696f6e73206f666665", "0x72207472616e73706172656e637920616e64207573657220636f6e74726f6c", "0x0", "0x0"},
			out: "Decentralized applications offer transparency and user control",
		},
		{
			in:  []string{"0x0", "0x3132333435", "0x5"},
			out: "12345",
		},
		{
			in:  []string{"0x1", "0x31323334353637383930313233343536373839303132333435363738393031", "0x0", "0x0"},
			out: "1234567890123456789012345678901",
		},
		{
			in:  []string{"0x1", "0x31323334353637383930313233343536373839303132333435363738393031", "0x32", "0x1"},
			out: "12345678901234567890123456789012",
		},
		{
			in:  []string{"0x0", "0x66656c741", "0x5"},
			out: "\x06\x66\x56\xc7\x41",
		},
	}

	for _, tc := range tests {
		in, err := HexArrToFelt(tc.in)
		require.NoError(t, err, "error returned from HexArrToFelt")
		res, err := ByteArrFeltToString(in)
		require.NoError(t, err, "error returned from ByteArrFeltToString")
		require.Equal(t, tc.out, res, "invalid conversion: output does not match")
	}
}

func TestHexToU256Felt(t *testing.T) {
	var tests = []struct {
		name       string
		hexInput   string
		wantLow    string
		wantHigh   string
		shouldFail bool
	}{
		{
			name:     "simple decimal 2",
			hexInput: "0x2",
			wantLow:  "0x2",
			wantHigh: "0x0",
		},
		{
			name:     "2^128",
			hexInput: "0x100000000000000000000000000000000",
			wantLow:  "0x0",
			wantHigh: "0x1",
		},
		{
			name:     "2^129 + 2^128 + 20",
			hexInput: "0x300000000000000000000000000000014",
			wantLow:  "0x14",
			wantHigh: "0x3",
		},
		{
			name:     "max uint128 in low part",
			hexInput: "0xffffffffffffffffffffffffffffffff",
			wantLow:  "0xffffffffffffffffffffffffffffffff",
			wantHigh: "0x0",
		},
		{
			name:     "max uint128 in high part",
			hexInput: "0xffffffffffffffffffffffffffffffff00000000000000000000000000000000",
			wantLow:  "0x0",
			wantHigh: "0xffffffffffffffffffffffffffffffff",
		},
		{
			name:     "hex without 0x prefix",
			hexInput: "abcdef",
			wantLow:  "0xabcdef",
			wantHigh: "0x0",
		},
		{
			name:       "invalid hex string",
			hexInput:   "0xZZZ",
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := HexToU256Felt(tc.hexInput)

			if tc.shouldFail {
				require.Error(t, err, "expected error but got none")
				return
			}

			require.NoError(t, err, "unexpected error")
			require.Len(t, result, 2, "result should contain exactly 2 felt values")

			// Convert the result felts to hex strings for comparison
			lowHex := result[0].String()
			highHex := result[1].String()

			assert.Equal(t, tc.wantLow, lowHex, "low bits do not match expected value")
			assert.Equal(t, tc.wantHigh, highHex, "high bits do not match expected value")
		})
	}
}

func TestU256FeltToHex(t *testing.T) {
	var tests = []struct {
		name       string
		lowHex     string
		highHex    string
		wantHex    string
		shouldFail bool
	}{
		{
			name:    "simple decimal 2",
			lowHex:  "0x2",
			highHex: "0x0",
			wantHex: "0x2",
		},
		{
			name:    "2^128",
			lowHex:  "0x0",
			highHex: "0x1",
			wantHex: "0x100000000000000000000000000000000",
		},
		{
			name:    "2^129 + 2^128 + 20",
			lowHex:  "0x14",
			highHex: "0x3",
			wantHex: "0x300000000000000000000000000000014",
		},
		{
			name:    "max uint128 in low part",
			lowHex:  "0xffffffffffffffffffffffffffffffff",
			highHex: "0x0",
			wantHex: "0xffffffffffffffffffffffffffffffff",
		},
		{
			name:    "max uint128 in high part",
			lowHex:  "0x0",
			highHex: "0xffffffffffffffffffffffffffffffff",
			wantHex: "0xffffffffffffffffffffffffffffffff00000000000000000000000000000000",
		},
		{
			name:    "zero value",
			lowHex:  "0x0",
			highHex: "0x0",
			wantHex: "0x0",
		},
		{
			name:       "invalid input length",
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var input []*felt.Felt

			if tc.shouldFail {
				// Test with empty slice for error case
				result, err := U256FeltToHex(input)
				assert.Error(t, err, "expected error but got none")
				assert.Empty(t, result, "expected empty result on error")
				return
			}

			// Parse the low and high hex values into felt.Felt objects
			lowFelt, err := HexToFelt(tc.lowHex)
			require.NoError(t, err, "failed to parse low hex")

			highFelt, err := HexToFelt(tc.highHex)
			require.NoError(t, err, "failed to parse high hex")

			input = []*felt.Felt{lowFelt, highFelt}

			// Call the function under test
			result, err := U256FeltToHex(input)
			require.NoError(t, err, "unexpected error")

			// Normalize expected value for comparison
			// Convert hex to big.Int and back to normalized hex for consistent comparison
			expectedBigInt, _ := new(big.Int).SetString(tc.wantHex[2:], 16)
			expectedHex := fmt.Sprintf("%#x", expectedBigInt)

			assert.Equal(t, expectedHex, result, "hex representation does not match expected value")

			// Round-trip test: Convert back to u256 felt representation
			roundTrip, err := HexToU256Felt(result)
			require.NoError(t, err, "failed in round-trip conversion")
			assert.Equal(t, lowFelt.String(), roundTrip[0].String(), "round-trip low bits do not match")
			assert.Equal(t, highFelt.String(), roundTrip[1].String(), "round-trip high bits do not match")
		})
	}
}
