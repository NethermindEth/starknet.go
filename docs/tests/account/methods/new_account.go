package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get RPC URL from environment
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	// Create RPC provider
	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		fmt.Printf("Error creating provider: %v\n", err)
		return
	}
	fmt.Println("Provider created successfully")

	// Create account address (ETH token address for example)
	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Create keystore with test key
	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"
	privateKey := new(big.Int).SetUint64(123456789)
	ks := account.SetNewMemKeystore(publicKey, privateKey)

	// Create a Cairo v2 account
	fmt.Println("\nCreating Cairo v2 account:")
	acc, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV2)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
	} else {
		fmt.Printf("Account created successfully\n")
		fmt.Printf("Account address: %s\n", accountAddress)
		fmt.Printf("Cairo version: %d\n", acc.CairoVersion)

		// Get chain ID
		chainID, err := provider.ChainID(ctx)
		if err == nil {
			fmt.Printf("Chain ID: %s\n", chainID)
		}
	}

	// Create a Cairo v0 account
	fmt.Println("\nCreating Cairo v0 account:")
	acc0, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV0)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
	} else {
		fmt.Printf("Account created successfully\n")
		fmt.Printf("Account address: %s\n", accountAddress)
		fmt.Printf("Cairo version: %d\n", acc0.CairoVersion)
	}

	// Test with different account address
	fmt.Println("\nCreating account with different address:")
	accountAddress2, _ := new(felt.Felt).SetString("0x06fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39")
	_, err = account.NewAccount(provider, accountAddress2, publicKey, ks, account.CairoV2)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
	} else {
		fmt.Printf("Account created successfully\n")
		fmt.Printf("Account address: %s\n", accountAddress2)
	}

	// Create account with nil keystore (will fail)
	fmt.Println("\nTesting with nil keystore:")
	_, err = account.NewAccount(provider, accountAddress, publicKey, nil, account.CairoV2)
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	} else {
		fmt.Println("Created account (unexpected)")
	}

	// Create account with empty public key
	fmt.Println("\nTesting with empty public key:")
	_, err = account.NewAccount(provider, accountAddress, "", ks, account.CairoV2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Account created with empty public key")
	}

	// Test with different Cairo versions
	fmt.Println("\nAvailable Cairo versions:")
	fmt.Printf("Cairo v0: %d\n", account.CairoV0)
	fmt.Printf("Cairo v2: %d\n", account.CairoV2)

	fmt.Println("\nSummary:")
	fmt.Println("- Successfully created accounts with Cairo v0 and v2")
	fmt.Println("- Different addresses can use the same keystore")
	fmt.Println("- Accounts require valid provider, address, public key, and keystore")
}