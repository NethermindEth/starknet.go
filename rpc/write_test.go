package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/dontpanicdao/caigo/rpc/types"
)

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
