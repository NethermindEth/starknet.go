package rpcv02

import (
	"context"
	"fmt"
	"time"

	types "github.com/NethermindEth/caigo/types"
)

// TransactionByHash gets the details and status of a submitted transaction.
func (provider *Provider) TransactionByHash(ctx context.Context, hash types.Felt) (Transaction, error) {
	var tx UnknownTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		// TODO: Bind Pathfinder/Devnet Error to
		// TXN_HASH_NOT_FOUND
		return nil, err
	}
	return tx.Transaction, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (provider *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error) {
	var tx UnknownTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		// TODO: Bind Pathfinder/Devnet Error to
		// INVALID_TXN_INDEX and INVALID_TXN_INDEX
		return nil, err
	}
	return tx.Transaction, nil
}

// PendingTransaction returns the transactions in the transaction pool, recognized by this sequencer.
func (provider *Provider) PendingTransaction(ctx context.Context) ([]Transaction, error) {
	txs := []Transaction{}
	if err := do(ctx, provider.c, "starknet_pendingTransactions", &txs, []interface{}{}); err != nil {
		return nil, err
	}
	return txs, nil
}

// TxnReceipt gets the transaction receipt by the transaction hash.
func (provider *Provider) TransactionReceipt(ctx context.Context, transactionHash types.Felt) (TransactionReceipt, error) {
	var receipt UnknownTransactionReceipt
	err := do(ctx, provider.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		// TODO: check Pathfinder/Devnet for error
		// TXN_HASH_NOT_FOUND
		return nil, err
	}
	return receipt.TransactionReceipt, nil
}

// WaitForTransaction waits for the transaction to succeed or fail
func (provider *Provider) WaitForTransaction(ctx context.Context, transactionHash types.Felt, pollInterval time.Duration) (types.TransactionState, error) {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-t.C:
			_, err := provider.TransactionByHash(ctx, transactionHash)
			if err != nil {
				break
			}
			receipt, err := provider.TransactionReceipt(ctx, transactionHash)
			if err != nil {
				continue
			}
			switch r := receipt.(type) {
			case DeclareTransactionReceipt:
				if r.Status.IsTransactionFinal() {
					return r.Status, nil
				}
			case DeployTransactionReceipt:
				if r.Status.IsTransactionFinal() {
					return r.Status, nil
				}
			case InvokeTransactionReceipt:
				if r.Status.IsTransactionFinal() {
					return r.Status, nil
				}
			case L1HandlerTransactionReceipt:
				if r.Status.IsTransactionFinal() {
					return r.Status, nil
				}
			default:
				return "", fmt.Errorf("unknown receipt %T", receipt)
			}
		}
	}
}
