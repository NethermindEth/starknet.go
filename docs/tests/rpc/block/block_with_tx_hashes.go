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

	// Get block with transaction hashes using latest block tag
	blockID := rpc.WithBlockTag("latest")
	block, err := client.BlockWithTxHashes(ctx, blockID)
	if err != nil {
		log.Fatal("Failed to get block with tx hashes:", err)
	}

	// Marshal to JSON for readable output
	blockJSON, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal block:", err)
	}

	fmt.Printf("Block with Transaction Hashes:\n%s\n", string(blockJSON))
}
