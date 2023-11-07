package merkle

import (
	"math/big"
	"testing"
)

// debugProof is a function used for debugging purposes. It logs the proofs to the testing logger.
//
// Parameters:
// - t: a pointer to the testing.T object
// - proofs: a slice of pointers to big.Int objects representing the proofs
// Returns:
//   none
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
// - t: A testing.T object used for reporting test failures and logging.
// Returns:
//   none
func TestGeneral_FixedSizeMerkleTree_Check1(t *testing.T) {
	leaves := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5), big.NewInt(6), big.NewInt(7)}
	merkleTree, err := NewFixedSizeMerkleTree(leaves...)
	proof_7_0, _ := MerkleHash(big.NewInt(7), big.NewInt(0))
	proof_1_2, _ := MerkleHash(big.NewInt(1), big.NewInt(2))
	proof_3_4, _ := MerkleHash(big.NewInt(3), big.NewInt(4))
	proof_1_2_3_4, _ := MerkleHash(proof_1_2, proof_3_4)
	manualProof := []*big.Int{
		big.NewInt(6),
		proof_7_0,
		proof_1_2_3_4,
	}
	leaf := big.NewInt(5)
	if err != nil {
		t.Fatal("should generate merkle tree, error", err)
	}
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
