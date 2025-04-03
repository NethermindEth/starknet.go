---
sidebar_position: 1
---

# Simple Call Example

This example demonstrates how to call contract functions on StarkNet without modifying the state. It uses an ERC20 token contract, but the same approach can be applied to any smart contract.

## Overview

The simple call example shows how to:

1. Initialize a connection to a StarkNet provider
2. Create function calls with and without calldata
3. Execute the calls and retrieve the results

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the `.env` file
3. Uncomment and assign your account address to the `ACCOUNT_ADDRESS` variable in the `.env` file

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"os"

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

	// Get the RPC provider URL and account address from environment variables
	providerURL := os.Getenv("RPC_PROVIDER_URL")
	if providerURL == "" {
		panic("RPC_PROVIDER_URL environment variable is not set")
	}

	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	if accountAddress == "" {
		panic("ACCOUNT_ADDRESS environment variable is not set")
	}

	// Initialize connection to provider
	provider, err := rpc.NewProvider(providerURL)
	if err != nil {
		panic(fmt.Sprintf("Error initializing provider: %s", err))
	}

	// ERC20 contract address (STRK token on Sepolia)
	contractAddress := "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"

	// Example 1: Call without calldata - get token name
	nameCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromName("name"),
		Calldata:           []string{},
	}

	nameResult, err := provider.Call(context.Background(), nameCall, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(fmt.Sprintf("Error calling name function: %s", err))
	}

	// Parse the name result (string in felt array format)
	name := utils.HexArrToString(nameResult)
	fmt.Printf("Token name: %s\n", name)

	// Example 2: Call with calldata - get token balance
	balanceCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromName("balanceOf"),
		Calldata:           []string{accountAddress},
	}

	balanceResult, err := provider.Call(context.Background(), balanceCall, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(fmt.Sprintf("Error calling balanceOf function: %s", err))
	}

	// Parse the balance result
	if len(balanceResult) > 0 {
		fmt.Printf("Balance of %s: %s\n", accountAddress, balanceResult[0])
	} else {
		fmt.Println("No balance returned")
	}

	// Example 3: Call with calldata - get token decimals
	decimalsCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromName("decimals"),
		Calldata:           []string{},
	}

	decimalsResult, err := provider.Call(context.Background(), decimalsCall, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(fmt.Sprintf("Error calling decimals function: %s", err))
	}

	// Parse the decimals result
	if len(decimalsResult) > 0 {
		fmt.Printf("Token decimals: %s\n", decimalsResult[0])
	} else {
		fmt.Println("No decimals returned")
	}
}
```

## Running the Example

To run this example:

1. Make sure you are in the "simpleCall" directory
2. Execute `go run main.go`

The call outputs will be returned at the end of the execution.

## Expected Output

```
Token name: Starknet Token
Balance of 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef: 1000000000000000000
Token decimals: 18
```

## Key Concepts

### Function Calls

A function call in StarkNet.go is represented by the `FunctionCall` struct:

```go
type FunctionCall struct {
    ContractAddress    string   // The address of the contract to call
    EntryPointSelector string   // The selector of the function to call
    Calldata           []string // The arguments to pass to the function
}
```

### Entry Point Selectors

Entry point selectors are computed from function names using the `GetSelectorFromName` function:

```go
selector := utils.GetSelectorFromName("balanceOf")
```

### Parsing Results

Results from contract calls may need to be parsed depending on the data type:

- For strings: `utils.HexArrToString(result)`
- For numbers: The result is returned as a hex string that can be converted to a number
- For complex types: Custom parsing logic may be required

## Next Steps

After understanding how to make simple calls, you can:

- Learn how to [Deploy an Account](./deploy-account.md)
- Explore how to [Invoke Transactions](./invoke-transaction.md) to modify contract state
- Understand how to [Declare Contracts](./declare-transaction.md) on StarkNet
