package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

// Call a starknet function without creating a StarkNet transaction.
func (sc *Provider) Call(ctx context.Context, call types.FunctionCall, block types.BlockID) ([]string, error) {
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.CallData) == 0 {
		call.CallData = make([]string, 0)
	}
	var result []string
	if err := do(ctx, sc.c, "starknet_call", &result, call, block); err != nil {
		return nil, err
	}
	return result, nil
}
