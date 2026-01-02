package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

func (provider *Provider) AddInvokeTransaction(
	ctx context.Context,
	invokeTxn *BroadcastInvokeTxnV3,
) (AddInvokeTransactionResponse, error) {
	return AddInvokeTransactionResponse{}, nil
}

func (provider *Provider) AddDeclareTransaction(
	ctx context.Context,
	declareTransaction *BroadcastDeclareTxnV3,
) (AddDeclareTransactionResponse, error) {
	return AddDeclareTransactionResponse{}, nil
}

func (provider *Provider) AddDeployAccountTransaction(
	ctx context.Context,
	deployAccountTransaction *BroadcastDeployAccountTxnV3,
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
	blockID BlockID,
) (uint64, error) {
	return 0, nil
}

func (provider *Provider) BlockWithReceipts(
	ctx context.Context,
	blockID BlockID,
) (interface{}, error) {
	return nil, nil
}

func (provider *Provider) BlockWithTxHashes(
	ctx context.Context,
	blockID BlockID,
) (interface{}, error) {
	return nil, nil
}

func (provider *Provider) BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error) {
	return nil, nil
}

func (provider *Provider) Call(
	ctx context.Context,
	call FunctionCall,
	block BlockID,
) ([]*felt.Felt, error) {
	return nil, nil
}

func (provider *Provider) ChainID(ctx context.Context) (string, error) {
	if provider.chainID != "" {
		return provider.chainID, nil
	}

	chainID, err := ChainID(ctx, provider.c)
	if err != nil {
		return "", err
	}
	provider.chainID = internalUtils.HexToShortStr(chainID)

	return provider.chainID, nil
}

func (provider *Provider) Class(
	ctx context.Context,
	blockID BlockID,
	classHash *felt.Felt,
) (ClassOutput, error) {
	return nil, nil
}

func (provider *Provider) ClassAt(
	ctx context.Context,
	blockID BlockID,
	contractAddress *felt.Felt,
) (ClassOutput, error) {
	return nil, nil
}

func (provider *Provider) ClassHashAt(
	ctx context.Context,
	blockID BlockID,
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
	requests []BroadcastTxn,
	simulationFlags []SimulationFlag,
	blockID BlockID,
) ([]FeeEstimation, error) {
	return nil, nil
}

func (provider *Provider) EstimateMessageFee(
	ctx context.Context,
	msg MsgFromL1,
	blockID BlockID,
) (MessageFeeEstimation, error) {
	return MessageFeeEstimation{}, nil
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
	blockID BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	return nil, nil
}

func (provider *Provider) SimulateTransactions(
	ctx context.Context,
	blockID BlockID,
	txns []BroadcastTxn,
	simulationFlags []SimulationFlag,
) ([]SimulatedTransaction, error) {
	return nil, nil
}

func (provider *Provider) SpecVersion(ctx context.Context) (string, error) {
	return "", nil
}

func (provider *Provider) StateUpdate(
	ctx context.Context,
	blockID BlockID,
) (*StateUpdateOutput, error) {
	return nil, nil
}

func (provider *Provider) StorageAt(
	ctx context.Context,
	contractAddress *felt.Felt,
	key string,
	blockID BlockID,
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
	blockID BlockID,
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
	blockID BlockID,
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
