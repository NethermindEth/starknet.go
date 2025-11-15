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

	// Use ETH token contract address on Sepolia
	contractAddress, err := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	if err != nil {
		log.Fatal("Failed to parse contract address:", err)
	}

	blockID := rpc.WithBlockTag("latest")
	classHash, err := client.ClassHashAt(ctx, blockID, contractAddress)
	if err != nil {
		log.Fatal("Failed to get class hash at address:", err)
	}

	fmt.Printf("Contract Address: %s\n", contractAddress)
	fmt.Printf("Class Hash: %s\n", classHash)
}
