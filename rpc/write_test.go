package rpc

import (
	"context"
	"encoding/json"
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

		var declareTx BroadcastedDeclareTransactionV1
		err = json.Unmarshal(declareTxJSON, &declareTx)
		require.Nil(t, err, "Error unmarshalling decalreTx")

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy

		// To do: test transaction against client that supports RPC method (currently Sequencer uses
		// "sierra_program" instead of "program" in BroadcastedDeclareTransactionV2)
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
func TestAddInvokeTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		InvokeTx      BroadcastedInvokeTransaction
		ExpectedResp  AddInvokeTransactionResponse
		ExpectedError RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {
			{
				InvokeTx:     BroadcastedInvokeV1Transaction{SenderAddress: new(felt.Felt).SetUint64(123)},
				ExpectedResp: AddInvokeTransactionResponse{&felt.Zero},
				ExpectedError: RPCError{
					code:    ErrUnexpectedError.code,
					message: ErrUnexpectedError.message,
					data:    "Something crazy happened"},
			},
			{
				InvokeTx:      BroadcastedInvokeV1Transaction{},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0xdeadbeef")},
				ExpectedError: RPCError{},
			},
		},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.AddInvokeTransaction(context.Background(), test.InvokeTx)
		if err != nil {
			require.Equal(t, err, &test.ExpectedError, "AddInvokeTransaction did not give expected error")
		} else {
			require.Equal(t, *resp, test.ExpectedResp)
		}

	}
}
