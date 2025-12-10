package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestBlockNumber tests the BlockNumber function.
func TestBlockNumber(t *testing.T) {
	tests.RunTestOn(t,
		tests.DevnetEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_blockNumber",
			).
			DoAndReturn(
				func(_, result, _ any, _ ...any) error {
					rawResp := result.(*json.RawMessage)
					rawBlockNumber := json.RawMessage("1234")
					*rawResp = rawBlockNumber

					return nil
				},
			).
			Times(1)
	}

	blockNumber, err := provider.BlockNumber(t.Context())
	require.NoError(t, err)

	rawExpectedResp := testConfig.Spy.LastResponse()
	rawActualResp, err := json.Marshal(blockNumber)
	require.NoError(t, err)
	assert.JSONEq(t, string(rawExpectedResp), string(rawActualResp))
}

// TestBlockHashAndNumber tests the BlockHashAndNumber function.
func TestBlockHashAndNumber(t *testing.T) {
	tests.RunTestOn(t,
		tests.DevnetEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_blockHashAndNumber",
			).
			DoAndReturn(
				func(_, result, _ any, _ ...any) error {
					rawResp := result.(*json.RawMessage)
					rawBlockHashAndNumber := json.RawMessage(
						`{
							"block_hash": "0x7fcc97a2e4e4a328582326254baca628cad2a82b17a711e7a8e5c9edd8022e6",
							"block_number": 3640605
						}`,
					)
					*rawResp = rawBlockHashAndNumber

					return nil
				},
			).
			Times(1)
	}

	blockHashAndNumber, err := provider.BlockHashAndNumber(t.Context())
	require.NoError(t, err, "BlockHashAndNumber should not return an error")

	rawExpectedResp := testConfig.Spy.LastResponse()
	rawActualResp, err := json.Marshal(blockHashAndNumber)
	require.NoError(t, err)
	assert.JSONEq(t, string(rawExpectedResp), string(rawActualResp))
}

// TestBlockWithTxHashes tests the BlockWithTxHashes function.
func TestBlockWithTxHashes(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	type testSetType struct {
		BlockID     BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	if tests.TEST_ENV != tests.MockEnv {
		// add the common block IDs to the test set of network tests
		blockIDs := GetCommonBlockIDs(t, provider)
		for _, blockID := range blockIDs {
			testSet = append(testSet, testSetType{
				BlockID: blockID,
			})
		}
	}

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				blockSepolia3100000 := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithHashes/sepoliaBlockWithHashes3100000.json", "result",
				)

				blockSepoliaPreConfirmed := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithHashes/sepoliaPreConfirmed.json", "result",
				)

				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getBlockWithTxHashes",
						test.BlockID,
					).
					DoAndReturn(
						func(_, result, _ any, args ...any) error {
							rawResp := result.(*json.RawMessage)
							blockID := args[0].(BlockID)

							switch blockID.Tag {
							case BlockTagPreConfirmed:
								*rawResp = blockSepoliaPreConfirmed
							case BlockTagLatest:
								*rawResp = blockSepolia3100000
							}

							if blockID.Number != nil && *blockID.Number == 99999999999999999 {
								return RPCError{
									Code:    24,
									Message: "Block not found",
								}
							}

							return nil
						},
					).
					Times(1)
			}

			result, err := provider.BlockWithTxHashes(t.Context(), test.BlockID)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)
			rawExpectedBlock := testConfig.Spy.LastResponse()

			switch block := result.(type) {
			case *BlockTxHashes:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			case *PreConfirmedBlockTxHashes:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			default:
				t.Fatalf("unexpected block type, found: %T\n", block)
			}
		})
	}
}

