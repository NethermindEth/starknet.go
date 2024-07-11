package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

const (
	CONTRACT_ADDRESS = "0x059a943ca214c10234b9a3b61c558ac20c005127d183b86a99a8f3c60a08b4ff" 
)

var (
	FROM_BLOCK uint64 = 657100
	TO_BLOCK   uint64 = 657101
)

// main entry point of the program.
//
// It initializes the environment and establishes a connection with the client.
// It then retrieves events from the contract and prints how many it found with
// a series of more selective filters (all event types, just 2 event types, just
// those 2 event types but with a specified key value).

// Parameters:
//
//	none
//
// Returns:
//
//	none
func main() {
	fmt.Println("Starting readEvents example")

	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Initialize connection to RPC provider
	provider, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the RPC provider")

	// contractAddress is the address of the token contract whose events we want to read
	// it is read from the .env file

	contractAddress, err := utils.HexToFelt(CONTRACT_ADDRESS)
	if err != nil {
		msg := fmt.Errorf("failed to transform the token contract address %s, error %w", CONTRACT_ADDRESS, err)
		panic(msg)
	}

	// fromBlock and toBlock are the block numbers between which we want to read events
	fromBlock := &FROM_BLOCK
	toBlock := &TO_BLOCK

	// create an EventFilter that specifies which events we want
	eventFilter := rpc.EventFilter{
		FromBlock: rpc.BlockID{Number: fromBlock},
		ToBlock:   rpc.BlockID{Number: toBlock},
		Address:   contractAddress,
	}
	// ChunkSize is the maximum number of events to read in one call.
	// The readEvents function will make multiple calls to the provider
	// when the number of matching events is larger than ChunkSize
	resPageReq := rpc.ResultPageRequest{
		ChunkSize: 1000,
	}
	eventsInput := rpc.EventsInput{
		EventFilter:       eventFilter,
		ResultPageRequest: resPageReq,
	}

	var events []rpc.EmittedEvent

	// read all the events emitted by the contract in this range of blocks
	events = readEvents(eventsInput, provider)
	fmt.Printf("num events from contract %d\n", len(events))

	// narrow the scope to event types InterestStateSetBorrowingRate and InterestStateUpdated
	eventTypes := []string{"InterestStateSetBorrowingRate", "InterestStateUpdated"}
	// eventData is an empty slice, so we are not filtering by any key data, just event type
	eventData := [][]*felt.Felt{}
	keyFilter := buildKeyFilter(eventTypes, eventData)
	// set the filter on the Keys field of the eventsInput
	eventsInput.Keys = keyFilter

	// read the events with the new filter
	events = readEvents(eventsInput, provider)
	fmt.Printf("num events of specified types %d\n", len(events))

	// narrow the scope to events with particular key values
	data, _ := new(felt.Felt).SetString("0x1258eae3eae5002125bebf062d611a772e8aea3a1879b64a19f363ebd00947")
	keyFilter = buildKeyFilter(eventTypes, [][]*felt.Felt{{data}})
	// set the new filter on the Keys field of the eventsInput
	eventsInput.Keys = keyFilter

	// read the events with this filter
	events = readEvents(eventsInput, provider)
	fmt.Printf("num events of specified types with given keys %d\n", len(events))
}

func readEvents(eventsInput rpc.EventsInput, provider *rpc.Provider) []rpc.EmittedEvent {

	ctx := context.Background()
	events := make([]rpc.EmittedEvent, 0)
	haveMoreToRead := true
	for haveMoreToRead {
		eventChunk, err := provider.Events(ctx, eventsInput)
		if err != nil {
			msg := fmt.Errorf("error retrieving events: %w", err)
			panic(msg)
		}
		events = append(events, eventChunk.Events...)
		if eventChunk.ContinuationToken == "" {
			haveMoreToRead = false
		} else {
			eventsInput.ContinuationToken = eventChunk.ContinuationToken
		}
	}

	return events

}

// function to build a filter for particular event types and/or key values
// for a description of how the key filter values are interpreted see:
// https://community.starknet.io/t/snip-14-index-transfer-and-approval-events-in-erc20s/114212
func buildKeyFilter(eventTypes []string, eventData [][]*felt.Felt) [][]*felt.Felt {

	eventKeys := make([]*felt.Felt, 0, len(eventTypes))
	for _, eventName := range eventTypes {
		eventId := utils.GetSelectorFromName(eventName)
		eventKeys = append(eventKeys, utils.BigIntToFelt(eventId))
	}
	keyFilter := make([][]*felt.Felt, 0, 2)
	// first element of the key data is the event selector
	keyFilter = append(keyFilter, eventKeys)
	// this is followed by the key values
	keyFilter = append(keyFilter, eventData...)
	return keyFilter
}
