package rpc

import (
	"errors"
)

var (
	errBadRequest      = errors.New("bad request")
	errBadTxType       = errors.New("bad transaction type")
	errInvalidBlockTag = errors.New("invalid blocktag")
	errNotImplemented  = errors.New("not implemented")
)

// BlockHashAndNumberOutput is a struct that is returned by BlockHashAndNumber.
type BlockHashAndNumberOutput struct {
	BlockNumber uint64 `json:"block_number,omitempty"`
	BlockHash   string `json:"block_hash,omitempty"`
}

// blockID is an unexposed struct that is used in a OneOf for
// starknet_getBlockWithTxHashes.
type blockID struct {
	BlockNumber *BlockNumber `json:"block_number,omitempty"`
	BlockHash   *BlockHash   `json:"block_hash,omitempty"`
	BlockTag    *string      `json:"block_tag,omitempty"`
}

// BlockID is an options that can be used as a parameter for
// starknet_getBlockWithTxHashes.
type BlockID func(b *blockID) error

func (bid BlockID) isValid() bool {
	b := &blockID{}
	err := bid(b)
	if err != nil {
		return false
	}
	if b.BlockTag != nil &&
		(*b.BlockTag != "pending" && *b.BlockTag != "latest") {
		return false
	}
	return true
}

func (bid BlockID) tag() (string, bool) {
	b := &blockID{}
	bid(b)
	if b.BlockTag != nil &&
		(*b.BlockTag == "pending" || *b.BlockTag == "latest") {
		return *b.BlockTag, true
	}
	return "", false
}

func (bid BlockID) isPending() bool {
	if v, ok := bid.tag(); ok {
		if v == "pending" {
			return true
		}
	}
	return false
}

func (bid BlockID) isLatest() bool {
	if v, ok := bid.tag(); ok {
		if v == "latest" {
			return true
		}
	}
	return false
}

func (bid BlockID) getWithoutTag() (*blockID, error) {
	b := &blockID{}
	bid(b)
	if b.BlockTag != nil {
		return nil, errors.New("blockid is a tag")
	}
	return b, nil
}

type BlockHash string

type BlockNumber uint64

type PendingBlockWithTxHashes struct {
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`

	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress string `json:"sequencer_address"`

	// ParentHash The hash of this block's parent
	ParentHash BlockHash `json:"parent_hash"`

	BlockBodyWithTxHashes
}

type BlockStatus string

// TxnHash a transaction's hash
type TxnHash string

// BlockBodyWithTxHashes the hashes of the transactions included in this block.
type BlockBodyWithTxHashes struct {
	// Transactions The hashes of the transactions included in this block
	Transactions []TxnHash `json:"transactions"`
}

type BlockHeader struct {
	// BlockHash The hash of this block
	BlockHash BlockHash `json:"block_hash"`

	// ParentHash The hash of this block's parent
	ParentHash BlockHash `json:"parent_hash"`

	// BlockNumber the block number (its height)
	BlockNumber BlockNumber `json:"block_number"`

	// NewRoot The new global state root
	NewRoot string `json:"new_root"`

	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`

	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress string `json:"sequencer_address"`
}

// BlockWithTxHashes The block object
type BlockWithTxHashes struct {
	Status BlockStatus `json:"status"`
	BlockHeader
	BlockBodyWithTxHashes
}

type TxnType string

type NumAsHex string

type Signature []string

// BroadcastedCommonTxnProperties common properties of a transaction that is sent to the sequencer (but is not yet in a block)
type BroadcastedCommonTxnProperties struct {
	Type TxnType `json:"type,omitempty"`

	// MaxFee maximal fee that can be charged for including the transaction
	MaxFee string `json:"max_fee"`

	// Version of the transaction scheme
	Version NumAsHex `json:"version"`

	// Signature
	Signature Signature `json:"signature"`

	// Nonce
	Nonce string `json:"nonce,omitempty"`
}

type CommonTxnProperties struct {
	TransactionHash TxnHash `json:"transaction_hash"`
	BroadcastedCommonTxnProperties
}

type Address string

// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    Address `json:"contract_address"`
	EntryPointSelector string  `json:"entry_point_selector"`

	// CallData The parameters passed to the function
	CallData []string `json:"calldata"`
}

// InvokeV0 version 0 invoke transaction
type InvokeV0 FunctionCall

