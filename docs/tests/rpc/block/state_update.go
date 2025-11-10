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

	// Get state update for latest block
	blockID := rpc.WithBlockTag("latest")
	stateUpdate, err := client.StateUpdate(ctx, blockID)
	if err != nil {
		log.Fatal("Failed to get state update:", err)
	}

	fmt.Printf("Block Hash: %s\n", stateUpdate.BlockHash)
	fmt.Printf("New Root: %s\n", stateUpdate.NewRoot)
	fmt.Printf("Old Root: %s\n", stateUpdate.OldRoot)

	stateDiff := stateUpdate.StateDiff
	fmt.Printf("\nState Diff Summary:\n")
	fmt.Printf("  Storage Diffs: %d\n", len(stateDiff.StorageDiffs))
	fmt.Printf("  Deployed Contracts: %d\n", len(stateDiff.DeployedContracts))
	fmt.Printf("  Declared Classes: %d\n", len(stateDiff.DeclaredClasses))
	fmt.Printf("  Nonces: %d\n", len(stateDiff.Nonces))
}
