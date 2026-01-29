package rpc

import (
	"context"
	"time"

	"github.com/NethermindEth/juno/core/felt"
)

// @new
// WaitForTransactionReceipt waits for the transaction receipt of the given
// transaction hash to succeed or fail.
//
// Parameters:
//   - provider: The RPC provider to use.
//   - ctx: The context
//   - transactionHash: The hash of the transaction to wait for
//   - pollInterval: The time interval to poll the transaction receipt
//
// Returns:
//   - *TxReceipt: the transaction receipt
//   - error: an error if any
func WaitForTransactionReceipt[
	TxReceipt any,
	P interface {
		TransactionReceipt(
			ctx context.Context,
			transactionHash *felt.Felt,
		) (*TxReceipt, error)
	},
](
	provider P,
	ctx context.Context,
	transactionHash *felt.Felt,
	pollInterval time.Duration,
) (*TxReceipt, error) {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
			receiptWithBlockInfo, err := provider.TransactionReceipt(ctx, transactionHash)
			if err != nil {
				return nil, err
			}

			return receiptWithBlockInfo, nil
		}
	}
}
