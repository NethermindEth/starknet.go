package rpcv02

import (
	"encoding/json"

	types "github.com/dontpanicdao/caigo/types"
)

type ResultPageRequest struct {
	// a pointer to the last element of the delivered page, use this token in a subsequent query to obtain the next page
	ContinuationToken *string `json:"continuation_token,omitempty"`
	ChunkSize         int     `json:"chunk_size"`
}

// StorageEntry The changes in the storage of the contract
type StorageEntry struct {
}

// ContractStorageDiffItem is a change in a single storage item
type ContractStorageDiffItem struct {
	// ContractAddress is the contract address for which the state changed
	Address string `json:"address"`

	// Key returns the key of the changed value
	Key string `json:"key"`
	// Value is the new value applied to the given address
	Value string `json:"value"`
}

// DeployedContractItem A new contract deployed as part of the new state
type DeployedContractItem struct {
	// ContractAddress is the address of the contract
	Address string `json:"address"`
	// ClassHash is the hash of the contract code
	ClassHash string `json:"class_hash"`
}

// ContractNonce is a the updated nonce per contract address
type ContractNonce struct {
	// ContractAddress is the address of the contract
	ContractAddress types.Felt `json:"contract_address"`
	// Nonce is the nonce for the given address at the end of the block"
	Nonce string `json:"nonce"`
}

// StateDiff is the change in state applied in this block, given as a
// mapping of addresses to the new values and/or new contracts.
type StateDiff struct {
	// StorageDiffs list storage changes
	StorageDiffs []ContractStorageDiffItem `json:"storage_diffs"`
	// Contracts list new contracts added as part of the new state
	DeclaredContractHashes []string `json:"declared_contract_hashes"`
	// Nonces provides the updated nonces per contract addresses
	DeployedContracts []DeployedContractItem `json:"deployed_contracts"`
	// Nonces provides the updated nonces per contract addresses
	Nonces []ContractNonce `json:"nonces"`
}

type StateUpdateOutput struct {
	// BlockHash is the block identifier,
	BlockHash types.Felt `json:"block_hash"`
	// NewRoot is the new global state root.
	NewRoot string `json:"new_root"`
	// OldRoot is the previous global state root.
	OldRoot string `json:"old_root"`
	// AcceptedTime is when the block was accepted on L1.
	StateDiff StateDiff `json:"state_diff"`
}

// SyncStatus is An object describing the node synchronization status
type SyncStatus struct {
	SyncStatus        bool
	StartingBlockHash string         `json:"starting_block_hash,omitempty"`
	StartingBlockNum  types.NumAsHex `json:"starting_block_num,omitempty"`
	CurrentBlockHash  string         `json:"current_block_hash,omitempty"`
	CurrentBlockNum   types.NumAsHex `json:"current_block_num,omitempty"`
	HighestBlockHash  string         `json:"highest_block_hash,omitempty"`
	HighestBlockNum   types.NumAsHex `json:"highest_block_num,omitempty"`
}

func (s SyncStatus) MarshalJSON() ([]byte, error) {
	if !s.SyncStatus {
		return []byte("false"), nil
	}
	output := map[string]interface{}{}
	output["starting_block_hash"] = s.StartingBlockHash
	output["starting_block_num"] = s.StartingBlockNum
	output["current_block_hash"] = s.CurrentBlockHash
	output["current_block_num"] = s.CurrentBlockNum
	output["highest_block_hash"] = s.HighestBlockHash
	output["highest_block_num"] = s.HighestBlockNum
	return json.Marshal(output)
}

func (s *SyncStatus) UnmarshalJSON(data []byte) error {
	if string(data) == "false" {
		s.SyncStatus = false
		return nil
	}
	s.SyncStatus = true
	output := map[string]interface{}{}
	err := json.Unmarshal(data, &output)
	if err != nil {
		return err
	}
	s.StartingBlockHash = output["starting_block_hash"].(string)
	s.StartingBlockNum = types.NumAsHex(output["starting_block_num"].(string))
	s.CurrentBlockHash = output["current_block_hash"].(string)
	s.CurrentBlockNum = types.NumAsHex(output["current_block_num"].(string))
	s.HighestBlockHash = output["highest_block_hash"].(string)
	s.HighestBlockNum = types.NumAsHex(output["highest_block_num"].(string))
	return nil
}

// AddDeclareTransactionOutput provides the output for AddDeclareTransaction.
type AddDeclareTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
	ClassHash       string `json:"class_hash"`
}

// AddDeployTransactionOutput provides the output for AddDeployTransaction.
type AddDeployTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
	ContractAddress string `json:"contract_address"`
}
