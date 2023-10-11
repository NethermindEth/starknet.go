package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// Call calls the Starknet Provider's function with the given (Starknet) request and block ID.
//
// Parameters:
// - ctx: the context.Context object for the function call
// - request: the FunctionCall object representing the request
// - blockID: the BlockID object representing the block ID
// Returns
// - []*felt.Felt: the result of the function call
// - error: an error if any occurred during the execution
func (provider *Provider) Call(ctx context.Context, request FunctionCall, blockID BlockID) ([]*felt.Felt, error) {

	if len(request.Calldata) == 0 {
		request.Calldata = make([]*felt.Felt, 0)
	}
	var result []*felt.Felt
	if err := do(ctx, provider.c, "starknet_call", &result, request, blockID); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrContractError):
			return nil, ErrContractError
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return result, nil
}
