package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

func do(ctx context.Context, call callCloser, method string, data interface{}, args ...interface{}) error {
	var raw json.RawMessage
	err := call.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errNotFound
	}
	fmt.Println(string(raw))
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	return nil
}

func NewClient(url string) (*ethrpc.Client, error) {
	return ethrpc.DialContext(context.Background(), url)
}
