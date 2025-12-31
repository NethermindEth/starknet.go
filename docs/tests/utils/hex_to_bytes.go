package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert hex string to bytes
	hexStr := "0x48656c6c6f20576f726c64"
	bytes, err := utils.HexToBytes(hexStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Hex: %s\n", hexStr)
	fmt.Printf("Bytes: %v\n", bytes)
	fmt.Printf("String: %s\n", string(bytes))
	// Output:
	// Hex: 0x48656c6c6f20576f726c64
	// Bytes: [72 101 108 108 111 32 87 111 114 108 100]
	// String: Hello World
}
