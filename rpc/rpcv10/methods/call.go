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

// Call calls the Starknet Provider's function with the given (Starknet) request and block ID.
//
// Parameters:
//   - ctx: the context.Context object for the function call
//   - request: the FunctionCall object representing the request
//   - blockID: the BlockID object representing the block ID
//
// Returns
//   - []*felt.Felt: the result of the function call
//   - error: an error if any occurred during the execution
func Call(
	ctx context.Context,
	c callers.Caller,
	request types.FunctionCall,
	blockID types.BlockID,
) ([]*felt.Felt, error) {
	if request.Calldata == nil {
		request.Calldata = []*felt.Felt{}
	}

	var result []*felt.Felt
	if err := internal.Do(ctx, c, "starknet_call", &result, request, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrContractNotFound,
			rpcv10.ErrEntrypointNotFound,
			rpcv10.ErrContractError,
			rpcv10.ErrBlockNotFound,
		)
	}

	return result, nil
}
