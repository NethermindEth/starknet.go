package rpc

import "context"

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

// Version represents the current version of the Starknet RPC client.
type Version struct {
	// The major version number.
	Major int
	// The minor version number.
	Minor int
	// The patch version number.
	Patch int
}

// SDKVersion is the current version of the Starknet Go SDK.
// This should be updated when releasing new versions of the SDK.
const SDKVersion = "0.8.0"
