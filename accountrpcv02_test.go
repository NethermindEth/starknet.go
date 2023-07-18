package caigo

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	rpc "github.com/NethermindEth/caigo/rpcv02"
	"github.com/NethermindEth/caigo/types"
	"github.com/NethermindEth/caigo/utils"
	"github.com/NethermindEth/juno/core/felt"
)

// TestAccountNonce tests the account Nonce
func TestRPCv02Account_Nonce(t *testing.T) {
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
		// Shim a keystore into existing tests.
		// Use string representation of the PK as a fake sender address for the keystore.
		ks := NewMemKeystore()
		pk := os.Getenv(test.PrivateKeyEnvVar)
		fakeSenderAddress := pk
		k := types.SNValToBN(pk)
		ks.Put(fakeSenderAddress, k)

		// Convert the fake sender address and the test address to Felt.
		fakeSenderFelt, err := utils.HexToFelt(fakeSenderAddress)
		if err != nil {
			t.Fatalf("Failed to convert fake sender address to Felt: %v", err)
		}
		testAddressFelt, err := utils.HexToFelt(test.Address)
		if err != nil {
			t.Fatalf("Failed to convert test address to Felt: %v", err)
		}

		account, err := NewRPCAccount(fakeSenderFelt, testAddressFelt, ks, testConfig.providerv02, AccountVersion1)
		fmt.Println(account.rpcv02, testConfig.providerv02)
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
func TestRPCv02Account_EstimateFee(t *testing.T) {
	testConfig := beforeRPCEach(t)

	type testSetType struct {
		Address          string
		PrivateKeyEnvVar string
		Call             types.FunctionCall
	}

	devnetContractAddress, err := utils.HexToFelt("0x07704fb2d72fcdae1e6f658ef8521415070a01a3bd3cc5788f7b082126922b7b")
	if err != nil {
		t.Fatalf("Failed to convert devnet contract address to Felt: %v", err)
	}

	testnetContractAddress, err := utils.HexToFelt("0x07704fb2d72fcdae1e6f658ef8521415070a01a3bd3cc5788f7b082126922b7b")
	if err != nil {
		t.Fatalf("Failed to convert devnet contract address to Felt: %v", err)
	}

	testSet := map[string][]testSetType{
		"devnet": {
			{
				Address:          DevNetAccount032Address,
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
				Call: types.FunctionCall{
					ContractAddress:    devnetContractAddress,
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
					ContractAddress:    testnetContractAddress,
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
		fakeSenderFelt, err := utils.HexToFelt(fakeSenderAddress)
		if err != nil {
			t.Fatalf("Failed to convert fake sender address to Felt: %v", err)
		}
		testAddressFelt, err := utils.HexToFelt(test.Address)
		if err != nil {
			t.Fatalf("Failed to convert test address to Felt: %v", err)
		}

		account, err := NewRPCAccount(fakeSenderFelt, testAddressFelt, ks, testConfig.providerv02, AccountVersion1)
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
func TestRPCv02Account_Execute(t *testing.T) {
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
		fakeSenderFelt, err := utils.HexToFelt(fakeSenderAddress)
		if err != nil {
			t.Fatalf("Failed to convert fake sender address to Felt: %v", err)
		}
		testAddressFelt, err := utils.HexToFelt(test.Address)
		if err != nil {
			t.Fatalf("Failed to convert test address to Felt: %v", err)
		}

		account, err := NewRPCAccount(fakeSenderFelt, testAddressFelt, ks, testConfig.providerv02)
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		execute, err := account.Execute(ctx, []types.FunctionCall{test.Call}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("transaction_hash:", execute.TransactionHash)
		ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()
		status, err := account.rpcv02.WaitForTransaction(ctx, execute.TransactionHash, 8*time.Second)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if status != "PENDING" && status != "ACCEPTED_ON_L1" && status != "ACCEPTED_ON_L2" {
			t.Fatalf("tx %s wrong status: %s", execute.TransactionHash, status)
		}
	}
}
