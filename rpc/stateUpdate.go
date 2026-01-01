package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
	c callCloser,
	blockID BlockID,
) (*StateUpdateOutput, error) {
	var state StateUpdateOutput
	if err := do(ctx, c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return &state, nil
}
