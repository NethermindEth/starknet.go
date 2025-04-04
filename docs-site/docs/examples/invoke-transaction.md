---
sidebar_position: 3
---

# Invoke Transaction Example

This example demonstrates how to invoke a transaction on StarkNet, which allows you to call functions on deployed contracts and modify their state.

## Overview

The invoke transaction example shows how to:

1. Initialize a connection to a StarkNet provider
2. Load an existing account
3. Create a function call to transfer ERC20 tokens
4. Execute the transaction and wait for it to be accepted

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the `.env` file
3. Uncomment and assign your account address to the `ACCOUNT_ADDRESS` variable in the `.env` file
4. Uncomment and assign your public key to the `PUBLIC_KEY` variable in the `.env` file
5. Uncomment and assign your private key to the `PRIVATE_KEY` variable in the `.env` file

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
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

	// Get the RPC provider URL and account details from environment variables
	providerURL := os.Getenv("RPC_PROVIDER_URL")
	if providerURL == "" {
		panic("RPC_PROVIDER_URL environment variable is not set")
	}

	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	if accountAddress == "" {
		panic("ACCOUNT_ADDRESS environment variable is not set")
	}

	publicKey := os.Getenv("PUBLIC_KEY")
	if publicKey == "" {
		panic("PUBLIC_KEY environment variable is not set")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		panic("PRIVATE_KEY environment variable is not set")
	}

	// Initialize connection to provider
	provider, err := rpc.NewProvider(providerURL)
	if err != nil {
		panic(fmt.Sprintf("Error initializing provider: %s", err))
	}

	// Convert private key from hex string to big.Int
	privateKeyBigInt, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Error converting private key to big.Int")
	}

	// Create a keystore for the account
	keystore := account.NewMemKeystore()
	if err := keystore.Put(publicKey, privateKeyBigInt); err != nil {
		panic(fmt.Sprintf("Error putting key in keystore: %s", err))
	}

	// Create the account instance
	acc, err := account.NewAccount(
		provider,
		accountAddress,
		publicKey,
		keystore,
		1, // Cairo version (1 for Cairo 1.0)
	)
	if err != nil {
		panic(fmt.Sprintf("Error creating account: %s", err))
	}

	// ERC20 contract address (STRK token on Sepolia)
	contractAddress := "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"

	// Recipient address (replace with the address you want to send tokens to)
	recipientAddress := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	// Amount to transfer (1 token with 18 decimals = 1 * 10^18)
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	// Convert amount to hex string
	amountHex := utils.GetHexString(amount)

	// Create a function call to transfer tokens
	functionCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromName("transfer"),
		Calldata: []string{
			recipientAddress, // Recipient address
			amountHex,        // Amount (low bits)
			"0",              // Amount (high bits, 0 for amounts < 2^128)
		},
	}

	// Estimate the fee for the transaction
	feeEstimate, err := acc.EstimateFee(context.Background(), []rpc.FunctionCall{functionCall}, nil)
	if err != nil {
		panic(fmt.Sprintf("Error estimating fee: %s", err))
	}

	fmt.Printf("Estimated fee: %s wei\n", feeEstimate.OverallFee.Text(10))

	// Execute the transaction
	tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall}, nil)
	if err != nil {
		panic(fmt.Sprintf("Error executing transaction: %s", err))
	}

	fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))

	// Wait for the transaction to be accepted
	fmt.Println("Waiting for the transaction to be accepted...")
	receipt, err := provider.WaitForTransaction(context.Background(), tx.TransactionHash, 5, 10)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for transaction: %s", err))
	}

	fmt.Printf("Transaction status: %s\n", receipt.Status)
	fmt.Printf("Actual fee: %s wei\n", receipt.ActualFee.Text(10))

	// Check for events
	if len(receipt.Events) > 0 {
		fmt.Println("Events:")
		for i, event := range receipt.Events {
			fmt.Printf("  Event %d: From %s\n", i, event.FromAddress)
			fmt.Printf("    Keys: %v\n", event.Keys)
			fmt.Printf("    Data: %v\n", event.Data)
		}
	}

	fmt.Println("Transaction completed successfully!")
}
```

## Running the Example

To run this example:

1. Make sure you are in the "invoke" directory
2. Execute `go run main.go`

## Expected Output

```
Estimated fee: 12345678901234 wei
Transaction hash: 0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
Waiting for the transaction to be accepted...
Transaction status: ACCEPTED_ON_L2
Actual fee: 12345678901234 wei
Events:
  Event 0: From 0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d
    Keys: [0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9]
    Data: [0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef 0xde0b6b3a7640000 0x0]
Transaction completed successfully!
```

## Key Concepts

### Account Setup

Before invoking transactions, you need to set up an account:

```go
// Create a keystore for the account
keystore := account.NewMemKeystore()
keystore.Put(publicKey, privateKeyBigInt)

// Create the account instance
acc, err := account.NewAccount(
    provider,
    accountAddress,
    publicKey,
    keystore,
    1, // Cairo version (1 for Cairo 1.0)
)
```

### Function Calls

A function call in StarkNet.go is represented by the `FunctionCall` struct:

```go
functionCall := rpc.FunctionCall{
    ContractAddress:    contractAddress,
    EntryPointSelector: utils.GetSelectorFromName("transfer"),
    Calldata: []string{
        recipientAddress, // Recipient address
        amountHex,        // Amount (low bits)
        "0",              // Amount (high bits, 0 for amounts < 2^128)
    },
}
```

### Fee Estimation

Before executing a transaction, you can estimate the fee:

```go
feeEstimate, err := acc.EstimateFee(context.Background(), []rpc.FunctionCall{functionCall}, nil)
```

### Transaction Execution

To execute a transaction:

```go
tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall}, nil)
```

### Waiting for Transactions

StarkNet transactions can take some time to be accepted. You can wait for a transaction to be accepted using:

```go
receipt, err := provider.WaitForTransaction(context.Background(), txHash, retryInterval, maxRetries)
```

### Transaction Receipt

The transaction receipt contains information about the executed transaction:

- Status: The status of the transaction (e.g., ACCEPTED_ON_L2)
- Actual fee: The actual fee paid for the transaction
- Events: Events emitted during the transaction execution

## Multiple Function Calls

You can include multiple function calls in a single transaction:

```go
functionCall1 := rpc.FunctionCall{
    ContractAddress:    contractAddress1,
    EntryPointSelector: utils.GetSelectorFromName("approve"),
    Calldata: []string{
        spenderAddress,
        amountHex,
        "0",
    },
}

functionCall2 := rpc.FunctionCall{
    ContractAddress:    contractAddress2,
    EntryPointSelector: utils.GetSelectorFromName("transfer"),
    Calldata: []string{
        recipientAddress,
        amountHex,
        "0",
    },
}

// Execute both function calls in a single transaction
tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall1, functionCall2}, nil)
```

## Next Steps

After understanding how to invoke transactions, you can:

- Learn how to [Declare Contracts](./declare-transaction.md) on StarkNet
- Explore how to [Deploy Contracts](./deploy-contract-udc.md) using the Universal Deployer Contract
- Understand how to work with [Typed Data](./typed-data.md) for structured message signing
