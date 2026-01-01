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
func (provider *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := do(ctx, provider.c, "starknet_blockNumber", &blockNumber); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrNoBlocks
		}

		return 0, rpcerr.UnwrapToRPCErr(err)
	}

	return blockNumber, nil
}
