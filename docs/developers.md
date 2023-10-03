# A few words for developers

This document points some things you might want to know when onboarding on
starknet.go. For now, we have listed a number of questions that we would like you
know before you contribute:

<!-- - What version of `rpc` should I use -->
- What is the difference between `rpc` and `gateway`?
- Where are the resources about the protocol?
- How to interact with accounts?
- How to better understand the protocol?

## What is the difference between `rpc` and `gateway`?

starknet.go provides an access to Starknet through several interfaces:
- the `gateway` remains a primary access to Starknet. It is managed by
  Starkware and is stable. It evolves with the protocol features and is
  actually split in 2: the gateway allows to run transactions, and the feeder
  allows to read the state of the blockchain. It can be accessed directly from
  the Internet.
- the `rpc` access is the next access to Starknet. It requires you have run a 
  node like [eqlabs/pathfinder](https://github.com/eqlabs/pathfinder). This
  access has several benefits, including the fact that it will implement a
  peer-to-peer connection and is based on `openrpc` 2.0. As a result, the
  protocol definition is better standardize and starknet.go uses the
  [go-ethereum/rpc](https://pkg.go.dev/github.com/ethereum/go-ethereum/rpc)
  client to access Starknet.

If you wonder which one to use, the short answer is it depends. However, starknet.go
provide an access to both and tries to share a common interface between the 2.
Be careful that when you are doing something on the project, you should pay
attention to the 2 interfaces, even if the implementations are specific.

<!-- ## What version of `rpc` should I use?

`rpc` is currently being upgrade from `v0.1` to `v0.2`. Right now `starknet.go` only
supports `v0.1`. However to provide a smooth upgrade the plan is to support
both in parallel. That is why the current package names are: -->

## Where do I find resources about the protocol specifications?

There are several efforts to document how the API is working. You should
know them and you can contribute to the associated projects:

- [starknet-edu postman reference](https://www.postman.com/starknet-edu/workspace/starknet-edu/collection/20082312-a5291c43-a4e5-4a6d-9c51-125c6acd3b41?ctx=documentation)
  is a project that lists all the existing API and provides some examples
  of how to use them.
- There is no true specification of the API we are aware of, however the
  implementation is opensource and can be read from the
  starkware-libs/cairo-lang
  [api/gateway directory](https://github.com/starkware-libs/cairo-lang/tree/master/src/starkware/starknet/services/api/gateway)
- Starknet `openrpc` specification is available in the
  `starkware-libs/starknet-specs` project; check the content of
  [this directory](https://github.com/starkware-libs/starknet-specs/blob/master/api)

If you want to go deeper into the protocol, there are other interesting
resources available, check:

- [a set of examples](https://github.com/eqlabs/pathfinder/blob/main/crates/pathfinder/rpc_examples.sh)
  from the eqlabs/pathfinder project that provides RPC tests with the project.
  Note that the [go-ethereum] implementation of openrpc only supports positional
  arguments and, some of these examples must be changed into positional calls.
- `#🦦| starknet.go` channel in the discord from Starknet.

## How to interact with accounts?

An important part of interacting with Starknet consists in interacting
accounts. The reason is that you need to go through an account to run transactions,
including to create a contract from a class. That is why you should understand
what an account is and how to interact with it.

The devil living in the details, there are specificities associated with
account implementations. For instance, you might find that one account does not
support the same signature as an other. To start with account, read:

- [Starknet Account Abstraction Model - Part 1](https://community.starknet.io/t/starknet-account-abstraction-model-part-1/781)
- [Starknet Account Abstraction Model - Part 2](https://community.starknet.io/t/starknet-account-abstraction-model-part-2/839)
- [Learn how to build and deploy Starknet Accounts](https://github.com/starknet-edu/starknet-accounts)
  and the companion [Starknet workshop](https://www.youtube.com/watch?v=51Qb3TLpNro)
- Openzeppelin Cairo contracts that include an account
  [implementation](https://github.com/OpenZeppelin/cairo-contracts/tree/main/src/openzeppelin/account)
- The argent-x contract
  [implementation](https://github.com/argentlabs/argent-contracts-starknet)

And also, read the [Building on Starknet](https://starknet.io/building-on-starknet/)
section of the documentation.

## How to better understand the protocol?

Sometimes reading the documentation is not enough. Check the examples, check
integration tests and use tools to capture paylods:

- [MITM reverse](https://docs.mitmproxy.org/stable/concepts-modes/#reverse-proxy)
  mode allows to create a reverse proxy to the API and capture the workload. You
  then should be able to run transaction with tools like the Starknet CLI and
  better understand that is happening under the hood.
- If you interact from Dapps, like Voyager and want to understand how the
  interactions with contract are happening, you can use:
  - Chrome [debugging extensions](https://developer.chrome.com/docs/extensions/mv3/tut_debugging/)
  - Firefox [debugging extensions](https://extensionworkshop.com/documentation/develop/debugging/)
- Another most valuable resource is reading the [starknet.js](https://github.com/0xs34n/starknet.js)
  project.
