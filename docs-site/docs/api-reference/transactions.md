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
    
    "github.com/NethermindEth/juno/core/felt"
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
    internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
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
    contractAddressHex := "0x..." // ERC20 contract address
    contractAddress, err := new(felt.Felt).SetString(contractAddressHex)
    if err != nil {
        panic(err)
    }
    
    // Create an invoke function call
    functionCall := rpc.InvokeFunctionCall{
        ContractAddress:    contractAddress,
        EntryPointSelector: utils.GetSelectorFromName("transfer"),
        Calldata: []*felt.Felt{
            internalUtils.HexToFelt("0x..."), // Recipient address
            internalUtils.HexToFelt("1000"),  // Amount (in wei)
            internalUtils.HexToFelt("0"),     // Amount high bits (for large numbers)
        },
    }
    
    // Build and send the invoke transaction
    tx, err := acc.BuildAndSendInvokeTxn(
        context.Background(),
        []rpc.InvokeFunctionCall{functionCall},
        1.5, // Fee multiplier (50% buffer)
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
}
```

### Invoke Transaction Structure

The `InvokeTxnV3` struct represents the latest version of an invoke transaction:

```go
type InvokeTxnV3 struct {
    Type           TransactionType       `json:"type"`
    SenderAddress  *felt.Felt            `json:"sender_address"`
    Calldata       []*felt.Felt          `json:"calldata"`
    Version        TransactionVersion    `json:"version"`
    Signature      []*felt.Felt          `json:"signature"`
    Nonce          *felt.Felt            `json:"nonce"`
    ResourceBounds ResourceBoundsMapping `json:"resource_bounds"`
    Tip            U64                   `json:"tip"`
    // The data needed to allow the paymaster to pay for the transaction in native tokens
    PayMasterData []*felt.Felt `json:"paymaster_data"`
    // The data needed to deploy the account contract from which this tx will be initiated
    AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
    // The storage domain of the account's nonce (an account has a nonce per DA mode)
    NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
    // The storage domain of the account's balance from which fee will be charged
    FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}
```

## Declare Transactions

Declare transactions are used to register new contract classes on StarkNet:

```go
package main

