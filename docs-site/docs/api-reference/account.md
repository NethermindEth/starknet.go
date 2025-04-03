---
sidebar_position: 2
---

# Account API Reference

StarkNet.go provides a comprehensive account management API that allows you to create, deploy, and manage StarkNet accounts. This reference documents the available account methods and how to use them.

## Account Types

StarkNet.go supports different account implementations:

- Standard account (OpenZeppelin account contract)
- Argent account
- Braavos account
- Custom account implementations

## Account Interface

The `account.AccountInterface` defines the methods that all account implementations must provide:

```go
type AccountInterface interface {
    BuildAndEstimateDeployAccountTxn(ctx context.Context, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt, multiplier float64) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error)
    BuildAndSendInvokeTxn(ctx context.Context, functionCalls []rpc.InvokeFunctionCall, multiplier float64) (*rpc.AddInvokeTransactionResponse, error)
    BuildAndSendDeclareTxn(ctx context.Context, casmClass *contracts.CasmClass, contractClass *contracts.ContractClass, multiplier float64) (*rpc.AddDeclareTransactionResponse, error)
    Nonce(ctx context.Context) (*felt.Felt, error)
    SendTransaction(ctx context.Context, txn rpc.BroadcastTxn) (*rpc.TransactionResponse, error)
    Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
    SignInvokeTransaction(ctx context.Context, tx rpc.InvokeTxnType) error
    SignDeployAccountTransaction(ctx context.Context, tx rpc.DeployAccountType, precomputeAddress *felt.Felt) error
    SignDeclareTransaction(ctx context.Context, tx rpc.DeclareTxnType) error
    TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error)
    TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error)
    TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error)
    WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceiptWithBlockInfo, error)
}
```

## Creating a New Account

To create a new account instance:

```go
package main

import (
    "fmt"
    
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/rpc"
)

func main() {
    // Initialize a provider
    provider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
    
    // Account details
    accountAddressHex := "0x..." // Your account address as hex string
    publicKey := "0x..."         // Your public key
    privateKeyHex := "0x..."     // Your private key as hex string
    
    // Convert address from hex string to felt.Felt
    accountAddress, err := new(felt.Felt).SetString(accountAddressHex)
    if err != nil {
        panic(err)
    }
    
    // Convert private key from hex string to big.Int
    privateKey, ok := new(big.Int).SetString(privateKeyHex, 0)
    if !ok {
        panic("Error converting private key to big.Int")
    }
    
    // Create a keystore for the account
    keystore := account.NewMemKeystore()
    if err := keystore.Put(publicKey, privateKey); err != nil {
        panic(err)
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
        panic(err)
    }
    
    fmt.Println("Account created successfully!")
}
```

## Precomputing Account Address

Before deploying an account, you can precompute its address:

```go
// Precompute the account address
accountAddress, err := account.PrecomputeAddress(
    provider,
    publicKey,
    salt,
    classHash,
    constructorCalldata,
)
if err != nil {
    panic(err)
}

fmt.Printf("Precomputed account address: 0x%s\n", accountAddress.Text(16))
```

## Deploying an Account

To deploy an account contract:

```go
// Prepare deployment parameters
salt, err := new(felt.Felt).SetString("0x01")
if err != nil {
    panic(err)
}

classHash, err := new(felt.Felt).SetString("0x123...") // Class hash of the account contract
if err != nil {
    panic(err)
}

// Convert constructor calldata to []*felt.Felt
constructorCalldata := []*felt.Felt{
    internalUtils.HexToFelt("0x456..."), // Public key
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
```

## Executing Transactions

To execute a transaction (invoke a contract function):

```go
// Create a function call
contractAddressFelt, err := new(felt.Felt).SetString("0x...")
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
```

## Declaring Contracts

To declare a new contract class:

```go
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
```

## Estimating Fees

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

## Signing Messages

To sign a message:

```go
// Message to sign
message := "0x..."

// Sign the message
signature, err := acc.SignMessage(context.Background(), []string{message})
if err != nil {
    panic(err)
}

fmt.Printf("Signature: %v\n", signature)
```

## Verifying Signatures

To verify a signature:

```go
// Verify the signature
isValid, err := acc.VerifyMessageSignature(context.Background(), []string{message}, signature)
if err != nil {
    panic(err)
}

if isValid {
    fmt.Println("Signature is valid")
} else {
    fmt.Println("Signature is invalid")
}
```

## Keystore Interface

StarkNet.go provides a keystore interface for managing private keys:

```go
type Keystore interface {
    Sign(ctx context.Context, id string, msgHash *big.Int) (x *big.Int, y *big.Int, err error)
}
```

### Memory Keystore

The memory keystore stores keys in memory:

```go
// Create a memory keystore
keystore := account.NewMemKeystore()

// Add a key pair
keystore.Put(publicKey, privateKey)

// Get a private key
privateKey, err := keystore.Get(publicKey)
if err != nil {
    panic(err)
}
```

### File Keystore

The file keystore stores keys in encrypted files:

```go
// Create a file keystore
keystore, err := account.NewFileKeystore("/path/to/keystore", "password")
if err != nil {
    panic(err)
}

// Add a key pair
keystore.Put(publicKey, privateKey)

// Get a private key
privateKey, err := keystore.Get(publicKey)
if err != nil {
    panic(err)
}
```

## Account Factory

The account factory provides a convenient way to create accounts:

```go
// Create an account factory
factory := account.NewAccountFactory(provider, keystore)

// Create a new account
acc, err := factory.Create(accountAddress, publicKey, 1)
if err != nil {
    panic(err)
}
```
