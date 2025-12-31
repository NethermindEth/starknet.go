package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/merkle"
)

func main() {
	a := big.NewInt(5)
	b := big.NewInt(10)

	// Hash in both orders
	hash1 := merkle.MerkleHash(a, b)
	hash2 := merkle.MerkleHash(b, a)

	fmt.Printf("Value A: %d\n", a.Int64())
	fmt.Printf("Value B: %d\n\n", b.Int64())

	fmt.Printf("Hash(A, B): 0x%s\n", hash1.Text(16))
	fmt.Printf("Hash(B, A): 0x%s\n\n", hash2.Text(16))

	// Verify they are equal
	fmt.Printf("Hashes are equal: %v\n", hash1.Cmp(hash2) == 0)
	fmt.Println("\nThis deterministic ordering ensures MerkleHash(a,b) always equals MerkleHash(b,a)")
}
