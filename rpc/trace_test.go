package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransactionTrace is a function that tests the TransactionTrace function.
//
// It sets up the necessary test configuration and expected response. Then it performs a series of test sets,
// each with a different transaction hash. For each test set, it calls the TransactionTrace function and compares
// the response with the expected response.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestTransactionTrace(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)
	testConfig := BeforeEach(t, false)

	expectedFile1 := "./testData/trace/sepoliaInvokeTrace_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json"

	type testSetType struct {
		TransactionHash  *felt.Felt
		ExpectedRespFile string
		ExpectedError    error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			testSetType{
				TransactionHash:  internalUtils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
				ExpectedRespFile: expectedFile1,
				ExpectedError:    nil,
			},
			testSetType{
				TransactionHash:  internalUtils.TestHexToFelt(t, "0xc0ffee"),
				ExpectedRespFile: expectedFile1,
				ExpectedError:    ErrHashNotFound,
			},
			testSetType{
				TransactionHash:  internalUtils.TestHexToFelt(t, "0xf00d"),
				ExpectedRespFile: expectedFile1,
				ExpectedError: &RPCError{
					Code:    10,
					Message: "No trace available for transaction",
					Data:    &TraceStatusErrData{Status: TraceStatusRejected},
				},
			},
		},
		tests.TestnetEnv: {
			testSetType{ // with 5 out of 6 fields (without state diff)
				TransactionHash:  internalUtils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
				ExpectedRespFile: expectedFile1,
				ExpectedError:    nil,
			},
			testSetType{ // with 6 out of 6 fields
				TransactionHash:  internalUtils.TestHexToFelt(t, "0x49d98a0328fee1de19d43d950cbaeb973d080d0c74c652523371e034cc0bbb2"),
				ExpectedRespFile: "./testData/trace/sepoliaInvokeTrace_0x49d98a0328fee1de19d43d950cbaeb973d080d0c74c652523371e034cc0bbb2.json",
				ExpectedError:    nil,
			},
		},
		tests.IntegrationEnv: {
			testSetType{
				TransactionHash:  internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
				ExpectedRespFile: "./testData/trace/integrationInvokeTrace_0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019.json",
				ExpectedError:    nil,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.TransactionHash.String(), func(t *testing.T) {
			expectedResp := *internalUtils.TestUnmarshalJSONFileToType[InvokeTxnTrace](t, test.ExpectedRespFile, "")

			resp, err := testConfig.Provider.TraceTransaction(
				context.Background(),
				test.TransactionHash,
			)
			if test.ExpectedError != nil {
				assert.EqualError(t, test.ExpectedError, err.Error())

				return
			}
			compareTraceTxs(t, expectedResp, resp)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			rawExpectedResp, err := os.ReadFile(test.ExpectedRespFile)
			require.NoError(t, err)

			compareTraceTxnsJSON(t, rawExpectedResp, rawResp)
		})
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
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestSimulateTransaction(t *testing.T) {
	testConfig := BeforeEach(t, false)

	type testSetType struct {
		SimulateTxnInputFile string
		ExpectedRespFile     string

		// quick fix for calls with other block ID than the one in the input file
		AnotherBlockID *BlockID
	}

	expectedInputFile := "./testData/trace/sepoliaSimulateInvokeTx.json"
	expectedRespFile := "./testData/trace/sepoliaSimulateInvokeTxResp.json"

	testSet := map[tests.TestEnv][]testSetType{
		tests.DevnetEnv: {},
		tests.MockEnv: {testSetType{
			SimulateTxnInputFile: expectedInputFile,
			ExpectedRespFile:     expectedRespFile,
		}},
		tests.TestnetEnv: {
			{
				SimulateTxnInputFile: expectedInputFile,
				ExpectedRespFile:     expectedRespFile,
			},
			{
				SimulateTxnInputFile: expectedInputFile,
				ExpectedRespFile:     expectedRespFile,
				AnotherBlockID:       &BlockID{Tag: BlockTagLatest},
			},
			{
				SimulateTxnInputFile: expectedInputFile,
				ExpectedRespFile:     expectedRespFile,
				AnotherBlockID:       &BlockID{Tag: BlockTagPreConfirmed},
			},
			{
				SimulateTxnInputFile: expectedInputFile,
				ExpectedRespFile:     expectedRespFile,
				AnotherBlockID:       &BlockID{Tag: BlockTagL1Accepted},
			},
		},
		// TODO: add mainnet test cases. I couldn't find a valid v3 transaction on mainnet with all resource bounds fields filled
		tests.MainnetEnv: {},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(fmt.Sprintf("blockID: %v", test.AnotherBlockID), func(t *testing.T) {
			simulateTxIn := *internalUtils.TestUnmarshalJSONFileToType[SimulateTransactionInput](t, test.SimulateTxnInputFile, "params")
			expectedResp := *internalUtils.TestUnmarshalJSONFileToType[[]SimulatedTransaction](t, test.ExpectedRespFile, "result")

			if test.AnotherBlockID != nil {
				simulateTxIn.BlockID = *test.AnotherBlockID
			}

			resp, err := testConfig.Provider.SimulateTransactions(
				context.Background(),
				simulateTxIn.BlockID,
				simulateTxIn.Txns,
				simulateTxIn.SimulationFlags)
			require.NoError(t, err)

			if test.AnotherBlockID != nil {
				// since the block ID is not the same as the one in the input file, we only check that the response is not empty
				assert.NotEmpty(t, resp)

				return
			}

			// read file to compare JSONs
			expectedRespArr := *internalUtils.TestUnmarshalJSONFileToType[[]any](t, test.ExpectedRespFile, "result")

			//nolint:dupl // Similar to TestTraceBlockTransactions, but they're testing different things.
			for i, trace := range resp {
				assert.Equal(t, expectedResp[i].FeeEstimation, trace.FeeEstimation)
				compareTraceTxs(t, expectedResp[i].TxnTrace, trace.TxnTrace)

				// compare JSONs
				// get fee_estimation and transaction_trace from expected response JSON file
				expectedRespMap, ok := expectedRespArr[i].(map[string]any)
				require.True(t, ok)
				expectedFeeEstimation, ok := expectedRespMap["fee_estimation"]
				require.True(t, ok)
				expectedTxnTrace, ok := expectedRespMap["transaction_trace"]
				require.True(t, ok)

				// compare fee_estimation
				rawExpectedFeeEstimation, err := json.Marshal(expectedFeeEstimation)
				require.NoError(t, err)
				rawActualFeeEstimation, err := json.Marshal(trace.FeeEstimation)
				require.NoError(t, err)

				assert.JSONEq(t, string(rawExpectedFeeEstimation), string(rawActualFeeEstimation))

				// compare transaction_trace
				rawExpectedTxnTrace, err := json.Marshal(expectedTxnTrace)
				require.NoError(t, err)
				rawActualTxnTrace, err := json.Marshal(trace.TxnTrace)
				require.NoError(t, err)

				compareTraceTxnsJSON(t, rawExpectedTxnTrace, rawActualTxnTrace)
			}
		})
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
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestTraceBlockTransactions(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.MockEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		BlockID          BlockID
		ExpectedRespFile string
		ExpectedErr      *RPCError
	}

	expectedRespFile := "./testData/trace/sepoliaBlockTrace_0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2.json"

	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {
			{
				BlockID:          WithBlockNumber(99433),
				ExpectedRespFile: expectedRespFile,
				ExpectedErr:      nil,
			},
			{
				BlockID:          WithBlockTag(BlockTagLatest),
				ExpectedRespFile: "",
				ExpectedErr:      nil,
			},
			{
				BlockID:          WithBlockTag(BlockTagL1Accepted),
				ExpectedRespFile: "",
				ExpectedErr:      nil,
			},
		},
		tests.MockEnv: {
			testSetType{
				BlockID:          WithBlockHash(internalUtils.TestHexToFelt(t, "0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2")),
				ExpectedRespFile: expectedRespFile,
				ExpectedErr:      nil,
			},
			testSetType{
				BlockID:          WithBlockNumber(0),
				ExpectedRespFile: expectedRespFile,
				ExpectedErr:      ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(fmt.Sprintf("blockID: %v", test.BlockID), func(t *testing.T) {
			resp, err := testConfig.Provider.TraceBlockTransactions(
				context.Background(),
				test.BlockID,
			)
			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)

				return
			}
			require.NoError(t, err)

			if test.ExpectedRespFile == "" {
				assert.NotEmpty(t, resp)

				return
			}
			expectedTrace := *internalUtils.TestUnmarshalJSONFileToType[[]Trace](t, test.ExpectedRespFile, "result")

			// read file to compare JSONs
			expectedRespArr := *internalUtils.TestUnmarshalJSONFileToType[[]any](t, test.ExpectedRespFile, "result")

			//nolint:dupl // Similar to TestSimulateTransaction, but they're testing different things.
			for i, actualTrace := range resp {
				require.Equal(t, expectedTrace[i].TxnHash, actualTrace.TxnHash)
				compareTraceTxs(t, expectedTrace[i].TraceRoot, actualTrace.TraceRoot)

				// compare JSONs
				// get transaction_hash and trace_root from expected response JSON file
				expectedRespMap, ok := expectedRespArr[i].(map[string]any)
				require.True(t, ok)
				expectedTxHash, ok := expectedRespMap["transaction_hash"]
				require.True(t, ok)
				expectedTxnTrace, ok := expectedRespMap["trace_root"]
				require.True(t, ok)

				// compare transaction_hash
				rawExpectedTxHash, err := json.Marshal(expectedTxHash)
				require.NoError(t, err)
				rawActualTxHash, err := json.Marshal(actualTrace.TxnHash)
				require.NoError(t, err)

				assert.JSONEq(t, string(rawExpectedTxHash), string(rawActualTxHash))

				// compare trace_root
				rawExpectedTxnTrace, err := json.Marshal(expectedTxnTrace)
				require.NoError(t, err)
				rawActualTxnTrace, err := json.Marshal(actualTrace.TraceRoot)
				require.NoError(t, err)

				compareTraceTxnsJSON(t, rawExpectedTxnTrace, rawActualTxnTrace)
			}
		})
	}
}

// compareTraceTxs compares two transaction traces.
// It is necessary because the order of the fields in the transaction trace is not deterministic.
// Hence, we need to compare the traces field by field.
func compareTraceTxs(t *testing.T, traceTx1, traceTx2 TxnTrace) {
	switch traceTx := traceTx1.(type) {
	case DeclareTxnTrace:
		assert.Equal(t, traceTx.ValidateInvocation, traceTx2.(DeclareTxnTrace).ValidateInvocation)
		assert.Equal(t, traceTx.FeeTransferInvocation, traceTx2.(DeclareTxnTrace).FeeTransferInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(DeclareTxnTrace).StateDiff)
		assert.Equal(t, traceTx.Type, traceTx2.(DeclareTxnTrace).Type)
		assert.Equal(t, traceTx.ExecutionResources, traceTx2.(DeclareTxnTrace).ExecutionResources)
	case DeployAccountTxnTrace:
		assert.Equal(t, traceTx.ValidateInvocation, traceTx2.(DeployAccountTxnTrace).ValidateInvocation)
		assert.Equal(t, traceTx.ConstructorInvocation, traceTx2.(DeployAccountTxnTrace).ConstructorInvocation)
		assert.Equal(t, traceTx.FeeTransferInvocation, traceTx2.(DeployAccountTxnTrace).FeeTransferInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(DeployAccountTxnTrace).StateDiff)
		assert.Equal(t, traceTx.Type, traceTx2.(DeployAccountTxnTrace).Type)
		assert.Equal(t, traceTx.ExecutionResources, traceTx2.(DeployAccountTxnTrace).ExecutionResources)
	case InvokeTxnTrace:
		assert.Equal(t, traceTx.ValidateInvocation, traceTx2.(InvokeTxnTrace).ValidateInvocation)
		assert.Equal(t, traceTx.ExecuteInvocation, traceTx2.(InvokeTxnTrace).ExecuteInvocation)
		assert.Equal(t, traceTx.FeeTransferInvocation, traceTx2.(InvokeTxnTrace).FeeTransferInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(InvokeTxnTrace).StateDiff)
		assert.Equal(t, traceTx.Type, traceTx2.(InvokeTxnTrace).Type)
		assert.Equal(t, traceTx.ExecutionResources, traceTx2.(InvokeTxnTrace).ExecutionResources)
	case L1HandlerTxnTrace:
		assert.Equal(t, traceTx.FunctionInvocation, traceTx2.(L1HandlerTxnTrace).FunctionInvocation)
		compareStateDiffs(t, traceTx.StateDiff, traceTx2.(L1HandlerTxnTrace).StateDiff)
		assert.Equal(t, traceTx.Type, traceTx2.(L1HandlerTxnTrace).Type)
	default:
		require.Failf(t, "unknown trace", "type: %T", traceTx)
	}
}

// compareStateDiffs compares two StateDiff objects.
// It is necessary because the order of the 'storage_entries' fields in the StateDiff is not deterministic.
// Hence, we need to compare the StateDiff field by field.
func compareStateDiffs(t *testing.T, stateDiff1, stateDiff2 *StateDiff) {
	if stateDiff1 == nil {
		return
	}

	assert.ElementsMatch(
		t,
		stateDiff1.DeprecatedDeclaredClasses,
		stateDiff2.DeprecatedDeclaredClasses,
	)
	assert.ElementsMatch(t, stateDiff1.DeclaredClasses, stateDiff2.DeclaredClasses)
	assert.ElementsMatch(t, stateDiff1.DeployedContracts, stateDiff2.DeployedContracts)
	assert.ElementsMatch(t, stateDiff1.ReplacedClasses, stateDiff2.ReplacedClasses)
	assert.ElementsMatch(t, stateDiff1.Nonces, stateDiff2.Nonces)

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
			addressFelt := internalUtils.TestHexToFelt(t, address.(string))

			if *addressFelt != *diff1.Address {
				continue
			}

			err = remarshal(diffElem, &diff2)
			require.NoError(t, err)
		}
		assert.NotEmpty(t, diff2)

		assert.Equal(t, diff1.Address, diff2.Address)
		assert.ElementsMatch(t, diff1.StorageEntries, diff2.StorageEntries)
	}
}

