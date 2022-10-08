package main

import (
	"context"
	_ "embed"

	"github.com/dontpanicdao/caigo/rpcv01"
)

//go:embed artifacts/counter.json
var counterCompiled []byte

func (ap *accountPlugin) installCounter(ctx context.Context, provider rpcv01.Provider) (string, error) {
	return deployContract(ctx, provider, counterCompiled, ap.PublicKey, []string{})
}
