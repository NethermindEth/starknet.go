# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/NethermindEth/starknet.go/compare/v0.15.0...HEAD) <!-- Update the version number on each new release -->
<!-- template to copy:
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->

### Added
- New `client/rpcerr` package for handling RPC errors.
- New `curve.SignFelts` function for signing messages with felt.Felt parameters.

### Changed
- Major refactoring in types/functions names to match the Go naming conventions.
  - All type/function names that contained `_` have been renamed to remove the underscore.
  - The `typedData` pkg was renamed to `typedata`.
  - Other renames in exported types/fields/variables:
    - `typedData.Domain.ChainId` -> `typedData.Domain.ChainID`
    - `rpc.SKIP_FEE_CHARGE` -> `rpc.SkipFeeCharge`
    - `rpc.SKIP_VALIDATE` -> `rpc.SkipValidate`
    - `account.Account.ChainId` -> `account.Account.ChainID`
    - Variables `hash.PREFIX_TRANSACTION`, `hash.PREFIX_DECLARE`, and `hash.PREFIX_DEPLOY_ACCOUNT` were renamed and
      are no longer exported.
    - Variable `contracts.PREFIX_CONTRACT_ADDRESS` was renamed and is no longer exported.
- The `rpc.RPCError` type and logic was refactored and moved to the new `client/rpcerr` package.
  There are some changes in the new package:
  - The internal `tryUnwrapToRPCErr` func of the `rpc` pkg was renamed to `UnwrapToRPCErr` and moved to the new package.
  - The `Err` function now have a specific case for the `InternalError` code.

### Fixed
- The `typedData.TypedData` was not being marshaled exactly as it is in the original JSON. Now, the original JSON is preserved,
  so the output of `TypedData.MarshalJSON()` is exactly as the original JSON.
