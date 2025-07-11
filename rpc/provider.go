package rpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/gorilla/websocket"
	"golang.org/x/net/publicsuffix"
)

// rpcVersion is the version of the Starknet JSON-RPC specification that this SDK is compatible with.
// This should be updated when supporting new versions of the RPC specification.
const rpcVersion = "0.9.0"

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")

	// Warning messages for version compatibility
	warnVersionCheckFailed = "warning: could not check RPC version compatibility"
	//nolint:lll
	warnVersionMismatch = "warning: the RPC provider version is %s, and is different from the version %s implemented by the SDK. This may cause unexpected behaviour."
)

// Checks if the RPC provider version is compatible with the SDK version
// and prints a warning if they don't match.
func checkVersionCompatibility(provider *Provider) {
	version, err := provider.SpecVersion(context.Background())
	if err != nil {
		// Print a warning but don't fail
		fmt.Println(warnVersionCheckFailed, err)

		return
	}

	if !strings.Contains(version, rpcVersion) {
		fmt.Println(fmt.Sprintf(warnVersionMismatch, rpcVersion, version))
	}
}

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
func (p *WsProvider) Close() {
	p.c.Close()
}

// NewProvider creates a new HTTP rpc Provider instance.
func NewProvider(url string, options ...client.ClientOption) (*Provider, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Jar: jar} //nolint:exhaustruct
	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithHTTPClient(httpClient)}, options...)
	c, err := client.DialOptions(context.Background(), url, options...)
	if err != nil {
		return nil, err
	}

	provider := &Provider{c: c, chainID: ""}

	// Check version compatibility
	checkVersionCompatibility(provider)

	return provider, nil
}

// NewWebsocketProvider creates a new Websocket rpc Provider instance.
func NewWebsocketProvider(url string, options ...client.ClientOption) (*WsProvider, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{Jar: jar} //nolint:exhaustruct

	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithWebsocketDialer(dialer)}, options...)
	c, err := client.DialOptions(context.Background(), url, options...)
	if err != nil {
		return nil, err
	}

	return &WsProvider{c: c}, nil
}

//go:generate mockgen -destination=../mocks/mock_rpc_provider.go -package=mocks -source=provider.go api
type RpcProvider interface {
	AddInvokeTransaction(ctx context.Context, invokeTxn *BroadcastInvokeTxnV3) (*AddInvokeTransactionResponse, error)
	AddDeclareTransaction(ctx context.Context, declareTransaction *BroadcastDeclareTxnV3) (*AddDeclareTransactionResponse, error)
	AddDeployAccountTransaction(
		ctx context.Context,
		deployAccountTransaction *BroadcastDeployAccountTxnV3,
	) (*AddDeployAccountTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error)
	BlockWithReceipts(ctx context.Context, blockID BlockID) (interface{}, error)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error)
	BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error)
	Call(ctx context.Context, call FunctionCall, block BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error)
	ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	CompiledCasm(ctx context.Context, classHash *felt.Felt) (*contracts.CasmClass, error)
	EstimateFee(ctx context.Context, requests []BroadcastTxn, simulationFlags []SimulationFlag, blockID BlockID) ([]FeeEstimation, error)
	EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (MessageFeeEstimation, error)
	Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	GetStorageProof(ctx context.Context, storageProofInput StorageProofInput) (*StorageProofResult, error)
	GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResult, error)
	GetMessagesStatus(ctx context.Context, transactionHash NumAsHex) ([]MessageStatus, error)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	SimulateTransactions(
		ctx context.Context,
		blockID BlockID,
		txns []BroadcastTxn,
		simulationFlags []SimulationFlag,
	) ([]SimulatedTransaction, error)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error)
	StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error)
	SpecVersion(ctx context.Context) (string, error)
	Syncing(ctx context.Context) (*SyncStatus, error)
	TraceBlockTransactions(ctx context.Context, blockID BlockID) ([]Trace, error)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (*BlockTransaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (*BlockTransaction, error)
	TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (*TransactionReceiptWithBlockInfo, error)
	TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error)
}

type WebsocketProvider interface {
	SubscribeEvents(ctx context.Context, events chan<- *EmittedEvent, options *EventSubscriptionInput) (*client.ClientSubscription, error)
	SubscribeNewHeads(ctx context.Context, headers chan<- *BlockHeader, blockID BlockID) (*client.ClientSubscription, error)
	SubscribePendingTransactions(
		ctx context.Context,
		pendingTxns chan<- *PendingTxn,
		options *SubPendingTxnsInput,
	) (*client.ClientSubscription, error)
	SubscribeTransactionStatus(
		ctx context.Context,
		newStatus chan<- *NewTxnStatus,
		transactionHash *felt.Felt,
	) (*client.ClientSubscription, error)
}

//nolint:exhaustruct
var (
	_ RpcProvider       = &Provider{}
	_ WebsocketProvider = &WsProvider{}
)
