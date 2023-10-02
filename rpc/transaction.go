package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
)

var (
	feltZero = new(felt.Felt).SetUint64(0)
	feltOne  = new(felt.Felt).SetUint64(1)
	feltTwo  = new(felt.Felt).SetUint64(2)
)

// adaptTransaction converts a TXN object into a Transaction object.
//
// It takes a TXN object as a parameter and returns a Transaction object
// and an error. The TXN object represents a transaction, and the
// Transaction object is the converted version of the TXN object. If
// the conversion is successful, the converted Transaction object is
// returned along with a nil error. If there is an error during the
// conversion, the function returns a nil Transaction object and an
// error describing the reason for the failure.
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
		switch {
		case t.Version.Equal(feltZero):
			var tx DeclareTxnV0
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		case t.Version.Equal(feltOne):
			var tx DeclareTxnV1
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		case t.Version.Equal(feltTwo):
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
	}
	return nil, errors.New(fmt.Sprint("internal error with adaptTransaction() : unknown transaction type ", t.Type))

}

// TransactionByHash retrieves the details and status of a transaction by its hash.
//
// ctx - The context.Context object for the request.
// hash - The hash of the transaction to retrieve.
// Returns the Transaction object representing the retrieved transaction, or an error if the transaction is not found.
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

// TransactionByBlockIdAndIndex returns the details of the transaction with the given block ID and index.
//
// ctx - The context.Context object for the request.
// blockID - The ID of the block.
// index - The index of the transaction within the block.
// Returns a Transaction object and an error if no transaction is found.

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

// PendingTransaction returns a list of pending transactions in the transaction pool, recognized by this sequencer.
//
// The function takes a context.Context as a parameter and returns a slice of Transaction and an error.
func (provider *Provider) PendingTransaction(ctx context.Context) ([]Transaction, error) {
	txs := []Transaction{}
	if err := do(ctx, provider.c, "starknet_pendingTransactions", &txs, []interface{}{}); err != nil {
		return nil, err
	}
	return txs, nil
}

// TransactionReceipt retrieves the transaction receipt for a given transaction hash.
//
// ctx - The context.Context object for the request.
// transactionHash - The transaction hash for which to retrieve the receipt.
// Returns the TransactionReceipt object and an error if any.
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

// WaitForTransaction waits for a transaction to be executed and returns its execution status.
//
// It takes a context.Context object as the first parameter to control the execution flow and cancellation.
// The second parameter, transactionHash, is the hash of the transaction that needs to be monitored.
// The third parameter, pollInterval, is the time interval between each poll to check the transaction status.
// The function returns a TxnExecutionStatus and an error. TxnExecutionStatus represents the execution status of the transaction,
// and the error indicates any error occurred during the execution.
func (provider *Provider) WaitForTransaction(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (TxnExecutionStatus, error) {
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
			return receipt.GetExecutionStatus(), nil
		}
	}
}
