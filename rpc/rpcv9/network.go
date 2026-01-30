package rpcv9

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
)

// ChainID returns the chain ID for transaction replay protection.
//
// Parameters:
//   - ctx: The context.Context object for the function
//
// Returns:
//   - string: The chain ID
//   - error: An error if any occurred during the execution
func ChainID(ctx context.Context, c callers.Caller) (string, error) {
	var result string
	if err := internal.Do(ctx, c, "starknet_chainId", &result); err != nil {
		return "", rpcerr.UnwrapToRPCErr(err)
	}

	return internalUtils.HexToShortStr(result), nil
}

// SpecVersion returns the version of the Starknet JSON-RPC specification being
// implemented by the node.
//
// Parameters:
//   - ctx: The context for the function.
//
// Returns:
//   - string: The version of the Starknet JSON-RPC specification
//     implemented by the node.
//   - error: An error if the request fails.
func SpecVersion(ctx context.Context, c callers.Caller) (string, error) {
	var result string
	err := internal.Do(ctx, c, "starknet_specVersion", &result)
	if err != nil {
		return "", rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	return result, nil
}

// Syncing retrieves the synchronisation status of the provider.
//
// Parameters:
//   - ctx: The context.Context object for the function
//
// Returns:
//   - *SyncStatus: The synchronisation status
//   - error: An error if any occurred during the execution
func Syncing(ctx context.Context, c callers.Caller) (SyncStatus, error) {
	var result SyncStatus
	if err := internal.Do(ctx, c, "starknet_syncing", &result); err != nil {
		return SyncStatus{}, rpcerr.UnwrapToRPCErr(err)
	}

	return result, nil
}
