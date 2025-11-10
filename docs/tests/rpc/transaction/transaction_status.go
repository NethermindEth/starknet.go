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

	// Use a known transaction hash from Sepolia
	txHash, err := new(felt.Felt).SetString("0x2acec06e389c6bba5c4d693befd73377828facadf16eadb3f49f0a6ec201408")
	if err != nil {
		log.Fatal("Failed to parse transaction hash:", err)
	}

	status, err := client.TransactionStatus(ctx, txHash)
	if err != nil {
		log.Fatal("Failed to get transaction status:", err)
	}

	fmt.Printf("Transaction Status:\n")
	fmt.Printf("  Finality Status: %s\n", status.FinalityStatus)

	if status.ExecutionStatus != "" {
		fmt.Printf("  Execution Status: %s\n", status.ExecutionStatus)
	}

	if status.FailureReason != "" {
		fmt.Printf("  Failure Reason: %s\n", status.FailureReason)
	}
}
