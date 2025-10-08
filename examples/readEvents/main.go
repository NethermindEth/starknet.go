package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// main is the entry point of the program that demonstrates how to query Starknet events.
//
// This example shows how to:
// 1. Connect to a Starknet RPC provider
// 2. Query events with pagination using ChunkSize and ContinuationToken
// 3. Filter events by block range and contract address
// 4. Filter events by specific event keys
// 5. Combine multiple filters for precise event retrieval
//
// The program progressively applies more selective filters to demonstrate
// different ways of narrowing down event queries for efficient data retrieval.
func main() {
	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	wsProviderUrl := setup.GetWsProviderUrl()

	// Initialise connection to RPC provider
	provider, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the RPC provider: %v", err))
	}
	fmt.Println("Established connection with the RPC provider")

	// Now we will call some functions to demonstrate the different filters we can apply.
	// Enter each function declaration to see the filters in action.

	// 1. call with ChunkSize and ContinuationToken
	callWithChunkSizeAndContinuationToken(provider)
	// 2. call with Block and Address filters
	callWithBlockAndAddressFilters(provider)
	// 3. call with Keys filter
	callWithKeysFilter(provider)
	// optional: filter with websocket
	filterWithWebsocket(
		provider,
		wsProviderUrl,
	) // if the wsProviderUrl is empty, the websocket example will be skipped

	// after all, here is a call with all filters combined
	fmt.Println("\n ----- 4. all filters -----")

	contractAddress, err := utils.HexToFelt(
		"0x1948e239f559bcbdf9388938a3c46bc79f52bcba7c4d5c9732568cb8eb6a53d",
	) // a random contract address for our example
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the contract address, error %v", err))
	}
	key3, err := utils.HexToFelt(
		"0x1bfc84464f990c09cc0e5d64d18f54c3469fd5c467398bf31293051bade1c39",
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the provided key, error %v", err))
	}
	key5, err := utils.HexToFelt("0x0")
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the provided key, error %v", err))
	}

	eventChunk, err := provider.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.WithBlockNumber(660000), // from block 660000
			ToBlock:   rpc.WithBlockNumber(660100), // to block 660100
			Address:   contractAddress,             // sent from this contract address
			Keys: [][]*felt.Felt{
				// Here we are filtering all 'Transfer', 'Approval' and 'GameStarted' events.
				// (all events that have one of these selectors as the first key)
				{
					utils.GetSelectorFromNameFelt("Transfer"),
					utils.GetSelectorFromNameFelt("Approval"),
					utils.GetSelectorFromNameFelt("GameStarted"),
				},
				{},     // here we are saying that the second key is unconstrained, it can take any value
				{key3}, // the third key must be equal to key3
				{},     // the fourth key is also unconstrained
				{key5}, // the fifth key must be equal to key5
			},
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	}) // so this will return all events, between block 660000 and 660100, sent from the specified contract address,
	// that have one of the 'Transfer', 'Approval' or 'GameStarted' selectors as the first key,
	// and the third key is equal to key3, and the fifth key is equal to key5; the second and fourth keys can be any value
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}

	fmt.Printf("number of returned events: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the first event: %d\n", eventChunk.Events[0].BlockNumber)
	fmt.Printf(
		"block number of the last event: %d\n",
		eventChunk.Events[len(eventChunk.Events)-1].BlockNumber,
	)
	randomEvent := eventChunk.Events[rand.Intn(len(eventChunk.Events))] // get a random event from the chunk
	fmt.Printf("random event block number: %d\n", randomEvent.BlockNumber)
	fmt.Printf("random event tx hash: %s\n", randomEvent.TransactionHash.String())
	fmt.Printf("random event sender address: %s\n", randomEvent.FromAddress.String())
	fmt.Printf("random event first key: %v\n", randomEvent.Keys[0].String())
	fmt.Printf("random event third key: %v\n", randomEvent.Keys[2].String())
	fmt.Printf("random event fifth key: %v\n", randomEvent.Keys[4].String())
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
	fmt.Printf(
		"block number of the first event in the first chunk: %d\n",
		eventChunk.Events[0].BlockNumber,
	)
	fmt.Printf(
		"block number of the last event in the first chunk: %d\n",
		eventChunk.Events[len(eventChunk.Events)-1].BlockNumber,
	)

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
	fmt.Printf(
		"block number of the first event in the second chunk: %d\n",
		secondEventChunk.Events[0].BlockNumber,
	)
	fmt.Printf(
		"block number of the last event in the second chunk: %d\n",
		secondEventChunk.Events[len(secondEventChunk.Events)-1].BlockNumber,
	)
}

