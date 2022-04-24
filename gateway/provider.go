package gateway

import (
	"context"

	"github.com/dontpanicdao/caigo/types"
)

type GatewayProvider struct {
	Gateway
}

func NewProvider(opts ...Option) *GatewayProvider {
	return &GatewayProvider{
		*NewClient(opts...),
	}
}

func (p *GatewayProvider) BlockByHash(ctx context.Context, hash string) (*types.Block, error) {
	b, err := p.Block(ctx, &BlockOptions{BlockHash: hash})
	if err != nil {
		return nil, err
	}

	return b.Normalize(), nil
}

func (p *GatewayProvider) BlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	b, err := p.Block(ctx, &BlockOptions{BlockNumber: number})
	if err != nil {
		return nil, err
	}

	return b.Normalize(), nil
}

func (p *GatewayProvider) TransactionByHash(ctx context.Context, hash string) (*types.Transaction, error) {
	t, err := p.Transaction(ctx, TransactionOptions{TransactionHash: hash})
	if err != nil {
		return nil, err
	}

	return t.Transaction.Normalize(), nil
}
