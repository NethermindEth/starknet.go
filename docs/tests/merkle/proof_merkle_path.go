package main
 
import (
	"fmt"
	"math/big"
 
	"github.com/NethermindEth/starknet.go/merkle"
)
 
func main() {
	// Create a simple Merkle tree structure
	leaf1 := big.NewInt(1)
	leaf2 := big.NewInt(2)
	leaf3 := big.NewInt(3)
	leaf4 := big.NewInt(4)
 
	// Build a simple tree
	hash12 := merkle.MerkleHash(leaf1, leaf2)
	hash34 := merkle.MerkleHash(leaf3, leaf4)
	root := merkle.MerkleHash(hash12, hash34)
 
	// Create proof path for leaf1
	// Path: [leaf2, hash34] to reach root
	path := []*big.Int{leaf2, hash34}
 
	// Verify the proof
	isValid := merkle.ProofMerklePath(root, leaf1, path)
 
	fmt.Println("ProofMerklePath:")
	fmt.Printf("  Root: %s\n", root.String())
	fmt.Printf("  Leaf: %s\n", leaf1.String())
	fmt.Printf("  Path length: %d\n", len(path))
	fmt.Printf("  Proof valid: %v\n", isValid)
}