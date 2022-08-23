package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

var (
	errBadRequest      = errors.New("bad request")
	errBadTxType       = errors.New("bad transaction type")
	errInvalidBlockTag = errors.New("invalid blocktag")
)

type Events struct {
	Events []Event `json:"events"`
}

type Event struct {
	*types.Event
	FromAddress     string `json:"from_address"`
	BlockHash       string `json:"block_hash"`
	BlockNumber     int    `json:"block_number"`
	TransactionHash string `json:"transaction_hash"`
}

type EventParams struct {
	FromBlock  uint64 `json:"fromBlock"`
	ToBlock    uint64 `json:"toBlock"`
	PageSize   uint64 `json:"page_size"`
	PageNumber uint64 `json:"page_number"`
}

// Call a starknet function without creating a StarkNet transaction.
func (sc *Client) Call(ctx context.Context, call types.FunctionCall, hash string) ([]string, error) {
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.Calldata) == 0 {
		call.Calldata = make([]string, 0)
	}

	var result []string
	if err := sc.do(ctx, "starknet_call", &result, call, hash); err != nil {
		return nil, err
	}

	return result, nil
}

// BlockNumber gets the most recent accepted block number.
func (sc *Client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var blockNumber big.Int
	if err := sc.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		return nil, err
	}

	return &blockNumber, nil
}

// BlockHashAndNumberOutput is a struct that is returned by BlockHashAndNumber.
type BlockHashAndNumberOutput struct {
	BlockNumber uint64 `json:"block_number,omitempty"`
	BlockHash   string `json:"block_hash,omitempty"`
}

// BlockHashAndNumber gets block information given the block number or its hash.
func (sc *Client) BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := sc.do(ctx, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// blockID is an unexposed struct that is used in a OneOf for
// starknet_getBlockWithTxHashes.
type blockID struct {
	BlockNumber *BlockNumber `json:"block_number,omitempty"`
	BlockHash   *BlockHash   `json:"block_hash,omitempty"`
	BlockTag    *string      `json:"block_tag,omitempty"`
}

// BlockIDOption is an options that can be used as a parameter for
// starknet_getBlockWithTxHashes.
type BlockIDOption func(b *blockID) error

func WithBlockIDNumber(blockNumber BlockNumber) BlockIDOption {
	return BlockIDOption(func(b *blockID) error {
		b.BlockNumber = &blockNumber
		return nil
	})
}

func WithBlockIDHash(blockHash BlockHash) BlockIDOption {
	return BlockIDOption(func(b *blockID) error {
		b.BlockHash = &blockHash
		return nil
	})
}

func WithBlockIDTag(blockTag string) BlockIDOption {
	return BlockIDOption(func(b *blockID) error {
		if blockTag != "latest" && blockTag != "pending" {
			return errInvalidBlockTag
		}
		b.BlockTag = &blockTag
		return nil
	})
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

// BlockWithTxHashes gets block information given the block id.
func (sc *Client) BlockWithTxHashes(ctx context.Context, blockIDOption BlockIDOption) (interface{}, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag == "pending" {
		var block PendingBlockWithTxHashes
		if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &block, "pending"); err != nil {
			return nil, err
		}
		return &block, nil
	}
	var block BlockWithTxHashes
	if opt.BlockTag != nil && *opt.BlockTag == "latest" {
		if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &block, "latest"); err != nil {
			return nil, err
		}
		return &block, nil
	}
	if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &block, opt); err != nil {
		return nil, err
	}
	return &block, nil
}

type TxnType string

type NumAsHex string

type Signature []string

