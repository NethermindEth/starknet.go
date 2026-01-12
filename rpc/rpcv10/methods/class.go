package methods

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// Class retrieves the class information from the Provider with the given hash.
//
// Parameters:
//   - ctx: The context.Context object
//   - blockID: The BlockID object
//   - classHash: The *felt.Felt object
//
// Returns:
//   - ClassOutput: The output of the class.
//   - error: An error if any occurred during the execution.
func Class(
	ctx context.Context,
	c callers.Caller,
	blockID types.BlockID,
	classHash *felt.Felt,
) (rpcv10.ClassOutput, error) {
	var rawClass map[string]any
	if err := internal.Do(ctx, c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrClassHashNotFound, rpcv10.ErrBlockNotFound)
	}

	return typecastClassOutput(rawClass)
}
