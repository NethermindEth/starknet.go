package types

type Event struct {
	FromAddress Hash     `json:"from_address"`
	Keys        []string `json:"keys"`
	Data        []string `json:"data"`
}

type EmittedEvent struct {
	Event
	BlockHash       Hash   `json:"block_hash"`
	BlockNumber     uint64 `json:"block_number"`
	TransactionHash Hash   `json:"transaction_hash"`
}

type EventFilter struct {
	FromBlock BlockID `json:"from_block"`
	ToBlock   BlockID `json:"to_block"`
	Address   Hash    `json:"address"`
	// Keys the values used to filter the events
	Keys []string `json:"keys"`

	// ContinuationToken a pointer to the last element of the delivered page, use this token in a subsequent query to obtain the next page
	ContinuationToken string `json:"continuation_token,omitempty"`
	ChunkSize         uint64 `json:"chunk_size,omitempty"`
}

type EventsOutput struct {
	Events            []EmittedEvent `json:"events"`
	ContinuationToken string         `json:"continuation_token,omitempty"`
}
