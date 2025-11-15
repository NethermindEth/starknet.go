package main
 
import (
    "context"
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
 
    // Get latest block hash and number
    result, err := provider.BlockHashAndNumber(context.Background())
    if err != nil {
        log.Fatal(err)
    }
 
    fmt.Printf("Block Hash: %s\n", result.Hash)
    fmt.Printf("Block Number: %d\n", result.Number)
}