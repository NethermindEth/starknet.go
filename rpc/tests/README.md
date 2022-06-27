## About test manifests

This directory contains test manifests that are used by the tests:

- `counter.json` is a compile version of `rpc/contracts/counter.cairo`

To regenerate it, clone the repository and build the contract like below:

```shell
cd rpc
starknet-compile contracts/counter.cairo \
  --output tests/counter.json
```
