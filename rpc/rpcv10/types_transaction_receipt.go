package rpcv10

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

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
	MsgToL1 types.MsgToL1
}

// The response of the starknet_subscribeTransactionStatus subscription.
type NewTxnStatus struct {
	TransactionHash *felt.Felt            `json:"transaction_hash"`
	Status          types.TxnStatusResult `json:"status"`
}
