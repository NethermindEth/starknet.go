package types

import ctypes "github.com/dontpanicdao/caigo/types"

type Event struct {
	FromAddress ctypes.Hash `json:"from_address"`
	Keys        []string    `json:"keys"`
	Data        []string    `json:"data"`
}

type EmittedEvent struct {
	Event
	BlockHash       ctypes.Hash `json:"block_hash"`
	BlockNumber     uint64      `json:"block_number"`
	TransactionHash ctypes.Hash `json:"transaction_hash"`
}

type EventFilter struct {
	FromBlock BlockID     `json:"from_block"`
	ToBlock   BlockID     `json:"to_block,omitempty"`
	Address   ctypes.Hash `json:"address,omitempty"`
	// Keys the values used to filter the events
	Keys []string `json:"keys,omitempty"`

	PageSize   uint64 `json:"page_size,omitempty"`
	PageNumber uint64 `json:"page_number"`
}

type EventsOutput struct {
	Events     []EmittedEvent `json:"events"`
	PageNumber uint64         `json:"page_number"`
	IsLastPage bool           `json:"is_last_page"`
}
