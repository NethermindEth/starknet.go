package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
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

	// Use ETH token contract to read total supply storage
	contractAddress, err := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	if err != nil {
		log.Fatal("Failed to parse contract address:", err)
	}

	// Storage key for total supply (this is contract-specific)
	// For ERC20, total supply is typically at a specific storage slot
	storageKey := "0x0"

	blockID := rpc.WithBlockTag("latest")
	value, err := client.StorageAt(ctx, contractAddress, storageKey, blockID)
	if err != nil {
		log.Fatal("Failed to get storage at:", err)
	}

	fmt.Printf("Contract Address: %s\n", contractAddress)
	fmt.Printf("Storage Key: %s\n", storageKey)
	fmt.Printf("Storage Value: %s\n", value)
}
