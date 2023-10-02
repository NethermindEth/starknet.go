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

// TestTransactionTrace tests the TransactionTrace function.
//
// It sets up the necessary test configuration, including the expected response from a mocked environment.
// It then iterates through a test set, calling the TransactionTrace function with different transaction hashes.
// It compares the returned response with the expected response for each test case.
// If the response is not nil, it checks if it matches the expected response.
// If the response is nil, it checks if the expected error matches the actual error.
func TestTransactionTrace(t *testing.T) {
	testConfig := beforeEach(t)

	var expectedResp InvokeTxnTrace
	if testEnv == "mock" {
		var rawjson struct {
			Result InvokeTxnTrace `json:"result"`
		}
		expectedrespRaw, err := os.ReadFile("./tests/0xff66e14fc6a96f3289203690f5f876cb4b608868e8549b5f6a90a21d4d6329.json")
		require.NoError(t, err, "Error ReadFile for TestTraceTransaction")

		err = json.Unmarshal(expectedrespRaw, &rawjson)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceTransaction")

		txnTrace, err := json.Marshal(rawjson.Result)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceTransaction")
		err = json.Unmarshal(txnTrace, &expectedResp)
	}

	type testSetType struct {
		TransactionHash *felt.Felt
		ExpectedResp    TxnTrace
		ExpectedError   *RPCError
	}
	testSet := map[string][]testSetType{
		"mock": {
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xff66e14fc6a96f3289203690f5f876cb4b608868e8549b5f6a90a21d4d6329"),
				ExpectedResp:    expectedResp,
				ExpectedError:   nil,
			},
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xc0ffee"),
				ExpectedResp:    nil,
				ExpectedError:   ErrInvalidTxnHash,
			},
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xf00d"),
				ExpectedResp:    nil,
				ExpectedError: &RPCError{
					code:    10,
					message: "No trace available for transaction",
					data:    "REJECTED",
				},
			},
		},
		"devnet":  {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TransactionTrace(context.Background(), test.TransactionHash)
		if err != nil {
			require.Equal(t, test.ExpectedError, err)
		} else {
			require.Equal(t, test.ExpectedResp, resp)
		}
	}
}

// TestSimulateTransaction tests the SimulateTransaction function.
//
// It sets up the necessary test configuration and variables. If the test environment is "mainnet", it reads the simulateInvokeTx.json file and unmarshals the data into the `simulateTxIn` variable, as well as the simulateInvokeTxResp.json file and unmarshals the data into the `expectedResp` variable. 
//
// The function then iterates through the testSet and calls the `SimulateTransactions` function of the provider with the given inputs. It compares the response with the expected response and ensures that they are equal.
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

// TestTraceBlockTransactions tests the TraceBlockTransactions function.
//
// It initializes the test configuration and sets up the expected response. It then iterates over a test set and calls the TraceBlockTransactions function for each test case. It compares the actual response with the expected response and asserts the equality. If there is an error, it also asserts the equality of the error with the expected error. The function uses the testing.T object for assertions.
//
// Parameters:
// - t: The testing.T object for assertions.
//
// Return type: None.
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
			require.Equal(t, test.ExpectedResp, resp)
		}

	}
}
