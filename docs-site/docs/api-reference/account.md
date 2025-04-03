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

The `account.Account` interface defines the methods that all account implementations must provide:

```go
type Account interface {
    Address() string
    PublicKey() string
    ChainID() string
    Signer() Signer
    Deployer() Deployer
    Executor() Executor
    Declarer() Declarer
    Estimator() Estimator
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
    accountAddress := "0x..." // Your account address
    publicKey := "0x..."      // Your public key
    privateKey := "0x..."     // Your private key
    
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
// Deploy the account
deployTx, err := acc.Deploy(context.Background(), salt, classHash, constructorCalldata)
if err != nil {
    panic(err)
}

fmt.Printf("Deploy transaction hash: 0x%s\n", deployTx.TransactionHash.Text(16))
```

## Executing Transactions

To execute a transaction (invoke a contract function):

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

// Execute the transaction
tx, err := acc.Execute(context.Background(), []rpc.FunctionCall{functionCall}, nil)
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
    Get(pubKey string) (*big.Int, error)
    Put(pubKey string, privKey *big.Int) error
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

## Deployer Interface

The deployer interface defines methods for deploying accounts:

```go
type Deployer interface {
    Deploy(ctx context.Context, salt *big.Int, classHash string, constructorCalldata []string) (*rpc.DeployAccountTxnResponse, error)
}
```

## Executor Interface

The executor interface defines methods for executing transactions:

```go
type Executor interface {
    Execute(ctx context.Context, calls []rpc.FunctionCall, opts *ExecuteOpts) (*rpc.InvokeTxnResponse, error)
}
```

## Declarer Interface

The declarer interface defines methods for declaring contracts:

```go
type Declarer interface {
    Declare(ctx context.Context, contractClass *contracts.SierraContractClass, compiledClassHash *big.Int) (*rpc.DeclareTxnResponse, error)
}
```

## Estimator Interface

The estimator interface defines methods for estimating fees:

```go
type Estimator interface {
    EstimateFee(ctx context.Context, calls []rpc.FunctionCall, opts *ExecuteOpts) (*rpc.FeeEstimate, error)
}
```

## Signer Interface

The signer interface defines methods for signing messages:

```go
type Signer interface {
    SignMessage(ctx context.Context, message []string) ([]string, error)
    VerifyMessageSignature(ctx context.Context, message []string, signature []string) (bool, error)
}
```
