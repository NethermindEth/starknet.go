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

# Download counter.cairo example
wget https://github.com/starknet-edu/ultimate-env/blob/main/counter.cairo
# or if your machine doesn't have wget
curl https://raw.githubusercontent.com/starknet-edu/ultimate-env/main/counter.cairo > counter.cairo

python3.7 -m venv ~/cairo_venv; source ~/cairo_venv/bin/activate; export STARKNET_NETWORK=alpha-goerli
starknet-compile counter.cairo --output counter_compiled.json --abi counter_abi.json
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
