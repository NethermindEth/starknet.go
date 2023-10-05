package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/starknet.go/utils"
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
	provider.chainID = utils.HexToShortStr(result)
	return provider.chainID, nil
}

// Syncing checks the syncing status of the node.
func (provider *Provider) Syncing(ctx context.Context) (*SyncStatus, error) {
	var result interface{}
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := provider.c.CallContext(ctx, &result, "starknet_syncing", []interface{}{}...); err != nil {
		return nil, err
	}
	switch res := result.(type) {
	case bool:
		return &SyncStatus{SyncStatus: res}, nil
	case SyncStatus:
		return &res, nil
	default:
		return nil, errors.New("internal error with starknet_syncing")
	}

}
