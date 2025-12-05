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
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
 
	wsURL := os.Getenv("STARKNET_WS_URL")
	if wsURL == "" {
		log.Fatal("STARKNET_WS_URL not found in .env file")
	}
 
	ctx := context.Background()
 
	// Connect to WebSocket endpoint
	wsProvider, err := rpc.NewWebsocketProvider(ctx, wsURL)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer wsProvider.Close()
 
	// Transaction hash to monitor (replace with actual transaction hash)
	txHash, err := new(felt.Felt).SetString("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		log.Fatal("Invalid transaction hash:", err)
	}
 
	// Create channel to receive status updates
	statusChan := make(chan *rpc.NewTxnStatus)
 
	// Subscribe to transaction status updates
	sub, err := wsProvider.SubscribeTransactionStatus(ctx, statusChan, txHash)
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()
 
	fmt.Printf("Successfully subscribed to transaction status\n")
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
				fmt.Println("Transaction finalized!")
				if status.Status.ExecutionStatus == "SUCCEEDED" {
					fmt.Println("Transaction executed successfully")
				} else if status.Status.ExecutionStatus == "REVERTED" {
					fmt.Printf("Transaction reverted: %s\n", status.Status.FailureReason)
				}
				fmt.Printf("\nTest completed - received %d status update(s)\n", updateCount)
				return
			}
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case <-timeout:
			if updateCount > 0 {
				fmt.Printf("Test completed - received %d status update(s)\n", updateCount)
			} else {
				fmt.Println("Timeout - no updates received (transaction hash may not exist)")
			}
			return
		}
	}
}