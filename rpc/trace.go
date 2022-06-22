package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo/types"
)

// not implemented for testing yet
func (sc *Client) TransactionTrace(ctx context.Context, hash string) (*types.TransactionTrace, error) {
	var tx types.TransactionTrace
	if err := sc.do(ctx, "starknet_traceTransaction", &tx, hash); err != nil {
		return nil, err
	}

	return &tx, nil
}

// not implemented for testing yet
// func (sc *Client) TraceBlockTransactions(ctx context.Context, hash string) ([]*types.TransactionTrace, error) {

// }
