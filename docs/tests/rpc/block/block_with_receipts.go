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

	// Create RPC client
	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Get block with receipts using latest block tag
	blockID := rpc.WithBlockTag("latest")
	block, err := client.BlockWithReceipts(ctx, blockID)
	if err != nil {
		log.Fatal("Failed to get block with receipts:", err)
	}

	// Marshal to JSON for readable output (show first 2000 chars)
	blockJSON, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal block:", err)
	}

	if len(blockJSON) > 2000 {
		fmt.Printf("Block with Receipts (first 2000 chars):\n%s...\n", string(blockJSON[:2000]))
		fmt.Printf("\nTotal JSON size: %d bytes\n", len(blockJSON))
	} else {
		fmt.Printf("Block with Receipts:\n%s\n", string(blockJSON))
	}
}
