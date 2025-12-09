package main

import (
	"fmt"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert byte slice to big.Int
	bytes := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	bigInt := utils.BytesToBig(bytes)

	fmt.Printf("Bytes: %x\n", bytes)
	fmt.Printf("Big Int: %s\n", bigInt.String())
	// Output:
	// Bytes: 0123456789abcdef
	// Big Int: 81985529216486895
}
