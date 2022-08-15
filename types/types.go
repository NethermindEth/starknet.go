package types

import "github.com/dontpanicdao/caigo/felt"

type Bytecode []string

type Block struct {
	BlockHash       felt.Felt     `json:"block_hash"`
	ParentBlockHash *felt.Felt    `json:"parent_hash"`
	BlockNumber     uint64        `json:"block_number"`
	NewRoot         string        `json:"new_root"`
	OldRoot         string        `json:"old_root"`
	Status          string        `json:"status"`
	AcceptedTime    uint64        `json:"accepted_time"`
	GasPrice        *felt.Felt    `json:"gas_price"`
	Transactions    []Transaction `json:"transactions"`
}

type Code struct {
	Bytecode Bytecode `json:"bytecode"`
	Abi      []ABI    `json:"abi"`
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
	TxStatus        string     `json:"tx_status"`
	BlockHash       *felt.Felt `json:"block_hash,omitempty"`
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
	Code            string    `json:"code"`
	TransactionHash felt.Felt `json:"transaction_hash"`
}

type DeployRequest struct {
	Type                string        `json:"type"`
	ContractAddressSalt *felt.Felt    `json:"contract_address_salt"`
	ConstructorCalldata []felt.Felt   `json:"constructor_calldata"`
	ContractDefinition  ContractClass `json:"contract_definition"`
}

type ContractClass struct {
	ABI               []ABI             `json:"abi"`
	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`
	Program           interface{}       `json:"program"`
}

type DeclareRequest struct {
	Type          string         `json:"type"`
	SenderAddress felt.Felt      `json:"sender_address"`
	MaxFee        *felt.Felt     `json:"max_fee"`
	Nonce         *felt.Felt     `json:"nonce"`
	Signature     felt.Signature `json:"signature"`
	ContractClass ContractClass  `json:"contract_class"`
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
	ContractAddress    felt.Felt   `json:"contract_address"`
	EntryPointSelector *felt.Felt  `json:"entry_point_selector"`
	Calldata           []felt.Felt `json:"calldata"`
}

type FunctionInvoke struct {
	FunctionCall
	MaxFee          *felt.Felt     `json:"max_fee,omitempty"`
	Nonce           *felt.Felt     `json:"nonce,omitempty"`
	Version         *felt.Felt     `json:"version,omitempty"`
	Signature       felt.Signature `json:"signature,omitempty"`
	TransactionHash felt.Felt      `json:"txn_hash,omitempty"`
}

// FeeEstimate provides a set of properties to understand fee estimations.
type FeeEstimate struct {
	OverallFee *felt.Felt `json:"overall_fee,omitempty"`
	GasUsage   *felt.Felt `json:"gas_usage,omitempty"`
	GasPrice   *felt.Felt `json:"gas_price,omitempty"`
	Unit       string     `json:"unit,omitempty"`
}

type ContractAddresses struct {
	Starknet             string `json:"Starknet"`
	GpsStatementVerifier string `json:"GpsStatementVerifier"`
}
