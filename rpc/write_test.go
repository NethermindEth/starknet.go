package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestAddInvokeTransaction tests AddInvokeTransaction
func TestAddInvokeTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		InvokeTxn               types.InvokeV0
		Signature               []string
		MaxFee                  string
		Version                 string
		ExpectedTransactionHash string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			InvokeTxn: types.InvokeV0{
				ContractAddress:    types.HexToHash("0x23371b227eaecd8e8920cd429d2cd0f3fee6abaacca08d3ab82a7cdd"),
				EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
				CallData: []string{
					"0x1",
					"0x677bb1cdc050e8d63855e8743ab6e09179138def390676cc03c484daf112ba1",
					"0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320",
					"0x0",
					"0x1",
					"0x1",
					"0x2b",
					"0x0",
				},
			},
			MaxFee:                  "0x4f388496839",
			Signature:               []string{"0x7dd3a55d94a0de6f3d6c104d7e6c88ec719a82f4e2bbc12587c8c187584d3d5", "0x71456dded17015d1234779889d78f3e7c763ddcfd2662b19e7843c7542614f8"},
			Version:                 "0x0",
			ExpectedTransactionHash: "0x72242e7e366576ff33fbe8f772ab8a58f34273126a145e40e27f432471d1471",
		}},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy

		txHash, err := testConfig.client.AddInvokeTransaction(context.Background(), test.InvokeTxn, test.Signature, test.MaxFee, test.Version)
		if err != nil {
			t.Fatalf("invoke should succeed, instead: %v", err)
		}
		if txHash.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("transaction hash does not match expected, instead: %s", txHash.TransactionHash)
		}

		if diff, err := spy.Compare(txHash, false); err != nil || diff != "FullMatch" {
			spy.Compare(txHash, true)
			t.Fatalf("expected to match, instead: %v", err)
		}
	}
}

// TestChainID checks the chainId matches the one for the environment
func TestDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename          string
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:          "./tests/counter.json",
			ExpectedClassHash: "0x71ee80b6ec623ffe1adb80a0f0cd1a3012b633d7312fe604816ce7aa75cf209",
		}},
	}[testEnv]

	for _, test := range testSet {
		content, err := os.ReadFile(test.Filename)
		if err != nil {
			t.Fatal("should read file with success, instead:", err)
		}
		v := map[string]interface{}{}
		if err := json.Unmarshal(content, &v); err != nil {
			t.Fatal("should parse file with success, instead:", err)
		}

		program := ""
		if data, ok := v["program"]; ok {
			dataProgram, err := json.Marshal(data)
			if err != nil {
				t.Fatal("should read file, instead:", err)
			}
			if program, err = encodeProgram(dataProgram); err != nil {
				t.Fatal("should encode file, instead:", err)
			}
		}
		entryPointsByType := types.EntryPointsByType{}
		if data, ok := v["entry_points_by_type"]; ok {
			dataEntryPointsByType, err := json.Marshal(data)
			if err != nil {
				t.Fatal("should marshall entryPointsByType, instead:", err)
			}
			err = json.Unmarshal(dataEntryPointsByType, &entryPointsByType)
			if err != nil {
				t.Fatal("should unmarshall entryPointsByType, instead:", err)
			}
		}
		var abiPointer *types.ABI
		if data, ok := v["abi"]; ok {
			if abis, ok := data.([]interface{}); ok {
				abiPointer, err = guessABI(abis)
				if err != nil {
					t.Fatal("should read ABI, instead:", err)
				}
			}
		}

		contractClass := types.ContractClass{
			EntryPointsByType: entryPointsByType,
			Program:           program,
			Abi:               abiPointer,
		}
		version := "0x0"

		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		dec, err := testConfig.client.AddDeclareTransaction(context.Background(), contractClass, version)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if dec.ClassHash != test.ExpectedClassHash {
			t.Fatalf("classHash does not match expected, current: %s", dec.ClassHash)
		}
		if diff, err := spy.Compare(dec, false); err != nil || diff != "FullMatch" {
			spy.Compare(dec, true)
			t.Fatal("expecting to match", err)
		}
	}
}
