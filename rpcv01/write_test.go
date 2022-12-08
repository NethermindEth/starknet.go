package rpcv01

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename          string
		Version           string
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			Filename:          "./tests/counter.json",
			Version:           "0x0",
			ExpectedClassHash: "0x01649a376a9aa5ccb5ddf2f59c267de5fb6b3b177056a53f45d42877c856a051",
		}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:          "./tests/counter.json",
			Version:           "0x0",
			ExpectedClassHash: "0x4484265a6e003e8afe272e6c9bf3e7d0d8e343b2df57763a995828285fdfbbd",
		}},
	}[testEnv]

	for _, test := range testSet {
		content, err := os.ReadFile(test.Filename)
		if err != nil {
			t.Fatal("should read file with success, instead:", err)
		}
		contractClass := ctypes.ContractClass{}
		if err := json.Unmarshal(content, &contractClass); err != nil {
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
		Filename                string
		Salt                    string
		ConstructorCall         []string
		ExpectedContractAddress string
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Filename:                "./tests/counter.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{"0x1"},
				ExpectedContractAddress: "0x035a55a64238b776664d7723de1f6b50350116a1ab1ca1fe154320a0eba53d3a",
			},
			{
				Filename:                "./tests/oz_v0.3.2_account.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: DevNetAccount032Address,
			},
			{
				Filename:                "./tests/oz_v0.4.0b_account.json",
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
		content, err := os.ReadFile(test.Filename)
		if err != nil {
			t.Fatal("should read file with success, instead:", err)
		}
		contractClass := ctypes.ContractClass{}
		if err := json.Unmarshal(content, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		dec, err := testConfig.provider.AddDeployTransaction(context.Background(), test.Salt, test.ConstructorCall, contractClass)
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
		fmt.Println("transaction hash:", dec.TransactionHash)
	}
}
