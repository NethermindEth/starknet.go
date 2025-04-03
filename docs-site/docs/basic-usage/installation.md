---
sidebar_position: 1
---

# Installation

This guide will help you install and set up StarkNet.go in your Go project.

## Prerequisites

- Go 1.19 or higher
- Basic understanding of Go modules
- Basic understanding of StarkNet concepts

## Installing StarkNet.go

To install StarkNet.go, use the standard Go package installation command:

```bash
go get github.com/NethermindEth/starknet.go
```

This will download and install the latest version of StarkNet.go and its dependencies.

## Setting Up Your Project

Create a new Go module for your project if you haven't already:

```bash
mkdir my-starknet-project
cd my-starknet-project
go mod init github.com/yourusername/my-starknet-project
```

Then, create a simple `main.go` file to verify the installation:

```go
package main

import (
    "fmt"
    
    "github.com/NethermindEth/starknet.go/rpc"
)

func main() {
    fmt.Println("StarkNet.go is installed!")
    
    // Initialize a provider (this is just a check, it won't connect to any node)
    _, err := rpc.NewProvider("https://example.com")
    if err != nil {
        fmt.Printf("Error initializing provider: %v\n", err)
        return
    }
    
    fmt.Println("Provider initialized successfully!")
}
```

Run your program to verify the installation:

```bash
go run main.go
```

If everything is set up correctly, you should see the message "StarkNet.go is installed!" followed by "Provider initialized successfully!".

## Environment Setup

For most examples and real-world usage, you'll need to set up environment variables for connecting to StarkNet nodes. Create a `.env` file in your project root:

```
# StarkNet Node URLs
RPC_PROVIDER_URL=https://starknet-sepolia.infura.io/v3/YOUR_API_KEY
WS_PROVIDER_URL=wss://starknet-sepolia.infura.io/ws/v3/YOUR_API_KEY

# Account Information (for transaction examples)
ACCOUNT_ADDRESS=0x...
PUBLIC_KEY=0x...
PRIVATE_KEY=0x...
```

Make sure to replace the placeholder values with your actual API keys and account information.

## Next Steps

Now that you have StarkNet.go installed, you can:

- Learn about [Account Management](./account-management.md)
- Explore how to work with [Transactions](./transactions.md)
- Understand [Contract Interaction](./contract-interaction.md)
