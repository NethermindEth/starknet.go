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

	fmt.Println("AddDeployAccountTransaction RPC Method")
	fmt.Println("=======================================")
	fmt.Println()
	fmt.Println("Method Signature:")
	fmt.Println("  func (provider *Provider) AddDeployAccountTransaction(")
	fmt.Println("      ctx context.Context,")
	fmt.Println("      deployAccountTransaction *BroadcastDeployAccountTxnV3,")
	fmt.Println("  ) (AddDeployAccountTransactionResponse, error)")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - Account class hash (pre-declared)")
	fmt.Println("  - Constructor calldata")
	fmt.Println("  - Properly signed transaction")
	fmt.Println("  - Sufficient funds at computed address for fees")
	fmt.Println()
	fmt.Println("Practical Usage:")
	fmt.Println("  1. Pre-fund the computed account address")
	fmt.Println("  2. Use account.NewAccount to deploy:")
	fmt.Println()
	fmt.Println("  ks := account.NewMemKeystore()")
	fmt.Println("  ks.Put(address, privateKey)")
	fmt.Println("  acct, _ := account.NewAccount(client, address, classHash, ks, 2)")
	fmt.Println("  // Account deployment happens automatically on first transaction")
	fmt.Println()
	fmt.Println("Expected Response Structure:")
	fmt.Println("  AddDeployAccountTransactionResponse{")
	fmt.Println("    TransactionHash: *felt.Felt  // Hash of deploy transaction")
	fmt.Println("    ContractAddress: *felt.Felt  // Deployed account address")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("Example Output:")
	fmt.Println("  Transaction Hash: 0x1234567890abcdef...")
	fmt.Println("  Contract Address: 0x0517c64b48079568a30...")
	fmt.Println()
	fmt.Printf("RPC Provider connected: %T\n", client)
}
