package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

// OpenZeppelin Account Class Hash in Sepolia
var predeployedClassHash = "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	rpcProviderURL := os.Getenv("STARKNET_RPC_URL")
	if rpcProviderURL == "" {
		panic("STARKNET_RPC_URL environment variable is not set")
	}

	ctx := context.Background()

	// Initialise the client
	client, err := rpc.NewProvider(ctx, rpcProviderURL)
	if err != nil {
		panic(err)
	}

	// Get random keys for the new account
	ks, pub, privKey := account.GetRandomKeys()
	fmt.Printf("Generated public key: %v\n", pub)
	fmt.Printf("Generated private key: %v\n", privKey)

	// Set up the account (using pub as temporary address since we'll compute the real one)
	accnt, err := account.NewAccount(client, pub, pub.String(), ks, account.CairoV2)
	if err != nil {
		panic(err)
	}

	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		panic(err)
	}

	// Build and estimate fees for the deploy account transaction
	deployAccountTxn, precomputedAddress, err := accnt.BuildAndEstimateDeployAccountTxn(
		ctx,
		pub,
		classHash,
		[]*felt.Felt{pub},
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Precomputed address: %s\n", precomputedAddress.String())

	// Save the generated credentials to the .env file
	if saveErr := saveCredentialsToEnv(privKey, pub, precomputedAddress, classHash); saveErr != nil {
		fmt.Printf("Warning: Failed to save credentials to .env: %v\n", saveErr)
	} else {
		fmt.Println("Credentials saved to .env file")
	}

	// Calculate fee in STRK
	overallFee, err := utils.ResBoundsMapToOverallFee(
		deployAccountTxn.ResourceBounds,
		1,
		deployAccountTxn.Tip,
	)
	if err != nil {
		panic(err)
	}
	feeInSTRK := utils.FRIToSTRK(overallFee)

	// Wait for user to fund the account
	var input string
	fmt.Println("\nThe account needs STRK to deploy.")
	fmt.Printf("Send approximately %f STRK to: %s\n", feeInSTRK, precomputedAddress.String())
	fmt.Println("You can use the Starknet faucet: https://starknet-faucet.vercel.app/")
	fmt.Println("\nPress Enter after funding the account...")
	_, _ = fmt.Scanln(&input)

	// Send transaction to the network
	resp, err := accnt.SendTransaction(ctx, deployAccountTxn)
	if err != nil {
		fmt.Println("Error sending transaction:")
		panic(err)
	}

	fmt.Println("Deploy transaction submitted!")
	fmt.Printf("Transaction hash: %s\n", resp.Hash.String())
	fmt.Printf("Contract address: %s\n", resp.ContractAddress.String())

	// Wait for transaction confirmation
	fmt.Print("Waiting for confirmation")
	receipt, err := waitForTransaction(ctx, client, resp.Hash)
	if err != nil {
		fmt.Printf("\nWarning: Could not confirm transaction: %v\n", err)
		fmt.Println("Check the transaction status on Voyager or Starkscan.")
	} else {
		fmt.Printf("\n\nAccount deployed successfully!\n")
		fmt.Printf("Block number: %d\n", receipt.BlockNumber)
		fmt.Printf("Status: %s\n", receipt.FinalityStatus)
	}
}

// saveCredentialsToEnv saves the generated account credentials to the .env file
func saveCredentialsToEnv(privKey, pubKey, address, classHash *felt.Felt) error {
	// Load existing env vars
	envMap := make(map[string]string)
	if data, err := os.ReadFile(".env"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}
	}

	// Add new credentials
	envMap["ACCOUNT_PRIVATE_KEY"] = privKey.String()
	envMap["ACCOUNT_PUBLIC_KEY"] = pubKey.String()
	envMap["ACCOUNT_ADDRESS"] = address.String()
	envMap["ACCOUNT_CLASS_HASH"] = classHash.String()

	// Write back to file
	var content strings.Builder
	for key, value := range envMap {
		content.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

	return os.WriteFile(".env", []byte(content.String()), 0o644)
}

// waitForTransaction polls the network until the transaction is confirmed or times out
func waitForTransaction(
	ctx context.Context,
	client *rpc.Provider,
	txHash *felt.Felt,
) (*rpc.TransactionReceiptWithBlockInfo, error) {
	timeout := 60 // 5 seconds * 60 = 5 minutes

	for range timeout {
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

	return nil, errors.New("transaction confirmation timeout")
}
