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

	// Get syncing status
	syncStatus, err := client.Syncing(ctx)
	if err != nil {
		log.Fatal("Failed to get syncing status:", err)
	}

	// Check if syncing
	if syncStatus.IsSyncing {
		fmt.Println("Node is syncing:")
		fmt.Printf("  Starting Block Hash: %s\n", syncStatus.StartingBlockHash)
		fmt.Printf("  Starting Block Num: %d\n", syncStatus.StartingBlockNum)
		fmt.Printf("  Current Block Hash: %s\n", syncStatus.CurrentBlockHash)
		fmt.Printf("  Current Block Num: %d\n", syncStatus.CurrentBlockNum)
		fmt.Printf("  Highest Block Hash: %s\n", syncStatus.HighestBlockHash)
		fmt.Printf("  Highest Block Num: %d\n", syncStatus.HighestBlockNum)
	} else {
		fmt.Println("Node is not syncing (fully synchronized)")
	}
}
