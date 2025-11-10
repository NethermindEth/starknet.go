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

	// Create context
	ctx := context.Background()

	// Create RPC client
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Get latest block number
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal("Failed to get block number:", err)
	}

	fmt.Printf("Block Number: %d\n", blockNumber)
}
