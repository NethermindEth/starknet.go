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
 
    ctx := context.Background()
 
    // Get synchronization status
    syncStatus, err := provider.Syncing(ctx)
    if err != nil {
        log.Fatal(err)
    }
 
    // Check if node is syncing
    if syncStatus.IsSyncing {
        fmt.Printf("Node is syncing:\n")
        fmt.Printf("  Starting Block: %d (Hash: %s)\n",
            syncStatus.StartingBlockNum, syncStatus.StartingBlockHash)
        fmt.Printf("  Current Block: %d (Hash: %s)\n",
            syncStatus.CurrentBlockNum, syncStatus.CurrentBlockHash)
        fmt.Printf("  Highest Block: %d (Hash: %s)\n",
            syncStatus.HighestBlockNum, syncStatus.HighestBlockHash)
 
        // Calculate sync progress
        if syncStatus.HighestBlockNum > syncStatus.StartingBlockNum {
            progress := float64(syncStatus.CurrentBlockNum-syncStatus.StartingBlockNum) /
                       float64(syncStatus.HighestBlockNum-syncStatus.StartingBlockNum) * 100
            fmt.Printf("  Progress: %.2f%%\n", progress)
        }
    } else {
        fmt.Println("Node is fully synchronized")
    }
}