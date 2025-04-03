---
sidebar_position: 3
---

# Transactions API Reference

StarkNet.go provides a comprehensive API for creating, signing, and sending transactions on StarkNet. This reference documents the available transaction methods and how to use them.

## Transaction Types

StarkNet supports several transaction types:

1. **Invoke Transactions**: Call functions on deployed contracts
2. **Declare Transactions**: Register new contract classes
3. **Deploy Account Transactions**: Deploy new account contracts

## Invoke Transactions

Invoke transactions are used to call functions on deployed contracts:

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
}
```

### Invoke Transaction Structure

The `InvokeTxnV1` struct represents an invoke transaction:

```go
type InvokeTxnV1 struct {
    MaxFee         *big.Int
    Version        uint64
    Signature      []string
    Nonce          *big.Int
    Type           string
    SenderAddress  *big.Int
    CallData       []string
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
}
```

### Declare Transaction Structure

The `DeclareTxnV2` struct represents a declare transaction:

```go
type DeclareTxnV2 struct {
    MaxFee           *big.Int
    Version          uint64
    Signature        []string
    Nonce            *big.Int
    Type             string
    SenderAddress    *big.Int
    CompiledClassHash *big.Int
    ContractClass    *contracts.SierraContractClass
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
}
```

### Deploy Account Transaction Structure

The `DeployAccountTxn` struct represents a deploy account transaction:

```go
type DeployAccountTxn struct {
    MaxFee              *big.Int
    Version             uint64
    Signature           []string
    Nonce               *big.Int
    Type                string
    ContractAddressSalt *big.Int
    ConstructorCalldata []string
    ClassHash           *big.Int
}
```

## Transaction Responses

### Invoke Transaction Response

The `InvokeTxnResponse` struct represents the response from an invoke transaction:

```go
type InvokeTxnResponse struct {
    TransactionHash *big.Int
}
```

### Declare Transaction Response

The `DeclareTxnResponse` struct represents the response from a declare transaction:

```go
type DeclareTxnResponse struct {
    TransactionHash *big.Int
    ClassHash       *big.Int
}
```

### Deploy Account Transaction Response

The `DeployAccountTxnResponse` struct represents the response from a deploy account transaction:

```go
type DeployAccountTxnResponse struct {
    TransactionHash  *big.Int
    ContractAddress  *big.Int
}
```

## Transaction Receipt

The `TxnReceipt` struct represents a transaction receipt:

```go
type TxnReceipt struct {
    TransactionHash  *big.Int
    Status           string
    ActualFee        *big.Int
    ExecutionStatus  string
    FinalityStatus   string
    BlockHash        *big.Int
    BlockNumber      uint64
    Type             string
    Events           []Event
    ExecutionResources *ExecutionResources
}
```

### Getting Transaction Receipts

To get a transaction receipt:

```go
// Get a transaction receipt
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

## Transaction Status

The `TxnStatus` struct represents the status of a transaction:

```go
type TxnStatus struct {
    Status         string
    FinalityStatus string
    ExecutionStatus string
}
```

### Getting Transaction Status

To get the status of a transaction:

```go
// Get the status of a transaction
status, err := provider.TransactionStatus(context.Background(), txHash)
if err != nil {
    panic(err)
}

fmt.Printf("Transaction status: %s\n", status.Status)
fmt.Printf("Transaction finality status: %s\n", status.FinalityStatus)
```

## Fee Estimation

The `FeeEstimate` struct represents a fee estimate:

```go
type FeeEstimate struct {
    GasConsumed *big.Int
    GasPrice    *big.Int
    OverallFee  *big.Int
}
```

### Estimating Transaction Fees

To estimate the fee for a transaction:

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

## Transaction Simulation

The `SimulatedTransaction` struct represents a simulated transaction:

```go
type SimulatedTransaction struct {
    TransactionTrace *TransactionTrace
    FeeEstimate      *FeeEstimate
}
```

### Simulating Transactions

To simulate a transaction:

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

// Create an invoke transaction
invokeTx := rpc.InvokeTxnV1{
    MaxFee:         maxFee,
    Version:        1,
    Signature:      signature,
    Nonce:          nonce,
    Type:           "INVOKE",
    SenderAddress:  senderAddress,
    CallData:       calldata,
}

// Simulate the transaction
result, err := provider.SimulateTransaction(
    context.Background(),
    []rpc.BroadcastTxn{invokeTx},
    rpc.BlockID{Tag: "latest"},
    rpc.SimulationFlagSkipValidate,
)
if err != nil {
    panic(err)
}

fmt.Printf("Simulation result: %+v\n", result)
```

## Transaction Tracing

The `TransactionTrace` struct represents a transaction trace:

```go
type TransactionTrace struct {
    ValidateInvocation *FunctionInvocation
    FunctionInvocation *FunctionInvocation
    FeeTransferInvocation *FunctionInvocation
    ExecuteInvocation *FunctionInvocation
}
```

### Tracing Transactions

To trace a transaction:

```go
// Trace a transaction
trace, err := provider.TraceTransaction(context.Background(), txHash)
if err != nil {
    panic(err)
}

fmt.Printf("Transaction trace: %+v\n", trace)
```

## Waiting for Transactions

To wait for a transaction to be accepted:

```go
// Wait for the transaction to be accepted
receipt, err := provider.WaitForTransaction(context.Background(), txHash, 5, 2)
if err != nil {
    panic(err)
}

fmt.Printf("Transaction status: %s\n", receipt.Status)
```

The `WaitForTransaction` method takes the following parameters:

- `ctx`: The context
- `txHash`: The transaction hash
- `retryInterval`: The interval in seconds between retries
- `maxRetries`: The maximum number of retries