// TestBlockWithTxs tests the BlockWithTxs function.
func TestBlockWithTxs(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	type testSetType struct {
		BlockID     BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	if tests.TEST_ENV != tests.MockEnv {
		// add the common block IDs to the test set of network tests
		blockIDs := GetCommonBlockIDs(t, provider)
		for _, blockID := range blockIDs {
			testSet = append(testSet, testSetType{
				BlockID: blockID,
			})
		}
	}

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				blockSepolia3100000 := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithTxns/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithTxns/sepoliaPreConfirmed.json", "result",
				)

				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getBlockWithTxs",
						test.BlockID,
					).
					DoAndReturn(
						func(_, result, _ any, args ...any) error {
							rawResp := result.(*json.RawMessage)
							blockID := args[0].(BlockID)

							switch blockID.Tag {
							case BlockTagPreConfirmed:
								*rawResp = blockSepoliaPreConfirmed
							case BlockTagLatest:
								*rawResp = blockSepolia3100000
							}

							if blockID.Number != nil && *blockID.Number == 99999999999999999 {
								return RPCError{
									Code:    24,
									Message: "Block not found",
								}
							}

							return nil
						},
					).
					Times(1)
			}

			blockWithTxsInterface, err := provider.BlockWithTxs(
				t.Context(),
				test.BlockID,
			)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedBlock := testConfig.Spy.LastResponse()

			switch block := blockWithTxsInterface.(type) {
			case *PreConfirmedBlock:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			case *Block:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			}
		})
	}
}

// TestBlockTransactionCount tests the BlockTransactionCount function.
func TestBlockTransactionCount(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	type testSetType struct {
		BlockID     BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	if tests.TEST_ENV != tests.MockEnv {
		// add the common block IDs to the test set of network tests
		blockIDs := GetCommonBlockIDs(t, provider)
		for _, blockID := range blockIDs {
			testSet = append(testSet, testSetType{
				BlockID: blockID,
			})
		}
	}

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getBlockTransactionCount",
						test.BlockID,
					).
					DoAndReturn(
						func(_, result, _ any, args ...any) error {
							rawResp := result.(*json.RawMessage)
							blockID := args[0].(BlockID)

							if blockID.Tag == BlockTagLatest {
								*rawResp = json.RawMessage("100")
							}

							if blockID.Number != nil && *blockID.Number == 99999999999999999 {
								return RPCError{
									Code:    24,
									Message: "Block not found",
								}
							}

							return nil
						},
					).
					Times(1)
			}

			count, err := testConfig.Provider.BlockTransactionCount(
				t.Context(),
				test.BlockID,
			)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedCount := testConfig.Spy.LastResponse()

			rawCount, err := json.Marshal(count)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedCount), string(rawCount))
		})
	}
}

// TestCaptureUnsupportedBlockTxn tests the functionality of capturing unsupported block transactions.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestCaptureUnsupportedBlockTxn(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		StartBlock uint64
		EndBlock   uint64
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {
			{
				StartBlock: 52959,
				EndBlock:   52960,
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		for i := test.StartBlock; i < test.EndBlock; i++ {
			blockWithTxsInterface, err := testConfig.Provider.BlockWithTxs(
				context.Background(),
				WithBlockNumber(i),
			)
			require.NoError(t, err)
			blockWithTxs, ok := blockWithTxsInterface.(*Block)
			require.True(t, ok, "expecting *rpc.Block, instead %T", blockWithTxsInterface)

			for k, v := range blockWithTxs.Transactions {
				_, okv0 := v.Transaction.(InvokeTxnV0)
				_, okv1 := v.Transaction.(InvokeTxnV1)
				_, okv3 := v.Transaction.(InvokeTxnV3)
				_, okl1 := v.Transaction.(L1HandlerTxn)
				_, okdec0 := v.Transaction.(DeclareTxnV0)
				_, okdec1 := v.Transaction.(DeclareTxnV1)
				_, okdec2 := v.Transaction.(DeclareTxnV2)
				_, okdec3 := v.Transaction.(DeclareTxnV3)
				_, okdep := v.Transaction.(DeployTxn)
				_, okdepac := v.Transaction.(DeployAccountTxnV1)
				if !okv0 && !okv1 && !okv3 && !okl1 && !okdec0 && !okdec1 && !okdec2 && !okdec3 &&
					!okdep &&
					!okdepac {
					t.Fatalf("New Type Detected %T at Block(%d)/Txn(%d)", v, i, k)
				}
			}
		}
	}
}

