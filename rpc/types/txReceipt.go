package types

import (
	"fmt"
	"strconv"

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

type TxnExecutionStatus string

const (
	TxnExecutionStatusSUCCEEDED TxnExecutionStatus = "SUCCEEDED"
	TxnExecutionStatusREVERTED  TxnExecutionStatus = "REVERTED"
)

// UnmarshalJSON unmarshals the JSON data into a TxnExecutionStatus struct.
func (ex *TxnExecutionStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "SUCCEEDED":
		*ex = TxnExecutionStatusSUCCEEDED
	case "REVERTED":
		*ex = TxnExecutionStatusREVERTED
	default:
		return fmt.Errorf("unsupported execution status: %s", data)
	}

	return nil
}

// MarshalJSON returns the JSON encoding of the TxnExecutionStatus.
func (ex TxnExecutionStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(ex))), nil
}

// String returns the string representation of the TxnExecutionStatus.
func (ex TxnExecutionStatus) String() string {
	return string(ex)
}

type TxnFinalityStatus string

const (
	TxnFinalityStatusPreConfirmed TxnFinalityStatus = "PRE_CONFIRMED"
	TxnFinalityStatusAcceptedOnL2 TxnFinalityStatus = "ACCEPTED_ON_L2"
	TxnFinalityStatusAcceptedOnL1 TxnFinalityStatus = "ACCEPTED_ON_L1"
)

// UnmarshalJSON unmarshals the JSON data into a TxnFinalityStatus.
//
// Parameters:
//   - data: It takes a byte slice as a parameter, which represents the JSON data to
//     be unmarshalled
//
// Returns:
//   - error: an error if the unmarshaling fails
func (fs *TxnFinalityStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "PRE_CONFIRMED":
		*fs = TxnFinalityStatusPreConfirmed
	case "ACCEPTED_ON_L2":
		*fs = TxnFinalityStatusAcceptedOnL2
	case "ACCEPTED_ON_L1":
		*fs = TxnFinalityStatusAcceptedOnL1
	default:
		return fmt.Errorf("unsupported finality status: %s", data)
	}

	return nil
}

// MarshalJSON marshals the TxnFinalityStatus into JSON.
func (fs TxnFinalityStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(fs))), nil
}

// String returns the string representation of the TxnFinalityStatus.
func (fs TxnFinalityStatus) String() string {
	return string(fs)
}
