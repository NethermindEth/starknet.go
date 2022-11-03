package rpcv02

import "github.com/dontpanicdao/caigo/types"

type Event struct {
	FromAddress types.Hash `json:"from_address"`
	Keys        []string   `json:"keys"`
	Data        []string   `json:"data"`
}

// EmittedEvent an event emitted as a result of transaction execution
type EmittedEvent struct {
	Event
	// BlockHash the hash of the block in which the event was emitted
	BlockHash types.Hash `json:"block_hash"`
	// BlockNumber the number of the block in which the event was emitted
	BlockNumber uint64 `json:"block_number"`
	// TransactionHash the transaction that emitted the event
	TransactionHash types.Hash `json:"transaction_hash"`
}

type EventFilter struct {
	// FromBlock from block
	FromBlock BlockID `json:"from_block"`
	// ToBlock to block
	ToBlock BlockID `json:"to_block,omitempty"`
	// Address from contract
	Address types.Hash `json:"address,omitempty"`
	// Keys the values used to filter the events
	Keys []string `json:"keys,omitempty"`
}

type EventsOutput struct {
	Events     []EmittedEvent `json:"events"`
	PageNumber uint64         `json:"page_number"`
	IsLastPage bool           `json:"is_last_page"`
}
