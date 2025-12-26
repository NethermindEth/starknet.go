package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Encode "Hello"
	hex := utils.StrToHex("Hello")
	fmt.Printf("Hello = %s\n", hex)
	// Output: 0x48656c6c6f

	// Encode "World"
	hex2 := utils.StrToHex("World")
	fmt.Printf("World = %s\n", hex2)
	// Output: 0x576f726c64

	// Encode "Starknet"
	hex3 := utils.StrToHex("Starknet")
	fmt.Printf("Starknet = %s\n", hex3)
	// Output: 0x537461726b6e6574
}
