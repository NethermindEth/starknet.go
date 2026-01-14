package rpcv10

import (
	"github.com/NethermindEth/juno/core/felt"
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
	FinalityStatus TxnFinalityStatus `json:"finality_status"`
	// The execution status of the L1_HANDLER transaction
	ExecutionStatus TxnExecutionStatus `json:"execution_status"`
	// The failure reason. Only appears if `execution_status` is REVERTED
	FailureReason string `json:"failure_reason,omitempty"`
}

type OrderedMsg struct {
	// The order of the message within the transaction
	Order   int `json:"order"`
	MsgToL1 MsgToL1
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

type TxnStatus string

const (
	TxnStatusReceived     TxnStatus = "RECEIVED"
	TxnStatusCandidate    TxnStatus = "CANDIDATE"
	TxnStatusPreConfirmed TxnStatus = "PRE_CONFIRMED"
	TxnStatusAcceptedOnL2 TxnStatus = "ACCEPTED_ON_L2"
	TxnStatusAcceptedOnL1 TxnStatus = "ACCEPTED_ON_L1"
)

// Transaction status result, including finality status and execution status
type TxnStatusResult struct {
	FinalityStatus  TxnStatus          `json:"finality_status"`
	ExecutionStatus TxnExecutionStatus `json:"execution_status,omitempty"`
	// the failure reason, only appears if execution_status is REVERTED
	FailureReason string `json:"failure_reason,omitempty"`
}

// The response of the starknet_subscribeTransactionStatus subscription.
type NewTxnStatus struct {
	TransactionHash *felt.Felt      `json:"transaction_hash"`
	Status          TxnStatusResult `json:"status"`
}
