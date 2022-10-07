package account

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo/rpc"
	"github.com/dontpanicdao/caigo/rpc/types"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// TestAccountNonce tests the account Nonce
func TestAccountNonce(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Provider         *rpc.Provider
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
		"testnet": {
			{
				Address:          TestNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address, testConfig.provider)
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
		Call             ctypes.FunctionCall
	}

	testSet := map[string][]testSetType{
		"devnet": {
			{
				Address:          DevNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: ctypes.FunctionCall{
					ContractAddress:    ctypes.HexToHash("0x035a55a64238b776664d7723de1f6b50350116a1ab1ca1fe154320a0eba53d3a"),
					EntryPointSelector: "increment",
					Calldata:           []string{},
				},
			},
		},
		"testnet": {
			{
				Address:          TestNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: ctypes.FunctionCall{
					ContractAddress:    ctypes.HexToHash("0x357b37bf12f59dd04c4da4933dcadf4a104e158365886d64ca0e554ada68fef"),
					EntryPointSelector: "increment",
					Calldata:           []string{},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address, testConfig.provider)
		if err != nil {
			t.Fatal(err)
		}
		estimate, err := account.EstimateFee(context.Background(), []ctypes.FunctionCall{test.Call}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		if ctypes.HexToBN(string(estimate.OverallFee)).Cmp(big.NewInt(1000000)) < 0 {
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
		Call             ctypes.FunctionCall
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
			// 		Calldata:           []string{},
			// 	},
			// },
		},
		"testnet": {
			// Disabled tests due to the fact it is taking ages on the CI. It should
			// work on demand though...
			// {
			// 	Address:          TestNetAccount032Address,
			// 	PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			// 	Call: types.FunctionCall{
			// 		ContractAddress:    types.HexToHash("0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167"),
			// 		EntryPointSelector: "increment",
			// 		Calldata:           []string{},
			// 	},
			// },
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := NewAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address, testConfig.provider)
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		execute, err := account.Execute(ctx, []ctypes.FunctionCall{test.Call}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(execute.TransactionHash, "0x") {
			t.Fatal("TransactionHash start with 0x, instead:", execute.TransactionHash)
		}
		fmt.Println("transaction_hash:", execute.TransactionHash)
		ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()
		status, err := account.Provider.WaitForTransaction(ctx, ctypes.HexToHash(execute.TransactionHash), 8*time.Second)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if status != "PENDING" && status != "ACCEPTED_ON_L1" && status != "ACCEPTED_ON_L2" {
			t.Fatalf("tx %s wrong status: %s", execute.TransactionHash, status)
		}
	}
}
