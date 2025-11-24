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
 
	// Create channel to receive new transactions
	txnsChan := make(chan *rpc.TxnWithHashAndStatus)
 
	// Subscribe to all new transactions (no filter)
	sub, err := wsProvider.SubscribeNewTransactions(ctx, txnsChan, &rpc.SubNewTxnsInput{})
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()
 
	fmt.Println("Successfully subscribed to new transactions")
	fmt.Println("Monitoring new transactions on the network...\n")
 
	// Listen for new transactions with timeout
	timeout := time.After(60 * time.Second)
	txnCount := 0
 
	for {
		select {
		case txn := <-txnsChan:
			txnCount++
			fmt.Printf("Transaction #%d:\n", txnCount)
			fmt.Printf("  Finality Status: %s\n", txn.FinalityStatus)
			fmt.Printf("  Type: %T\n\n", txn.BlockTransaction)
 
			if txnCount >= 3 {
				fmt.Printf("Successfully received %d transactions\n", txnCount)
				return
			}
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case <-timeout:
			if txnCount > 0 {
				fmt.Printf("Test completed - received %d transaction(s)\n", txnCount)
			} else {
				fmt.Println("Timeout - no transactions received in 60s (this is normal if network is quiet)")
			}
			return
		}
	}
}