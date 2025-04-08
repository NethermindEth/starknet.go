package merkle

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

type FixedSizeMerkleTree struct {
	Leaves   []*big.Int
	Branches [][]*big.Int
	Root     *big.Int
}

// NewFixedSizeMerkleTree creates a new fixed-size Merkle tree.
//
// It takes a variable number of *big.Int leaves as input and returns a pointer to a FixedSizeMerkleTree and an error.
// The function builds the Merkle tree using the given leaves and sets the tree's root.
//
// Parameters:
// - leaves: a slice of *big.Int representing the leaves of the tree.
// Returns:
// - *FixedSizeMerkleTree: a pointer to a FixedSizeMerkleTree
func NewFixedSizeMerkleTree(leaves ...*big.Int) *FixedSizeMerkleTree {
	mt := &FixedSizeMerkleTree{
		Leaves:   leaves,
		Branches: [][]*big.Int{},
	}
	mt.Root = mt.build(leaves)
	return mt
}

// MerkleHash calculates the Merkle hash of two big integers.
//
// Parameters:
// - x: the first big integer
// - y: the second big integer
// Returns:
// - *big.Int: the Merkle hash of the two big integers
func MerkleHash(x, y *big.Int) *big.Int {
	if x.Cmp(y) <= 0 {
		return curve.HashPedersenElements([]*big.Int{x, y})
	}
	return curve.HashPedersenElements([]*big.Int{y, x})
}

// build recursively constructs a Merkle tree from the given leaves.
//
// Parameter(s):
// - leaves: a slice of *big.Int representing the leaves of the tree
// Return type(s):
// - *big.Int: the root hash of the Merkle tree
func (mt *FixedSizeMerkleTree) build(leaves []*big.Int) *big.Int {
	if len(leaves) == 1 {
		return leaves[0]
	}
	mt.Branches = append(mt.Branches, leaves)
	newLeaves := []*big.Int{}
	for i := 0; i < len(leaves); i += 2 {
		if i+1 == len(leaves) {
			hash := MerkleHash(leaves[i], big.NewInt(0))
			newLeaves = append(newLeaves, hash)
			break
		}
		hash := MerkleHash(leaves[i], leaves[i+1])
		newLeaves = append(newLeaves, hash)
	}
	return mt.build(newLeaves)
}

// Proof calculates the Merkle proof for a given leaf in the FixedSizeMerkleTree.
//
// Parameters:
// - leaf: The leaf for which the Merkle proof is calculated
// Returns:
// - []*big.Int: The Merkle proof for the given leaf
// - error: An error if the calculation of the Merkle proof fails
func (mt *FixedSizeMerkleTree) Proof(leaf *big.Int) ([]*big.Int, error) {
	return mt.recursiveProof(leaf, 0, []*big.Int{})
}

// recursiveProof calculates the proof of a leaf in the fixed-size Merkle tree.
//
// Parameters:
// - leaf: is the value to be proven
// - branchIndex: the index of the current branch
// - hashPath: the path from the leaf to the root of the tree
// Returns:
// - []*big.Int: the Merkle proof for the given leaf
// - error: if the key is not found in the branch or there is an error in the nextproof calculation.
func (mt *FixedSizeMerkleTree) recursiveProof(leaf *big.Int, branchIndex int, hashPath []*big.Int) ([]*big.Int, error) {
	if branchIndex >= len(mt.Branches) {
		return hashPath, nil
	}
	branch := mt.Branches[branchIndex]
	index := -1
	for k, v := range branch {
		if v.Cmp(leaf) == 0 {
			index = k
			break
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("key 0x%s not found in branch", leaf.Text(16))
	}
	nextProof := big.NewInt(0)
	if index%2 == 0 && index < len(branch) {
		nextProof = branch[index+1]
	}
	if index%2 != 0 {
		nextProof = branch[index-1]
	}
	newLeaf := MerkleHash(leaf, nextProof)
	newHashPath := append(hashPath, nextProof)
	return mt.recursiveProof(newLeaf, branchIndex+1, newHashPath)
}

// ProofMerklePath checks if a given leaf node is part of a Merkle tree path.
//
// It takes the root node of the Merkle tree, the leaf node to be checked, and the
// path of nodes from the leaf to the root. The function recursively traverses the
// path, verifying each node against the expected hash value. If the path is valid,
// the function returns true; otherwise, it returns false.
//
// Parameters:
// - root: The root node of the Merkle tree as a *big.Int
// - leaf: The leaf node to be checked as a *big.Int
// - path: The path of nodes from the leaf to the root as a slice of *big.Int
// Returns:
// - bool: True if the leaf node is part of the Merkle tree path, false otherwise
func ProofMerklePath(root *big.Int, leaf *big.Int, path []*big.Int) bool {
	if len(path) == 0 {
		return root.Cmp(leaf) == 0
	}
	nexLeaf := MerkleHash(leaf, path[0])

	return ProofMerklePath(root, nexLeaf, path[1:])
}
