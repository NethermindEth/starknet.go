package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/rpcerr"
)

// SpecVersion returns the version of the Starknet JSON-RPC specification being used
// Parameters: None
// Returns: String of the Starknet JSON-RPC specification
func (provider *Provider) SpecVersion(ctx context.Context) (string, error) {
	var result string
	err := do(ctx, provider.c, "starknet_specVersion", &result)
	if err != nil {
		return "", rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	return result, nil
}
