package rpcv02

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/dontpanicdao/caigo/artifacts"
	"github.com/dontpanicdao/caigo/types"
	ctypes "github.com/dontpanicdao/caigo/types"
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
				ConstructorCall:         []string{"0x1"},
				ExpectedContractAddress: "0x035a55a64238b776664d7723de1f6b50350116a1ab1ca1fe154320a0eba53d3a",
			},
			{
				Filename:                artifacts.AccountV0Compiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: DevNetAccount032Address,
			},
			{
				Filename:                artifacts.AccountCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{TestPublicKey},
				ExpectedContractAddress: DevNetAccount040Address,
			},
		},
		"mainnet": {},
		"mock":    {},
		"testnet": {
			{
				Filename:                artifacts.CounterCompiled,
				Salt:                    "0xdeadbeef",
				ConstructorCall:         []string{"0x1"},
				ExpectedContractAddress: "0x357b37bf12f59dd04c4da4933dcadf4a104e158365886d64ca0e554ada68fef",
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
		contractClass := ctypes.ContractClass{}
		if err := json.Unmarshal(test.Filename, &contractClass); err != nil {
			t.Fatal(err)
		}

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		broadcastedDeployTransaction := BroadcastedDeployTransaction{
			Version:             big.NewInt(0),
			ContractAddressSalt: test.Salt,
			ConstructorCalldata: test.ConstructorCall,
			ContractClass:       contractClass,
		}
		dec, err := testConfig.provider.AddDeployTransaction(context.Background(), broadcastedDeployTransaction)
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
