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
 
    // Example L1 transaction hash that sent messages to Starknet
    // This should be an Ethereum transaction hash from StarkGate bridge or similar
    l1TxHashStr := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    l1TxHash := rpc.NumAsHex(l1TxHashStr)
 
    // Get messages status
    messagesStatus, err := provider.MessagesStatus(ctx, l1TxHash)
    if err != nil {
        log.Fatal(err)
    }
 
    fmt.Printf("L1 Transaction: %s\n", l1TxHashStr)
    fmt.Printf("Messages Count: %d\n\n", len(messagesStatus))
 
    for i, status := range messagesStatus {
        fmt.Printf("Message %d:\n", i+1)
        fmt.Printf("  L2 Transaction Hash: %s\n", status.Hash)
        fmt.Printf("  Finality Status: %s\n", status.FinalityStatus)
        fmt.Printf("  Execution Status: %s\n", status.ExecutionStatus)
        if status.FailureReason != "" {
            fmt.Printf("  Failure Reason: %s\n", status.FailureReason)
        }
        fmt.Println()
    }
}