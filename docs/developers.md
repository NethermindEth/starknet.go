# Developer Guide

This guide explains the starknet.go project structure and development workflow.

## Project Structure

```
starknet.go/
├── account/       # Account management
├── client/        # RPC client infrastructure
├── contracts/     # Contract utilities
├── curve/         # Cryptographic operations
├── devnet/        # Devnet integration
├── docs/          # Documentation (Vocs-based)
├── examples/      # Usage examples
├── hash/          # Hashing utilities
├── internal/      # Internal utilities
├── merkle/        # Merkle tree operations
├── paymaster/     # Paymaster integration
├── rpc/           # RPC types and provider
├── typeddata/     # Typed data signing
└── utils/         # Utility functions
```

## Package Descriptions

### `rpc/`
Core RPC functionality:
- `Provider` - main RPC client for Starknet nodes
- Type definitions for blocks, transactions, events
- WebSocket support for subscriptions

### `account/`
Account management:
- Account creation and import
- Transaction signing
- Nonce management
- Multi-call support

### `client/`
Low-level client infrastructure:
- HTTP/WebSocket client
- Error handling (`rpcerr/`)
- Logging (`log/`)

### `utils/`
Utility functions:
- Felt conversions (`HexToFelt`, `FeltToHex`)
- Selector generation (`GetSelectorFromNameFelt`)
- Transaction utilities

### `hash/`
Cryptographic hashing:
- Pedersen hash
- Poseidon hash

### `typeddata/`
Typed data for off-chain signing (similar to EIP-712)

### `paymaster/`
Paymaster integration for sponsored transactions

### `examples/`
Working examples demonstrating SDK usage:
- `deployAccount/` - Account deployment
- `invoke/` - Transaction invocation
- `readEvents/` - Event filtering
- `simpleCall/` - Contract calls
- `websocket/` - WebSocket subscriptions

## Development Workflow

### 1. Setup
```bash
make install-deps
```

### 2. Make Changes
Edit code in relevant package.

### 3. Test
```bash
# Quick feedback with mock tests
make mock-test

# Full test against devnet
make devnet-test
```

### 4. Lint
```bash
make lint
```

### 5. Submit PR
Create a pull request with your changes.

## Testing Strategy

### Mock Tests
- Located alongside source files (`*_test.go`)
- Use mocked RPC responses
- Fast, no network required
- Run with: `make mock-test`

### Integration Tests
- Test against real Starknet nodes
- Require environment configuration
- Run with: `make devnet-test`, `make testnet-test`, `make mainnet-test`

### Test Data
- `testData/` folders contain fixtures
- JSON files with sample responses

## Makefile Commands

| Command                | Description                     |
|------------------------|---------------------------------|
| `make test`            | Run all tests                   |
| `make mock-test`       | Run mock tests only             |
| `make devnet-test`     | Run devnet integration tests    |
| `make testnet-test`    | Run testnet integration tests   |
| `make mainnet-test`    | Run mainnet integration tests   |
| `make lint`            | Run golangci-lint with auto-fix |
| `make install-deps`    | Install dev dependencies        |
| `make clean-testcache` | Clear Go test cache             |

## Dependencies

### Required
- Go 1.25+
- Make

### Development Tools
- golangci-lint v2.7.2 - Linting
- mockgen v0.6.0 - Mock generation

Install with:
```bash
make install-deps
```

## Environment Variables

For integration tests, create a `.env` file in `examples/internal/`:

```bash
# RPC URLs
INTEGRATION_BASE=https://your-rpc-url
WS_INTEGRATION_BASE=wss://your-ws-url
```

See existing examples in `examples/` for reference.
