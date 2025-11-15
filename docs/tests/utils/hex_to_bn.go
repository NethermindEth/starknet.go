package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("HexToBN:")
	
	hexValues := []string{
		"0x123",
		"0xabcdef",
		"0x1234567890abcdef",
	}
	
	for _, hex := range hexValues {
		bn := utils.HexToBN(hex)
		fmt.Printf("  Input: %s\n", hex)
		fmt.Printf("  Output: %s\n\n", bn.String())
	}
}
