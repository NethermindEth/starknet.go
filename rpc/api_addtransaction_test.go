package rpc

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

//go:embed tests/counter.json
var counterFile []byte

// TestAddDeployTransaction tests AddDeployTransaction
func TestAddDeployTransaction(t *testing.T) {
	_ = beforeEach(t)

	type testSetType struct {
		BroadcastedDeployTxn    BroadcastedDeployTxn
		ExpectedTransactionHash string
		ExpectedContractAddress string
	}
	var contract types.ContractClass

	if err := json.Unmarshal(counterFile, &contract); err != nil {
		t.Fatal("error loading contract:", err)
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BroadcastedDeployTxn: BroadcastedDeployTxn{
					ContractClass: ContractClass{
						Program: "",
						EntryPointsByType: struct {
							CONSTRUCTOR ContractEntryPointList "json:\"CONSTRUCTOR\""
							EXTERNAL    ContractEntryPointList "json:\"EXTERNAL\""
							L1_HANDLER  ContractEntryPointList "json:\"L1_HANDLER\""
						}{},
						Abi: &ContractABI{},
					},
				},
				ExpectedTransactionHash: "0xdeadbeef",
				ExpectedContractAddress: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				BroadcastedDeployTxn: BroadcastedDeployTxn{
					ContractClass: ContractClass{
						Program: "",
						EntryPointsByType: struct {
							CONSTRUCTOR ContractEntryPointList "json:\"CONSTRUCTOR\""
							EXTERNAL    ContractEntryPointList "json:\"EXTERNAL\""
							L1_HANDLER  ContractEntryPointList "json:\"L1_HANDLER\""
						}{},
						Abi: &ContractABI{},
					},
				},
				ExpectedTransactionHash: "0x2149bf99d96ed687a488091ea0d2b1e0b24f73fd7ab96809c2640ae2fc0c791",
				ExpectedContractAddress: "0x30b0fc513edb49b5602f985f5515540a63c8884f3c23a9a9b70f3c14eab7255",
			},
		},
		"devnet":  {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		_ = test
	}
}

// TestAddDeclareTransaction tests AddDeclareTransaction
func TestAddDeclareTransaction(t *testing.T) {
	_ = beforeEach(t)

	type testSetType struct {
		BroadcastedDeclareTxn   BroadcastedDeclareTxn
		ExpectedTransactionHash string
		ExpectedClassHash       string
	}
	var contract types.ContractClass

	if err := json.Unmarshal(counterFile, &contract); err != nil {
		t.Fatal("error loading contract:", err)
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		// TODO: add tests for mainnet when possible or when figure out how to
		// create a white-listed contract.
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		_ = test
	}
}

// TestAddInvokeTransaction tests AddInvokeTransaction
func TestAddInvokeTransaction(t *testing.T) {
	_ = beforeEach(t)

	type testSetType struct {
		BroadcastedInvokeTxn    BroadcastedInvokeTxn
		ExpectedTransactionHash string
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		_ = test
	}
}
