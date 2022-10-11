#!/bin/bash

set -ex

go run . -command install -with-plugin -with-proxy -account-version v0  -provider rpcv01
go run . -command execute -with-plugin  -provider rpcv01
go run . -command execute  -provider rpcv01
rm -f .starknet-account.json
go run . -command install -account-version v0 -provider rpcv01
go run . -command execute -provider rpcv01
rm -f .starknet-account.json
# TODO: Monitor https://github.com/Shard-Labs/starknet-devnet/pull/303 progress
# go run . -command install -account-version v1 -provider rpcv01
# go run . -command execute -provider rpcv01
# rm -f .starknet-account.json
go run . -command install -account-version v0 -provider gateway
go run . -command execute -provider gateway
rm -f .starknet-account.json
go run . -command install -with-plugin -with-proxy -account-version v0 -provider gateway
go run . -command execute -provider gateway
rm -f .starknet-account.json
go run . -command install -account-version v1 -provider gateway
go run . -command execute -provider gateway
rm -f .starknet-account.json
