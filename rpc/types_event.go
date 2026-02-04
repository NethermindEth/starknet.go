package rpc

import (
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

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
	// The index of the event in the transaction by which it was emitted
	EventIndex uint64 `json:"event_index"`
	// TransactionHash the transaction that emitted the event
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// The index of the transaction in the block by which the event was emitted
	TransactionIndex uint64 `json:"transaction_index"`
}

type EventFilter struct {
	// FromBlock from block
	FromBlock BlockID `json:"from_block,omitempty"`
	// ToBlock to block
	ToBlock BlockID `json:"to_block,omitempty"`
	// A contract address or a list of addresses from which events should originate"
	Address AddressList `json:"address,omitempty"`
	// Keys the values used to filter the events
	Keys [][]*felt.Felt `json:"keys,omitempty"`
}

// AddressList is a list of addresses from which events should originate.
// If the list contains a single address, it will be marshalled
// as a single string in the JSON response, otherwise as an array.
type AddressList []*felt.Felt

// MarshalJSON marshals the AddressList into JSON.
func (al AddressList) MarshalJSON() ([]byte, error) {
	if len(al) == 1 {
		return json.Marshal(al[0])
	}

	return json.Marshal([]*felt.Felt(al))
}

// UnmarshalJSON unmarshals the JSON data into an AddressList.
func (al *AddressList) UnmarshalJSON(data []byte) error {
	var arr []*felt.Felt
	err := json.Unmarshal(data, &arr)
	if err != nil {
		var singleF *felt.Felt
		err2 := json.Unmarshal(data, &singleF)
		if err2 != nil {
			return errors.Join(errors.New("failed to unmarshal address list"), err, err2)
		}
		*al = AddressList{singleF}
	} else {
		*al = arr
	}

	return nil
}

// EventsInput is the input for the 'starknet_getEvents' method.
// All fields are optional, except for the 'chunk_size' field.
type EventsInput struct {
	EventFilter
	ResultPageRequest
}

// EventSubscriptionInput is the input for the 'starknet_subscribeEvents' method.

type EventSubscriptionInput struct {
	// (Optional) A contract address or a list of addresses from which events
	// should originate
	FromAddress AddressList `json:"from_address,omitempty"`
	// (Optional) Per key (by position), designate the possible values to be
	// matched for events to be returned. Empty array designates 'any' value
	Keys [][]*felt.Felt `json:"keys,omitempty"`
	// (Optional) The block to get notifications from, default is latest, limited
	// to 1024 blocks back
	SubBlockID SubscriptionBlockID `json:"block_id,omitzero"`
	// (Optional) The finality status of the most recent events to include.
	// Only `PRE_CONFIRMED` and `ACCEPTED_ON_L2` are supported. Default is `ACCEPTED_ON_L2`.
	// If PRE_CONFIRMED finality is selected, events might appear multiple times,
	// once for each finality status update.
	FinalityStatus TxnFinalityStatus `json:"finality_status,omitempty"`
}

// Notification from the server about a new event.
// The event also includes the finality status of the transaction emitting the
// event.
type EmittedEventWithFinalityStatus struct {
	EmittedEvent
	FinalityStatus TxnFinalityStatus `json:"finality_status"`
}
