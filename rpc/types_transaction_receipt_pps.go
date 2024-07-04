package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

// FeePayment represents the fee payment details.
type FeePayment struct {
	Amount *felt.Felt     `json:"amount"`
	Unit   FeePaymentUnit `json:"unit"`
}

// FeePaymentUnit is the unit for fee payment.
type FeePaymentUnit string

const (
	UnitWei  FeePaymentUnit = "WEI"
	UnitStrk FeePaymentUnit = "FRI"
)

// TransactionReceipt represents the common structure of a transaction receipt.
type TransactionReceipt struct {
	TransactionHash    *felt.Felt         `json:"transaction_hash"`
	ActualFee          FeePayment         `json:"actual_fee"`
	ExecutionStatus    TxnExecutionStatus `json:"execution_status"`
	FinalityStatus     TxnFinalityStatus  `json:"finality_status"`
	Type               TransactionType    `json:"type,omitempty"`
	MessagesSent       []MsgToL1          `json:"messages_sent"`
	RevertReason       string             `json:"revert_reason,omitempty"`
	Events             []Event            `json:"events"`
	ExecutionResources ExecutionResources `json:"execution_resources"`
	ContractAddress    *felt.Felt         `json:"contract_address,omitempty"`
	MessageHash        NumAsHex           `json:"message_hash,omitempty"`
}

// Hash returns the transaction hash associated with the TransactionReceipt.
func (tr TransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

// GetExecutionStatus returns the execution status of the TransactionReceipt.
func (tr TransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

// TransactionType represents the type of transaction.
type TransactionType string

const (
	TransactionType_Declare       TransactionType = "DECLARE"
	TransactionType_DeployAccount TransactionType = "DEPLOY_ACCOUNT"
	TransactionType_Deploy        TransactionType = "DEPLOY"
	TransactionType_Invoke        TransactionType = "INVOKE"
	TransactionType_L1Handler     TransactionType = "L1_HANDLER"
)

// UnmarshalJSON unmarshals the JSON data into a TransactionType.
func (tt *TransactionType) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "DECLARE":
		*tt = TransactionType_Declare
	case "DEPLOY_ACCOUNT":
		*tt = TransactionType_DeployAccount
	case "DEPLOY":
		*tt = TransactionType_Deploy
	case "INVOKE":
		*tt = TransactionType_Invoke
	case "L1_HANDLER":
		*tt = TransactionType_L1Handler
	default:
		return fmt.Errorf("unsupported type: %s", data)
	}

	return nil
}

// MarshalJSON marshals the TransactionType to JSON.
func (tt TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(tt))), nil
}

// ComputationResources represents the computation resources used by the transaction.
type ComputationResources struct {
	Steps               int `json:"steps"`
	MemoryHoles         int `json:"memory_holes,omitempty"`
	RangeCheckApps      int `json:"range_check_builtin_applications,omitempty"`
	PedersenApps        int `json:"pedersen_builtin_applications,omitempty"`
	PoseidonApps        int `json:"poseidon_builtin_applications,omitempty"`
	ECOPApps            int `json:"ec_op_builtin_applications,omitempty"`
	ECDSAApps           int `json:"ecdsa_builtin_applications,omitempty"`
	BitwiseApps         int `json:"bitwise_builtin_applications,omitempty"`
	KeccakApps          int `json:"keccak_builtin_applications,omitempty"`
	SegmentArenaBuiltin int `json:"segment_arena_builtin,omitempty"`
}

// Validate checks if the fields are non-zero to match the starknet-specs.
func (er *ComputationResources) Validate() bool {
	if er.Steps == 0 || er.MemoryHoles == 0 || er.RangeCheckApps == 0 || er.PedersenApps == 0 ||
		er.PoseidonApps == 0 || er.ECOPApps == 0 || er.ECDSAApps == 0 || er.BitwiseApps == 0 ||
		er.KeccakApps == 0 || er.SegmentArenaBuiltin == 0 {
		return false
	}
	return true
}

// ExecutionResources represents the resources consumed by the transaction.
type ExecutionResources struct {
	ComputationResources
	DataAvailability `json:"data_availability"`
}

// DataAvailability represents the data availability details.
type DataAvailability struct {
	L1Gas     uint `json:"l1_gas"`
	L1DataGas uint `json:"l1_data_gas"`
}

// TransactionReceiptWithBlockInfo represents a transaction receipt with block info.
type TransactionReceiptWithBlockInfo struct {
	TransactionReceipt
	BlockHash   *felt.Felt `json:"block_hash,omitempty"`
	BlockNumber uint       `json:"block_number,omitempty"`
}

// MarshalJSON marshals the TransactionReceiptWithBlockInfo to JSON.
func (t *TransactionReceiptWithBlockInfo) MarshalJSON() ([]byte, error) {
	aux := &struct {
		TransactionReceipt
		BlockHash   string `json:"block_hash,omitempty"`
		BlockNumber uint   `json:"block_number,omitempty"`
	}{
		TransactionReceipt: t.TransactionReceipt,
		BlockHash:          t.BlockHash.String(),
		BlockNumber:        t.BlockNumber,
	}

	return json.Marshal(aux)
}

// UnmarshalJSON unmarshals the given JSON data into a TransactionReceiptWithBlockInfo.
func (tr *TransactionReceiptWithBlockInfo) UnmarshalJSON(data []byte) error {
	var aux struct {
		TransactionHash    *felt.Felt         `json:"transaction_hash"`
		ActualFee          FeePayment         `json:"actual_fee"`
		ExecutionStatus    TxnExecutionStatus `json:"execution_status"`
		FinalityStatus     TxnFinalityStatus  `json:"finality_status"`
		Type               TransactionType    `json:"type,omitempty"`
		MessagesSent       []MsgToL1          `json:"messages_sent"`
		RevertReason       string             `json:"revert_reason,omitempty"`
		Events             []Event            `json:"events"`
		ExecutionResources ExecutionResources `json:"execution_resources"`
		ContractAddress    *felt.Felt         `json:"contract_address,omitempty"`
		MessageHash        NumAsHex           `json:"message_hash,omitempty"`
		BlockHash          string             `json:"block_hash,omitempty"`
		BlockNumber        uint               `json:"block_number,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	tr.TransactionReceipt = TransactionReceipt{
		TransactionHash:    aux.TransactionHash,
		ActualFee:          aux.ActualFee,
		ExecutionStatus:    aux.ExecutionStatus,
		FinalityStatus:     aux.FinalityStatus,
		Type:               aux.Type,
		MessagesSent:       aux.MessagesSent,
		RevertReason:       aux.RevertReason,
		Events:             aux.Events,
		ExecutionResources: aux.ExecutionResources,
		ContractAddress:    aux.ContractAddress,
		MessageHash:        aux.MessageHash,
	}

	blockHash, err := new(felt.Felt).SetString(aux.BlockHash)
	if err != nil {
		return err
	}
	tr.BlockHash = blockHash
	tr.BlockNumber = aux.BlockNumber

	return nil
}

type MsgToL1 struct {
	// FromAddress The address of the L2 contract sending the message
	FromAddress *felt.Felt `json:"from_address"`
	// ToAddress The target L1 address the message is sent to
	ToAddress *felt.Felt `json:"to_address"`
	//Payload  The payload of the message
	Payload []*felt.Felt `json:"payload"`
}
