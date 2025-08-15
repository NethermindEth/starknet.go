package rpc

import "github.com/NethermindEth/juno/core/felt"

type OrderedEvent struct {
	// The order of the event within the transaction
	Order int `json:"order"`
	*EventContent
}

type Event struct {
	FromAddress *felt.Felt `json:"from_address"`
	EventContent
}

type EventContent struct {
	Keys []*felt.Felt `json:"keys"`
	Data []*felt.Felt `json:"data"`
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
	FromBlock BlockID `json:"from_block,omitempty"`
	// ToBlock to block
	ToBlock BlockID `json:"to_block,omitempty"`
	// Address from contract
	Address *felt.Felt `json:"address,omitempty"`
	// Keys the values used to filter the events
	Keys [][]*felt.Felt `json:"keys,omitempty"`
}

// EventsInput is the input for the 'starknet_getEvents' method.
// All fields are optional, except for the 'chunk_size' field.
type EventsInput struct {
	EventFilter
	ResultPageRequest
}

// EventSubscriptionInput is the input for the 'starknet_subscribeEvents' method.

type EventSubscriptionInput struct {
	// (Optional) Filter events by from_address which emitted the event
	FromAddress *felt.Felt `json:"from_address,omitempty"`
	// (Optional) Per key (by position), designate the possible values to be
	// matched for events to be returned. Empty array designates 'any' value
	Keys [][]*felt.Felt `json:"keys,omitempty"`
	// (Optional) The block to get notifications from, default is latest, limited to 1024 blocks back
	SubBlockID SubscriptionBlockID `json:"block_id,omitzero"`
	// (Optional) The finality status of the most recent events to include.
	// Only `PRE_CONFIRMED` and `ACCEPTED_ON_L2` are supported. Default is `ACCEPTED_ON_L2`.
	// If PRE_CONFIRMED finality is selected, events might appear multiple times, once for each finality status update.
	FinalityStatus TxnFinalityStatus `json:"finality_status,omitempty"`
}

// Notification from the server about a new event.
// The event also includes the finality status of the transaction emitting the event.
type EmittedEventWithFinalityStatus struct {
	EmittedEvent
	FinalityStatus TxnFinalityStatus `json:"finality_status"`
}
