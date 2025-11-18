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

	fmt.Println("AddInvokeTransaction RPC Method")
	fmt.Println("================================")
	fmt.Println()
	fmt.Println("Method Signature:")
	fmt.Println("  func (provider *Provider) AddInvokeTransaction(")
	fmt.Println("      ctx context.Context,")
	fmt.Println("      invokeTxn *BroadcastInvokeTxnV3,")
	fmt.Println("  ) (AddInvokeTransactionResponse, error)")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - Properly signed BroadcastInvokeTxnV3 transaction")
	fmt.Println("  - Valid account with private key for signing")
	fmt.Println("  - Correct nonce and resource bounds")
	fmt.Println()
	fmt.Println("Practical Usage:")
	fmt.Println("  Use the account package's Execute method which handles")
	fmt.Println("  transaction signing and submission automatically:")
	fmt.Println()
	fmt.Println("  acct, _ := account.NewAccount(client, address, address, ks, 2)")
	fmt.Println("  tx, _ := acct.Execute(ctx, calls, details)")
	fmt.Println()
	fmt.Println("Expected Response Structure:")
	fmt.Println("  AddInvokeTransactionResponse{")
	fmt.Println("    TransactionHash: *felt.Felt  // Hash of submitted transaction")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("Example Output:")
	fmt.Println("  Transaction Hash: 0x1234567890abcdef...")
	fmt.Println()
	fmt.Printf("RPC Provider connected: %T\n", client)
}
