package merkle

import (
	"math/big"
	"testing"
)

// debugProof is a function used for debugging purposes. It logs the proofs to the testing logger.
//
// Parameters:
//   - t: a pointer to the testing.T object
//   - proofs: a slice of pointers to big.Int objects representing the proofs
//
// Returns:
//
//	none
func debugProof(t *testing.T, proofs []*big.Int) {
	t.Log("...proof")
	for k, v := range proofs {
		t.Logf("key[%d] 0x%s\n", k, v.Text(16))
	}
}

// TestGeneral_FixedSizeMerkleTree_Check1 is a Go function that tests the functionality of the FixedSizeMerkleTree.Check1 method.
//
// It creates a fixed-size Merkle tree with the given leaves and calculates the Merkle proof for a specific leaf. It then compares the manual proof generated with the expected proof and checks if the Merkle tree root matches the proof.
//
// Parameters:
//   - t: A testing.T object used for reporting test failures and logging.
//
// Returns:
//
//	none
//
//nolint:staticcheck // Best readability
func TestGeneral_FixedSizeMerkleTree_Check1(t *testing.T) {
	leaves := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
		big.NewInt(5),
		big.NewInt(6),
		big.NewInt(7),
	}
	merkleTree := NewFixedSizeMerkleTree(leaves...)
	proof_7_0 := MerkleHash(big.NewInt(7), big.NewInt(0))
	proof_1_2 := MerkleHash(big.NewInt(1), big.NewInt(2))
	proof_3_4 := MerkleHash(big.NewInt(3), big.NewInt(4))
	proof_1_2_3_4 := MerkleHash(proof_1_2, proof_3_4)
	manualProof := []*big.Int{
		big.NewInt(6),
		proof_7_0,
		proof_1_2_3_4,
	}
	leaf := big.NewInt(5)
	proof, err := merkleTree.Proof(leaf)
	if err != nil {
		t.Fatal("should generate merkle proof, error", err)
	}
	if len(manualProof) != len(proof) {
		debugProof(t, proof)
		t.Fatalf("tree length should match, expected: %d, got: %d", len(manualProof), len(proof))
	}
	for i, p := range manualProof {
		if p.Cmp(proof[i]) != 0 {
			t.Fatalf("proof should match, expected: 0x%s, got: 0x%s", p.Text(16), proof[i].Text(16))
		}
	}
	if ok := ProofMerklePath(merkleTree.Root, leaf, proof); !ok {
		t.Logf("MerkleTree Root, 0x%s\n", merkleTree.Root.Text(16))
		t.Fatal("root should match proof. it does not")
	}
}

// TestRecursiveProofFixedSizeMerkleTree is a Go function that tests the correctness of the recursive proof generation in the FixedSizeMerkleTree.
//
// It creates a Merkle tree with multiple leaves, selects a specific leaf, and generates a Merkle proof using the recursiveProof method.
// The test then reconstructs the Merkle root using the generated proof and verifies that it matches the original root of the Merkle tree.
//
// Parameters:
//   - t: A testing.T object used for reporting test failures and logging.
//
// Returns:
//
//	none
func TestRecursiveProofFixedSizeMerkleTree(t *testing.T) {
	// Create a Merkle tree with multiple leaves
	leaves := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)}
	merkleTree := NewFixedSizeMerkleTree(leaves...)

	// Choose a leaf for which to generate the Merkle proof
	targetLeaf := leaves[2] // Replace with the desired leaf

	// Generate the Merkle proof using the recursiveProof method
	proof, err := merkleTree.recursiveProof(targetLeaf, 0, []*big.Int{})
	if err != nil {
		t.Fatalf("Error generating Merkle proof: %v", err)
	}

	// Verify the correctness of the generated proof
	reconstructedRoot := reconstructRootFromProof(targetLeaf, proof)

	// Verify that the reconstructed root matches the original root
	if merkleTree.Root.Cmp(reconstructedRoot) != 0 {
		t.Fatalf(
			"Reconstructed Merkle root does not match the original root. Expected: 0x%s, Got: 0x%s",
			merkleTree.Root.Text(16),
			reconstructedRoot.Text(16),
		)
	}
}

// reconstructRootFromProof is a helper function that reconstructs the Merkle
// root from a given root, leaf, and Merkle proof.
func reconstructRootFromProof(leaf *big.Int, proof []*big.Int) *big.Int {
	currentHash := leaf
	for _, sibling := range proof {
		currentHash = MerkleHash(currentHash, sibling)
	}

	return currentHash
}
