package rpcv10

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

type MsgToL1 struct {
	// FromAddress The address of the L2 contract sending the message
	FromAddress *felt.Felt `json:"from_address"`
	// ToAddress The target L1 address the message is sent to
	ToAddress *felt.Felt `json:"to_address"`
	// Payload  The payload of the message
	Payload []*felt.Felt `json:"payload"`
}

type MsgFromL1 struct {
	// FromAddress The address of the L1 contract sending the message
	FromAddress string `json:"from_address"`
	// ToAddress The target L2 address the message is sent to
	ToAddress *felt.Felt `json:"to_address"`
	// EntryPointSelector The selector of the l1_handler in invoke in the target contract
	Selector *felt.Felt `json:"entry_point_selector"`
	// Payload  The payload of the message
	Payload []*felt.Felt `json:"payload"`
}

// MessageStatus represents the status of a message sent from an L1 transaction
// to an L2 contract.
type MessageStatus struct {
	// The hash of the L1_HANDLER transaction in L2 that contains the message
	Hash *felt.Felt `json:"transaction_hash"`
	// The finality status of the L1_HANDLER transaction, including the case the txn
	// is still in the mempool or
	// failed validation during the block construction phase
	FinalityStatus types.TxnFinalityStatus `json:"finality_status"`
	// The execution status of the L1_HANDLER transaction
	ExecutionStatus types.TxnExecutionStatus `json:"execution_status"`
	// The failure reason. Only appears if `execution_status` is REVERTED
	FailureReason string `json:"failure_reason,omitempty"`
}

type OrderedMsg struct {
	// The order of the message within the transaction
	Order   int `json:"order"`
	MsgToL1 MsgToL1
}

type FeePayment struct {
	Amount *felt.Felt      `json:"amount"`
	Unit   types.PriceUnit `json:"unit"`
}

// TransactionReceipt represents the common structure of a transaction receipt.
type TransactionReceipt struct {
	Hash               *felt.Felt               `json:"transaction_hash"`
	Type               types.TransactionType    `json:"type"`
	ActualFee          FeePayment               `json:"actual_fee"`
	FinalityStatus     types.TxnFinalityStatus  `json:"finality_status"`
	MessagesSent       []MsgToL1                `json:"messages_sent"`
	Events             []Event                  `json:"events"`
	ExecutionResources ExecutionResources       `json:"execution_resources"`
	ExecutionStatus    types.TxnExecutionStatus `json:"execution_status"`
	// Only present in case of a Deploy or DeployAccount transaction receipt
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
	// Only appears if the transaction is a L1Handler transaction
	MessageHash types.NumAsHex `json:"message_hash,omitempty"`
	// Only appears if execution_status is REVERTED
	RevertReason string `json:"revert_reason,omitempty"`
}

type ExecutionResources struct {
	// l1 gas consumed by this transaction, used for l2-->l1 messages and state
	// updates if blobs are not used
	L1Gas uint `json:"l1_gas"`
	// data gas consumed by this transaction, 0 if blobs are not used
	L1DataGas uint `json:"l1_data_gas"`
	// l2 gas consumed by this transaction, used for computation and calldata
	L2Gas uint `json:"l2_gas"`
}

// Transaction status result, including finality status and execution status
type TxnStatusResult struct {
	FinalityStatus  types.TxnStatus          `json:"finality_status"`
	ExecutionStatus types.TxnExecutionStatus `json:"execution_status,omitempty"`
	// the failure reason, only appears if execution_status is REVERTED
	FailureReason string `json:"failure_reason,omitempty"`
}

// The response of the starknet_subscribeTransactionStatus subscription.
type NewTxnStatus struct {
	TransactionHash *felt.Felt      `json:"transaction_hash"`
	Status          TxnStatusResult `json:"status"`
}

type TransactionReceiptWithBlockInfo struct {
	TransactionReceipt
	// If this field is missing, it means the receipt belongs to the pre-confirmed block
	BlockHash   *felt.Felt `json:"block_hash,omitempty"`
	BlockNumber uint       `json:"block_number"`
}
