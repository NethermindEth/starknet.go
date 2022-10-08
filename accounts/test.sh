#!/bin/bash

set -ex

go run . -command install -with-plugin -with-proxy -account-version v0
go run . -command execute -with-plugin
go run . -command execute
rm -f .starknet-account.json
go run . -command install -account-version v0
go run . -command execute
rm -f .starknet-account.json

