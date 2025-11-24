package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Compute Starknet Keccak hash
	input := []byte("Hello Starknet")
	hash := curve.StarknetKeccak(input)

	fmt.Println("StarknetKeccak Hash:")
	fmt.Printf("  Input: %s\n", string(input))
	fmt.Printf("  Hash: %s\n", hash.String())
}
