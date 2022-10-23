# CHANGELOG

## v0.4.0 - Support for Starknet v0.10.1 and RPC v0.1

This release **breaks** many of the 0.3.x API to support the `gateway` and the
`rpc` v0.1 protocol together. It is being tested against
[starknet-devnet](https://github.com/Shard-Labs/starknet-devnet),
[eqlabs/pathfinder](https://github.com/eqlabs/pathfinder) and the standard
[starknet gateway](https://github.com/starkware-libs/cairo-lang). It includes
the following features:

- Support for Starknet v0.10.1 protocol on the gateway
- Support for Starknet v0.10.0 protocol on the rpc v0.1
- Support for both V0 and V1 accounts and invoke transaction
- An account manager to install many account versions, including the
  [plugin account](https://github.com/argentlabs/starknet-plugin-account) and
  the [openzeppelin accounts](https://github.com/OpenZeppelin/cairo-contracts)
- An account interface and an implementation with the gateway and rpc v0.1
- Support for plugin extensions in the account and an implementation of the
  session key.
- Some tools to work with devnet and transform help with converting data
- An implementation of a fixed-size merkle tree algorithm

**known issues**: due to the fact rpc v0.2 is not yet fully supported by
Pathfinder and Devnet, the implementation is lacking from this release and
will be part of v0.5.0.
