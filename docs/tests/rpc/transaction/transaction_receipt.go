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
	txHash, err := new(felt.Felt).SetString("0x2acec06e389c6bba5c4d693befd73377828facadf16eadb3f49f0a6ec201408")
	if err != nil {
		log.Fatal("Failed to parse transaction hash:", err)
	}

	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Fatal("Failed to get transaction receipt:", err)
	}

	// Marshal to JSON for readable output
	receiptJSON, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal receipt:", err)
	}

	fmt.Printf("Transaction Receipt:\n%s\n", string(receiptJSON))
}
