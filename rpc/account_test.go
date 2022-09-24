package rpc

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestAccountNonce tests the account Nonce
func TestAccountNonce(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Provider         *Provider
		Address          string
		PrivateKeyEnvVar string
	}

	testSet := map[string][]testSetType{
		"devnet": {
			{
				Address:          DevNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			},
		},
		"mock": {},
		"testnet": {
			{
				Address:          TestNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := testConfig.provider.NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address)
		if err != nil {
			t.Fatal(err)
		}
		_, err = account.Nonce(context.Background())
		if err != nil {
			t.Fatal(err)
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
		"devnet": {
			{
				Address:          DevNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x035a55a64238b776664d7723de1f6b50350116a1ab1ca1fe154320a0eba53d3a"),
					EntryPointSelector: "increment",
					CallData:           []string{},
				},
			},
		},
		"mock": {},
		"testnet": {
			{
				Address:          TestNetAccount032Address,
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
		spy := NewSpy(testConfig.provider.c, false)
		testConfig.provider.c = spy
		account, err := testConfig.provider.NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address)
		if err != nil {
			t.Fatal(err)
		}
		estimate, err := account.EstimateFee(context.Background(), []types.FunctionCall{test.Call}, types.ExecuteDetails{})
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
		"devnet": {
			// TODO: there is a problem with devnet 0.3.1 that does not implement
			// positional argument properly for starknet_addInvokeTransaction. I
			// have proposed a PR https://github.com/Shard-Labs/starknet-devnet/pull/283
			// to address the issue. Meanwhile, we are stuck with that feature on devnet.
			// {
			// 	Address:          DevNetAccountAddress,
			// 	PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			// 	Call: types.FunctionCall{
			// 		ContractAddress:    types.HexToHash("0x035a55a64238b776664d7723de1f6b50350116a1ab1ca1fe154320a0eba53d3a"),
			// 		EntryPointSelector: "increment",
			// 		CallData:           []string{},
			// 	},
			// },
		},
		"mock": {},
		"testnet": {
			{
				Address:          TestNetAccount032Address,
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
		spy := NewSpy(testConfig.provider.c, false)
		testConfig.provider.c = spy
		account, err := testConfig.provider.NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address)
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		execute, err := account.Execute(ctx, []types.FunctionCall{test.Call}, types.ExecuteDetails{})
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
		fmt.Println("transaction_hash:", execute.TransactionHash)
		ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()
		status, err := account.Provider.WaitForTransaction(ctx, types.HexToHash(execute.TransactionHash), 8*time.Second)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if status != "PENDING" && status != "ACCEPTED_ON_L1" && status != "ACCEPTED_ON_L2" {
			t.Fatalf("tx %s wrong status: %s", execute.TransactionHash, status)
		}
	}
}
