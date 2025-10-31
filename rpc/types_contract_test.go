package rpc_test

import (
	"math/big"
	"testing"

	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
)

// TestU128_ToBigInt tests the ToBigInt method of the U128 type.
func TestU128_ToBigInt(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		u128    rpc.U128
		want    *big.Int
		wantErr bool
	}{
		{
			name: "within the range",
			u128: "0xabcdef",
			want: internalUtils.HexToBN("0xabcdef"),
		},
		{
			name: "max uint128",
			u128: "0xffffffffffffffffffffffffffffffff",
			want: internalUtils.HexToBN("0xffffffffffffffffffffffffffffffff"),
		},
		{
			name:    "out of range",
			u128:    "0x100000000000000000000000000000000",
			wantErr: true,
		},
		{
			name:    "invalid hex string",
			u128:    "56yrty45",
			wantErr: true,
		},
		{
			name:    "without 0x prefix",
			u128:    "abcdef",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.u128
			got, gotErr := u.ToBigInt()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToBigInt() failed: %v", gotErr)
				}

				return
			}
			if tt.wantErr {
				t.Fatal("ToBigInt() succeeded unexpectedly")
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
