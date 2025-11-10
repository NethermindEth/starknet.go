package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("BigToHex:")
	
	bigValues := []*big.Int{
		big.NewInt(123),
		big.NewInt(11259375),
		new(big.Int).SetUint64(1234567890),
	}
	
	for _, bigInt := range bigValues {
		hex := utils.BigToHex(bigInt)
		fmt.Printf("  Input: %s\n", bigInt.String())
		fmt.Printf("  Output: %s\n\n", hex)
	}
}