// compareTraceTxnsJSON compares two Marshalled JSON transaction traces to assert Marshalled JSON equality.
// It is necessary because the order of the fields in the 'storage_diffs' > 'storage_entries' is not deterministic.
func compareTraceTxnsJSON(t *testing.T, expectedResp, actualResp []byte) {
	t.Helper()

	expectedTxn, expectedStorageDiffs := splitJSONTraceTxn(t, expectedResp)
	actualTxn, actualStorageDiffs := splitJSONTraceTxn(t, actualResp)

	assert.JSONEq(t, string(expectedTxn), string(actualTxn))
	compareStorageDiffs(t, expectedStorageDiffs, actualStorageDiffs)
}

// splitJSONTraceTxn splits a transaction trace into a transaction without storage diffs and the storage diffs.
func splitJSONTraceTxn(t *testing.T, txn []byte) (txnWithoutStorageDiffs []byte, storageDiffs any) {
	var txnMap map[string]any
	require.NoError(t, json.Unmarshal(txn, &txnMap))

	if txnMap["state_diff"] == nil {
		return txn, nil
	}

	stateDiffMap, err := internalUtils.UnwrapJSON(txnMap, "state_diff")
	require.NoError(t, err)

	storageDiffs, ok := stateDiffMap["storage_diffs"]
	require.True(t, ok)

	delete(stateDiffMap, "storage_diffs")
	txnMap["state_diff"] = stateDiffMap
	txnWithoutStateDiff, err := json.Marshal(txnMap)
	require.NoError(t, err)

	return txnWithoutStateDiff, storageDiffs
}

