## About test manifests

This directory contains test manifests that are used by the tests:

- `account.json` is a compile version of the OpenZeppelin contract found in 
  [OpenZeppelin/cairo-contracts](https://github.com/OpenZeppelin/cairo-contracts).

To regenerate it, clone the repository and build the contract like below:

```shell
git clone https://github.com/OpenZeppelin/cairo-contracts.git
cd cairo-contracts/src
starknet-compile openzeppelin/account/Account.cairo \
   --output account.json --account_contract
```
