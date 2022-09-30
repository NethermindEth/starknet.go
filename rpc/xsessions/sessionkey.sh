#!/bin/bash

set -e 

go clean -testcache
go test -v -run TestSessionKey_RegisterPlugin
go test -v -run TestSessionKey_DeployAccount
go test -v -run TestSessionKey_MintEth
go test -v -run TestSessionKey_CheckEth
go test -v -run TestCounter_DeployContract
go test -v -run TestCounter_IncrementWithSessionKeyPlugin
