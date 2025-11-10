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

	// Use ETH token contract address
	contractAddress, err := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	if err != nil {
		log.Fatal("Failed to parse contract address:", err)
	}

	// Create storage proof input
	storageKey, _ := new(felt.Felt).SetString("0x0")
	proofInput := rpc.StorageProofInput{
		BlockID:           rpc.WithBlockTag("latest"),
		ContractAddresses: []*felt.Felt{contractAddress},
		ContractsStorageKeys: []rpc.ContractStorageKeys{
			{
				ContractAddress: contractAddress,
				StorageKeys:     []*felt.Felt{storageKey},
			},
		},
	}

	proof, err := client.StorageProof(ctx, proofInput)
	if err != nil {
		log.Fatal("Failed to get storage proof:", err)
	}

	// Marshal to JSON for readable output
	proofJSON, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal proof:", err)
	}

	fmt.Printf("Storage Proof:\n%s\n", string(proofJSON))
}
