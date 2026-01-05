package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
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
	c rpc.Caller,
) (*rpcv10.BlockHashAndNumberOutput, error) {
	var block rpcv10.BlockHashAndNumberOutput
	if err := internal.Do(ctx, c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrNoBlocks)
	}

	return &block, nil
}
