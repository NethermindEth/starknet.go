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
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
 
    // Get RPC URL from environment variable
    rpcURL := os.Getenv("STARKNET_RPC_URL")
    if rpcURL == "" {
        log.Fatal("STARKNET_RPC_URL not found in .env file")
    }
 
    // Initialize provider
    provider, err := rpc.NewProvider(context.Background(), rpcURL)
    if err != nil {
        log.Fatal(err)
    }
 
    ctx := context.Background()
 
    // Get a transaction hash from the latest block
    blockID := rpc.WithBlockTag("latest")
    blockResult, err := provider.BlockWithTxHashes(ctx, blockID)
    if err != nil {
        log.Fatal(err)
    }

    var transactionHashes   []*felt.Felt

    switch {
	case blockResult.Block != nil:
		transactionHashes = blockResult.Block.Transactions
	case blockResult.PreConfirmed != nil:
		transactionHashes = blockResult.PreConfirmed.Transactions
	default:
		log.Fatal("Unexpected block output")
	}

    if len(transactionHashes) == 0 {
        log.Fatal("No transaction in latest block")
    }

    txHash := transactionHashes[0]

	// Get transaction status
	status, err := provider.TransactionStatus(ctx, txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction Hash: %s\n", txHash)
	fmt.Printf("Finality Status: %s\n", status.FinalityStatus)
	fmt.Printf("Execution Status: %s\n", status.ExecutionStatus)
	if status.FailureReason != "" {
		fmt.Printf("Failure Reason: %s\n", status.FailureReason)
	}
}