package types

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

type Transaction struct {
	TransactionReceipt
	TransactionHash    string   `json:"txn_hash,omitempty"`
	ContractAddress    string   `json:"contract_address,omitempty"`
	EntryPointSelector string   `json:"entry_point_selector,omitempty"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
	Nonce              string   `json:"nonce,omitempty"`
	Type               string   `json:"type,omitempty"`
}

type L1Message struct {
	ToAddress string  `json:"to_address,omitempty"`
	Payload   []*Felt `json:"payload,omitempty"`
}

type L2Message struct {
	FromAddress string  `json:"from_address,omitempty"`
	Payload     []*Felt `json:"payload,omitempty"`
}

type Event struct {
	FromAddress string  `json:"from_address,omitempty"`
	Keys        []*Felt `json:"keys,omitempty"`
	Values      []*Felt `json:"values,omitempty"`
}

type TransactionReceipt struct {
	TransactionHash string       `json:"txn_hash,omitempty"`
	Status          string       `json:"status,omitempty"`
	StatusData      string       `json:"status_data,omitempty"`
	MessagesSent    []*L1Message `json:"messages_sent,omitempty"`
	L1OriginMessage *L2Message   `json:"l1_origin_message,omitempty"`
	Events          []*Event     `json:"events,omitempty"`
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
