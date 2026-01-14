package types

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// TransactionReceipt represents the common structure of a transaction receipt.
type TransactionReceipt struct {
	Hash               *felt.Felt            `json:"transaction_hash"`
	Type               types.TransactionType `json:"type"`
	ActualFee          types.FeePayment      `json:"actual_fee"`
	FinalityStatus     TxnFinalityStatus     `json:"finality_status"`
	MessagesSent       []MsgToL1             `json:"messages_sent"`
	Events             []Event               `json:"events"`
	ExecutionResources ExecutionResources    `json:"execution_resources"`
	ExecutionStatus    TxnExecutionStatus    `json:"execution_status"`
	// Only present in case of a Deploy or DeployAccount transaction receipt
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
	// Only appears if the transaction is a L1Handler transaction
	MessageHash NumAsHex `json:"message_hash,omitempty"`
	// Only appears if execution_status is REVERTED
	RevertReason string `json:"revert_reason,omitempty"`
}

type TransactionReceiptWithBlockInfo struct {
	TransactionReceipt
	// If this field is missing, it means the receipt belongs to the pre-confirmed block
	BlockHash   *felt.Felt `json:"block_hash,omitempty"`
	BlockNumber uint       `json:"block_number"`
}
