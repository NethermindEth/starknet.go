package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("HexToBytes:")

	// Convert hex string to bytes
	hexStr := "0x48656c6c6f"
	bytes, err := utils.HexToBytes(hexStr)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("  Input hex: %s\n", hexStr)
	fmt.Printf("  Bytes: %v\n", bytes)
	fmt.Printf("  ASCII: %s\n", string(bytes))
}
