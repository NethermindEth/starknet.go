package utils

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

const ethDecimals = 18

// WeiToETH converts a Wei amount to ETH
// Returns the ETH value as a float64
func WeiToETH(wei *felt.Felt) float64 {
	// Convert to ETH (1 ETH = 10^18 wei)

	// Create a big.Float from the wei amount
	weiFloat := new(big.Float).SetInt(wei.BigInt(new(big.Int)))
	// Create a big.Float for the divisor (10^18)
	divisor := new(big.Float).SetFloat64(math.Pow10(ethDecimals))
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
	multiplier := new(big.Float).SetFloat64(math.Pow10(ethDecimals))
	// Multiply to get Wei
	weiFloat := new(big.Float).Mul(ethFloat, multiplier)

	// Convert to big.Int
	weiInt := new(big.Int)
	weiFloat.Int(weiInt)

	return new(felt.Felt).SetBigInt(weiInt)
}

// FillHexWithZeroes normalizes a hex string to have a '0x' prefix and pads it with leading zeros
// to a total length of 66 characters (including the '0x' prefix).
func FillHexWithZeroes(hex string) string {
	trimHex := strings.TrimPrefix(hex, "0x")
	return strings.Replace(fmt.Sprintf("0x%064s", trimHex), " ", "0", -1)
}
