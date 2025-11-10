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

	// Example L1 transaction hash (Ethereum mainnet format)
	// NOTE: This is a demo hash. To test this method, you need a real L1 transaction hash
	// from Ethereum that sent L1->L2 messages to Starknet.
	// Format: "0x..." as NumAsHex
	l1TxHash := rpc.NumAsHex("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	statuses, err := client.MessagesStatus(ctx, l1TxHash)
	if err != nil {
		// Expected to fail with demo hash - this shows the method signature and usage
		fmt.Printf("Note: This demo hash doesn't exist. To test with a real hash:\n")
		fmt.Printf("1. Find an L1 transaction that sent messages to Starknet\n")
		fmt.Printf("2. Replace the l1TxHash with the real Ethereum transaction hash\n")
		fmt.Printf("3. Run the test again\n\n")
		fmt.Printf("Error (expected): %v\n", err)
		return
	}

	// Marshal to JSON for readable output
	statusJSON, err := json.MarshalIndent(statuses, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal statuses:", err)
	}

	fmt.Printf("L1->L2 Message Statuses:\n%s\n", string(statusJSON))
	fmt.Printf("\nNumber of messages: %d\n", len(statuses))
}
