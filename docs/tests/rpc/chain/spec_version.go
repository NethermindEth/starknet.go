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

	// Get spec version
	specVersion, err := client.SpecVersion(ctx)
	if err != nil {
		log.Fatal("Failed to get spec version:", err)
	}

	fmt.Printf("Spec Version: %s\n", specVersion)
}
