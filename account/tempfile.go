package account

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// @todo update docs for the entire package
// add tests where needed.

type providerWrapper struct {
	chainID string
	version RPCVersion

	rpcv9  rpcv9.RPCProvider
	rpcv10 rpcv10.RPCProvider
}

// NewProviderWrapper creates a new HTTP rpc Provider instance.
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
func NewProviderWrapper(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*providerWrapper, error) {
	c, err := internal.NewHTTPClient(ctx, url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	rawNodeVersion, err := rpcv10.SpecVersion(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to get the node's RPC spec version: %w", err)
	}

	var RPCVersion RPCVersion
	err = RPCVersion.UnmarshalJSON([]byte(rawNodeVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node version: %w", err)
	}
	var provider providerWrapper
	provider.version = RPCVersion

	switch RPCVersion {
	case RPCVersion9:
		rpcv9, err := rpcv9.NewProvider(ctx, url, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create RPCv9 provider: %w", err)
		}
		provider.rpcv9 = rpcv9
	case RPCVersion10:
		rpcv10, err := rpcv10.NewProvider(ctx, url, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create RPCv10 provider: %w", err)
		}
		provider.rpcv10 = rpcv10
	}

	return &provider, nil
}

// @new
type RPCProvider interface {
	*rpcv10.Provider | *rpcv9.Provider
}

func NewProviderWrapperFrom[P RPCProvider](provider P) *providerWrapper {
	var wrapper providerWrapper

	switch p := any(provider).(type) {
	case *rpcv9.Provider:
		wrapper.rpcv9 = p
		wrapper.version = RPCVersion9
	case *rpcv10.Provider:
		wrapper.rpcv10 = p
		wrapper.version = RPCVersion10
	}

	return &wrapper
}

type RPCProviderV10Copy interface {
	// AddInvokeTransaction(
	// 	ctx context.Context,
	// 	invokeTxn *types.BroadcastInvokeTxnV3,
	// ) (AddInvokeTransactionResponse, error)
	// AddDeclareTransaction(
	// 	ctx context.Context,
	// 	declareTransaction *types.BroadcastDeclareTxnV3,
	// ) (AddDeclareTransactionResponse, error)
	// AddDeployAccountTransaction(
	// 	ctx context.Context,
	// 	deployAccountTransaction *types.BroadcastDeployAccountTxnV3,
	// ) (AddDeployAccountTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (uint64, *felt.Felt, error) //****modified****
	// BlockNumber(ctx context.Context) (uint64, error)
	// BlockTransactionCount(ctx context.Context, blockID types.BlockID) (uint64, error)
	// BlockWithReceipts(ctx context.Context, blockID types.BlockID) (interface{}, error)
	// BlockWithTxHashes(ctx context.Context, blockID types.BlockID) (interface{}, error)
	// BlockWithTxs(ctx context.Context, blockID types.BlockID) (interface{}, error)
	Call(ctx context.Context, call types.FunctionCall, block types.BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	// Class(ctx context.Context, blockID types.BlockID, classHash *felt.Felt) (ClassOutput, error)
	// ClassAt(
	// 	ctx context.Context,
	// 	blockID types.BlockID,
	// 	contractAddress *felt.Felt,
	// ) (ClassOutput, error)
	// ClassHashAt(
	// 	ctx context.Context,
	// 	blockID types.BlockID,
	// 	contractAddress *felt.Felt,
	// ) (*felt.Felt, error)
	// CompiledCasm(ctx context.Context, classHash *felt.Felt) (*contracts.CasmClass, error)
	EstimateFee(
		ctx context.Context,
		requests []types.BroadcastTxn,
		simulationFlags []types.SimulationFlag,
		blockID types.BlockID,
	) ([]types.FeeEstimation, error)
	// EstimateMessageFee(
	// 	ctx context.Context,
	// 	msg MsgFromL1,
	// 	blockID types.BlockID,
	// ) (types.MessageFeeEstimation, error)
	// Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	// MessagesStatus(ctx context.Context, transactionHash NumAsHex) ([]MessageStatus, error)
	Nonce(
		ctx context.Context,
		blockID types.BlockID,
		contractAddress *felt.Felt,
	) (*felt.Felt, error)
	// SimulateTransactions(
	// 	ctx context.Context,
	// 	blockID types.BlockID,
	// 	txns []types.BroadcastTxn,
	// 	simulationFlags []types.SimulationFlag,
	// ) ([]SimulatedTransaction, error)
	// SpecVersion(ctx context.Context) (string, error)
	// StateUpdate(ctx context.Context, blockID types.BlockID) (*StateUpdateOutput, error)
	// StorageAt(
	// 	ctx context.Context,
	// 	contractAddress *felt.Felt,
	// 	key string,
	// 	blockID types.BlockID,
	// ) (string, error)
	// StorageProof(
	// 	ctx context.Context,
	// 	storageProofInput StorageProofInput,
	// ) (*StorageProofResult, error)
	IsSyncing(ctx context.Context) (bool, error) //****modified****
	// TraceBlockTransactions(ctx context.Context, blockID types.BlockID) ([]Trace, error)
	// TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error)
	// TransactionByBlockIDAndIndex(
	// 	ctx context.Context,
	// 	blockID types.BlockID,
	// 	index uint64,
	// ) (*BlockTransaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (types.BlockTransaction, error) //****modified****
	TransactionReceipt(
		ctx context.Context,
		transactionHash *felt.Felt,
	) (types.TransactionReceiptWithBlockInfo, error) //****modified****
	TransactionStatus(
		ctx context.Context,
		transactionHash *felt.Felt,
	) (types.TxnStatusResult, error) //****modified****
}

type OtherMethods interface {
	AsV9() rpcv9.RPCProvider
	AsV10() rpcv10.RPCProvider
	EstimateTip(ctx context.Context, multiplier float64) (tip types.U64, err error)
	SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
}

// implementing the methods
func (p *providerWrapper) BlockHashAndNumber(ctx context.Context) (uint64, *felt.Felt, error)
func (p *providerWrapper) Call(ctx context.Context, call types.FunctionCall, block types.BlockID) ([]*felt.Felt, error)
func (p *providerWrapper) ChainID(ctx context.Context) (string, error)
func (p *providerWrapper) EstimateFee(
	ctx context.Context,
	requests []types.BroadcastTxn,
	simulationFlags []types.SimulationFlag,
	blockID types.BlockID,
) ([]types.FeeEstimation, error)
func (p *providerWrapper) Nonce(
	ctx context.Context,
	blockID types.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error)
func (p *providerWrapper) IsSyncing(ctx context.Context) (bool, error)
func (p *providerWrapper) TransactionByHash(ctx context.Context, hash *felt.Felt) (types.BlockTransaction, error)
func (p *providerWrapper) TransactionReceipt(
	ctx context.Context,
	transactionHash *felt.Felt,
) (types.TransactionReceiptWithBlockInfo, error)
func (p *providerWrapper) TransactionStatus(
	ctx context.Context,
	transactionHash *felt.Felt,
) (types.TxnStatusResult, error)

func (p *providerWrapper) EstimateTip(ctx context.Context, multiplier float64) (tip types.U64, err error)
func (p *providerWrapper) SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
func (p *providerWrapper) AsV9() rpcv9.RPCProvider {
	return p.rpcv9
}

func (p *providerWrapper) AsV10() rpcv10.RPCProvider {
	return p.rpcv10
}

func (p *providerWrapper) Version() RPCVersion {
	return p.version
}

// @todo add tests for this type

type RPCVersion int

const (
	RPCVersion9 RPCVersion = iota
	RPCVersion10
)

func (v RPCVersion) String() string {
	switch v {
	case RPCVersion9:
		return "0.9.0"
	case RPCVersion10:
		return "0.10.1"
	}
	return ""
}

func (v RPCVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *RPCVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	semversion, err := semver.NewVersion(s)
	if err != nil {
		return fmt.Errorf("failed to parse version to semver: %w", err)
	}

	switch {
	case semversion.Compare(semver.MustParse(RPCVersion9.String())) == 0:
		*v = RPCVersion9
	case semversion.Compare(semver.MustParse(RPCVersion10.String())) == 0:
		*v = RPCVersion10
	default:
		return fmt.Errorf("invalid RPC version: %s", s)
	}
	return nil
}
