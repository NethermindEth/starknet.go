package rpcv02

import (
	"context"

	"github.com/smartcontractkit/caigo/types"
)

// Call a starknet function without creating a StarkNet transaction.
func (provider *Provider) Call(ctx context.Context, request types.FunctionCall, blockID BlockID) ([]string, error) {
	request.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(request.EntryPointSelector))
	if len(request.Calldata) == 0 {
		request.Calldata = make([]string, 0)
	}
	var result []string
	if err := do(ctx, provider.c, "starknet_call", &result, request, blockID); err != nil {
		// TODO: Bind Pathfinder/Devnet Error to
		// CONTRACT_NOT_FOUND, INVALID_MESSAGE_SELECTOR, INVALID_CALL_DATA, CONTRACT_ERROR, BLOCK_NOT_FOUND
		return nil, err
	}
	return result, nil
}
