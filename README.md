<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/docs/public/starknetgo_vertical_dark.png">
    <img src="docs/docs/public/starknetgo_vertical_dark.png" height="256">
  </picture>
</div>

<p align="center">
    <a href="https://pkg.go.dev/github.com/NethermindEth/starknet.go">
        <img src="https://pkg.go.dev/badge/github.com/NethermindEth/starknet.go.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/nethermindeth/starknet.go/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-black">
    </a>
    <a href="https://github.com/NethermindEth/starknet.go/actions/workflows/test_testnet.yml">
        <img src="https://github.com/NethermindEth/starknet.go/actions/workflows/test_testnet.yml/badge.svg" alt="Main tests">
    </a>
    <a href="https://github.com/NethermindEth/starknet.go/stargazers">
      <img src="https://img.shields.io/github/stars/nethermindeth/starknet.go?style=social"/>
    </a>
</p>

</p>
<p align="center">
  <a href="https://x.com/NethermindStark">
    <img src="https://img.shields.io/twitter/follow/NethermindStark?style=social"/>
  </a>
  <a href="https://discord.com/invite/PaCMRFdvWT">
    <img src="https://img.shields.io/discord/629004402170134531?label=Nethermind%20Discord&logo=discord">
  </a>  
  <a href="https://t.me/StarknetGo">
    <img src="https://img.shields.io/badge/Starknet.go_group-gray?logo=telegram">
  </a>
</p>

<h1 align="center">Get the gopher Starkpilled</h1>

<a href="https://pkg.go.dev/github.com/NethermindEth/starknet.go">
<img src="https://img.shields.io/badge/Documentation-Website-yellow"
 height="50" />
</a>
<br><br>

**starknet.go** will get your golang backends and WASM frontends to Starknet easily.
With simple abstractions for the Starknet RPC, account management and common
operations on the wallets. The package has excellent documentation for a smooth
0 to 1 experience.

# 🌟 Features

- Seamless interaction with the Starknet RPC
- Tight integration with Juno
- Account management: Deploy accounts easily
- Good concurrency support

# Getting Started

- library documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/NethermindEth/starknet.go).
- [simple call example](./examples/simpleCall) to make a contract call to a testnet contract.
- [deploy account example](./examples/deployAccount) to deploy a new account contract on testnet.
- [invoke transaction example](./examples/invoke) to add a new invoke transaction on testnet.
- [declare transaction example](./examples/simpleDeclare) to add a new contract on testnet.
- [deploy contract UDC example](./examples/deployContractUDC) to deploy an ERC20 token using [UDC (Universal Deployer Contract)](https://docs.openzeppelin.com/contracts-cairo/1.0.0/udc) on testnet.
- [typed data example](./examples/typedData) to sign and verify a typed data.
- [websocket example](./examples/websocket) to learn how to subscribe to WebSocket methods.

Check [here](https://github.com/NethermindEth/starknet.go/tree/main/examples) for some FAQ.


### RPC

`starknet.go` RPC implements the Starknet [RPC v0.9.0 spec](https://github.com/starkware-libs/starknet-specs/releases/tag/v0.9.0)

| Method                                     | Implemented (*)    |
| ------------------------------------------ | ------------------ |
| `starknet_getBlockWithReceipts`            | :heavy_check_mark: |
| `starknet_getBlockWithTxHashes`            | :heavy_check_mark: |
| `starknet_getBlockWithTxs`                 | :heavy_check_mark: |
| `starknet_getStateUpdate`                  | :heavy_check_mark: |
| `starknet_getStorageAt`                    | :heavy_check_mark: |
| `starknet_getTransactionByHash`            | :heavy_check_mark: |
| `starknet_getTransactionByBlockIdAndIndex` | :heavy_check_mark: |
| `starknet_getTransactionReceipt`           | :heavy_check_mark: |
| `starknet_getTransactionStatus`            | :heavy_check_mark: |
| `starknet_getClass`                        | :heavy_check_mark: |
| `starknet_getClassHashAt`                  | :heavy_check_mark: |
| `starknet_getClassAt`                      | :heavy_check_mark: |
| `starknet_getBlockTransactionCount`        | :heavy_check_mark: |
| `starknet_call`                            | :heavy_check_mark: |
| `starknet_estimateFee`                     | :heavy_check_mark: |
| `starknet_estimateMessageFee`              | :heavy_check_mark: |
| `starknet_blockNumber`                     | :heavy_check_mark: |
| `starknet_blockHashAndNumber`              | :heavy_check_mark: |
| `starknet_chainId`                         | :heavy_check_mark: |
| `starknet_syncing`                         | :heavy_check_mark: |
| `starknet_getEvents`                       | :heavy_check_mark: |
| `starknet_getNonce`                        | :heavy_check_mark: |
| `starknet_addInvokeTransaction`            | :heavy_check_mark: |
| `starknet_addDeclareTransaction`           | :heavy_check_mark: |
| `starknet_addDeployAccountTransaction`     | :heavy_check_mark: |
| `starknet_traceTransaction`                | :heavy_check_mark: |
| `starknet_simulateTransaction`             | :heavy_check_mark: |
| `starknet_specVersion`                     | :heavy_check_mark: |
| `starknet_traceBlockTransactions`          | :heavy_check_mark: |
| `starknet_getStorageProof`                 | :heavy_check_mark: |
| `starknet_getMessagesStatus`               | :heavy_check_mark: |
| `starknet_getCompiledCasm`                 | :heavy_check_mark: |

#### WebSocket Methods

| Method                                     | Implemented (*)    |
| ------------------------------------------ | ------------------ |
| `starknet_subscribeEvents`                 | :heavy_check_mark: |
| `starknet_subscribeNewHeads`               | :heavy_check_mark: |
| `starknet_subscribeNewTransactions`        | :heavy_check_mark: |
| `starknet_subscribeNewTransactionReceipts` | :heavy_check_mark: |

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

## 📖 License

This project is licensed under the **MIT license**.

See [LICENSE](LICENSE) for more information.

Happy coding! 🎉

<!--- TODO: set it to the docs
## 🤝 Contribute

We're always looking for passionate developers to join our community and
contribute to `starknet.go`. Check out our [contributing guide](./docs/CONTRIBUTING.md)
for more information on how to get started.
-->

## Contributors ✨

Thanks goes to these wonderful people

<a href="https://github.com/NethermindEth/starknet.go/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=NethermindEth/starknet.go" />
</a>
<!-- Made with [contrib.rocks](https://contrib.rocks). -->
