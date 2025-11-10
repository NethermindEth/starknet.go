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
	// Load environment variables
	godotenv.Load(".env")
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	// Create RPC client
	ctx := context.Background()

	// Create RPC client
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Get transaction count for latest block
	blockID := rpc.WithBlockTag("latest")
	count, err := client.BlockTransactionCount(ctx, blockID)
	if err != nil {
		log.Fatal("Failed to get transaction count:", err)
	}

	fmt.Printf("Transaction Count: %d\n", count)
}
