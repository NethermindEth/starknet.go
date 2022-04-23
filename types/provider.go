package types

import (
	"context"
	"math/big"
)

type Provider interface {
	AccountNonce(context.Context, string) (*big.Int, error)
	BlockByHash(context.Context, string) (*Block, error)
	BlockByNumber(context.Context, uint64) (*Block, error)
	ChainID(context.Context) (string, error)
	Invoke(context.Context, Transaction) (*AddTxResponse, error)
}
