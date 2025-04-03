---
sidebar_position: 1
---

# Getting Started with StarkNet.go

StarkNet.go is a comprehensive Go library that facilitates interaction with the StarkNet blockchain network. StarkNet is a permissionless, decentralized ZK-Rollup operating as a Layer 2 scaling solution for Ethereum.

## Overview

StarkNet.go enables Go backend applications and WASM frontends to seamlessly connect with StarkNet by providing abstractions for the StarkNet RPC, account management, and wallet operations.

The library serves developers building applications on StarkNet by offering:

- A complete implementation of the StarkNet RPC v0.8.0 specification
- Account creation, management, and transaction signing
- Transaction building and submission (invoke, declare, deploy)
- Contract interaction and querying
- Cryptographic operations for StarkNet (Stark curve, Pedersen hashing)
- WebSocket subscriptions for real-time updates

## Features

- **Seamless interaction with the StarkNet RPC**: Complete implementation of the StarkNet RPC v0.8.0 specification
- **Tight integration with Juno**: Works seamlessly with the Juno StarkNet sequencer implementation in Go
- **Account management**: Deploy and manage accounts easily
- **Good concurrency support**: Designed with Go's concurrency patterns in mind

## Installation

To install StarkNet.go, use the standard Go package installation command:

```bash
go get github.com/NethermindEth/starknet.go
```

## Quick Start

Here's a simple example to get you started with StarkNet.go:

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
    
    // Get the latest block number
    blockNumber, err := provider.BlockNumber(context.Background())
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Latest block number: %d\n", blockNumber)
}
```

## Next Steps

Explore the following sections to learn more about StarkNet.go:

- [Basic Usage](./basic-usage/installation.md): Learn how to install and set up StarkNet.go
- [API Reference](./api-reference/rpc.md): Detailed documentation of the StarkNet.go API
- [Examples](./examples/simple-call.md): Practical examples demonstrating StarkNet.go functionality
