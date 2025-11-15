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
	fmt.Println("Testing SubscribeNewTransactionReceipts WebSocket method...")

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

	// Create channel to receive transaction receipts
	receiptsChan := make(chan *rpc.TransactionReceiptWithBlockInfo)

	// Subscribe to all transaction receipts (no filter)
	sub, err := wsProvider.SubscribeNewTransactionReceipts(ctx, receiptsChan, nil)
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("✅ Successfully subscribed to transaction receipts")
	fmt.Println("Waiting for receipts (will timeout after 60 seconds)...\n")

	// Listen for receipts with timeout
	timeout := time.After(60 * time.Second)
	receiptCount := 0

	for {
		select {
		case receipt := <-receiptsChan:
			receiptCount++
			fmt.Printf("Receipt #%d:\n", receiptCount)
			fmt.Printf("  Block Number: %d\n", receipt.BlockNumber)
			fmt.Printf("  Block Hash: %s\n", receipt.BlockHash)

			// Access receipt details through the embedded TransactionReceipt interface
			fmt.Printf("  Type: %T\n", receipt.TransactionReceipt)
			fmt.Println()

			if receiptCount >= 3 {
				fmt.Printf("✅ Successfully received %d receipts, test passed!\n", receiptCount)
				return
			}
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case <-timeout:
			if receiptCount > 0 {
				fmt.Printf("✅ Test completed - received %d receipt(s)\n", receiptCount)
			} else {
				fmt.Println("⏱️  Timeout - no receipts received in 60s (this is normal if network is quiet)")
			}
			return
		}
	}
}
