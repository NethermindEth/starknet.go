package rpcv10

// @todo remove this file later

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

var (
	// rpcVersion is the version of the Starknet JSON-RPC specification that
	// this SDK is compatible with.
	// This should be updated when supporting new versions of the RPC specification.
	rpcVersion = semver.MustParse("0.10.0")

	// ErrIncompatibleVersion is returned when the JSON-RPC specification  implemented
	// by the node is different from the version implemented by the Provider type.
	ErrIncompatibleVersion = errors.New("incompatible JSON-RPC specification version")
)

// Provider provides the provider for starknet.go/rpc implementation.
type Provider struct {
	c       callers.Caller
	chainID string
}

// WsProvider provides the provider for websocket starknet.go/rpc implementation.
type WsProvider struct {
	s callers.Subscriber
}

// Close closes the client, aborting any in-flight requests.
func (ws *WsProvider) Close() {
	ws.s.Close()
}

// NewProvider creates a new HTTP rpc Provider instance.
//
// Parameters:
//   - ctx: The context for the function.
//   - url: The URL of the RPC endpoint.
//   - options: The options for the client.
//
// Returns:
//   - *Provider: The new Provider instance.
//   - error: An error if any.
//     If the node JSON-RPC specification version is different from the version
//     implemented by the Provider type, the ErrIncompatibleVersion will be returned,
//     but the returned Provider instance is valid.
func NewProvider(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*Provider, error) {
	c, err := internal.NewHTTPClient(ctx, url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	provider := &Provider{c: c, chainID: ""}

	// Check version compatibility
	isCompatible, nodeVersion, err := IsCompatible(ctx, provider)
	if err != nil {
		return nil, err
	}
	if !isCompatible {
		return provider, errors.Join(
			ErrIncompatibleVersion,
			fmt.Errorf("expected version: %s, got: %s", rpcVersion, nodeVersion),
		)
	}

	return provider, nil
}

// NewWebsocketProvider creates a new Websocket rpc Provider instance.
func NewWebsocketProvider(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*WsProvider, error) {
	s, err := internal.NewWSClient(ctx, url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Websocket client: %w", err)
	}

	return &WsProvider{s: s}, nil
}

//go:generate mockgen -destination=../../internal/tests/mocks/rpcv10mock/rpc.go -package=rpcv10mock -source=provider.go
type RPCProvider interface {
	AddInvokeTransaction(
		ctx context.Context,
		invokeTxn *types.BroadcastInvokeTxnV3,
	) (AddInvokeTransactionResponse, error)
	AddDeclareTransaction(
		ctx context.Context,
		declareTransaction *types.BroadcastDeclareTxnV3,
	) (AddDeclareTransactionResponse, error)
	AddDeployAccountTransaction(
		ctx context.Context,
		deployAccountTransaction *types.BroadcastDeployAccountTxnV3,
	) (AddDeployAccountTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID types.BlockID) (uint64, error)
	BlockWithReceipts(ctx context.Context, blockID types.BlockID) (interface{}, error)
	BlockWithTxHashes(ctx context.Context, blockID types.BlockID) (interface{}, error)
	BlockWithTxs(ctx context.Context, blockID types.BlockID) (interface{}, error)
	Call(ctx context.Context, call types.FunctionCall, block types.BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, blockID types.BlockID, classHash *felt.Felt) (ClassOutput, error)
	ClassAt(
		ctx context.Context,
		blockID types.BlockID,
		contractAddress *felt.Felt,
	) (ClassOutput, error)
	ClassHashAt(
		ctx context.Context,
		blockID types.BlockID,
		contractAddress *felt.Felt,
	) (*felt.Felt, error)
	CompiledCasm(ctx context.Context, classHash *felt.Felt) (*contracts.CasmClass, error)
	EstimateFee(
		ctx context.Context,
		requests []types.BroadcastTxn,
		simulationFlags []types.SimulationFlag,
		blockID types.BlockID,
	) ([]types.FeeEstimation, error)
	EstimateMessageFee(
		ctx context.Context,
		msg MsgFromL1,
		blockID types.BlockID,
	) (types.MessageFeeEstimation, error)
	Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	MessagesStatus(ctx context.Context, transactionHash NumAsHex) ([]MessageStatus, error)
	Nonce(
		ctx context.Context,
		blockID types.BlockID,
		contractAddress *felt.Felt,
	) (*felt.Felt, error)
	SimulateTransactions(
		ctx context.Context,
		blockID types.BlockID,
		txns []types.BroadcastTxn,
		simulationFlags []types.SimulationFlag,
	) ([]SimulatedTransaction, error)
	SpecVersion(ctx context.Context) (string, error)
	StateUpdate(ctx context.Context, blockID types.BlockID) (*StateUpdateOutput, error)
	StorageAt(
		ctx context.Context,
		contractAddress *felt.Felt,
		key string,
		blockID types.BlockID,
	) (string, error)
	StorageProof(
		ctx context.Context,
		storageProofInput StorageProofInput,
	) (*StorageProofResult, error)
	Syncing(ctx context.Context) (SyncStatus, error)
	TraceBlockTransactions(ctx context.Context, blockID types.BlockID) ([]Trace, error)
	TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error)
	TransactionByBlockIDAndIndex(
		ctx context.Context,
		blockID types.BlockID,
		index uint64,
	) (*types.BlockTransaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (*types.BlockTransaction, error)
	TransactionReceipt(
		ctx context.Context,
		transactionHash *felt.Felt,
	) (*types.TransactionReceiptWithBlockInfo, error)
	TransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResult, error)
}

type WebsocketProvider interface {
	SubscribeEvents(
		ctx context.Context,
		events chan<- *EmittedEventWithFinalityStatus,
		options *EventSubscriptionInput,
	) (*client.ClientSubscription, error)
	SubscribeNewHeads(
		ctx context.Context,
		headers chan<- *BlockHeader,
		subBlockID types.SubscriptionBlockID,
	) (*client.ClientSubscription, error)
	SubscribeNewTransactions(
		ctx context.Context,
		newTxns chan<- *TxnWithHashAndStatus,
		options *SubNewTxnsInput,
	) (*client.ClientSubscription, error)
	SubscribeNewTransactionReceipts(
		ctx context.Context,
		txnReceipts chan<- *types.TransactionReceiptWithBlockInfo,
		options *SubNewTxnReceiptsInput,
	) (*client.ClientSubscription, error)
	SubscribeTransactionStatus(
		ctx context.Context,
		newStatus chan<- *NewTxnStatus,
		transactionHash *felt.Felt,
	) (*client.ClientSubscription, error)
}

var (
	_ RPCProvider       = (*Provider)(nil)
	_ WebsocketProvider = (*WsProvider)(nil)
)
