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
	fmt.Println("Testing SubscribeNewHeads WebSocket method...")

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

	// Create channel to receive block headers
	headersChan := make(chan *rpc.BlockHeader)

	// Subscribe from latest block (pass nil or empty struct)
	sub, err := wsProvider.SubscribeNewHeads(ctx, headersChan, rpc.SubscriptionBlockID{Tag: "latest"})
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("✅ Successfully subscribed to new block headers")
	fmt.Println("Waiting for new blocks (will timeout after 120 seconds)...\n")

	// Listen for new blocks with timeout
	timeout := time.After(120 * time.Second)
	blockCount := 0

	for {
		select {
		case header := <-headersChan:
			blockCount++
			fmt.Printf("Block #%d:\n", blockCount)
			fmt.Printf("  Block Number: %d\n", header.Number)
			fmt.Printf("  Block Hash: %s\n", header.Hash)
			fmt.Printf("  Timestamp: %d\n", header.Timestamp)
			fmt.Printf("  Sequencer: %s\n", header.SequencerAddress)
			fmt.Printf("  Parent Hash: %s\n\n", header.ParentHash)

			if blockCount >= 2 {
				fmt.Printf("✅ Successfully received %d blocks, test passed!\n", blockCount)
				return
			}
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case <-timeout:
			if blockCount > 0 {
				fmt.Printf("✅ Test completed - received %d block(s)\n", blockCount)
			} else {
				fmt.Println("⏱️  Timeout - no blocks received in 120s")
			}
			return
		}
	}
}
