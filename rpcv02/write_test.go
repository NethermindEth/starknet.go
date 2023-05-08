package rpcv02

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/caigo/artifacts"
	"github.com/smartcontractkit/caigo/types"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Filename          []byte
		Version           TransactionVersion
		Signature         []string
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			Filename:          artifacts.CounterCompiled,
			Version:           TransactionV1,
			Signature:         []string{"0x62fe3eb3f30824c8039e5a8b9f16e02f6ad71e01d5593b330451fccd97a2576", "0x33eda2c4b010a3365a9ed453c0f21b796f84e03cc33d3602f49f57cd50d2270"},
			ExpectedClassHash: "0x029c64881bf658fae000fa6d5112f379eb4fc9c629a5cd7455eafc0744e34a8a",
		}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			Filename:          artifacts.CounterCompiled,
			Version:           TransactionV1,
			Signature:         []string{"0x62fe3eb3f30824c8039e5a8b9f16e02f6ad71e01d5593b330451fccd97a2576", "0x33eda2c4b010a3365a9ed453c0f21b796f84e03cc33d3602f49f57cd50d2270"},
			ExpectedClassHash: "0x29c64881bf658fae000fa6d5112f379eb4fc9c629a5cd7455eafc0744e34a8a",
		}},
	}[testEnv]

	for _, test := range testSet {
		contractClass := types.ContractClass{}
		if err := json.Unmarshal(test.Filename, &contractClass); err != nil {
			t.Fatal(err)
		}
		maxFee, _ := big.NewInt(0).SetString("10000000000000", 0)
		nonce := big.NewInt(1)
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		declareTransaction := BroadcastedDeclareTransaction{
			BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
				Version:   test.Version,
				MaxFee:    maxFee,
				Nonce:     nonce,
				Signature: test.Signature,
			},
			ContractClass: contractClass,
			SenderAddress: types.HexToHash(TestNetAccount040Address),
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
