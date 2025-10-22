package rpc

import (
	"context"

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
