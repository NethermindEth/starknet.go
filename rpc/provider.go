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
	AddInvokeTransaction(ctx context.Context, invokeTxn AddInvokeTxnInput) (*AddInvokeTransactionResponse, error)
	AddDeclareTransaction(ctx context.Context, declareTransaction AddDeclareTxnInput) (*AddDeclareTransactionResponse, error)
	AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction AddDeployAccountTxnInput) (*AddDeployTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error)
	BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error)
	Call(ctx context.Context, call FunctionCall, block BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error)
	ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	EstimateFee(ctx context.Context, requests []EstimateFeeInput, blockID BlockID) ([]FeeEstimate, error)
	EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, error)
	Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error)
	SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error)
	StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error)
	Syncing(ctx context.Context) (*SyncStatus, error)
	TraceBlockTransactions(ctx context.Context, blockHash *felt.Felt) ([]Trace, error)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (Transaction, error)
	TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (TransactionReceipt, error)
	TransactionTrace(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error)
}

var _ api = &Provider{}
