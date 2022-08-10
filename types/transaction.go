package types

type StarknetTransaction struct {
	TransactionIndex int         `json:"transaction_index"`
	BlockNumber      int         `json:"block_number"`
	Transaction      Transaction `json:"transaction"`
	BlockHash        *Felt       `json:"block_hash"`
	Status           string      `json:"status"`
}

type Transaction struct {
	TransactionReceipt
	TransactionHash    string  `json:"txn_hash,omitempty"`
	ClassHash          string  `json:"class_hash,omitempty"`
	ContractAddress    *Felt   `json:"contract_address,omitempty"`
	SenderAddress      *Felt   `json:"sender_address,omitempty"`
	EntryPointSelector string  `json:"entry_point_selector,omitempty"`
	Calldata           []*Felt `json:"calldata"`
	Signature          []*Felt `json:"signature"`
	MaxFee             *Felt   `json:"max_fee,omitempty"`
	Nonce              *Felt   `json:"nonce,omitempty"`
	Version            string  `json:"version,omitempty"`
	Type               string  `json:"type,omitempty"`
}

type L1Message struct {
	ToAddress *Felt   `json:"to_address,omitempty"`
	Payload   []*Felt `json:"payload,omitempty"`
}

type L2Message struct {
	FromAddress *Felt   `json:"from_address,omitempty"`
	Payload     []*Felt `json:"payload,omitempty"`
}

type Event struct {
	Order       int     `json:"order,omitempty"`
	FromAddress *Felt   `json:"from_address,omitempty"`
	Keys        []*Felt `json:"keys,omitempty"`
	Data        []*Felt `json:"data,omitempty"`
}

type TransactionReceipt struct {
	TransactionHash string       `json:"txn_hash,omitempty"`
	Status          string       `json:"status,omitempty"`
	StatusData      string       `json:"status_data,omitempty"`
	MessagesSent    []*L1Message `json:"messages_sent,omitempty"`
	L1OriginMessage *L2Message   `json:"l1_origin_message,omitempty"`
	Events          []*Event     `json:"events,omitempty"`
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
	Signature          []*Felt            `json:"signature"`
}

type FunctionInvocation struct {
	CallerAddress      *Felt                `json:"caller_address"`
	ContractAddress    *Felt                `json:"contract_address"`
	Calldata           []*Felt              `json:"calldata"`
	CallType           string               `json:"call_type"`
	ClassHash          string               `json:"class_hash"`
	Selector           *Felt                `json:"selector"`
	EntryPointType     string               `json:"entry_point_type"`
	Result             []string             `json:"result"`
	ExecutionResources ExecutionResources   `json:"execution_resources"`
	InternalCalls      []FunctionInvocation `json:"internal_calls"`
	Events             []Event              `json:"events"`
	Messages           []interface{}        `json:"messages"`
}
