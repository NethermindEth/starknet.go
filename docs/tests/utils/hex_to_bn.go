package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert simple hex
	bn := utils.HexToBN("0x123")
	fmt.Printf("0x123 = %s\n", bn.String())
	// Output: 291

	// Convert larger hex
	bn2 := utils.HexToBN("0xabcdef")
	fmt.Printf("0xabcdef = %s\n", bn2.String())
	// Output: 11259375

	// Convert 8-byte hex
	bn3 := utils.HexToBN("0x1234567890abcdef")
	fmt.Printf("0x1234567890abcdef = %s\n", bn3.String())
	// Output: 1311768467294899695
}
