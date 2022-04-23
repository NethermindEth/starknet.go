package caigo

type Block struct {
	BlockHash           string               `json:"block_hash"`
	ParentBlockHash     string               `json:"parent_block_hash"`
	BlockNumber         int                  `json:"block_number"`
	StateRoot           string               `json:"state_root"`
	Status              string               `json:"status"`
	Transactions        []Transaction        `json:"transactions"`
	Timestamp           int                  `json:"timestamp"`
	TransactionReceipts []TransactionReceipt `json:"transaction_receipts"`
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

// Starknet transaction composition
type Transaction struct {
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
	EntryPointType     string   `json:"entry_point_type,omitempty"`
	TransactionHash    string   `json:"transaction_hash,omitempty"`
	Type               string   `json:"type,omitempty"`
	Nonce              string   `json:"nonce,omitempty"`
}

type StarknetTransaction struct {
	TransactionIndex int         `json:"transaction_index"`
	BlockNumber      int         `json:"block_number"`
	Transaction      Transaction `json:"transaction"`
	BlockHash        string      `json:"block_hash"`
	Status           string      `json:"status"`
}

type TransactionReceipt struct {
	Status                string `json:"status"`
	BlockHash             string `json:"block_hash"`
	BlockNumber           int    `json:"block_number"`
	TransactionIndex      int    `json:"transaction_index"`
	TransactionHash       string `json:"transaction_hash"`
	L1ToL2ConsumedMessage struct {
		FromAddress string   `json:"from_address"`
		ToAddress   string   `json:"to_address"`
		Selector    string   `json:"selector"`
		Payload     []string `json:"payload"`
	} `json:"l1_to_l2_consumed_message"`
	L2ToL1Messages     []interface{} `json:"l2_to_l1_messages"`
	Events             []interface{} `json:"events"`
	ExecutionResources struct {
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
	} `json:"execution_resources"`
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

type RawContractDefinition struct {
	ABI               []ABI                  `json:"abi"`
	EntryPointsByType EntryPointsByType      `json:"entry_points_by_type"`
	Program           map[string]interface{} `json:"program"`
}

type DeployRequest struct {
	Type                string   `json:"type"`
	ContractAddressSalt string   `json:"contract_address_salt"`
	ConstructorCalldata []string `json:"constructor_calldata"`
	ContractDefinition  struct {
		ABI               []ABI             `json:"abi"`
		EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`
		Program           string            `json:"program"`
	} `json:"contract_definition"`
}

type EntryPointsByType struct {
	Constructor []struct {
		Offset   string `json:"offset"`
		Selector string `json:"selector"`
	} `json:"CONSTRUCTOR"`
	External []struct {
		Offset   string `json:"offset"`
		Selector string `json:"selector"`
	} `json:"EXTERNAL"`
	L1Handler []interface{} `json:"L1_HANDLER"`
}
