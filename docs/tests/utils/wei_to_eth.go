package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert 1 ETH in Wei to ETH
	wei, _ := utils.HexToFelt("0xde0b6b3a7640000")
	eth := utils.WeiToETH(wei)
	fmt.Printf("%s Wei = %v ETH\n", wei.String(), eth)
	// Output: 1 ETH

	// Convert 0.5 ETH in Wei to ETH
	wei2, _ := utils.HexToFelt("0x6f05b59d3b20000")
	eth2 := utils.WeiToETH(wei2)
	fmt.Printf("%s Wei = %v ETH\n", wei2.String(), eth2)
	// Output: 0.5 ETH
}
