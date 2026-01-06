package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// Syncing retrieves the synchronisation status of the provider.
//
// Parameters:
//   - ctx: The context.Context object for the function
//
// Returns:
//   - *SyncStatus: The synchronisation status
//   - error: An error if any occurred during the execution
func Syncing(ctx context.Context, c callers.Caller) (rpcv10.SyncStatus, error) {
	var result rpcv10.SyncStatus
	if err := internal.Do(ctx, c, "starknet_syncing", &result); err != nil {
		return rpcv10.SyncStatus{}, rpcerr.UnwrapToRPCErr(err)
	}

	return result, nil
}
