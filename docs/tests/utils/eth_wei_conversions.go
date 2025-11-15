package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("ETH/Wei Conversions:")
	
	// ETH to Wei
	ethValues := []float64{1.0, 0.5, 2.5}
	
	fmt.Println("\nETHToWei:")
	for _, eth := range ethValues {
		wei := utils.ETHToWei(eth)
		fmt.Printf("  Input ETH: %v\n", eth)
		fmt.Printf("  Output Wei: %s\n\n", wei.String())
	}
	
	// Wei to ETH
	weiValues := []string{
		"0xde0b6b3a7640000",   // 1 ETH
		"0x6f05b59d3b20000",   // 0.5 ETH
		"0x22b1c8c1227a00000", // 2.5 ETH
	}
	
	fmt.Println("WeiToETH:")
	for _, weiHex := range weiValues {
		wei, _ := utils.HexToFelt(weiHex)
		eth := utils.WeiToETH(wei)
		fmt.Printf("  Input Wei: %s\n", wei.String())
		fmt.Printf("  Output ETH: %v\n\n", eth)
	}
}
