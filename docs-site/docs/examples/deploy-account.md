---
sidebar_position: 2
---

# Deploy Account Example

This example demonstrates how to deploy a new account contract on StarkNet. In StarkNet, accounts are smart contracts that must be deployed before they can be used.

## Overview

The deploy account example shows how to:

1. Generate a new key pair
2. Precompute the account address
3. Fund the account with ETH
4. Deploy the account contract

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the `.env` file
3. Fund the precomputed account address with ETH using a faucet (e.g., [StarkNet Sepolia Faucet](https://faucet.goerli.starknet.io/))

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load variables from '.env' file
	err := godotenv.Load("../.env")
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err))
	}

	// Get the RPC provider URL from environment variables
	providerURL := os.Getenv("RPC_PROVIDER_URL")
	if providerURL == "" {
		panic("RPC_PROVIDER_URL environment variable is not set")
	}

	// Initialize connection to provider
	provider, err := rpc.NewProvider(providerURL)
	if err != nil {
		panic(fmt.Sprintf("Error initializing provider: %s", err))
	}

	// Generate a new key pair
	starkCurve := curve.NewStarkCurve()
	privateKey, err := utils.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("Error generating private key: %s", err))
	}

	publicKey, err := starkCurve.GetPublicKey(privateKey)
	if err != nil {
		panic(fmt.Sprintf("Error getting public key: %s", err))
	}

	fmt.Printf("Generated private key: 0x%s\n", privateKey.Text(16))
	fmt.Printf("Generated public key: 0x%s\n", publicKey.Text(16))

	// Class hash of the account contract (OpenZeppelin account contract for Cairo 1.0)
	// This is the class hash of the OpenZeppelin account contract on Sepolia
	classHash := "0x4d07e40e93398ed3c76981e72dd1fd22557a78ce36c0515f679e27f0bb5bc5f"

	// Constructor calldata (public key)
	constructorCalldata := []string{publicKey.Text(16)}

	// Salt for address generation (can be any value, using 0 for simplicity)
	salt := big.NewInt(0)

	// Precompute the account address
	accountAddress, err := account.PrecomputeAddress(
		provider,
		publicKey.Text(16),
		salt,
		classHash,
		constructorCalldata,
	)
	if err != nil {
		panic(fmt.Sprintf("Error precomputing address: %s", err))
	}

	fmt.Printf("Precomputed account address: 0x%s\n", accountAddress.Text(16))
	fmt.Println("\nIMPORTANT: Fund this address with ETH before continuing!")
	fmt.Println("You can use the StarkNet Sepolia Faucet: https://faucet.goerli.starknet.io/")
	fmt.Println("Press Enter to continue after funding the account...")
	fmt.Scanln() // Wait for user input

	// Create a keystore for the account
	keystore := account.NewMemKeystore()
	if err := keystore.Put(publicKey.Text(16), privateKey); err != nil {
		panic(fmt.Sprintf("Error putting key in keystore: %s", err))
	}

	// Create a new account instance
	acc, err := account.NewAccount(
		provider,
		accountAddress.Text(16),
		publicKey.Text(16),
		keystore,
		1, // Cairo version (1 for Cairo 1.0)
	)
	if err != nil {
		panic(fmt.Sprintf("Error creating account: %s", err))
	}

	// Deploy the account
	deployTx, err := acc.Deploy(context.Background(), salt, classHash, constructorCalldata)
	if err != nil {
		panic(fmt.Sprintf("Error deploying account: %s", err))
	}

	fmt.Printf("Deploy transaction hash: 0x%s\n", deployTx.TransactionHash.Text(16))

	// Wait for the transaction to be accepted
	fmt.Println("Waiting for the transaction to be accepted...")
	receipt, err := provider.WaitForTransaction(context.Background(), deployTx.TransactionHash, 5, 10)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for transaction: %s", err))
	}

	fmt.Printf("Transaction status: %s\n", receipt.Status)
	fmt.Printf("Account deployed successfully at address: 0x%s\n", accountAddress.Text(16))
	fmt.Println("\nIMPORTANT: Save your private and public keys securely!")
	fmt.Printf("Private key: 0x%s\n", privateKey.Text(16))
	fmt.Printf("Public key: 0x%s\n", publicKey.Text(16))
}
```

## Running the Example

To run this example:

1. Make sure you are in the "deployAccount" directory
2. Execute `go run main.go`
3. When prompted, fund the precomputed account address with ETH using a faucet
4. Press Enter to continue with the deployment

## Expected Output

```
Generated private key: 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
Generated public key: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890

Precomputed account address: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

IMPORTANT: Fund this address with ETH before continuing!
You can use the StarkNet Sepolia Faucet: https://faucet.goerli.starknet.io/
Press Enter to continue after funding the account...

Deploy transaction hash: 0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
Waiting for the transaction to be accepted...
Transaction status: ACCEPTED_ON_L2
Account deployed successfully at address: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

IMPORTANT: Save your private and public keys securely!
Private key: 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
Public key: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
```

## Key Concepts

### Account Contracts

In StarkNet, accounts are smart contracts that must be deployed before they can be used. The account contract handles transaction signing and verification.

### Precomputing Account Addresses

StarkNet allows you to precompute the address of an account before deploying it. This is useful for funding the account before deployment.

```go
accountAddress, err := account.PrecomputeAddress(
    provider,
    publicKey.Text(16),
    salt,
    classHash,
    constructorCalldata,
)
```

### Account Deployment

To deploy an account, you need:

1. A provider connected to a StarkNet node
2. A key pair (private and public keys)
3. A class hash for the account contract
4. Constructor calldata (typically just the public key)
5. A salt for address generation

The deployment process involves:

```go
deployTx, err := acc.Deploy(context.Background(), salt, classHash, constructorCalldata)
```

### Waiting for Transactions

StarkNet transactions can take some time to be accepted. You can wait for a transaction to be accepted using:

```go
receipt, err := provider.WaitForTransaction(context.Background(), txHash, retryInterval, maxRetries)
```

## Next Steps

After deploying an account, you can:

- Learn how to [Invoke Transactions](./invoke-transaction.md) to interact with contracts
- Explore how to [Declare Contracts](./declare-transaction.md) on StarkNet
- Understand how to [Deploy Contracts](./deploy-contract-udc.md) using the Universal Deployer Contract
