package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
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

	// Read provider URL from .env file
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Initialize connection to RPC provider
	provider, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %v", err))
	}
	fmt.Println("Established connection with the RPC provider")

	// TODO: the plan is to make one function for each example
	// and then call it with the provider as argument

	// **********
	// 1. call with ChunkSize and ContinuationToken
	// **********
	callWithChunkSizeAndContinuationToken(provider)
	// **********
	// 2. call with ChunkSize only
	// **********
	fmt.Println(" ----- 2. call with ChunkSize only -----")

	simpleExample(provider)

	moreComplexExample(provider)
}

func callWithChunkSizeAndContinuationToken(provider *rpc.Provider) {
	fmt.Println(" ----- 1. call with ChunkSize and ContinuationToken -----")

	// The only required field is the 'ChunkSize' field, so let's fill it. This field is used
	// to limit the number of events returned in one call. If the number of events is greater than
	// the ChunkSize, the provider will return a continuation token in the 'ContinuationToken' field
	// that can be used to retrieve the next chunk.
	//
	// This will return 1000 events starting from the block 0.
	eventChunk, err := provider.Events(context.Background(), rpc.EventsInput{
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}
	fmt.Printf("number of returned events in the first chunk: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the last event in the first chunk: %d\n", eventChunk.Events[len(eventChunk.Events)-1].BlockNumber)

	// Now we will get the second chunk
	secondEventChunk, err := provider.Events(context.Background(), rpc.EventsInput{
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize:         1000,
			ContinuationToken: eventChunk.ContinuationToken,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}
	fmt.Printf("number of returned events in the second chunk: %d\n", len(secondEventChunk.Events))
	fmt.Printf("block number of the last event in the second chunk: %d\n", secondEventChunk.Events[len(secondEventChunk.Events)-1].BlockNumber)
}

func simpleExample(provider *rpc.Provider) {
	const (
		CONTRACT_ADDRESS        = "0x04c1d9da136846ab084ae18cf6ce7a652df7793b666a16ce46b1bf5850cc739d"
		FROM_BLOCK       uint64 = 301886
		TO_BLOCK         uint64 = 301887
	)

	fmt.Println("Starting readEvents simple example")

	// create an EventFilter that specifies which events we want
	eventFilter := rpc.EventFilter{
		FromBlock: rpc.WithBlockNumber(FROM_BLOCK),
		ToBlock:   rpc.WithBlockNumber(TO_BLOCK),
	}

	// ChunkSize is the maximum number of events to read in one call.
	resPageReq := rpc.ResultPageRequest{
		ChunkSize: 1000,
	}

	// Create an EventsInput object with the event filter and result page request
	eventsInput := rpc.EventsInput{
		EventFilter:       eventFilter,
		ResultPageRequest: resPageReq,
	}

	// Read the events from the contract
	eventChunk, err := provider.Events(context.Background(), eventsInput)
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}
	events := eventChunk.Events
	fmt.Printf("number of events from contract: %d\n", len(events))

	// print out the events read
	fmt.Println("events read:")
	printEvents(events)
	fmt.Println()
}

func moreComplexExample(provider *rpc.Provider) {
	const (
		CONTRACT_ADDRESS        = "0x04c1d9da136846ab084ae18cf6ce7a652df7793b666a16ce46b1bf5850cc739d"
		FROM_BLOCK       uint64 = 16206
		TO_BLOCK         uint64 = 16208
	)

	fmt.Println("Starting readEvents more complex example")

	// contractAddress is the address of the contract whose events we want to read
	contractAddress, err := utils.HexToFelt(CONTRACT_ADDRESS)
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the contract address %s, error %v", CONTRACT_ADDRESS, err))
	}

	// create an EventFilter that specifies which events we want
	eventFilter := rpc.EventFilter{
		FromBlock: rpc.WithBlockNumber(FROM_BLOCK),
		ToBlock:   rpc.WithBlockNumber(TO_BLOCK),
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
	fmt.Printf("number of events from contract: %d\n", len(events))

	// narrow the scope to event types AccountCreated and TransactionExecuted
	eventTypes := []string{"AccountCreated", "TransactionExecuted"}
	// eventData is an empty slice, so we are not filtering by any key data, just event type
	eventData := [][]*felt.Felt{}
	keyFilter := buildKeyFilter(eventTypes, eventData)
	// set the filter on the Keys field of the eventsInput
	eventsInput.Keys = keyFilter

	// read the events with the new filter
	events = readEvents(eventsInput, provider)
	fmt.Printf("number of events of specified types: %d\n", len(events))

	// narrow the scope to events with a particular key value
	var data *felt.Felt
	data, err = new(felt.Felt).SetString(
		"0x2f1d2a0070a008fd312a2368776aca5b57c4a3cd734efdb619c616af7ab64f5",
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create felt, error: %v", err))
	}
	keyFilter = buildKeyFilter(eventTypes, [][]*felt.Felt{{data}})
	// set the new filter on the Keys field of the eventsInput
	eventsInput.Keys = keyFilter

	// read the events with this filter
	events = readEvents(eventsInput, provider)
	fmt.Printf("number of events of specified types with given key: %d\n", len(events))

	fmt.Println("events read:")
	printEvents(events)
}

func readEvents(eventsInput rpc.EventsInput, provider *rpc.Provider) []rpc.EmittedEvent {

	ctx := context.Background()
	events := []rpc.EmittedEvent{}
	haveMoreToRead := true
	for haveMoreToRead {
		eventChunk, err := provider.Events(ctx, eventsInput)
		if err != nil {
			panic(fmt.Sprintf("error retrieving events: %v", err))
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
//
// for a description of how the key filter values are interpreted see:
//
// https://community.starknet.io/t/snip-14-index-transfer-and-approval-events-in-erc20s/114212
//
// in particular this part:
//
// "if the user sent an event filter containing [[k_1, k_2], [], [k_3]]
// then the node should return events whose first key is k_1 or k_2
// and the third key is k_3, and the second key is unconstrained and can
// take any value"

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

func printEvents(events []rpc.EmittedEvent) {
	accountCreatedFelt := utils.BigIntToFelt(utils.GetSelectorFromName("AccountCreated"))
	transactionExecutedFelt := utils.BigIntToFelt(utils.GetSelectorFromName("TransactionExecuted"))
	for _, event := range events {
		if event.Keys[0].Cmp(accountCreatedFelt) == 0 {
			fmt.Printf("AccountCreated %s\n", accountCreatedFelt.String())
		} else if event.Keys[0].Cmp(transactionExecutedFelt) == 0 {
			fmt.Printf("TransactionExecuted event %s\n", transactionExecutedFelt.String())
		} else {
			fmt.Printf("event type %s\n", event.Keys[0].String())
		}
		fmt.Printf("from: %s\n", event.FromAddress.String())
		fmt.Printf("tx: %s\n", event.TransactionHash.String())
		for i, key := range event.Keys {
			fmt.Printf("key %d: %s\n", i, key.String())
		}
		for i, data := range event.Data {
			fmt.Printf("data %d: %s\n", i, data.String())
		}
	}
}
