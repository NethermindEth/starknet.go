<div align="center">
  <img src="docs/images/golang_starknet_repo_banner.png" height="256">
</div>

<p align="center">
    <a href="https://pkg.go.dev/github.com/NethermindEth/starknet.go">
        <img src="https://pkg.go.dev/badge/github.com/NethermindEth/starknet.go.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/nethermindeth/starknet.go/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-black">
    </a>
    <a href="https://github.com/nethermindeth/starknet.go/actions/workflows/test.yml">
        <img src="https://github.com/nethermindeth/starknet.go/actions/workflows/test.yml/badge.svg?branch=main" alt="test">
    </a>
    <a href="https://twitter.com/NethermindStark">
      <img src="https://img.shields.io/twitter/follow/NethermindStark?style=social"/>
    </a>
    <a href="https://github.com/nethermindeth/starknet.go">
      <img src="https://img.shields.io/github/stars/nethermindeth/starknet.go?style=social"/>
    </a>
</p>

<h1 align="center">Get the gopher high on StarkNet</h1>

<a href="https://pkg.go.dev/github.com/NethermindEth/starknet.go">
<img src="https://img.shields.io/badge/Documentation-Website-yellow"
 height="50" />
</a>

#### :warning: `starknet.go` is currently under active development and is experiencing a rearchitecture. It will undergo breaking changes.

`starknet.go` will get your golang backends and WASM frontends to Starknet easily.
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

- library documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/NethermindEth/starknet.go).
- [curve example](./examples/curve) initializing the StarkCurve for signing and verification
- [contract example](./examples/contract) for smart contract deployment and function call
- [account example](./examples/contract) for Account initialization and invocation call

### Run Examples

***starknet simpleCall***

```sh
cd examples/simpleCall
go mod tidy
go run main.go
```


### RPC

`starknet.go` RPC implements the [StarkNet RPC Spec](https://github.com/starkware-libs/starknet-specs):

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
contribute to `starknet.go`. Check out our [contributing guide](./docs/CONTRIBUTING.md)
for more information on how to get started.

## 📖 License

This project is licensed under the **MIT license**.

See [LICENSE](LICENSE) for more information.

Happy coding! 🎉
## Contributors ✨

Thanks goes to these wonderful people
([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/drspacemn"><img src="https://avatars.githubusercontent.com/u/16685321?v=4?s=100" width="100px;" alt="drspacemn"/><br /><sub><b>drspacemn</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=drspacemn" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/gregoryguillou"><img src="https://avatars.githubusercontent.com/u/10611760?v=4?s=100" width="100px;" alt="Gregory Guillou"/><br /><sub><b>Gregory Guillou</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=gregoryguillou" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/tarrencev"><img src="https://avatars.githubusercontent.com/u/4740651?v=4?s=100" width="100px;" alt="Tarrence van As"/><br /><sub><b>Tarrence van As</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=tarrencev" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/alex-sumner"><img src="https://avatars.githubusercontent.com/u/46249612?v=4?s=100" width="100px;" alt="Alex Sumner"/><br /><sub><b>Alex Sumner</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=alex-sumner" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/broody"><img src="https://avatars.githubusercontent.com/u/610224?v=4?s=100" width="100px;" alt="Yun"/><br /><sub><b>Yun</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=broody" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/rzmahmood"><img src="https://avatars.githubusercontent.com/u/35128199?v=4?s=100" width="100px;" alt="Zoraiz Mahmood"/><br /><sub><b>Zoraiz Mahmood</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=rzmahmood" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/LucasLvy"><img src="https://avatars.githubusercontent.com/u/70894690?v=4?s=100" width="100px;" alt="Lucas @ StarkWare"/><br /><sub><b>Lucas @ StarkWare</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=LucasLvy" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/coburn24"><img src="https://avatars.githubusercontent.com/u/29192260?v=4?s=100" width="100px;" alt="Coburn"/><br /><sub><b>Coburn</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=coburn24" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Larkooo"><img src="https://avatars.githubusercontent.com/u/59736843?v=4?s=100" width="100px;" alt="Larko"/><br /><sub><b>Larko</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=Larkooo" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/oxlime"><img src="https://avatars.githubusercontent.com/u/93354898?v=4?s=100" width="100px;" alt="oxlime"/><br /><sub><b>oxlime</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=oxlime" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://mxxn.io"><img src="https://avatars.githubusercontent.com/u/1372918?v=4?s=100" width="100px;" alt="Blaž Hrastnik"/><br /><sub><b>Blaž Hrastnik</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=archseer" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/FlorianRichardSMT"><img src="https://avatars.githubusercontent.com/u/110891350?v=4?s=100" width="100px;" alt="Florian"/><br /><sub><b>Florian</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=FlorianRichardSMT" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/greged93"><img src="https://avatars.githubusercontent.com/u/82421016?v=4?s=100" width="100px;" alt="greged93"/><br /><sub><b>greged93</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=greged93" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/jney"><img src="https://avatars.githubusercontent.com/u/747?v=4?s=100" width="100px;" alt="Jean-Sébastien Ney"/><br /><sub><b>Jean-Sébastien Ney</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=jney" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://runningbeta.io"><img src="https://avatars.githubusercontent.com/u/615877?v=4?s=100" width="100px;" alt="Kristijan Rebernisak"/><br /><sub><b>Kristijan Rebernisak</b></sub></a><br /><a href="https://github.com/NethermindEth/starknet.go/commits?author=krebernisak" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the
[all-contributors](https://github.com/all-contributors/all-contributors)
specification. Contributions of any kind welcome!
