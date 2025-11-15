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

	fmt.Println("AddDeclareTransaction RPC Method")
	fmt.Println("=================================")
	fmt.Println()
	fmt.Println("Method Signature:")
	fmt.Println("  func (provider *Provider) AddDeclareTransaction(")
	fmt.Println("      ctx context.Context,")
	fmt.Println("      declareTransaction *BroadcastDeclareTxnV3,")
	fmt.Println("  ) (AddDeclareTransactionResponse, error)")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - Compiled contract class (Sierra and CASM)")
	fmt.Println("  - Properly signed BroadcastDeclareTxnV3 transaction")
	fmt.Println("  - Valid account with private key for signing")
	fmt.Println()
	fmt.Println("Practical Usage:")
	fmt.Println("  Use the account package's Declare method:")
	fmt.Println()
	fmt.Println("  acct, _ := account.NewAccount(client, address, address, ks, 2)")
	fmt.Println("  tx, _ := acct.Declare(ctx, classHash, compiledClass, details)")
	fmt.Println()
	fmt.Println("Expected Response Structure:")
	fmt.Println("  AddDeclareTransactionResponse{")
	fmt.Println("    TransactionHash: *felt.Felt  // Hash of declare transaction")
	fmt.Println("    ClassHash:       *felt.Felt  // Hash of declared class")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("Example Output:")
	fmt.Println("  Transaction Hash: 0x1234567890abcdef...")
	fmt.Println("  Class Hash: 0xabcdef1234567890...")
	fmt.Println()
	fmt.Printf("RPC Provider connected: %T\n", client)
}
