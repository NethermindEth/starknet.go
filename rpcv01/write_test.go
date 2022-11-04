package rpcv01

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dontpanicdao/caigo/artifacts"
	"github.com/dontpanicdao/caigo/types"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename          []byte
		Version           string
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			Filename:          artifacts.CounterCompiled,
			Version:           "0x0",
			ExpectedClassHash: "0x07c9716ab7c3e398ff5bf5b90366e18687915a57080c0e9f9c58ffea9bdf02e2",
		}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:          artifacts.CounterCompiled,
			Version:           "0x0",
			ExpectedClassHash: "0x4484265a6e003e8afe272e6c9bf3e7d0d8e343b2df57763a995828285fdfbbd",
		}},
	}[testEnv]

	for _, test := range testSet {
		contractClass := types.ContractClass{}
		if err := json.Unmarshal(test.Filename, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		dec, err := testConfig.provider.AddDeclareTransaction(context.Background(), contractClass, test.Version)
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
		fmt.Println("transaction hash:", dec.TransactionHash)
	}
}

// TestDeployTransaction tests starknet_addDeployTransaction
func TestDeployTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename                []byte
		Salt                    string
		ConstructorCall         []string
		ExpectedContractAddress string
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Filename:                artifacts.CounterCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{},
				ExpectedContractAddress: "0x07704fb2d72fcdae1e6f658ef8521415070a01a3bd3cc5788f7b082126922b7b",
			},
			{
				Filename:                artifacts.AccountV0Compiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: DevNetAccount032Address,
			},
			{
				Filename:                artifacts.AccountCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: DevNetAccount040Address,
			},
		},
		"mainnet": {},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		contractClass := types.ContractClass{}
		if err := json.Unmarshal(test.Filename, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		dec, err := testConfig.provider.AddDeployTransaction(context.Background(), test.Salt, test.ConstructorCall, contractClass)
		if err != nil {
			t.Fatal("deploy should succeed, instead:", err)
		}
		if dec.ContractAddress != test.ExpectedContractAddress {
			t.Fatalf("contractAddress does not match expected %s, got: %s", test.ExpectedContractAddress, dec.ContractAddress)
		}
		if diff, err := spy.Compare(dec, false); err != nil || diff != "FullMatch" {
			spy.Compare(dec, true)
			t.Fatal("expecting to match", err)
		}
		fmt.Println("transaction hash:", dec.TransactionHash)
	}
}
