# RPC Test Files

This directory contains test files for all RPC methods. Each test file can be run independently using:

```bash
cd /docs/tests
go run rpc/<category>/<test_file>.go
```

## Directory Structure

- `block/` - Block-related RPC methods (7 tests)
- `chain/` - Chain information methods (3 tests)
- `call/` - Contract call method (1 test)
- `fee/` - Fee estimation methods (2 tests)
- `events/` - Event query method (1 test)
- `trace/` - Transaction tracing methods (3 tests)
- `write/` - Write operations (3 tests)
- `transaction/` - Transaction query methods (5 tests)
- `contract/` - Contract state methods (6 tests)
- `other/` - Other methods (1 test)

## Prerequisites

1. Set up `.env` file in `/docs/tests/` with:
   ```
   STARKNET_RPC_URL=your_rpc_url_here
   ```

2. Run from `/docs/tests/` directory

## Example Usage

```bash
# Run block number test
go run rpc/block/block_number.go

# Run chain ID test
go run rpc/chain/chain_id.go
```

## Total Tests: 32
