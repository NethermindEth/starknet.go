---
sidebar_position: 6
---

# Typed Data Example

This example demonstrates how to work with typed data in StarkNet.go. Typed data is a structured format for signing messages, similar to Ethereum's EIP-712, which provides better security and user experience when signing messages.

## Overview

The typed data example shows how to:

1. Initialize a connection to a StarkNet provider
2. Load an existing account
3. Define typed data in JSON format
4. Parse the typed data
5. Get the message hash
6. Sign the typed data
7. Verify the signature

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the `.env` file
3. Uncomment and assign your account address to the `ACCOUNT_ADDRESS` variable in the `.env` file
4. Uncomment and assign your public key to the `PUBLIC_KEY` variable in the `.env` file
5. Uncomment and assign your private key to the `PRIVATE_KEY` variable in the `.env` file

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/typedData"
	"github.com/joho/godotenv"
)

func main() {
	// Load variables from '.env' file
	err := godotenv.Load("../.env")
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err))
	}

	// Get the RPC provider URL and account details from environment variables
	providerURL := os.Getenv("RPC_PROVIDER_URL")
	if providerURL == "" {
		panic("RPC_PROVIDER_URL environment variable is not set")
	}

	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	if accountAddress == "" {
		panic("ACCOUNT_ADDRESS environment variable is not set")
	}

	publicKey := os.Getenv("PUBLIC_KEY")
	if publicKey == "" {
		panic("PUBLIC_KEY environment variable is not set")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		panic("PRIVATE_KEY environment variable is not set")
	}

	// Initialize connection to provider
	provider, err := rpc.NewProvider(providerURL)
	if err != nil {
		panic(fmt.Sprintf("Error initializing provider: %s", err))
	}

	// Convert private key from hex string to big.Int
	privateKeyBigInt, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Error converting private key to big.Int")
	}

	// Create a keystore for the account
	keystore := account.NewMemKeystore()
	if err := keystore.Put(publicKey, privateKeyBigInt); err != nil {
		panic(fmt.Sprintf("Error putting key in keystore: %s", err))
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
		panic(fmt.Sprintf("Error creating account: %s", err))
	}

	// Define typed data in JSON format
	typedDataJSON := `{
		"types": {
			"StarkNetDomain": {
				"name": "string",
				"version": "felt",
				"chainId": "felt"
			},
			"Person": {
				"name": "string",
				"wallet": "felt"
			},
			"Mail": {
				"from": "Person",
				"to": "Person",
				"contents": "string"
			}
		},
		"primaryType": "Mail",
		"domain": {
			"name": "StarkNet Mail",
			"version": "1",
			"chainId": "1"
		},
		"message": {
			"from": {
				"name": "Alice",
				"wallet": "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
			},
			"to": {
				"name": "Bob",
				"wallet": "0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
			},
			"contents": "Hello, Bob!"
		}
	}`

	// Parse typed data
	td, err := typedData.TypedDataFromJSON([]byte(typedDataJSON))
	if err != nil {
		panic(fmt.Sprintf("Error parsing typed data: %s", err))
	}

	// Get the message hash
	messageHash, err := td.GetMessageHash(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Error getting message hash: %s", err))
	}

	fmt.Printf("Message hash: 0x%s\n", messageHash.Text(16))

	// Sign the typed data
	signature, err := acc.SignMessage(context.Background(), []string{messageHash.Text(16)})
	if err != nil {
		panic(fmt.Sprintf("Error signing message: %s", err))
	}

	fmt.Printf("Signature: %v\n", signature)

	// Verify the signature
	isValid, err := acc.VerifyMessageSignature(context.Background(), []string{messageHash.Text(16)}, signature)
	if err != nil {
		panic(fmt.Sprintf("Error verifying signature: %s", err))
	}

	if isValid {
		fmt.Println("Signature is valid")
	} else {
		fmt.Println("Signature is invalid")
	}

	// You can also get the struct hash and domain hash separately
	structHash, err := td.GetStructHash(context.Background(), "Mail", td.Message)
	if err != nil {
		panic(fmt.Sprintf("Error getting struct hash: %s", err))
	}

	fmt.Printf("Struct hash: 0x%s\n", structHash.Text(16))

	domainHash, err := td.GetDomainHash(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Error getting domain hash: %s", err))
	}

	fmt.Printf("Domain hash: 0x%s\n", domainHash.Text(16))
}
```

## Running the Example

To run this example:

1. Make sure you are in the "typedData" directory
2. Execute `go run main.go`

## Expected Output

```
Message hash: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
Signature: [0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210]
Signature is valid
Struct hash: 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
Domain hash: 0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
```

## Key Concepts

### Typed Data Structure

Typed data in StarkNet.go is represented by the `TypedData` struct:

```go
type TypedData struct {
    Types       map[string][]TypedDataField
    PrimaryType string
    Domain      map[string]interface{}
    Message     map[string]interface{}
}
```

### Typed Data JSON Format

Typed data can be defined in JSON format with the following structure:

```json
{
    "types": {
        "StarkNetDomain": {
            "name": "string",
            "version": "felt",
            "chainId": "felt"
        },
        "CustomType1": {
            "field1": "string",
            "field2": "felt"
        },
        "CustomType2": {
            "field1": "CustomType1",
            "field2": "string"
        }
    },
    "primaryType": "CustomType2",
    "domain": {
        "name": "Example App",
        "version": "1",
        "chainId": "1"
    },
    "message": {
        "field1": {
            "field1": "Value 1",
            "field2": "0x123"
        },
        "field2": "Value 2"
    }
}
```

### Parsing Typed Data

To parse typed data from JSON:

```go
td, err := typedData.TypedDataFromJSON([]byte(typedDataJSON))
```

### Getting the Message Hash

To get the message hash of typed data:

```go
messageHash, err := td.GetMessageHash(context.Background())
```

### Signing Typed Data

To sign typed data:

```go
signature, err := acc.SignMessage(context.Background(), []string{messageHash.Text(16)})
```

### Verifying Signatures

To verify a signature:

```go
isValid, err := acc.VerifyMessageSignature(context.Background(), []string{messageHash.Text(16)}, signature)
```

### Getting Struct and Domain Hashes

You can also get the struct hash and domain hash separately:

```go
structHash, err := td.GetStructHash(context.Background(), "CustomType", td.Message)
domainHash, err := td.GetDomainHash(context.Background())
```

## Use Cases for Typed Data

Typed data is useful for:

1. **Structured Message Signing**: Provides a clear structure for what the user is signing
2. **Off-chain Signatures**: Enables off-chain signing for various applications
3. **Meta-transactions**: Allows for gasless transactions where a relayer pays the gas
4. **Authentication**: Can be used for authentication purposes
5. **Multi-signature Wallets**: Facilitates coordination between multiple signers

## Next Steps

After understanding how to work with typed data, you can:

- Explore [WebSocket](./websocket.md) for real-time updates
- Learn how to integrate typed data signing into your dApp
- Implement meta-transactions using typed data
