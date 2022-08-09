package types

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
	TxStatus        string `json:"tx_status"`
	BlockHash       string `json:"block_hash,omitempty"`
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
	Program           interface{}       `json:"program"`
}

type DeclareRequest struct {
	Type          string        `json:"type"`
	SenderAddress string        `json:"sender_address"`
	MaxFee        string        `json:"max_fee"`
	Nonce         string        `json:"nonce"`
	Signature     []string      `json:"signature"`
	ContractClass ContractClass `json:"contract_class"`
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
}

type Signature []*Felt10

type FunctionInvoke struct {
	FunctionCall
	MaxFee          *Felt10   `json:"max_fee,omitempty"`
	Nonce           *Felt10   `json:"nonce,omitempty"`
	Version         uint64    `json:"version,omitempty"`
	Signature       Signature `json:"signature,omitempty"`
	TransactionHash *Felt10   `json:"txn_hash,omitempty"`
}

// FeeEstimate provides a set of properties to understand fee estimations.
type FeeEstimate struct {
	OverallFee uint64 `json:"overall_fee,omitempty"`
	GasUsage   uint64 `json:"gas_usage,omitempty"`
	GasPrice   uint64 `json:"gas_price,omitempty"`
	Unit       string `json:"unit,omitempty"`
}

type ContractAddresses struct {
	Starknet             string `json:"Starknet"`
	GpsStatementVerifier string `json:"GpsStatementVerifier"`
}
