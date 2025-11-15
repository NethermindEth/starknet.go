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
 
    // Get latest block with receipts
    blockID := rpc.WithBlockTag("latest")
    block, err := provider.BlockWithReceipts(context.Background(), blockID)
    if err != nil {
        log.Fatal(err)
    }
 
    // Pretty print the result
    blockJSON, _ := json.MarshalIndent(block, "", "  ")
    fmt.Printf("Block with receipts:\n%s\n", blockJSON)
}