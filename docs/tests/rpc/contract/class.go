package main
 
import (
    "context"
    "encoding/json"
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
 
    // ETH contract class hash on Sepolia
    classHash, err := utils.HexToFelt("0x046ded64ae2dead6448e247234bab192a9c483644395b66f2155f2614e5804b0")
    if err != nil {
        log.Fatal(err)
    }
 
    // Get class definition
    blockID := rpc.WithBlockTag("latest")
    class, err := provider.Class(ctx, blockID, classHash)
    if err != nil {
        log.Fatal(err)
    }
 
    // Print class information (first 500 characters)
    classJSON, _ := json.MarshalIndent(class, "", "  ")
    classStr := string(classJSON)
    if len(classStr) > 500 {
        classStr = classStr[:500] + "..."
    }
    fmt.Printf("Class definition:\n%s\n", classStr)
}