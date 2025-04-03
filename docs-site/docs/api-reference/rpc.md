---
sidebar_position: 1
---

# RPC API Reference

StarkNet.go provides a comprehensive implementation of the StarkNet RPC v0.8.0 specification. This reference documents the available RPC methods and how to use them.

## Provider Initialization

To interact with the StarkNet RPC API, you first need to initialize a provider:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/NethermindEth/starknet.go/rpc"
)

func main() {
    // Initialize a provider with a StarkNet node URL
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Now you can use the provider to interact with StarkNet
    fmt.Println("Provider initialized successfully!")
}
```

## WebSocket Provider

For real-time updates, you can use the WebSocket provider:

```go
// Initialize a WebSocket provider
wsProvider, err := rpc.NewWebsocketProvider("wss://starknet-sepolia.infura.io/ws/v3/YOUR_API_KEY")
if err != nil {
    panic(err)
}
defer wsProvider.Close() // Close the connection when done
```

## Block Methods

### Get Block Number

```go
// Get the latest block number
blockNumber, err := provider.BlockNumber(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("Latest block number: %d\n", blockNumber)
```

### Get Block Hash and Number

```go
// Get the latest block hash and number
blockHashAndNumber, err := provider.BlockHashAndNumber(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("Block number: %d\n", blockHashAndNumber.BlockNumber)
fmt.Printf("Block hash: 0x%s\n", blockHashAndNumber.BlockHash.Text(16))
```

### Get Block with Transactions

```go
// Get a block with its transactions
block, err := provider.BlockWithTxs(context.Background(), rpc.BlockID{Number: 12345})
if err != nil {
    panic(err)
}
fmt.Printf("Block number: %d\n", block.BlockNumber)
fmt.Printf("Block hash: 0x%s\n", block.BlockHash.Text(16))
fmt.Printf("Number of transactions: %d\n", len(block.Transactions))
```

### Get Block with Transaction Hashes

```go
// Get a block with transaction hashes
block, err := provider.BlockWithTxHashes(context.Background(), rpc.BlockID{Tag: "latest"})
if err != nil {
    panic(err)
}
fmt.Printf("Block number: %d\n", block.BlockNumber)
fmt.Printf("Block hash: 0x%s\n", block.BlockHash.Text(16))
fmt.Printf("Number of transaction hashes: %d\n", len(block.TransactionHashes))
```

### Get Block with Receipts

```go
// Get a block with transaction receipts
block, err := provider.BlockWithReceipts(context.Background(), rpc.BlockID{Hash: blockHash})
if err != nil {
    panic(err)
}
fmt.Printf("Block number: %d\n", block.BlockNumber)
fmt.Printf("Block hash: 0x%s\n", block.BlockHash.Text(16))
fmt.Printf("Number of transactions with receipts: %d\n", len(block.Transactions))
```

### Get Block Transaction Count

```go
// Get the number of transactions in a block
count, err := provider.BlockTransactionCount(context.Background(), rpc.BlockID{Number: 12345})
if err != nil {
    panic(err)
}
fmt.Printf("Number of transactions in block: %d\n", count)
```

## Transaction Methods

### Get Transaction by Hash

```go
// Get a transaction by its hash
tx, err := provider.Transaction(context.Background(), txHash)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction type: %s\n", tx.Type)
fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
```

### Get Transaction by Block ID and Index

```go
// Get a transaction by block ID and index
tx, err := provider.TransactionByBlockIdAndIndex(context.Background(), rpc.BlockID{Number: 12345}, 0)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction type: %s\n", tx.Type)
fmt.Printf("Transaction hash: 0x%s\n", tx.TransactionHash.Text(16))
```

### Get Transaction Receipt

```go
// Get a transaction receipt
receipt, err := provider.TransactionReceipt(context.Background(), txHash)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction status: %s\n", receipt.Status)
fmt.Printf("Actual fee: %s wei\n", receipt.ActualFee.Text(10))
```

### Get Transaction Status

```go
// Get the status of a transaction
status, err := provider.TransactionStatus(context.Background(), txHash)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction status: %s\n", status.Status)
fmt.Printf("Transaction finality status: %s\n", status.FinalityStatus)
```

### Add Invoke Transaction

```go
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

// Add the invoke transaction
response, err := provider.AddInvokeTransaction(context.Background(), invokeTx)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction hash: 0x%s\n", response.TransactionHash.Text(16))
```

### Add Declare Transaction

```go
// Create a declare transaction
declareTx := rpc.DeclareTxnV2{
    MaxFee:           maxFee,
    Version:          2,
    Signature:        signature,
    Nonce:            nonce,
    Type:             "DECLARE",
    SenderAddress:    senderAddress,
    CompiledClassHash: compiledClassHash,
    ContractClass:    contractClass,
}

// Add the declare transaction
response, err := provider.AddDeclareTransaction(context.Background(), declareTx)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction hash: 0x%s\n", response.TransactionHash.Text(16))
fmt.Printf("Class hash: 0x%s\n", response.ClassHash.Text(16))
```

### Add Deploy Account Transaction

```go
// Create a deploy account transaction
deployAccountTx := rpc.DeployAccountTxn{
    MaxFee:           maxFee,
    Version:          1,
    Signature:        signature,
    Nonce:            nonce,
    Type:             "DEPLOY_ACCOUNT",
    ContractAddressSalt: salt,
    ConstructorCalldata: constructorCalldata,
    ClassHash:        classHash,
}

// Add the deploy account transaction
response, err := provider.AddDeployAccountTransaction(context.Background(), deployAccountTx)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction hash: 0x%s\n", response.TransactionHash.Text(16))
fmt.Printf("Contract address: 0x%s\n", response.ContractAddress.Text(16))
```

## Contract Methods

### Call Contract

```go
// Create a function call
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
fmt.Printf("Result: %v\n", result)
```

### Estimate Fee

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

// Estimate the fee
feeEstimate, err := provider.EstimateFee(context.Background(), []rpc.BroadcastTxn{invokeTx}, rpc.BlockID{Tag: "latest"})
if err != nil {
    panic(err)
}
fmt.Printf("Estimated fee: %s wei\n", feeEstimate[0].OverallFee.Text(10))
```

### Get Storage At

```go
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
```

### Get Class

```go
// Get a contract class by its hash
class, err := provider.Class(context.Background(), rpc.BlockID{Tag: "latest"}, classHash)
if err != nil {
    panic(err)
}
fmt.Printf("Contract class: %+v\n", class)
```

### Get Class Hash At

```go
// Get the class hash of a contract at a specific address
classHash, err := provider.ClassHashAt(context.Background(), rpc.BlockID{Tag: "latest"}, contractAddress)
if err != nil {
    panic(err)
}
fmt.Printf("Class hash: 0x%s\n", classHash.Text(16))
```

### Get Class At

```go
// Get the contract class at a specific address
class, err := provider.ClassAt(context.Background(), rpc.BlockID{Tag: "latest"}, contractAddress)
if err != nil {
    panic(err)
}
fmt.Printf("Contract class: %+v\n", class)
```

## State Methods

### Get Nonce

```go
// Get the nonce of an account
nonce, err := provider.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accountAddress)
if err != nil {
    panic(err)
}
fmt.Printf("Nonce: %s\n", nonce.Text(16))
```

### Get State Update

```go
// Get the state update of a block
stateUpdate, err := provider.StateUpdate(context.Background(), rpc.BlockID{Number: 12345})
if err != nil {
    panic(err)
}
fmt.Printf("State update: %+v\n", stateUpdate)
```

## Chain Methods

### Get Chain ID

```go
// Get the chain ID
chainID, err := provider.ChainID(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("Chain ID: %s\n", chainID)
```

### Get Syncing Status

```go
// Get the syncing status
syncing, err := provider.Syncing(context.Background())
if err != nil {
    panic(err)
}
if syncing.Syncing {
    fmt.Printf("Node is syncing. Current block: %d, Highest block: %d\n", syncing.StartingBlockNumber, syncing.HighestBlockNumber)
} else {
    fmt.Println("Node is not syncing")
}
```

## Event Methods

### Get Events

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

## WebSocket Methods

### Subscribe to New Heads

```go
// Create a channel to receive block headers
newHeadsChan := make(chan *rpc.BlockHeader)

// Subscribe to new block headers
sub, err := wsProvider.SubscribeNewHeads(context.Background(), newHeadsChan, rpc.BlockID{})
if err != nil {
    panic(err)
}

// Process new block headers
go func() {
    for {
        select {
        case newHead := <-newHeadsChan:
            fmt.Printf("New block header received: %d\n", newHead.BlockNumber)
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

### Subscribe to Events

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

### Subscribe to Pending Transactions

```go
// Create a channel to receive pending transactions
pendingTxsChan := make(chan *string)

// Subscribe to pending transactions
sub, err := wsProvider.SubscribePendingTransactions(context.Background(), pendingTxsChan)
if err != nil {
    panic(err)
}

// Process pending transactions
go func() {
    for {
        select {
        case txHash := <-pendingTxsChan:
            fmt.Printf("Pending transaction received: %s\n", *txHash)
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

### Subscribe to Transaction Status

```go
// Create a channel to receive transaction status updates
txStatusChan := make(chan *rpc.TxnStatus)

// Subscribe to transaction status
sub, err := wsProvider.SubscribeTransactionStatus(context.Background(), txStatusChan, txHash)
if err != nil {
    panic(err)
}

// Process transaction status updates
go func() {
    for {
        select {
        case status := <-txStatusChan:
            fmt.Printf("Transaction status update: %s\n", status.Status)
            fmt.Printf("Transaction finality status: %s\n", status.FinalityStatus)
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

## Trace Methods

### Trace Transaction

```go
// Trace a transaction
trace, err := provider.TraceTransaction(context.Background(), txHash)
if err != nil {
    panic(err)
}
fmt.Printf("Transaction trace: %+v\n", trace)
```

### Trace Block Transactions

```go
// Trace all transactions in a block
traces, err := provider.TraceBlockTransactions(context.Background(), rpc.BlockID{Number: 12345})
if err != nil {
    panic(err)
}
fmt.Printf("Number of transaction traces: %d\n", len(traces))
```

## Simulation Methods

### Simulate Transaction

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

## Utility Methods

### Get Spec Version

```go
// Get the RPC spec version
version, err := provider.SpecVersion(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("RPC spec version: %s\n", version)
```
