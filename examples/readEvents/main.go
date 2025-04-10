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
	CONTRACT_ADDRESS = "0x049D36570D4e46f48e99674bd3fcc84644DdD6b96F7C741B1562B82f9e004dC7" // StarkGate: ETH Token
)

// main entry point of the program.
//
// It initializes the environment and establishes a connection with the client.
// It then retrieves events from the contract and prints how many it found with
// a series of more selective filters (all event types, just 2 event types, just
// those 2 event types but with a specified key value).
func main() {
	// Read provider URL from .env file
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Initialize connection to RPC provider
	provider, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %v", err))
	}
	fmt.Println("Established connection with the RPC provider")

	contractAddress, err := utils.HexToFelt(CONTRACT_ADDRESS)
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the contract address %s, error %v", CONTRACT_ADDRESS, err))
	}

	// 1. call with ChunkSize and ContinuationToken
	callWithChunkSizeAndContinuationToken(provider)
	// 2. call with Block and Address filters
	callWithBlockAndAddressFilters(provider, contractAddress)
	// 3. call with Keys filter
	callWithKeysFilter(provider, contractAddress)

	// moreComplexExample(provider)
}

func callWithChunkSizeAndContinuationToken(provider *rpc.Provider) {
	fmt.Println()
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
	fmt.Printf("block number of the first event in the first chunk: %d\n", eventChunk.Events[0].BlockNumber)
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
	fmt.Printf("block number of the first event in the second chunk: %d\n", secondEventChunk.Events[0].BlockNumber)
	fmt.Printf("block number of the last event in the second chunk: %d\n", secondEventChunk.Events[len(secondEventChunk.Events)-1].BlockNumber)
}

func callWithBlockAndAddressFilters(provider *rpc.Provider, contractAddress *felt.Felt) {
	fmt.Println()
	fmt.Println(" ----- 2. call with Block and Address filters -----")
	fmt.Println("Contract Address: ", contractAddress.String())

	// We are using the following filters:
	// - FromBlock: The starting block number (inclusive)
	// - ToBlock: The ending block number (inclusive)
	// - Address: The contract address to filter events from
	//
	// So, we are filtering events from block 0 to block 100 and only from the provided contract address.
	eventChunk, err := provider.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.WithBlockNumber(0),
			ToBlock:   rpc.WithBlockNumber(100),
			Address:   contractAddress,
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}
	fmt.Printf("number of returned events: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the first event: %d\n", eventChunk.Events[0].BlockNumber)
	fmt.Printf("block number of the last event: %d\n", eventChunk.Events[len(eventChunk.Events)-1].BlockNumber)
	fmt.Printf("contract address of the first event: %s\n", eventChunk.Events[0].FromAddress.String())
}

func callWithKeysFilter(provider *rpc.Provider, contractAddress *felt.Felt) {
	fmt.Println()
	fmt.Println(" ----- 3. call with Keys filter -----")
	fmt.Println(" --- step 1: filter all 'Transfer' events ---")

	// Firstly, we need to understand how the 'keys' filter works.
	//
	// If the we send an event filter containing [[k_1, k_2], [], [k_3]], then the node should return
	// events whose first key is k_1 or k_2, the third key is k_3, and the second key is unconstrained and can take any value.
	// Ref: https://community.starknet.io/t/snip-13-index-transfer-and-approval-events-in-erc20s/114212
	//
	// The keys are interpreted as follows:
	// - The first key usually is the event selector
	// - The remaining keys will vary depending on the event
	//
	// So here we are filtering all 'Transfer' events (to be more precise, all events with the 'Transfer' selector as the first key)
	// from all addresses and contracts, from block 600000 to block 600100.
	eventChunk, err := provider.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.WithBlockNumber(600000),
			ToBlock:   rpc.WithBlockNumber(600100),
			Keys: [][]*felt.Felt{
				{
					utils.GetSelectorFromNameFelt("Transfer"),
				},
			},
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}

	fmt.Printf("'Transfer' hash selector: %s\n", utils.GetSelectorFromNameFelt("Transfer").String())

	fmt.Printf("number of returned events: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the first event: %d\n", eventChunk.Events[0].BlockNumber)
	fmt.Printf("block number of the last event: %d\n", eventChunk.Events[len(eventChunk.Events)-1].BlockNumber)
	fmt.Printf("first key of the first event: %s\n", eventChunk.Events[0].Keys[0].String())

	fmt.Println()
	fmt.Println(" --- step 2: filter multiple events types ---")

	// Here we are filtering all 'Transfer', 'Approval' and 'GameStarted' events.
	eventChunk, err = provider.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.WithBlockNumber(600000),
			ToBlock:   rpc.WithBlockNumber(600100),
			Keys: [][]*felt.Felt{
				{
					utils.GetSelectorFromNameFelt("Transfer"),
					utils.GetSelectorFromNameFelt("Approval"),
					utils.GetSelectorFromNameFelt("GameStarted"),
				},
			},
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}

	fmt.Printf("'Transfer' hash selector: %s\n", utils.GetSelectorFromNameFelt("Transfer").String())
	fmt.Printf("'Approval' hash selector: %s\n", utils.GetSelectorFromNameFelt("Approval").String())
	fmt.Printf("'GameStarted' hash selector: %s\n", utils.GetSelectorFromNameFelt("GameStarted").String())

	fmt.Printf("number of returned events: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the first event: %d\n", eventChunk.Events[0].BlockNumber)
	fmt.Printf("block number of the last event: %d\n", eventChunk.Events[len(eventChunk.Events)-1].BlockNumber)
	transferEvent := findEventInChunk(eventChunk, "Transfer")
	fmt.Printf("'Transfer' event found in block %d, tx hash: %s\n", transferEvent.BlockNumber, transferEvent.TransactionHash.String())
	gameStartedEvent := findEventInChunk(eventChunk, "GameStarted")
	fmt.Printf("'GameStarted' event found in block %d, tx hash: %s\n", gameStartedEvent.BlockNumber, gameStartedEvent.TransactionHash.String())
	approvalEvent := findEventInChunk(eventChunk, "Approval")
	fmt.Printf("'Approval' event found in block %d, tx hash: %s\n", approvalEvent.BlockNumber, approvalEvent.TransactionHash.String())
}

func findEventInChunk(eventChunk *rpc.EventChunk, eventName string) rpc.EmittedEvent {
	selector := utils.GetSelectorFromNameFelt(eventName)

	for _, event := range eventChunk.Events {
		if event.Keys[0].String() == selector.String() {
			return event
		}
	}
	return rpc.EmittedEvent{}
}
