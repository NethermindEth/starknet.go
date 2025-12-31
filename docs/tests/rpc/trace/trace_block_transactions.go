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
 
	// Trace all transactions in a specific block
	blockID := rpc.WithBlockNumber(99433)
	traces, err := provider.TraceBlockTransactions(ctx, blockID)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Printf("Block has %d transactions\n\n", len(traces))
 
	// Analyze each transaction trace
	for i, trace := range traces {
		fmt.Printf("Transaction %d:\n", i+1)
		fmt.Printf("  Hash: %s\n", trace.TxnHash)
 
		// Type assert to analyze specific trace types
		if invokeTrace, ok := trace.TraceRoot.(rpc.InvokeTxnTrace); ok {
			fmt.Printf("  Type: %s\n", invokeTrace.Type)
			fmt.Printf("  Total L2 Gas: %d\n", invokeTrace.ExecutionResources.L2Gas)
 
			// Count nested calls
			nestedCallCount := len(invokeTrace.ExecuteInvocation.NestedCalls)
			fmt.Printf("  Nested Calls: %d\n", nestedCallCount)
 
			// Check for state changes
			if invokeTrace.StateDiff != nil {
				fmt.Printf("  Storage Updates: %d contracts\n", len(invokeTrace.StateDiff.StorageDiffs))
				fmt.Printf("  Nonce Updates: %d accounts\n", len(invokeTrace.StateDiff.Nonces))
			}
		}
 
		fmt.Println()
	}
}