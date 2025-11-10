package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Create Wei amounts
	wei1, _ := new(felt.Felt).SetString("1000000000000000000") // 1 ETH
	wei2, _ := new(felt.Felt).SetString("500000000000000000")  // 0.5 ETH
	wei3, _ := new(felt.Felt).SetString("1000000000000000")    // 0.001 ETH

	fmt.Println("WeiToETH:")
	
	eth1 := utils.WeiToETH(wei1)
	fmt.Printf("  Wei: %s\n", wei1.String())
	fmt.Printf("  ETH: %.18f\n", eth1)
	fmt.Println()
	
	eth2 := utils.WeiToETH(wei2)
	fmt.Printf("  Wei: %s\n", wei2.String())
	fmt.Printf("  ETH: %.18f\n", eth2)
	fmt.Println()
	
	eth3 := utils.WeiToETH(wei3)
	fmt.Printf("  Wei: %s\n", wei3.String())
	fmt.Printf("  ETH: %.18f\n", eth3)
}
