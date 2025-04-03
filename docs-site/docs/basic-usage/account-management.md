---
sidebar_position: 2
---

# Account Management

StarkNet.go provides comprehensive account management capabilities, allowing you to create, deploy, and manage StarkNet accounts.

## Understanding StarkNet Accounts

In StarkNet, accounts are smart contracts that manage user identity and transaction signing. Unlike Ethereum, where accounts are derived directly from private keys, StarkNet accounts must be deployed as smart contracts.

## Creating a New Account

To create a new account, you first need to generate a key pair and then deploy an account contract:

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
    
    fmt.Printf("Private key: 0x%s\n", privateKey.Text(16))
    fmt.Printf("Public key: 0x%s\n", publicKey.Text(16))
    
    // Precompute the account address
    accountAddress, err := account.PrecomputeAddress(
        provider,
        publicKey.Text(16),
        big.NewInt(0), // Salt
        "0x...", // Class hash of the account contract
        []string{}, // Constructor calldata
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Account address: 0x%s\n", accountAddress.Text(16))
    
    // Note: At this point, you need to fund this address with ETH before deploying
    // You can use a faucet for testnet: https://starknet-faucet.vercel.app/
}
```

## Deploying an Account

After generating the key pair and precomputing the address, you need to fund the address with ETH and then deploy the account:

```go
// Continuing from the previous example...

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
deployTx, err := acc.Deploy(context.Background(), big.NewInt(0), "0x...", []string{})
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
```

## Using an Existing Account

If you already have an account deployed on StarkNet, you can use it with StarkNet.go:

```go
package main

import (
    "context"
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
    
    // Now you can use the account for transactions
    fmt.Println("Account loaded successfully!")
}
```

## Getting Account Nonce

Before sending transactions, you need to get the current nonce for your account:

```go
nonce, err := provider.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accountAddress)
if err != nil {
    panic(err)
}

fmt.Printf("Current nonce: %s\n", nonce.Text(16))
```

## Next Steps

Now that you understand account management, you can:

- Learn about [Transactions](./transactions.md)
- Explore [Contract Interaction](./contract-interaction.md)