// BroadcastedCommonTxnProperties common properties of a transaction that is sent to the sequencer (but is not yet in a block)
type BroadcastedCommonTxnProperties struct {
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

type CommonTxnProperties struct {
	TransactionHash TxnHash
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

// InvokeTxnV0 version 0 invoke transaction
type InvokeTxnV0 struct {
	CommonTxnProperties
	FunctionCall
}

// InvokeTxnV1 version 1 invoke transaction
type InvokeTxnV1 struct {
	CommonTxnProperties
	AccountAddress Address `json:"account_address"`
	// CallData The parameters passed to the function
	CallData []string `json:"calldata"`
}

// InvokeTxnDuck is a type used to understand the Invoke Version
type InvokeTxnDuck struct {
	AccountAddress     Address `json:"account_address"`
	ContractAddress    Address `json:"contract_address"`
	EntryPointSelector string  `json:"entry_point_selector"`
}

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
	TransactionHash TxnHash

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

func guessTxWithType(i interface{}) (interface{}, error) {
	switch local := i.(type) {
	case map[string]interface{}:
		typeValue, ok := local["type"]
		if !ok {
			return nil, errBadTxType
		}
		value, ok := typeValue.(string)
		if !ok {
			return nil, errBadTxType
		}
		switch value {
		case "DECLARE":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := DeclareTxn{}
			err = json.Unmarshal(data, &tx)
			return tx, err
		case "DEPLOY":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := DeployTxn{}
			err = json.Unmarshal(data, &tx)
			return tx, err
		case "L1_HANDLER":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := L1HandlerTxn{}
			err = json.Unmarshal(data, &tx)
			return tx, err
		case "INVOKE":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := InvokeTxnDuck{}
			err = json.Unmarshal(data, &tx)
			if err != nil {
				return nil, err
			}
			if tx.AccountAddress != "" {
				txv1 := InvokeTxnV1{}
				err = json.Unmarshal(data, &txv1)
				return txv1, err
			}
			if tx.ContractAddress != "" && tx.EntryPointSelector != "" {
				txv0 := InvokeTxnV0{}
				err = json.Unmarshal(data, &txv0)
				return txv0, err
			}
			return nil, errBadTxType
		}
		return nil, errBadTxType
	}
	return nil, errBadTxType
}

func guessTxsWithType(txs []Txn) ([]Txn, error) {
	for k, v := range txs {
		tv, err := guessTxWithType(v)
		if err != nil {
			return nil, errBadTxType
		}
		txs[k] = tv
	}
	return txs, nil
}

// BlockWithTxs get block information with full transactions given the block id.
func (sc *Client) BlockWithTxs(ctx context.Context, blockIDOption BlockIDOption) (interface{}, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag == "pending" {
		var block PendingBlockWithTxs
		if err := sc.do(ctx, "starknet_getBlockWithTxs", &block, "pending"); err != nil {
			return nil, err
		}
		txns, err := guessTxsWithType(block.Transactions)
		if err != nil {
			return nil, err
		}
		block.Transactions = txns
		return &block, nil
	}
	var block BlockWithTxs
	if opt.BlockTag != nil && *opt.BlockTag == "latest" {
		if err := sc.do(ctx, "starknet_getBlockWithTxs", &block, "latest"); err != nil {
			return nil, err
		}
		txns, err := guessTxsWithType(block.Transactions)
		if err != nil {
			return nil, err
		}
		block.Transactions = txns
		return &block, nil
	}
	if err := sc.do(ctx, "starknet_getBlockWithTxs", &block, opt); err != nil {
		return nil, err
	}
	txns, err := guessTxsWithType(block.Transactions)
	if err != nil {
		return nil, err
	}
	block.Transactions = txns
	return &block, nil
}

// CodeAt returns the contract and class associated with the an address.
// Deprecated: you should use ClassAt and TransactionByHash to access the
// associated values.
func (sc *Client) CodeAt(ctx context.Context, address string) (*types.Code, error) {
	var contractRaw struct {
		Bytecode []string `json:"bytecode"`
		AbiRaw   string   `json:"abi"`
		Abi      types.ABI
	}
	if err := sc.do(ctx, "starknet_getCode", &contractRaw, address); err != nil {
		return nil, err
	}

	contract := types.Code{
		Bytecode: contractRaw.Bytecode,
	}
	if err := json.Unmarshal([]byte(contractRaw.AbiRaw), &contract.Abi); err != nil {
		return nil, err
	}

	return &contract, nil
}

// Class gets the contract class definition associated with the given hash.
func (sc *Client) Class(ctx context.Context, hash string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := sc.do(ctx, "starknet_getClass", &rawClass, hash); err != nil {
		return nil, err
	}

	return &rawClass, nil
}

// ClassAt get the contract class definition at the given address.
func (sc *Client) ClassAt(ctx context.Context, address string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := sc.do(ctx, "starknet_getClassAt", &rawClass, address); err != nil {
		return nil, err
	}

	return &rawClass, nil
}

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (sc *Client) ClassHashAt(ctx context.Context, blockIDOption BlockIDOption, contractAddress Address) (*string, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag != "pending" && *opt.BlockTag != "latest" {
		return nil, errInvalidBlockTag
	}
	var result string
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getClassHashAt", &result, *opt.BlockTag, contractAddress); err != nil {
			return nil, err
		}
		return &result, nil
	}
	if err := sc.do(ctx, "starknet_getStateUpdate", &result, opt, contractAddress); err != nil {
		return nil, err
	}
	return &result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (sc *Client) StorageAt(ctx context.Context, contractAddress Address, key string, blockIDOption BlockIDOption) (string, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return "", err
	}
	var value string
	hashKey := fmt.Sprintf("0x%s", caigo.GetSelectorFromName(key).Text(16))
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getStorageAt", &value, string(contractAddress), hashKey, *opt.BlockTag); err != nil {
			return "", err
		}
		return value, nil
	}
	if err := sc.do(ctx, "starknet_getStorageAt", &value, string(contractAddress), hashKey, opt); err != nil {
		return "", err
	}
	return value, nil
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

	// StorageEntries the changes in the storage of the contract
	StorageEntries []StorageEntry `json:"storage_entries"`
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

