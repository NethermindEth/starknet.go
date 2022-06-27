package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

//go:embed tests/counter.json
var counterFile []byte

// TestAddDeployTransaction tests AddDeployTransaction
func TestAddDeployTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Salt                    string
		Contract                types.ContractClass
		ConstructorCallData     []string
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
				Salt:                    "0xffffff",
				Contract:                contract,
				ConstructorCallData:     []string{},
				ExpectedTransactionHash: "0xdeadbeef",
				ExpectedContractAddress: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				Salt:                    "0xffffff",
				Contract:                contract,
				ConstructorCallData:     []string{},
				ExpectedTransactionHash: "0x2149bf99d96ed687a488091ea0d2b1e0b24f73fd7ab96809c2640ae2fc0c791",
				ExpectedContractAddress: "0x30b0fc513edb49b5602f985f5515540a63c8884f3c23a9a9b70f3c14eab7255",
			},
		},
		// TODO: add tests for mainnet when possible or when figure out how to
		// create a white-listed contract. For now, the output is:
		//   - code: NotPermittedContract
		//   - message: The contract class attempted to be deployed is not permitted.
		// This behavior isis on purpose due to the fact mainnet is under limited
		// access. For more details, check this discord
		// [conversation](https://discord.com/channels/793094838509764618/793094838987128844/990692360608444508)
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		output, err := testConfig.client.AddDeployTransaction(context.Background(), test.Salt, test.ConstructorCallData, test.Contract)
		if err != nil {
			t.Fatal(err)
		}
		if output.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("tx expected %s, got %s", test.ExpectedTransactionHash, output.TransactionHash)
		}
		if output.ContractAddress != test.ExpectedContractAddress {
			t.Fatalf("contract expected %s, got %s", test.ExpectedContractAddress, output.ContractAddress)
		}
	}
}

// TestAddDeclareTransaction tests AddDeclareTransaction
func TestAddDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Contract                types.ContractClass
		Version                 string
		ExpectedTransactionHash string
		ExpectedClassHash       string
	}
	var contract types.ContractClass

	if err := json.Unmarshal(counterFile, &contract); err != nil {
		t.Fatal("error loading contract:", err)
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				Contract:                contract,
				Version:                 "0x0",
				ExpectedTransactionHash: "0xdeadbeef",
				ExpectedClassHash:       "0xdeadbeef",
			},
		},
		"testnet": {
			{
				Contract:                contract,
				Version:                 "0x0",
				ExpectedTransactionHash: "0x3d570dbde5ed56ddcb5f69578fb5f83b362c4af8b2a60e2be33ed229148e10a",
				ExpectedClassHash:       "0x646552d8029a8fe940dbbe2847bce558d3d1b3e78a5519e970395df6a2b2cc9",
			},
		},
		// TODO: add tests for mainnet when possible or when figure out how to
		// create a white-listed contract.
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		output, err := testConfig.client.AddDeclareTransaction(context.Background(), test.Contract, test.Version)
		if err != nil {
			t.Fatal(err)
		}
		if output.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("tx expected %s, got %s", test.ExpectedTransactionHash, output.TransactionHash)
		}
		if output.ClassHash != test.ExpectedClassHash {
			t.Fatalf("class expected %s, got %s", test.ExpectedClassHash, output.ClassHash)
		}
	}
}
