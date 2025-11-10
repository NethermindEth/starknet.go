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

	// Use a known transaction hash from Sepolia
	// This is a real transaction from a recent block
	txHash, err := new(felt.Felt).SetString("0x2acec06e389c6bba5c4d693befd73377828facadf16eadb3f49f0a6ec201408")
	if err != nil {
		log.Fatal("Failed to parse transaction hash:", err)
	}

	transaction, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		log.Fatal("Failed to get transaction:", err)
	}

	// Marshal to JSON for readable output
	txJSON, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal transaction:", err)
	}

	fmt.Printf("Transaction Details:\n%s\n", string(txJSON))
}
