package rpc

import (
	"context"
	"errors"

	"github.com/dontpanicdao/caigo/rpc/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")
)

// Provider provides the provider for caigo/rpc implementation.
type Provider struct {
	c callCloser
}

// NewProvider creates a *Provider from an existing `go-ethereum/rpc` *Client.
func NewProvider(c *rpc.Client) *Provider {
	return &Provider{c: c}
}

type api interface {
	BlockHashAndNumber(ctx context.Context) (*types.BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID types.BlockID) (uint64, error)
	BlockWithTxHashes(ctx context.Context, blockID types.BlockID) (types.Block, error)
	BlockWithTxs(ctx context.Context, blockID types.BlockID) (interface{}, error)
	Call(ctx context.Context, call types.FunctionCall, block types.BlockID) ([]string, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, classHash string) (*types.ContractClass, error)
	ClassAt(ctx context.Context, blockID types.BlockID, contractAddress types.Hash) (*types.ContractClass, error)
	ClassHashAt(ctx context.Context, blockID types.BlockID, contractAddress types.Hash) (*string, error)
	EstimateFee(ctx context.Context, request types.Call, blockID types.BlockID) (*types.FeeEstimate, error)
	Events(ctx context.Context, filter types.EventFilter) (*types.EventsOutput, error)
	Nonce(ctx context.Context, contractAddress types.Hash) (*string, error)
	PendingTransactions(ctx context.Context) (types.Transactions, error)
	StateUpdate(ctx context.Context, blockID types.BlockID) (*types.StateUpdateOutput, error)
	StorageAt(ctx context.Context, contractAddress types.Hash, key string, blockID types.BlockID) (string, error)
	Syncing(ctx context.Context) (*types.SyncResponse, error)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID types.BlockID, index uint64) (types.Transaction, error)
	TransactionByHash(ctx context.Context, hash types.Hash) (types.Transaction, error)
	TransactionReceipt(ctx context.Context, transactionHash types.Hash) (types.TransactionReceipt, error)
}

var _ api = &Provider{}
