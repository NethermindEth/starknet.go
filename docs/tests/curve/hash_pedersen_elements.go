package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create array of big.Int elements
	elem1 := big.NewInt(111)
	elem2 := big.NewInt(222)
	elem3 := big.NewInt(333)
	elems := []*big.Int{elem1, elem2, elem3}

	// Compute Pedersen hash of elements
	hash := curve.HashPedersenElements(elems)
	if hash == nil {
		log.Fatal("Failed to compute Pedersen hash on elements")
	}

	fmt.Println("HashPedersenElements:")
	fmt.Printf("  Input elements: [%s, %s, %s]\n", elem1.String(), elem2.String(), elem3.String())
	fmt.Printf("  Hash: %s\n", hash.String())
}
