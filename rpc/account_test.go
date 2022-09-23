package rpc

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestAccountNonce tests the account Nonce
func TestAccountNonce(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Provider         *Client
		Address          string
		PrivateKeyEnvVar string
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				Address:          "0x19e63006d7df131737f5222283da28de2d9e2f0ee92fdc4c4c712d1659826b0",
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := testConfig.client.NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address)
		if err != nil {
			t.Fatal(err)
		}
		nonce, err := account.Nonce(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if nonce.Uint64() <= 1 {
			t.Fatal("nonce should be > 1", nonce.Uint64())
		}
	}
}

// TestAccountEstimateFee tests the account EstimateFee
func TestAccountEstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Address          string
		PrivateKeyEnvVar string
		Call             types.FunctionCall
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				Address:          "0x19e63006d7df131737f5222283da28de2d9e2f0ee92fdc4c4c712d1659826b0",
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167"),
					EntryPointSelector: "increment",
					CallData:           []string{},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c, false)
		testConfig.client.c = spy
		account, err := testConfig.client.NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address)
		if err != nil {
			t.Fatal(err)
		}
		estimate, err := account.EstimateFee(context.Background(), []types.FunctionCall{test.Call}, ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(estimate, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(estimate, true)
			t.Fatal("expecting to match, instead:", diff)
		}
		if caigo.HexToBN(string(estimate.OverallFee)).Cmp(big.NewInt(1000000)) < 0 {
			t.Fatal("OverallFee should be > 1000000, instead:", estimate.OverallFee)
		}
	}
}

// TestAccountExecute tests the account Execute method
func TestAccountExecute(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Address          string
		PrivateKeyEnvVar string
		Call             types.FunctionCall
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				Address:          "0x19e63006d7df131737f5222283da28de2d9e2f0ee92fdc4c4c712d1659826b0",
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167"),
					EntryPointSelector: "increment",
					CallData:           []string{},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c, false)
		testConfig.client.c = spy
		account, err := testConfig.client.NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address)
		if err != nil {
			t.Fatal(err)
		}
		execute, err := account.Execute(context.Background(), []types.FunctionCall{test.Call}, ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(execute, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(execute, true)
			t.Fatal("expecting to match, instead:", diff)
		}
		if !strings.HasPrefix(execute.TransactionHash, "0x") {
			t.Fatal("TransactionHash start with 0x, instead:", execute.TransactionHash)
		}
	}
}
