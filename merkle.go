package starknetgo

import (
	"fmt"
	"math/big"
)

type FixedSizeMerkleTree struct {
	Leaves   []*big.Int
	Branches [][]*big.Int
	Root     *big.Int
}

// NewFixedSizeMerkleTree creates a new fixed-size Merkle tree.
//
// It takes a variable number of *big.Int leaves as input and returns a
// pointer to a FixedSizeMerkleTree and an error. The function constructs
// a FixedSizeMerkleTree object with the given leaves and an empty set of
// branches. It then builds the Merkle tree and sets the root of the tree.
// If there is an error during the build process, it returns nil and the error.
//
// Parameters:
//   - leaves: A variable number of *big.Int leaves representing the data
//     to be stored in the Merkle tree.
//
// Returns:
//   - *FixedSizeMerkleTree: A pointer to the constructed FixedSizeMerkleTree
//     object.
//   - error: An error object if there was an error during the build process.
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
// It takes in two parameters, x and y, which are pointers to big.Int.
// It returns a pointer to big.Int and an error.
func MerkleHash(x, y *big.Int) (*big.Int, error) {
	if x.Cmp(y) <= 0 {
		return Curve.HashElements([]*big.Int{x, y})
	}
	return Curve.HashElements([]*big.Int{y, x})
}

// build builds the FixedSizeMerkleTree.
//
// It takes in a slice of *big.Int leaves as a parameter.
// It returns a *big.Int and an error.
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

// Proof calculates the Merkle tree proof for a given leaf.
//
// Parameters:
// - leaf: The leaf for which the proof needs to be calculated.
//
// Returns:
// - A slice of big.Int values representing the Merkle tree proof.
// - An error if the proof calculation fails.
func (mt *FixedSizeMerkleTree) Proof(leaf *big.Int) ([]*big.Int, error) {
	return mt.recursiveProof(leaf, 0, []*big.Int{})
}

// recursiveProof calculates the proof of a leaf in a fixed-size Merkle tree.
//
// Parameters:
// - leaf: the leaf value to calculate the proof for.
// - branchIndex: the index of the branch to start the proof from.
// - hashPath: the array of hash values representing the Merkle path.
//
// Returns:
// - []*big.Int: the proof values for the leaf.
// - error: an error if the proof calculation fails.
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

// ProofMerklePath checks if a given leaf node belongs to a Merkle tree with the given root and path.
//
// Parameters:
// - root: a pointer to a big.Int representing the root of the Merkle tree.
// - leaf: a pointer to a big.Int representing the leaf node to be checked.
// - path: a slice of pointers to big.Int representing the path from the leaf to the root.
//
// Returns:
// - a boolean indicating whether the leaf node belongs to the Merkle tree.
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