// InvokeV1 version 1 invoke transaction
type InvokeV1 struct {
	SenderAddress Address `json:"sender_address"`
	// CallData The parameters passed to the function
	CallData []string `json:"calldata"`
}

// InvokeTxnDuck is a type used to understand the Invoke Version
type InvokeTxnDuck struct {
	AccountAddress     Address `json:"account_address"`
	ContractAddress    Address `json:"contract_address"`
	EntryPointSelector string  `json:"entry_point_selector"`
}

type InvokeTxnV0 struct {
	CommonTxnProperties
	InvokeV0
}

type InvokeTxnV1 struct {
	CommonTxnProperties
	InvokeV1
}

type InvokeTxn interface{}

type L1HandlerTxn struct {
	// TransactionHash The hash identifying the transaction
	TransactionHash TxnHash

	Type TxnType `json:"type"`

	// MaxFee maximal fee that can be charged for including the transaction
	MaxFee string `json:"max_fee"`

	// Version of the transaction scheme
	Version NumAsHex `json:"version"`

	// Signature
	Signature Signature `json:"signature"`

	// Nonce
	Nonce string `json:"nonce"`
}

type DeclareTxn struct {
	CommonTxnProperties

	// ClassHash the hash of the declared class
	ClassHash string `json:"class_hash"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress string `json:"sender_address"`
}

// DeployTxn The structure of a deploy transaction. Note that this transaction type is deprecated and will no longer be supported in future versions
type DeployTxn struct {
	// TransactionHash The hash identifying the transaction
	TransactionHash TxnHash `json:"transaction_hash"`

	// ClassHash The hash of the deployed contract's class
	ClassHash string `json:"class_hash"`

	// Version of the transaction scheme
	Version NumAsHex `json:"version"`

	Type TxnType `json:"type"`

	// ContractAddress The address of the deployed contract
	ContractAddress string `json:"contract_address"`

	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt string `json:"contract_address_salt"`

	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []string `json:"constructor_calldata"`
}

type Txn interface{}

// BlockBodyWithTxs the hashes of the transactions included in this block.
type BlockBodyWithTxs struct {
	// Transactions The hashes of the transactions included in this block
	Transactions []Txn `json:"transactions"`
}

type BlockWithTxs struct {
	Status BlockStatus `json:"status"`
	BlockHeader
	BlockBodyWithTxs
}

type PendingBlockWithTxs struct {
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`

	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress string `json:"sequencer_address"`

	// ParentHash The hash of this block's parent
	ParentHash BlockHash `json:"parent_hash"`

	BlockBodyWithTxs
}

type BroadcastedTxn interface{}

type BroadcastedInvokeTxnDuck struct {
	BroadcastedCommonTxnProperties
	InvokeTxnDuck
}

type BroadcastedInvokeTxn interface {
	Version() uint64
}

type BroadcastedInvokeTxnV0 struct {
	BroadcastedCommonTxnProperties
	InvokeV0
}

type BroadcastedInvokeTxnV1 struct {
	BroadcastedCommonTxnProperties
	InvokeV1
}

type ContractEntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique identifier of the entry point (function) in the program
	Selector string `json:"selector"`
}

type ContractEntryPointList []ContractEntryPoint

type ContractABI []ContractABIEntry

type EntryPointsByType struct {
	CONSTRUCTOR ContractEntryPointList `json:"CONSTRUCTOR"`
	EXTERNAL    ContractEntryPointList `json:"EXTERNAL"`
	L1_HANDLER  ContractEntryPointList `json:"L1_HANDLER"`
}

type ContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`

	Abi *ContractABI `json:"abi,omitempty"`
}

type ContractABIEntry interface {
	IsType() string
}

type StructABIType string

const (
	StructABITypeEvent StructABIType = "struct"
)

type EventABIType string

const (
	EventABITypeEvent EventABIType = "event"
)

type FunctionABIType string

const (
	FunctionABITypeFunction  FunctionABIType = "function"
	FunctionABITypeL1Handler FunctionABIType = "l1_handler"
)

type StructABIEntry struct {
	// The event type
	Type StructABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Size uint64 `json:"size"`

	Members []StructMember `json:"members"`
}

type StructMember struct {
	TypedParameter
	Offset uint64 `json:"offset"`
}

type EventABIEntry struct {
	// The event type
	Type EventABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Keys []TypedParameter `json:"keys"`

	Data TypedParameter `json:"data"`
}

type FunctionABIEntry struct {
	// The function type
	Type FunctionABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}

type BroadcastedDeclareTxn struct {
	BroadcastedCommonTxnProperties
	ContractClass ContractClass `json:"contract_class"`
	SenderAddress Address       `json:"sender_address"`
}

type DeployTxnProperties struct {
	// Version of the transaction scheme
	Version NumAsHex `json:"version"`

	Type TxnType `json:"type"`

	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt string `json:"contract_address_salt"`

	// ConstructorCallData The parameters passed to the constructor
	ConstructorCallData []string `json:"constructor_calldata"`
}

type BroadcastedDeployTxn struct {
	ContractClass ContractClass `json:"contract_class"`
	DeployTxnProperties
}

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
	StorageEntry
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
	ContractAddress Address `json:"contract_address"`
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
	BlockHash BlockHash `json:"block_hash"`
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

type TxnStatus string

type CommonReceiptProperties struct {
	TransactionHash TxnHash `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee   string      `json:"actual_fee"`
	Status      TxnStatus   `json:"status"`
	BlockHash   BlockHash   `json:"block_hash"`
	BlockNumber BlockNumber `json:"block_number"`
	Type        *TxnType    `json:"type,omitempty"`
}

type MsgToL1 struct {
	// ToAddress The target L1 address the message is sent to
	ToAddress string `json:"to_address"`
	//Payload  The payload of the message
	Payload []string `json:"payload"`
}

type EventContent struct {
	Keys []string `json:"keys"`
	Data []string `json:"data"`
}

type Event struct {
	FromAddress Address `json:"from_address"`
	//payload  The payload of the message
	EventContent
}

type InvokeTxnReceiptProperties struct {
	MessageSent []MsgToL1 `json:"messages_sent"`
	// A list of events assocuated with the Invoke Transaction
	Events []Event `json:"events"`
}

// InvokeTxnReceipt Invoke Transaction Receipt
type InvokeTxnReceipt struct {
	CommonReceiptProperties
	// ActualFee The fee that was charged by the sequencer
	*InvokeTxnReceiptProperties `json:",omitempty"`
}

// DeclareTxnReceipt Declare Transaction Receipt
type DeclareTxnReceipt struct {
	CommonReceiptProperties
}

// DeployTxnReceipt Deploy Transaction Receipt
type DeployTxnReceipt struct {
	CommonReceiptProperties
	// ContractAddress The address of the deployed contract
	ContractAddress string `json:"contract_address"`
}

// L1HandlerTxnReceipt L1 Handler Transaction Receipt
type L1HandlerTxnReceipt struct {
	CommonReceiptProperties
}

type TxnReceipt interface{}

type PendingCommonReceiptProperties struct {
	TransactionHash TxnHash `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee string  `json:"actual_fee"`
	Type      TxnType `json:"type"`
}

type PendingInvokeTxnReceipt struct {
	PendingCommonReceiptProperties
	InvokeTxnReceiptProperties
}

type PendingTxnReceipt interface{}

type EmittedEvent struct {
	Event
	BlockHash       BlockHash   `json:"block_hash"`
	BlockNumber     BlockNumber `json:"block_number"`
	TransactionHash TxnHash     `json:"transaction_hash"`
}

type EventFilter struct {
	FromBlock BlockID `json:"from_block"`
	ToBlock   BlockID `json:"to_block"`
	Address   Address `json:"address"`
	// Keys the values used to filter the events
	Keys []string `json:"keys"`
}

type ResultPageRequest struct {
	// ContinuationToken a pointer to the last element of the delivered page, use this token in a subsequent query to obtain the next page
	ContinuationToken *string `json:"continuation_token"`

	ChunkSize uint64 `json:"chunk_size"`
}

type EventFilterParams struct {
	EventFilter
	ResultPageRequest
}

type EventsOutput struct {
	Events            []EmittedEvent `json:"events"`
	ContinuationToken string         `json:"continuation_token"`
}

type FeeEstimate struct {
	GasConsumed NumAsHex `json:"gas_consumed"`
	GasPrice    NumAsHex `json:"gas_price"`
	OverallFee  NumAsHex `json:"overall_fee"`
}
