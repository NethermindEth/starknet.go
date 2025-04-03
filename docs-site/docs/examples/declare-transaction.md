---
sidebar_position: 4
---

# Declare Transaction Example

This example demonstrates how to declare a new contract class on StarkNet. Declaring a contract is the process of registering a contract class on the network before it can be deployed.

## Overview

The declare transaction example shows how to:

1. Initialize a connection to a StarkNet provider
2. Load an existing account
3. Load and parse Sierra and CASM contract files
4. Calculate the compiled class hash
5. Declare the contract class on StarkNet

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the `.env` file
3. Uncomment and assign your account address to the `ACCOUNT_ADDRESS` variable in the `.env` file
4. Uncomment and assign your public key to the `PUBLIC_KEY` variable in the `.env` file
5. Uncomment and assign your private key to the `PRIVATE_KEY` variable in the `.env` file
6. Have Sierra and CASM contract files ready for declaration

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
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

	// Load the contract Sierra file
	sierraContent, err := os.ReadFile("./contract.sierra.json")
	if err != nil {
		panic(fmt.Sprintf("Error reading Sierra file: %s", err))
	}

	// Load the contract CASM file
	casmContent, err := os.ReadFile("./contract.casm.json")
	if err != nil {
		panic(fmt.Sprintf("Error reading CASM file: %s", err))
	}

	// Parse the contract files
	sierra, err := contracts.NewSierraContractClass(sierraContent)
	if err != nil {
		panic(fmt.Sprintf("Error parsing Sierra contract: %s", err))
	}

	casm, err := contracts.NewCasmContractClass(casmContent)
	if err != nil {
		panic(fmt.Sprintf("Error parsing CASM contract: %s", err))
	}

	// Calculate the compiled class hash
	compiledClassHash, err := contracts.ComputeCompiledClassHash(casm)
	if err != nil {
		panic(fmt.Sprintf("Error computing compiled class hash: %s", err))
	}

	fmt.Printf("Compiled class hash: 0x%s\n", compiledClassHash.Text(16))

	// Declare the contract
	declareTx, err := acc.Declare(context.Background(), sierra, compiledClassHash)
	if err != nil {
		panic(fmt.Sprintf("Error declaring contract: %s", err))
	}

	fmt.Printf("Transaction hash: 0x%s\n", declareTx.TransactionHash.Text(16))
	fmt.Printf("Class hash: 0x%s\n", declareTx.ClassHash.Text(16))

	// Wait for the transaction to be accepted
	fmt.Println("Waiting for the transaction to be accepted...")
	receipt, err := provider.WaitForTransaction(context.Background(), declareTx.TransactionHash, 5, 10)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for transaction: %s", err))
	}

	fmt.Printf("Transaction status: %s\n", receipt.Status)
	fmt.Printf("Actual fee: %s wei\n", receipt.ActualFee.Text(10))

	fmt.Println("Contract declared successfully!")
	fmt.Printf("Class hash: 0x%s\n", declareTx.ClassHash.Text(16))
	fmt.Println("You can now deploy instances of this contract using the class hash.")
}
```

## Running the Example

To run this example:

1. Make sure you are in the "simpleDeclare" directory
2. Place your Sierra and CASM contract files in the same directory
3. Execute `go run main.go`

## Expected Output

```
Compiled class hash: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
Transaction hash: 0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
Class hash: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
Waiting for the transaction to be accepted...
Transaction status: ACCEPTED_ON_L2
Actual fee: 12345678901234 wei
Contract declared successfully!
Class hash: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
You can now deploy instances of this contract using the class hash.
```

## Key Concepts

### Contract Classes

In StarkNet, contracts are represented by contract classes. A contract class consists of:

1. **Sierra Contract Class**: The high-level representation of the contract
2. **CASM Contract Class**: The compiled assembly version of the contract

### Contract Declaration

Before deploying a contract, you must declare its class on the network. This is a one-time operation per contract class.

```go
declareTx, err := acc.Declare(context.Background(), sierra, compiledClassHash)
```

### Class Hash and Compiled Class Hash

Each contract class has two important hashes:

1. **Class Hash**: Uniquely identifies the Sierra contract class
2. **Compiled Class Hash**: Uniquely identifies the CASM (compiled) contract class

```go
// Calculate the compiled class hash
compiledClassHash, err := contracts.ComputeCompiledClassHash(casm)
```

### Contract Files

StarkNet contracts are typically distributed as two files:

1. **Sierra File**: Contains the high-level representation of the contract (`.sierra.json`)
2. **CASM File**: Contains the compiled assembly version of the contract (`.casm.json`)

### Declaration vs. Deployment

It's important to understand the difference between:

1. **Declaration**: Registering a contract class on the network (done once per class)
2. **Deployment**: Creating an instance of a declared contract class (can be done multiple times)

## Next Steps

After declaring a contract, you can:

- [Deploy the Contract](./deploy-contract-udc.md) using the Universal Deployer Contract
- Learn how to work with [Typed Data](./typed-data.md) for structured message signing
- Explore [WebSocket](./websocket.md) for real-time updates
