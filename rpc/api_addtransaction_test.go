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
	defer testConfig.client.Close()

	type testSetType struct {
		Salt                    string
		Contract                types.ContractClass
		ConstructorCallData     []string
		ExpectedTransactionHash string
		ExpectedContractAddress string
	}
	var contract types.ContractClass
	err := json.Unmarshal(counterFile, &contract)
	if err != nil {
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
