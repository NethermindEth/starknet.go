package rpc

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10/methods"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// @todo update docs for the entire package
// add tests where needed

type BasicProvider struct {
	c       callers.Caller
	chainID string
	version RPCVersion

	rpcv9  *rpcv9.Provider
	rpcv10 *rpcv10.Provider
}

// NewBasicProvider creates a new HTTP rpc Provider instance.
//
// Parameters:
//   - ctx: The context for the function.
//   - url: The URL of the RPC endpoint.
//   - options: The options for the client.
//
// Returns:
//   - *Provider: The new Provider instance.
//   - error: An error if any.
//     If the node JSON-RPC specification version is different from the version
//     implemented by the Provider type, the ErrIncompatibleVersion will be returned,
//     but the returned Provider instance is valid.
func NewBasicProvider(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*BasicProvider, error) {
	c, err := internal.NewHTTPClient(ctx, url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	rawNodeVersion, err := methods.SpecVersion(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to get the node's RPC spec version: %w", err)
	}

	var RPCVersion RPCVersion
	err = RPCVersion.UnmarshalJSON([]byte(rawNodeVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node version: %w", err)
	}
	var provider BasicProvider
	provider.version = RPCVersion

	switch RPCVersion {
	case RPCVersion9:
		rpcv9, err := rpcv9.NewProvider(ctx, url, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create RPCv9 provider: %w", err)
		}
		provider.rpcv9 = rpcv9
	case RPCVersion10:
		rpcv10, err := rpcv10.NewProvider(ctx, url, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create RPCv10 provider: %w", err)
		}
		provider.rpcv10 = rpcv10
	}

	return &provider, nil
}
