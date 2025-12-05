package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	hexArr := []string{"0x1", "0x2", "0x3", "0xabc", "0xdef"}

	feltArr, err := utils.HexArrToFelt(hexArr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Input: %v\n", hexArr)
	fmt.Println("Output:")
	for i, felt := range feltArr {
		fmt.Printf("  [%d]: %s\n", i, felt.String())
	}
	// Output:
	// [0]: 0x1
	// [1]: 0x2
	// [2]: 0x3
	// [3]: 0xabc
	// [4]: 0xdef
}
