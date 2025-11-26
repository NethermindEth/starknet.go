package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert small number
	felt := utils.Uint64ToFelt(123)
	fmt.Printf("123 = %s\n", felt.String())
	// Output: 0x7b

	// Convert large number
	felt2 := utils.Uint64ToFelt(999999)
	fmt.Printf("999999 = %s\n", felt2.String())
	// Output: 0xf423f

	// Convert max uint64
	felt3 := utils.Uint64ToFelt(18446744073709551615)
	fmt.Printf("Max uint64 = %s\n", felt3.String())
	// Output: 0xffffffffffffffff
}
