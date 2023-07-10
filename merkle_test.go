package starknetgo

import (
	"math/big"
	"testing"
)

func debugProof(t *testing.T, proofs []*big.Int) {
	t.Log("...proof")
	for k, v := range proofs {
		t.Logf("key[%d] 0x%s\n", k, v.Text(16))
	}
}

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
