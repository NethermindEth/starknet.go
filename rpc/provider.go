package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")
)

type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// Provider provides the provider for caigo/rpc implementation.
type Provider struct {
	c callCloser
}

// Dial connects a client to the given URL. It creates a `go-ethereum/rpc` *Client and relies on context.Background().
func Dial(rawurl string) (*Provider, error) {
	return DialContext(context.Background(), rawurl)
}

// DialContext connects a Provider to the given URL with an existing context. It creates a `go-ethereum/rpc` *Client.
func DialContext(ctx context.Context, rawurl string) (*Provider, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewProvider(c), nil
}

// NewProvider creates a *Provider from an existing `go-ethereum/rpc` *Client.
func NewProvider(c *rpc.Client) *Provider {
	return &Provider{c: c}
}

// Close closes the underlying client.
func (sc *Provider) Close() {
	sc.c.Close()
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (sc *Provider) ChainID(ctx context.Context) (string, error) {
	var result string
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := sc.c.CallContext(ctx, &result, "starknet_chainId", []interface{}{}...); err != nil {
		return "", err
	}
	return caigo.HexToShortStr(result), nil
}

// Syncing checks the syncing status of the node.
func (sc *Provider) Syncing(ctx context.Context) (*types.SyncResponse, error) {
	var result types.SyncResponse
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := sc.c.CallContext(ctx, &result, "starknet_syncing", []interface{}{}...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (sc *Provider) do(ctx context.Context, method string, data interface{}, args ...interface{}) error {
	var raw json.RawMessage
	err := sc.c.CallContext(ctx, &raw, method, args...)
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
