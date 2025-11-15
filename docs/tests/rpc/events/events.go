package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Create event filter (query recent blocks for any events)
	eventInput := rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.BlockID{Tag: "latest"},
			ToBlock:   rpc.BlockID{Tag: "latest"},
			Address:   nil, // Query all contracts
			Keys:      nil, // Get all events
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 5, // Limit to 5 events
		},
	}

	eventChunk, err := client.Events(ctx, eventInput)
	if err != nil {
		log.Fatal("Failed to get events:", err)
	}

	// Marshal to JSON for readable output
	eventsJSON, err := json.MarshalIndent(eventChunk, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal events:", err)
	}

	fmt.Printf("Events from ETH Contract:\n%s\n", string(eventsJSON))
	fmt.Printf("\nNumber of events: %d\n", len(eventChunk.Events))
}
