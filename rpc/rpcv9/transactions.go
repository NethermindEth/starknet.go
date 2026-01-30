package rpcv9

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
)

// TransactionByBlockIDAndIndex retrieves a transaction by its block ID and index.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - blockID: The ID of the block containing the transaction.
//   - index: The index of the transaction within the block.
//
// Returns:
//   - BlockTransaction: The retrieved Transaction object
//   - error: An error, if any
func TransactionByBlockIDAndIndex(
	ctx context.Context,
	c callers.Caller,
	blockID BlockID,
	index uint64,
) (*BlockTransaction, error) {
	var tx BlockTransaction
	if err := internal.Do(
		ctx, c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrInvalidTxnIndex, ErrBlockNotFound)
	}

	return &tx, nil
}

// TransactionByHash retrieves the details and status of a transaction by its hash.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - hash: The hash of the transaction.
//
// Returns:
//   - BlockTransaction: The retrieved Transaction
//   - error: An error if any
func TransactionByHash(
	ctx context.Context,
	c callers.Caller,
	hash *felt.Felt,
) (*BlockTransaction, error) {
	var tx BlockTransaction
	if err := internal.Do(ctx, c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound)
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
func GetTransactionReceipt(
	ctx context.Context,
	c callers.Caller,
	transactionHash *felt.Felt,
) (*TransactionReceiptWithBlockInfo, error) {
	var receipt TransactionReceiptWithBlockInfo
	err := internal.Do(ctx, c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound)
	}

	return &receipt, nil
}

// TransactionStatus gets the transaction status (possibly reflecting that
// the tx is still in the mempool, or dropped from it)
// Parameters:
//   - ctx: the context.Context object for cancellation and timeouts.
//   - transactionHash: The hash of the requested transaction
//
// Returns:
//   - *TxnStatusResult: Transaction status result, including finality status
//     and execution status
//   - error, if one arose.
func TransactionStatus(
	ctx context.Context,
	c callers.Caller,
	transactionHash *felt.Felt,
) (*TxnStatusResult, error) {
	var receipt TxnStatusResult
	err := internal.Do(ctx, c, "starknet_getTransactionStatus", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound)
	}

	return &receipt, nil
}
