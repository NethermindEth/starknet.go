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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("STARKNET_RPC_URL")
	testAccountAddress := os.Getenv("TEST_ACCOUNT_ADDRESS")
	testAccountPrivateKey := os.Getenv("TEST_ACCOUNT_PRIVATE_KEY")
	testAccountPublicKey := os.Getenv("TEST_ACCOUNT_PUBLIC_KEY")
	testAccountClassHash := os.Getenv("TEST_ACCOUNT_CLASS_HASH")

	if rpcURL == "" || testAccountAddress == "" || testAccountPrivateKey == "" || testAccountPublicKey == "" || testAccountClassHash == "" {
		log.Fatal("Test account credentials not set in .env")
	}

	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Parse values
	accountAddress, _ := new(felt.Felt).SetString(testAccountAddress)
	privateKeyFelt, _ := new(felt.Felt).SetString(testAccountPrivateKey)
	privateKey := privateKeyFelt.BigInt(new(big.Int))
	classHash, _ := new(felt.Felt).SetString(testAccountClassHash)

	// Create keystore
	ks := account.SetNewMemKeystore(testAccountPublicKey, privateKey)

	fmt.Println("=== DEPLOYING CAIRO V2 ACCOUNT ===")
	fmt.Println()
	fmt.Printf("Account Address: %s\n", testAccountAddress)
	fmt.Printf("Class Hash: %s\n", testAccountClassHash)
	fmt.Println()
	fmt.Println("Checking if account is already deployed...")

	// Check if account is already deployed
	_, err = provider.Nonce(ctx, rpc.BlockID{Tag: "latest"}, accountAddress)
	if err == nil {
		fmt.Println()
		fmt.Println("✅ Account is ALREADY DEPLOYED!")
		fmt.Println()
		fmt.Println("You can now use this Cairo v2 account for testing.")
		return
	}

	fmt.Println("Account not deployed yet. Deploying Cairo v2 account...")
	fmt.Println()

	// Salt and constructor calldata
	salt := new(felt.Felt).SetUint64(0)
	publicKeyFelt, _ := new(felt.Felt).SetString(testAccountPublicKey)
	constructorCalldata := []*felt.Felt{publicKeyFelt}

	// Create a temporary account for signing
	tempAcc, err := account.NewAccount(provider, accountAddress, testAccountPublicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatalf("Error creating temp account: %v", err)
	}

	// Build and estimate the deploy account transaction
	fmt.Println("Building and estimating deploy account transaction...")
	deployTx, precomputedAddr, err := tempAcc.BuildAndEstimateDeployAccountTxn(
		ctx,
		salt,
		classHash,
		constructorCalldata,
		nil,
	)
	if err != nil {
		log.Fatalf("Error building deploy tx: %v", err)
	}

	fmt.Printf("Precomputed address matches: %v\n", precomputedAddr.Equal(accountAddress))
	fmt.Println()

	// Send the transaction
	fmt.Println("Sending deploy account transaction...")
	resp, err := provider.AddDeployAccountTransaction(ctx, deployTx)
	if err != nil {
		log.Fatalf("Error sending deploy tx: %v", err)
	}

	fmt.Printf("Deploy transaction sent!\n")
	fmt.Printf("Transaction Hash: %s\n", resp.Hash)
	fmt.Println()
	fmt.Println("Waiting for transaction to be accepted...")

	// Wait for the transaction
	for i := 0; i < 60; i++ {
		time.Sleep(2 * time.Second)
		
		receipt, err := provider.TransactionReceipt(ctx, resp.Hash)
		if err != nil {
			continue
		}

		status := receipt.ExecutionStatus
		if status == rpc.TxnExecutionStatusSUCCEEDED {
			fmt.Println()
			fmt.Println("✅ CAIRO V2 Account DEPLOYED successfully!")
			fmt.Println()
			fmt.Printf("Transaction Hash: %s\n", resp.Hash)
			fmt.Printf("Account Address:  %s\n", testAccountAddress)
			fmt.Println()
			fmt.Println("You can now use this Cairo v2 account for all transaction testing!")
			return
		} else if status == rpc.TxnExecutionStatusREVERTED {
			fmt.Println()
			fmt.Println("❌ Deploy transaction REVERTED")
			return
		}

		fmt.Printf(".")
	}

	fmt.Println()
	fmt.Println("Transaction is taking longer than expected. Check status:")
	fmt.Printf("https://sepolia.voyager.online/tx/%s\n", resp.Hash)
}
