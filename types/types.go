package types

import (
	"math/big"
)

type Block struct {
	BlockHash       string         `json:"block_hash"`
	ParentBlockHash string         `json:"parent_hash"`
	BlockNumber     int            `json:"block_number"`
	NewRoot         string         `json:"new_root"`
	OldRoot         string         `json:"old_root"`
	Status          string         `json:"status"`
	AcceptedTime    uint64         `json:"accepted_time"`
	Transactions    []*Transaction `json:"transactions"`
}

type Code struct {
	Bytecode []string `json:"bytecode"`
	Abi      []struct {
		Inputs []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"inputs"`
		Name            string        `json:"name"`
		Outputs         []interface{} `json:"outputs"`
		Type            string        `json:"type"`
		StateMutability string        `json:"stateMutability,omitempty"`
	} `json:"abi"`
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
	BlockHash       string `json:"block_hash"`
	TxFailureReason struct {
		ErrorMessage string `json:"error_message,omitempty"`
	} `json:"tx_failure_reason,omitempty"`
}

type ABI struct {
	Members []struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"`
		Type   string `json:"type"`
	} `json:"members,omitempty"`
	Name   string `json:"name"`
	Size   int    `json:"size,omitempty"`
	Type   string `json:"type"`
	Inputs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"inputs,omitempty"`
	Outputs         []interface{} `json:"outputs,omitempty"`
	StateMutability string        `json:"stateMutability,omitempty"`
}

type AddTxResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
}

type DeployRequest struct {
	Type                string        `json:"type"`
	ContractAddressSalt string        `json:"contract_address_salt"`
	ConstructorCalldata []string      `json:"constructor_calldata"`
	ContractDefinition  ContractClass `json:"contract_definition"`
}

type ContractClass struct {
	ABI               []ABI             `json:"abi"`
	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`
	Program           string            `json:"program"`
}

type DeclareRequest struct {
	Type          string   `json:"type"`
	SenderAddress string   `json:"sender_address"`
	MaxFee        string   `json:"max_fee"`
	Nonce         string   `json:"nonce"`
	Signature     []string `json:"signature"`
	ContractClass struct {
		ABI               []ABI             `json:"abi"`
		EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`
		Program           string            `json:"program"`
	} `json:"contract_class"`
}

type EntryPointsByType struct {
	Constructor []EntryPointList `json:"CONSTRUCTOR"`
	External    []EntryPointList `json:"EXTERNAL"`
	L1Handler   []EntryPointList `json:"L1_HANDLER"`
}

type EntryPointList struct {
	Offset   string `json:"offset"`
	Selector string `json:"selector"`
}

type FunctionCall struct {
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
}

type FeeEstimate struct {
	Amount *big.Int `json:"amount"`
	Unit   string   `json:"unit"`
}

type ContractAddresses struct {
	Starknet             string `json:"Starknet"`
	GpsStatementVerifier string `json:"GpsStatementVerifier"`
}
