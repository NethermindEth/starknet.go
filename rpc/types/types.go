package types

// StorageEntry The changes in the storage of the contract
type StorageEntry struct {
	// Key returns the key of the changed value
	Key string `json:"key"`
	// Value is the new value applied to the given address
	Value string `json:"value"`
}

// ContractStorageDiffItem is a change in a single storage item
type ContractStorageDiffItem struct {
	// ContractAddress is the contract address for which the state changed
	Address string `json:"address"`

	// StorageEntry the changes in the storage of the contract
	Entries []StorageEntry `json:"storage_entries"`
}

// DeclaredContractItem A new contract declared as part of the new state
type DeclaredContractItem struct {
	// ClassHash the hash of the contract code
	ClassHash string `json:"class_hash"`
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
	ContractAddress Hash `json:"contract_address"`
	// Nonce is the nonce for the given address at the end of the block"
	Nonce string `json:"nonce"`
}

// StateDiff is the change in state applied in this block, given as a
// mapping of addresses to the new values and/or new contracts.
type StateDiff struct {
	// StorageDiffs list storage changes
	StorageDiffs []ContractStorageDiffItem `json:"storage_diffs"`
	// Contracts list new contracts added as part of the new state
	DeclaredContracts []DeclaredContractItem `json:"declared_contracts"`
	// Nonces provides the updated nonces per contract addresses
	DeployedContracts []DeployedContractItem `json:"deployed_contracts"`
	// Nonces provides the updated nonces per contract addresses
	Nonces []ContractNonce `json:"nonces"`
}

type StateUpdateOutput struct {
	// BlockHash is the block identifier,
	BlockHash Hash `json:"block_hash"`
	// NewRoot is the new global state root.
	NewRoot string `json:"new_root"`
	// OldRoot is the previous global state root.
	OldRoot string `json:"old_root"`
	// AcceptedTime is when the block was accepted on L1.
	AcceptedTime int `json:"accepted_time,omitempty"`
	// StateDiff is the change in state applied in this block, given as a
	// mapping of addresses to the new values and/or new contracts.
	StateDiff StateDiff `json:"state_diff"`
}

type PendingTxnReceipt interface{}

type FeeEstimate struct {
	GasConsumed NumAsHex `json:"gas_consumed"`
	GasPrice    NumAsHex `json:"gas_price"`
	OverallFee  NumAsHex `json:"overall_fee"`
}

func (i BroadcastedInvokeTxnV0) Version() uint64 {
	return 0
}

func (i BroadcastedInvokeTxnV1) Version() uint64 {
	return 1
}

func (s *StructABIEntry) IsType() string {
	return string(s.Type)
}

func (e *EventABIEntry) IsType() string {
	return string(e.Type)
}

func (f *FunctionABIEntry) IsType() string {
	return string(f.Type)
}
