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
		},
		"mock": {},
		"testnet": {
			{MissingMethod: "starknet_traceTransaction"},
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

// TestTraceBlockTransactions tests the TraceBlockTransactions method
func TestTraceBlockTransactions(t *testing.T) {
	testConfig := beforeEach(t)

	var expectedResp []Trace
	if testEnv == "mock" {
		var rawjson struct {
			Result []Trace `json:"result"`
		}
		expectedrespRaw, err := os.ReadFile("./tests/0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2.json")
		require.NoError(t, err, "Error ReadFile for TestTraceBlockTransactions")

		err = json.Unmarshal(expectedrespRaw, &rawjson)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceBlockTransactions")
		expectedResp = rawjson.Result
	}

	type testSetType struct {
		BlockHash    *felt.Felt
		ExpectedResp []Trace
		ExpectedErr  *RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {}, // devenet doesn't support TraceBlockTransactions https://0xspaceshard.github.io/starknet-devnet/docs/guide/json-rpc-api#trace-api
		"mainnet": {},
		"mock": {
			testSetType{
				BlockHash:    utils.TestHexToFelt(t, "0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2"),
				ExpectedResp: expectedResp,
				ExpectedErr:  nil,
			},
			testSetType{
				BlockHash:    utils.TestHexToFelt(t, "0x0"),
				ExpectedResp: nil,
				ExpectedErr:  ErrInvalidBlockHash,
			}},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceBlockTransactions(context.Background(), test.BlockHash)

		if err != nil {
			require.Equal(t, test.ExpectedErr, err)
		} else {
			require.Equal(t, test.ExpectedResp, *resp)
		}

	}
}
