package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

type SyncResponse struct {
	StartingBlockHash string `json:"starting_block_hash"`
	StartingBlockNum  string `json:"starting_block_num"`
	CurrentBlockHash  string `json:"current_block_hash"`
	CurrentBlockNum   string `json:"current_block_num"`
	HighestBlockHash  string `json:"highest_block_hash"`
	HighestBlockNum   string `json:"highest_block_num"`
}

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

func (sc *Client) Syncing(ctx context.Context) (*SyncResponse, error) {
	var result SyncResponse
	if err := sc.c.CallContext(ctx, &result, "starknet_syncing"); err != nil {
		return nil, err
	}

	return &result, nil
}

func (sc *Client) do(ctx context.Context, method string, data interface{}, args ...interface{}) error {
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
