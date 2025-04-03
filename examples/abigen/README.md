# Cairo ABI Source Code Generation Example

This example demonstrates how to use the `abigen` tool to generate Go bindings for Cairo contracts and interact with them.

## Prerequisites

- Go 1.19 or later
- A StarkNet node (testnet, mainnet, or local)

## Steps

### 1. Install the abigen tool

```bash
go install github.com/NethermindEth/starknet.go/abigen/cmd@latest
```

### 2. Generate bindings from a Cairo ABI

```bash
abigen --abi=contract.json --pkg=main --out=contract.go --type=SimpleContract
```

### 3. Use the generated bindings

See the [main.go](./main.go) file for a complete example of how to use the generated bindings.

## Running the Example

```bash
go run main.go
```

This will:
1. Connect to a StarkNet node
2. Load the generated contract bindings
3. Call a view function on the contract
4. Send a transaction to the contract
5. Wait for the transaction to be confirmed
