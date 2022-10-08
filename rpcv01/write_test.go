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
		"testnet": {
			{
				Filename:                "./tests/counter.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{"0x1"},
				ExpectedContractAddress: "0x357b37bf12f59dd04c4da4933dcadf4a104e158365886d64ca0e554ada68fef",
			},
			{
				Filename:                "./tests/oz_v0.3.2_account.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: TestNetAccount032Address,
			},
			{
				Filename:                "./tests/oz_v0.4.0b_account.json",
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: TestNetAccount040Address,
			},
		},
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

// TestInvokeTransaction_InvokeTxvV0 tests starknet_addInvokeTransaction with a V0 account
// func TestInvokeTransaction_InvokeTxvV0(t *testing.T) {
// 	testConfig := beforeEach(t)

// 	type testSetType struct {
// 		AccountPrivateKeyEnvVar string
// 		AccountAddress          string
// 		Call                    ctypes.FunctionCall
// 		MaxFee                  string
// 	}
// 	testSet := map[string][]testSetType{
// 		"devnet":  {},
// 		"mainnet": {},
// 		"mock":    {},
// 		"testnet": {
// 			// Disabled tests due to the fact it is taking ages on the CI. It should
// 			// work on demand though...
// 			// {
// 			// 	AccountPrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
// 			// 	Version:                 AccountVersion0,
// 			// 	AccountAddress:          TestNetAccount032Address,
// 			// 	Call: types.FunctionCall{
// 			// 		ContractAddress:    types.HexToHash("0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167"),
// 			// 		EntryPointSelector: "increment",
// 			// 		Calldata:           []string{},
// 			// 	},
// 			// 	MaxFee: "0x200000000",
// 			// },
// 		},
// 	}[testEnv]

