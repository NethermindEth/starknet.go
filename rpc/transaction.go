package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
)

func adaptTransaction(t TXN) (Transaction, error) {
	txMarshalled, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case TransactionType_Invoke:
		var tx InvokeTxnV1
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	case TransactionType_Declare:

		switch t.Version {
		case &felt.Zero:
			var tx DeclareTxnV0
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		case felt.Zero.SetUint64(1):
			var tx DeclareTxnV1
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		case felt.Zero.SetUint64(2):
			var tx DeclareTxnV2
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		}
	case TransactionType_Deploy:
		var tx DeployTxn
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	case TransactionType_DeployAccount:
		var tx DeployAccountTxn
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	case TransactionType_L1Handler:
		var tx L1HandlerTxn
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	default:
		return nil, errors.New("Internal error with adaptTransaction() : unknown transaction type")
	}
	return nil, errors.New("Internal error with adaptTransaction()")

}

// TransactionByHash gets the details and status of a submitted transaction.
func (provider *Provider) TransactionByHash(ctx context.Context, hash *felt.Felt) (Transaction, error) {
	// todo: update to return a custom Transaction type, then use adapt function
	var tx TXN
	if err := do(ctx, provider.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		if errors.Is(err, ErrHashNotFound) {
			return nil, ErrHashNotFound
		}
		return nil, err
	}
	return adaptTransaction(tx)
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (provider *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error) {
	var tx TXN
	if err := do(ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		switch {
		case errors.Is(err, ErrInvalidTxnIndex):
			return nil, ErrInvalidTxnIndex
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return adaptTransaction(tx)
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
func (provider *Provider) TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (TransactionReceipt, error) {
	var receipt UnknownTransactionReceipt
	err := do(ctx, provider.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		if errors.Is(err, ErrHashNotFound) {
			return nil, ErrHashNotFound
		}
		return nil, err
	}
	return receipt.TransactionReceipt, nil
}

// WaitForTransaction waits for the transaction to succeed or fail
func (provider *Provider) WaitForTransaction(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (TransactionState, error) {
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
			case DeployAccountTransactionReceipt:
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
