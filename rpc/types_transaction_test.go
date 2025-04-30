package rpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionVersionBigInt(t *testing.T) {
	tests := []struct {
		name     string
		version  TransactionVersion
		expected string
		wantErr  bool
	}{
		{
			name:     "TransactionV0",
			version:  TransactionV0,
			expected: string(TransactionV0),
			wantErr:  false,
		},
		{
			name:     "TransactionV1",
			version:  TransactionV1,
			expected: string(TransactionV1),
			wantErr:  false,
		},
		{
			name:     "TransactionV2",
			version:  TransactionV2,
			expected: string(TransactionV2),
			wantErr:  false,
		},
		{
			name:     "TransactionV3",
			version:  TransactionV3,
			expected: string(TransactionV3),
			wantErr:  false,
		},
		{
			name:     "TransactionV0WithQueryBit",
			version:  TransactionV0WithQueryBit,
			expected: string(TransactionV0WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "TransactionV1WithQueryBit",
			version:  TransactionV1WithQueryBit,
			expected: string(TransactionV1WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "TransactionV2WithQueryBit",
			version:  TransactionV2WithQueryBit,
			expected: string(TransactionV2WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "TransactionV3WithQueryBit",
			version:  TransactionV3WithQueryBit,
			expected: string(TransactionV3WithQueryBit),
			wantErr:  false,
		},
		{
			name:     "InvalidVersion",
			version:  "0xinvalid",
			expected: "-1",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
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
	tests := []struct {
		name     string
		version  TransactionVersion
		expected int
	}{
		{
			name:     "TransactionV0",
			version:  TransactionV0,
			expected: 0,
		},
		{
			name:     "TransactionV1",
			version:  TransactionV1,
			expected: 1,
		},
		{
			name:     "TransactionV2",
			version:  TransactionV2,
			expected: 2,
		},
		{
			name:     "TransactionV3",
			version:  TransactionV3,
			expected: 3,
		},
		{
			name:     "TransactionV0WithQueryBit",
			version:  TransactionV0WithQueryBit,
			expected: 0,
		},
		{
			name:     "TransactionV1WithQueryBit",
			version:  TransactionV1WithQueryBit,
			expected: 1,
		},
		{
			name:     "TransactionV2WithQueryBit",
			version:  TransactionV2WithQueryBit,
			expected: 2,
		},
		{
			name:     "TransactionV3WithQueryBit",
			version:  TransactionV3WithQueryBit,
			expected: 3,
		},
		{
			name:     "InvalidVersion",
			version:  "0xinvalid",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.Int()
			assert.Equal(t, tt.expected, got)
		})
	}
}
