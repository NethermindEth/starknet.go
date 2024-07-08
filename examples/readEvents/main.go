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
// It then retrieves events from the contract and prints them.
//
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
	// accountAddress := setup.GetAccountAddress()

	// Initialize connection to RPC provider
	provider, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the client")

	contractAddrStr := setup.GetContractAddress()
	contractAddress, err := utils.HexToFelt(contractAddrStr)
	if err != nil {
		msg := fmt.Errorf("failed to transform the token contract address %s, error %w", contractAddrStr, err)
		panic(msg)
	}

	fromBlock, toBlock := setup.GetFromAndToBlocks()
	fmt.Printf("FROM_BLOCK: %d\n", fromBlock)
	eventTypes := []string{"Transfer", "Approve"}

	eventIds := make([][]*felt.Felt, 0, len(eventTypes))
	for _, eventType := range eventTypes {
		eventId := utils.GetSelectorFromName(eventType)
		eventFelt := utils.BigIntToFelt(eventId)
		eventIds = append(eventIds, []*felt.Felt{eventFelt})
	}
	fmt.Printf("eventIds %v\n", eventIds)
	eventFilter := rpc.EventFilter{
		FromBlock: rpc.BlockID{Number: &fromBlock},
		ToBlock:   rpc.BlockID{Number: &toBlock},
		// Keys:      eventIds,

	}
	resPageReq := rpc.ResultPageRequest{
		ChunkSize: 100,
	}
	eventInput := rpc.EventsInput{
		EventFilter:       eventFilter,
		ResultPageRequest: resPageReq,
	}
	eventInput.EventFilter.Address = contractAddress
	events := make([]rpc.EmittedEvent, 0, 10)
	eventChunk, err := provider.Events(context.Background(), eventInput)
	if err != nil {
		msg := fmt.Errorf("error retrieving events: %w", err)
		panic(msg)
	}
	events = append(events, eventChunk.Events...)
	chunkCount := 1
	for eventChunk.ContinuationToken != "" {
		eventInput.ContinuationToken = eventChunk.ContinuationToken
		eventChunk, err = provider.Events(context.Background(), eventInput)
		if err != nil {
			msg := fmt.Errorf("error retrieving events: %w", err)
			panic(msg)
		}
		events = append(events, eventChunk.Events...)
		chunkCount++
	}


	// fromJson, err := eventInput.EventFilter.FromBlock.MarshalJSON()
	// if err != nil {
	// 	msg := fmt.Errorf("error marshalling from block to JSON: %w", err)
	// 	panic(msg)
	// }
	// toJson, err := eventInput.EventFilter.ToBlock.MarshalJSON()
	// if err != nil {
	// 	msg := fmt.Errorf("error marshalling to block to JSON: %w", err)
	// 	panic(msg)
	// }
	// fmt.Printf("event filter From %s\n", fromJson)
	// fmt.Printf("event filter To %s\n", toJson)
	// fmt.Println("continuation token...")
	// fmt.Println(eventChunk.ContinuationToken)
	// fmt.Println("...continuation token")
	fmt.Printf("num events %d, num chunks %d\n", len(events), chunkCount)
}
