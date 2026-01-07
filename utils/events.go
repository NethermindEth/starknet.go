package utils

import (
	"github.com/NethermindEth/starknet.go/rpc"
)

// EventWith filters events by event key/selector name.
// Returns the first matching event or nil if not found.
//
// The function converts the eventKey string to a selector using GetSelectorFromNameFelt,
// then searches through the events slice for an event whose first key matches the selector.
//
// Parameters:
//   - events: A slice of rpc.Event to search through
//   - eventKey: The event name to search for (e.g., "Transfer", "Approval")
//
// Returns:
//   - *rpc.Event: A pointer to the first matching event, or nil if not found
//
// Example:
//
//	events := []rpc.Event{...}
//	transferEvent := utils.EventWith(events, "Transfer")
//	if transferEvent != nil {
//	    fmt.Println("Found transfer event from:", transferEvent.FromAddress)
//	}
func EventWith(events []rpc.Event, eventKey string) *rpc.Event {
	if len(events) == 0 {
		return nil
	}

	selector := GetSelectorFromNameFelt(eventKey)
	selectorStr := selector.String()

	for i := range events {
		if len(events[i].Keys) > 0 && events[i].Keys[0].String() == selectorStr {
			return &events[i]
		}
	}

	return nil
}

// EventsWithKey filters and returns all events matching the given key/selector name.
//
// The function converts the eventKey string to a selector using GetSelectorFromNameFelt,
// then collects all events whose first key matches the selector.
//
// Parameters:
//   - events: A slice of rpc.Event to search through
//   - eventKey: The event name to search for (e.g., "Transfer", "Approval")
//
// Returns:
//   - []rpc.Event: A slice of all matching events (empty slice if none found)
//
// Example:
//
//	events := []rpc.Event{...}
//	allTransfers := utils.EventsWithKey(events, "Transfer")
//	fmt.Printf("Found %d transfer events\n", len(allTransfers))
func EventsWithKey(events []rpc.Event, eventKey string) []rpc.Event {
	if len(events) == 0 {
		return []rpc.Event{}
	}

	selector := GetSelectorFromNameFelt(eventKey)
	selectorStr := selector.String()

	var result []rpc.Event
	for _, event := range events {
		if len(event.Keys) > 0 && event.Keys[0].String() == selectorStr {
			result = append(result, event)
		}
	}

	if result == nil {
		return []rpc.Event{}
	}

	return result
}
