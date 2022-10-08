## Accounts

This directory contains a number of **experimental** accounts and links related
to them. In particular, you will find:

- A link to the OZ v0.3.2 preset account you might not want to use
- A link to the OZ v0.4.0b preset account you might not want to use
- The implementation of the `YeaSayer` plugin that always agree on the
  transactions
- A simple OZ ERC-20 preset account that you probably can use
- A very basic counter account for demonstration purposes

## Building Accounts and Contracts

Since those contracts are being improved, we did not want to include them in
this repository but link them. To help you, we provide a `Makefile` that
shows how to build those contracts. You should:

- install cairo v0.9 with Python 3.9 and run the following commands

```shell
make
make ozv0
```

- install cairo v0.10 with Python 3.9 and run the following commands

```shell
make ozlatest
```

> Explaination:
> Starknet v0.9 and v.10 contract bytecodes should be compatible. However, the
> v0 transaction implementations of protocol tend to be developed with Cairo
> v0.9 when v1 transaction implementations of protocol tend to be developed
> with Cairo v0.10. That is why we need the 2 versions for now. 

## Testing Accounts and Contracts

> Note 2: RPC v0.1 is not supposed to work with v1 of transactions and v0.2 is
> not yet implemented neither by Pathfinder neither by Devnet. As a result, we 
> did not upgrade to the new implementations yet. We will do the change when
> we have implementations that are good enough on the 2 fronts!

