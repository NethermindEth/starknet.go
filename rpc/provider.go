package rpc

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/juno/core/felt"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/publicsuffix"
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

// NewProvider creates a new rpc Provider instance.
func NewProvider(url string) (*Provider, error) {

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	c, err := ethrpc.DialHTTPWithClient(url, client)
	if err != nil {
		return nil, err
	}

	return &Provider{c: c}, nil
}

//go:generate mockgen -destination=../mocks/mock_rpc_provider.go -package=mocks -source=provider.go api
type RpcProvider interface {
	AddInvokeTransaction(ctx context.Context, invokeTxn BroadcastInvokeTxnType) (*AddInvokeTransactionResponse, *RPCError)
	AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastDeclareTxnType) (*AddDeclareTransactionResponse, *RPCError)
	AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastAddDeployTxnType) (*AddDeployAccountTransactionResponse, *RPCError)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, *RPCError)
	BlockNumber(ctx context.Context) (uint64, *RPCError)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, *RPCError)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, *RPCError)
	BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, *RPCError)
	Call(ctx context.Context, call FunctionCall, block BlockID) ([]*felt.Felt, *RPCError)
	ChainID(ctx context.Context) (string, *RPCError)
	Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, *RPCError)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, *RPCError)
	ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, *RPCError)
	EstimateFee(ctx context.Context, requests []BroadcastTxn, simulationFlags []SimulationFlag, blockID BlockID) ([]FeeEstimate, *RPCError)
	EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, *RPCError)
	Events(ctx context.Context, input EventsInput) (*EventChunk, *RPCError)
	BlockWithReceipts(ctx context.Context, blockID BlockID) (interface{}, *RPCError)
	GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResp, *RPCError)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, *RPCError)
	SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, *RPCError)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, *RPCError)
	StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, *RPCError)
	SpecVersion(ctx context.Context) (string, *RPCError)
	Syncing(ctx context.Context) (*SyncStatus, *RPCError)
	TraceBlockTransactions(ctx context.Context, blockID BlockID) ([]Trace, *RPCError)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, *RPCError)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (Transaction, *RPCError)
	TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (*TransactionReceiptWithBlockInfo, *RPCError)
	TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, *RPCError)
}

var _ RpcProvider = &Provider{}
