package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename          string
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			Filename:          "./tests/counter.json",
			ExpectedClassHash: "0x01649a376a9aa5ccb5ddf2f59c267de5fb6b3b177056a53f45d42877c856a051",
		}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:          "./tests/counter.json",
			ExpectedClassHash: "0x7cca67b54cd7edfcdd45ceef4e43636b926101a26a99af003722f7ef10b08b3",
		}},
	}[testEnv]

	for _, test := range testSet {
		content, err := os.ReadFile(test.Filename)
		if err != nil {
			t.Fatal("should read file with success, instead:", err)
		}

		contractClass := types.ContractClass{}
		if err := json.Unmarshal(content, &contractClass); err != nil {
			t.Fatal(err)
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
				Filename:                "./tests/account.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: DevNetAccountAddress,
			},
		},
		"mainnet": {},
		"mock":    {},
		"testnet": {
			{
				Filename:                "./tests/counter.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{"0x1"},
				ExpectedContractAddress: "0x6a57b89a061930d1141bbfec7c4afecffa8dc8f75174420161991b994a9ad4f",
			},
			{
				Filename:                "./tests/account.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: "0x4916cb2ef37f886d7e35f6bdbb38d20917057efc4de7fad73143566f8db73a1",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		content, err := os.ReadFile(test.Filename)
		if err != nil {
			t.Fatal("should read file with success, instead:", err)
		}
		contractClass := types.ContractClass{}
		if err := json.Unmarshal(content, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		dec, err := testConfig.client.AddDeployTransaction(context.Background(), test.Salt, test.ConstructorCall, contractClass)
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

// TestInvokeTransaction tests starknet_addDeployTransaction
func TestInvokeTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		AccountPrivateKeyEnvVar string
		AccountPublicKey        string
		AccountAddress          string
		Call                    types.FunctionCall
		MaxFee                  string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			AccountPrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			AccountPublicKey:        TestPublicKey,
			AccountAddress:          TestNetAccountAddress,
			Call: types.FunctionCall{
				ContractAddress:    types.HexToHash("0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167"),
				EntryPointSelector: "increment",
				CallData:           []string{},
			},
			MaxFee: "0x200000001",
		}},
	}[testEnv]

	for _, test := range testSet {
		privateKey := os.Getenv(test.AccountPrivateKeyEnvVar)
		if privateKey == "" {
			t.Fatal("should have a private key for the account")
		}
		account, err := testConfig.client.NewAccount(privateKey, test.AccountAddress)
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		n, err := account.Nonce(context.Background())
		if err != nil {
			t.Fatal("should return nonce, instead", err)
		}
		maxFee, _ := big.NewInt(0).SetString(test.MaxFee, 0)
		spy := NewSpy(testConfig.client.c, false)
		testConfig.client.c = spy
		txHash, err := account.HashMultiCall(
			[]types.FunctionCall{test.Call},
			ExecuteDetails{
				Nonce:   n,
				MaxFee:  maxFee,
				Version: big.NewInt(0),
			},
		)
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		s1, s2, err := account.Sign(txHash)
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		calldata := fmtExecuteCalldataStrings(n, []types.FunctionCall{test.Call})
		output, err := testConfig.client.AddInvokeTransaction(
			context.Background(),
			types.FunctionCall{
				ContractAddress:    types.HexToHash(test.AccountAddress),
				EntryPointSelector: "__execute__",
				CallData:           calldata,
			},
			[]string{s1.Text(10), s2.Text(10)},
			test.MaxFee,
			"0x0",
		)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		if output.TransactionHash != fmt.Sprintf("0x%s", txHash.Text(16)) {
			t.Log("transaction error...")
			t.Logf("- computed:  %s", fmt.Sprintf("0x%s", txHash.Text(16)))
			t.Logf("- collected: %s", output.TransactionHash)
			t.FailNow()
		}
		if diff, err := spy.Compare(output, false); err != nil || diff != "FullMatch" {
			spy.Compare(output, true)
			t.Fatal("expecting to match", err)
		}
		fmt.Println("transaction_hash:", output.TransactionHash)
	}
}
