package utils

import (
	"math/big"
	"testing"
)

func TestWeiToETH(t *testing.T) {
	tests := []struct {
		name     string
		wei      *big.Int
		expected float64
	}{
		{
			name:     "1 ETH",
			wei:      new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 10^18 wei = 1 ETH
			expected: 1.0,
		},
		{
			name:     "0 ETH",
			wei:      big.NewInt(0),
			expected: 0.0,
		},
		{
			name:     "Small amount",
			wei:      big.NewInt(1000000000000000), // 0.001 ETH
			expected: 0.001,
		},
		{
			name:     "Large amount",
			wei:      new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 1000 ETH
			expected: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeiToETH(tt.wei)

			// For floating point comparison, use a small delta
			delta := 0.0000001
			if diff := abs(got - tt.expected); diff > delta {
				t.Errorf("WeiToETH() = %v, want %v, diff %v", got, tt.expected, diff)
			}
		})
	}
}

func TestETHToWei(t *testing.T) {
	tests := []struct {
		name     string
		eth      float64
		expected *big.Int
	}{
		{
			name:     "1 ETH",
			eth:      1.0,
			expected: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 10^18 wei = 1 ETH
		},
		{
			name:     "0 ETH",
			eth:      0.0,
			expected: big.NewInt(0),
		},
		{
			name:     "Small amount",
			eth:      0.001,
			expected: big.NewInt(1000000000000000), // 0.001 ETH
		},
		{
			name:     "Large amount",
			eth:      1000.0,
			expected: new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 1000 ETH
		},
		{
			name: "Fractional amount",
			eth:  1.5,
			expected: new(big.Int).Add(
				new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),                                   // 1 ETH
				new(big.Int).Div(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), big.NewInt(2))), // 0.5 ETH
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ETHToWei(tt.eth)

			// For big.Int comparison, we can use Cmp
			if got.Cmp(tt.expected) != 0 {
				t.Errorf("ETHToWei() = %v, want %v", got.String(), tt.expected.String())
			}
		})
	}
}

// Helper function to get absolute value of float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
