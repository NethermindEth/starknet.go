package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create array of felt values
	felt1 := new(felt.Felt).SetUint64(1)
	felt2 := new(felt.Felt).SetUint64(2)
	felt3 := new(felt.Felt).SetUint64(3)
	felt4 := new(felt.Felt).SetUint64(4)

	// Compute Poseidon hash of array
	hash := curve.PoseidonArray(felt1, felt2, felt3, felt4)
	if hash == nil {
		log.Fatal("Failed to compute PoseidonArray hash")
	}

	fmt.Println("PoseidonArray Hash:")
	fmt.Printf("  Input: [%s, %s, %s, %s]\n", felt1.String(), felt2.String(), felt3.String(), felt4.String())
	fmt.Printf("  Hash: %s\n", hash.String())
}
