package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ErrNotFound is returned by API methods if the requested item does not exist.
var ErrNotFound = errors.New("not found")

type Client struct {
	c *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) *Client {
	return &Client{c: c}
}

func (sc *Client) Close() {
	sc.c.Close()
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (sc *Client) ChainID(ctx context.Context) (string, error) {
	var result string
	err := sc.c.CallContext(ctx, &result, "starknet_chainId")
	if err != nil {
		return "", err
	}
	return result, err
}

func (sc *Client) AccountNonce(context.Context, string) (*big.Int, error) {
	panic("not implemented")
}

func (sc *Client) BlockByHash(ctx context.Context, hash string, scope string) (*types.Block, error) {
	return sc.getBlock(ctx, "starknet_getBlockByHash", hash, scope)
}

func (sc *Client) BlockByNumber(ctx context.Context, number *big.Int, scope string) (*types.Block, error) {
	return sc.getBlock(ctx, "starknet_getBlockByNumber", toBlockNumArg(number), scope)
}

func (sc *Client) Invoke(context.Context, types.Transaction) (*types.AddTxResponse, error) {
	panic("not implemented")
}

func (sc *Client) TransactionByHash(context.Context, string) (*types.Transaction, error) {
	panic("not implemented")
}

func (sc *Client) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	err := sc.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, ErrNotFound
	}

	var block types.Block
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

func toBlockNumArg(number *big.Int) interface{} {
	var numOrTag interface{}

	if number == nil {
		numOrTag = "latest"
	} else if number.Cmp(big.NewInt(-1)) == 0 {
		numOrTag = "pending"
	} else {
		numOrTag = number.Uint64()
	}

	return numOrTag
}
