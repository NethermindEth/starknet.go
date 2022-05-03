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

func (sc *Client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var blockNumber big.Int
	err := sc.c.CallContext(ctx, &blockNumber, "starknet_blockNumber")
	if err != nil {
		return nil, err
	}
	return &blockNumber, err
}

func (sc *Client) LatestBlock(ctx context.Context, scope string) (*types.Block, error) {
	blockNumber, err := sc.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	block, err := sc.BlockByNumber(ctx, blockNumber, scope)
	if err != nil {
		return nil, err
	}

	return block, err
}

func (sc *Client) BlockByHash(ctx context.Context, hash string, scope string) (*types.Block, error) {
	var block types.Block
	err := sc.fetchStructure(ctx, "starknet_getBlockByHash", block, hash, scope)

	return &block, err
}

func (sc *Client) BlockByNumber(ctx context.Context, number *big.Int, scope string) (*types.Block, error) {
	var block types.Block
	err := sc.fetchStructure(ctx, "starknet_getBlockByNumber", block, toBlockNumArg(number), scope)

	return &block, err
}

func (sc *Client) CodeByAddress(ctx context.Context, address string) (*types.Code, error) {
	var contract types.Code
	err := sc.fetchStructure(ctx, "starknet_getCode", contract, address)

	return &contract, err
}

func (sc *Client) Invoke(context.Context, types.Transaction) (*types.AddTxResponse, error) {
	panic("not implemented")
}

func (sc *Client) TransactionByHash(ctx context.Context, hash string) (*types.Transaction, error) {
	var tx types.Transaction
	err := sc.fetchStructure(ctx, "starknet_getTransactionByHash", tx, hash)

	return &tx, err
}

func (sc *Client) fetchStructure(ctx context.Context, method string, data interface{}, args ...interface{}) error {
	var raw json.RawMessage
	err := sc.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return err
	} else if len(raw) == 0 {
		return ErrNotFound
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return nil
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
