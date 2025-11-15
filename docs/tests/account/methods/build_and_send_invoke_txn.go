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

	// Create keystore and account (Cairo v2)
	ks := account.SetNewMemKeystore(testAccountPublicKey, privateKey)
	acc, err := account.NewAccount(provider, accountAddress, testAccountPublicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== TESTING BuildAndSendInvokeTxn ===")
	fmt.Println()
	fmt.Printf("Account Address: %s\n", testAccountAddress)
	fmt.Println("Account Type: Cairo v2")
	fmt.Println()

	// Make a simple transfer of 1 wei to ourselves  
	// ETH token contract on Starknet Sepolia
	ethContract, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	
	// Transfer: recipient (felt), amount (u256 as low, high)
	functionCall := rpc.InvokeFunctionCall{
		ContractAddress: ethContract,
		FunctionName:    "transfer",
		CallData: []*felt.Felt{
			accountAddress,              // recipient (ourselves)
			new(felt.Felt).SetUint64(1), // amount (u256 low part)
			new(felt.Felt).SetUint64(0), // amount (u256 high part)
		},
	}

	fmt.Println("Sending invoke transaction (transfer 1 wei ETH to self)...")
	fmt.Println()
	fmt.Println("BuildAndSendInvokeTxn handles:")
	fmt.Println("  1. Getting the current nonce")
	fmt.Println("  2. Formatting the calldata")
	fmt.Println("  3. Estimating the fee") 
	fmt.Println("  4. Signing the transaction")
	fmt.Println("  5. Sending it to the network")
	fmt.Println()
	
	// Build and send the invoke transaction
	resp, err := acc.BuildAndSendInvokeTxn(ctx, []rpc.InvokeFunctionCall{functionCall}, nil)
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
	}

	fmt.Println("✅ Transaction sent successfully!")
	fmt.Printf("Transaction Hash: %s\n", resp.Hash)
	fmt.Println()
	fmt.Printf("View on Voyager: https://sepolia.voyager.online/tx/%s\n", resp.Hash)
	fmt.Println()
	
	// Wait a bit and check status
	fmt.Println("Waiting for transaction to be processed...")
	time.Sleep(5 * time.Second)
	
	receipt, err := provider.TransactionReceipt(ctx, resp.Hash)
	if err == nil {
		fmt.Printf("Execution Status: %s\n", receipt.ExecutionStatus)
		fmt.Printf("Finality Status: %s\n", receipt.FinalityStatus)
	} else {
		fmt.Println("Transaction is still pending (check Voyager link above)")
	}
}
