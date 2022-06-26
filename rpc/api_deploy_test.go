package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

//go:embed tests/account.json
var accountFile []byte

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
	err := json.Unmarshal(accountFile, &contract)
	if err != nil {
		t.Fatal("error loading contract:", err)
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				Salt:                    "0xffffff",
				Contract:                contract,
				ConstructorCallData:     []string{"5ff5eff3bed10c5109c25ab3618323d74a436e7e0b66a512ca6dbab27f08a6"},
				ExpectedTransactionHash: "0xdeadbeef",
				ExpectedContractAddress: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				Salt:                    "0xffffff",
				Contract:                contract,
				ConstructorCallData:     []string{"5ff5eff3bed10c5109c25ab3618323d74a436e7e0b66a512ca6dbab27f08a6"},
				ExpectedTransactionHash: "0x60206d9eaa2466235831009a8cd12c8ad019566aa4dccb5fd82c2f5c706cbf1",
				ExpectedContractAddress: "0x2f2ce1da30dbd3727f5e68b4f2ca7ee0e0108a0c9350a08a7b4d46564e9a21c",
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
		fmt.Printf("%+v\n", output)
		if output.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("tx expected %s, got %s", test.ExpectedTransactionHash, output.TransactionHash)
		}
		if output.ContractAddress != test.ExpectedContractAddress {
			t.Fatalf("contract expected %s, got %s", test.ExpectedContractAddress, output.ContractAddress)
		}
	}
}
