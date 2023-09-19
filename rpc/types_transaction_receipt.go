package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

// CommonTransactionReceipt Common properties for a transaction receipt
type CommonTransactionReceipt struct {
	// TransactionHash The hash identifying the transaction
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee    *felt.Felt       `json:"actual_fee"`
	Status       TransactionState `json:"status"`
	BlockHash    *felt.Felt       `json:"block_hash"`
	BlockNumber  uint64           `json:"block_number"`
	Type         TransactionType  `json:"type"`
	MessagesSent []MsgToL1        `json:"messages_sent"`
	// Events The events emitted as part of this transaction
	Events []Event `json:"events"`
}

func (tr CommonTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
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

// DeclareTransactionReceipt Declare Transaction Receipt
type DeclareTransactionReceipt CommonTransactionReceipt

func (tr DeclareTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
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

// DeployAccountTransactionReceipt Deploy Account Transaction Receipt
type DeployAccountTransactionReceipt struct {
	CommonTransactionReceipt
	// ContractAddress The address of the deployed contract
	ContractAddress *felt.Felt `json:"contract_address"`
}

func (tr DeployAccountTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

// L1HandlerTransactionReceipt L1 Handler Transaction Receipt
type L1HandlerTransactionReceipt CommonTransactionReceipt

func (tr L1HandlerTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

type PendingDeployTransactionReceipt struct {
	CommonTransactionReceipt
	// The address of the deployed contract
	ContractAddress *felt.Felt `json:"contract_address"`
}

func (tr PendingDeployTransactionReceipt) Hash() *felt.Felt {
	return tr.TransactionHash
}

type PendingCommonTransactionReceiptProperties struct {
	// TransactionHash The hash identifying the transaction
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee    *felt.Felt      `json:"actual_fee"`
	Type         TransactionType `json:"type"`
	MessagesSent []MsgToL1       `json:"messages_sent"`
	// Events The events emitted as part of this transaction
	Events []Event `json:"events"`
}

func (tr PendingCommonTransactionReceiptProperties) Hash() *felt.Felt {
	return tr.TransactionHash
}

type TransactionReceipt interface {
	Hash() *felt.Felt
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
	case string:
		var txn InvokeTransactionReceipt
		txhash, err := utils.HexToFelt(casted)
		if err != nil {
			return txn, err
		}
		return TransactionHash{txhash}, nil
	case map[string]interface{}:
		// NOTE(tvanas): Pathfinder 0.3.3 does not return
		// transaction receipt types. We handle this by
		// naively marshalling into an invoke type. Once it
		// is supported, this condition can be removed.
		typ, ok := casted["type"]
		if !ok {
			var txn InvokeTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		}

		// Pending doesn't have a block number
		if casted["block_hash"] == nil {
			switch TransactionType(typ.(string)) {
			case TransactionType_Deploy:
				var txn PendingDeployTransactionReceipt
				remarshal(casted, &txn)
				return txn, nil
			default:
				var txn PendingCommonTransactionReceiptProperties
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
