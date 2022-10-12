package contracts

import (
	"context"
	_ "embed"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
)

//go:embed artifacts/counter.json
var CounterCompiled []byte

func (ap *AccountPlugin) InstallCounterWithRPCv01(ctx context.Context, provider *rpcv01.Provider) (string, error) {
	p := RPCProvider(*provider)
	return (&p).deployContract(ctx, CounterCompiled, ap.PublicKey, []string{})
}

func (ap *AccountPlugin) InstallCounterWithGateway(ctx context.Context, provider *gateway.Gateway) (string, error) {
	p := GatewayProvider(*provider)
	return (&p).deployContract(ctx, CounterCompiled, ap.PublicKey, []string{})
}
