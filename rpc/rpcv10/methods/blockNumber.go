package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// BlockNumber returns the block number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request
//
// Returns:
//   - uint64: The block number
//   - error: An error if any
func BlockNumber(ctx context.Context, c rpc.Caller) (uint64, error) {
	var blockNumber uint64
	if err := internal.Do(ctx, c, "starknet_blockNumber", &blockNumber); err != nil {
		return 0, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrNoBlocks)
	}

	return blockNumber, nil
}
