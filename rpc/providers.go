package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// @new
func NewProviderV10(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*rpcv10.Provider, error) {
	provider, err := rpcv10.NewProvider(ctx, url, options...)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// @new
func NewWSProviderV10(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*rpcv10.WsProvider, error) {
	provider, err := rpcv10.NewWebsocketProvider(ctx, url, options...)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// @new
func NewProviderV9(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*rpcv9.Provider, error) {
	provider, err := rpcv9.NewProvider(ctx, url, options...)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// @new
func NewWSProviderV9(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*rpcv9.WsProvider, error) {
	provider, err := rpcv9.NewWebsocketProvider(ctx, url, options...)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
