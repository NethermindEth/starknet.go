package rpc

import (
	"context"
	"fmt"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

// Call a starknet function without creating a StarkNet transaction.
func (sc *Client) Call(ctx context.Context, call types.FunctionCall, block types.BlockID) ([]string, error) {
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.CallData) == 0 {
		call.CallData = make([]string, 0)
	}
	var result []string
	if err := sc.do(ctx, "starknet_call", &result, call, block); err != nil {
		return nil, err
	}
	return result, nil
}

// BlockNumber gets the most recent accepted block number.
func (sc *Client) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := sc.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil
}

// BlockHashAndNumber gets block information given the block number or its hash.
func (sc *Client) BlockHashAndNumber(ctx context.Context) (*types.BlockHashAndNumberOutput, error) {
	var block types.BlockHashAndNumberOutput
	if err := sc.do(ctx, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func WithBlockNumber(n uint64) types.BlockID {
	return types.BlockID{
		Number: n,
	}
}

func WithBlockHash(h types.Hash) types.BlockID {
	return types.BlockID{
		Hash: &h,
	}
}

func WithBlockTag(tag string) types.BlockID {
	return types.BlockID{
		Tag: tag,
	}
}

// PendingTransactions returns the list of pending transactions.
func (sc *Client) PendingTransactions(ctx context.Context) (types.Transactions, error) {
	var txns types.Transactions
	if err := sc.do(ctx, "starknet_pendingTransactions", &txns); err != nil {
		return nil, err
	}
	return txns, nil
}

// BlockWithTxHashes gets block information given the block id.
func (sc *Client) BlockWithTxHashes(ctx context.Context, blockID types.BlockID) (types.Block, error) {
	var result types.Block
	if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return types.Block{}, err
	}
	return result, nil
}

// BlockTransactionCount gets the number of transactions in a block
func (sc *Client) BlockTransactionCount(ctx context.Context, blockID types.BlockID) (uint64, error) {
	var result uint64
	if err := sc.do(ctx, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, err
	}
	return result, nil
}

// Nonce returns the Nonce of a contract
func (sc *Client) Nonce(ctx context.Context, contractAddress types.Hash) (*string, error) {
	nonce := ""
	if err := sc.do(ctx, "starknet_getNonce", &nonce, contractAddress); err != nil {
		return nil, err
	}
	return &nonce, nil
}

// BlockWithTxs get block information with full transactions given the block id.
func (sc *Client) BlockWithTxs(ctx context.Context, blockID types.BlockID) (interface{}, error) {
	var result types.Block
	if err := sc.do(ctx, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, err
	}
	return &result, nil
}

// Class gets the contract class definition associated with the given hash.
func (sc *Client) Class(ctx context.Context, classHash string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := sc.do(ctx, "starknet_getClass", &rawClass, classHash); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassAt get the contract class definition at the given address.
func (sc *Client) ClassAt(ctx context.Context, blockID types.BlockID, contractAddress types.Hash) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := sc.do(ctx, "starknet_getClassAt", &rawClass, blockID, contractAddress); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (sc *Client) ClassHashAt(ctx context.Context, blockID types.BlockID, contractAddress types.Hash) (*string, error) {
	var result string
	if err := sc.do(ctx, "starknet_getClassHashAt", &result, blockID, contractAddress); err != nil {
		return nil, err
	}
	return &result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (sc *Client) StorageAt(ctx context.Context, contractAddress types.Hash, key string, blockID types.BlockID) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%s", caigo.GetSelectorFromName(key).Text(16))
	if err := sc.do(ctx, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID); err != nil {
		return "", err
	}
	return value, nil
}

// StateUpdate gets the information about the result of executing the requested block.
func (sc *Client) StateUpdate(ctx context.Context, blockID types.BlockID) (*types.StateUpdateOutput, error) {
	var state types.StateUpdateOutput
	if err := sc.do(ctx, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, err
	}
	return &state, nil
}

// TransactionByHash gets the details and status of a submitted transaction.
func (sc *Client) TransactionByHash(ctx context.Context, hash types.Hash) (types.Transaction, error) {
	var tx types.UnknownTransaction
	if err := sc.do(ctx, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (sc *Client) TransactionByBlockIdAndIndex(ctx context.Context, blockID types.BlockID, index uint64) (types.Transaction, error) {
	var tx types.UnknownTransaction
	if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TxnReceipt gets the transaction receipt by the transaction hash.
func (sc *Client) TransactionReceipt(ctx context.Context, transactionHash types.Hash) (types.TransactionReceipt, error) {
	var receipt types.UnknownTransactionReceipt
	err := sc.do(ctx, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, err
	}
	return receipt.TransactionReceipt, nil
}

// Events returns all events matching the given filter
func (sc *Client) Events(ctx context.Context, filter types.EventFilter) (*types.EventsOutput, error) {
	var result types.EventsOutput
	if err := sc.do(ctx, "starknet_getEvents", &result, filter); err != nil {
		return nil, err
	}

	return &result, nil
}

// EstimateFee estimates the fee for a given StarkNet transaction.
func (sc *Client) EstimateFee(ctx context.Context, request types.BroadcastedTxn, blockID types.BlockID) (*types.FeeEstimate, error) {
	var raw types.FeeEstimate
	if err := sc.do(ctx, "starknet_estimateFee", &raw, request, blockID); err != nil {
		return nil, err
	}
	return &raw, nil
}