- Wrong encoding of the `selector` type in the `typedData` pkg for a specific case, when the value was already a hashed selector.
  More details in the PR [793](https://github.com/NethermindEth/starknet.go/pull/793).
- Not using the provided `context.Context` in the `account.Nonce` method when calling the `rpc.Nonce` method.

### Dev updates
- New `internal/tests/jsonrpc_spy.go` file containing a `Spy` type for spying JSON-RPC calls in tests. The
old `rpc/spy_test.go` file was removed.
- New `mocks/mock_client.go` file containing a mock of the `client.Client` type (`client.ClientI` interface).
- New benchmarks and tests for the `typedData` pkg.
- New linter rules in the `.golangci.yaml` file, thus, a lot of changes in the codebase to fix the new rules.

## [0.15.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.15.0) - 2025-09-03
### Changed
- The following functions now return the response as a value instead of a pointer:
  - `rpc.AddInvokeTransaction`
  - `rpc.AddDeclareTransaction`
  - `rpc.AddDeployAccountTransaction`
  - `account.BuildAndSendInvokeTxn`
  - `account.BuildAndSendDeclareTxn`
  - `account.DeployContractWithUDC`
  - `account.SendTransaction`
- The `account.AccountInterface` and the `rpc.RpcProvider` were updated to reflect the new return types.
- Reorg events are now supported in all Starknet subscriptions.

### Fixed
- Bug when receiving reorg events from subscriptions, which was panicking in some cases.

### Dev updates
- Regenerated mocks for the `account` and `rpc` packages
- Tests updated accordingly

## [0.14.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.14.0) - 2025-08-15
### Added
- New WebSocket subscription endpoints:
  - `rpc.SubscribeNewTransactions`
  - `rpc.SubscribeNewTransactionReceipts`
- `rpc.SubscriptionBlockID` type for websocket subscriptions, which is a restricted version of `BlockID` that doesn't allow `pre_confirmed` or `l1_accepted` tags
- Helper methods for `SubscriptionBlockID`: `BlockID()`, `WithBlockNumber()`, `WithBlockHash()`, and `WithLatestTag()`
- `l1_accepted` and `pre_confirmed` block tags in the `rpc.BlockTag` type
- `rpc.ErrFeeBelowMinimum` and `rpc.ErrReplacementTransactionUnderpriced` rpc errors
- Types:
  - `rpc.EmittedEventWithFinalityStatus`: the return value for the `rpc.SubscribeEvents` endpoint
  - `rpc.PriceUnitWei` and `rpc.PriceUnitFri` enums: representing the `WEI` and `FRI` units
  - `rpc.SubNewTxnReceiptsInput`: input for the `rpc.SubscribeNewTransactionReceipts` endpoint
  - `rpc.SubNewTxnsInput`: input for the `rpc.SubscribeNewTransactions` endpoint
  - `rpc.TxnWithHashAndStatus`: return value for the `rpc.SubscribeNewTransactions` endpoint
  - `rpc.MessageFeeEstimation`: return value for the `rpc.EstimateMessageFee` endpoint
  - `rpc.FeeEstimationCommon`: common fields for `rpc.FeeEstimation` and `rpc.MessageFeeEstimation`

### Changed
- `pending` terminology replaced by `pre_confirmed` across RPC headers, statuses and documentation
- `rpc.SubscribeNewHeads` endpoint now accepts `rpc.SubscriptionBlockID` instead of `rpc.BlockID` parameter
- `rpc.SubscribeEvents` endpoint now accepts the `rpc.EmittedEventWithFinalityStatus` type as parameter
- `rpc.EstimateMessageFee` endpoint now returns `MessageFeeEstimation` instead of a `FeeEstimation` pointer
- Improved RPC error handling: only code `0` is considered invalid in `tryUnwrapToRPCErr`
- `rpc.TraceBlockTransactions` endpoint now checks for `pre_confirmed` tag in the `BlockID` parameter and returns an error if it is set
- Small change in the error returned by the `UnmarshalJSON` method of the `rpc.TxnExecutionStatus` and `rpc.TxnFinalityStatus` types
- New errors returned by the `rpc.AddInvokeTransaction`, `rpc.AddDeclareTransaction`, and `rpc.AddDeployAccountTransaction` endpoints
- RPCErrors now returns more data than before when the error is not a known RPC error
- Types:
  - `rpc.RpcProvider` interface: change in the `EstimateMessageFee` method return value
  - `rpc.WebsocketProvider` interface: endpoints added/removed
  - `rpc.ErrInvalidTransactionNonce`: new string `data` field
  - `rpc.EventSubscriptionInput`: new/renamed fields
  - `rpc.MessageStatus`: new/renamed fields
  - `rpc.FeePayment`: change in a field type
  - `rpc.FeePaymentUnit` enum: renamed to `rpc.PriceUnit`
  - `rpc.TxnStatus` enum: new values
  - `rpc.StateUpdateOutput`: changed field type
  - `rpc.PendingStateUpdate`: renamed to `rpc.Pre_confirmedStateUpdate`
  - `rpc.FeeEstimation`: new/changed fields
  - `rpc.TxnFinalityStatus` enum: new value

### Removed
- `rpc.(*WsProvider).SubscribePendingTransactions` endpoint in favor of `rpc.SubscribeNewTransactions`
- `REJECTED` block tag in the `rpc.BlockStatus` type
- `rpc.checkForPre_confirmed` function calls in websocket methods since `rpc.SubscriptionBlockID` type prevents invalid tags at the type level
- `rpc.TxnStatus` enum: removed `REJECTED` value
- Types:
  - `rpc.SubPendingTxnsInput`
  - `rpc.PendingTxn`

### Fixed
- Wrong displayed versions in the warning message when using a different RPC version than the one implemented by starknet.go
- Error when subscribing to the `rpc.SubscribeEvents` and `rpc.SubscribeNewHeads` endpoints with an empty `rpc.SubscriptionBlockID` on Pathfinder node

### Dev updates
- Big refactor in the tests organization
- Add integration environment for testing
- A lot of tests refactorings and updates
- Removed `.vscode/launch.json` file
- Expanded integration tests and datasets; CI and `Makefile` targets updated to include integration runs

## [0.13.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.13.1) - 2025-08-05
### Fixed
- `rpc.Syncing` was crashing when the node response was that it was syncing.
- Removed wrong `omitempty` tags in the RPC error data types.

### Changed
- `rpc.Syncing` method now returns a `rpc.SyncStatus` struct instead of a pointer to it.
- `rpc.SyncStatus` type: `StartingBlockNum`, `HighestBlockNum` and `CurrentBlockNum` are now of type `uint64` instead of `NumAsHex`
- `rpc.SyncStatus` type: `SyncStatus` field was renamed to `IsSyncing`

### Removed
- `curve.g1Affline` variable, users should use new instances of `starkcurve.G1Affine` instead. It was causing bugs in the tests when running in parallel.

### Dev updates
- Renamed account/tests, rpc/tests, contracts/tests, hash/tests, and typedData/tests folders to testData
- Migrate internal/test.go file to the new internal/tests pkg
- New tests.TestEnv enum type representing test environments
- New tests.RunTestOn func for environment validation
- Updated all testing to use the new enum and the RunTestOn when necessary

## [0.13.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.13.0) - 2025-06-27
### Added
- `account` pkg
  - `Verify` method to the `Account` type and `AccountInterface` interface
  - `CairoVersion` type
  - `TxnOptions` type, allowing optional settings when building/estimating/sending a transaction with the Build* methods
  - `DeployContractWithUDC` method for deploying contracts using the Universal Deployer Contract (UDC)
- `utils` pkg
  - `TxnOptions` type, allowing optional settings when building a transaction (set tip and query bit version)
  - `BuildUDCCalldata` function to build calldata for UDC contract deployments
  - `PrecomputeAddressForUDC` function to compute contract addresses deployed with UDC
  - `UDCOptions` type and `UDCVersion` enum for configuring UDC deployments
- A warning message when calling `rpc.NewProvider` with a provider using a different RPC version than the one implemented by starknet.go.

### Removed
- `rpc.NewClient` function
- `rpc.SKIP_EXECUTE` tag

### Changed
- In `account` package
  - for the `BuildAndEstimateDeployAccountTxn`, `BuildAndSendInvokeTxn`, and `BuildAndSendDeclareTxn` methods
    - removed `multiplier` and `withQueryBitVersion` parameters
    - added `opts` parameter of type `*account.TxnOptions`
  - In `NewAccount` function, the `cairoVersion` parameter is now of type `account.CairoVersion` 
- In `utils` package, added `opts` parameter of type `*utils.TxnOptions` to the `BuildInvokeTxn`, `BuildDeclareTxn`, and `BuildEstimateDeployAccountTxn` methods
- `rpc.WithBlockTag` now accepts `BlockTag` instead of `string` as parameter
- `setup.GetAccountCairoVersion` now returns `account.CairoVersion` instead of `int`
- examples in `examples` folder updated to use `utils.TxnOptions` and `account.CairoVersion`
- Updated `examples/typedData/main.go` to use the new `Verify` method
- Updated `examples/deployContractUDC/main.go` to use the new `DeployContractWithUDC` method and UDC utilities

### Dev updates
- Added tests:
  - `utils.TestTxnOptions`
  - `rpc.TestVersionCompatibility`
  - `account_test.TestTxnOptions`
  - `account_test.TestVerify`
  - `account_test.TestDeployContractWithUDC` with comprehensive test cases for ERC20 and no-constructor deployments
  - `utils_test.TestBuildUDCCalldata` + others, covering various UDC scenarios and options
  - `utils_test.TestPrecomputeAddressForUDC` for origin-dependent and independent address computation
- RPC pkg:
  - Added "Warning" word in the logs when missing the .env file on `internal/test.go`
  - New `warnVersionCheckFailed` and `warnVersionMismatch` variables in `rpc/provider.go`
  - New `checkVersionCompatibility()` function in `rpc/provider.go` to check the version of the RPC provider. It is called inside `rpc.NewProvider`
  - `TestCookieManagement` modified to handle the `specVersion` call when creating a new provider
  - New `rpcVersion` constant in `rpc/provider.go`, representing the version of the RPC spec that starknet.go is compatible with
  - Updated `TestSpecVersion` to use the `rpcVersion` constant
- In `account.NewAccount` function, the `cairoVersion` parameter is now of type `account.CairoVersion` 
- New `testConfig` struct to the `account_test` package, for easier test configuration
- The `account` package has been refactored and split into multiple files.
- Regenerated mocks for Account interface
- `rpc.WithBlockTag` now accepts `BlockTag` instead of `string` as parameter
- Updated `examples/typedData/main.go` to use the new `Verify` method
- `setup.GetAccountCairoVersion` now returns `account.CairoVersion` instead of `int`

## What's Changed
* refactor: split account pkg into multiple files by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/750
* chore(deps): bump brace-expansion from 2.0.1 to 2.0.2 in /www in the npm_and_yarn group across 1 directory by @dependabot in https://github.com/NethermindEth/starknet.go/pull/754
* Thiagodeev/feat account verify by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/753
* Feature/version compatibility check by @RafieAmandio in https://github.com/NethermindEth/starknet.go/pull/725
* Thiagodeev/small random changes by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/755
* Fix: increase block range for websocket tests by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/757
* Add CODEOWNERS file by @nethermind-oss-compliance in https://github.com/NethermindEth/starknet.go/pull/761
* replace dead link README.md by @eeemmmmmm in https://github.com/NethermindEth/starknet.go/pull/762
* Add DeployContractWithUDC method to improve contract deployment experience. by @HACKER097 in https://github.com/NethermindEth/starknet.go/pull/760
* release v0.13.0 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/763

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.12.0...v0.13.0

## [0.12.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.12.0) - 2025-06-02
### Added
- `utils.FillHexWithZeroes` utility function

### Changed
A reestructuring on the entire `curve` pkg. See more in PR [744](https://github.com/NethermindEth/starknet.go/pull/744).

- `StarkCurve.GetRandomPrivateKey()` -> `GetRandomKeys()`
- `StarkCurve.GetYCoordinate()` -> `GetYCoordinate()`
- `StarkCurve.Sign()` -> `Sign()`
- `StarkCurve.PrivateToPoint()` -> `PrivateKeyToPoint()`
- `StarkCurve.Verify()` -> `Verify()`
- `VerifySignature()` -> `Verify()`

### Removed:
- `StarkCurve` type and all its methods
  - `StarkCurve.Add()`
  - `StarkCurve.Double()`
  - `StarkCurve.ScalarMult()`
  - `StarkCurve.ScalarBaseMult()`
  - `StarkCurve.IsOnCurve()`
  - `StarkCurve.InvModCurveSize()`
  - `StarkCurve.MimicEcMultAir()`
  - `StarkCurve.EcMult()`
  - `StarkCurve.SignFelt()`
  - `StarkCurve.GenerateSecret()`
  - `StarkCurve.GetYCoordinate()` -> func `GetYCoordinate()`
  - `StarkCurve.Verify()` -> func `Verify()`
  - `StarkCurve.Sign()` -> func `Sign()`
  - `StarkCurve.PrivateToPoint()` -> func `PrivateKeyToPoint()`
  - `StarkCurve.GetRandomPrivateKey()` -> func `GetRandomKeys()`
- `StarkCurvePayload` type
- `PedersenParamsRaw` variable
- `PedersenParams` variable
- `Curve` variable implementing the `StarkCurve`
- `VerifySignature()` function -> func `Verify()`
- `CurveOption` type
- `WithConstants` function
- `curve.DivMod` function
- `curve.FmtKecBytes` function
- `curve.MaskBits` function
- `pedersen_params.json` file

### What's Changed
* docs: update Cairo account package link by @dizer-ti in https://github.com/NethermindEth/starknet.go/pull/737
* Documentation Website by @WiseMrMusa in https://github.com/NethermindEth/starknet.go/pull/732
* Small bug on the starknet_getStorageProof method test by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/740
* Update mainnetBlockReceipts588763.json data by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/742
* chore(deps): bump vite from 6.3.3 to 6.3.5 in /www in the npm_and_yarn group across 1 directory by @dependabot in https://github.com/NethermindEth/starknet.go/pull/741
* Update some txn test data with new Juno data by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/743
* `curve` pkg migration to `gnark-crypto` implementation by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/744
* docs: add CHANGELOG.md by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/748


**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.11.1...v0.12.0

**Note: starting CHANGELOG.md from v0.12.0. All descriptions before v0.12.0 were taken directly from the Github releases description**

## [0.11.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.11.1) - 2025-05-15
This release fixes the bug when calling the `BuildAndSendDeclareTxn`, `BuildAndSendInvokeTxn`, and `BuildAndEstimateDeployAccountTxn` methods with the `withQueryBitVersion` parameter set to true.
The txn with the queryBit version was also being sent to the node instead of being only used when estimating the fees, which is incorrect.
This has been fixed.

### What's Changed
* Thiagodeev/fix sending querybit txn by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/733

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.11.0...v0.11.1

## [0.11.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.11.0) - 2025-05-12
This release implements Starknet RPC v0.8.1

Key changes:

- Multiple code refactorings in the transaction types (see PR [700](https://github.com/NethermindEth/starknet.go/pull/700) for more)
- Improve error handling in the `utils.SplitFactStr` utility

For contributors: new Makefile, and linter improvements

### What's Changed
* Thiagodeev/refactor txn types by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/700
* Fix: Add spaces before dashes in func descriptions by @devin-ai-integration in https://github.com/NethermindEth/starknet.go/pull/717
* `Makefile` and new linters for starknet.go by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/598
* Fix bug report of starknet.go by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/727

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.10.0...v0.11.0

## [0.10.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.10.0) - 2025-04-30
This release implements Starknet RPC v0.8.1

Key changes:
- new `withQueryBitVersion` param to the `BuildAndEstimateDeployAccountTxn`, `BuildAndSendInvokeTxn`, and `BuildAndSendDeclareTxn` methods
- small bug fix when unsubscribing from websocket

### What's Changed
* Update CallContext to CallContextWithSliceArgs in `requestUnsubscribe()` by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/705
* Hotfix/not estimating querybitversion txns by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/714


**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.9.0...v0.10.0

## [0.9.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.9.0) - 2025-04-24
This release implements Starknet RPC v0.8.1 ([ref](https://github.com/starkware-libs/starknet-specs/releases/tag/v0.8.1)).

Key changes:
- Go v1.23.4.
- Updated Juno dependency, fixing bug in `curve` pkg regarding a renamed method ( `NewFelt` > `New` ). Ref: https://github.com/NethermindEth/starknet.go/issues/669.
- A warning is now returned when instantiating a Braavos account in the `NewAccount` function since it doesn't support Starknet RPCv.8.0 (ref: https://community.starknet.io/t/starknet-devtools-for-0-13-5/115495#p-2359168-braavos-compatibility-issues-3).
- Some fields have been renamed (like `Account.AccountAddress` > `Account.Address`, `BlockHeader.BlockHash` > `BlockHeader.Hash`, etc).

Also, we now have a new `readEvents` in the `examples` folder teaching how to read Starknet events with Starknet.go (thanks to @alex-sumner).

### What's Changed
* Fix ByteArrFeltToString hex decoding error by @Hyodar in https://github.com/NethermindEth/starknet.go/pull/668
* Upgrade juno version to v0.14.0, and Go to 1.23.4 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/692
* Allow empty block_id, and update code to use it by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/695
* Remove fillEmptyFeeEstimation workaround as Juno issue is fixed by @devin-ai-integration in https://github.com/NethermindEth/starknet.go/pull/694
* Braavos fix by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/691
* Eventexample by @alex-sumner in https://github.com/NethermindEth/starknet.go/pull/593
* Refactor: Rename redundant struct field names by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/699
* chore(deps): bump golang.org/x/net from 0.36.0 to 0.38.0 in the go_modules group across 1 directory by @dependabot in https://github.com/NethermindEth/starknet.go/pull/701
* docs: remove rpc/README.md file by @zeroprooff in https://github.com/NethermindEth/starknet.go/pull/702
* ci: bump actions/checkout to v4 by @dizer-ti in https://github.com/NethermindEth/starknet.go/pull/703
* fix typo in `account.go` by @detrina in https://github.com/NethermindEth/starknet.go/pull/704
* Fix typos throughout the codebase by @devin-ai-integration in https://github.com/NethermindEth/starknet.go/pull/707

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.8.1...v0.9.0

## [0.8.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.8.1) - 2025-04-08
This release implements Starknet RPC v0.8.1 ([ref](https://github.com/starkware-libs/starknet-specs/releases/tag/v0.8.1)).

Key change: 
- WebSocket `subscription_id`s are now strings instead of integers.

### What's Changed
* Test slices functions by @estensen in https://github.com/NethermindEth/starknet.go/pull/496
* CI Refactor by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/686
* Bump the go_modules group across 1 directory with 2 updates by @dependabot in https://github.com/NethermindEth/starknet.go/pull/683
* fix: typos in documentation files by @kilavvy in https://github.com/NethermindEth/starknet.go/pull/676
* docs(readme): replaced the non-working badge with a working one by @operagxsasha in https://github.com/NethermindEth/starknet.go/pull/682
* Add TestRecursiveProofFixedSizeMerkleTree for recursive proof generation by @byksy in https://github.com/NethermindEth/starknet.go/pull/459
* v0.8.1 release changes by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/687

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.8.0...v0.8.1

## [0.8.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.8.0) - 2025-04-02
### Summary
This release fully implements Starknet RPC v0.8.0 ([ref](https://github.com/starkware-libs/starknet-specs/releases/tag/v0.8.0)).

Key features:
- Full Starknet RPC v0.8.0 compliance (new RPC and WebSocket methods).
- New helper methods for building and sending transactions in the `account` pkg.
- New utilities for building and handling transactions in the `utils` pkg.
- All examples updated to use V3 transactions and the new utilities for building and sending transactions.
- New `simpleDeclare` and `websocket` examples.
- New hash functions to hash all transaction versions.
- An important change in the SNIP-12 enum hashing. A new release called [`v0.8.0-snip12`](https://github.com/NethermindEth/starknet.go/releases/tag/v0.8.0-snip12) has been created specifically for this.

#### New RPC methods:
- `starknet_getCompiledCasm`
- `starknet_getMessagesStatus`
- `starknet_getStorageProof`

#### New RPC WebSocket methods:
- `starknet_subscribeEvents`
- `starknet_subscribeNewHeads`
- `starknet_subscribePendingTransactions`
- `starknet_subscribeTransactionStatus`

Detailed changes:
- Removed V0/1/2 broadcast txns. Now, only V3 transactions can be sent.
- Removed duplicates `rpc.Provider` methods in the `Account` type. They should now be called via the `acc.Provider` field.
- New `NewWebsocketProvider` function for creating a WebSocket provider.
- Transaction hash functions moved from the `account` pkg to the `hash` pkg.
- Contract class types moved from the `rpc` pkg to the `contracts` pkg.
- Bug fixed in the `CompiledClassHash` function (thanks to @baitcode).
- Improvements in the RPCError `data` field.
- New `client` pkg.
- New utils:
  - BuildInvokeTxn
  - BuildDeclareTxn
  - BuildDeployAccountTxn
  - WeiToETH/ETHToWei
  - FRIToSTRK/STRKToFRI
  - FeeEstToResBoundsMap
  - ResBoundsMapToOverallFee
  - UnmarshalJSONFileToType
- New hash functions (one for each version):
  - TransactionHashInvokeV0/V1/V3
  - TransactionHashDeclareV1/V2/V3
  - TransactionHashDeployAccountV1/V3

### What's Changed
* Thiagodeev/rpcv08 write methods, blockHeader and error 53 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/626
* Creates a description for the multi-call feature by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/648
* Thiagodeev/rpcv08 final changes by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/644
* New rpcdata + RPCv08 errors by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/649
* docs: changed the link to shield stars by @Olexandr88 in https://github.com/NethermindEth/starknet.go/pull/652
* Update README.md to include typedData example by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/653
* Ensure starknet_call calldata is not null by @Hyodar in https://github.com/NethermindEth/starknet.go/pull/658
* fix outdated documentation links by @youyyytrok in https://github.com/NethermindEth/starknet.go/pull/659
* chore: fix minor grammar and wording issues by @famouswizard in https://github.com/NethermindEth/starknet.go/pull/660
* chore: Fix Typo in Error Message for Unsupported Transaction Type by @0xwitty in https://github.com/NethermindEth/starknet.go/pull/661
* Add Test Case Names for Improved CI/CD Debugging by @AnkushinDaniil in https://github.com/NethermindEth/starknet.go/pull/663
* Add actions scanning and go community pack in codeql-analysis.yml by @yevh in https://github.com/NethermindEth/starknet.go/pull/664
* Fix contract class unmarshall for modern cairo compilers by @baitcode in https://github.com/NethermindEth/starknet.go/pull/656
* thiagodeev/rpcv08-websocket by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/651
* Thiagodeev/rpcv08 - implement/fix/update remaining tests for rpcv08 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/662
* Thiagodeev/rpcv08 refactor tests and utils organization by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/672
* Thiagodeev/rpcv08 new txn utilities by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/673
* Thiagodeev/rpcv08 update examples by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/677
* Thiagodeev/rpcv08 final adjustments by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/678
* Thiagodeev/rpcv08 final adjustments v2 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/680
* Thiagodeev/rpcv08 final adjustments v3 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/681
* Revert SNIP-12 code to be compatible with starknet.js by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/667

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.7.3...v0.8.0

## [0.7.3](https://github.com/NethermindEth/starknet.go/releases/tag/v0.7.3) - 2024-12-13
This version implements RPCv0.7.1

### What's Changed
* Implements SNIP-12 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/637
* New "ValidateSignature" function in 'curve' pkg by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/637
* Fix trace/simulation tests for sepolia by @wojciechos in https://github.com/NethermindEth/starknet.go/pull/639
* Add Dependency review workflow by @yevh in https://github.com/NethermindEth/starknet.go/pull/643
* General "SendTransaction" method for 'account' pkg by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/647

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.7.2...v0.7.3

## [0.7.2](https://github.com/NethermindEth/starknet.go/releases/tag/v0.7.2) - 2024-10-07
This version implements RPCv0.7.1

### What's Changed
* Upgrades Juno and Go version by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/625
* Create codeql.yml by @yevh in https://github.com/NethermindEth/starknet.go/pull/627
* fix definition of invoke transaction trace by @dkatzan in https://github.com/NethermindEth/starknet.go/pull/631
* fix test validation bug by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/632
* Update CONTRIBUTING.md by @RuneRogue in https://github.com/NethermindEth/starknet.go/pull/633

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.7.1...v0.7.2

## [0.7.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.7.1) - 2024-09-02
### Summary
This release primarily replaces Goerli test data with Sepolia test data, improves test coverage, fixes several bugs related to RPC methods and the Pedersen hash function, adds support for Byte array serialization, refactors certain types for a better user experience (UX), updates examples, and implements RPCv0.7.1.

### What's Changed
* RPC methods to return error by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/542
* Allow configuring the RPC client by @archseer in https://github.com/NethermindEth/starknet.go/pull/547
* Update SpecVersion to stop returning InternalError for nil errors by @AryanGodara in https://github.com/NethermindEth/starknet.go/pull/549
* #545 Add golangci-lint workflow by @nsiregar in https://github.com/NethermindEth/starknet.go/pull/546
* only load .env for non-mock envs by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/559
* fix account testnet environment load by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/560
* set base by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/561
* provider_test.go upgraded to Sepolia by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/563
* RPC/block tests upgraded to sepolia, fixed not running correct RPC url passed from env flag on CLI bug by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/566
* Update transaction tests to sepolia by @AryanGodara in https://github.com/NethermindEth/starknet.go/pull/557
* Update account tests to sepolia by @AryanGodara in https://github.com/NethermindEth/starknet.go/pull/558
* Update block tests to Sepolia and add more tests by @AryanGodara in https://github.com/NethermindEth/starknet.go/pull/554
* Update rpc: Call, Contract, Event and mock tests by @AryanGodara in https://github.com/NethermindEth/starknet.go/pull/555
* Update README.md by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/568
* correct json tag price_in_fri by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/572
* add missing l1 transaction receipt field by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/577
* Upgrade examples to use Sepolia network by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/579
* Move PrecomputeAddress function to 'contract' package by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/574
* Migrate from test-go to stretchr by @mdantonio in https://github.com/NethermindEth/starknet.go/pull/582
* chore ::> method on the BroadcastInvokeTxnType interface to constrain it. by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/584
* chore::> Constrain BroadcastDeclare and BroadcastDeployAccount  by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/585
* Update README.md by @RV12R in https://github.com/NethermindEth/starknet.go/pull/587
* chore ::> Merging all TransactionReceipt types into a single type by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/588
* pedersen x range bug fix by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/595
* Adding support for Byte array serialization by @krishnateja262 in https://github.com/NethermindEth/starknet.go/pull/583
* `deployContractUDC` added in CI test. by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/606
* Reverted changes in `BroadcastDeclareTxnV3` by @PsychoPunkSage in https://github.com/NethermindEth/starknet.go/pull/604
* Update tests to Sepolia by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/589
* Add Juno Pederson hash by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/610
* Thiagodeev/fix unreturned tx hash by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/614
* Fixes before 7.1 release by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/617
* Thiagodeev/fixes before 7 1 release v2 by @thiagodeev in https://github.com/NethermindEth/starknet.go/pull/619
* update readme for 0.7.1 release by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/621

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.7.0...v0.7.1

## [0.7.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.7.0) - 2024-03-08
The main changes include upgrading the rpc and accounts packages to support the new Starknet RPC spec v0.7.0-rc2

### Added
* RPCv0.7-rc2 support
* New examples (estimateFee and deployAccountUDC)
* RPC handlers now return RPC errors

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.6.1...v0.7.0

## [0.6.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.6.1) - 2024-01-24
### What's Changed
* Fix FmtCallData by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/498
* Fix starknet_EstimateFee by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/500
* Add an example that fetches token balance of account by @ChiHaoLu in https://github.com/NethermindEth/starknet.go/pull/509

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.6.0...v0.6.1

## [0.6.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.6.0) - 2023-12-07
The main changes include upgrading the rpc and accounts pkg to support the new starknet rpc spec v0.6.0

Note: we still need to add e2e tests for v3 transactions. We decided to wait until Goerli supports v0.13 to reduce friction in generating tests. These tests will be added soon.

### What's Changed
* Rpcv06 trace file by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/471
* rpcv06 update exec resources by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/466
* rpcv06 update emitted event by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/473
* rpcv06 update receipts by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/464
* rpcv06 implement v3 transactions by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/462
* Update Felt.go by @Bitcoinnoobie in https://github.com/NethermindEth/starknet.go/pull/482
* Rpcv06 rv2 to rc4 by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/487
* Rianhughes/rpcv06 integration tests by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/485
* update readme by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/480
* rpcv06 update readme to rc-5 by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/494
* update readme for rpcv0.6.0 non rc by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/495


**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.5.0...v0.6.0

## [0.5.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.5.0) - 2023-11-02
Full support for Starknet [RPC v.0.5](https://github.com/starkware-libs/starknet-specs/tree/v0.5.0).

### :wrench: What's Changed
- Implement rpc05 `specVersion` method by @rianhughes in #416 
- Fix domain name in typed message by @fico308 in #417 
- Implement `GetTransactionStatus` method by @rianhughes in #418 
- Remove `pendingTransactions` method by @rianhughes in #419 
- Update `DeclareTxn` types by @rianhughes in #421 
- Update pending types by @rianhughes in #424 
- Update `traceBlockTransactions` method by @rianhughes in #428 
- Update Function Invocation type by @rianhughes and @cicr99 in #430 
- Update transaction trace types by @rianhughes and @cicr99 in #432 
- Update BlockHeader and PendingBlockHeader by @rianhughes and @cicr99 in #435 
-  Implement BroadcastedTxn types and update related methods by @rianhughes and @cicr99 in #437, #438, #454 
- Replace deadlinks by @jelilat in #439 
- Update go version by @cicr99 and @joshklop in #446, #447  
- Implement new example for adding an invoke transaction and update README and workflow accordingly by @Akashneelesh and @cicr99 in #450, #456, #457 
- Fix `Nonce` method return type by @cicr99 in #453 
- Handle errors by @cicr99 in #455
- Add code comments by @aquental in #415  

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.4.6...v0.5.0

### Note:
FmtCalldata() may be unstable.

## [0.4.6](https://github.com/NethermindEth/starknet.go/releases/tag/v0.4.6) - 2023-10-11
A "housekeeping" release including:
- Dedicated 'account' package
- Other library packages and code repository restructuring

IMPORTANT:exclamation::eyes:: This release does include breaking changes, and will require refactoring of application code which depends on **starknet.go** library.

### :wrench: What's Changed
* New account implementation by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/302
* Test Declare transaction by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/386
* Implement `WaitForTransactionReceipt` in accounts by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/389
* Implement `formatCallDAta` for Cairo2 by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/399
* Update examples to use account interface by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/414
* Remove deprecated packages by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/373
* Restructuring rpc tests folder by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/379
* Remove types package by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/384
* Create contracts package and hash logic by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/377
* Create curve package and organize remaining files by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/394, https://github.com/NethermindEth/starknet.go/pull/398
* Fix tests package and rename to `devnet` by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/395
* GH Action improvements by @stranger80 and @rianhughes in https://github.com/NethermindEth/starknet.go/pull/396, https://github.com/NethermindEth/starknet.go/pull/397, https://github.com/NethermindEth/starknet.go/pull/413 and https://github.com/NethermindEth/starknet.go/pull/362
* Fix typos by @omahs in https://github.com/NethermindEth/starknet.go/pull/363
* Update contract class `abi` to string by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/368

**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.4.5...v0.4.6

## [0.4.5](https://github.com/NethermindEth/starknet.go/releases/tag/v0.4.5) - 2023-09-27
Full support for Starknet [RPC v.0.4](https://github.com/starkware-libs/starknet-specs/tree/v0.4.0).

### What's Changed
* rpcv04 update trace and read errros by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/331
* rpcv04 update trace api by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/330
* Merge changes in Rpcv04 branch by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/338
* Implement DeclareTxn rpcv04 by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/288
* Update readme by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/339
* Rpcv04 tx status by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/289
* rpcv4_write_errors by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/332
* fix unknown type by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/343
* RPCv04 Remove BroadcastTx by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/290
* Modification to InvokeTxnV0 by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/340
* fix node error handling by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/352
* Block with tx hashes by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/354
* Create deployAccount example using rpc by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/349
* Readme rpcv04 update by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/357


**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.4.3...v0.4.5

## [0.4.4](https://github.com/NethermindEth/starknet.go/releases/tag/v0.4.4) - 2023-09-26
Full support for Starknet [RPC v.0.4](https://github.com/starkware-libs/starknet-specs/tree/v0.4.0). 

### What's Changed
* rpcv04 update trace and read errros by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/331
* rpcv04 update trace api by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/330
* Merge changes in Rpcv04 branch by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/338
* Implement DeclareTxn rpcv04 by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/288
* Update readme by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/339
* Rpcv04 tx status by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/289
* rpcv4_write_errors by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/332
* fix unknown type by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/343
* RPCv04 Remove BroadcastTx by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/290
* Modification to InvokeTxnV0 by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/340
* fix node error handling by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/352
* Block with tx hashes by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/354
* Create deployAccount example using rpc by @rianhughes in https://github.com/NethermindEth/starknet.go/pull/349
* Readme rpcv04 update by @cicr99 in https://github.com/NethermindEth/starknet.go/pull/357


**Full Changelog**: https://github.com/NethermindEth/starknet.go/compare/v0.4.3...v0.4.4

## [0.4.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.4.1) - 2022-10-24
A preview of go-starknet, the caigo cli to work with Starknet.

## [0.4.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.4.0) - 2022-10-23
This release **breaks** many of the 0.3.x API to support the `gateway` and the `rpc` v0.1 protocol together. It is being tested against [starknet-devnet](https://github.com/Shard-Labs/starknet-devnet), [eqlabs/pathfinder](https://github.com/eqlabs/pathfinder) and the standard [starknet gateway](https://github.com/starkware-libs/cairo-lang). It includes the following:

- support for Starknet v0.10.1 protocol on the `gateway`
- support for Starknet v0.10.0 protocol on the rpc v0.1 because it is not compatible with the latest release
- support for both V0 and V1 accounts for invoke transaction
- an account manager to install many account versions, including the [plugin account](https://github.com/argentlabs/starknet-plugin-account) and the [openzeppelin accounts](https://github.com/OpenZeppelin/cairo-contracts)
- an account interface and the associated implementation for both the gateway and rpc v0.1
- support for plugin extensions in the account and an implementation of the session key.
- some tools to work with devnet and transform help with converting data
- an implementation of a fixed-size merkle tree algorithm

**known issues**: due to the fact rpc v0.2 is not yet fully supported by Pathfinder and Devnet, the implementation is lacking. It remains a short term goal to implement it. Hopefully it will be part of v0.5.0.

## [0.3.1](https://github.com/NethermindEth/starknet.go/releases/tag/v0.3.1) - 2022-09-08
**EventParams:**

- nodes have ammended their `starknet_getEvents` params interface from `..."fromBlock": 800 ...` to `..."fromBlock": {"block_number": 250000}...` in order to support both `block_number` and `block_hash` values
- https://github.com/dontpanicdao/caigo/pull/104 addresses these changes

**BlockWithTxs:**
- in order for this release to still access blocks from https://github.com/starkware-libs/starknet-specs/releases/tag/v0.2.0 we've added a `BlockWithTxs` utility function
- please note `BlockWithTxHashes` and other RPC v0.2.0 functionality will be added in the caigo v0.4.0 release

## [0.3.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.3.0) - 2022-06-28
Gateway Package:
- Starknet client/provider implementation for feeder and gateway apis
- Endpoints/Types including: `estimate_fee`, `get_block`, `get_block_hash_by_id`, `get_id_by_hash`, `get_class_by_hash`, `get_class_hash_at`, `get_full_contract`, `call_contract`, `add_transaction`, `get_state_update`, `get_contract_addresses`, `get_storage_at`, `get_transaction`, `get_transaction_status`, `get_transaction_id_by_hash`, `get_transaction_receipt`, `get_transaction_trace`

RPC Package:
- RPC client implementing StarkNet RPC [Spec](https://github.com/starkware-libs/starknet-specs)
- Method implementation status included in README 

Types Package:
- Common provider types
- Post v0.3.0 `*big.Int` args/returns will be changed to `*types.Felt`s

StarkNet v0.9.0 [Support](https://medium.com/starkware/starknet-alpha-0-9-0-dce43cf13490):
- provider getters for class hashes and definitions
- estimate fee getters for transaction type
- Accounts with signature schemes adhering to SN `get_tx_info -> tx.hash`
- Accounts compatible w/ Argent [v0.9.0](https://github.com/argentlabs/argent-contracts-starknet/blob/54e5da0e7e2a69d1d56cef9bfe968dc6369a6579/contracts/ArgentAccount.cairo#L222)
- NOTE: devnet uses preloaded accounts based on [OZ v0.1.0](https://github.com/OpenZeppelin/cairo-contracts/releases/tag/v0.1.0) and signatures are incompatible with caigo Accounts

## [0.2.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.2.0) - 2022-04-01
Support for StarkNet actions:

- Gateway and Signer structures(similar to starknet.js)
- Setters - Deploy, invoke, execute transactions on StarkNet contracts(support for OpenZeppelin Account Multicall)
- Getters - transaction, transaction status, transaction receipt, storage, code, block
- Feeder Gateway poller for desired transaction status
- Contract call - post calldata to contract at desired entry point
- Typed Data and an encoding interface

Misc:

- proper randomness generation for private
- proper curve boundary checking

## [0.1.0](https://github.com/NethermindEth/starknet.go/releases/tag/v0.1.0) - 2022-01-10
Implements the crypto primitives for initiating and using various StarkNet functionality:

- Sign Message
- Verify Signature
- Hash array of elements
- Obtain Pedersen hash of elements
- Parse and hash starknet transaction in accordance with 'starknet.js'
- Initialize the Starknet Elliptic Curve
- Initialize the Starknet Elliptic Curve w/ Constant Points
- Elliptic Curve math primitives: Add, Double, DivMod, InvMod, EcMult, ScalarMult, ScalarBaseMult, MimicEcMultAir