// TestStateUpdate is a test function for the StateUpdate method.
func TestStateUpdate(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		BlockID                       BlockID
		ExpectedStateUpdateOutputPath string
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:                       WithBlockNumber(30000),
				ExpectedStateUpdateOutputPath: "testData/stateUpdate/sepolia_30000.json",
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID:                       WithBlockNumber(30000),
				ExpectedStateUpdateOutputPath: "testData/stateUpdate/sepolia_30000.json",
			},
			{
				BlockID:                       WithBlockNumber(1060000),
				ExpectedStateUpdateOutputPath: "testData/stateUpdate/sepolia_1_060_000.json",
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID:                       WithBlockNumber(30000),
				ExpectedStateUpdateOutputPath: "testData/stateUpdate/integration_30000.json",
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		t.Run(fmt.Sprintf("BlockID: %v", test.BlockID), func(t *testing.T) {
			stateUpdate, err := testConfig.Provider.StateUpdate(context.Background(), test.BlockID)
			require.NoError(t, err, "Unable to fetch the given block.")

			if test.ExpectedStateUpdateOutputPath != "" {
				rawExpectedCasmClass, err := os.ReadFile(test.ExpectedStateUpdateOutputPath)
				require.NoError(t, err)

				rawActualCasmClass, err := json.Marshal(stateUpdate)
				require.NoError(t, err)

				assertStateUpdateJSONEquality(t, "result", rawExpectedCasmClass, rawActualCasmClass)

				return
			}
			assert.NotEmpty(t, stateUpdate)
		})
	}
}

func assertStateUpdateJSONEquality(t *testing.T, subfield string, expectedResult, result []byte) {
	// unmarshal to map[string]any
	var expectedResultMap, resultMap map[string]any
	require.NoError(t, json.Unmarshal(expectedResult, &expectedResultMap))
	require.NoError(t, json.Unmarshal(result, &resultMap))

	if subfield != "" {
		var ok bool
		expectedResultMap, ok = expectedResultMap[subfield].(map[string]any)
		require.True(t, ok, "expected result map should have a subfield %s", subfield)
	}

	assert.Equal(t, expectedResultMap["block_hash"], resultMap["block_hash"])
	assert.Equal(t, expectedResultMap["new_root"], resultMap["new_root"])
	assert.Equal(t, expectedResultMap["old_root"], resultMap["old_root"])

	// ********** compare 'state_diff' **********
	expectedStateDiff, ok := expectedResultMap["state_diff"].(map[string]any)
	require.True(t, ok)
	resultStateDiff, ok := resultMap["state_diff"].(map[string]any)
	require.True(t, ok)

	// compare 'state_diff.storage_diffs'
	expectedStorageDiffs, ok := expectedStateDiff["storage_diffs"].([]any)
	require.True(t, ok)
	resultStorageDiffs, ok := resultStateDiff["storage_diffs"].([]any)
	require.True(t, ok)

	expectedStorageDiffsMap := make(map[string]any)
	resultStorageDiffsMap := make(map[string]any)

	for i, expectedStorageDiff := range expectedStorageDiffs {
		expectedStorageDiffMap, ok := expectedStorageDiff.(map[string]any)
		require.True(t, ok)
		address, ok := expectedStorageDiffMap["address"].(string)
		require.True(t, ok)
		storageEntries, ok := expectedStorageDiffMap["storage_entries"].([]any)
		require.True(t, ok)

		expectedStorageDiffsMap[address] = storageEntries

		resultStorageDiffMap, ok := resultStorageDiffs[i].(map[string]any)
		require.True(t, ok)
		address2, ok := resultStorageDiffMap["address"].(string)
		require.True(t, ok)
		storageEntries2, ok := resultStorageDiffMap["storage_entries"].([]any)
		require.True(t, ok)

		resultStorageDiffsMap[address2] = storageEntries2
	}

	assert.Len(t, resultStorageDiffsMap, len(expectedStorageDiffsMap))
	for address, expectedStorageEntries := range expectedStorageDiffsMap {
		resultStorageEntries, ok := resultStorageDiffsMap[address]
		require.True(t, ok, "address %s not found in resultStorageDiffsMap", address)
		assert.ElementsMatch(t, expectedStorageEntries, resultStorageEntries)
	}

	// other state diff fields
	assert.ElementsMatch(t, expectedResultMap["nonces"], resultMap["nonces"])
	assert.ElementsMatch(
		t,
		expectedResultMap["deployed_contracts"],
		resultMap["deployed_contracts"],
	)
	assert.ElementsMatch(
		t,
		expectedResultMap["deprecated_declared_classes"],
		resultMap["deprecated_declared_classes"],
	)
	assert.ElementsMatch(t, expectedResultMap["declared_classes"], resultMap["declared_classes"])
	assert.ElementsMatch(t, expectedResultMap["replaced_classes"], resultMap["replaced_classes"])
}
