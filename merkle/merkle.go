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
// If there is an error during the tree building process, the function returns nil and the error.
func NewFixedSizeMerkleTree(leaves ...*big.Int) (*FixedSizeMerkleTree, error) {
	mt := &FixedSizeMerkleTree{
		Leaves:   leaves,
		Branches: [][]*big.Int{},
	}
	root, err := mt.build(leaves)
	if err != nil {
		return nil, err
	}
	mt.Root = root
	return mt, err
}

// MerkleHash calculates the Merkle hash of two big integers.
//
// Parameters:
// - x: the first big integer
// - y: the second big integer
//
// Returns:
// - the Merkle hash of the two big integers
// - an error if the calculation fails
func MerkleHash(x, y *big.Int) (*big.Int, error) {
	if x.Cmp(y) <= 0 {
		return curve.Curve.HashElements([]*big.Int{x, y})
	}
	return curve.Curve.HashElements([]*big.Int{y, x})
}

// build recursively constructs a Merkle tree from the given leaves.
//
// Parameter(s):
// - leaves: a slice of *big.Int representing the leaves of the tree.
//
// Return type(s):
// - *big.Int: the root hash of the Merkle tree.
// - error: any error that occurred during the construction of the tree.
func (mt *FixedSizeMerkleTree) build(leaves []*big.Int) (*big.Int, error) {
	if len(leaves) == 1 {
		return leaves[0], nil
	}
	mt.Branches = append(mt.Branches, leaves)
	newLeaves := []*big.Int{}
	for i := 0; i < len(leaves); i += 2 {
		if i+1 == len(leaves) {
			hash, err := MerkleHash(leaves[i], big.NewInt(0))
			if err != nil {
				return nil, err
			}
			newLeaves = append(newLeaves, hash)
			break
		}
		hash, err := MerkleHash(leaves[i], leaves[i+1])
		if err != nil {
			return nil, err
		}
		newLeaves = append(newLeaves, hash)
	}
	return mt.build(newLeaves)
}

// Proof calculates the Merkle proof for a given leaf in the FixedSizeMerkleTree.
//
// Parameters:
// - leaf: The leaf for which the Merkle proof is calculated.
//
// Returns:
// - []*big.Int: The Merkle proof for the given leaf.
// - error: An error if the calculation of the Merkle proof fails.
func (mt *FixedSizeMerkleTree) Proof(leaf *big.Int) ([]*big.Int, error) {
	return mt.recursiveProof(leaf, 0, []*big.Int{})
}

// recursiveProof calculates the proof of a leaf in the fixed-size Merkle tree.
//
// It takes a leaf, branch index, and a hash path as input parameters.
// The leaf is the value to be proven.
// The branch index is the index of the current branch.
// The hash path is the path from the leaf to the root of the tree.
//
// It returns the hash path and an error if any.
// The hash path is the updated path from the leaf to the root of the tree,
// including the proofs for the intermediate nodes.
// The error is returned if the key is not found in the branch or there is an error in the nextproof calculation.
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
	newLeaf, err := MerkleHash(leaf, nextProof)
	if err != nil {
		return nil, fmt.Errorf("nextproof error: %v", err)
	}
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
// - root (*big.Int): The root node of the Merkle tree.
// - leaf (*big.Int): The leaf node to be checked.
// - path ([]*big.Int): The path of nodes from the leaf to the root.
//
// Returns:
// - bool: True if the leaf node is part of the Merkle tree path, false otherwise.
func ProofMerklePath(root *big.Int, leaf *big.Int, path []*big.Int) bool {
	if len(path) == 0 {
		return root.Cmp(leaf) == 0
	}
	nexLeaf, err := MerkleHash(leaf, path[0])
	if err != nil {
		return false
	}
	return ProofMerklePath(root, nexLeaf, path[1:])
}
