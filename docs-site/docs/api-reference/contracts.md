---
sidebar_position: 5
---

# Contracts API Reference

StarkNet.go provides a comprehensive API for working with StarkNet contracts. This reference documents the available contract methods and how to use them.

## Contract Classes

StarkNet contracts are represented by contract classes. The `contracts` package provides methods for working with contract classes:

```go
package main

import (
    "fmt"
    "io/ioutil"
    
    "github.com/NethermindEth/starknet.go/contracts"
)

func main() {
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
    
    fmt.Printf("Compiled class hash: 0x%s\n", compiledClassHash.Text(16))
}
```

### Sierra Contract Class

The `SierraContractClass` struct represents a Sierra contract class:

```go
type SierraContractClass struct {
    Contract struct {
        ABI                []interface{} `json:"abi"`
        EntryPointsByType  interface{}   `json:"entry_points_by_type"`
        Sierra_Program     []string      `json:"sierra_program"`
        ContractClassVersion string      `json:"contract_class_version"`
    } `json:"contract"`
}
```

### CASM Contract Class

The `CasmContractClass` struct represents a CASM (compiled) contract class:

```go
type CasmContractClass struct {
    Prime         string                 `json:"prime"`
    Bytecode      []string               `json:"bytecode"`
    Hints         map[string]interface{} `json:"hints"`
    Pythonic_Hints map[string]interface{} `json:"pythonic_hints"`
    Compiler_Version string              `json:"compiler_version"`
    EntryPoints   struct {
        External []struct {
            Selector string `json:"selector"`
            Offset   int    `json:"offset"`
        } `json:"external"`
        L1Handler []struct {
            Selector string `json:"selector"`
            Offset   int    `json:"offset"`
        } `json:"l1_handler"`
        Constructor []struct {
            Selector string `json:"selector"`
            Offset   int    `json:"offset"`
        } `json:"constructor"`
    } `json:"entry_points"`
}
```

## Computing Class Hashes

To compute the hash of a contract class:

```go
// Compute the class hash of a Sierra contract
classHash, err := contracts.ComputeClassHash(sierra)
if err != nil {
    panic(err)
}

fmt.Printf("Class hash: 0x%s\n", classHash.Text(16))

// Compute the compiled class hash of a CASM contract
compiledClassHash, err := contracts.ComputeCompiledClassHash(casm)
if err != nil {
    panic(err)
}

fmt.Printf("Compiled class hash: 0x%s\n", compiledClassHash.Text(16))
```

## Contract ABI

The contract ABI (Application Binary Interface) defines the methods and events of a contract. The `contracts` package provides methods for working with contract ABIs:

```go
package main

import (
    "fmt"
    "io/ioutil"
    
    "github.com/NethermindEth/starknet.go/contracts"
)

func main() {
    // Load the contract Sierra file
    sierraContent, err := ioutil.ReadFile("contract.sierra.json")
    if err != nil {
        panic(err)
    }
    
    // Parse the contract file
    sierra, err := contracts.NewSierraContractClass(sierraContent)
    if err != nil {
        panic(err)
    }
    
    // Get the contract ABI
    abi := sierra.Contract.ABI
    
    // Print the ABI
    fmt.Printf("Contract ABI: %+v\n", abi)
}
```

## Contract Deployment

To deploy a contract using the Universal Deployer Contract (UDC):

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

## Contract Interaction

### Calling Contract Functions

To call a contract function without modifying the state:

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

### Invoking Contract Functions

To invoke a contract function and modify the state:

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

## Function Call Structure

The `FunctionCall` struct represents a function call to a contract:

```go
type FunctionCall struct {
    ContractAddress    string
    EntryPointSelector string
    Calldata           []string
}
```

### Getting Selectors

To get a selector from a function name:

```go
// Get a selector from a function name
selector := utils.GetSelectorFromName("transfer")
fmt.Printf("Selector: 0x%s\n", selector)
```

## Events

StarkNet contracts can emit events during execution. The `Event` struct represents an event:

```go
type Event struct {
    FromAddress string
    Keys        []string
    Data        []string
    BlockNumber uint64
    BlockHash   string
    TransactionHash string
}
```

### Querying Events

To query events emitted by a contract:

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

### Subscribing to Events

To subscribe to events in real-time using WebSocket:

```go
// Create a channel to receive events
eventsChan := make(chan *rpc.Event)

// Subscribe to events
sub, err := wsProvider.SubscribeEvents(
    context.Background(),
    eventsChan,
    rpc.EventsInput{
        Address: contractAddress,
        Keys:    [][]string{{"0x..."}}, // Event key to filter by
    },
)
if err != nil {
    panic(err)
}

// Process events
go func() {
    for {
        select {
        case event := <-eventsChan:
            fmt.Printf("Event received from %s:\n", event.FromAddress)
            fmt.Printf("  Keys: %v\n", event.Keys)
            fmt.Printf("  Data: %v\n", event.Data)
        case err := <-sub.Err():
            if err != nil {
                fmt.Printf("Subscription error: %v\n", err)
                return
            }
        }
    }
}()

// Unsubscribe when done
defer sub.Unsubscribe()
```

## Storage

StarkNet contracts store their state in storage slots. The `rpc` package provides methods for accessing contract storage:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/NethermindEth/starknet.go/rpc"
)

func main() {
    // Initialize a provider
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Contract address
    contractAddress := "0x..."
    
    // Storage key
    storageKey := "0x..."
    
    // Get the value of a storage slot
    value, err := provider.StorageAt(
        context.Background(),
        contractAddress,
        storageKey,
        rpc.BlockID{Tag: "latest"},
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Storage value: 0x%s\n", value.Text(16))
}
```

## Contract Simulation

To simulate a contract function call without sending a transaction:

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
}
```

## Contract Tracing

To trace a contract function call:

```go
// Trace a transaction
trace, err := provider.TraceTransaction(context.Background(), txHash)
if err != nil {
    panic(err)
}

fmt.Printf("Transaction trace: %+v\n", trace)
```

The `TransactionTrace` struct represents a transaction trace:

```go
type TransactionTrace struct {
    ValidateInvocation *FunctionInvocation
    FunctionInvocation *FunctionInvocation
    FeeTransferInvocation *FunctionInvocation
    ExecuteInvocation *FunctionInvocation
}
```

The `FunctionInvocation` struct represents a function invocation in a trace:

```go
type FunctionInvocation struct {
    CallType         string
    ContractAddress  string
    EntryPointType   string
    EntryPointSelector string
    Calldata         []string
    Result           []string
    Calls            []*FunctionInvocation
    Events           []Event
    Messages         []Message
}
```
