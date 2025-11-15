package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("STARKNET_RPC_URL")
	testAccountAddress := os.Getenv("TEST_ACCOUNT_ADDRESS")
	testAccountPrivateKey := os.Getenv("TEST_ACCOUNT_PRIVATE_KEY")
	testAccountPublicKey := os.Getenv("TEST_ACCOUNT_PUBLIC_KEY")

	if rpcURL == "" || testAccountAddress == "" || testAccountPrivateKey == "" || testAccountPublicKey == "" {
		log.Fatal("Required environment variables not set")
	}

	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Parse account details
	accountAddress, _ := new(felt.Felt).SetString(testAccountAddress)
	privateKeyFelt, _ := new(felt.Felt).SetString(testAccountPrivateKey)
	privateKey := privateKeyFelt.BigInt(new(big.Int))

	// Create account (Cairo v2)
	ks := account.SetNewMemKeystore(testAccountPublicKey, privateKey)
	acc, err := account.NewAccount(provider, accountAddress, testAccountPublicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== TESTING WaitForTransactionReceipt ===")
	fmt.Println()

	// First, send a transaction
	ethContract, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	functionCall := rpc.InvokeFunctionCall{
		ContractAddress: ethContract,
		FunctionName:    "transfer",
		CallData: []*felt.Felt{
			accountAddress,              // recipient (self)
			new(felt.Felt).SetUint64(1), // amount (u256 low)
			new(felt.Felt).SetUint64(0), // amount (u256 high)
		},
	}

	fmt.Println("Sending a transaction...")
	response, err := acc.BuildAndSendInvokeTxn(ctx, []rpc.InvokeFunctionCall{functionCall}, nil)
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
	}

	fmt.Printf("Transaction Hash: %s\n", response.Hash)
	fmt.Println()

	// Now wait for the transaction receipt
	fmt.Println("Waiting for transaction receipt (polling every 2 seconds)...")
	fmt.Println()

	receipt, err := acc.WaitForTransactionReceipt(
		ctx,
		response.Hash,
		2*time.Second, // Poll every 2 seconds
	)
	if err != nil {
		log.Fatalf("Error waiting for receipt: %v", err)
	}

	fmt.Println("✅ Transaction confirmed!")
	fmt.Println()
	fmt.Printf("Block Hash: %s\n", receipt.BlockHash)
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Execution Status: %s\n", receipt.ExecutionStatus)
	fmt.Printf("Finality Status: %s\n", receipt.FinalityStatus)
	fmt.Println()
	fmt.Printf("View on Voyager: https://sepolia.voyager.online/tx/%s\n", response.Hash)
}
