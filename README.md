<h1 align="center">Golang Library for StarkNet</h1>

<p align="center">
    <a href="https://pkg.go.dev/github.com/dontpanicdao/caigo">
        <img src="https://pkg.go.dev/badge/github.com/dontpanicdao/caigo.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/dontpanicdao/caigo/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-black">
    </a>
    <a href="https://starkware.co/">
        <img src="https://img.shields.io/badge/powered_by-StarkWare-navy">
    </a>
</p>

Caigo is an MIT-licensed Go library for interacting with [StarkNet](https://docs.starknet.io/docs/intro).

### Getting Started

- library documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/dontpanicdao/caigo).
- [example](./examples/starkcurve) initializing the StarkCurve for signing and verification
- [example](./examples/starknet) for StarkNet interactions including deploy, call, invoke, and poll transaction

### Compatibility and stability

Caigo is currently under active development and will under go breaking changes until the initial stable(v1.0.0) release. The example directories and *_test.go files should always be applicable for the latest commitment on the main branch.
*NOTE: examples and tests may be out of sync with tagged versions and pkg.go.dev documentation*

### RPC

Caigo RPC implements the [StarkNet RPC Spec](https://github.com/starkware-libs/starknet-specs).

Implementation status:

| Method                        | Implemented           |
| ----------------------------- | --------------------- |
| `starknet_getBlockByHash` | :heavy_check_mark: |
| `starknet_getBlockByNumber` | :heavy_check_mark: |
| `starknet_getTransactionByHash` | :heavy_check_mark: |
| `starknet_getTransactionReceipt` | :heavy_check_mark: |
| `starknet_getClass` | :heavy_check_mark: |
| `starknet_getClassHashAt` | :heavy_check_mark: |
| `starknet_getClassAt` | :heavy_check_mark: |
| `starknet_call` | :heavy_check_mark: |
| `starknet_blockNumber` | :heavy_check_mark: |
| `starknet_chainId` | :heavy_check_mark: |
| `starknet_syncing` | :heavy_check_mark: |
| `starknet_getEvents` | :heavy_check_mark: |
| `starknet_addInvokeTransaction` | :x: |
| `starknet_addDeployTransaction` | :heavy_check_mark: |
| `starknet_addDeclareTransaction` | :x: |
| `starknet_traceTransaction` | :x: |
| `starknet_traceBlockTransactions` | :x: |
| `starknet_getNonce` | :x: |
| `starknet_protocolVersion` | :x: |
| `starknet_pendingTransactions` | :x: |
| `starknet_estimateFee` | :x: |
| `starknet_getBlockTransactionCountByHash` | :heavy_check_mark: |
| `starknet_getBlockTransactionCountByNumber` | :heavy_check_mark: |
| `starknet_getTransactionByBlockNumberAndIndex` | :heavy_check_mark: |
| `starknet_getTransactionByBlockHashAndIndex` | :heavy_check_mark: |
| `starknet_getStorageAt` | :heavy_check_mark: |
| `starknet_getStateUpdateByHash` | :x: |

### Run Examples

starkcurve

```sh
cd examples/starkcurve
go mod tidy
go run main.go
```

starknet

```sh
cd examples/starknet
go mod tidy

python3.7 -m venv ~/cairo_venv; source ~/cairo_venv/bin/activate; export STARKNET_NETWORK=alpha-goerli
starknet-compile ../../gateway/contracts/counter.cairo --output counter_compiled.json
go run main.go
```

### Run Tests

```go
go test -v ./...
```

### Run Benchmarks

```go
go test -bench=.
```

## Issues

If you find an issue/bug or have a feature request please submit an issue here
[Issues](https://github.com/dontpanicdao/caigo/issues)

## Contributing

If you are looking to contribute, please head to the
[Contributing](https://github.com/dontpanicdao/caigo/blob/main/CONTRIBUTING.md) section.
