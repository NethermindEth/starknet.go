package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

type ResultPageRequest struct {
	// a pointer to the last element of the delivered page, use this token in a subsequent query to obtain the next page
	ContinuationToken *string `json:"continuation_token,omitempty"`
	ChunkSize         int     `json:"chunk_size"`
}

type StorageEntry struct {
	Key   *felt.Felt `json:"key"`
	Value *felt.Felt `json:"value"`
}

// type StorageEntries struct {
// 	StorageEntry []StorageEntry
// }

// ContractStorageDiffItem is a change in a single storage item
type ContractStorageDiffItem struct {
	// ContractAddress is the contract address for which the state changed
	Address        *felt.Felt     `json:"address"`
	StorageEntries []StorageEntry `json:"storage_entries"`
}

// DeployedContractItem A new contract deployed as part of the new state
type DeployedContractItem struct {
	// ContractAddress is the address of the contract
	Address *felt.Felt `json:"address"`
	// ClassHash is the hash of the contract code
	ClassHash *felt.Felt `json:"class_hash"`
}

// ContractNonce is a the updated nonce per contract address
type ContractNonce struct {
	// ContractAddress is the address of the contract
	ContractAddress *felt.Felt `json:"contract_address"`
	// Nonce is the nonce for the given address at the end of the block"
	Nonce *felt.Felt `json:"nonce"`
}

// StateDiff is the change in state applied in this block, given as a
// mapping of addresses to the new values and/or new contracts.
type StateDiff struct {
	// StorageDiffs list storage changes
	StorageDiffs []ContractStorageDiffItem `json:"storage_diffs"`
	// Contracts list new contracts added as part of the new state
	DeclaredContractHashes []*felt.Felt `json:"declared_contract_hashes"`
	// Nonces provides the updated nonces per contract addresses
	DeployedContracts []DeployedContractItem `json:"deployed_contracts"`
	// Nonces provides the updated nonces per contract addresses
	Nonces []ContractNonce `json:"nonces"`
}

// STATE_UPDATE in spec
type StateUpdateOutput struct {
	// BlockHash is the block identifier,
	BlockHash *felt.Felt `json:"block_hash"`
	// NewRoot is the new global state root.
	NewRoot *felt.Felt `json:"new_root"`
	// OldRoot is the previous global state root.
	OldRoot *felt.Felt `json:"old_root"`
	// AcceptedTime is when the block was accepted on L1.
	StateDiff StateDiff `json:"state_diff"`
}

// SyncStatus is An object describing the node synchronization status
type SyncStatus struct {
	SyncStatus        bool       // todo(remove? not in spec)
	StartingBlockHash *felt.Felt `json:"starting_block_hash,omitempty"`
	StartingBlockNum  NumAsHex   `json:"starting_block_num,omitempty"`
	CurrentBlockHash  *felt.Felt `json:"current_block_hash,omitempty"`
	CurrentBlockNum   NumAsHex   `json:"current_block_num,omitempty"`
	HighestBlockHash  *felt.Felt `json:"highest_block_hash,omitempty"`
	HighestBlockNum   NumAsHex   `json:"highest_block_num,omitempty"`
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
	return json.Unmarshal(data, s)

	// if string(data) == "false" {
	// 	s.SyncStatus = false
	// 	return nil
	// }
	// s.SyncStatus = true
	// output := map[string]interface{}{}
	// err := json.Unmarshal(data, &output)
	// if err != nil {
	// 	return err
	// }
	// s.StartingBlockHash = output["starting_block_hash"].(string)
	// s.StartingBlockNum = types.NumAsHex(output["starting_block_num"].(string))
	// s.CurrentBlockHash = output["current_block_hash"].(string)
	// s.CurrentBlockNum = types.NumAsHex(output["current_block_num"].(string))
	// s.HighestBlockHash = output["highest_block_hash"].(string)
	// s.HighestBlockNum = types.NumAsHex(output["highest_block_num"].(string))
	// return nil
}

// AddDeclareTransactionOutput provides the output for AddDeclareTransaction.
type AddDeclareTransactionOutput struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ClassHash       *felt.Felt `json:"class_hash"`
}

// AddDeployTransactionOutput provides the output for AddDeployTransaction.
type AddDeployTransactionOutput struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ContractAddress *felt.Felt `json:"contract_address"`
}

// AddDeployAccountTransactionOutput provides the output for AddDeployTransaction.
type AddDeployAccountTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
}

// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    *felt.Felt `json:"contract_address"`
	EntryPointSelector *felt.Felt `json:"entry_point_selector,omitempty"`

	// Calldata The parameters passed to the function
	Calldata []*felt.Felt `json:"calldata"`
}

type FeeEstimate struct {
	// GasConsumed the Ethereum gas cost of the transaction (see https://docs.starknet.io/docs/Fees/fee-mechanism for more info)
	GasConsumed NumAsHex `json:"gas_consumed"`

	// GasPrice the gas price (in gwei) that was used in the cost estimation
	GasPrice NumAsHex `json:"gas_price"`

	// OverallFee the estimated fee for the transaction (in gwei), product of gas_consumed and gas_price
	OverallFee NumAsHex `json:"overall_fee"`
}

type TransactionState string

const (
	TransactionAcceptedOnL1 TransactionState = "ACCEPTED_ON_L1"
	TransactionAcceptedOnL2 TransactionState = "ACCEPTED_ON_L2"
	TransactionNotReceived  TransactionState = "NOT_RECEIVED"
	TransactionPending      TransactionState = "PENDING"
	TransactionReceived     TransactionState = "RECEIVED"
	TransactionRejected     TransactionState = "REJECTED"
)

func (ts *TransactionState) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "ACCEPTED_ON_L2":
		*ts = TransactionAcceptedOnL2
	case "ACCEPTED_ON_L1":
		*ts = TransactionAcceptedOnL1
	case "NOT_RECEIVED":
		*ts = TransactionNotReceived
	case "PENDING":
		*ts = TransactionPending
	case "RECEIVED":
		*ts = TransactionReceived
	case "REJECTED":
		*ts = TransactionRejected
	default:
		return fmt.Errorf("unsupported status: %s", data)
	}
	return nil
}

func (ts TransactionState) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(ts))), nil
}

func (s TransactionState) String() string {
	return string(s)
}

func (s TransactionState) IsTransactionFinal() bool {
	if s == TransactionAcceptedOnL2 ||
		s == TransactionAcceptedOnL1 ||
		s == TransactionRejected {
		return true
	}
	return false
}
