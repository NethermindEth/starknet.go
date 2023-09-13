package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// Call a starknet function without creating a Starknet transaction.
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
