package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create array of big.Int values
	elems := []*big.Int{
		big.NewInt(111),
		big.NewInt(222),
		big.NewInt(333),
	}

	// Compute Pedersen hash
	hash := curve.HashPedersenElements(elems)

	fmt.Println("HashPedersenElements:")
	fmt.Printf("  Input elements: %v\n", elems)
	fmt.Printf("  Hash: %s\n", hash.String())
}
