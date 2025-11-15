package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	feltValues := []string{
		"0x123",
		"0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a",
		"0x1",
	}

	fmt.Println("FeltToBigInt:")
	for _, hexStr := range feltValues {
		felt, err := utils.HexToFelt(hexStr)
		if err != nil {
			log.Fatal(err)
		}
		bigInt := utils.FeltToBigInt(felt)
		fmt.Printf("  Input Felt: %s\n", felt.String())
		fmt.Printf("  Output BigInt: %s\n\n", bigInt.String())
	}
}
