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

// TestDeclareTransaction is a test function for declaring a transaction.
//
// It sets up different test sets based on the test environment and runs the tests accordingly.
// Each test set contains a transaction hash, class hash, and an expected error.
// The function reads a JSON file, unmarshals it into an `AddDeclareTxnInput` struct,
// and performs various tests using the `AddDeclareTransaction` function.
// It checks if the error matches the expected error and if the transaction hash matches the expected hash.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//  none
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

		declareTxJSON, err := os.ReadFile("./tests/write/declareTx.json")
		if err != nil {
			t.Fatal("should be able to read file", err)
		}

		var declareTx AddDeclareTxnInput
		err = json.Unmarshal(declareTxJSON, &declareTx)
		require.Nil(t, err, "Error unmarshalling decalreTx")

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy

		// TODO: test transaction against client that supports RPC method (currently Sequencer uses
		// "sierra_program" instead of "program"
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

// TestAddInvokeTransaction is a test function that checks the AddInvokeTransaction functionality.
//
// It initializes a test configuration and defines sets of test cases for different environments.
// Each test case includes an InvokeTxnV1 object, an expected AddInvokeTransactionResponse, and an expected RPCError.
// The function iterates through the test cases, invokes AddInvokeTransaction, and compares the response and error with the expected values.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//  none
func TestAddInvokeTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		InvokeTx      InvokeTxnV1
		ExpectedResp  AddInvokeTransactionResponse
		ExpectedError RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {
			{
				InvokeTx:     InvokeTxnV1{SenderAddress: new(felt.Felt).SetUint64(123)},
				ExpectedResp: AddInvokeTransactionResponse{&felt.Zero},
				ExpectedError: RPCError{
					code:    ErrUnexpectedError.code,
					message: ErrUnexpectedError.message,
					data:    "Something crazy happened"},
			},
			{
				InvokeTx:      InvokeTxnV1{},
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
