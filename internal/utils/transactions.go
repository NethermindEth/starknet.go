package utils

import (
	"math"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

// WeiToETH converts a Wei amount to ETH
// Returns the ETH value as a float64
func WeiToETH(wei *felt.Felt) float64 {
	// Convert to ETH (1 ETH = 10^18 wei)

	// Create a big.Float from the wei amount
	weiFloat := new(big.Float).SetInt(wei.BigInt(new(big.Int)))
	// Create a big.Float for the divisor (10^18)
	divisor := new(big.Float).SetFloat64(math.Pow10(18))
	// Divide to get ETH
	ethFloat := new(big.Float).Quo(weiFloat, divisor)
	// Convert to float64
	ethValue, _ := ethFloat.Float64()

	return ethValue
}

// ETHToWei converts an ETH amount to Wei
// Returns the Wei value as a *felt.Felt
func ETHToWei(eth float64) *felt.Felt {
	// Convert to Wei (1 ETH = 10^18 wei)

	// Create a big.Float from the eth amount
	ethFloat := new(big.Float).SetFloat64(eth)
	// Create a big.Float for the multiplier (10^18)
	multiplier := new(big.Float).SetFloat64(math.Pow10(18))
	// Multiply to get Wei
	weiFloat := new(big.Float).Mul(ethFloat, multiplier)

	// Convert to big.Int
	weiInt := new(big.Int)
	weiFloat.Int(weiInt)

	return new(felt.Felt).SetBigInt(weiInt)
}
