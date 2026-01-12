package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

type Provider[P RPCProvider] struct {
	provider P
}

func (p *Provider[P]) addInvokeTransaction(
	ctx context.Context,
	invokeTxn *types.BroadcastInvokeTxnV3,
) (types.TransactionResponse, error) {
	switch p.provider.(type) {
	case *rpcv10.Provider:
		resp, err := p.provider.AddInvokeTransaction(ctx, invokeTxn)
		if err != nil {
			return types.TransactionResponse{}, err
		}
		return types.TransactionResponse{Hash: resp.Hash}, nil
	}
}

func (p *Provider[P]) addDeclareTransaction(
	ctx context.Context,
	declareTxn *types.BroadcastDeclareTxnV3,
) (types.TransactionResponse, error) {
	return p.v10.AddDeclareTransaction(ctx, declareTxn)
}

func (p *Provider[P]) addDeployAccountTransaction(
	ctx context.Context,
	deployAccountTxn *types.BroadcastDeployAccountTxnV3,
) (types.TransactionResponse, error) {
	return p.v10.AddDeployAccountTransaction(ctx, deployAccountTxn)
}
