package rpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/Masterminds/semver/v3"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/gorilla/websocket"
	"golang.org/x/net/publicsuffix"
)

// rpcVersion is the version of the Starknet JSON-RPC specification that
// this SDK is compatible with.
// This should be updated when supporting new versions of the RPC specification.
var rpcVersion = semver.MustParse("0.10.0")

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")

	// ErrIncompatibleVersion is returned when the JSON-RPC specification  implemented
	// by the node is different from the version implemented by the Provider type.
	ErrIncompatibleVersion = errors.New("incompatible JSON-RPC specification version")
)

// Provider provides the provider for starknet.go/rpc implementation.
type Provider struct {
	c       callCloser
	chainID string
}

// WsProvider provides the provider for websocket starknet.go/rpc implementation.
type WsProvider struct {
	c wsConn
}

// Close closes the client, aborting any in-flight requests.
func (ws *WsProvider) Close() {
	ws.c.Close()
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
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Jar: jar} //nolint:exhaustruct // Only the Jar field is used.
	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithHTTPClient(httpClient)}, options...)
	c, err := client.DialOptions(ctx, url, options...)
	if err != nil {
		return nil, err
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
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{Jar: jar} //nolint:exhaustruct // Only the Jar field is used.

	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithWebsocketDialer(dialer)}, options...)
	c, err := client.DialOptions(ctx, url, options...)
	if err != nil {
		return nil, err
	}

	return &WsProvider{c: c}, nil
}

//go:generate mockgen -destination=../internal/tests/mocks/rpcv10mock/rpc.go -package=rpcv10mock -source=provider.go
type RPCProvider interface {
	AddInvokeTransaction(
		ctx context.Context,
		invokeTxn *BroadcastInvokeTxnV3,
	) (AddInvokeTransactionResponse, error)
	AddDeclareTransaction(
		ctx context.Context,
		declareTransaction *BroadcastDeclareTxnV3,
	) (AddDeclareTransactionResponse, error)
	AddDeployAccountTransaction(
		ctx context.Context,
		deployAccountTransaction *BroadcastDeployAccountTxnV3,
	) (AddDeployAccountTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error)
	BlockWithReceipts(ctx context.Context, blockID BlockID) (interface{}, error)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (*BlockWithTxHashesOutput, error)
	BlockWithTxs(ctx context.Context, blockID BlockID) (*BlockWithTxsOutput, error)
	Call(ctx context.Context, call FunctionCall, block BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error)
	ClassHashAt(
		ctx context.Context,
		blockID BlockID,
		contractAddress *felt.Felt,
	) (*felt.Felt, error)
	CompiledCasm(ctx context.Context, classHash *felt.Felt) (*contracts.CasmClass, error)
	EstimateFee(
		ctx context.Context,
		requests []BroadcastTxn,
		simulationFlags []SimulationFlag,
		blockID BlockID,
	) ([]FeeEstimation, error)
	EstimateMessageFee(
		ctx context.Context,
		msg MsgFromL1,
		blockID BlockID,
	) (MessageFeeEstimation, error)
	Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	MessagesStatus(ctx context.Context, transactionHash NumAsHex) ([]MessageStatus, error)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	SimulateTransactions(
		ctx context.Context,
		blockID BlockID,
		txns []BroadcastTxn,
		simulationFlags []SimulationFlag,
	) ([]SimulatedTransaction, error)
	SpecVersion(ctx context.Context) (string, error)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error)
	StorageAt(
		ctx context.Context,
		contractAddress *felt.Felt,
		key string,
		blockID BlockID,
	) (string, error)
	StorageProof(
		ctx context.Context,
		storageProofInput StorageProofInput,
	) (*StorageProofResult, error)
	Syncing(ctx context.Context) (SyncStatus, error)
	TraceBlockTransactions(ctx context.Context, blockID BlockID) ([]Trace, error)
	TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error)
	TransactionByBlockIDAndIndex(
		ctx context.Context,
		blockID BlockID,
		index uint64,
	) (*BlockTransaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (*BlockTransaction, error)
	TransactionReceipt(
		ctx context.Context,
		transactionHash *felt.Felt,
	) (*TransactionReceiptWithBlockInfo, error)
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
		subBlockID SubscriptionBlockID,
	) (*client.ClientSubscription, error)
	SubscribeNewTransactions(
		ctx context.Context,
		newTxns chan<- *TxnWithHashAndStatus,
		options *SubNewTxnsInput,
	) (*client.ClientSubscription, error)
	SubscribeNewTransactionReceipts(
		ctx context.Context,
		txnReceipts chan<- *TransactionReceiptWithBlockInfo,
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
