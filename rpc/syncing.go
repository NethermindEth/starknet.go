package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// Syncing retrieves the synchronisation status of the provider.
//
// Parameters:
//   - ctx: The context.Context object for the function
//
// Returns:
//   - *SyncStatus: The synchronisation status
//   - error: An error if any occurred during the execution
func Syncing(ctx context.Context, c callCloser) (SyncStatus, error) {
	var result SyncStatus
	if err := do(ctx, c, "starknet_syncing", &result); err != nil {
		return SyncStatus{}, rpcerr.UnwrapToRPCErr(err)
	}

	return result, nil
}
