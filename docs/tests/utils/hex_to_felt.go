package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert hex strings to Felt
	felt, err := utils.HexToFelt("0x123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output: %s\n", felt.String())
	// Output: 0x123

	// Convert contract address
	felt2, err := utils.HexToFelt("0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output: %s\n", felt2.String())
	// Output: 0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a
}
