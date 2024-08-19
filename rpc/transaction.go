package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

// TransactionByHash retrieves the details and status of a transaction by its hash.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - hash: The hash of the transaction.
// Returns:
// - BlockTransaction: The retrieved Transaction
// - error: An error if any
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
// - ctx: The context.Context object for the request.
// - blockID: The ID of the block containing the transaction.
// - index: The index of the transaction within the block.
// Returns:
// - BlockTransaction: The retrieved Transaction object
// - error: An error, if any
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
// - ctx: the context.Context object for the request
// - transactionHash: the hash of the transaction as a Felt
// Returns:
// - TransactionReceipt: the transaction receipt
// - error: an error if any
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
// - ctx: the context.Context object for cancellation and timeouts.
// - transactionHash: the transaction hash as a felt
// Returns:
// - *GetTxnStatusResp: The transaction status
// - error, if one arose.
func (provider *Provider) GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResp, error) {
	var receipt TxnStatusResp
	err := do(ctx, provider.c, "starknet_getTransactionStatus", &receipt, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound)
	}
	return &receipt, nil
}
