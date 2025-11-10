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
	// Load environment variables
	godotenv.Load(".env")
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	// Create context
	ctx := context.Background()

	// Create RPC client
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Get block with full transactions using latest block tag
	blockID := rpc.WithBlockTag("latest")
	block, err := client.BlockWithTxs(ctx, blockID)
	if err != nil {
		log.Fatal("Failed to get block with transactions:", err)
	}

	// Convert to JSON for readable output
	blockJSON, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal block:", err)
	}

	fmt.Printf("Block with Transactions:\n%s\n", string(blockJSON))
}
