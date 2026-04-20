# Contributing to starknet.go

Thank you for your interest in contributing to starknet.go!

## Prerequisites

- Go 1.25 or higher
- Make

## Setup Development Environment

1. Clone the repository:
```bash
git clone https://github.com/NethermindEth/starknet.go.git
cd starknet.go
```

2. Install development dependencies:
```bash
make install-deps
```

This installs:
- golangci-lint v2.7.2 - for linting
- mockgen v0.6.0 - for generating mocks

## Running Tests

### Test Environments

The SDK supports multiple test environments:

| Environment | Command             | Description                          |
|-------------|---------------------|--------------------------------------|
| Mock        | `make mock-test`    | Unit tests with mocked RPC responses |
| Devnet      | `make devnet-test`  | Tests against local devnet           |
| Testnet     | `make testnet-test` | Tests against Starknet testnet       |
| Mainnet     | `make mainnet-test` | Tests against Starknet mainnet       |
| All         | `make test`         | Run all tests                        |

### Running Specific Tests

```bash
# Run tests for specific package
go test -v ./rpc/...
go test -v ./account/...
go test -v ./utils/...

# Run with specific environment
go test -v ./rpc/... -env testnet
```

## Linting

Run linter before submitting PR:

```bash
make lint
```

This runs golangci-lint with auto-fix enabled. Configuration is in `.golangci.yaml`.

## Code Style

- Follow standard Go conventions
- Use gofmt for formatting (handled by linter)
- Add comments for exported functions
- Write unit tests for new functionality

## Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Commit with descriptive message
7. Push and create Pull Request

## Pull Request Guidelines

- Describe what the PR does
- Reference related issues
- Include test coverage for new code
- Ensure all CI checks pass

## Documentation

When updating documentation, follow the rules in `docs/DOCUMENTATION_RULES.md`:
- Always verify struct definitions against source code
- Include GitHub source links
- Test that code examples compile

## Questions?

Open an issue or reach out to maintainers.
