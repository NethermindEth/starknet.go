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
	expectedrespRaw, err := os.ReadFile("./tests/trace/sepoliaInvokeTrace_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json")
	require.NoError(t, err, "Error ReadFile for TestTraceTransaction")
	require.NoError(t, json.Unmarshal(expectedrespRaw, &expectedResp), "Error unmarshalling testdata TestTraceTransaction")

	type testSetType struct {
		TransactionHash *felt.Felt
		ExpectedResp    TxnTrace
		ExpectedError   error
	}
	testSet := map[string][]testSetType{
		"mock": {
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
				ExpectedResp:    expectedResp,
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
		"devnet": {},
		"testnet": {
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
				ExpectedResp:    expectedResp,
				ExpectedError:   nil,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceTransaction(context.Background(), test.TransactionHash)
		require.Equal(t, test.ExpectedError, err)
		compareTraceTxs(t, test.ExpectedResp, resp)
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
		simulateTxnRaw, err := os.ReadFile("./tests/trace/mainnetSimulateInvokeTx.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTx")
		require.NoError(t, json.Unmarshal(simulateTxnRaw, &simulateTxIn), "Error unmarshalling simulateInvokeTx")

		expectedrespRaw, err := os.ReadFile("./tests/trace/mainnetSimulateInvokeTxResp.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTxResp")
		require.NoError(t, json.Unmarshal(expectedrespRaw, &expectedResp), "Error unmarshalling simulateInvokeTxResp")
	}

	if testEnv == "testnet" || testEnv == "mock" {
		simulateTxnRaw, err := os.ReadFile("./tests/trace/sepoliaSimulateInvokeTx.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTx")
		require.NoError(t, json.Unmarshal(simulateTxnRaw, &simulateTxIn), "Error unmarshalling simulateInvokeTx")

		expectedrespRaw, err := os.ReadFile("./tests/trace/sepoliaSimulateInvokeTxResp.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTxResp")
		require.NoError(t, json.Unmarshal(expectedrespRaw, &expectedResp), "Error unmarshalling simulateInvokeTxResp")
	}

	type testSetType struct {
		SimulateTxnInput SimulateTransactionInput
		ExpectedResp     SimulateTransactionOutput
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mock": {testSetType{
			SimulateTxnInput: simulateTxIn,
			ExpectedResp:     expectedResp,
		}},
		"testnet": {testSetType{
			SimulateTxnInput: simulateTxIn,
			ExpectedResp:     expectedResp,
		}},
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

		for i, trace := range resp {
			require.Equal(t, test.ExpectedResp.Txns[i].FeeEstimation, trace.FeeEstimation)
			compareTraceTxs(t, test.ExpectedResp.Txns[i].TxnTrace, trace.TxnTrace)
		}
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
	require := require.New(t)

	var blockTraceSepolia []Trace

	expectedrespRaw, err := os.ReadFile("./tests/trace/sepoliaBlockTrace_0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2.json")
	require.NoError(err, "Error ReadFile for TestTraceBlockTransactions")
	require.NoError(json.Unmarshal(expectedrespRaw, &blockTraceSepolia), "Error unmarshalling testdata TestTraceBlockTransactions")

	type testSetType struct {
		BlockID      BlockID
		ExpectedResp []Trace
		ExpectedErr  *RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {}, // devenet doesn't support TraceBlockTransactions https://0xspaceshard.github.io/starknet-devnet/docs/guide/json-rpc-api#trace-api
		"mainnet": {},
		"testnet": {
			testSetType{
				BlockID:      WithBlockNumber(99433),
				ExpectedResp: blockTraceSepolia,
				ExpectedErr:  nil,
			},
		},
		"mock": {
			testSetType{
				BlockID:      WithBlockHash(utils.TestHexToFelt(t, "0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2")),
				ExpectedResp: blockTraceSepolia,
				ExpectedErr:  nil,
			},
			testSetType{
				BlockID:      WithBlockNumber(0),
				ExpectedResp: nil,
				ExpectedErr:  ErrBlockNotFound,
			}},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceBlockTransactions(context.Background(), test.BlockID)

		if err != nil {
			require.Equal(test.ExpectedErr, err)
		} else {
			for i, trace := range resp {
				require.Equal(test.ExpectedResp[i].TxnHash, trace.TxnHash)
				compareTraceTxs(t, test.ExpectedResp[i].TraceRoot, trace.TraceRoot)
			}
		}

	}
}

