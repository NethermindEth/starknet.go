package main

import (
	"context"
	_ "embed"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
)

//go:embed artifacts/counter.json
var counterCompiled []byte

func (ap *accountPlugin) installCounterWithRPCv01(ctx context.Context, provider *rpcv01.Provider) (string, error) {
	p := RPCProvider(*provider)
	return (&p).deployContract(ctx, counterCompiled, ap.PublicKey, []string{})
}

func (ap *accountPlugin) installCounterWithGateway(ctx context.Context, provider *gateway.Gateway) (string, error) {
	p := GatewayProvider(*provider)
	return (&p).deployContract(ctx, counterCompiled, ap.PublicKey, []string{})
}
