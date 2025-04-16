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
