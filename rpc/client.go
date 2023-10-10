package rpc

import (
	"context"
	"encoding/json"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// do is a function that performs a remote procedure call (RPC) using the provided callCloser.
//
// It takes a context.Context object, a callCloser object, a method string, a data interface{}, and optional args ...interface{} as parameters.
// It returns an error.
func do(ctx context.Context, call callCloser, method string, data interface{}, args ...interface{}) error {
	var raw json.RawMessage
	err := call.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errNotFound
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	return nil
}

// NewClient creates a new ethrpc.Client instance.
//
// It takes a URL string as a parameter and returns a pointer to ethrpc.Client and an error.
func NewClient(url string) (*ethrpc.Client, error) {
	return ethrpc.DialContext(context.Background(), url)
}
