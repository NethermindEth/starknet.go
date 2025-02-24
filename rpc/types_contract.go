package rpc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

// An integer number in hex format (0x...)
type NumAsHex string

// 64 bit integers, represented by hex string of length at most 16
type U64 string

// ToUint64 converts the U64 type to a uint64.
func (u U64) ToUint64() (uint64, error) {
	hexStr := strings.TrimPrefix(string(u), "0x")

	val, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex string: %v", err)
	}

	return val, nil
}

// 64 bit integers, represented by hex string of length at most 32
type U128 string

type DeprecatedCairoEntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type ClassOutput interface{}

var _ ClassOutput = &DeprecatedContractClass{}
var _ ClassOutput = &ContractClass{}

type ABI []ABIEntry

type DeprecatedEntryPointsByType struct {
	Constructor []DeprecatedCairoEntryPoint `json:"CONSTRUCTOR"`
	External    []DeprecatedCairoEntryPoint `json:"EXTERNAL"`
	L1Handler   []DeprecatedCairoEntryPoint `json:"L1_HANDLER"`
}

type DeprecatedContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	DeprecatedEntryPointsByType DeprecatedEntryPointsByType `json:"entry_points_by_type"`

	ABI *ABI `json:"abi,omitempty"`
}

type NestedString string

func (ns *NestedString) UnmarshalJSON(data []byte) error {
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if value, ok := temp.(string); ok {
		// For cairo compiler prior to 2.7.0, the ABI is a string
		*ns = NestedString(value)
	} else {

		var out bytes.Buffer
		err := json.Indent(&out, data, "", "")

		if err != nil {
			return err
		}

		// Replace '\n' to ''
		out_str := bytes.ReplaceAll(out.Bytes(), []byte{10}, []byte{})
		// Replace ',"' to ', "'
		out_str = bytes.ReplaceAll(out_str, []byte{44, 34}, []byte{44, 32, 34})
		// Replace ',{' to ', {'
		out_str = bytes.ReplaceAll(out_str, []byte{44, 123}, []byte{44, 32, 123})

		*ns = NestedString(out_str)
	}

	return nil
}

type ContractClass struct {
	// The list of Sierra instructions of which the program consists
	SierraProgram []*felt.Felt `json:"sierra_program"`

	// The version of the contract class object. Currently, the Starknet OS supports version 0.1.0
	ContractClassVersion string `json:"contract_class_version"`

	EntryPointsByType SierraEntryPointsByType `json:"entry_points_by_type"`

	ABI NestedString `json:"abi,omitempty"`
}

type StorageProofInput struct {
	// Required. The hash of the requested block, or number (height) of the requested block, or a block tag
	BlockID BlockID `json:"block_id"`
	// Optional. A list of the class hashes for which we want to prove membership in the classes trie
	ClassHashes []*felt.Felt `json:"class_hashes,omitempty"`
	// Optional. A list of contracts for which we want to prove membership in the global state trie
	ContractAddresses []*felt.Felt `json:"contract_addresses,omitempty"`
	// Optional. A list of (contract_address, storage_keys) pairs
	ContractsStorageKeys []ContractStorageKeys `json:"contracts_storage_keys,omitempty"`
}

type ContractStorageKeys struct {
	ContractAddress *felt.Felt   `json:"contract_address"`
	StorageKeys     []*felt.Felt `json:"storage_keys"`
}

// The requested storage proofs. Note that if a requested leaf has the default value,
// the path to it may end in an edge node whose path is not a prefix of the requested leaf,
// thus effecitvely proving non-membership
type StorageProofResult struct {
	ClassesProof           NodeHashToNode   `json:"classes_proof"`
	ContractsProof         ContractsProof   `json:"contracts_proof"`
	ContractsStorageProofs []NodeHashToNode `json:"contracts_storage_proofs"`
	GlobalRoots            []NodeHashToNode `json:"global_roots"`
}

type ContractsProof struct {
	// The nodes in the union of the paths from the contracts tree root to the requested leaves
	Nodes              NodeHashToNode       `json:"nodes"`
	ContractLeavesData []ContractLeavesData `json:"contract_leaves_data"`
}

