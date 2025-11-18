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

	// Use a real transaction hash
	txHash, err := new(felt.Felt).SetString("0x2acec06e389c6bba5c4d693befd73377828facadf16eadb3f49f0a6ec201408")
	if err != nil {
		log.Fatal("Failed to parse transaction hash:", err)
	}

	trace, err := client.TraceTransaction(ctx, txHash)
	if err != nil {
		log.Fatal("Failed to get transaction trace:", err)
	}

	// Marshal to JSON (show first 2000 chars)
	traceJSON, err := json.MarshalIndent(trace, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal trace:", err)
	}

	if len(traceJSON) > 2000 {
		fmt.Printf("Transaction Trace (first 2000 chars):\n%s...\n", string(traceJSON[:2000]))
		fmt.Printf("\nTotal JSON size: %d bytes\n", len(traceJSON))
	} else {
		fmt.Printf("Transaction Trace:\n%s\n", string(traceJSON))
	}
}
