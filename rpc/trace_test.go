package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/test-go/testify/require"
)

// TestNotImplemented checks the method not yet implemented
func TestNotImplemented(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		MissingMethod string
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mainnet": {
			{MissingMethod: "starknet_traceTransaction"},
			{MissingMethod: "starknet_traceBlockTransactions"},
		},
		"mock": {},
		"testnet": {
			{MissingMethod: "starknet_traceTransaction"},
			{MissingMethod: "starknet_traceBlockTransactions"},
		},
	}[testEnv]

	for _, test := range testSet {
		var out string
		err := do(context.Background(), testConfig.provider.c, test.MissingMethod, &out)

		if err == nil || err.Error() != "Method Not Found" {
			t.Fatalf("Method %s is now available, got %v\n", test.MissingMethod, err)
		}
	}
}

// TestSimulateTransaction tests the SimulateTransaction method
func TestSimulateTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	var simulateTxIn SimulateTransactionInput
	var expectedResp SimulateTransactionOutput
	if testEnv == "mainnet" {
		simulateTxnRaw, err := os.ReadFile("./tests/simulateInvokeTx.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTx")

		err = json.Unmarshal(simulateTxnRaw, &simulateTxIn)
		require.NoError(t, err, "Error unmarshalling simulateInvokeTx")

		expectedrespRaw, err := os.ReadFile("./tests/simulateInvokeTxResp.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTxResp")

		err = json.Unmarshal(expectedrespRaw, &expectedResp)
		require.NoError(t, err, "Error unmarshalling simulateInvokeTxResp")
	}

	type testSetType struct {
		SimulateTxnInput SimulateTransactionInput
		ExpectedResp     SimulateTransactionOutput
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mock":    {},
		"testnet": {},
		"mainnet": {testSetType{
			SimulateTxnInput: simulateTxIn,
			ExpectedResp:     expectedResp,
		}},
	}[testEnv]

	for _, test := range testSet {

		resp, err := testConfig.provider.SimulateTransactions(
			context.Background(),
			test.SimulateTxnInput.BlockID,
			test.SimulateTxnInput.Txns,
			test.SimulateTxnInput.SimulationFlags)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedResp.Txns, resp)
	}
}