func compareTraceTxs(t *testing.T, traceTx1, traceTx2 TxnTrace) {
	require := require.New(t)

	switch traceTx := traceTx1.(type) {
	case DeclareTxnTrace:
		require.Equal(traceTx.ValidateInvocation, traceTx2.(DeclareTxnTrace).ValidateInvocation)
		require.Equal(traceTx.FeeTransferInvocation, traceTx2.(DeclareTxnTrace).FeeTransferInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(DeclareTxnTrace).StateDiff)
		require.Equal(traceTx.Type, traceTx2.(DeclareTxnTrace).Type)
		require.Equal(traceTx.ExecutionResources, traceTx2.(DeclareTxnTrace).ExecutionResources)
	case DeployAccountTxnTrace:
		require.Equal(traceTx.ValidateInvocation, traceTx2.(DeployAccountTxnTrace).ValidateInvocation)
		require.Equal(traceTx.ConstructorInvocation, traceTx2.(DeployAccountTxnTrace).ConstructorInvocation)
		require.Equal(traceTx.FeeTransferInvocation, traceTx2.(DeployAccountTxnTrace).FeeTransferInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(DeployAccountTxnTrace).StateDiff)
		require.Equal(traceTx.Type, traceTx2.(DeployAccountTxnTrace).Type)
		require.Equal(traceTx.ExecutionResources, traceTx2.(DeployAccountTxnTrace).ExecutionResources)
	case InvokeTxnTrace:
		require.Equal(traceTx.ValidateInvocation, traceTx2.(InvokeTxnTrace).ValidateInvocation)
		require.Equal(traceTx.ExecuteInvocation, traceTx2.(InvokeTxnTrace).ExecuteInvocation)
		require.Equal(traceTx.FeeTransferInvocation, traceTx2.(InvokeTxnTrace).FeeTransferInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(InvokeTxnTrace).StateDiff)
		require.Equal(traceTx.Type, traceTx2.(InvokeTxnTrace).Type)
		require.Equal(traceTx.ExecutionResources, traceTx2.(InvokeTxnTrace).ExecutionResources)
	case L1HandlerTxnTrace:
		require.Equal(traceTx.FunctionInvocation, traceTx2.(L1HandlerTxnTrace).FunctionInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(L1HandlerTxnTrace).StateDiff)
		require.Equal(traceTx.Type, traceTx2.(L1HandlerTxnTrace).Type)
	}
}

func compareStateDiffs(t *testing.T, stateDiff1, stateDiff2 StateDiff) {
	require.ElementsMatch(t, stateDiff1.DeprecatedDeclaredClasses, stateDiff2.DeprecatedDeclaredClasses)
	require.ElementsMatch(t, stateDiff1.DeclaredClasses, stateDiff2.DeclaredClasses)
	require.ElementsMatch(t, stateDiff1.DeployedContracts, stateDiff2.DeployedContracts)
	require.ElementsMatch(t, stateDiff1.ReplacedClasses, stateDiff2.ReplacedClasses)
	require.ElementsMatch(t, stateDiff1.Nonces, stateDiff2.Nonces)

	// compares storage diffs (they come in a random order)
	rawStorageDiff, err := json.Marshal(stateDiff2.StorageDiffs)
	require.NoError(t, err)
	var mapDiff []map[string]interface{}
	require.NoError(t, json.Unmarshal(rawStorageDiff, &mapDiff))

	for _, diff1 := range stateDiff1.StorageDiffs {
		var diff2 ContractStorageDiffItem

		for _, diffElem := range mapDiff {
			address, ok := diffElem["address"]
			require.True(t, ok)
			addressFelt := utils.TestHexToFelt(t, address.(string))

			if *addressFelt != *diff1.Address {
				continue
			}

			err = remarshal(diffElem, &diff2)
			require.NoError(t, err)
		}
		require.NotEmpty(t, diff2)

		require.Equal(t, diff1.Address, diff2.Address)
		require.ElementsMatch(t, diff1.StorageEntries, diff2.StorageEntries)
	}
}
