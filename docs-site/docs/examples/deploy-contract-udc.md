---
sidebar_position: 5
---

# Deploy Contract Using UDC Example

This example demonstrates how to deploy a contract on StarkNet using the Universal Deployer Contract (UDC). The UDC is a utility contract that allows you to deploy other contracts on StarkNet.

## Overview

The deploy contract UDC example shows how to:

1. Initialize a connection to a StarkNet provider
2. Load an existing account
3. Prepare the deployment parameters
4. Call the UDC to deploy the contract
5. Extract the deployed contract address from the events

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the `.env` file
3. Uncomment and assign your account address to the `ACCOUNT_ADDRESS` variable in the `.env` file
4. Uncomment and assign your public key to the `PUBLIC_KEY` variable in the `.env` file
5. Uncomment and assign your private key to the `PRIVATE_KEY` variable in the `.env` file
6. Have a declared contract class hash ready for deployment

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"

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

	// UDC contract address (same on all networks)
	udcAddress := "0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf"

	// Class hash of the contract to deploy (replace with your contract's class hash)
	// This should be a class hash that has already been declared on the network
	classHash := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	// Constructor calldata (replace with your contract's constructor arguments)
	constructorCalldata := []string{
		"0x1", // Example constructor argument 1
		"0x2", // Example constructor argument 2
	}

	// Salt for address generation (can be any value, using a random value for uniqueness)
	salt := utils.GetHexString(big.NewInt(12345))

	// Prepare the calldata for the UDC
	udcCalldata := []string{
		classHash,           // Class hash
		salt,                // Salt
		"0",                 // Unique (0 for normal deployment)
		strconv.Itoa(len(constructorCalldata)), // Constructor calldata length
	}

	// Append constructor calldata
	udcCalldata = append(udcCalldata, constructorCalldata...)

	// Create a function call to the UDC
	functionCall := rpc.FunctionCall{
		ContractAddress:    udcAddress,
		EntryPointSelector: utils.GetSelectorFromName("deployContract"),
		Calldata:           udcCalldata,
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

	// Extract the deployed contract address from the events
	var contractAddress string
	if len(receipt.Events) > 0 {
		for _, event := range receipt.Events {
			if event.FromAddress == udcAddress {
				contractAddress = event.Data[0]
				break
			}
		}
	}

	if contractAddress != "" {
		fmt.Printf("Contract deployed successfully at address: %s\n", contractAddress)
	} else {
		fmt.Println("Could not find deployed contract address in events")
	}
}
```

## Running the Example

To run this example:

1. Make sure you are in the "deployContractUDC" directory
2. Replace the `classHash` and `constructorCalldata` with your own values
3. Execute `go run main.go`

## Expected Output

```
Estimated fee: 12345678901234 wei
Transaction hash: 0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
Waiting for the transaction to be accepted...
Transaction status: ACCEPTED_ON_L2
Actual fee: 12345678901234 wei
Contract deployed successfully at address: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
```

## Key Concepts

### Universal Deployer Contract (UDC)

The UDC is a utility contract that allows you to deploy other contracts on StarkNet. It has a fixed address on all StarkNet networks:

```
0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf
```

### Contract Deployment Process

The process of deploying a contract using the UDC involves:

1. Preparing the deployment parameters (class hash, salt, constructor calldata)
2. Calling the `deployContract` function on the UDC
3. Extracting the deployed contract address from the events

### UDC Calldata Structure

The calldata for the UDC's `deployContract` function has the following structure:

```
[
    classHash,           // Class hash of the contract to deploy
    salt,                // Salt for address generation
    unique,              // 0 for normal deployment, 1 for unique deployment
    constructorCalldataLength, // Length of the constructor calldata
    ...constructorCalldata     // Constructor arguments
]
```

### Contract Address Extraction

The UDC emits an event containing the deployed contract address. You can extract this address from the transaction receipt:

```go
// Extract the deployed contract address from the events
var contractAddress string
if len(receipt.Events) > 0 {
    for _, event := range receipt.Events {
        if event.FromAddress == udcAddress {
            contractAddress = event.Data[0]
            break
        }
    }
}
```

### Salt and Address Determinism

The contract address is deterministically computed from:

1. The class hash
2. The salt
3. The constructor calldata
4. The deployer address (UDC address)

Using the same parameters will result in the same contract address. To deploy multiple instances of the same contract, use different salt values.

## Next Steps

After deploying a contract, you can:

- Learn how to work with [Typed Data](./typed-data.md) for structured message signing
- Explore [WebSocket](./websocket.md) for real-time updates
- Interact with your deployed contract using [Simple Call](./simple-call.md) or [Invoke Transaction](./invoke-transaction.md)
