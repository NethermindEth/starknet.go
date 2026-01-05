package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
func TransactionByBlockIDAndIndex(
	ctx context.Context,
	c rpc.Caller,
	blockID rpcv10.BlockID,
	index uint64,
) (*rpcv10.BlockTransaction, error) {
	var tx rpcv10.BlockTransaction
	if err := internal.Do(
		ctx, c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrInvalidTxnIndex, rpcv10.ErrBlockNotFound)
	}

	return &tx, nil
}
