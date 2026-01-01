package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

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
	c callCloser,
	hash *felt.Felt,
) (*BlockTransaction, error) {
	var tx BlockTransaction
	if err := do(ctx, c, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound)
	}

	return &tx, nil
}