import (
    "context"
    "fmt"
    "os"
    
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
    sierraContent, err := os.ReadFile("contract.sierra.json")
    if err != nil {
        panic(err)
    }
    
    casmContent, err := os.ReadFile("contract.casm.json")
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
    
    // Build and send the declare transaction
    tx, err := acc.BuildAndSendDeclareTxn(
        context.Background(),
        casm,
        sierra,
        1.5, // Fee multiplier (50% buffer)
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
    fmt.Printf("Class hash: 0x%s\n", tx.ClassHash.Text(16))
}
```

### Declare Transaction Structure

The `DeclareTxnV3` struct represents the latest version of a declare transaction:

```go
type DeclareTxnV3 struct {
    Type              TransactionType       `json:"type"`
    SenderAddress     *felt.Felt            `json:"sender_address"`
    CompiledClassHash *felt.Felt            `json:"compiled_class_hash"`
    Version           TransactionVersion    `json:"version"`
    Signature         []*felt.Felt          `json:"signature"`
    Nonce             *felt.Felt            `json:"nonce"`
    ClassHash         *felt.Felt            `json:"class_hash"`
    ResourceBounds    ResourceBoundsMapping `json:"resource_bounds"`
    Tip               U64                   `json:"tip"`
    // The data needed to allow the paymaster to pay for the transaction in native tokens
    PayMasterData []*felt.Felt `json:"paymaster_data"`
    // The data needed to deploy the account contract from which this tx will be initiated
    AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
    // The storage domain of the account's nonce (an account has a nonce per DA mode)
    NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
    // The storage domain of the account's balance from which fee will be charged
    FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
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
    
    "github.com/NethermindEth/juno/core/felt"
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/curve"
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
    internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
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
    classHashHex := "0x..." 
    classHash, err := new(felt.Felt).SetString(classHashHex)
    if err != nil {
        panic(err)
    }
    
    // Constructor calldata (public key)
    publicKeyFelt, err := new(felt.Felt).SetString(publicKey.Text(16))
    if err != nil {
        panic(err)
    }
    constructorCalldata := []*felt.Felt{publicKeyFelt}
    
    // Salt for address generation
    salt, err := new(felt.Felt).SetUint64(0)
    if err != nil {
        panic(err)
    }
    
    // Create a keystore for the account
    keystore := account.NewMemKeystore()
    if err := keystore.Put(publicKey.Text(16), privateKey); err != nil {
        panic(err)
    }
    
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
    
    // Build and estimate the deploy account transaction
    deployTx, precomputedAddress, err := acc.BuildAndEstimateDeployAccountTxn(
        context.Background(),
        salt,
        classHash,
        constructorCalldata,
        1.5, // Fee multiplier (50% buffer)
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Precomputed account address: 0x%s\n", precomputedAddress.Text(16))
    fmt.Println("Fund this address with STRK tokens before sending the transaction")
    
    // After funding, send the transaction
    txResponse, err := acc.SendTransaction(context.Background(), deployTx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Deploy transaction hash: 0x%s\n", txResponse.TransactionHash.Text(16))
}
```

### Deploy Account Transaction Structure

The `DeployAccountTxnV3` struct represents the latest version of a deploy account transaction:

```go
type DeployAccountTxnV3 struct {
    Type                TransactionType       `json:"type"`
    Version             TransactionVersion    `json:"version"`
    Signature           []*felt.Felt          `json:"signature"`
    Nonce               *felt.Felt            `json:"nonce"`
    ContractAddressSalt *felt.Felt            `json:"contract_address_salt"`
    ConstructorCalldata []*felt.Felt          `json:"constructor_calldata"`
    ClassHash           *felt.Felt            `json:"class_hash"`
    ResourceBounds      ResourceBoundsMapping `json:"resource_bounds"`
    Tip                 U64                   `json:"tip"`
    // The data needed to allow the paymaster to pay for the transaction in native tokens
    PayMasterData []*felt.Felt `json:"paymaster_data"`
    // The storage domain of the account's nonce (an account has a nonce per DA mode)
    NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
    // The storage domain of the account's balance from which fee will be charged
    FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}
```

## Transaction Responses

### Invoke Transaction Response

The `AddInvokeTransactionResponse` struct represents the response from an invoke transaction:

```go
type AddInvokeTransactionResponse struct {
    TransactionHash *felt.Felt `json:"transaction_hash"`
}
```

### Declare Transaction Response

The `AddDeclareTransactionResponse` struct represents the response from a declare transaction:

```go
type AddDeclareTransactionResponse struct {
    TransactionHash *felt.Felt `json:"transaction_hash"`
    ClassHash       *felt.Felt `json:"class_hash"`
}
```

### Deploy Account Transaction Response

The `AddDeployAccountTransactionResponse` struct represents the response from a deploy account transaction:

```go
type AddDeployAccountTransactionResponse struct {
    TransactionHash  *felt.Felt `json:"transaction_hash"`
    ContractAddress  *felt.Felt `json:"contract_address"`
}
```

## Transaction Receipt

The `TransactionReceiptWithBlockInfo` struct represents a transaction receipt:

```go
type TransactionReceiptWithBlockInfo struct {
    TransactionHash  *felt.Felt                `json:"transaction_hash"`
    Status           TransactionStatus         `json:"status"`
    ActualFee        *felt.Felt                `json:"actual_fee"`
    ExecutionStatus  TxnExecutionStatus       `json:"execution_status"`
    FinalityStatus   TxnFinalityStatus        `json:"finality_status"`
    BlockHash        *felt.Felt                `json:"block_hash"`
    BlockNumber      U64                       `json:"block_number"`
    Type             TransactionType           `json:"type"`
    Events           []Event                   `json:"events"`
    ExecutionResources *ExecutionResources     `json:"execution_resources"`
    MessagesSent     []MsgToL1                 `json:"messages_sent"`
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
contractAddressFelt, err := new(felt.Felt).SetString("0x...") // Contract address
if err != nil {
    panic(err)
}

// Create an invoke function call
functionCall := rpc.InvokeFunctionCall{
    ContractAddress:    contractAddressFelt,
    EntryPointSelector: utils.GetSelectorFromName("transfer"),
    Calldata: []*felt.Felt{
        internalUtils.HexToFelt("0x..."), // Recipient address
        internalUtils.HexToFelt("1000"),  // Amount (in wei)
        internalUtils.HexToFelt("0"),     // Amount high bits (for large numbers)
    },
}

// Build the invoke transaction for fee estimation
invokeTx, err := acc.BuildInvokeTxn(
    context.Background(),
    []rpc.InvokeFunctionCall{functionCall},
    1.5, // Fee multiplier (50% buffer)
)
if err != nil {
    panic(err)
}

// Estimate the fee
feeEstimate, err := provider.EstimateFee(
    context.Background(),
    []rpc.BroadcastTxn{invokeTx},
    rpc.BlockID{Tag: "latest"},
)
if err != nil {
    panic(err)
}

fmt.Printf("Estimated fee: %s wei\n", feeEstimate[0].OverallFee.Text(10))
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
contractAddressFelt, err := new(felt.Felt).SetString("0x...") // Contract address
if err != nil {
    panic(err)
}

// Create an invoke function call
functionCall := rpc.InvokeFunctionCall{
    ContractAddress:    contractAddressFelt,
    EntryPointSelector: utils.GetSelectorFromName("transfer"),
    Calldata: []*felt.Felt{
        internalUtils.HexToFelt("0x..."), // Recipient address
        internalUtils.HexToFelt("1000"),  // Amount (in wei)
        internalUtils.HexToFelt("0"),     // Amount high bits (for large numbers)
    },
}

// Build the invoke transaction for simulation
invokeTx, err := acc.BuildInvokeTxn(
    context.Background(),
    []rpc.InvokeFunctionCall{functionCall},
    1.5, // Fee multiplier (50% buffer)
)
if err != nil {
    panic(err)
}

// Simulate the transaction
result, err := provider.SimulateTransaction(
    context.Background(),
    []rpc.BroadcastTxn{invokeTx},
    rpc.BlockID{Tag: "latest"},
    []rpc.SimulationFlag{rpc.SimulationFlagSkipValidate},
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
// Convert transaction hash from hex string to felt.Felt
txHashHex := "0x..."
txHash, err := new(felt.Felt).SetString(txHashHex)
if err != nil {
    panic(err)
}

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
// Convert transaction hash from hex string to felt.Felt
txHashHex := "0x..."
txHash, err := new(felt.Felt).SetString(txHashHex)
if err != nil {
    panic(err)
}

// Wait for the transaction to be accepted
receipt, err := acc.WaitForTransactionReceipt(
    context.Background(),
    txHash,
    5 * time.Second, // Poll interval
)
if err != nil {
    panic(err)
}

fmt.Printf("Transaction status: %s\n", receipt.Status)
```

The `WaitForTransactionReceipt` method takes the following parameters:

- `ctx`: The context
- `transactionHash`: The transaction hash as a `*felt.Felt`
- `pollInterval`: The interval between polling attempts as a `time.Duration`
