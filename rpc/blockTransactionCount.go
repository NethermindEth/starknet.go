package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// BlockTransactionCount returns the number of transactions in a specific block.
//
// Parameters:
//   - ctx: The context.Context object to handle cancellation signals and timeouts
//   - blockID: The ID of the block to retrieve the number of transactions from
//
// Returns:
//   - uint64: The number of transactions in the block
//   - error: An error, if any
func BlockTransactionCount(
	ctx context.Context,
	c callCloser,
	blockID BlockID,
) (uint64, error) {
	var result uint64
	if err := do(ctx, c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return result, nil
}
