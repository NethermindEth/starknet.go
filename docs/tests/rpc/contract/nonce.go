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

	// Use a test account address from .env
	accountAddress := os.Getenv("TEST_ACCOUNT_ADDRESS")
	if accountAddress == "" {
		log.Fatal("TEST_ACCOUNT_ADDRESS not set in .env")
	}

	contractAddress, err := new(felt.Felt).SetString(accountAddress)
	if err != nil {
		log.Fatal("Failed to parse contract address:", err)
	}

	blockID := rpc.WithBlockTag("latest")
	nonce, err := client.Nonce(ctx, blockID, contractAddress)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	fmt.Printf("Contract Address: %s\n", contractAddress)
	fmt.Printf("Nonce: %s\n", nonce)
}
