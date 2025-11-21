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
 
	// Get account address and current nonce
	accountAddress, _ := new(felt.Felt).SetString("0x36d67ab362562a97f9fba8a1051cf8e37ff1a1449530fb9f1f0e32ac2da7d06")
	nonce, err := provider.Nonce(ctx, rpc.WithBlockTag(rpc.BlockTagLatest), accountAddress)
	if err != nil {
		log.Fatal(err)
	}
 
	// Build an invoke transaction
	contractAddress, _ := new(felt.Felt).SetString("0x669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54")
	entrypoint, _ := new(felt.Felt).SetString("0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354")
 
	invokeTx := rpc.BroadcastInvokeTxnV3{
		Type:          rpc.TransactionTypeInvoke,
		Version:       rpc.TransactionV3,
		Nonce:         nonce,
		SenderAddress: accountAddress,
		Signature:     []*felt.Felt{},
		Calldata: []*felt.Felt{
			new(felt.Felt).SetUint64(1),
			contractAddress,
			entrypoint,
			new(felt.Felt).SetUint64(2),
			new(felt.Felt).SetUint64(256),
			new(felt.Felt).SetUint64(0),
		},
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x1e0",
				MaxPricePerUnit: "0x922",
			},
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0xfbfdefe2186",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x16eea0",
				MaxPricePerUnit: "0x1830e58f7",
			},
		},
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}
 
	// Simulate with SKIP_VALIDATE and SKIP_FEE_CHARGE flags
	// SKIP_VALIDATE skips signature and nonce validation
	// SKIP_FEE_CHARGE skips fee charging, useful when resource bounds may be outdated
	simulationFlags := []rpc.SimulationFlag{rpc.SkipValidate, rpc.SkipFeeCharge}
	results, err := provider.SimulateTransactions(
		ctx,
		rpc.WithBlockTag(rpc.BlockTagLatest),
		[]rpc.BroadcastTxn{&invokeTx},
		simulationFlags,
	)
	if err != nil {
		log.Fatal(err)
	}
 
	// Analyze simulation results
	for i, result := range results {
		fmt.Printf("\nTransaction %d:\n", i+1)
 
		// Check fee estimate
		fmt.Printf("Fee Estimate:\n")
		fmt.Printf("  Overall Fee: %s FRI\n", result.FeeEstimation.OverallFee)
		fmt.Printf("  L2 Gas Consumed: %s\n", result.FeeEstimation.L2GasConsumed)
		fmt.Printf("  L2 Gas Price: %s\n", result.FeeEstimation.L2GasPrice)
 
		// Check for reverts
		if invokeTrace, ok := result.TxnTrace.(rpc.InvokeTxnTrace); ok {
			if invokeTrace.ExecuteInvocation.IsReverted {
				fmt.Printf("  Status: REVERTED\n")
				fmt.Printf("  Revert Reason: %s\n", invokeTrace.ExecuteInvocation.RevertReason)
			} else {
				fmt.Printf("  Status: SUCCESS\n")
				fmt.Printf("  Nested Calls: %d\n", len(invokeTrace.ExecuteInvocation.NestedCalls))
			}
		}
	}
}