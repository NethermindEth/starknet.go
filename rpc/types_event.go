package rpc

import "github.com/NethermindEth/juno/core/felt"

type OrderedEvent struct {
	// The order of the event within the transaction
	Order int `json:"order"`
	*Event
}

type Event struct {
	FromAddress *felt.Felt   `json:"from_address"`
	Keys        []*felt.Felt `json:"keys"`
	Data        []*felt.Felt `json:"data"`
}

type EventChunk struct {
	Events            []EmittedEvent `json:"events"`
	ContinuationToken string         `json:"continuation_token,omitempty"`
}

// EmittedEvent an event emitted as a result of transaction execution
type EmittedEvent struct {
	Event
	// BlockHash the hash of the block in which the event was emitted
	BlockHash *felt.Felt `json:"block_hash,omitempty"`
	// BlockNumber the number of the block in which the event was emitted
	BlockNumber uint64 `json:"block_number,omitempty"`
	// TransactionHash the transaction that emitted the event
	TransactionHash *felt.Felt `json:"transaction_hash"`
}

type EventFilter struct {
	// FromBlock from block
	FromBlock BlockID `json:"from_block"`
	// ToBlock to block
	ToBlock BlockID `json:"to_block,omitempty"`
	// Address from contract
	Address *felt.Felt `json:"address,omitempty"`
	// Keys the values used to filter the events
	Keys [][]*felt.Felt `json:"keys,omitempty"`
}

type EventsInput struct {
	EventFilter
	ResultPageRequest
}

type EventSubscriptionInput struct {
	FromAddress *felt.Felt     `json:"from_address,omitempty"` // Optional. Filter events by from_address which emitted the event
	Keys        [][]*felt.Felt `json:"keys,omitempty"`         // Optional. Per key (by position), designate the possible values to be matched for events to be returned. Empty array designates 'any' value
	BlockID     BlockID        `json:"block_id,omitempty"`     // Optional. The block to get notifications from, default is latest, limited to 1024 blocks back
}
