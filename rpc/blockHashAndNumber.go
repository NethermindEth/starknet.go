package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// BlockHashAndNumber retrieves the hash and number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request.
//
// Returns:
//   - *BlockHashAndNumberOutput: The hash and number of the current block
//   - error: An error if any
func (provider *Provider) BlockHashAndNumber(
	ctx context.Context,
) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := do(ctx, provider.c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrNoBlocks)
	}

	return &block, nil
}
