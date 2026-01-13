package rpc

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10/methods"
)

// @todo update docs for the entire package
// add tests where needed

type Provider struct {
	c       callers.Caller
	chainID string
	version RPCVersion
}

// NewProvider creates a new HTTP rpc Provider instance.
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
func NewProvider(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*Provider, error) {
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

	return &Provider{c: c, chainID: "", version: RPCVersion}, nil
}
