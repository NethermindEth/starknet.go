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

	// Get transaction at index 0 in the latest block
	blockID := rpc.WithBlockTag("latest")
	index := uint64(0)

	transaction, err := client.TransactionByBlockIDAndIndex(ctx, blockID, index)
	if err != nil {
		log.Fatal("Failed to get transaction:", err)
	}

	// Marshal to JSON for readable output
	txJSON, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal transaction:", err)
	}

	fmt.Printf("Transaction at index %d in latest block:\n%s\n", index, string(txJSON))
}
