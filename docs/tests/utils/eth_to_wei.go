package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert ETH amounts to Wei
	ethAmounts := []float64{1.0, 0.5, 0.001, 2.5}

	fmt.Println("ETHToWei:")
	for _, eth := range ethAmounts {
		wei := utils.ETHToWei(eth)
		fmt.Printf("  ETH: %.4f\n", eth)
		fmt.Printf("  Wei: %s\n", wei.String())
		fmt.Println()
	}
}
