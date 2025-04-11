package rpc

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/gorilla/websocket"
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
	httpClient := &http.Client{Jar: jar}
	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithHTTPClient(httpClient)}, options...)
	c, err := client.DialOptions(context.Background(), url, options...)

	if err != nil {
		return nil, err
	}

	return &Provider{c: c}, nil
}

// NewWebsocketProvider creates a new Websocket rpc Provider instance.
func NewWebsocketProvider(url string, options ...client.ClientOption) (*WsProvider, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{Jar: jar}

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
	AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction *BroadcastDeployAccountTxnV3) (*AddDeployAccountTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	Number(ctx context.Context) (uint64, error)
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
	EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimation, error)
	Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	GetStorageProof(ctx context.Context, storageProofInput StorageProofInput) (*StorageProofResult, error)
	GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResp, error)
	GetMessagesStatus(ctx context.Context, transactionHash NumAsHex) ([]MessageStatusResp, error)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	SimulateTransactions(ctx context.Context, blockID BlockID, txns []BroadcastTxn, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error)
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
	SubscribePendingTransactions(ctx context.Context, pendingTxns chan<- *SubPendingTxns, options *SubPendingTxnsInput) (*client.ClientSubscription, error)
	SubscribeTransactionStatus(ctx context.Context, newStatus chan<- *NewTxnStatusResp, transactionHash *felt.Felt) (*client.ClientSubscription, error)
}

var _ RpcProvider = &Provider{}
var _ WebsocketProvider = &WsProvider{}
