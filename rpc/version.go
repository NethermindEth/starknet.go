package rpc

import "context"

// RPCVersion is the version of the Starknet JSON-RPC specification that this SDK is compatible with.
// This should be updated when supporting new versions of the RPC specification.
const RPCVersion = "0.8.1"

// SpecVersion returns the version of the Starknet JSON-RPC specification being used
// Parameters: None
// Returns: String of the Starknet JSON-RPC specification
func (provider *Provider) SpecVersion(ctx context.Context) (string, error) {
	var result string
	err := do(ctx, provider.c, "starknet_specVersion", &result)
	if err != nil {
		return "", Err(InternalError, StringErrData(err.Error()))
	}

	return result, nil
}
