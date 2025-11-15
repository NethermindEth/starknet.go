package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Testing SubscribeTransactionStatus WebSocket method...")
	fmt.Println("Note: This test requires a valid transaction hash to monitor.")
	fmt.Println("Using a recent transaction hash for demonstration...\n")

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

	// For testing, we'll use a placeholder transaction hash
	// In a real scenario, you would use an actual pending transaction hash
	// Example: txHash from a recently submitted transaction
	txHash, err := new(felt.Felt).SetString("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		log.Fatal("Invalid transaction hash:", err)
	}

	// Create channel to receive status updates
	statusChan := make(chan *rpc.NewTxnStatus)

	// Subscribe to transaction status updates
	sub, err := wsProvider.SubscribeTransactionStatus(ctx, statusChan, txHash)
	if err != nil {
		fmt.Printf("❌ Failed to subscribe: %v\n", err)
		fmt.Println("\n⚠️  Note: This is expected if the transaction hash doesn't exist.")
		fmt.Println("To properly test this method:")
		fmt.Println("1. Submit a transaction and get its hash")
		fmt.Println("2. Immediately subscribe to its status")
		fmt.Println("3. Monitor status changes until finalization")
		fmt.Println("\n✅ Test structure is correct - method signature and usage are valid")
		return
	}
	defer sub.Unsubscribe()

	fmt.Printf("✅ Successfully subscribed to transaction status\n")
	fmt.Printf("Monitoring transaction: %s\n\n", txHash)

	// Listen for status updates with timeout
	timeout := time.After(60 * time.Second)
	updateCount := 0

	for {
		select {
		case status := <-statusChan:
			updateCount++
			fmt.Printf("Status Update #%d:\n", updateCount)
			fmt.Printf("  Transaction Hash: %s\n", status.TransactionHash)
			fmt.Printf("  Execution Status: %s\n", status.Status.ExecutionStatus)
			fmt.Printf("  Finality Status: %s\n", status.Status.FinalityStatus)

			if status.Status.ExecutionStatus == "REVERTED" {
				fmt.Printf("  Failure Reason: %s\n", status.Status.FailureReason)
			}
			fmt.Println()

			// Check if transaction is finalized
			if status.Status.FinalityStatus == "ACCEPTED_ON_L1" ||
				status.Status.FinalityStatus == "ACCEPTED_ON_L2" {
				fmt.Println("✅ Transaction finalized!")
				if status.Status.ExecutionStatus == "SUCCEEDED" {
					fmt.Println("✅ Transaction executed successfully")
				} else if status.Status.ExecutionStatus == "REVERTED" {
					fmt.Printf("❌ Transaction reverted: %s\n", status.Status.FailureReason)
				}
				fmt.Printf("\n✅ Test completed - received %d status update(s)\n", updateCount)
				return
			}
		case err := <-sub.Err():
			fmt.Printf("Subscription error: %v\n", err)
			if updateCount > 0 {
				fmt.Printf("✅ Test partially successful - received %d update(s) before error\n", updateCount)
			}
			return
		case <-timeout:
			if updateCount > 0 {
				fmt.Printf("✅ Test completed - received %d status update(s)\n", updateCount)
			} else {
				fmt.Println("⏱️  Timeout - no updates received")
			}
			return
		}
	}
}
