## Accounts

This directory contains a number of **experimental** accounts and links related
to them. In particular, you will find:

- A link to the OZ v0.3.2 preset account you might not want to use
- A link to the OZ v0.4.0b preset account you might not want to use
- A link to the Argentlabs/Ledger/Catridge plugin account, you might not want
  to use either
- The implementation of the `YeaSayer` plugin that always agree on the
  transactions and the implementation of the SessionKey plugin that comes with
  the plugin account
- A very basic counter account for demonstration purposes
- A simple OZ ERC-20 preset account that you probably can use

## Building Accounts and Contracts

Since those contracts are being improved, we did not want to include them in
this repository but link them. To help you, we provide a `Makefile` that
shows how to build those contracts. You should:

- install cairo v0.9.1 with Python 3.9 and run the following commands

```shell
# this command reset the OZ submodule to a version that was providing the v0
# account. It is using the cairo v0.9 syntax
make v0 -f MakefileV0
# compile the contracts
make -f MakefileV0
```

- install cairo v0.10 with Python 3.9 and run the following commands

```shell
make latest
make
```

> Explanation:
> Starknet v0.9 and v.10 contract bytecodes should be compatible. However, the
> v0 transaction implementations of protocol tend to be developed with Cairo
> v0.9 when v1 transaction implementations of protocol tend to be developed
> with Cairo v0.10. That is why we need the 2 versions for now. 
