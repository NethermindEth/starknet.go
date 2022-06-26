package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo/types"
)

func (sc *Client) Invoke(context.Context, types.Transaction) (*types.AddTxResponse, error) {
	panic("'starknet_addInvokeTransaction' not implemented")
}

func (sc *Client) Declare(context.Context, types.Transaction) (*types.AddTxResponse, error) {
	panic("'starknet_addDeclareTransaction' not implemented")
}

func (sc *Client) Deploy(context.Context, types.Transaction) (*types.AddTxResponse, error) {
	panic("'starknet_addDeployTransaction' not implemented")
}
