package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// BlockHashAndNumber retrieves the hash and number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request.
//
// Returns:
//   - *BlockHashAndNumberOutput: The hash and number of the current block
//   - error: An error if any
func BlockHashAndNumber(
	ctx context.Context,
	c callers.Caller,
) (*rpcv10.BlockHashAndNumberOutput, error) {
	var block rpcv10.BlockHashAndNumberOutput
	if err := internal.Do(ctx, c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrNoBlocks)
	}

	return &block, nil
}
