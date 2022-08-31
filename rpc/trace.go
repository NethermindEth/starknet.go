package rpc

import (
	"context"
)

// not implemented for testing yet
func (sc *Client) TransactionTrace(ctx context.Context, hash string) error {
	return errNotImplemented
}

// not implemented for testing yet
func (sc *Client) TraceBlockTransactions(ctx context.Context, hash string) error {
	return errNotImplemented
}
