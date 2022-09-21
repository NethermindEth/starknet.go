package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
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
			ExpectedClassHash: "0x7944a315da387dcfa0c3ca204b81d836e37415e994834334e2a2d0c632344f0",
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
		dec, err := testConfig.apiv010.AddDeclareTransaction(context.Background(), contractClass, version)
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

// TestDeployTransaction tests starknet_addDeployTransaction
func TestDeployTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename                string
		Salt                    string
		ConstructorCall         []string
		ExpectedContractAddress string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:                "./tests/counter.json",
			Salt:                    "0xdeadbeef",
			ConstructorCall:         []string{"0x1"},
			ExpectedContractAddress: "0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167",
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

		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		dec, err := testConfig.apiv010.AddDeployTransaction(context.Background(), test.Salt, test.ConstructorCall, contractClass)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if dec.ContractAddress != test.ExpectedContractAddress {
			t.Fatalf("contractAddress does not match expected, current: %s", dec.ContractAddress)
		}
		if diff, err := spy.Compare(dec, false); err != nil || diff != "FullMatch" {
			spy.Compare(dec, true)
			t.Fatal("expecting to match", err)
		}
	}
}

// TestInvokeTransaction tests starknet_addDeployTransaction
func TesInvokeTransaction(t *testing.T) {
	// testConfig := beforeEach(t)

	type testSetType struct {
		AccountPrivateKeyEnvVar string
		AccountPublicKey        string
		AccountAddress          string
		ContractAddress         string
		ContractEntryPoint      string
		ContractCallData        []string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			AccountPrivateKeyEnvVar: "ACCOUNT_PRIVATE_KEY",
			AccountPublicKey:        "0x783318b2cc1067e5c06d374d2bb9a0382c39aabd009b165d7a268b882971d6",
			AccountAddress:          "0x19e63006d7df131737f5222283da28de2d9e2f0ee92fdc4c4c712d1659826b0",
			ContractAddress:         "0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167",
			ContractEntryPoint:      "incrementCounter",
			ContractCallData:        []string{"0x1"},
		}},
	}[testEnv]

	for _, test := range testSet {
		privateKey := os.Getenv(test.AccountPrivateKeyEnvVar)
		if privateKey == "" {
			t.Fatal("should have a private key for the account")
		}
		// 	caigo.GetSelectorFromName()
		// 	program := ""
		// 	if data, ok := v["program"]; ok {
		// 		dataProgram, err := json.Marshal(data)
		// 		if err != nil {
		// 			t.Fatal("should read file, instead:", err)
		// 		}
		// 		if program, err = encodeProgram(dataProgram); err != nil {
		// 			t.Fatal("should encode file, instead:", err)
		// 		}
		// 	}
		// 	entryPointsByType := types.EntryPointsByType{}
		// 	if data, ok := v["entry_points_by_type"]; ok {
		// 		dataEntryPointsByType, err := json.Marshal(data)
		// 		if err != nil {
		// 			t.Fatal("should marshall entryPointsByType, instead:", err)
		// 		}
		// 		err = json.Unmarshal(dataEntryPointsByType, &entryPointsByType)
		// 		if err != nil {
		// 			t.Fatal("should unmarshall entryPointsByType, instead:", err)
		// 		}
		// 	}
		// 	var abiPointer *types.ABI
		// 	if data, ok := v["abi"]; ok {
		// 		if abis, ok := data.([]interface{}); ok {
		// 			abiPointer, err = guessABI(abis)
		// 			if err != nil {
		// 				t.Fatal("should read ABI, instead:", err)
		// 			}
		// 		}
		// 	}

		// 	contractClass := types.ContractClass{
		// 		EntryPointsByType: entryPointsByType,
		// 		Program:           program,
		// 		Abi:               abiPointer,
		// 	}

		// 	spy := NewSpy(testConfig.client.c)
		// 	testConfig.client.c = spy
		// 	dec, err := testConfig.client.AddDeployTransaction(context.Background(), test.Salt, test.ConstructorCall, contractClass)
		// 	if err != nil {
		// 		t.Fatal("declare should succeed, instead:", err)
		// 	}
		// 	if dec.ContractAddress != test.ExpectedContractAddress {
		// 		t.Fatalf("contractAddress does not match expected, current: %s", dec.ContractAddress)
		// 	}
		// 	if diff, err := spy.Compare(dec, false); err != nil || diff != "FullMatch" {
		// 		spy.Compare(dec, true)
		// 		t.Fatal("expecting to match", err)
		// 	}
		// }
	}
}
