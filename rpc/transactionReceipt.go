package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

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
	c callCloser,
	transactionHash *felt.Felt,
) (*TransactionReceiptWithBlockInfo, error) {
	var receipt TransactionReceiptWithBlockInfo
	err := do(ctx, c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound)
	}

	return &receipt, nil
}
