package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Sample data to hash
	data := []byte("Hello Starknet")

	// Compute Starknet Keccak hash
	hash := curve.StarknetKeccak(data)

	fmt.Println("StarknetKeccak Hash:")
	fmt.Printf("  Input: %s\n", string(data))
	fmt.Printf("  Hash: %s\n", hash.String())
}
