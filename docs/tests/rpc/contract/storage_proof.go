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
 
    // Contract address to get storage proof for
    contractAddr, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
 
    // Storage keys to prove
    storageKey1, _ := new(felt.Felt).SetString("0x1")
    storageKey2, _ := new(felt.Felt).SetString("0x2")
 
    // Create storage proof input
    input := rpc.StorageProofInput{
        BlockID: rpc.WithBlockTag("latest"),
        ContractsStorageKeys: []rpc.ContractStorageKeys{
            {
                ContractAddress: contractAddr,
                StorageKeys:     []*felt.Felt{storageKey1, storageKey2},
            },
        },
    }
 
    // Get storage proof
    proof, err := provider.StorageProof(ctx, input)
    if err != nil {
        log.Fatal(err)
    }
 
    // Access proof data
    fmt.Printf("Contracts Tree Root: %s\n", proof.GlobalRoots.ContractsTreeRoot)
    fmt.Printf("Classes Tree Root: %s\n", proof.GlobalRoots.ClassesTreeRoot)
    fmt.Printf("Block Hash: %s\n", proof.GlobalRoots.BlockHash)
    fmt.Printf("Number of contract proof nodes: %d\n", len(proof.ContractsProof.Nodes))
    fmt.Printf("Number of storage proofs: %d\n", len(proof.ContractsStorageProofs))
}