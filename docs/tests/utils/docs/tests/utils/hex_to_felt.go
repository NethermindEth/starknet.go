package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	hexValues := []string{
		"0x123",
		"0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a",
		"0x0",
		"0x1",
	}

	fmt.Println("HexToFelt:")
	for _, hex := range hexValues {
		felt, err := utils.HexToFelt(hex)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Input: %s\n", hex)
		fmt.Printf("  Output: %s\n\n", felt.String())
	}
}
