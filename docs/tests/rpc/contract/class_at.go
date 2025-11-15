package main

import (
	"context"
	"encoding/json"
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

	// Use ETH token contract address on Sepolia
	contractAddress, err := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	if err != nil {
		log.Fatal("Failed to parse contract address:", err)
	}

	blockID := rpc.WithBlockTag("latest")
	classOutput, err := client.ClassAt(ctx, blockID, contractAddress)
	if err != nil {
		log.Fatal("Failed to get class at address:", err)
	}

	// Marshal to JSON for readable output (show first 1000 chars)
	classJSON, err := json.MarshalIndent(classOutput, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal class:", err)
	}

	fmt.Printf("Contract Class at Address (first 1000 chars):\n%s...\n", string(classJSON[:1000]))
	fmt.Printf("\nTotal JSON size: %d bytes\n", len(classJSON))
}
