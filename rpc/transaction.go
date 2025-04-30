package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

// TransactionByHash retrieves the details and status of a transaction by its hash.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - hash: The hash of the transaction.
//
// Returns:
//   - BlockTransaction: The retrieved Transaction
//   - error: An error if any
func (provider *Provider) TransactionByHash(ctx context.Context, hash *felt.Felt) (*BlockTransaction, error) {
	var tx BlockTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound)
	}
	return &tx, nil
}

// TransactionByBlockIdAndIndex retrieves a transaction by its block ID and index.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - blockID: The ID of the block containing the transaction.
//   - index: The index of the transaction within the block.
//
// Returns:
//   - BlockTransaction: The retrieved Transaction object
//   - error: An error, if any
func (provider *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (*BlockTransaction, error) {
	var tx BlockTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrInvalidTxnIndex, ErrBlockNotFound)
	}
	return &tx, nil
}

// TransactionReceipt fetches the transaction receipt for a given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the hash of the transaction as a Felt
//
// Returns:
//   - TransactionReceipt: the transaction receipt
//   - error: an error if any
func (provider *Provider) TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (*TransactionReceiptWithBlockInfo, error) {
	var receipt TransactionReceiptWithBlockInfo
	err := do(ctx, provider.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound)
	}
	return &receipt, nil
}

// GetTransactionStatus gets the transaction status (possibly reflecting that the tx is still in the mempool, or dropped from it)
// Parameters:
//   - ctx: the context.Context object for cancellation and timeouts.
//   - transactionHash: The hash of the requested transaction
//
// Returns:
//   - *TxnStatusResult: Transaction status result, including finality status and execution status
//   - error, if one arose.
func (provider *Provider) GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResult, error) {
	var receipt TxnStatusResult
	err := do(ctx, provider.c, "starknet_getTransactionStatus", &receipt, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound)
	}
	return &receipt, nil
}

// Given an L1 tx hash, returns the associated l1_handler tx hashes and statuses for all L1 -> L2 messages sent by the l1 transaction, ordered by the L1 tx sending order
//
// Parameters:
//   - ctx: the context.Context object for cancellation and timeouts.
//   - transactionHash: The hash of the L1 transaction that sent L1->L2 messages
//
// Returns:
//   - [] MessageStatusResp: An array containing the status of the messages sent by the L1 transaction
//   - error, if one arose.
func (provider *Provider) GetMessagesStatus(ctx context.Context, transactionHash NumAsHex) ([]MessageStatus, error) {
	var response []MessageStatus
	err := do(ctx, provider.c, "starknet_getMessagesStatus", &response, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound)
	}
	return response, nil
}
