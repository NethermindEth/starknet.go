package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	// Get traces for all transactions in latest block
	blockID := rpc.WithBlockTag("latest")
	traces, err := client.TraceBlockTransactions(ctx, blockID)
	if err != nil {
		log.Fatal("Failed to get block transaction traces:", err)
	}

	// Show summary
	fmt.Printf("Number of transaction traces: %d\n", len(traces))

	if len(traces) > 0 {
		// Show first trace in detail
		traceJSON, err := json.MarshalIndent(traces[0], "", "  ")
		if err != nil {
			log.Fatal("Failed to marshal trace:", err)
		}

		if len(traceJSON) > 1000 {
			fmt.Printf("\nFirst transaction trace (first 1000 chars):\n%s...\n", string(traceJSON[:1000]))
		} else {
			fmt.Printf("\nFirst transaction trace:\n%s\n", string(traceJSON))
		}
	}
}