// The nonce and class hash for each requested contract address, in the order in which
// they appear in the request. These values are needed to construct the associated leaf node
type ContractLeavesData struct {
	Nonce     *felt.Felt `json:"nonce"`
	ClassHash *felt.Felt `json:"class_hash"`
}

type GlobalRoots struct {
	ContractsTreeRoot *felt.Felt `json:"contracts_tree_root"`
	ClassesTreeRoot   *felt.Felt `json:"classes_tree_root"`
	// the associated block hash (needed in case the caller used a block tag for the block_id parameter)
	BlockHash *felt.Felt `json:"block_hash"`
}

// A node_hash -> node mapping of all the nodes in the union of the paths between the requested leaves and the root
type NodeHashToNode struct {
	NodeHash *felt.Felt `json:"node_hash"`
	Node     MerkleNode `json:"node"`
}

// A node in the Merkle-Patricia tree, can be a leaf, binary node, or an edge node
type MerkleNode struct {
	EdgeNode   `json:",omitempty"`
	BinaryNode `json:",omitempty"`
}

// Represents a path to the highest non-zero descendant node
type EdgeNode struct {
	// an integer whose binary representation represents the path from the current node to its highest non-zero descendant (bounded by 2^251)
	Path NumAsHex `json:"path"`
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

// UnmarshalJSON unmarshals the JSON content into the DeprecatedContractClass struct.
//
// It takes a byte array `content` as a parameter and returns an error if there is any.
// The function processes the `program` field in the JSON object.
// If `program` is a string, it is assigned to the `Program` field in the struct.
// Otherwise, it is encoded and assigned to the `Program` field.
// The function then processes the `entry_points_by_type` field in the JSON object.
// The value is unmarshaled into the `DeprecatedEntryPointsByType` field in the struct.
// Finally, the function processes the `abi` field in the JSON object.
// If it is missing, the function returns nil.
// Otherwise, it unmarshals the value into a slice of interfaces.
// For each element in the slice, it checks the type and assigns it to the appropriate field in the `ABI` field in the struct.
//
// Parameters:
// - content: byte array
// Returns:
// - error: error if there is any
func (c *DeprecatedContractClass) UnmarshalJSON(content []byte) error {
	v := map[string]json.RawMessage{}
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	// process 'program'. If it is a string, keep it, otherwise encode it.
	data, ok := v["program"]
	if !ok {
		return fmt.Errorf("missing program in json object")
	}
	program := ""
	if err := json.Unmarshal(data, &program); err != nil {
		if program, err = encodeProgram(data); err != nil {
			return err
		}
	}
	c.Program = program

	// process 'entry_points_by_type'
	data, ok = v["entry_points_by_type"]
	if !ok {
		return fmt.Errorf("missing entry_points_by_type in json object")
	}

	depEntryPointsByType := DeprecatedEntryPointsByType{}
	if err := json.Unmarshal(data, &depEntryPointsByType); err != nil {
		return err
	}
	c.DeprecatedEntryPointsByType = depEntryPointsByType

	// process 'abi'
	data, ok = v["abi"]
	if !ok {
		// contractClass can have an empty ABI for instance with ClassAt
		return nil
	}

	abis := []interface{}{}
	if err := json.Unmarshal(data, &abis); err != nil {
		return err
	}

	abiPointer := ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab ABIEntry
			abiType, ok := checkABI["type"].(string)
			if !ok {
				return fmt.Errorf("unknown abi type %v", checkABI["type"])
			}
			switch abiType {
			case string(ABITypeConstructor), string(ABITypeFunction), string(ABITypeL1Handler):
				ab = &FunctionABIEntry{}
			case string(ABITypeStruct):
				ab = &StructABIEntry{}
			case string(ABITypeEvent):
				ab = &EventABIEntry{}
			default:
				return fmt.Errorf("unknown ABI type %v", checkABI["type"])
			}
			data, err := json.Marshal(checkABI)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, ab)
			if err != nil {
				return err
			}
			abiPointer = append(abiPointer, ab)
		}
	}

	c.ABI = &abiPointer
	return nil
}

