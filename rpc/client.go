package rpc

import (
	"context"
	"encoding/json"
)

type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// do is a function that performs a remote call and populates the provided data structure.
//
// It takes a context.Context object, a callCloser object, a method string, a data interface{},
// and an optional variadic argument args. It returns an error.
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
