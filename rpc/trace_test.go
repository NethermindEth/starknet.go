package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestTraceTransaction tests the TraceTransaction function.
func TestTransactionTrace(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)
	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TransactionHash *felt.Felt
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TransactionHash: internalUtils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
			},
			{
				TransactionHash: internalUtils.DeadBeef,
				ExpectedError:   ErrHashNotFound,
			},
			{
				TransactionHash: &felt.Zero,
				ExpectedError: &RPCError{
					Code:    10,
					Message: "No trace available for transaction",
					Data: &TraceStatusErrData{
						Status: TraceStatusReceived,
					},
				},
			},
		},
		tests.TestnetEnv: {
			{ // with 5 out of 6 fields (without state diff)
				TransactionHash: internalUtils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
			},
			{ // with 6 out of 6 fields
				TransactionHash: internalUtils.TestHexToFelt(t, "0x49d98a0328fee1de19d43d950cbaeb973d080d0c74c652523371e034cc0bbb2"),
			},
			{
				TransactionHash: internalUtils.DeadBeef,
				ExpectedError:   ErrHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				TransactionHash: internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
			},
			{
				TransactionHash: internalUtils.DeadBeef,
				ExpectedError:   ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.TransactionHash.String(), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
					testConfig.MockClient.EXPECT().
						CallContextWithSliceArgs(
							t.Context(),
							gomock.Any(),
							"starknet_traceTransaction",
							test.TransactionHash,
						).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						transactionHash := args[0].(*felt.Felt)

						if transactionHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    29,
								Message: "Transaction hash not found",
							}
						}

						if *transactionHash == felt.Zero {
							return RPCError{
								Code:    10,
								Message: "No trace available for transaction",
								Data: &TraceStatusErrData{
									Status: TraceStatusReceived,
								},
							}
						}

						*rawResp = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
						t,
						"./testData/trace/sepoliaInvokeTrace_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json",
						"result",
					)

							return nil
						}).
						Times(1)
				}

			resp, err := testConfig.Provider.TraceTransaction(
				t.Context(),
				test.TransactionHash,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}

// TestSimulateTransaction tests the SimulateTransaction function.
func TestSimulateTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)
	testConfig := BeforeEach(t, false)

	type simulateTxnInput struct {
		BlockID         BlockID          `json:"block_id"`
		Txns            []BroadcastTxn   `json:"transactions"`
		SimulationFlags []SimulationFlag `json:"simulation_flags"`
	}
	input := *internalUtils.TestUnmarshalJSONFileToType[simulateTxnInput](
		t, "./testData/trace/sepoliaSimulateInvokeTx.json", "params")

	type testSetType struct {
		Description     string
		BlockID         BlockID
		Txns            []BroadcastTxn
		SimulationFlags []SimulationFlag
		ExpectedError   *RPCError
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "valid call, all flags",
				BlockID:         input.BlockID,
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{SkipValidate, SkipFeeCharge},
			},
			{
				Description:     "block not found",
				BlockID:         WithBlockHash(internalUtils.DeadBeef),
				Txns:            input.Txns,
				SimulationFlags: input.SimulationFlags,
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "exec error, pre confirmed",
				BlockID:         WithBlockTag(BlockTagPreConfirmed),
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{},
				ExpectedError:   ErrTxnExec, // due to invalid nonce
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "valid call, no flags",
				BlockID:         input.BlockID,
				Txns:            input.Txns,
				SimulationFlags: input.SimulationFlags,
			},
			{
				Description:     "valid call, skip fee charge",
				BlockID:         input.BlockID,
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{SkipFeeCharge},
			},
			{
				Description:     "valid call, latest + skip validate + skip fee charge",
				BlockID:         WithBlockTag(BlockTagLatest),
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{SkipValidate, SkipFeeCharge},
			},
			{
				Description:     "valid call, l1 accepted + skip validate",
				BlockID:         WithBlockTag(BlockTagL1Accepted),
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{SkipValidate},
			},
			{
				Description:     "exec error, pre confirmed",
				BlockID:         WithBlockTag(BlockTagPreConfirmed),
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{},
				ExpectedError:   ErrTxnExec, // due to invalid nonce
			},
			{
				Description:     "block not found",
				BlockID:         WithBlockHash(internalUtils.DeadBeef),
				Txns:            input.Txns,
				SimulationFlags: input.SimulationFlags,
				ExpectedError:   ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_simulateTransactions",
						test.BlockID,
						test.Txns,
						test.SimulationFlags,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if blockID.Tag == BlockTagPreConfirmed {
							return RPCError{
								Code:    41,
								Message: "Transaction execution error",
								Data:    &TransactionExecErrData{},
							}
						}
						*rawResp = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/trace/sepoliaSimulateInvokeTxResp.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.SimulateTransactions(
				t.Context(),
				test.BlockID,
				test.Txns,
				test.SimulationFlags)

			if test.ExpectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)

				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

				return
			}
				require.NoError(t, err)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawResp, err := json.Marshal(resp)
				require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
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