func (nodeHashToNode *NodeHashToNode) UnmarshalJSON(bytes []byte) error {
	valueMap := make(map[string]any)
	if err := json.Unmarshal(bytes, &valueMap); err != nil {
		return err
	}

	nodeHash, ok := valueMap["node_hash"]
	if !ok {
		return fmt.Errorf("missing 'node_hash' in json object")
	}
	nodeHashFelt, ok := nodeHash.(felt.Felt)
	if !ok {
		return fmt.Errorf("error casting 'node_hash' to felt.Felt")
	}

	node, ok := valueMap["node"]
	if !ok {
		return fmt.Errorf("missing 'node' in json object")
	}
	var merkleNode MerkleNode
	switch nodeT := node.(type) {
	case BinaryNode:
		merkleNode = MerkleNode{BinaryNode: nodeT}
	case EdgeNode:
		merkleNode = MerkleNode{EdgeNode: nodeT}
	default:
		return fmt.Errorf("'node' should be an EdgeNode or BinaryNode")
	}

	*nodeHashToNode = NodeHashToNode{
		NodeHash: &nodeHashFelt,
		Node:     merkleNode,
	}

	return nil
}

type SierraEntryPoint struct {
	// The index of the function in the program
	FunctionIdx int `json:"function_idx"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type SierraEntryPointsByType struct {
	Constructor []SierraEntryPoint `json:"CONSTRUCTOR"`
	External    []SierraEntryPoint `json:"EXTERNAL"`
	L1Handler   []SierraEntryPoint `json:"L1_HANDLER"`
}

type ABIEntry interface {
	IsType() ABIType
}

type ABIType string

const (
	ABITypeConstructor ABIType = "constructor"
	ABITypeFunction    ABIType = "function"
	ABITypeL1Handler   ABIType = "l1_handler"
	ABITypeEvent       ABIType = "event"
	ABITypeStruct      ABIType = "struct"
)

type StructABIEntry struct {
	// The event type
	Type ABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	// todo(minumum size should be 1)
	Size uint64 `json:"size"`

	Members []Member `json:"members"`
}

type Member struct {
	TypedParameter
	Offset int64 `json:"offset"`
}

type EventABIEntry struct {
	// The event type
	Type ABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Keys []TypedParameter `json:"keys"`

	Data []TypedParameter `json:"data"`
}

type FunctionStateMutability string

const (
	FuncStateMutVIEW FunctionStateMutability = "view"
)

type FunctionABIEntry struct {
	// The function type
	Type ABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	StateMutability FunctionStateMutability `json:"stateMutability,omitempty"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

// IsType returns the ABIType of the StructABIEntry.
//
// Parameters:
//
//	none
//
// Returns:
// - ABIType: the ABIType
func (s *StructABIEntry) IsType() ABIType {
	return s.Type
}

// IsType returns the ABIType of the EventABIEntry.
//
// Parameters:
//
//	none
//
// Returns:
// - ABIType: the ABIType
func (e *EventABIEntry) IsType() ABIType {
	return e.Type
}

// IsType returns the ABIType of the FunctionABIEntry.
//
// Parameters:
//
//	none
//
// Returns:
// - ABIType: the ABIType
func (f *FunctionABIEntry) IsType() ABIType {
	return f.Type
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}

// encodeProgram encodes the content byte array using gzip compression and base64 encoding.
//
// It takes a content byte array as a parameter and returns the encoded program string and an error.
//
// Parameters:
// - content: byte array to be encoded
// Returns:
// - string: the encoded program
// - error: the error if any
func encodeProgram(content []byte) (string, error) {
	buf := bytes.NewBuffer(nil)
	gzipContent := gzip.NewWriter(buf)
	_, err := gzipContent.Write(content)
	if err != nil {
		return "", err
	}
	gzipContent.Close()
	program := base64.StdEncoding.EncodeToString(buf.Bytes())
	return program, nil
}
