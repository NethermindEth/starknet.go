# Devnet Tests

Test files for the DevNet package.

## Prerequisites

These tests require a running Starknet DevNet instance.

### Starting DevNet

**Option 1: Using Docker** (Recommended)
```bash
docker run -p 5050:5050 shardlabs/starknet-devnet-rs:latest
```

**Option 2: Using Cargo**
```bash
cargo install starknet-devnet
starknet-devnet
```

**Option 3: Using Docker Compose**
```yaml
# docker-compose.yml
version: '3.8'
services:
  starknet-devnet:
    image: shardlabs/starknet-devnet-rs:latest
    ports:
      - "5050:5050"
    command: --seed 42
```

```bash
docker-compose up
```

## Configuration

Create a `.env` file in `/docs/tests/`:

```bash
DEVNET_URL=http://localhost:5050
```

If not set, defaults to `http://localhost:5050`.

## Running Tests

```bash
cd /path/to/starknet.go/docs/tests

# Run individual tests
go run devnet/is_alive.go
go run devnet/accounts.go
go run devnet/mint.go
go run devnet/fee_token.go
```

## Test Files

### 1. new_devnet.go
Creates a DevNet instance and verifies connection.

**Expected Output:**
```
DevNet instance created for: http://localhost:5050
✓ DevNet is running
```

### 2. is_alive.go
Checks if DevNet is running.

**Expected Output:**
```
DevNet Status: Running ✓
```

### 3. accounts.go
Retrieves pre-funded test accounts.

**Expected Output:**
```
Found 10 pre-funded accounts

First Account:
  Address:     0x64b48806902a367c8598f4f95c305e8c1a1acba5f082d294a43793113115691
  Public Key:  0x39d9e6ce352ad4530a0ef5d5a18fd3303c3606a7fa6ac5b620020ad681cc33b
  Private Key: 0xe1406455b7d66b1690803be066cbe5e
```

### 4. mint.go
Mints tokens to an address.

**Expected Output:**
```
Tokens minted successfully!
New Balance: 1001000000000000000000 WEI
Transaction Hash: 0x5f4f16b0df35ee7059a7e5d4f84bd70b0e5e0b8a9c2b74b8acab5e8ff5f88e48
```

### 5. fee_token.go
Gets fee token information.

**Expected Output:**
```
Fee Token Information:
  Symbol:  ETH
  Address: 0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7
```

## Notes

- All tests use environment variables (no hardcoded URLs)
- Tests will fail gracefully if DevNet is not running
- Default DevNet provides 10 pre-funded accounts
- Each account starts with 1,000 ETH (1e21 wei)
- DevNet runs on localhost:5050 by default

## Troubleshooting

### DevNet not running
```
Error: DevNet Status: Not Running ✗
Solution: Start DevNet with one of the methods above
```

### Connection refused
```
Error: connection refused
Solution: Check if DevNet is running on the correct port
```

### Port already in use
```
Solution: Use a different port:
docker run -p 5051:5050 shardlabs/starknet-devnet-rs:latest

Then set DEVNET_URL=http://localhost:5051
```

## Integration Testing

These test files can be used in CI/CD pipelines:

```yaml
# .github/workflows/test.yml
name: Test
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      devnet:
        image: shardlabs/starknet-devnet-rs:latest
        ports:
          - 5050:5050
    steps:
      - uses: actions/checkout@v2
      - name: Run DevNet tests
        run: |
          cd docs/tests
          go run devnet/is_alive.go
          go run devnet/accounts.go
```

## Related Documentation

- [DevNet Package Documentation](../docs/pages/docs/devnet/)
- [DevNet Methods](../docs/pages/docs/devnet/methods/)
- [RPC Tests](../rpc/) - Tests for RPC methods
