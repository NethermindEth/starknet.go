package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert 1 ETH to Wei
	wei := utils.ETHToWei(1.0)
	fmt.Printf("1 ETH = %s Wei\n", wei.String())
	// Output: 0xde0b6b3a7640000

	// Convert 0.5 ETH to Wei
	wei2 := utils.ETHToWei(0.5)
	fmt.Printf("0.5 ETH = %s Wei\n", wei2.String())
	// Output: 0x6f05b59d3b20000

	// Convert 2.5 ETH to Wei
	wei3 := utils.ETHToWei(2.5)
	fmt.Printf("2.5 ETH = %s Wei\n", wei3.String())
	// Output: 0x22b1c8c1227a0000
}
