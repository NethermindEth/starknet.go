package main
 
import (
	"context"
	"fmt"
	"log"
	"os"
 
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)
 
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
 
	// Get RPC URL from environment variable
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not found in .env file")
	}
 
	// Initialize provider
	provider, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Fatal(err)
	}
 
	ctx := context.Background()
 
	// Query events from a specific block range
	eventsInput := rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.WithBlockNumber(1676500),
			ToBlock:   rpc.WithBlockNumber(1676510),
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 10,
		},
	}
 
	// Get first page of events
	eventChunk, err := provider.Events(ctx, eventsInput)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Printf("Found %d events\n", len(eventChunk.Events))
 
	// Display event details
	for i, event := range eventChunk.Events {
		fmt.Printf("\nEvent %d:\n", i+1)
		fmt.Printf("  From Contract: %s\n", event.FromAddress)
		fmt.Printf("  Block Number: %d\n", event.BlockNumber)
		fmt.Printf("  Transaction Hash: %s\n", event.TransactionHash)
		fmt.Printf("  Number of Keys: %d\n", len(event.Keys))
		fmt.Printf("  Number of Data Fields: %d\n", len(event.Data))
	}
 
	// Fetch next page if available
	if eventChunk.ContinuationToken != "" {
		fmt.Printf("\nMore events available. Continuation token: %s\n", eventChunk.ContinuationToken)
 
		// Use continuation token for next page
		eventsInput.ContinuationToken = eventChunk.ContinuationToken
		nextPage, err := provider.Events(ctx, eventsInput)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Next page has %d events\n", len(nextPage.Events))
	}
}