package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

// ChainID retrieves the current chain ID for transaction replay protection.
func (sc *Provider) ChainID(ctx context.Context) (string, error) {
	var result string
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := sc.c.CallContext(ctx, &result, "starknet_chainId", []interface{}{}...); err != nil {
		return "", err
	}
	return caigo.HexToShortStr(result), nil
}

// Syncing checks the syncing status of the node.
func (sc *Provider) Syncing(ctx context.Context) (*types.SyncResponse, error) {
	var result types.SyncResponse
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := sc.c.CallContext(ctx, &result, "starknet_syncing", []interface{}{}...); err != nil {
		return nil, err
	}
	return &result, nil
}

// StateUpdate gets the information about the result of executing the requested block.
func (sc *Provider) StateUpdate(ctx context.Context, blockID types.BlockID) (*types.StateUpdateOutput, error) {
	var state types.StateUpdateOutput
	if err := do(ctx, sc.c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, err
	}
	return &state, nil
}
