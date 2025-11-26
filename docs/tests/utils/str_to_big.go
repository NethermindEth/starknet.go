package main

import (
	"fmt"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert numeric string to big.Int
	numStr := "123456789012345678901234567890"
	bigInt := utils.StrToBig(numStr)

	fmt.Printf("String: %s\n", numStr)
	fmt.Printf("Big Int: %s\n", bigInt.String())
	fmt.Printf("Type: %T\n", bigInt)
	// Output:
	// String: 123456789012345678901234567890
	// Big Int: 123456789012345678901234567890
	// Type: *big.Int
}
