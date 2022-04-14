package caigo

import (
	"net/http"
)

type gatewayOptions struct {
	client  *http.Client
	chainID string
}

// funcGatewayOption wraps a function that modifies gatewayOptions into an
// implementation of the GatewayOption interface.
type funcGatewayOption struct {
	f func(*gatewayOptions)
}

func (fso *funcGatewayOption) apply(do *gatewayOptions) {
	fso.f(do)
}

func newFuncGatewayOption(f func(*gatewayOptions)) *funcGatewayOption {
	return &funcGatewayOption{
		f: f,
	}
}

// GatewayOption configures how we set up the connection.
type GatewayOption interface {
	apply(*gatewayOptions)
}

func WithHttpClient(client http.Client) GatewayOption {
	return newFuncGatewayOption(func(o *gatewayOptions) {
		o.client = &client
	})
}

func WithChain(chainID string) GatewayOption {
	return newFuncGatewayOption(func(o *gatewayOptions) {
		o.chainID = chainID
	})
}
