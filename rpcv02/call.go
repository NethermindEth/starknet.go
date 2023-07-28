package rpcv02

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// Call a starknet function without creating a StarkNet transaction.
func (provider *Provider) Call(ctx context.Context, request FunctionCall, blockID BlockID) ([]string, error) {

	if len(request.Calldata) == 0 {
		request.Calldata = make([]*felt.Felt, 0)
	}
	var result []string
	if err := do(ctx, provider.c, "starknet_call", &result, request, blockID); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrInvalidMessageSelector):
			return nil, ErrInvalidMessageSelector
		case errors.Is(err, ErrInvalidCallData):
			return nil, ErrInvalidCallData
		case errors.Is(err, ErrContractError):
			return nil, ErrContractError
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return result, nil
}