func callWithBlockAndAddressFilters(provider *rpc.Provider) {
	fmt.Println()
	fmt.Println(" ----- 2. call with Block and Address filters -----")
	contractAddress, err := utils.HexToFelt(
		"0x049D36570D4e46f48e99674bd3fcc84644DdD6b96F7C741B1562B82f9e004dC7",
	) // StarkGate: ETH Token
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the contract address, error %v", err))
	}

	fmt.Println("Contract Address: ", contractAddress.String())

	// We are using the following filters:
	//   - FromBlock: The starting block number (inclusive)
	//   - ToBlock: The ending block number (inclusive)
	//   - Address: The contract address to filter events from
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
	fmt.Printf(
		"block number of the last event: %d\n",
		eventChunk.Events[len(eventChunk.Events)-1].BlockNumber,
	)
	fmt.Printf(
		"contract address of the first event: %s\n",
		eventChunk.Events[0].FromAddress.String(),
	)
}

func callWithKeysFilter(provider *rpc.Provider) {
	fmt.Println()
	fmt.Println(" ----- 3. call with Keys filter -----")
	fmt.Println(" --- step 1: filter all events with the 'Transfer' name ---")

	// Firstly, we need to understand how the 'keys' filter works.
	// (and of course, we need to know what 'keys' are, right? read more here:
	// 	https://book.cairo-lang.org/ch101-03-contract-events.html,
	// 	https://docs.starknet.io/architecture-and-concepts/smart-contracts/starknet-events/
	// )
	//
	// If the we send an event filter containing [[k_1, k_2], [], [k_3]], then the node should return
	// events whose first key is k_1 or k_2, the third key is k_3, and the second key is unconstrained and can take any value.
	// Ref: https://community.starknet.io/t/snip-13-index-transfer-and-approval-events-in-erc20s/114212
	//
	// The keys are interpreted as follows:
	//   - The first key usually is the event selector
	//   - The remaining keys will vary depending on the event
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
	// NOTE: An event can be nested in a Cairo component (See the Cairo code of the contract to verify).
	// In this case, the array of keys will start with additional hashes, and you will have to adapt your code in consequence
	// Ref: https://starknetjs.com/docs/guides/events#without-transaction-hash
	if err != nil {
		panic(fmt.Sprintf("error retrieving events: %v", err))
	}

	fmt.Printf("'Transfer' hash selector: %s\n", utils.GetSelectorFromNameFelt("Transfer").String())

	fmt.Printf("number of returned events: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the first event: %d\n", eventChunk.Events[0].BlockNumber)
	fmt.Printf(
		"block number of the last event: %d\n",
		eventChunk.Events[len(eventChunk.Events)-1].BlockNumber,
	)
	fmt.Printf("first key of the first event: %s\n", eventChunk.Events[0].Keys[0].String())

	fmt.Println()
	fmt.Println(" --- step 2: filter multiple events types ---")

	// Here we are filtering all 'Transfer', 'Approval' and 'GameStarted' events.
	eventChunk, err = provider.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.WithBlockNumber(600000),
			ToBlock:   rpc.WithBlockNumber(600100),
			Keys: [][]*felt.Felt{
				// Notice that we are passing all selectors together in the same array, meaning that
				// the node will return events that match any of these values.
				// Also notice that the array is in the first position of the array, so basically
				// we are filtering all events that have one of these selectors as the first key.
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
	fmt.Printf(
		"'GameStarted' hash selector: %s\n",
		utils.GetSelectorFromNameFelt("GameStarted").String(),
	)

	fmt.Printf("number of returned events: %d\n", len(eventChunk.Events))
	fmt.Printf("block number of the first event: %d\n", eventChunk.Events[0].BlockNumber)
	fmt.Printf(
		"block number of the last event: %d\n",
		eventChunk.Events[len(eventChunk.Events)-1].BlockNumber,
	)
	transferEvent := findEventInChunk(eventChunk, "Transfer")
	fmt.Printf(
		"'Transfer' event found in block %d, tx hash: %s\n",
		transferEvent.BlockNumber,
		transferEvent.TransactionHash.String(),
	)
	gameStartedEvent := findEventInChunk(eventChunk, "GameStarted")
	fmt.Printf(
		"'GameStarted' event found in block %d, tx hash: %s\n",
		gameStartedEvent.BlockNumber,
		gameStartedEvent.TransactionHash.String(),
	)
	approvalEvent := findEventInChunk(eventChunk, "Approval")
	fmt.Printf(
		"'Approval' event found in block %d, tx hash: %s\n",
		approvalEvent.BlockNumber,
		approvalEvent.TransactionHash.String(),
	)
}

func filterWithWebsocket(provider *rpc.Provider, websocketUrl string) {
	if websocketUrl == "" {
		fmt.Println("\nNo websocket URL provided. Skipping websocket filter...")

		return
	}

	fmt.Println()
	fmt.Println(" ----- 4. filter with websocket -----")

	wsProvider, err := rpc.NewWebsocketProvider(websocketUrl)
	if err != nil {
		panic(fmt.Sprintf("error dialling the RPC provider: %v", err))
	}
	contractAddress, err := utils.HexToFelt(
		"0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D",
	) // StarkGate: ETH Token
	if err != nil {
		panic(fmt.Sprintf("failed to create felt from the contract address, error %v", err))
	}

	// Get the latest block number
	blockNumber, err := provider.BlockNumber(context.Background())
	if err != nil {
		panic(fmt.Sprintf("error getting the latest block number: %v", err))
	}

	// Create a channel to receive events
	eventsChan := make(chan *rpc.EmittedEventWithFinalityStatus)

	// Subscribe to events
	sub, err := wsProvider.SubscribeEvents(
		context.Background(),
		eventsChan,
		&rpc.EventSubscriptionInput{
			// Only events from this contract address
			FromAddress: contractAddress,
			// Subscribe to events from the latest block minus 10 (it'll return
			// events from the last 10 blocks and progressively update as new blocks are added)
			SubBlockID: new(rpc.SubscriptionBlockID).WithBlockNumber(blockNumber - 10),
			Keys: [][]*felt.Felt{
				// the 'keys'filter behaves the same way as the RPC provider `starknet_getEvents` explained above.
				// So this will return all events that have the 'Transfer' selector as the first key.
				{
					utils.GetSelectorFromNameFelt("Transfer"),
				},
			},
		},
	)
	if err != nil {
		panic(fmt.Sprintf("error subscribing to events: %v", err))
	}

	fmt.Println("Successfully subscribed to events")

	// Read events from the channel
	for {
		select {
		case event := <-eventsChan:
			// This case will be triggered when a new event is received.
			fmt.Printf(
				"New event received: Block %d, Event tx hash: %s\n",
				event.BlockNumber,
				event.TransactionHash.String(),
			)
		case err := <-sub.Err():
			// This case will be triggered when an error occurs.
			panic(err)
		case <-time.After(5 * time.Second):
			// stop the loop after 5 seconds
			fmt.Println("Exiting...")

			return
		}
	}
}

// simple function to find an event by name in a chunk of events
func findEventInChunk(eventChunk *rpc.EventChunk, eventName string) rpc.EmittedEvent {
	selector := utils.GetSelectorFromNameFelt(eventName)

	for _, event := range eventChunk.Events {
		if event.Keys[0].String() == selector.String() {
			return event
		}
	}

	return rpc.EmittedEvent{}
}
