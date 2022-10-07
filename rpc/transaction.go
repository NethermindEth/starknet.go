package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/dontpanicdao/caigo/rpc/types"
	ctypes "github.com/dontpanicdao/caigo/types"
)

// TransactionByHash gets the details and status of a submitted transaction.
func (provider *Provider) TransactionByHash(ctx context.Context, hash ctypes.Hash) (types.Transaction, error) {
	var tx types.UnknownTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (provider *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID types.BlockID, index uint64) (types.Transaction, error) {
	var tx types.UnknownTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TxnReceipt gets the transaction receipt by the transaction hash.
func (provider *Provider) TransactionReceipt(ctx context.Context, transactionHash ctypes.Hash) (types.TransactionReceipt, error) {
	var receipt types.UnknownTransactionReceipt
	err := do(ctx, provider.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, err
	}
	return receipt.TransactionReceipt, nil
}

// WaitForTransaction waits for the transaction to succeed or fail
func (provider *Provider) WaitForTransaction(ctx context.Context, transactionHash ctypes.Hash, pollInterval time.Duration) (types.TransactionStatus, error) {
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
			switch status := receipt.(type) {
			case types.DeclareTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			case types.DeployTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			case types.InvokeTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			case types.L1HandlerTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			default:
				return "", fmt.Errorf("unknown receipt %T", receipt)
			}
		}
	}
}

func isTransactionFinal(v types.TransactionStatus) bool {
	if v == types.TransactionStatus("ACCEPTED_ON_L1") ||
		v == types.TransactionStatus("ACCEPTED_ON_L2") ||
		v == types.TransactionStatus("PENDING") ||
		v == types.TransactionStatus("REJECTED") ||
		v == types.TransactionStatus("NOT_RECEIVED") {
		return true
	}
	return false
}
