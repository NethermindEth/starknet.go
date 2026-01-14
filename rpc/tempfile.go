package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

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
	TransactionByHash(ctx context.Context, hash *felt.Felt) (types.BlockTransaction, error)
	TransactionReceipt(
		ctx context.Context,
		transactionHash *felt.Felt,
	) (*TransactionReceiptWithBlockInfo, error)
	TransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResult, error)
}
