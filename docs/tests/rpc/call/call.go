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
	// Load environment variables
	godotenv.Load(".env")
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	// Create RPC client
	ctx := context.Background()

	// Create RPC client
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// ETH contract address on Sepolia
	ethContractAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Get selector for "name" function (simpler, always exists)
	nameSelector, _ := new(felt.Felt).SetString("0x361458367e696363fbcc70777d07ebbd2394e89fd0adcaf147faccd1d294d60")

	// Prepare function call
	functionCall := rpc.FunctionCall{
		ContractAddress:    ethContractAddress,
		EntryPointSelector: nameSelector,
		Calldata:           []*felt.Felt{}, // name takes no arguments
	}

	// Call the contract
	result, err := client.Call(ctx, functionCall, rpc.WithBlockTag("latest"))
	if err != nil {
		log.Fatal("Failed to call contract:", err)
	}

	fmt.Printf("Contract Name Result: %s\n", result[0])
}