// 	for _, test := range testSet {
// 		privateKey := os.Getenv(test.AccountPrivateKeyEnvVar)
// 		if privateKey == "" {
// 			t.Fatal("should have a private key for the account")
// 		}
// 		account, err := testConfig.provider.NewAccount(privateKey, test.AccountAddress, AccountVersion0)
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		if account.version.Cmp(big.NewInt(0)) != 0 {
// 			t.Fatalf("This test only supports version v0, current: %s", account.version.Text(10))
// 		}
// 		ctx := context.Background()
// 		n, err := account.Nonce(ctx)
// 		if err != nil {
// 			t.Fatal("should return nonce, instead", err)
// 		}
// 		maxFee, _ := big.NewInt(0).SetString(test.MaxFee, 0)
// 		spy := NewSpy(testConfig.provider.c, false)
// 		testConfig.provider.c = spy
// 		txHash, err := account.TransactionHash(
// 			[]ctypes.FunctionCall{test.Call},
// 			types.ExecuteDetails{
// 				Nonce:  n,
// 				MaxFee: maxFee,
// 			},
// 		)
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		s1, s2, err := account.Sign(txHash)
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		calldata := fmtV0CalldataStrings(n, []ctypes.FunctionCall{test.Call})
// 		output, err := testConfig.provider.AddInvokeTransaction(
// 			ctx,
// 			ctypes.FunctionCall{
// 				ContractAddress:    ctypes.HexToHash(test.AccountAddress),
// 				EntryPointSelector: EXECUTE_SELECTOR,
// 				Calldata:           calldata,
// 			},
// 			[]string{s1.Text(10), s2.Text(10)},
// 			test.MaxFee,
// 			"0x0",
// 		)
// 		if err != nil {
// 			t.Fatal("declare should succeed, instead:", err)
// 		}
// 		if output.TransactionHash != fmt.Sprintf("0x%s", txHash.Text(16)) {
// 			t.Log("transaction error...")
// 			t.Logf("- computed:  %s", fmt.Sprintf("0x%s", txHash.Text(16)))
// 			t.Logf("- collected: %s", output.TransactionHash)
// 			t.FailNow()
// 		}
// 		if diff, err := spy.Compare(output, false); err != nil || diff != "FullMatch" {
// 			spy.Compare(output, true)
// 			t.Fatal("expecting to match", err)
// 		}
// 		fmt.Println("transaction_hash:", output.TransactionHash)
// 		ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
// 		defer cancel()
// 		status, err := account.Provider.WaitForTransaction(ctx, ctypes.HexToHash(output.TransactionHash), 8*time.Second)
// 		if err != nil {
// 			t.Fatal("declare should succeed, instead:", err)
// 		}
// 		if status != "PENDING" && status != "ACCEPTED_ON_L1" && status != "ACCEPTED_ON_L2" {
// 			t.Fatalf("tx %s wrong status: %s", output.TransactionHash, status)
// 		}
// 	}
// }

// // TestInvokeTransaction_InvokeTxvV1 tests starknet_addInvokeTransaction with a V1 account
// func TestInvokeTransaction_InvokeTxvV1(t *testing.T) {
// 	testConfig := beforeEach(t)

// 	type testSetType struct {
// 		AccountPrivateKeyEnvVar string
// 		AccountAddress          string
// 		Call                    ctypes.FunctionCall
// 		MaxFee                  string
// 	}
// 	testSet := map[string][]testSetType{
// 		"devnet":  {},
// 		"mainnet": {},
// 		"mock":    {},
// 		"testnet": {
// 			// Disabled tests due to the fact it is taking ages on the CI. It should
// 			// work on demand though...
// 			// {
// 			// 	AccountPrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
// 			// 	AccountAddress:          TestNetAccount040Address,
// 			// 	Call: types.FunctionCall{
// 			// 		ContractAddress:    types.HexToHash("0x37a2490365294ef4bc896238642b9bcb0203f86e663f11688bb86c5e803c167"),
// 			// 		EntryPointSelector: "increment",
// 			// 		Calldata:           []string{},
// 			// 	},
// 			// 	MaxFee: "0x200000001",
// 			// },
// 		},
// 	}[testEnv]

// 	for _, test := range testSet {
// 		privateKey := os.Getenv(test.AccountPrivateKeyEnvVar)
// 		if privateKey == "" {
// 			t.Fatal("should have a private key for the account")
// 		}
// 		account, err := testConfig.provider.NewAccount(privateKey, test.AccountAddress, AccountVersion1)
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		if account.version.Cmp(big.NewInt(1)) != 0 {
// 			t.Fatalf("This test only supports version v1, current: %s", account.version.Text(10))
// 		}
// 		ctx := context.Background()
// 		n, err := account.Nonce(ctx)
// 		if err != nil {
// 			t.Fatal("should return nonce, instead", err)
// 		}
// 		maxFee, _ := big.NewInt(0).SetString(test.MaxFee, 0)
// 		spy := NewSpy(testConfig.provider.c, false)
// 		testConfig.provider.c = spy
// 		txHash, err := account.TransactionHash(
// 			[]ctypes.FunctionCall{test.Call},
// 			types.ExecuteDetails{
// 				Nonce:  n,
// 				MaxFee: maxFee,
// 			},
// 		)
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		s1, s2, err := account.Sign(txHash)
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		calldata := fmtCalldataStrings([]ctypes.FunctionCall{test.Call})
// 		output, err := testConfig.provider.AddInvokeTransaction(
// 			ctx,
// 			ctypes.FunctionCall{
// 				ContractAddress: ctypes.HexToHash(test.AccountAddress),
// 				Calldata:        calldata,
// 			},
// 			[]string{s1.Text(10), s2.Text(10)},
// 			test.MaxFee,
// 			"0x0",
// 		)
// 		if err != nil {
// 			t.Fatal("declare should succeed, instead:", err)
// 		}
// 		if output.TransactionHash != fmt.Sprintf("0x%s", txHash.Text(16)) {
// 			t.Log("transaction error...")
// 			t.Logf("- computed:  %s", fmt.Sprintf("0x%s", txHash.Text(16)))
// 			t.Logf("- collected: %s", output.TransactionHash)
// 			t.FailNow()
// 		}
// 		if diff, err := spy.Compare(output, false); err != nil || diff != "FullMatch" {
// 			spy.Compare(output, true)
// 			t.Fatal("expecting to match", err)
// 		}
// 		fmt.Println("transaction_hash:", output.TransactionHash)
// 		ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
// 		defer cancel()
// 		status, err := account.Provider.WaitForTransaction(ctx, ctypes.HexToHash(output.TransactionHash), 8*time.Second)
// 		if err != nil {
// 			t.Fatal("declare should succeed, instead:", err)
// 		}
// 		if status != "PENDING" && status != "ACCEPTED_ON_L1" && status != "ACCEPTED_ON_L2" {
// 			t.Fatalf("tx %s wrong status: %s", output.TransactionHash, status)
// 		}
// 	}
// }
