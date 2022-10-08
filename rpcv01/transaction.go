package rpcv01

import (
	"context"
	"fmt"
	"time"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// TransactionByHash gets the details and status of a submitted transaction.
func (provider *Provider) TransactionByHash(ctx context.Context, hash ctypes.Hash) (Transaction, error) {
	var tx UnknownTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (provider *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error) {
	var tx UnknownTransaction
	if err := do(ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

// TxnReceipt gets the transaction receipt by the transaction hash.
func (provider *Provider) TransactionReceipt(ctx context.Context, transactionHash ctypes.Hash) (TransactionReceipt, error) {
	var receipt UnknownTransactionReceipt
	err := do(ctx, provider.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, err
	}
	return receipt.TransactionReceipt, nil
}

// WaitForTransaction waits for the transaction to succeed or fail
func (provider *Provider) WaitForTransaction(ctx context.Context, transactionHash ctypes.Hash, pollInterval time.Duration) (TransactionStatus, error) {
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
			case DeclareTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			case DeployTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			case InvokeTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			case L1HandlerTransactionReceipt:
				if isTransactionFinal(status.Status) {
					return status.Status, nil
				}
			default:
				return "", fmt.Errorf("unknown receipt %T", receipt)
			}
		}
	}
}

func isTransactionFinal(v TransactionStatus) bool {
	if v == TransactionStatus("ACCEPTED_ON_L1") ||
		v == TransactionStatus("ACCEPTED_ON_L2") ||
		v == TransactionStatus("PENDING") ||
		v == TransactionStatus("REJECTED") ||
		v == TransactionStatus("NOT_RECEIVED") {
		return true
	}
	return false
}
