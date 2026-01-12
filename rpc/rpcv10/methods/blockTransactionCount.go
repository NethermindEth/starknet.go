package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
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
	c callers.Caller,
	blockID types.BlockID,
) (uint64, error) {
	var result uint64
	if err := internal.Do(ctx, c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}

	return result, nil
}
