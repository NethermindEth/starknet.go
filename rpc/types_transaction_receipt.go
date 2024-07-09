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
	//Payload  The payload of the message
	Payload []*felt.Felt `json:"payload"`
}

type MsgFromL1 struct {
	// FromAddress The address of the L1 contract sending the message
	FromAddress string `json:"from_address"`
	// ToAddress The target L2 address the message is sent to
	ToAddress *felt.Felt `json:"to_address"`
	// EntryPointSelector The selector of the l1_handler in invoke in the target contract
	Selector *felt.Felt `json:"entry_point_selector"`
	//Payload  The payload of the message
	Payload []*felt.Felt `json:"payload"`
}

type OrderedMsg struct {
	// The order of the message within the transaction
	Order   int `json:"order"`
	MsgToL1 MsgToL1
}

type FeePayment struct {
	Amount *felt.Felt     `json:"amount"`
	Unit   FeePaymentUnit `json:"unit"`
}

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
// - data: It takes a byte slice as input representing the JSON data to be unmarshaled
// Returns:
// - error: an error if the unmarshaling fails
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
//
// Returns:
// - []byte: a byte slice
// - error: an error if any
func (tt TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(tt))), nil
}

type ComputationResources struct {
	// The number of Cairo steps used
	Steps int `json:"steps"`
	// The number of unused memory cells (each cell is roughly equivalent to a step)
	MemoryHoles int `json:"memory_holes,omitempty"`
	// The number of RANGE_CHECK builtin instances
	RangeCheckApps int `json:"range_check_builtin_applications,omitempty"`
	// The number of Pedersen builtin instances
	PedersenApps int `json:"pedersen_builtin_applications,omitempty"`
	// The number of Poseidon builtin instances
	PoseidonApps int `json:"poseidon_builtin_applications,omitempty"`
	// The number of EC_OP builtin instances
	ECOPApps int `json:"ec_op_builtin_applications,omitempty"`
	// The number of ECDSA builtin instances
	ECDSAApps int `json:"ecdsa_builtin_applications,omitempty"`
	// The number of BITWISE builtin instances
	BitwiseApps int `json:"bitwise_builtin_applications,omitempty"`
	// The number of KECCAK builtin instances
	KeccakApps int `json:"keccak_builtin_applications,omitempty"`
	// The number of accesses to the segment arena
	SegmentArenaBuiltin int `json:"segment_arena_builtin,omitempty"`
}

func (er *ComputationResources) Validate() bool {
	if er.Steps == 0 || er.MemoryHoles == 0 || er.RangeCheckApps == 0 || er.PedersenApps == 0 ||
		er.PoseidonApps == 0 || er.ECOPApps == 0 || er.ECDSAApps == 0 || er.BitwiseApps == 0 ||
		er.KeccakApps == 0 || er.SegmentArenaBuiltin == 0 {
		return false
	}
	return true
}

type ExecutionResources struct {
	ComputationResources
	DataAvailability `json:"data_availability"`
}

type DataAvailability struct {
	// the gas consumed by this transaction's data, 0 if it uses data gas for DA
	L1Gas uint `json:"l1_gas"`
	// the data gas consumed by this transaction's data, 0 if it uses gas for DA
	L1DataGas uint `json:"l1_data_gas"`
}

type TxnStatus string

const (
	TxnStatus_Received       TxnStatus = "RECEIVED"
	TxnStatus_Rejected       TxnStatus = "REJECTED"
	TxnStatus_Accepted_On_L2 TxnStatus = "ACCEPTED_ON_L2"
	TxnStatus_Accepted_On_L1 TxnStatus = "ACCEPTED_ON_L1"
)

type TxnStatusResp struct {
	ExecutionStatus TxnExecutionStatus `json:"execution_status,omitempty"`
	FinalityStatus  TxnStatus          `json:"finality_status"`
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
	var aux struct {
		TransactionReceipt
		BlockHash   string `json:"block_hash,omitempty"`
		BlockNumber uint   `json:"block_number,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	tr.TransactionReceipt = aux.TransactionReceipt

	blockHash, err := new(felt.Felt).SetString(aux.BlockHash)
	if err != nil {
		return err
	}
	tr.BlockHash = blockHash
	tr.BlockNumber = aux.BlockNumber

	return nil
}
