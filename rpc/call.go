package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// Call calls the Provider's StarkNet function with the given context, request, and block ID without creating a Starknet transaction.
//
// ctx: The context.Context to be used for the function call.
// request: The FunctionCall representing the call data.
// blockID: The BlockID representing the block to execute the function on.
// []*felt.Felt: The result of the function call.
// error: An error if the function call fails.
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
