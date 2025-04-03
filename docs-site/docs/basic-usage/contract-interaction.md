---
sidebar_position: 4
---

# Contract Interaction

StarkNet.go provides several ways to interact with contracts deployed on StarkNet. This guide covers how to call contract functions and handle contract data.

## Reading Contract State

To read the state of a contract without modifying it, you can use the `Call` method:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
)

func main() {
    // Initialize a provider
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // ERC20 contract address
    contractAddress := "0x..."
    
    // Create a function call to get the balance
    functionCall := rpc.FunctionCall{
        ContractAddress:    contractAddress,
        EntryPointSelector: utils.GetSelectorFromName("balanceOf"),
        Calldata: []string{
            "0x...", // Account address to check balance for
        },
    }
    
    // Call the function
    result, err := provider.Call(context.Background(), functionCall, rpc.BlockID{Tag: "latest"})
    if err != nil {
        panic(err)
    }
    
    // Parse the result (balance)
    if len(result) > 0 {
        fmt.Printf("Balance: %s\n", result[0])
    } else {
        fmt.Println("No result returned")
    }
}
```

## Calling Contract Functions with Transactions

To modify the state of a contract, you need to send a transaction:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
)

func main() {
    // Initialize provider and account (see Account Management section)
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Load your account
    // ... (account setup code)
    
    // ERC20 contract address
    contractAddress := "0x..."
    
    // Create a function call to transfer tokens
    functionCall := rpc.FunctionCall{
        ContractAddress:    contractAddress,
        EntryPointSelector: utils.GetSelectorFromName("transfer"),
        Calldata: []string{
            "0x...", // Recipient address
            "1000",  // Amount (in wei)
            "0",     // Amount high bits (for large numbers)
        },
    }
    
    // Execute the transaction
    tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall}, nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
    
    // Wait for the transaction to be accepted
    receipt, err := provider.WaitForTransaction(context.Background(), tx.TransactionHash, 5, 2)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction status: %s\n", receipt.Status)
}
```

## Making Multiple Function Calls in a Single Transaction

You can include multiple function calls in a single transaction:

```go
// Create multiple function calls
functionCall1 := rpc.FunctionCall{
    ContractAddress:    contractAddress1,
    EntryPointSelector: utils.GetSelectorFromName("approve"),
    Calldata: []string{
        "0x...", // Spender address
        "1000",  // Amount (in wei)
        "0",     // Amount high bits (for large numbers)
    },
}

functionCall2 := rpc.FunctionCall{
    ContractAddress:    contractAddress2,
    EntryPointSelector: utils.GetSelectorFromName("transfer"),
    Calldata: []string{
        "0x...", // Recipient address
        "500",   // Amount (in wei)
        "0",     // Amount high bits (for large numbers)
    },
}

// Execute both function calls in a single transaction
tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall1, functionCall2}, nil)
if err != nil {
    panic(err)
}
```

## Deploying Contracts Using UDC

The Universal Deployer Contract (UDC) allows you to deploy contracts on StarkNet:

```go
package main

import (
    "context"
    "fmt"
    "io/ioutil"
    
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/contracts"
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
)

func main() {
    // Initialize provider and account (see Account Management section)
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Load your account
    // ... (account setup code)
    
    // UDC contract address
    udcAddress := "0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf"
    
    // Class hash of the contract to deploy
    classHash := "0x..."
    
    // Constructor calldata
    constructorCalldata := []string{
        "0x...", // Constructor argument 1
        "0x...", // Constructor argument 2
    }
    
    // Salt for address generation
    salt := "0x1"
    
    // Prepare the calldata for the UDC
    udcCalldata := []string{
        classHash,           // Class hash
        salt,                // Salt
        "0",                 // Unique (0 for normal deployment)
        fmt.Sprintf("%d", len(constructorCalldata)), // Constructor calldata length
    }
    
    // Append constructor calldata
    udcCalldata = append(udcCalldata, constructorCalldata...)
    
    // Create a function call to the UDC
    functionCall := rpc.FunctionCall{
        ContractAddress:    udcAddress,
        EntryPointSelector: utils.GetSelectorFromName("deployContract"),
        Calldata:           udcCalldata,
    }
    
    // Execute the transaction
    tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall}, nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
    
    // Wait for the transaction to be accepted
    receipt, err := provider.WaitForTransaction(context.Background(), tx.TransactionHash, 5, 2)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction status: %s\n", receipt.Status)
    
    // Extract the deployed contract address from the events
    if len(receipt.Events) > 0 {
        for _, event := range receipt.Events {
            if event.FromAddress == udcAddress {
                fmt.Printf("Deployed contract address: 0x%s\n", event.Data[0])
                break
            }
        }
    }
}
```

## Handling Events

Contracts can emit events during execution. You can query these events:

```go
// Query events
events, err := provider.Events(context.Background(), rpc.EventsInput{
    FromBlock: rpc.BlockID{Number: 10000},
    ToBlock:   rpc.BlockID{Tag: "latest"},
    Address:   contractAddress,
    Keys:      [][]string{{"0x..."}}, // Event key to filter by
})
if err != nil {
    panic(err)
}

// Process events
for _, event := range events.Events {
    fmt.Printf("Event from %s at block %d:\n", event.FromAddress, event.BlockNumber)
    fmt.Printf("  Keys: %v\n", event.Keys)
    fmt.Printf("  Data: %v\n", event.Data)
}
```

## Next Steps

Now that you understand contract interaction, you can:

- Explore the [API Reference](../api-reference/rpc.md)
- Check out the [Examples](../examples/simple-call.md) for more practical use cases
