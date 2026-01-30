package rpcv9

import (
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/rpc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionVersionBigInt(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testSet := []struct {
		name     string
		version  types.TransactionVersion
		expected string
		wantErr  bool
	}{
		{
			name:     "TransactionV0",
			version:  types.TransactionV0,
			expected: string(types.TransactionV0),
			wantErr:  false,
		},
		{
			name:     "TransactionV1",
			version:  types.TransactionV1,
			expected: string(types.TransactionV1),
			wantErr:  false,
		},
		{
			name:     "TransactionV2",
			version:  types.TransactionV2,
			expected: string(types.TransactionV2),
			wantErr:  false,
		},
		{
			name:     "TransactionV3",
			version:  types.TransactionV3,
			expected: string(types.TransactionV3),
			wantErr:  false,
		},
		{
			name:     "TransactionV0WithQueryBit",
			version:  types.TransactionV0WithQueryBit,
			expected: string(types.TransactionV0WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "TransactionV1WithQueryBit",
			version:  types.TransactionV1WithQueryBit,
			expected: string(types.TransactionV1WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "TransactionV2WithQueryBit",
			version:  types.TransactionV2WithQueryBit,
			expected: string(types.TransactionV2WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "TransactionV3WithQueryBit",
			version:  types.TransactionV3WithQueryBit,
			expected: string(types.TransactionV3WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "InvalidVersion",
			version:  "0xinvalid",
			expected: "-1",
			wantErr:  true,
		},
	}

	for _, tt := range testSet {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.version.BigInt()
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.expected, got.String())

				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.expected, fmt.Sprintf("%#x", got))
		})
	}
}

func TestTransactionVersionInt(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testSet := []struct {
		name     string
		version  types.TransactionVersion
		expected int
	}{
		{
			name:     "TransactionV0",
			version:  types.TransactionV0,
			expected: 0,
		},
		{
			name:     "TransactionV1",
			version:  types.TransactionV1,
			expected: 1,
		},
		{
			name:     "TransactionV2",
			version:  types.TransactionV2,
			expected: 2,
		},
		{
			name:     "TransactionV3",
			version:  types.TransactionV3,
			expected: 3,
		},
		{
			name:     "TransactionV0WithQueryBit",
			version:  types.TransactionV0WithQueryBit,
			expected: 0,
		},
		{
			name:     "TransactionV1WithQueryBit",
			version:  types.TransactionV1WithQueryBit,
			expected: 1,
		},
		{
			name:     "TransactionV2WithQueryBit",
			version:  types.TransactionV2WithQueryBit,
			expected: 2,
		},
		{
			name:     "TransactionV3WithQueryBit",
			version:  types.TransactionV3WithQueryBit,
			expected: 3,
		},
		{
			name:     "InvalidVersion",
			version:  "0xinvalid",
			expected: -1,
		},
	}

	for _, tt := range testSet {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.Int()
			assert.Equal(t, tt.expected, got)
		})
	}
}
