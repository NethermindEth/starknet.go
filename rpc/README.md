## RPC implementation

Caigo RPC implementation provides the RPC API to perform operations with 
StarkNet. It is currently being tested and maintained up-to-date with
Pathfinder and relies on [go-ethereum](github.com/ethereum/go-ethereum/rpc)
to provide the JSON RPC 2.0 client implementation.

If you need Caigo to support another API, open an issue on the project.

### Testing the RPC API

To test the RPC API, you should simply go the the rpc directory and run
`go test` like below:

```shell
cd rpc
go test -v .
```

We provide an additional `-env` flag to `go test` so that you can choose the
environment you want to test. For instance, if you plan to test with the
`testnet`, run:

```shell
cd rpc
go test -env testnet -v .
```

Supported environments are `mock`, `testnet` and `mainnet`. The support for
`devnet` is planned but might require some dedicated condition since it is empty. 

If you plan to specify an alternative URL to test the environment, you can set
the `INTEGRATION_BASE` environment variable. In addition, tests load `.env.${env}`,
and `.env` before relying on the environment variable. So for instanve if you want
the URL to change only for the testnet environment, you could add the line below
in `.env.testnet`:

```text
INTEGRATION_BASE=http://localhost:9546
```

### Test coverage

The table below shows the test coverage

The 
| Method                        | Implemented           | Test               | 
| ----------------------------- | --------------------- |--------------------|
| `starknet_getBlockByHash`     |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getBlockByNumber`   |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getTransactionByHash` |  :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getTransactionReceipt` | :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getClass`           |    :heavy_check_mark: |                :x: |
| `starknet_getClassHashAt`     |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getClassAt`         |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getCode`            |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_call`               |    :heavy_check_mark: |                :x: |
| `starknet_blockNumber`        |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_chainId`            |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_syncing`            |    :heavy_check_mark: | :heavy_check_mark: |
| `starknet_getEvents`          |    :heavy_check_mark: |                :x: |
| `starknet_addInvokeTransaction` |                 :x: |                :x: |
| `starknet_addDeployTransaction` |                 :x: |                :x: |
| `starknet_addDeclareTransaction` |                :x: |                :x: |
| `starknet_traceTransaction`   |                   :x: |                :x: |
| `starknet_traceBlockTransactions` |               :x: |                :x: |
| `starknet_getNonce`           |                   :x: |                :x: |
| `starknet_protocolVersion`    |                   :x: |                :x: |
| `starknet_pendingTransactions` |                  :x: |                :x: |
| `starknet_estimateFee`         |                  :x: |                :x: |
| `starknet_getBlockTransactionCountByHash` |       :x: |                :x: |
| `starknet_getBlockTransactionCountByNumber` |     :x: |                :x: |
| `starknet_getTransactionByBlockNumberAndIndex` |  :x: |                :x: |
| `starknet_getTransactionByBlockHashAndIndex` |    :x: |                :x: |
| `starknet_getStorageAt`        |                  :x: |                :x: |
| `starknet_getStateUpdateByHash` |                 :x: |                :x: |
| ----------------------------- | --------------------- |--------------------|
