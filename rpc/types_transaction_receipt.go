package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

// CommonTransactionReceipt Common properties for a transaction receipt
type CommonTransactionReceipt struct {
	// TransactionHash The hash identifying the transaction
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee       *felt.Felt         `json:"actual_fee"`
	ExecutionStatus TxnExecutionStatus `json:"execution_status"`
	FinalityStatus  TxnFinalityStatus  `json:"finality_status"`
	BlockHash       *felt.Felt         `json:"block_hash"`
	BlockNumber     uint64             `json:"block_number"`
	Type            TransactionType    `json:"type,omitempty"`
	MessagesSent    []MsgToL1          `json:"messages_sent"`
	RevertReason    string             `json:"revert_reason,omitempty"`
	// Events The events emitted as part of this transaction
	Events []Event `json:"events"`
}

func (tr CommonTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr CommonTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

// TODO: check how we can move that type up in starknet.go/types
type TransactionType string

const (
	TransactionType_Declare       TransactionType = "DECLARE"
	TransactionType_DeployAccount TransactionType = "DEPLOY_ACCOUNT"
	TransactionType_Deploy        TransactionType = "DEPLOY"
	TransactionType_Invoke        TransactionType = "INVOKE"
	TransactionType_L1Handler     TransactionType = "L1_HANDLER"
)

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

func (tt TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(tt))), nil
}

// InvokeTransactionReceipt Invoke Transaction Receipt
type InvokeTransactionReceipt CommonTransactionReceipt

