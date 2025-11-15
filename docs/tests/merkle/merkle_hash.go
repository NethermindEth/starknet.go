package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/merkle"
)

func main() {
	// Create two big integers to hash
	x := big.NewInt(100)
	y := big.NewInt(200)

	// Calculate Merkle hash
	hash := merkle.MerkleHash(x, y)

	fmt.Println("MerkleHash:")
	fmt.Printf("  Input X: %s\n", x.String())
	fmt.Printf("  Input Y: %s\n", y.String())
	fmt.Printf("  Merkle Hash: %s\n", hash.String())
}
