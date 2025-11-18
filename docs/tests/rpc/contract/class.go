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

	// Use ETH contract class hash
	classHash, err := new(felt.Felt).SetString("0x046ded64ae2dead6448e247234bab192a9c483644395b66f2155f2614e5804b0")
	if err != nil {
		log.Fatal("Failed to parse class hash:", err)
	}

	blockID := rpc.WithBlockTag("latest")
	classOutput, err := client.Class(ctx, blockID, classHash)
	if err != nil {
		log.Fatal("Failed to get class:", err)
	}

	// Marshal to JSON for readable output (truncate to first 100 lines)
	classJSON, err := json.MarshalIndent(classOutput, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal class:", err)
	}

	fmt.Printf("Contract Class (first 1000 chars):\n%s...\n", string(classJSON[:1000]))
	fmt.Printf("\nTotal JSON size: %d bytes\n", len(classJSON))
}
