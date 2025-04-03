---
sidebar_position: 3
---

# Transactions

StarkNet.go supports various transaction types for interacting with the StarkNet blockchain. This guide covers the main transaction types and how to use them.

## Transaction Types

StarkNet supports several transaction types:

1. **Invoke Transactions**: Call functions on deployed contracts
2. **Declare Transactions**: Register new contract classes
3. **Deploy Account Transactions**: Deploy new account contracts

## Invoke Transactions

Invoke transactions are used to call functions on deployed contracts. Here's how to create and send an invoke transaction:

```go
package main

import (
    "context"
    "fmt"
    "math/big"
    
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
    
    // Create a function call to an ERC20 transfer
    contractAddress := "0x..." // ERC20 contract address
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

## Declare Transactions

Declare transactions are used to register new contract classes on StarkNet:

```go
package main

import (
    "context"
    "fmt"
    "io/ioutil"
    
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/contracts"
    "github.com/NethermindEth/starknet.go/rpc"
)

func main() {
    // Initialize provider and account (see Account Management section)
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Load your account
    // ... (account setup code)
    
    // Load the contract Sierra and CASM files
    sierraContent, err := ioutil.ReadFile("contract.sierra.json")
    if err != nil {
        panic(err)
    }
    
    casmContent, err := ioutil.ReadFile("contract.casm.json")
    if err != nil {
        panic(err)
    }
    
    // Parse the contract files
    sierra, err := contracts.NewSierraContractClass(sierraContent)
    if err != nil {
        panic(err)
    }
    
    casm, err := contracts.NewCasmContractClass(casmContent)
    if err != nil {
        panic(err)
    }
    
    // Calculate the compiled class hash
    compiledClassHash, err := contracts.ComputeCompiledClassHash(casm)
    if err != nil {
        panic(err)
    }
    
    // Declare the contract
    tx, err := acc.Declare(context.Background(), sierra, compiledClassHash)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
    fmt.Printf("Class hash: 0x%s\n", tx.ClassHash.Text(16))
    
    // Wait for the transaction to be accepted
    receipt, err := provider.WaitForTransaction(context.Background(), tx.TransactionHash, 5, 2)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction status: %s\n", receipt.Status)
}
```

## Deploy Account Transactions

Deploy account transactions are used to deploy new account contracts:

```go
package main

import (
    "context"
    "fmt"
    "math/big"
    
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/curve"
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
)

func main() {
    // Initialize a provider
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Generate a new key pair
    starkCurve := curve.NewStarkCurve()
    privateKey, err := utils.GeneratePrivateKey()
    if err != nil {
        panic(err)
    }
    
    publicKey, err := starkCurve.GetPublicKey(privateKey)
    if err != nil {
        panic(err)
    }
    
    // Class hash of the account contract (OpenZeppelin account contract)
    classHash := "0x..." 
    
    // Constructor calldata (public key)
    constructorCalldata := []string{publicKey.Text(16)}
    
    // Salt for address generation
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
        panic(err)
    }
    
    fmt.Printf("Account address: 0x%s\n", accountAddress.Text(16))
    
    // Create a keystore for the account
    keystore := account.NewMemKeystore()
    keystore.Put(publicKey.Text(16), privateKey)
    
    // Create a new account instance
    acc, err := account.NewAccount(
        provider,
        accountAddress,
        publicKey.Text(16),
        keystore,
        1, // Cairo version (1 for Cairo 1.0)
    )
    if err != nil {
        panic(err)
    }
    
    // Deploy the account
    deployTx, err := acc.Deploy(context.Background(), salt, classHash, constructorCalldata)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Deploy transaction hash: 0x%s\n", deployTx.TransactionHash.Text(16))
    
    // Wait for the transaction to be accepted
    receipt, err := provider.WaitForTransaction(context.Background(), deployTx.TransactionHash, 5, 2)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction status: %s\n", receipt.Status)
}
```

## Estimating Transaction Fees

Before sending a transaction, you can estimate the fees:

```go
// Create a function call
functionCall := rpc.FunctionCall{
    ContractAddress:    contractAddress,
    EntryPointSelector: utils.GetSelectorFromName("transfer"),
    Calldata: []string{
        "0x...", // Recipient address
        "1000",  // Amount (in wei)
        "0",     // Amount high bits (for large numbers)
    },
}

// Estimate the fee
feeEstimate, err := acc.EstimateFee(context.Background(), []rpc.FunctionCall{functionCall}, nil)
if err != nil {
    panic(err)
}

fmt.Printf("Estimated fee: %s wei\n", feeEstimate.OverallFee.Text(10))
```

## Getting Transaction Receipts

After sending a transaction, you can get its receipt:

```go
receipt, err := provider.TransactionReceipt(context.Background(), txHash)
if err != nil {
    panic(err)
}

fmt.Printf("Transaction status: %s\n", receipt.Status)
fmt.Printf("Actual fee: %s wei\n", receipt.ActualFee.Text(10))

// Check for events
if len(receipt.Events) > 0 {
    fmt.Println("Events:")
    for i, event := range receipt.Events {
        fmt.Printf("  Event %d: From %s, Key %s\n", i, event.FromAddress, event.Keys[0])
    }
}
```

## Next Steps

Now that you understand transactions, you can:

- Learn about [Contract Interaction](./contract-interaction.md)
- Explore the [API Reference](../api-reference/rpc.md)
