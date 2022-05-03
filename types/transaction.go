package types

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
