package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// BlockNumber returns the block number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request
//
// Returns:
//   - uint64: The block number
//   - error: An error if any
func BlockNumber(ctx context.Context, c callCloser) (uint64, error) {
	var blockNumber uint64
	if err := do(ctx, c, "starknet_blockNumber", &blockNumber); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrNoBlocks
		}

		return 0, rpcerr.UnwrapToRPCErr(err)
	}

	return blockNumber, nil
}
