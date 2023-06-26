package rpcv02

import (
	"context"

	"github.com/smartcontractkit/caigo/types"
)

// ChainID retrieves the current chain ID for transaction replay protection.
func (provider *Provider) ChainID(ctx context.Context) (string, error) {
	if provider.chainID != "" {
		return provider.chainID, nil
	}
	var result string
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := provider.c.CallContext(ctx, &result, "starknet_chainId", []interface{}{}...); err != nil {
		return "", err
	}
	provider.chainID = types.HexToShortStr(result)
	return provider.chainID, nil
}

// Syncing checks the syncing status of the node.
func (provider *Provider) Syncing(ctx context.Context) (*SyncStatus, error) {
	// TODO: manage the fact SyncingStatus is set to false.
	var result SyncStatus
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := provider.c.CallContext(ctx, &result, "starknet_syncing", []interface{}{}...); err != nil {
		return nil, err
	}
	return &result, nil
}
