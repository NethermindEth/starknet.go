package main
 
import (
    "context"
    "fmt"
    "log"
    "os"
 
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
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
 
    // Example account address
    accountAddress, err := utils.HexToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
    if err != nil {
        log.Fatal(err)
    }
 
    // Get current nonce
    blockID := rpc.WithBlockTag("latest")
    nonce, err := provider.Nonce(ctx, blockID, accountAddress)
    if err != nil {
        log.Fatal(err)
    }
 
    fmt.Printf("Current nonce: %s\n", nonce)
}