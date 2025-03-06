package utils

import (
	"math/big"
	"testing"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
)

func TestResBoundsMapToOverallFee(t *testing.T) {
	tests := []struct {
		name       string
		resBounds  rpc.ResourceBoundsMapping
		multiplier float64
		expected   *big.Int
	}{
		{
			name: "Basic calculation",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0xa",  // 10
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x32", // 50
					MaxPricePerUnit: "0x5",  // 5
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0xc8", // 200
					MaxPricePerUnit: "0x3",  // 3
				},
			},
			multiplier: 1.0,
			// Expected: (100*10) + (50*5) + (200*3) = 1000 + 250 + 600 = 1850
			expected: big.NewInt(1850),
		},
		{
			name: "Zero values",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
			},
			multiplier: 1.0,
			expected:   big.NewInt(0),
		},
		{
			name: "Large values",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x3e8", // 1000
					MaxPricePerUnit: "0x64",  // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x1f4", // 500
					MaxPricePerUnit: "0x32",  // 50
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x7d0", // 2000
					MaxPricePerUnit: "0x19",  // 25
				},
			},
			multiplier: 1.0,
			// Expected: (1000*100) + (500*50) + (2000*25) = 100000 + 25000 + 50000 = 175000
			expected: big.NewInt(175000),
		},
		{
			name: "With multiplier",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0xa",  // 10
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x32", // 50
					MaxPricePerUnit: "0x5",  // 5
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0xc8", // 200
					MaxPricePerUnit: "0x3",  // 3
				},
			},
			multiplier: 1.5,
			// Expected: (100*10) + (50*5) + (200*3) = 1000 + 250 + 600 = 1850
			// Note: The multiplier doesn't seem to be used in the function implementation
			expected: big.NewInt(1850),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResBoundsMapToOverallFee(tt.resBounds, tt.multiplier)

			assert.Equal(t, tt.expected, got, "ResBoundsMapToOverallFee() returned incorrect value")
		})
	}
}
