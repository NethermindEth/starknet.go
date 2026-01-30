package rpcv9

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

type ClassOutput interface{}

var (
	_ ClassOutput = (*contracts.DeprecatedContractClass)(nil)
	_ ClassOutput = (*contracts.ContractClass)(nil)
)

type StorageProofInput struct {
	// Required. The hash of the requested block, or number (height) of the
	// requested block, or a block tag
	BlockID types.BlockID `json:"block_id"`
	// Optional. A list of the class hashes for which we want to prove
	// membership in the classes trie
	ClassHashes []*felt.Felt `json:"class_hashes,omitempty"`
	// Optional. A list of contracts for which we want to prove membership in the
	// global state trie
	ContractAddresses []*felt.Felt `json:"contract_addresses,omitempty"`
	// Optional. A list of (contract_address, storage_keys) pairs
	ContractsStorageKeys []ContractStorageKeys `json:"contracts_storage_keys,omitempty"`
}

type ContractStorageKeys struct {
	ContractAddress *felt.Felt         `json:"contract_address"`
	StorageKeys     []types.StorageKey `json:"storage_keys"`
}

// The requested storage proofs. Note that if a requested leaf has the default
// value, the path to it may end in an edge node whose path is not a prefix of
// the requested leaf, thus effectively proving non-membership
type StorageProofResult struct {
	ClassesProof           []NodeHashToNode   `json:"classes_proof"`
	ContractsProof         ContractsProof     `json:"contracts_proof"`
	ContractsStorageProofs [][]NodeHashToNode `json:"contracts_storage_proofs"`
	GlobalRoots            GlobalRoots        `json:"global_roots"`
}

type ContractsProof struct {
	// The nodes in the union of the paths from the contracts tree root to the
	// requested leaves
	Nodes              []NodeHashToNode     `json:"nodes"`
	ContractLeavesData []ContractLeavesData `json:"contract_leaves_data"`
}

// The nonce and class hash for each requested contract address, in the order
// in which they appear in the request. These values are needed to construct the
// associated leaf node
type ContractLeavesData struct {
	Nonce       *felt.Felt `json:"nonce"`
	ClassHash   *felt.Felt `json:"class_hash"`
	StorageRoot *felt.Felt `json:"storage_root,omitempty"`
}

type GlobalRoots struct {
	ContractsTreeRoot *felt.Felt `json:"contracts_tree_root"`
	ClassesTreeRoot   *felt.Felt `json:"classes_tree_root"`
	// the associated block hash (needed in case the caller used a block tag
	// for the block_id parameter)
	BlockHash *felt.Felt `json:"block_hash"`
}

// A node_hash -> node mapping of all the nodes in the union of the paths
// between the requested leaves and the root
type NodeHashToNode struct {
	NodeHash *felt.Felt `json:"node_hash"`
	Node     MerkleNode `json:"node"`
}

// A node in the Merkle-Patricia tree, can be a leaf, binary node, or an edge
// node (EdgeNode or BinaryNode types)
type MerkleNode struct {
	Type string
	Data any
}

// UnmarshalJSON implements the json.Unmarshaler interface for MerkleNode
// It unmarshals the data into an EdgeNode or BinaryNode depending on the type
func (m *MerkleNode) UnmarshalJSON(data []byte) error {
	// Create a decoder with DisallowUnknownFields
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var edgeNode EdgeNode
	if err := decoder.Decode(&edgeNode); err == nil {
		m.Type = "EdgeNode"
		m.Data = edgeNode

		return nil
	}

	// Create a decoder with DisallowUnknownFields
	decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var binaryNode BinaryNode
	if err := decoder.Decode(&binaryNode); err == nil {
		m.Type = "BinaryNode"
		m.Data = binaryNode

		return nil
	}

	return errors.New("invalid merkle node type")
}

// MarshalJSON implements the json.Marshaler interface for MerkleNode
// It marshals the data into an EdgeNode or BinaryNode depending on the type
func (m *MerkleNode) MarshalJSON() ([]byte, error) {
	if m.Type == "EdgeNode" {
		return json.Marshal(m.Data.(EdgeNode))
	}
	if m.Type == "BinaryNode" {
		return json.Marshal(m.Data.(BinaryNode))
	}

	return nil, errors.New("invalid merkle node type")
}

// Represents a path to the highest non-zero descendant node
type EdgeNode struct {
	// an unsigned integer whose binary representation represents the path from
	// the current node to its highest non-zero descendant (bounded by 2^251)
	Path types.NumAsHex `json:"path"`
	// the length of the path (bounded by 251)
	Length uint `json:"length"`
	// the hash of the unique non-zero maximal-height descendant node
	Child *felt.Felt `json:"child"`
}

// An internal node whose both children are non-zero
type BinaryNode struct {
	// the hash of the left child
	Left *felt.Felt `json:"left"`
	// the hash of the right child
	Right *felt.Felt `json:"right"`
}
