package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

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

// MessageStatus represents the status of a message sent from an L1 transaction to an L2 contract.
type MessageStatus struct {
	// The hash of the L1_HANDLER transaction in L2 that contains the message
	Hash *felt.Felt `json:"transaction_hash"`
	// The finality status of the L1_HANDLER transaction, including the case the txn is still in the mempool or
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

type FeePayment struct {
	Amount *felt.Felt `json:"amount"`
	Unit   PriceUnit  `json:"unit"`
}

// Units in which the fee is given
type PriceUnit string

const (
	UnitWei PriceUnit = "WEI"
	UnitFri PriceUnit = "FRI"
)

// Representation of the unit WEI
type PriceUnitWei string

const (
	WeiUnit PriceUnitWei = "WEI"
)

// Representation of the unit FRI
type PriceUnitFri string

const (
	FriUnit PriceUnitFri = "FRI"
)

// Unmarshals the JSON data into a PriceUnit.
func (f *PriceUnit) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "WEI":
		*f = UnitWei
	case "FRI":
		*f = UnitFri
	default:
		return fmt.Errorf("unsupported price unit: %s", data)
	}

	return nil
}

// Unmarshals the JSON data into a PriceUnitWei.
func (f *PriceUnitWei) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	if unquoted != string(WeiUnit) {
		return fmt.Errorf("price unit should be WEI, got: %s", data)
	}

	*f = WeiUnit

	return nil
}

// Unmarshals the JSON data into a PriceUnitFri.
func (f *PriceUnitFri) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	if unquoted != string(FriUnit) {
		return fmt.Errorf("price unit should be FRI, got: %s", data)
	}

	*f = FriUnit

	return nil
}

// TransactionReceipt represents the common structure of a transaction receipt.
type TransactionReceipt struct {
	Hash               *felt.Felt         `json:"transaction_hash"`
	Type               TransactionType    `json:"type"`
	ActualFee          FeePayment         `json:"actual_fee"`
	FinalityStatus     TxnFinalityStatus  `json:"finality_status"`
	MessagesSent       []MsgToL1          `json:"messages_sent"`
	Events             []Event            `json:"events"`
	ExecutionResources ExecutionResources `json:"execution_resources"`
	ExecutionStatus    TxnExecutionStatus `json:"execution_status"`
	// Only present in case of a Deploy or DeployAccount transaction receipt
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
	// Only appears if the transaction is a L1Handler transaction
	MessageHash NumAsHex `json:"message_hash,omitempty"`
	// Only appears if execution_status is REVERTED
	RevertReason string `json:"revert_reason,omitempty"`
}

type TransactionType string

const (
	TransactionType_Declare       TransactionType = "DECLARE"
	TransactionType_DeployAccount TransactionType = "DEPLOY_ACCOUNT"
	TransactionType_Deploy        TransactionType = "DEPLOY"
	TransactionType_Invoke        TransactionType = "INVOKE"
	TransactionType_L1Handler     TransactionType = "L1_HANDLER"
)

// UnmarshalJSON unmarshals the JSON data into a TransactionType.
//
// The function modifies the value of the TransactionType pointer tt based on the unmarshaled data.
// The supported JSON values and their corresponding TransactionType values are:
//   - "DECLARE" maps to TransactionType_Declare
//   - "DEPLOY_ACCOUNT" maps to TransactionType_DeployAccount
//   - "DEPLOY" maps to TransactionType_Deploy
//   - "INVOKE" maps to TransactionType_Invoke
//   - "L1_HANDLER" maps to TransactionType_L1Handler
//
// If none of the supported values match the input data, the function returns an error.
//
//	nil if the unmarshaling is successful.
//
// Parameters:
//   - data: It takes a byte slice as input representing the JSON data to be unmarshaled
//
// Returns:
//   - error: an error if the unmarshaling fails
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
		return fmt.Errorf("unsupported transaction type: %s", data)
	}

	return nil
}

// MarshalJSON marshals the TransactionType to JSON.
//
// Returns:
//   - []byte: a byte slice
//   - error: an error if any
func (tt TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(tt))), nil
}

type ExecutionResources struct {
	// l1 gas consumed by this transaction, used for l2-->l1 messages and state updates if blobs are not used
	L1Gas uint `json:"l1_gas"`
	// data gas consumed by this transaction, 0 if blobs are not used
	L1DataGas uint `json:"l1_data_gas"`
	// l2 gas consumed by this transaction, used for computation and calldata
	L2Gas uint `json:"l2_gas"`
}

type TxnStatus string

const (
	TxnStatus_Received       TxnStatus = "RECEIVED"
	TxnStatus_Candidate      TxnStatus = "CANDIDATE"
	TxnStatus_Pre_confirmed  TxnStatus = "PRE_CONFIRMED"
	TxnStatus_Accepted_On_L2 TxnStatus = "ACCEPTED_ON_L2"
	TxnStatus_Accepted_On_L1 TxnStatus = "ACCEPTED_ON_L1"
)

// Transaction status result, including finality status and execution status
type TxnStatusResult struct {
	FinalityStatus  TxnStatus          `json:"finality_status"`
	ExecutionStatus TxnExecutionStatus `json:"execution_status,omitempty"`
	FailureReason   string             `json:"failure_reason,omitempty"`
}

// The response of the starknet_subscribeTransactionStatus subscription.
type NewTxnStatus struct {
	TransactionHash *felt.Felt      `json:"transaction_hash"`
	Status          TxnStatusResult `json:"status"`
}

type TransactionReceiptWithBlockInfo struct {
	TransactionReceipt
	BlockHash   *felt.Felt `json:"block_hash,omitempty"`
	BlockNumber uint       `json:"block_number,omitempty"`
}

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

func (tr *TransactionReceiptWithBlockInfo) UnmarshalJSON(data []byte) error {
	type Alias TransactionReceiptWithBlockInfo
	var txnResp Alias

	if err := json.Unmarshal(data, &txnResp); err != nil {
		return err
	}

	// If the block hash is nil (txn from pre_confirmed block), set it to felt.Zero to avoid nil pointer dereference
	if txnResp.BlockHash == nil {
		txnResp.BlockHash = new(felt.Felt)
	}

	*tr = TransactionReceiptWithBlockInfo(txnResp)

	return nil
}
