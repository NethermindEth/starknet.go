package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("StrToHex:")
	
	strings := []string{"Hello", "World", "Starknet"}
	
	for _, str := range strings {
		hex := utils.StrToHex(str)
		fmt.Printf("  Input: %s\n", str)
		fmt.Printf("  Output: %s\n\n", hex)
	}
}
