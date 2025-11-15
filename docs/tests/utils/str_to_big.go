package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("StrToBig:")

	// Convert numeric string to big.Int
	numStr := "999999999999999999999999"
	bigInt := utils.StrToBig(numStr)

	fmt.Printf("  Input string: \"%s\"\n", numStr)
	fmt.Printf("  Big Int: %s\n", bigInt.String())
}
