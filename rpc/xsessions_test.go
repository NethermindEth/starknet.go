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

// TestXSessions_EstimateFee tests the xsessions EstimateFee
func TestXSessions_EstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Address          string
		PrivateKeyEnvVar string
		Call             types.FunctionCall
		PluginClassHash  string
		XSession         XSession
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c, false)
		testConfig.provider.c = spy
		account, err := testConfig.provider.NewAccount(
			os.Getenv(test.PrivateKeyEnvVar),
			test.Address,
			WithXSessionsPlugin(test.PluginClassHash, test.XSession))
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

// TestXSessions_Execute tests the xsessions Execute method
func TestXSessions_Execute(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Address          string
		PrivateKeyEnvVar string
		Call             types.FunctionCall
		PluginClassHash  string
		XSession         XSession
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c, false)
		testConfig.provider.c = spy
		account, err := testConfig.provider.NewAccount(
			os.Getenv(test.PrivateKeyEnvVar),
			test.Address,
			WithXSessionsPlugin(test.PluginClassHash, test.XSession),
		)
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
