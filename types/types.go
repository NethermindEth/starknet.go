package types

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type NumAsHex string

type Bytecode []string

type Block struct {
	BlockHash       string         `json:"block_hash"`
	ParentBlockHash string         `json:"parent_hash"`
	BlockNumber     int            `json:"block_number"`
	NewRoot         string         `json:"new_root"`
	OldRoot         string         `json:"old_root"`
	Status          string         `json:"status"`
	AcceptedTime    uint64         `json:"accepted_time"`
	GasPrice        string         `json:"gas_price"`
	Transactions    []*Transaction `json:"transactions"`
}

type Code struct {
	Bytecode Bytecode `json:"bytecode"`
	Abi      *ABI     `json:"abi"`
}

func (c *Code) UnmarshalJSON(content []byte) error {
	v := map[string]json.RawMessage{}
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	// process 'bytecode'.
	data, ok := v["bytecode"]
	if !ok {
		return fmt.Errorf("missing bytecode in json object")
	}
	bytecode := []string{}
	if err := json.Unmarshal(data, &bytecode); err != nil {
		return err
	}
	c.Bytecode = bytecode

	// process 'abi'
	data, ok = v["abi"]
	if !ok {
		// contractClass can have an empty ABI for instance with ClassAt
		return nil
	}

	abis := []interface{}{}
	if err := json.Unmarshal(data, &abis); err != nil {
		return err
	}

	abiPointer := ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab ABIEntry
			abiType, ok := checkABI["type"].(string)
			if !ok {
				return fmt.Errorf("unknown abi type %v", checkABI["type"])
			}
			switch abiType {
			case string(ABITypeConstructor), string(ABITypeFunction), string(ABITypeL1Handler):
				ab = &FunctionABIEntry{}
			case string(ABITypeStruct):
				ab = &StructABIEntry{}
			case string(ABITypeEvent):
				ab = &EventABIEntry{}
			default:
				return fmt.Errorf("unknown ABI type %v", checkABI["type"])
			}
			data, err := json.Marshal(checkABI)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, ab)
			if err != nil {
				return err
			}
			abiPointer = append(abiPointer, ab)
		}
	}

	c.Abi = &abiPointer
	return nil
}

/*
StarkNet transaction states
*/
const (
	NOT_RECIEVED = TxStatus(iota)
	REJECTED
	RECEIVED
	PENDING
	ACCEPTED_ON_L2
	ACCEPTED_ON_L1
)

var TxStatuses = []string{"NOT_RECEIVED", "REJECTED", "RECEIVED", "PENDING", "ACCEPTED_ON_L2", "ACCEPTED_ON_L1"}

type TxStatus int

func (s TxStatus) String() string {
	return TxStatuses[s]
}

type TransactionStatus struct {
	TxStatus        string `json:"tx_status"`
	BlockHash       string `json:"block_hash,omitempty"`
	TxFailureReason struct {
		ErrorMessage string `json:"error_message,omitempty"`
	} `json:"tx_failure_reason,omitempty"`
}

type AddInvokeTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
}

type AddDeclareResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
	ClassHash       string `json:"class_hash"`
}

type AddDeployResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
	ContractAddress string `json:"address"`
}

type DeployRequest struct {
	Type                string        `json:"type"`
	ContractAddressSalt string        `json:"contract_address_salt"`
	ConstructorCalldata []string      `json:"constructor_calldata"`
	ContractDefinition  ContractClass `json:"contract_definition"`
}

type DeclareRequest struct {
	Type          string        `json:"type"`
	SenderAddress string        `json:"sender_address"`
	MaxFee        string        `json:"max_fee"`
	Nonce         string        `json:"nonce"`
	Signature     []string      `json:"signature"`
	ContractClass ContractClass `json:"contract_class"`
}

type EntryPointList struct {
	Offset   string `json:"offset"`
	Selector string `json:"selector"`
}

// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    Hash   `json:"contract_address"`
	EntryPointSelector string `json:"entry_point_selector,omitempty"`

	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

type Signature []*big.Int

type FunctionInvoke struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version uint64 `json:"version"`
	// Signature
	Signature Signature `json:"signature"`
	// Nonce should only be set with Transaction V1
	Nonce *big.Int `json:"nonce,omitempty"`

	FunctionCall
}

type FeeEstimate struct {
	GasConsumed NumAsHex `json:"gas_consumed"`
	GasPrice    NumAsHex `json:"gas_price"`
	OverallFee  NumAsHex `json:"overall_fee"`
}

type ContractAddresses struct {
	Starknet             string `json:"Starknet"`
	GpsStatementVerifier string `json:"GpsStatementVerifier"`
}

// ExecuteDetails provides some details about the execution.
type ExecuteDetails struct {
	MaxFee *big.Int
	Nonce  *big.Int
}
