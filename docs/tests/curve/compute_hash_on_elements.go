package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create array of big.Int elements
	elem1 := big.NewInt(100)
	elem2 := big.NewInt(200)
	elem3 := big.NewInt(300)
	elems := []*big.Int{elem1, elem2, elem3}

	// Compute hash on elements
	hash := curve.ComputeHashOnElements(elems)
	if hash == nil {
		log.Fatal("Failed to compute hash on elements")
	}

	fmt.Println("ComputeHashOnElements:")
	fmt.Printf("  Input elements: [%s, %s, %s]\n", elem1.String(), elem2.String(), elem3.String())
	fmt.Printf("  Hash: %s\n", hash.String())
}