// Nonce is a the updated nonce per contract address
type Nonce struct {
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
	Nonces []Nonce `json:"nonces"`
}

type StateUpdateOutput struct {
	// BlockHash is the block identifier,
	BlockHash BlockHash `json:"block_hash"`
	// NewRoot is the new global state root.
	NewRoot string `json:"new_root"`
	// OldRoot is the previous global state root.
	OldRoot string `json:"old_root"`
	// AcceptedTime is when the block was accepted on L1.
	AcceptedTime int `json:"accepted_time"`
	// StateDiff is the change in state applied in this block, given as a
	// mapping of addresses to the new values and/or new contracts.
	StateDiff StateDiff `json:"state_diff"`
}

// StateUpdate gets the information about the result of executing the requested block.
func (sc *Client) StateUpdate(ctx context.Context, blockIDOption BlockIDOption) (*StateUpdateOutput, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag != "latest" {
		return nil, errInvalidBlockTag
	}
	var state StateUpdateOutput
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getStateUpdate", &state, "latest"); err != nil {
			return nil, err
		}
		return &state, nil
	}
	if err := sc.do(ctx, "starknet_getStateUpdate", &state, opt); err != nil {
		return nil, err
	}
	return &state, nil
}

// TransactionByHash gets the details and status of a submitted transaction.
func (sc *Client) TransactionByHash(ctx context.Context, hash TxnHash) (*Txn, error) {
	var tx interface{}
	if err := sc.do(ctx, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	}
	txWithType, err := guessTxWithType(tx)
	if err != nil {
		return nil, err
	}
	txTxn := Txn(txWithType)
	return &txTxn, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (sc *Client) TransactionByBlockIdAndIndex(ctx context.Context, blockIDOption BlockIDOption, index uint64) (*Txn, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag != "latest" {
		return nil, errInvalidBlockTag
	}
	var tx interface{}
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, *opt.BlockTag, index); err != nil {
			return nil, err
		}
		txWithType, err := guessTxWithType(tx)
		if err != nil {
			return nil, err
		}
		txTxn := Txn(txWithType)
		return &txTxn, nil
	}
	if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, opt, index); err != nil {
		return nil, err
	}
	txWithType, err := guessTxWithType(tx)
	if err != nil {
		return nil, err
	}
	txTxn := Txn(txWithType)
	return &txTxn, nil
}

// TransactionReceipt gets the transaction receipt by the transaction hash.
func (sc *Client) TransactionReceipt(ctx context.Context, hash string) (*types.TransactionReceipt, error) {
	var receipt types.TransactionReceipt
	err := sc.do(ctx, "starknet_getTransactionReceipt", &receipt, hash)
	if err != nil {
		return nil, err
	} else if receipt.TransactionHash == "" {
		return nil, ErrNotFound
	}

	return &receipt, nil
}

// Events returns all events matching the given filter
// TODO: check the query parameters as they include filter directives that have
// not been implemented. For more details, check the
// [specification](https://github.com/starkware-libs/starknet-specs/blob/master/api/starknet_api_openrpc.json)
func (sc *Client) Events(ctx context.Context, evParams EventParams) (*Events, error) {
	var result Events
	if err := sc.do(ctx, "starknet_getEvents", &result, evParams); err != nil {
		return nil, err
	}

	return &result, nil
}

type rpcFeeEstimate struct {
	GasUsage   string `json:"gas_usage"`
	GasPrice   string `json:"gas_price"`
	OverallFee string `json:"overall_fee"`
}

// EstimateFee estimates the fee for a given StarkNet transaction.
func (sc *Client) EstimateFee(ctx context.Context, call types.FunctionInvoke, blockHashOrTag string) (*types.FeeEstimate, error) {
	var raw rpcFeeEstimate
	if err := sc.do(ctx, "starknet_estimateFee", &raw, call, blockHashOrTag); err != nil {
		return nil, err
	}

	usage, err := strconv.ParseUint(strings.TrimPrefix(raw.GasUsage, "0x"), 16, 64)
	if err != nil {
		return nil, err
	}
	price, err := strconv.ParseUint(strings.TrimPrefix(raw.GasPrice, "0x"), 16, 64)
	if err != nil {
		return nil, err
	}
	fee, err := strconv.ParseUint(strings.TrimPrefix(raw.OverallFee, "0x"), 16, 64)
	if err != nil {
		return nil, err
	}

	return &types.FeeEstimate{
		GasUsage:   usage,
		GasPrice:   price,
		OverallFee: fee,
	}, nil
}

// AccountNonce gets the latest nonce associated with the given address
func (sc *Client) AccountNonce(ctx context.Context, contractAddress string) (*big.Int, error) {
	var nonce big.Int
	err := sc.do(ctx, "starknet_getNonce", &nonce, contractAddress)
	return &nonce, err
}

func (sc *Client) Invoke(context.Context, types.FunctionInvoke) (*types.AddTxResponse, error) {
	panic("not implemented")
}
