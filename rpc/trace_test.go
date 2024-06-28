package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

// TestTransactionTrace is a function that tests the TransactionTrace function.
//
// It sets up the necessary test configuration and expected response. Then it performs a series of test sets,
// each with a different transaction hash. For each test set, it calls the TransactionTrace function and compares
// the response with the expected response.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestTransactionTrace(t *testing.T) {
	testConfig := beforeEach(t)

	var expectedResp InvokeTxnTrace
	if testEnv == "mock" {
		var rawjson struct {
			Result InvokeTxnTrace `json:"result"`
		}
		expectedrespRaw, err := os.ReadFile("./tests/trace/0x4b861c47d0fbc4cc24dacf92cf155ad0a2f7e2a0fd9b057b90cdd64eba7e12e.json")
		require.NoError(t, err, "Error ReadFile for TestTraceTransaction")

		err = json.Unmarshal(expectedrespRaw, &rawjson)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceTransaction")

		txnTrace, err := json.Marshal(rawjson.Result)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceTransaction")
		require.NoError(t, json.Unmarshal(txnTrace, &expectedResp))
	}

	type testSetType struct {
		TransactionHash *felt.Felt
		ExpectedResp    *InvokeTxnTrace
		ExpectedError   *RPCError
	}
	testSet := map[string][]testSetType{
		"mock": {
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0x4b861c47d0fbc4cc24dacf92cf155ad0a2f7e2a0fd9b057b90cdd64eba7e12e"),
				ExpectedResp:    &expectedResp,
				ExpectedError:   nil,
			},
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xc0ffee"),
				ExpectedResp:    nil,
				ExpectedError:   ErrHashNotFound,
			},
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xf00d"),
				ExpectedResp:    nil,
				ExpectedError: &RPCError{
					Code:    10,
					Message: "No trace available for transaction",
					Data:    "REJECTED",
				},
			},
		},
		"devnet":  {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceTransaction(context.Background(), test.TransactionHash)
		if err != nil {
			require.Equal(t, test.ExpectedError, err)
		} else {
			invokeTrace := resp.(InvokeTxnTrace)
			require.Equal(t, invokeTrace, *test.ExpectedResp)
		}
	}
}

// TestSimulateTransaction is a function that tests the SimulateTransaction function in the codebase.
//
// It sets up the necessary test configuration and variables, and then performs a series of tests based on the test environment.
// The function reads input data from JSON files and performs JSON unmarshalling to set the values of the simulateTxIn and expectedResp variables.
// It then iterates over the testSet map, calling the SimulateTransactions function with the appropriate parameters and checking the response against the expected response.
// The function uses the testing.T type to report any errors or failures during the test execution.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestSimulateTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	var simulateTxIn SimulateTransactionInput
	var expectedResp SimulateTransactionOutput
	if testEnv == "mainnet" {
		simulateTxnRaw, err := os.ReadFile("./tests/trace/simulateInvokeTx.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTx")

		err = json.Unmarshal(simulateTxnRaw, &simulateTxIn)
		require.NoError(t, err, "Error unmarshalling simulateInvokeTx")

		expectedrespRaw, err := os.ReadFile("./tests/trace/simulateInvokeTxResp.json")
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

// TestTraceBlockTransactions tests the TraceBlockTransactions function.
//
// It sets up the test configuration and expected response. It then iterates
// through the test set, calling TraceBlockTransactions with the provided
// block hash. It checks if there is an error and compares the response with
// the expected response.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestTraceBlockTransactions(t *testing.T) {
	testConfig := beforeEach(t)

	var expectedResp []Trace
	if testEnv == "mock" {
		var rawjson struct {
			Result []Trace `json:"result"`
		}
		expectedrespRaw, err := os.ReadFile("./tests/trace/0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2.json")
		require.NoError(t, err, "Error ReadFile for TestTraceBlockTransactions")

		err = json.Unmarshal(expectedrespRaw, &rawjson)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceBlockTransactions")
		expectedResp = rawjson.Result
	}

	type testSetType struct {
		BlockID      BlockID
		ExpectedResp []Trace
		ExpectedErr  *RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {}, // devenet doesn't support TraceBlockTransactions https://0xspaceshard.github.io/starknet-devnet/docs/guide/json-rpc-api#trace-api
		"mainnet": {},
		"mock": {
			testSetType{
				BlockID:      BlockID{Hash: utils.TestHexToFelt(t, "0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2")},
				ExpectedResp: expectedResp,
				ExpectedErr:  nil,
			},
			testSetType{
				BlockID:      BlockID{Hash: utils.TestHexToFelt(t, "0x0")},
				ExpectedResp: nil,
				ExpectedErr:  ErrBlockNotFound,
			}},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceBlockTransactions(context.Background(), test.BlockID)

		if err != nil {
			require.Equal(t, test.ExpectedErr, err)
		} else {
			require.Equal(t, test.ExpectedResp, resp)
		}

	}
}
