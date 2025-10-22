package rpc

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

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
func (provider *Provider) SpecVersion(ctx context.Context) (string, error) {
	var result string
	err := do(ctx, provider.c, "starknet_specVersion", &result)
	if err != nil {
		return "", rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	return result, nil
}

// IsCompatible compares the version of the Starknet JSON-RPC Specification
// implemented by the node with the version implemented by the Provider type,
// and returns whether they are the same or not.
//
// Parameters:
//   - ctx: The context for the function.
//
// Returns:
//   - bool: True if the node version is compatible with the SDK version, false otherwise.
//   - error: An error if any.
func (provider *Provider) IsCompatible(ctx context.Context) (bool, error) {
	rawNodeVersion, err := provider.SpecVersion(ctx)
	if err != nil {
		// Print a warning but don't fail
		fmt.Println(warnVersionCheckFailed, err)

		return false, err
	}

	nodeVersion, err := semver.NewVersion(rawNodeVersion)
	if err != nil {
		return false, fmt.Errorf("failed to parse node version: %w", err)
	}

	return rpcVersion.Compare(nodeVersion) == 0, nil
}
