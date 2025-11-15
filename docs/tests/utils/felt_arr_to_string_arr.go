package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("FeltArrToStringArr:")
	
	hexArr := []string{"0x48656c6c6f", "0x576f726c64"}
	feltArr, err := utils.HexArrToFelt(hexArr)
	if err != nil {
		log.Fatal(err)
	}
	
	strArr := utils.FeltArrToStringArr(feltArr)
	
	fmt.Printf("  Input Felts: %v\n", hexArr)
	fmt.Printf("  Output Strings:\n")
	for i, str := range strArr {
		fmt.Printf("    [%d]: %s\n", i, str)
	}
}
