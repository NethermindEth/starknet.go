package rpc_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
)

// TestU128_ToBigInt tests the ToBigInt method of the U128 type.
func TestU128_ToBigInt(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testCases := []struct {
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
	}
	for _, tt := range testCases {
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

// TestU128_ToUint64 tests the ToUint64 method of the U128 type.
func TestU128_ToUint64(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testCases := []struct {
		name    string // description of this test case
		u64     rpc.U64
		want    uint64
		wantErr bool
	}{
		{
			name: "within the range",
			u64:  "0xabcdef",
			want: 11259375,
		},
		{
			name: "max uint64",
			u64:  "0xFFFFFFFFFFFFFFFF",
			want: math.MaxUint64,
		},
		{
			name:    "out of range",
			u64:     "0x10000000000000000",
			wantErr: true,
		},
		{
			name:    "invalid hex string",
			u64:     "56yrty45",
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.u64
			got, gotErr := u.ToUint64()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToUint64() failed: %v", gotErr)
				}

				return
			}
			if tt.wantErr {
				t.Fatal("ToUint64() succeeded unexpectedly")
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
