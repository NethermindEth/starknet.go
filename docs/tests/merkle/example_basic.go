package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/merkle"
)

func main() {
	// Create leaves
	leaves := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
	}

	// Build the tree
	tree := merkle.NewFixedSizeMerkleTree(leaves...)

	// Display tree information
	fmt.Printf("Merkle Tree Created:\n")
	fmt.Printf("  Number of leaves: %d\n", len(tree.Leaves))
	fmt.Printf("  Root: 0x%s\n", tree.Root.Text(16))
	fmt.Printf("  Branch levels: %d\n\n", len(tree.Branches))

	// Generate proof for leaf value 3
	targetLeaf := big.NewInt(3)
	proof, err := tree.Proof(targetLeaf)
	if err != nil {
		fmt.Printf("Error generating proof: %v\n", err)
		return
	}

	fmt.Printf("Proof for leaf %d:\n", targetLeaf.Int64())
	fmt.Printf("  Proof size: %d hashes\n", len(proof))
	for i, hash := range proof {
		fmt.Printf("  [%d]: 0x%s\n", i, hash.Text(16))
	}

	// Verify the proof
	isValid := merkle.ProofMerklePath(tree.Root, targetLeaf, proof)
	fmt.Printf("\nProof verification: %v\n", isValid)
}
