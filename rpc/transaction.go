package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// PendingTransactions returns the list of pending transactions.
func (sc *Provider) PendingTransactions(ctx context.Context) (types.Transactions, error) {
	var txns types.Transactions
	if err := do(ctx, sc.c, "starknet_pendingTransactions", &txns); err != nil {
		return nil, err
	}
	return txns, nil
}

// TransactionByHash gets the details and status of a submitted transaction.
func (sc *Provider) TransactionByHash(ctx context.Context, hash types.Hash) (types.Transaction, error) {
	var tx types.UnknownTransaction
	if err := do(ctx, sc.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (sc *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID types.BlockID, index uint64) (types.Transaction, error) {
	var tx types.UnknownTransaction
	if err := do(ctx, sc.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TxnReceipt gets the transaction receipt by the transaction hash.
func (sc *Provider) TransactionReceipt(ctx context.Context, transactionHash types.Hash) (types.TransactionReceipt, error) {
	var receipt types.UnknownTransactionReceipt
	err := do(ctx, sc.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, err
	}
	return receipt.TransactionReceipt, nil
}
