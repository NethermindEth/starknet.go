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
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get configuration from environment variables
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	privateKeyStr := os.Getenv("ACCOUNT_PRIVATE_KEY")
	publicKeyStr := os.Getenv("ACCOUNT_PUBLIC_KEY")
	accountAddressStr := os.Getenv("ACCOUNT_ADDRESS")

	if rpcURL == "" || privateKeyStr == "" || publicKeyStr == "" || accountAddressStr == "" {
		log.Fatal("Missing required environment variables")
	}

	// Parse account credentials
	privateKey, err := new(felt.Felt).SetString(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err := new(felt.Felt).SetString(publicKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	accountAddress, err := new(felt.Felt).SetString(accountAddressStr)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize provider
	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create keystore
	ks := account.NewMemKeystore()
	ks.Put(publicKey.String(), privateKey.BigInt(new(big.Int)))

	// Create account controller
	acct, err := account.NewAccount(client, accountAddress, publicKey.String(), ks, 2)
	if err != nil {
		log.Fatal(err)
	}

	// ETH token contract on Sepolia
	ethTokenAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Function selector for "balanceOf"
	balanceOfSelector := utils.GetSelectorFromNameFelt("balanceOf")

	// Build calldata for the invoke transaction
	// For account execute: [num_calls, contract_address, selector, calldata_len, ...calldata]
	calldata := []*felt.Felt{
		new(felt.Felt).SetUint64(1),  // Number of calls
		ethTokenAddress,               // Contract address
		balanceOfSelector,             // Function selector
		new(felt.Felt).SetUint64(1),  // Calldata length
		accountAddress,                // Account address as parameter to balanceOf
	}

	// Get the current nonce
	nonce, err := acct.Nonce(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create the invoke transaction (V3)
	invokeTx := rpc.InvokeTxnV3{
		Type:          rpc.TransactionTypeInvoke,
		Version:       rpc.TransactionV3,
		SenderAddress: accountAddress,
		Nonce:         nonce,
		Calldata:      calldata,
		Signature:     []*felt.Felt{}, // Will be filled by SignInvokeTransaction
		// Resource bounds - use high initial values for estimation
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x186a0",         // 100000
				MaxPricePerUnit: "0x33a937098d80", // High price for estimation
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x186a0",         // 100000
				MaxPricePerUnit: "0x33a937098d80", // High price for estimation
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x186a0",      // 100000
				MaxPricePerUnit: "0x10c388d00", // High price for estimation
			},
		},
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	// Estimate the fee using the provider (skip validation for estimation)
	simFlags := []rpc.SimulationFlag{rpc.SkipValidate}
	feeEstimate, err := client.EstimateFee(ctx, []rpc.BroadcastTxn{invokeTx}, simFlags, rpc.WithBlockTag(rpc.BlockTagLatest))
	if err != nil {
		log.Fatal(err)
	}

	if len(feeEstimate) > 0 {
		// Add buffer to the gas estimate (20% more)
		estimatedL1DataGas := feeEstimate[0].L1DataGasConsumed.Uint64()
		estimatedL2Gas := feeEstimate[0].L2GasConsumed.Uint64()

		// Set resource bounds with buffer
		invokeTx.ResourceBounds = &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", estimatedL1DataGas*12/10)), // 20% buffer
				MaxPricePerUnit: rpc.U128("0x33a937098d80"),                            // Max price in STRK
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", estimatedL1DataGas*12/10)), // 20% buffer
				MaxPricePerUnit: rpc.U128("0x33a937098d80"),                            // Max price in STRK
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", estimatedL2Gas*12/10)), // 20% buffer
				MaxPricePerUnit: rpc.U128("0x10c388d00"),                           // Max price in STRK
			},
		}
	}

	// Sign the transaction
	err = acct.SignInvokeTransaction(ctx, &invokeTx)
	if err != nil {
		log.Fatal(err)
	}

	// Submit the transaction using the provider
	resp, err := client.AddInvokeTransaction(ctx, &invokeTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction submitted successfully!\n")
	fmt.Printf("Transaction Hash: %s\n", resp.Hash.String())

	// Wait for transaction confirmation
	receipt, err := waitForTransaction(ctx, client, resp.Hash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction confirmed!\n")
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Status: %s\n", receipt.FinalityStatus)
	fmt.Printf("Actual Fee: %s\n", receipt.ActualFee.Amount.String())
}

// waitForTransaction waits for a transaction to be confirmed on the network
func waitForTransaction(ctx context.Context, client *rpc.Provider, txHash *felt.Felt) (*rpc.TransactionReceiptWithBlockInfo, error) {
	for i := 0; i < 60; i++ { // Wait up to 5 minutes
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