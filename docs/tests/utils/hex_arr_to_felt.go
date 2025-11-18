package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("HexArrToFelt:")
	
	hexArr := []string{"0x1", "0x2", "0x3", "0xabc", "0xdef"}
	
	feltArr, err := utils.HexArrToFelt(hexArr)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("  Input: %v\n", hexArr)
	fmt.Printf("  Output:\n")
	for i, felt := range feltArr {
		fmt.Printf("    [%d]: %s\n", i, felt.String())
	}
}
