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

	// Get latest block with transaction hashes
	blockID := rpc.WithBlockTag("latest")
	blockOutput, err := provider.BlockWithTxHashes(context.Background(), blockID)
	if err != nil {
		log.Fatal(err)
	}

	var block any
	switch {
	case blockOutput.Block != nil:
		block = blockOutput.Block
	case blockOutput.PreConfirmed != nil:
		block = blockOutput.PreConfirmed
	default:
		log.Fatal("Unexpected block output")
	}

	// Pretty print the result
	blockJSON, _ := json.MarshalIndent(block, "", "  ")
	fmt.Printf("Block with transaction hashes:\n%s\n", blockJSON)
}