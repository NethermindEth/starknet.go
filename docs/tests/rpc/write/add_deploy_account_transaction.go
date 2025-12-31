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
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	rpcURL := os.Getenv("STARKNET_RPC_URL")
	privateKeyStr := os.Getenv("ACCOUNT_PRIVATE_KEY")
	publicKeyStr := os.Getenv("ACCOUNT_PUBLIC_KEY")
	classHashStr := os.Getenv("ACCOUNT_CLASS_HASH")
	saltStr := os.Getenv("ACCOUNT_SALT")
	addressStr := os.Getenv("ACCOUNT_ADDRESS")

	if rpcURL == "" || privateKeyStr == "" || publicKeyStr == "" || classHashStr == "" || saltStr == "" || addressStr == "" {
		log.Fatal("Missing required environment variables")
	}

	// Parse credentials
	privateKey, err := new(felt.Felt).SetString(privateKeyStr)
	if err != nil {
		log.Fatal("Invalid private key:", err)
	}

	publicKey, err := new(felt.Felt).SetString(publicKeyStr)
	if err != nil {
		log.Fatal("Invalid public key:", err)
	}

	classHash, err := new(felt.Felt).SetString(classHashStr)
	if err != nil {
		log.Fatal("Invalid class hash:", err)
	}

	salt, err := new(felt.Felt).SetString(saltStr)
	if err != nil {
		log.Fatal("Invalid salt:", err)
	}

	address, err := new(felt.Felt).SetString(addressStr)
	if err != nil {
		log.Fatal("Invalid address:", err)
	}

	fmt.Printf("Account Address: %s\n", address.String())
	fmt.Printf("RPC URL: %s\n\n", rpcURL)

	// Connect to RPC provider
	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to RPC provider:", err)
	}

	// Check if account is already deployed
	fmt.Println("Checking if account is already deployed...")
	classHashResponse, err := client.ClassHashAt(ctx, rpc.BlockID{Tag: "latest"}, address)
	if err == nil && classHashResponse != nil {
		fmt.Println("Account is already deployed!")
		fmt.Printf("Class Hash: %s\n", classHashResponse.String())
		return
	}

	// Check STRK balance
	fmt.Println("Checking STRK balance...")
	strkBalance, err := getStrkBalance(ctx, client, address)
	if err != nil {
		log.Printf("Warning: Could not check STRK balance: %v\n", err)
	} else {
		// Convert to STRK (18 decimals)
		strkAmount := new(big.Float).Quo(
			new(big.Float).SetInt(strkBalance),
			new(big.Float).SetFloat64(1e18),
		)
		fmt.Printf("STRK Balance: %s STRK\n", strkAmount.Text('f', 4))
		if strkBalance.Cmp(big.NewInt(0)) == 0 {
			fmt.Println("\nWARNING: Account has 0 STRK balance!")
			fmt.Println("Please fund your account first:")
			fmt.Printf("Address: %s\n", address.String())
			fmt.Println("Faucet: https://starknet-faucet.vercel.app/")
			return
		}
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("Failed to get chain ID:", err)
	}
	fmt.Printf("Chain ID: %s\n\n", chainID)

	// Create account manager (undeployed)
	ks := account.NewMemKeystore()
	ks.Put(publicKey.String(), privateKey.BigInt(new(big.Int)))

	// Constructor calldata for OZ account is just the public key
	constructorCalldata := []*felt.Felt{publicKey}

	// Prepare deploy account transaction
	fmt.Println("Deploying account...")

	// Create account controller
	accnt, err := account.NewAccount(client, address, publicKey.String(), ks, 2)
	if err != nil {
		log.Fatal("Failed to create account controller:", err)
	}

	// Build and estimate deploy account transaction (V3 with STRK fees)
	deployTx, precomputedAddr, err := accnt.BuildAndEstimateDeployAccountTxn(
		ctx,
		salt,
		classHash,
		constructorCalldata,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to build deploy account transaction:", err)
	}

	if precomputedAddr.String() != address.String() {
		log.Fatalf("Address mismatch! Expected %s, got %s", address.String(), precomputedAddr.String())
	}

	fmt.Printf("Estimated gas: L1Gas=%d, L2Gas=%d\n",
		deployTx.ResourceBounds.L1Gas.MaxAmount,
		deployTx.ResourceBounds.L2Gas.MaxAmount)

	// Sign the transaction
	err = accnt.SignDeployAccountTransaction(ctx, &rpc.DeployAccountTxnV3{
		Type:                deployTx.Type,
		ClassHash:           deployTx.ClassHash,
		ContractAddressSalt: deployTx.ContractAddressSalt,
		ConstructorCalldata: deployTx.ConstructorCalldata,
		Version:             deployTx.Version,
		Signature:           deployTx.Signature,
		Nonce:               deployTx.Nonce,
		ResourceBounds:      deployTx.ResourceBounds,
		Tip:                 deployTx.Tip,
		PayMasterData:       deployTx.PayMasterData,
		NonceDataMode:       deployTx.NonceDataMode,
		FeeMode:             deployTx.FeeMode,
	}, address)
	if err != nil {
		log.Fatal("Failed to sign transaction:", err)
	}

	// Send the transaction
	resp, err := client.AddDeployAccountTransaction(ctx, deployTx)
	if err != nil {
		log.Fatal("Failed to deploy account:", err)
	}

	fmt.Printf("Deploy transaction sent!\n")
	fmt.Printf("Transaction Hash: %s\n", resp.Hash.String())

	// Wait for transaction confirmation
	fmt.Println("\nWaiting for transaction confirmation...")
	receipt, err := waitForTransaction(ctx, client, resp.Hash)
	if err != nil {
		log.Fatal("Transaction failed:", err)
	}

	fmt.Printf("Account deployed successfully!\n")
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Status: %s\n\n", receipt.FinalityStatus)

	fmt.Println("Your account is now ready to use!")
	fmt.Printf("Address: %s\n", address.String())
}

func getStrkBalance(ctx context.Context, client *rpc.Provider, address *felt.Felt) (*big.Int, error) {
	// STRK contract address on Sepolia
	strkContract, _ := new(felt.Felt).SetString("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")

	// Call balanceOf
	result, err := client.Call(ctx, rpc.FunctionCall{
		ContractAddress:    strkContract,
		EntryPointSelector: utils.GetSelectorFromNameFelt("balanceOf"),
		Calldata:           []*felt.Felt{address},
	}, rpc.BlockID{Tag: "latest"})

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return big.NewInt(0), nil
	}

	// balanceOf returns Uint256 (low, high)
	if len(result) >= 2 {
		low := result[0].BigInt(new(big.Int))
		high := result[1].BigInt(new(big.Int))
		// balance = low + high * 2^128
		balance := new(big.Int).Lsh(high, 128)
		balance.Add(balance, low)
		return balance, nil
	}

	return result[0].BigInt(new(big.Int)), nil
}

// waitForTransaction waits for a transaction to be confirmed on the network
func waitForTransaction(ctx context.Context, client *rpc.Provider, txHash *felt.Felt) (*rpc.TransactionReceiptWithBlockInfo, error) {
	for i := 0; i < 60; i++ {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			if receipt.FinalityStatus == rpc.TxnFinalityStatusAcceptedOnL2 ||
				receipt.FinalityStatus == rpc.TxnFinalityStatusAcceptedOnL1 {
				return receipt, nil
			}
		}

		time.Sleep(5 * time.Second)
		fmt.Print(".")
	}

	return nil, fmt.Errorf("transaction confirmation timeout")
}
