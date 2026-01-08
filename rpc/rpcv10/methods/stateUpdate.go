package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// StateUpdate is a function that performs a state update operation
// (gets the information about the result of executing the requested block).
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - blockID: The ID of the block to retrieve the transactions from
//
// Returns:
//   - *StateUpdateOutput: The retrieved state update
//   - error: An error, if any
func GetStateUpdate(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
) (*rpcv10.StateUpdateOutput, error) {
	var state rpcv10.StateUpdateOutput
	if err := internal.Do(ctx, c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}

	return &state, nil
}
