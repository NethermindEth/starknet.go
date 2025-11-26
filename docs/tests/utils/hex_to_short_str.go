package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Decode "Hello"
	str := utils.HexToShortStr("0x48656c6c6f")
	fmt.Printf("0x48656c6c6f = %s\n", str)
	// Output: Hello

	// Decode "World"
	str2 := utils.HexToShortStr("0x576f726c64")
	fmt.Printf("0x576f726c64 = %s\n", str2)
	// Output: World

	// Decode "Starknet"
	str3 := utils.HexToShortStr("0x537461726b6e6574")
	fmt.Printf("0x537461726b6e6574 = %s\n", str3)
	// Output: Starknet
}
