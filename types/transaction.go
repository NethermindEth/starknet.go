package types

import "github.com/dontpanicdao/caigo/felt"

type StarknetTransaction struct {
	TransactionIndex int         `json:"transaction_index"`
	BlockNumber      uint64      `json:"block_number"`
	Transaction      Transaction `json:"transaction"`
	BlockHash        felt.Felt   `json:"block_hash"`
	Status           string      `json:"status"`
}

type Transaction struct {
	TransactionReceipt *TransactionReceipt
	TransactionHash    felt.Felt   `json:"txn_hash,omitempty"`
	ClassHash          felt.Felt   `json:"class_hash,omitempty"`
	ContractAddress    felt.Felt   `json:"contract_address,omitempty"`
	SenderAddress      felt.Felt   `json:"sender_address,omitempty"`
	EntryPointSelector *felt.Felt  `json:"entry_point_selector,omitempty"`
	Calldata           []felt.Felt `json:"calldata"`
	Signature          []felt.Felt `json:"signature"`
	MaxFee             *felt.Felt  `json:"max_fee,omitempty"`
	Nonce              *felt.Felt  `json:"nonce,omitempty"`
	Version            string      `json:"version,omitempty"`
	Type               string      `json:"type,omitempty"`
}

type L1Message struct {
	ToAddress felt.Felt   `json:"to_address,omitempty"`
	Payload   []felt.Felt `json:"payload,omitempty"`
}

type L2Message struct {
	FromAddress felt.Felt   `json:"from_address,omitempty"`
	Payload     []felt.Felt `json:"payload,omitempty"`
}

type Event struct {
	Order       int         `json:"order,omitempty"`
	FromAddress felt.Felt   `json:"from_address,omitempty"`
	Keys        []felt.Felt `json:"keys,omitempty"`
	Data        []felt.Felt `json:"data,omitempty"`
}

type TransactionReceipt struct {
	TransactionHash felt.Felt   `json:"txn_hash,omitempty"`
	Status          string      `json:"status,omitempty"`
	StatusData      string      `json:"status_data,omitempty"`
	MessagesSent    []L1Message `json:"messages_sent,omitempty"`
	L1OriginMessage L2Message   `json:"l1_origin_message,omitempty"`
	Events          []Event     `json:"events,omitempty"`
}

type ExecutionResources struct {
	NSteps                 int `json:"n_steps"`
	BuiltinInstanceCounter struct {
		PedersenBuiltin   int `json:"pedersen_builtin"`
		RangeCheckBuiltin int `json:"range_check_builtin"`
		BitwiseBuiltin    int `json:"bitwise_builtin"`
		OutputBuiltin     int `json:"output_builtin"`
		EcdsaBuiltin      int `json:"ecdsa_builtin"`
		EcOpBuiltin       int `json:"ec_op_builtin"`
	} `json:"builtin_instance_counter"`
	NMemoryHoles int `json:"n_memory_holes"`
}

type TransactionTrace struct {
	FunctionInvocation FunctionInvocation `json:"function_invocation"`
	Signature          []felt.Felt        `json:"signature"`
}

type FunctionInvocation struct {
	CallerAddress      felt.Felt            `json:"caller_address"`
	ContractAddress    felt.Felt            `json:"contract_address"`
	Calldata           []felt.Felt          `json:"calldata"`
	CallType           string               `json:"call_type"`
	ClassHash          felt.Felt            `json:"class_hash"`
	Selector           felt.Felt            `json:"selector"`
	EntryPointType     felt.Felt            `json:"entry_point_type"`
	Result             []felt.Felt          `json:"result"`
	ExecutionResources *ExecutionResources  `json:"execution_resources"`
	InternalCalls      []FunctionInvocation `json:"internal_calls"`
	Events             []Event              `json:"events"`
	Messages           []interface{}        `json:"messages"`
}
