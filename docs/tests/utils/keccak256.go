package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("Keccak256:")
	
	// Test with simple strings
	data1 := []byte("hello")
	hash1 := utils.Keccak256(data1)
	fmt.Printf("  Input: %s\n", string(data1))
	fmt.Printf("  Output: %x\n\n", hash1)
	
	data2 := []byte("transfer")
	hash2 := utils.Keccak256(data2)
	fmt.Printf("  Input: %s\n", string(data2))
	fmt.Printf("  Output: %x\n\n", hash2)
}
