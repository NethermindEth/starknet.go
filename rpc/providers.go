package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// @new
func NewRPCv10Provider(
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
