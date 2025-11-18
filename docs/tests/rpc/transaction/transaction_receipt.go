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
 
    // Parse transaction hash
    txHashStr := "0x6dda0d2e2716227b87d912d654e1bc8b96441f043c29834e082413ae1320afa"
    txHash, err := new(felt.Felt).SetString(txHashStr)
    if err != nil {
        log.Fatal(err)
    }
 
    // Get transaction receipt
    receipt, err := provider.TransactionReceipt(context.Background(), txHash)
    if err != nil {
        log.Fatal(err)
    }
 
    // Pretty print the result
    receiptJSON, _ := json.MarshalIndent(receipt, "", "  ")
    fmt.Printf("Transaction receipt:\n%s\n", receiptJSON)
}