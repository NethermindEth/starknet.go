package caigo

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	rpc "github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

// TestAccountNonce tests the account Nonce
func TestAccountNonce(t *testing.T) {
	testConfig := beforeRPCEach(t)

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
		account, err := NewRPCAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address, testConfig.provider)
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
	testConfig := beforeRPCEach(t)

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
					Calldata:           []string{},
				},
			},
		},
		"testnet": {
			{
				Address:          TestNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x357b37bf12f59dd04c4da4933dcadf4a104e158365886d64ca0e554ada68fef"),
					EntryPointSelector: "increment",
					Calldata:           []string{},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := NewRPCAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address, testConfig.provider)
		if err != nil {
			t.Fatal(err)
		}
		estimate, err := account.EstimateFee(context.Background(), []types.FunctionCall{test.Call}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		if types.HexToBN(string(estimate.OverallFee)).Cmp(big.NewInt(1000000)) < 0 {
			t.Fatal("OverallFee should be > 1000000, instead:", estimate.OverallFee)
		}
	}
}

// TestRPCAccount_Execute tests the account Execute method
func TestRPCAccount_Execute(t *testing.T) {
	testConfig := beforeRPCEach(t)

	type testSetType struct {
		Address          string
		PrivateKeyEnvVar string
		Call             types.FunctionCall
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		account, err := NewRPCAccount(os.Getenv(test.PrivateKeyEnvVar), test.Address, testConfig.provider)
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		execute, err := account.Execute(ctx, []types.FunctionCall{test.Call}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
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
