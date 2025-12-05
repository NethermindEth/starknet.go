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
 
	ctx := context.Background()
 
	// Transaction hash to trace
	txHash, err := new(felt.Felt).SetString("0x0487bcdd78ea9f9512ba1c772114f851dc1dc057a23b27d6f2ffe2c84f97140d")
	if err != nil {
		log.Fatal(err)
	}
 
	// Get transaction trace
	trace, err := provider.TraceTransaction(ctx, txHash)
	if err != nil {
		log.Fatal(err)
	}
 
	// Print full trace as JSON
	traceJSON, _ := json.MarshalIndent(trace, "", "  ")
	fmt.Printf("Transaction Trace:\n%s\n", traceJSON)
 
	// Type assert to InvokeTxnTrace for detailed analysis
	if invokeTrace, ok := trace.(rpc.InvokeTxnTrace); ok {
		fmt.Printf("\nTransaction Type: %s\n", invokeTrace.Type)
		fmt.Printf("Total L2 Gas: %d\n", invokeTrace.ExecutionResources.L2Gas)
		fmt.Printf("Total L1 Data Gas: %d\n", invokeTrace.ExecutionResources.L1DataGas)
 
		// Analyze execute invocation
		fmt.Printf("\nExecute Invocation:\n")
		fmt.Printf("  Contract: %s\n", invokeTrace.ExecuteInvocation.ContractAddress)
		fmt.Printf("  Nested Calls: %d\n", len(invokeTrace.ExecuteInvocation.NestedCalls))
		fmt.Printf("  Events Emitted: %d\n", len(invokeTrace.ExecuteInvocation.InvocationEvents))
	}
}