// compareStorageDiffs compares the storage diffs of two Marshalled JSON transaction traces.
func compareStorageDiffs(t *testing.T, expectedStorageDiffs, actualStorageDiffs any) {
	t.Helper()

	if expectedStorageDiffs == nil {
		return
	}

	expectedStorageEntriesMap := getStorageEntries(t, expectedStorageDiffs)
	actualStorageEntriesMap := getStorageEntries(t, actualStorageDiffs)

	for address, expectedStorageEntries := range expectedStorageEntriesMap {
		actualStorageEntries, ok := actualStorageEntriesMap[address]
		require.True(t, ok)
		assert.ElementsMatch(t, expectedStorageEntries, actualStorageEntries)
	}
}

// getStorageEntries returns a map of storage entries for a given storage diffs array, classified by address.
func getStorageEntries(t *testing.T, storageDiffs any) (storageEntriesMap map[string][]any) {
	t.Helper()

	// e.g:
	// 	"storage_diffs": [
	// 		{
	//	   "address": "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d",
	//	   "storage_entries": [
	//	     {
	//	       "key": "0x5496768776e3db30053404f18067d81a6e06f5a2b0de326e21298fd9d569a9a",
	//	       "value": "0x1da15854fcce98a0660ba"
	//	     },
	//			...
	//	   },
	//	   ...
	anyArr, ok := storageDiffs.([]any)
	require.True(t, ok)
	storageDiffsMapArr := make([]map[string]any, len(anyArr))

	for i, diff := range anyArr {
		diffMap, ok := diff.(map[string]any)
		require.True(t, ok)

		storageDiffsMapArr[i] = diffMap
	}

	// e.g:
	//	   "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d": [
	//	     {
	//	       "key": "0x5496768776e3db30053404f18067d81a6e06f5a2b0de326e21298fd9d569a9a",
	//	       "value": "0x1da15854fcce98a0660ba"
	//	     },
	//			...
	//	   ]
	storageEntriesMap = make(map[string][]any)

	for _, diff := range storageDiffsMapArr {
		address, ok := diff["address"]
		require.True(t, ok)

		storageEntriesMap[address.(string)] = diff["storage_entries"].([]any)
	}

	return storageEntriesMap
}
