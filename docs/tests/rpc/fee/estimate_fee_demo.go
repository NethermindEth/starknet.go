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
	godotenv.Load(".env")
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("EstimateFee RPC Method")
	fmt.Println("======================")
	fmt.Println()
	fmt.Println("Method Signature:")
	fmt.Println("  func (provider *Provider) EstimateFee(")
	fmt.Println("      ctx context.Context,")
	fmt.Println("      requests []BroadcastTxn,")
	fmt.Println("      simulationFlags []SimulationFlag,")
	fmt.Println("      blockID BlockID,")
	fmt.Println("  ) ([]FeeEstimation, error)")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - Properly constructed BroadcastTxn (signed transaction)")
	fmt.Println("  - Valid account with nonce and resource bounds")
	fmt.Println()
	fmt.Println("Practical Usage:")
	fmt.Println("  Use the account package's EstimateFee method which handles")
	fmt.Println("  transaction construction automatically:")
	fmt.Println()
	fmt.Println("  acct, _ := account.NewAccount(client, address, address, ks, 2)")
	fmt.Println("  feeEstimate, _ := acct.EstimateFee(ctx, calls, details)")
	fmt.Println()
	fmt.Println("Expected Response Structure:")
	fmt.Println("  []FeeEstimation{")
	fmt.Println("    {")
	fmt.Println("      L1GasConsumed:    *felt.Felt")
	fmt.Println("      L1GasPrice:       *felt.Felt")
	fmt.Println("      L2GasConsumed:    *felt.Felt")
	fmt.Println("      L2GasPrice:       *felt.Felt")
	fmt.Println("      L1DataGasConsumed: *felt.Felt")
	fmt.Println("      L1DataGasPrice:   *felt.Felt")
	fmt.Println("      OverallFee:       *felt.Felt")
	fmt.Println("      Unit:             \"FRI\"")
	fmt.Println("    }")
	fmt.Println("  }")
	fmt.Println()
	fmt.Printf("RPC Provider: %s\n", client)
}
