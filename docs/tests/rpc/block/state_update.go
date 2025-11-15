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
 
    // Get state update for latest block
    blockID := rpc.WithBlockTag("latest")
    stateUpdate, err := provider.StateUpdate(context.Background(), blockID)
    if err != nil {
        log.Fatal(err)
    }
 
    fmt.Printf("Block Hash: %s\n", stateUpdate.BlockHash)
    fmt.Printf("New Root: %s\n", stateUpdate.NewRoot)
    fmt.Printf("Old Root: %s\n", stateUpdate.OldRoot)
    fmt.Printf("Storage Diffs: %d contracts\n", len(stateUpdate.StateDiff.StorageDiffs))
    fmt.Printf("Declared Classes: %d\n", len(stateUpdate.StateDiff.DeclaredClasses))
    fmt.Printf("Deployed Contracts: %d\n", len(stateUpdate.StateDiff.DeployedContracts))
    fmt.Printf("Nonce Updates: %d\n", len(stateUpdate.StateDiff.Nonces))
}