package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
func (provider *Provider) TransactionByBlockIDAndIndex(
	ctx context.Context,
	blockID BlockID,
	index uint64,
) (*BlockTransaction, error) {
	var tx BlockTransaction
	if err := do(
		ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrInvalidTxnIndex, ErrBlockNotFound)
	}

	return &tx, nil
}
