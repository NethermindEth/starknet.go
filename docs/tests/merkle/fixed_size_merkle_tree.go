package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/merkle"
)

func main() {
	// Create leaves for the tree
	leaf1 := big.NewInt(10)
	leaf2 := big.NewInt(20)
	leaf3 := big.NewInt(30)
	leaf4 := big.NewInt(40)

	// Create a fixed-size Merkle tree
	tree := merkle.NewFixedSizeMerkleTree(leaf1, leaf2, leaf3, leaf4)

	fmt.Println("FixedSizeMerkleTree:")
	fmt.Printf("  Number of leaves: %d\n", len(tree.Leaves))
	fmt.Printf("  Root: %s\n", tree.Root.String())
	fmt.Printf("  Number of branch levels: %d\n", len(tree.Branches))

	// Generate proof for leaf2
	proof, err := tree.Proof(leaf2)
	if err != nil {
		fmt.Printf("  Error generating proof: %v\n", err)
	} else {
		fmt.Printf("  Proof for leaf 20: %d elements\n", len(proof))
		for i, p := range proof {
			fmt.Printf("    Proof[%d]: %s\n", i, p.String())
		}
	}
}