func (tr InvokeTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr InvokeTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

// DeclareTransactionReceipt Declare Transaction Receipt
type DeclareTransactionReceipt CommonTransactionReceipt

func (tr DeclareTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr DeclareTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

// DeployTransactionReceipt Deploy  Transaction Receipt
type DeployTransactionReceipt struct {
	CommonTransactionReceipt
	// The address of the deployed contract
	ContractAddress *felt.Felt `json:"contract_address"`
}

func (tr DeployTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr DeployTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

// DeployAccountTransactionReceipt Deploy Account Transaction Receipt
type DeployAccountTransactionReceipt struct {
	CommonTransactionReceipt
	// ContractAddress The address of the deployed contract
	ContractAddress *felt.Felt `json:"contract_address"`
}

func (tr DeployAccountTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr DeployAccountTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

// L1HandlerTransactionReceipt L1 Handler Transaction Receipt
type L1HandlerTransactionReceipt CommonTransactionReceipt

func (tr L1HandlerTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr L1HandlerTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

type PendingL1HandlerTransactionReceipt struct {
	Type TransactionType `json:"type"`
	// The message hash as it appears on the L1 core contract
	MsgHash NumAsHex `json:"message_hash"`
	PendingCommonTransactionReceiptProperties
}

func (tr PendingL1HandlerTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr PendingL1HandlerTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

type PendingDeclareTransactionReceipt struct {
	Type TransactionType `json:"type"`
	PendingCommonTransactionReceiptProperties
}

func (tr PendingDeclareTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr PendingDeclareTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

type PendingDeployAccountTransactionReceipt struct {
	Type TransactionType `json:"type"`
	// The address of the deployed contract
	ContractAddress *felt.Felt `json:"contract_address"`
	PendingCommonTransactionReceiptProperties
}

func (tr PendingDeployAccountTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr PendingDeployAccountTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

type PendingInvokeTransactionReceipt struct {
	Type TransactionType `json:"type"`
	PendingCommonTransactionReceiptProperties
}

func (tr PendingInvokeTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr PendingInvokeTransactionReceipt) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

type PendingCommonTransactionReceiptProperties struct {
	// TransactionHash The hash identifying the transaction
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee       *felt.Felt         `json:"actual_fee"`
	Type            TransactionType    `json:"type,omitempty"`
	MessagesSent    []MsgToL1          `json:"messages_sent"`
	ExecutionStatus TxnExecutionStatus `json:"execution_status"`
	FinalityStatus  TxnFinalityStatus  `json:"finality_status"`
	RevertReason    string             `json:"revert_reason"`
	// Events The events emitted as part of this transaction
	Events             []Event            `json:"events"`
	ExecutionResources ExecutionResources `json:"execution_resources"`
}

type ExecutionResources struct {
	// The number of Cairo steps used
	Steps NumAsHex `json:"steps"`
	// The number of unused memory cells (each cell is roughly equivalent to a step)
	MemoryHoles NumAsHex `json:"memory_holes,omitempty"`
	// The number of RANGE_CHECK builtin instances
	RangeCheckApps NumAsHex `json:"range_check_builtin_applications"`
	// The number of Pedersen builtin instances
	PedersenApps NumAsHex `json:"pedersen_builtin_applications"`
	// The number of Poseidon builtin instances
	PoseidonApps NumAsHex `json:"poseidon_builtin_applications"`
	// The number of EC_OP builtin instances
	ECOPApps NumAsHex `json:"ec_op_builtin_applications"`
	// The number of ECDSA builtin instances
	ECDSAApps NumAsHex `json:"ecdsa_builtin_applications"`
	// The number of BITWISE builtin instances
	BitwiseApps NumAsHex `json:"bitwise_builtin_applications"`
	// The number of KECCAK builtin instances
	KeccakApps NumAsHex `json:"keccak_builtin_applications"`
}

func (tr PendingCommonTransactionReceiptProperties) Hash() *felt.Felt {
	return tr.TransactionHash
}

func (tr PendingCommonTransactionReceiptProperties) GetExecutionStatus() TxnExecutionStatus {
	return tr.ExecutionStatus
}

type TransactionReceipt interface {
	Hash() *felt.Felt
	GetExecutionStatus() TxnExecutionStatus
}

type OrderedMsg struct {
	// The order of the message within the transaction
	Order   int `json:"order"`
	MsgToL1 MsgToL1
}

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

type UnknownTransactionReceipt struct{ TransactionReceipt }

func (tr *UnknownTransactionReceipt) UnmarshalJSON(data []byte) error {
	var dec map[string]interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	t, err := unmarshalTransactionReceipt(dec)
	if err != nil {
		return err
	}
	*tr = UnknownTransactionReceipt{t}
	return nil
}

func unmarshalTransactionReceipt(t interface{}) (TransactionReceipt, error) {
	switch casted := t.(type) {
	case map[string]interface{}:
		// NOTE(tvanas): Pathfinder 0.3.3 does not return
		// transaction receipt types. We handle this by
		// naively marshalling into an invoke type. Once it
		// is supported, this condition can be removed.
		typ, ok := casted["type"]
		if !ok {
			return nil, fmt.Errorf("unknown transaction type: %v", t)
		}

		// Pending doesn't have a block number
		if casted["block_hash"] == nil {
			switch TransactionType(typ.(string)) {
			case TransactionType_Invoke:
				var txn PendingInvokeTransactionReceipt
				remarshal(casted, &txn)
				return txn, nil
			case TransactionType_DeployAccount:
				var txn PendingDeployAccountTransactionReceipt
				remarshal(casted, &txn)
				return txn, nil
			case TransactionType_L1Handler:
				var txn PendingL1HandlerTransactionReceipt
				remarshal(casted, &txn)
				return txn, nil
			case TransactionType_Declare:
				var txn PendingDeclareTransactionReceipt
				remarshal(casted, &txn)
				return txn, nil
			}
		}

		switch TransactionType(typ.(string)) {
		case TransactionType_Invoke:
			var txn InvokeTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_L1Handler:
			var txn L1HandlerTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_Declare:
			var txn DeclareTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_Deploy:
			var txn DeployTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_DeployAccount:
			var txn DeployAccountTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		}
	}

	return nil, fmt.Errorf("unknown transaction type: %v", t)
}

// The finality status of the transaction, including the case the txn is still in the mempool or failed validation during the block construction phase
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
