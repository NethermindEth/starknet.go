package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create two felt values to hash
	a := new(felt.Felt).SetUint64(123)
	b := new(felt.Felt).SetUint64(456)

	// Compute Poseidon hash
	hash := curve.Poseidon(a, b)
	if hash == nil {
		log.Fatal("Failed to compute Poseidon hash")
	}

	fmt.Println("Poseidon Hash:")
	fmt.Printf("  Input a: %s\n", a.String())
	fmt.Printf("  Input b: %s\n", b.String())
	fmt.Printf("  Hash: %s\n", hash.String())
}
