package rpcv10

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// @todo implement all the methods

func (provider *Provider) AddInvokeTransaction(
	ctx context.Context,
	invokeTxn *types.BroadcastInvokeTxnV3,
) (AddInvokeTransactionResponse, error) {
	return AddInvokeTransactionResponse{}, nil
}

func (provider *Provider) AddDeclareTransaction(
	ctx context.Context,
	declareTransaction *types.BroadcastDeclareTxnV3,
) (AddDeclareTransactionResponse, error) {
	return AddDeclareTransactionResponse{}, nil
}

func (provider *Provider) AddDeployAccountTransaction(
	ctx context.Context,
	deployAccountTransaction *types.BroadcastDeployAccountTxnV3,
) (AddDeployAccountTransactionResponse, error) {
	return AddDeployAccountTransactionResponse{}, nil
}

func (provider *Provider) BlockHashAndNumber(
	ctx context.Context,
) (*BlockHashAndNumberOutput, error) {
	return nil, nil
}

func (provider *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, nil
}

func (provider *Provider) BlockTransactionCount(
	ctx context.Context,
	blockID types.BlockID,
) (uint64, error) {
	return 0, nil
}

func (provider *Provider) BlockWithReceipts(
	ctx context.Context,
	blockID types.BlockID,
) (interface{}, error) {
	return nil, nil
}

func (provider *Provider) BlockWithTxHashes(
	ctx context.Context,
	blockID types.BlockID,
) (interface{}, error) {
	return nil, nil
}

func (provider *Provider) BlockWithTxs(
	ctx context.Context,
	blockID types.BlockID,
) (interface{}, error) {
	return nil, nil
}

func (provider *Provider) Call(
	ctx context.Context,
	call types.FunctionCall,
	block types.BlockID,
) ([]*felt.Felt, error) {
	return nil, nil
}

func (provider *Provider) ChainID(ctx context.Context) (string, error) {
	// if provider.chainID != "" {
	// 	return provider.chainID, nil
	// }

	// chainID, err := ChainID(ctx, provider.c)
	// if err != nil {
	// 	return "", err
	// }
	// provider.chainID = internalUtils.HexToShortStr(chainID)

	// return provider.chainID, nil
	return "", nil
}

func (provider *Provider) Class(
	ctx context.Context,
	blockID types.BlockID,
	classHash *felt.Felt,
) (ClassOutput, error) {
	return nil, nil
}

func (provider *Provider) ClassAt(
	ctx context.Context,
	blockID types.BlockID,
	contractAddress *felt.Felt,
) (ClassOutput, error) {
	return nil, nil
}

func (provider *Provider) ClassHashAt(
	ctx context.Context,
	blockID types.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	return nil, nil
}

func (provider *Provider) CompiledCasm(
	ctx context.Context,
	classHash *felt.Felt,
) (*contracts.CasmClass, error) {
	return nil, nil
}

func (provider *Provider) EstimateFee(
	ctx context.Context,
	requests []types.BroadcastTxn,
	simulationFlags []SimulationFlag,
	blockID types.BlockID,
) ([]types.FeeEstimation, error) {
	return nil, nil
}

func (provider *Provider) EstimateMessageFee(
	ctx context.Context,
	msg MsgFromL1,
	blockID types.BlockID,
) (types.MessageFeeEstimation, error) {
	return types.MessageFeeEstimation{}, nil
}

func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventChunk, error) {
	return nil, nil
}

func (provider *Provider) MessagesStatus(
	ctx context.Context,
	transactionHash NumAsHex,
) ([]MessageStatus, error) {
	return nil, nil
}

func (provider *Provider) Nonce(
	ctx context.Context,
	blockID types.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	return nil, nil
}

func (provider *Provider) SimulateTransactions(
	ctx context.Context,
	blockID types.BlockID,
	txns []types.BroadcastTxn,
	simulationFlags []SimulationFlag,
) ([]SimulatedTransaction, error) {
	return nil, nil
}

func (provider *Provider) SpecVersion(ctx context.Context) (string, error) {
	return "", nil
}

func (provider *Provider) StateUpdate(
	ctx context.Context,
	blockID types.BlockID,
) (*StateUpdateOutput, error) {
	return nil, nil
}

func (provider *Provider) StorageAt(
	ctx context.Context,
	contractAddress *felt.Felt,
	key string,
	blockID types.BlockID,
) (string, error) {
	return "", nil
}

func (provider *Provider) StorageProof(
	ctx context.Context,
	storageProofInput StorageProofInput,
) (*StorageProofResult, error) {
	return nil, nil
}

func (provider *Provider) Syncing(ctx context.Context) (SyncStatus, error) {
	return SyncStatus{}, nil
}

func (provider *Provider) TraceBlockTransactions(
	ctx context.Context,
	blockID types.BlockID,
) ([]Trace, error) {
	return nil, nil
}

func (provider *Provider) TraceTransaction(
	ctx context.Context,
	transactionHash *felt.Felt,
) (TxnTrace, error) {
	return nil, nil
}

func (provider *Provider) TransactionByBlockIDAndIndex(
	ctx context.Context,
	blockID types.BlockID,
	index uint64,
) (*BlockTransaction, error) {
	return nil, nil
}

func (provider *Provider) TransactionByHash(
	ctx context.Context,
	hash *felt.Felt,
) (*BlockTransaction, error) {
	return nil, nil
}

func (provider *Provider) TransactionReceipt(
	ctx context.Context,
	transactionHash *felt.Felt,
) (*TransactionReceiptWithBlockInfo, error) {
	return nil, nil
}

func (provider *Provider) TransactionStatus(
	ctx context.Context,
	transactionHash *felt.Felt,
) (*TxnStatusResult, error) {
	return nil, nil
}

// WSProvider methods

func (ws *WsProvider) SubscribeEvents(
	ctx context.Context,
	events chan<- *EmittedEventWithFinalityStatus,
	options *EventSubscriptionInput,
) (*client.ClientSubscription, error) {
	return nil, nil
}

func (ws *WsProvider) SubscribeNewHeads(
	ctx context.Context,
	headers chan<- *BlockHeader,
	subBlockID types.SubscriptionBlockID,
) (*client.ClientSubscription, error) {
	return nil, nil
}

func (ws *WsProvider) SubscribeNewTransactions(
	ctx context.Context,
	newTxns chan<- *TxnWithHashAndStatus,
	options *SubNewTxnsInput,
) (*client.ClientSubscription, error) {
	return nil, nil
}

func (ws *WsProvider) SubscribeNewTransactionReceipts(
	ctx context.Context,
	txnReceipts chan<- *TransactionReceiptWithBlockInfo,
	options *SubNewTxnReceiptsInput,
) (*client.ClientSubscription, error) {
	return nil, nil
}

func (ws *WsProvider) SubscribeTransactionStatus(
	ctx context.Context,
	newStatus chan<- *NewTxnStatus,
	transactionHash *felt.Felt,
) (*client.ClientSubscription, error) {
	return nil, nil
}
