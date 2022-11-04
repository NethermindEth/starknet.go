package rpcv02

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/dontpanicdao/caigo/artifacts"
	"github.com/dontpanicdao/caigo/types"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename          []byte
		Version           *big.Int
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			Filename:          artifacts.CounterCompiled,
			Version:           big.NewInt(0),
			ExpectedClassHash: "0x01649a376a9aa5ccb5ddf2f59c267de5fb6b3b177056a53f45d42877c856a051",
		}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:          artifacts.CounterCompiled,
			Version:           big.NewInt(0),
			ExpectedClassHash: "0x4484265a6e003e8afe272e6c9bf3e7d0d8e343b2df57763a995828285fdfbbd",
		}},
	}[testEnv]

	for _, test := range testSet {
		contractClass := types.ContractClass{}
		if err := json.Unmarshal(test.Filename, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		declareTransaction := BroadcastedDeclareTransaction{
			BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
				Version: test.Version,
			},
			ContractClass: contractClass,
		}
		dec, err := testConfig.provider.AddDeclareTransaction(context.Background(), declareTransaction)
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
		Filename                []byte
		Salt                    string
		ConstructorCall         []string
		ExpectedContractAddress string
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Filename:                artifacts.CounterCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{},
				ExpectedContractAddress: "0xbaaa96effb3564b6047e45944e8db9d9b0a056886d131038baabb56a959390",
			},
			{
				Filename:                artifacts.AccountV0Compiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: TestNetAccount032Address,
			},
			{
				Filename:                artifacts.AccountCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: TestNetAccount040Address,
			},
		},
		"mainnet": {},
		"mock":    {},
		"testnet": {
			{
				Filename:                artifacts.CounterCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{},
				ExpectedContractAddress: "0xbaaa96effb3564b6047e45944e8db9d9b0a056886d131038baabb56a959390",
			},
			{
				Filename:                artifacts.AccountV0Compiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: TestNetAccount032Address,
			},
			{
				Filename:                artifacts.AccountCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: TestNetAccount040Address,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		contractClass := types.ContractClass{}
		if err := json.Unmarshal(test.Filename, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		broadcastedDeployTransaction := BroadcastedDeployTransaction{
			Type:                "DEPLOY",
			Version:             big.NewInt(0),
			ContractAddressSalt: test.Salt,
			ConstructorCalldata: test.ConstructorCall,
			ContractClass:       contractClass,
		}
		dec, err := testConfig.provider.AddDeployTransaction(context.Background(), broadcastedDeployTransaction)
		if err != nil {
			t.Fatal("declare should succeed, instead:", err)
		}
		fmt.Printf("transaction hash: %s\n", dec.TransactionHash)
		expectedContractAddress, _ := big.NewInt(0).SetString(test.ExpectedContractAddress, 0)
		contractAddress, _ := big.NewInt(0).SetString(dec.ContractAddress, 0)
		if contractAddress.Cmp(expectedContractAddress) != 0 {
			t.Fatalf("contractAddress expecting %s, current: %s", test.ExpectedContractAddress, dec.ContractAddress)
		}
		if diff, err := spy.Compare(dec, false); err != nil || diff != "FullMatch" {
			spy.Compare(dec, true)
			t.Fatal("expecting to match", err)
		}
	}
}
