package starknetgo

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	rpc "github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

// TestAccountNonce tests the account Nonce
func TestRPCAccount_Nonce(t *testing.T) {
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
		// shim a keystore into existing tests.
		// use string representation of the PK as a fake sender address for the keystore
		ks := NewMemKeystore()
		pk := os.Getenv(test.PrivateKeyEnvVar)
		fakeSenderAddress := pk
		k := types.SNValToBN(pk)
		ks.Put(fakeSenderAddress, k)
		account, err := NewRPCAccount(utils.TestHexToFelt(t, fakeSenderAddress), utils.TestHexToFelt(t, test.Address), ks, testConfig.providerv02)
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
func TestRPCAccount_EstimateFee(t *testing.T) {
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
					ContractAddress:    utils.TestHexToFelt(t, "0x07704fb2d72fcdae1e6f658ef8521415070a01a3bd3cc5788f7b082126922b7b"),
					EntryPointSelector: types.GetSelectorFromNameFelt("increment"),
					Calldata:           []*felt.Felt{},
				},
			},
		},
		"testnet": {
			{
				Address:          TestNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: types.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x357b37bf12f59dd04c4da4933dcadf4a104e158365886d64ca0e554ada68fef"),
					EntryPointSelector: types.GetSelectorFromNameFelt("increment"),
					Calldata:           []*felt.Felt{},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		// shim a keystore into existing tests.
		// use string representation of the PK as a fake sender address for the keystore
		ks := NewMemKeystore()
		pk := os.Getenv(test.PrivateKeyEnvVar)
		fakeSenderAddress := pk
		k := types.SNValToBN(pk)
		ks.Put(fakeSenderAddress, k)
		account, err := NewRPCAccount(utils.TestHexToFelt(t, fakeSenderAddress), utils.TestHexToFelt(t, test.Address), ks, testConfig.providerv02)
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
		// shim a keystore into existing tests.
		// use string representation of the PK as a fake sender address for the keystore
		ks := NewMemKeystore()
		pk := os.Getenv(test.PrivateKeyEnvVar)
		fakeSenderAddress := pk
		k := types.SNValToBN(pk)
		ks.Put(fakeSenderAddress, k)
		account, err := NewRPCAccount(utils.TestHexToFelt(t, fakeSenderAddress), utils.TestHexToFelt(t, test.Address), ks, testConfig.providerv02)
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		execute, err := account.Execute(ctx, []types.FunctionCall{test.Call}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(execute.TransactionHash.String(), "0x") {
			t.Fatal("TransactionHash start with 0x, instead:", execute.TransactionHash)
		}
		fmt.Println("transaction_hash:", execute.TransactionHash)
		ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()
		status, err := account.rpc.WaitForTransaction(ctx, execute.TransactionHash, 8*time.Second)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if status != "PENDING" && status != "ACCEPTED_ON_L1" && status != "ACCEPTED_ON_L2" {
			t.Fatalf("tx %s wrong status: %s", execute.TransactionHash, status)
		}
	}
}
