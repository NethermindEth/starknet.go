package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/ethereum/go-ethereum/rpc"
)

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")
)

// Provider provides the provider for starknet.go/rpc implementation.
type Provider struct {
	c       callCloser
	chainID string
}

// NewProvider creates a *Provider from an existing `go-ethereum/rpc` *Client.
func NewProvider(c *rpc.Client) *Provider {
	return &Provider{c: c}
}

type api interface {
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error)
	BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error)
	Call(ctx context.Context, call FunctionCall, block BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, blockID BlockID, classHash string) (*ContractClass, error)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*ContractClass, error)
	ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error)
	EstimateFee(ctx context.Context, requests []BroadcastedTransaction, blockID BlockID) ([]FeeEstimate, error)
	EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, error)
	Events(ctx context.Context, input EventsInput) (*EventsOutput, error)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error)
	StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error)
	Syncing(ctx context.Context) (*SyncStatus, error)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (Transaction, error)
	TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (TransactionReceipt, error)
	AddInvokeTransaction(ctx context.Context, broadcastedInvoke BroadcastedInvokeTransaction) (*AddInvokeTransactionResponse, error)
	AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastedDeclareTransaction) (*AddDeclareTransactionResponse, error)
	AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastedDeployAccountTransaction) (*AddDeployTransactionResponse, error)
}

var _ api = &Provider{}
