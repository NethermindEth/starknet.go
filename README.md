<div align="center">
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
  <img src="docs/images/caigo-no-bg.png" height="256">
</div>

<p align="center">
    <a href="https://pkg.go.dev/github.com/dontpanicdao/caigo">
        <img src="https://pkg.go.dev/badge/github.com/dontpanicdao/caigo.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/nethermindeth/caigo/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-black">
    </a>
    <a href="https://github.com/nethermindeth/caigo/actions/workflows/test.yml">
        <img src="https://github.com/nethermindeth/caigo/actions/workflows/test.yml/badge.svg?branch=main" alt="test">
    </a>
    <a href="https://twitter.com/NethermindStark">
      <img src="https://img.shields.io/twitter/follow/NethermindStark?style=social"/>
    </a>
    <a href="https://github.com/nethermindeth/caigo">
      <img src="https://img.shields.io/github/stars/nethermindeth/caigo?style=social"/>
    </a>
</p>

<h1 align="center">Get the gopher high on StarkNet</h1>

<a href="https://pkg.go.dev/github.com/dontpanicdao/caigo">
<img src="https://img.shields.io/badge/Documentation-Website-yellow"
 height="50" />
</a>

#### :warning: `cai.go` is currently under active development and is experiencing a rearchitecture. It will undergo breaking changes.

`cai.go` will get your golang backends and WASM frontends to Starknet easily.
With simple abstractions for the Starknet RPC, account management and common
operations on the wallets. The package has excellent documentation for a smooth
0 to 1 experience.

# 🌟 Features

- Seamless interaction with the Starknet RPC
- Tight integration with Juno (uses the RPC types, resulting in updates and
  breaking changes landing quickly)
- Account management: Deploy accounts easily
- Good concurrency support

# Getting Started

- library documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/dontpanicdao/caigo).
- [curve example](./examples/curve) initializing the StarkCurve for signing and verification
- [contract example](./examples/contract) for smart contract deployment and function call
- [account example](./examples/contract) for Account initialization and invocation call

### Run Examples

***starkcurve***

```sh
cd examples/curve
go mod tidy
go run main.go
```

***starknet contract***

```sh
cd examples/contract
go mod tidy
go run main.go
```

***starknet account***

```sh
cd examples/account
go mod tidy
go run main.go
```

### RPC

`cai.go` RPC implements the [StarkNet RPC Spec](https://github.com/starkware-libs/starknet-specs):

| Method                                         | Implemented (*)    |
| ---------------------------------------------- | ------------------ |
| `starknet_getBlockByHash`                      | :heavy_check_mark: |
| `starknet_getBlockByNumber`                    | :heavy_check_mark: |
| `starknet_getTransactionByHash`                | :heavy_check_mark: |
| `starknet_getTransactionReceipt`               | :heavy_check_mark: |
| `starknet_getClass`                            | :heavy_check_mark: |
| `starknet_getClassHashAt`                      | :heavy_check_mark: |
| `starknet_getClassAt`                          | :heavy_check_mark: |
| `starknet_call`                                | :heavy_check_mark: |
| `starknet_blockNumber`                         | :heavy_check_mark: |
| `starknet_chainId`                             | :heavy_check_mark: |
| `starknet_syncing`                             | :heavy_check_mark: |
| `starknet_getEvents`                           | :heavy_check_mark: |
| `starknet_addInvokeTransaction`                | :heavy_check_mark: |
| `starknet_addDeployTransaction`                | :heavy_check_mark: |
| `starknet_addDeclareTransaction`               | :heavy_check_mark: |
| `starknet_estimateFee`                         | :heavy_check_mark: |
| `starknet_getBlockTransactionCountByHash`      | :heavy_check_mark: |
| `starknet_getBlockTransactionCountByNumber`    | :heavy_check_mark: |
| `starknet_getTransactionByBlockNumberAndIndex` | :heavy_check_mark: |
| `starknet_getTransactionByBlockHashAndIndex`   | :heavy_check_mark: |
| `starknet_getStorageAt`                        | :heavy_check_mark: |
| `starknet_getNonce`                            | :heavy_check_mark: |
| `starknet_getStateUpdate`                      | :heavy_check_mark: |
| *`starknet_traceBlockTransactions`             | :x:                |
| *`starknet_traceTransaction`                   | :x:                |

> (*) some methods are not implemented because they are not yet available
> from [eqlabs/pathfinder](https://github.com/eqlabs/pathfinder).

### Run Tests

```go
go test -v ./...
```

### Run RPC Tests

```go
go test -v ./rpc -env [mainnet|devnet|testnet|mock]
```

### Run Benchmarks

```go
go test -bench=.
```

### Compatibility and stability


## 🤝 Contribute

We're always looking for passionate developers to join our community and
contribute to `cai.go`. Check out our [contributing guide](./docs/CONTRIBUTING.md)
for more information on how to get started.

## 📖 License

This project is licensed under the **MIT license**.

See [LICENSE](LICENSE) for more information.

Happy coding! 🎉
## Contributors ✨

Thanks goes to these wonderful people

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/drspacemn"><img src="https://avatars.githubusercontent.com/u/16685321?v=4?s=100" width="100px;" alt="drspacemn"/><br /><sub><b>drspacemn</b></sub></a><br /><a href="https://github.com/NethermindEth/caigo/commits?author=drspacemn" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/gregoryguillou"><img src="https://avatars.githubusercontent.com/u/10611760?v=4?s=100" width="100px;" alt="Gregory Guillou"/><br /><sub><b>Gregory Guillou</b></sub></a><br /><a href="https://github.com/NethermindEth/caigo/commits?author=gregoryguillou" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->
([emoji key](https://allcontributors.org/docs/en/emoji-key)):

This project follows the
[all-contributors](https://github.com/all-contributors/all-contributors)
specification. Contributions of any kind welcome!
