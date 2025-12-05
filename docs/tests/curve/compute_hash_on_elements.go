package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create array of big.Int values
	elems := []*big.Int{
		big.NewInt(100),
		big.NewInt(200),
		big.NewInt(300),
	}

	// Compute hash with length prefix
	hash := curve.ComputeHashOnElements(elems)

	fmt.Println("ComputeHashOnElements:")
	fmt.Printf("  Input elements: %v\n", elems)
	fmt.Printf("  Hash: %s\n", hash.String())
}
