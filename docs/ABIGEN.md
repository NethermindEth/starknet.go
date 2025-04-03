# Cairo ABI Source Code Generation

The `abigen` package provides functionality to generate Go bindings for Cairo contracts, similar to how go-ethereum's abigen works for Solidity contracts.

## Overview

The Cairo ABI source code generation tool allows you to:

1. Generate type-safe Go bindings for Cairo 2 contracts
2. Create methods for calling contract functions with proper type conversion
3. Handle deployment and transaction creation
4. Provide event filtering and subscription

## Installation

The abigen CLI tool can be installed using:

```bash
go install github.com/NethermindEth/starknet.go/abigen/cmd@latest
```

## Usage

### Command Line Interface

```bash
abigen --abi=path/to/contract.json --pkg=mypackage --out=MyContract.go --type=MyContract
```

Options:
- `--abi`: Path to the Cairo contract ABI JSON file (required)
- `--bin`: Path to the Cairo contract bytecode (optional, for deploy methods)
- `--pkg`: Package name for the generated code (default: "main")
- `--type`: Contract struct name (default: derived from package or filename)
- `--out`: Output file for the generated code (default: stdout)

### Programmatic Usage

You can also use the abigen package programmatically in your Go code:

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NethermindEth/starknet.go/abigen/accounts/abi/abigen"
)

func main() {
	// Read ABI JSON
	abiBytes, err := os.ReadFile("path/to/contract.json")
	if err != nil {
		panic(err)
	}

	// Generate bindings
	code, err := abigen.BindCairoFixed(
		[]string{"MyContract"},     // Contract names
		[]string{string(abiBytes)}, // ABIs
		[]string{""},               // Bytecodes (optional)
		"mypackage",                // Package name
	)
	if err != nil {
		panic(err)
	}

	// Write to file
	err = os.WriteFile("MyContract.go", []byte(code), 0644)
	if err != nil {
		panic(err)
	}
}
```

## Type Mapping

Cairo types are mapped to Go types as follows:

| Cairo Type | Go Type |
|------------|---------|
| felt252 | *felt.Felt |
| u8, u16, u32 | uint32 |
| u64, u128 | uint64 |
| u256 | *big.Int |
| bool | bool |
| ContractAddress | *felt.Felt |
| Array<T> | []T |

## Generated Code Structure

The generated code includes:

1. A main contract struct with Caller, Transactor, and Filterer components
2. Methods for read-only calls (view functions)
3. Methods for state-changing transactions (external functions)
4. Event filtering and subscription methods
5. Helper methods for contract deployment

Example of generated code structure:

```go
// Main contract struct
type MyContract struct {
    MyContractCaller     // Read-only binding to the contract
    MyContractTransactor // Write-only binding to the contract
    MyContractFilterer   // Log filterer for contract events
}

// Read-only methods
func (_MyContract *MyContractCaller) GetBalance(opts *bind.CallOpts) (*felt.Felt, error) {
    // ...
}

// State-changing methods
func (_MyContract *MyContractTransactor) IncreaseBalance(opts *bind.TransactOpts, amount *felt.Felt) (*rpc.InvokeTxnResponse, error) {
    // ...
}
```

## Example

See the [abigen example](../examples/abigen/README.md) for a complete demonstration of using the generated bindings.
