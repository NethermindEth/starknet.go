package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Testing SubscribeEvents WebSocket method...")

	// Load environment variables
	godotenv.Load("../../.env")
	wsURL := os.Getenv("STARKNET_WS_URL")
	if wsURL == "" {
		log.Fatal("STARKNET_WS_URL not set in .env")
	}

	ctx := context.Background()

	// Connect to WebSocket endpoint
	wsProvider, err := rpc.NewWebsocketProvider(ctx, wsURL)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer wsProvider.Close()

	// Create channel to receive events
	eventsChan := make(chan *rpc.EmittedEventWithFinalityStatus)

	// Subscribe to all events (no filter)
	sub, err := wsProvider.SubscribeEvents(ctx, eventsChan, nil)
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("✅ Successfully subscribed to events")
	fmt.Println("Monitoring all events on the network...")
	fmt.Println("Waiting for events (will timeout after 45 seconds)...\n")

	// Listen for events with timeout
	timeout := time.After(45 * time.Second)
	eventCount := 0

	for {
		select {
		case event := <-eventsChan:
			eventCount++
			fmt.Printf("Event #%d received:\n", eventCount)
			fmt.Printf("  Block Number: %d\n", event.BlockNumber)
			fmt.Printf("  Transaction Hash: %s\n", event.TransactionHash)
			fmt.Printf("  From Address: %s\n", event.FromAddress)
			fmt.Printf("  Keys: %d\n", len(event.Keys))
			fmt.Printf("  Data: %d items\n", len(event.Data))
			fmt.Printf("  Finality: %s\n\n", event.FinalityStatus)

			if eventCount >= 3 {
				fmt.Printf("✅ Successfully received %d events, test passed!\n", eventCount)
				return
			}
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case <-timeout:
			if eventCount > 0 {
				fmt.Printf("✅ Test completed - received %d event(s)\n", eventCount)
			} else {
				fmt.Println("⏱️  Timeout - no events received in 45s (this is normal if network is quiet)")
			}
			return
		}
	}
}
