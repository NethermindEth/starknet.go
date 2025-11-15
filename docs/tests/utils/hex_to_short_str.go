package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("HexToShortStr:")
	
	hexValues := []string{
		"0x48656c6c6f",        // "Hello"
		"0x576f726c64",        // "World"
		"0x537461726b6e6574",  // "Starknet"
	}
	
	for _, hex := range hexValues {
		str := utils.HexToShortStr(hex)
		fmt.Printf("  Input: %s\n", hex)
		fmt.Printf("  Output: %s\n\n", str)
	}
}
