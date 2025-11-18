package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("Uint64ToFelt:")
	
	values := []uint64{0, 1, 123, 999999, 18446744073709551615}
	
	for _, val := range values {
		felt := utils.Uint64ToFelt(val)
		fmt.Printf("  Input: %d\n", val)
		fmt.Printf("  Output: %s\n\n", felt.String())
	}
}
