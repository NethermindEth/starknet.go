package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
)

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeclareTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionHash *felt.Felt
		ClassHash       *felt.Felt
		ExpectedError   string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			TransactionHash: utils.TestHexToFelt(t, "0x55b094dc5c84c2042e067824f82da90988674314d37e45cb0032aca33d6e0b9"),
			ClassHash:       utils.TestHexToFelt(t, "0xdeadbeef"),
			ExpectedError:   "Invalid Params",
		}},
	}[testEnv]

	for _, test := range testSet {

		declareTxJSON, err := os.ReadFile("./tests/declareTx.json")
		if err != nil {
			t.Fatal("should be able to read file", err)
		}

		var declareTx BroadcastedDeclareTransaction
		err = json.Unmarshal(declareTxJSON, &declareTx)
		require.Nil(t, err, "Error unmarshalling decalreTx")

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy

		// To do: test transaction against client that supports RPC method (currently Sequencer uses
		// "sierra_program" instead of "program" in BroadcastedDeclareTransaction
		dec, err := testConfig.provider.AddDeclareTransaction(context.Background(), declareTx)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
			continue
		}
		if dec.TransactionHash != test.TransactionHash {
			t.Fatalf("classHash does not match expected, current: %s", dec.ClassHash)
		}

	}
}

// TestDeclareTransaction tests starknet_addDeclareTransaction
func TestDeployAccountTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionHash *felt.Felt
		ContractAddress *felt.Felt
		ExpectedError   string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {{
			TransactionHash: utils.TestHexToFelt(t, "0xdeadbeef"),
			ContractAddress: utils.TestHexToFelt(t, "0xdeadbeef"),
			ExpectedError:   "",
		}},
		"testnet": {{
			TransactionHash: utils.TestHexToFelt(t, "0x1"),
			ContractAddress: utils.TestHexToFelt(t, "0x1"),
			ExpectedError:   "",
		}},
	}[testEnv]

	deployTxJSON := `{"type": "DEPLOY_ACCOUNT",
	"max_fee": "0x24423",
	"version": "0x1",
	"signature": [
		"0x41c3543008dd65ed98c767e5d218b0c0ce1bd0cd60877824951a6f87cc1637d",
		"0x7f803845aa7e43d183fd05cd553c64711b1c49af69a155fe8144e8da9a5a50d"
	],
	"nonce": "0x0",
	"class_hash": "0x1fac3074c9d5282f0acc5c69a4781a1c711efea5e73c550c5d9fb253cf7fd3d",
	"contract_address_salt": "0x14e2ae44cbb50dff0e18140e7c415c1f281207d06fd6a0106caf3ff21e130d8",
	"constructor_calldata": [
		"0x22577c8898de00dca9588b7761cb1f4f5590f4e596fddf24e4a1465d1432f88"
	]}`

	// Currently this just tests the client doesn't return an error, and so
	// we assume everything went well. We need to pre-compute the transaction
	// hash to check the response.
	// NOTE : unstable test
	for _, test := range testSet {
		var deployAcntTx BroadcastedDeployAccountTransaction
		err := json.Unmarshal([]byte(deployTxJSON), &deployAcntTx)
		require.NoError(t, err)

		randPub, err := new(felt.Felt).SetRandom()
		require.NoError(t, err)
		deployAcntTx.ConstructorCalldata = []*felt.Felt{randPub} // to prevent transaction hash clashes
		fmt.Println(randPub)

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy

		resp, err := testConfig.provider.AddDeployAccountTransaction(context.Background(), deployAcntTx)
		fmt.Println(resp)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
			continue
		}

	}
}
