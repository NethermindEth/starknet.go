package rpc

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

// not implemented for testing yet
func (sc *Client) TransactionTrace(ctx context.Context, hash string) error {
	return ErrNotImplemented
}

// not implemented for testing yet
func (sc *Client) TraceBlockTransactions(ctx context.Context, hash string) error {
	return ErrNotImplemented
}
