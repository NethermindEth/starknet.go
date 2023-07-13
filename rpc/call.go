package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/types"
)

// Call a starknet function without creating a StarkNet transaction.
func (provider *Provider) Call(ctx context.Context, request FunctionCall, blockID BlockID) ([]string, error) {
	request.EntryPointSelector = types.GetSelectorFromNameFelt(request.EntryPointSelector.String())
	if len(request.Calldata) == 0 {
		request.Calldata = make([]*felt.Felt, 0)
	}
	var result []string
	if err := do(ctx, provider.c, "starknet_call", &result, request, blockID); err != nil {
		// TODO: Bind Pathfinder/Devnet Error to
		// CONTRACT_NOT_FOUND, INVALID_MESSAGE_SELECTOR, INVALID_CALL_DATA, CONTRACT_ERROR, BLOCK_NOT_FOUND
		return nil, err
	}
	return result, nil
}
