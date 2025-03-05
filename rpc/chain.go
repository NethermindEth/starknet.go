package rpc

import (
	"context"

	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

// ChainID returns the chain ID for transaction replay protection.
//
// Parameters:
// - ctx: The context.Context object for the function
// Returns:
// - string: The chain ID
// - error: An error if any occurred during the execution
func (provider *Provider) ChainID(ctx context.Context) (string, error) {
	if provider.chainID != "" {
		return provider.chainID, nil
	}
	var result string
	if err := do(ctx, provider.c, "starknet_chainId", &result); err != nil {
		return "", tryUnwrapToRPCErr(err)
	}
	provider.chainID = internalUtils.HexToShortStr(result)
	return provider.chainID, nil
}

// Syncing retrieves the synchronization status of the provider.
//
// Parameters:
// - ctx: The context.Context object for the function
// Returns:
// - *SyncStatus: The synchronization status
// - error: An error if any occurred during the execution
func (provider *Provider) Syncing(ctx context.Context) (*SyncStatus, error) {
	var result interface{}
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := provider.c.CallContext(ctx, &result, "starknet_syncing", []interface{}{}...); err != nil {
		return nil, Err(InternalError, StringErrData(err.Error()))
	}
	switch res := result.(type) {
	case bool:
		return &SyncStatus{SyncStatus: &res}, nil
	case SyncStatus:
		return &res, nil
	default:
		return nil, Err(InternalError, StringErrData("internal error with starknet_syncing"))
	}

}
