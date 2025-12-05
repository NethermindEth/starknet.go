package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
					rawBlockNumber, err := json.Marshal(uint64(1234))
					if err != nil {
						return err
					}
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

// TestBlockHashAndNumber is a test function that tests the BlockHashAndNumber function and check if there is no errors.
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestBlockHashAndNumber(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	blockHashAndNumber, err := testConfig.Provider.BlockHashAndNumber(context.Background())
	require.NoError(t, err, "BlockHashAndNumber should not return an error")
	require.True(
		t,
		strings.HasPrefix(blockHashAndNumber.Hash.String(), "0x"),
		"current block hash should return a string starting with 0x",
	)

	if tests.TEST_ENV == tests.MockEnv {
		require.Equal(
			t,
			&BlockHashAndNumberOutput{Number: 1234, Hash: internalUtils.DeadBeef},
			blockHashAndNumber,
		)
	}
}

// TestBlockWithTxHashes tests the functionality of the BlockWithTxHashes function.
//
// The function takes a testing.T object as a parameter and initialises a testConfig object.
// It defines a testSetType struct that contains several fields including BlockID, ExpectedError, ExpectedBlockWithTxHashes, and ExpectedPre_confirmedBlockWithTxHashes.
// The function then initialises a blockSepolia64159 variable of type BlockTxHashes with a predefined set of values.
// It also initialises a txHashes variable of type []felt.Felt and a blockHash variable of type felt.Felt.
//
// The function defines a testSet map that has three keys: "mock", "testnet", and "mainnet".
// Each key corresponds to a slice of testSetType objects.
// The "mock" key has two testSetType objects with different field values.
// The "testnet" key has three testSetType objects with different field values.
// The "mainnet" key does not have any testSetType objects.
//
// The function then iterates over the testSet map and performs the following steps for each testSetType object:
//   - It creates a new Spy object and assigns it to the testConfig.provider.c field.
//   - It calls the BlockWithTxHashes function with the provided BlockID and stores the result in the result variable.
//   - It checks if the returned error matches the expected error. If not, it calls the Fatal function of the testing.T object with an error message.
//   - It checks the type of the result variable and performs specific assertions based on the type.
//   - If the result is of type *BlockTxHashes, it checks various fields of the BlockTxHashes object against the expected values.
//   - If the result is of type *Pre_confirmedBlockTxHashes, it checks various fields of the Pre_confirmedBlockTxHashes object against the expected values.
//   - If the result is of any other type, it calls the Fatal function of the testing.T object with an error message.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestBlockWithTxHashes(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider
	spy := tests.NewJSONRPCSpy(provider.c)
	provider.c = spy

	type testSetType struct {
		BlockID     BlockID
		ExpectedErr error
	}

	// TODO: use these blocks for mock tests
	// blockSepolia3100000 := *internalUtils.TestUnmarshalJSONFileToType[BlockTxHashes](t, "./testData/blockWithHashes/sepoliaBlockWithHashes3100000.json", "result")
	// blockIntegration1300000 := *internalUtils.TestUnmarshalJSONFileToType[BlockTxHashes](t, "./testData/blockWithHashes/integration1_300_000.json", "result")

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:     BlockID{Tag: BlockTagPreConfirmed},
				ExpectedErr: nil,
			},
			{
				BlockID: BlockID{Hash: internalUtils.DeadBeef},
			},
			{
				BlockID: BlockID{Tag: BlockTagL1Accepted},
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     WithBlockTag(BlockTagLatest),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagPreConfirmed),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagL1Accepted),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockHash(internalUtils.TestHexToFelt(t, "0x1640b846e71502526539c32c8420cd7cb0f28d83ece2e6e71aeaf7c97960bb2")),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockNumber(3100000),
				ExpectedErr: nil,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     WithBlockTag(BlockTagLatest),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagPreConfirmed),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagL1Accepted),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockHash(internalUtils.TestHexToFelt(t, "0x503e44c7d47a2e17022c52092e7dadd338b79df84f844b9f26dbdd1598a23e")),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockNumber(1300000),
				ExpectedErr: nil,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			result, err := provider.BlockWithTxHashes(context.Background(), test.BlockID)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)
			rawExpectedBlock := spy.LastResponse()

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
//
// The function tests the BlockWithTxs function by setting up a test configuration and a test set type.
// It then initialises Block type variables and invokes the BlockWithTxs function with different test scenarios.
// The function checks if the BlockWithTxs function returns the correct block data.
// It also verifies the block hash, the number of transactions in the block, and the details of a specific transaction.
//
// Parameters:
//   - t: The t testing object
//
// Returns:
//
//	none
func TestBlockWithTxs(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider
	spy := tests.NewJSONRPCSpy(provider.c)
	provider.c = spy

	type testSetType struct {
		BlockID BlockID
	}

	// TODO: use these blocks for mock tests
	// fullBlockSepolia65083 := *internalUtils.TestUnmarshalJSONFileToType[Block](t, "./testData/block/sepoliaBlockTxs65083.json", "result")
	// fullBlockSepolia122476 := *internalUtils.TestUnmarshalJSONFileToType[Block](t, "./testData/block/sepoliaBlockTxs122476.json", "result")
	// fullBlockIntegration1300000 := *internalUtils.TestUnmarshalJSONFileToType[Block](t, "./testData/block/integration1_300_000.json", "result")

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
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
				BlockID: WithBlockNumber(65083),
			},
			{
				BlockID: WithBlockHash(internalUtils.TestHexToFelt(t, "0x549770b5b74df90276277ff7a11af881c998dffa452f4156f14446db6005174")),
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
				BlockID: WithBlockNumber(65083),
			},
			{
				BlockID: WithBlockNumber(3100000),
			},
			{
				BlockID: WithBlockHash(internalUtils.TestHexToFelt(t, "0x56a71e0443d2fbfaa91b1000e830b516ca0d4a424abb9c970d23801957dbfa3")),
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
				BlockID: WithBlockNumber(1300000),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			blockWithTxsInterface, err := provider.BlockWithTxs(
				context.Background(),
				test.BlockID,
			)
			require.NoError(t, err, "Unable to fetch the given block.")
			rawExpectedBlock := spy.LastResponse()

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

// TestBlockTransactionCount tests the function that calculates the number of transactions in a block.
func TestBlockTransactionCount(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		BlockID       BlockID
		ExpectedCount int64
		ExpectedError error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:       WithBlockNumber(300000),
				ExpectedCount: 10,
			},
			{
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedCount: 10,
			},
			{
				BlockID:       WithBlockTag(BlockTagPreConfirmed),
				ExpectedCount: 10,
			},
			{
				BlockID:       WithBlockTag(BlockTagL1Accepted),
				ExpectedCount: 10,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:       WithBlockNumber(30000),
				ExpectedCount: 4,
			},
			{
				BlockID:       WithBlockNumber(52959),
				ExpectedCount: 58,
			},
			{
				BlockID:       WithBlockTag(BlockTagPreConfirmed),
				ExpectedCount: -1,
			},
			{
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedCount: -1,
			},
			{
				BlockID:       WithBlockTag(BlockTagL1Accepted),
				ExpectedCount: -1,
			},
			{
				BlockID:       WithBlockNumber(7338746823462834783),
				ExpectedError: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:       WithBlockNumber(30000),
				ExpectedCount: 4,
			},
			{
				BlockID:       WithBlockNumber(529590),
				ExpectedCount: 6,
			},
			{
				BlockID:       WithBlockNumber(7338746823462834783),
				ExpectedError: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		t.Run(
			fmt.Sprintf("Count: %v, BlockID: %v", test.ExpectedCount, test.BlockID),
			func(t *testing.T) {
				count, err := testConfig.Provider.BlockTransactionCount(
					context.Background(),
					test.BlockID,
				)
				if test.ExpectedError != nil {
					require.EqualError(t, test.ExpectedError, err.Error())

					return
				}
				require.NoError(t, err)

				if test.ExpectedCount == -1 {
					// since 0 is the default value of an int64 var, let's set the expected count to -1 when we want to skip the count check
					return
				}

				assert.Equal(t, uint64(test.ExpectedCount), count)
			},
		)
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
