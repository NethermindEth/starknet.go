package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	bigInt := big.NewInt(98765)
	felt := utils.BigIntToFelt(bigInt)
	fmt.Printf("BigInt: %s\n", bigInt.String())
	fmt.Printf("Felt: %s\n", felt.String())
	// Output:
	// BigInt: 98765
	// Felt: 0x181cd
}
