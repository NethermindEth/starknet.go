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
)

// TestBlockNumber is a test function to check the behaviour of the BlockNumber function and check if there is no errors.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestBlockNumber(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	blockNumber, err := testConfig.Provider.BlockNumber(context.Background())
	require.NoError(t, err, "BlockNumber should not return an error")
	if tests.TEST_ENV == tests.MockEnv {
		require.Equal(t, uint64(1234), blockNumber)
	}
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
			&BlockHashAndNumberOutput{Number: 1234, Hash: internalUtils.RANDOM_FELT},
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

	type testSetType struct {
		BlockID                                BlockID
		ExpectedErr                            error
		ExpectedBlockWithTxHashes              *BlockTxHashes
		ExpectedPre_confirmedBlockWithTxHashes *Pre_confirmedBlockTxHashes
	}

	blockSepolia64159 := *internalUtils.TestUnmarshalJSONFileToType[BlockTxHashes](t, "./testData/blockWithHashes/sepoliaBlockWithHashes64159.json", "result")
	blockIntegration1300000 := *internalUtils.TestUnmarshalJSONFileToType[BlockTxHashes](t, "./testData/blockWithHashes/integration1_300_000.json", "result")

	txHashesMock := internalUtils.TestHexArrToFelt(t, []string{
		"0x5754961d70d6f39d0e2c71a1a4ff5df0a26b1ceda4881ca82898994379e1e73",
		"0x692381bba0e8505a8e0b92d0f046c8272de9e65f050850df678a0c10d8781d",
	})
	blockMock := BlockTxHashes{
		BlockHeader: BlockHeader{
			Hash:             internalUtils.RANDOM_FELT,
			ParentHash:       internalUtils.RANDOM_FELT,
			Timestamp:        124,
			SequencerAddress: internalUtils.RANDOM_FELT,
		},
		Status:       BlockStatus_AcceptedOnL1,
		Transactions: txHashesMock,
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:     BlockID{Tag: BlockTagPre_confirmed},
				ExpectedErr: nil,
				ExpectedPre_confirmedBlockWithTxHashes: &Pre_confirmedBlockTxHashes{
					Pre_confirmedBlockHeader{
						Number:           1234,
						Timestamp:        123,
						SequencerAddress: internalUtils.RANDOM_FELT,
					},
					txHashesMock,
				},
			},
			{
				BlockID:                   BlockID{Hash: internalUtils.RANDOM_FELT},
				ExpectedBlockWithTxHashes: &blockMock,
			},
			{
				BlockID:                   BlockID{Tag: BlockTagL1Accepted},
				ExpectedBlockWithTxHashes: &blockMock,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     WithBlockTag(BlockTagLatest),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagPre_confirmed),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagL1Accepted),
				ExpectedErr: nil,
			},
			{
				BlockID:                   WithBlockHash(internalUtils.TestHexToFelt(t, "0x6df565874b2ea6a02d346a23f9efb0b26abbf5708b51bb12587f88a49052964")),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockSepolia64159,
			},
			{
				BlockID:                   WithBlockNumber(64159),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockSepolia64159,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     WithBlockTag(BlockTagLatest),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagPre_confirmed),
				ExpectedErr: nil,
			},
			{
				BlockID:     WithBlockTag(BlockTagL1Accepted),
				ExpectedErr: nil,
			},
			{
				BlockID:                   WithBlockHash(internalUtils.TestHexToFelt(t, "0x503e44c7d47a2e17022c52092e7dadd338b79df84f844b9f26dbdd1598a23e")),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockIntegration1300000,
			},
			{
				BlockID:                   WithBlockNumber(1300000),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockIntegration1300000,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			result, err := testConfig.Provider.BlockWithTxHashes(context.Background(), test.BlockID)
			require.Equal(t, test.ExpectedErr, err, "Error in BlockWithTxHashes")
			switch resultType := result.(type) {
			case *BlockTxHashes:
				block, ok := result.(*BlockTxHashes)
				require.Truef(t, ok, "should return *BlockTxHashes, instead: %T\n", result)

				if test.ExpectedErr != nil {
					return
				}

				assert.Truef(t, strings.HasPrefix(block.Hash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.Hash)

				if test.ExpectedBlockWithTxHashes != nil {
					assert.Exactly(t, test.ExpectedBlockWithTxHashes, block)
				}
			case *Pre_confirmedBlockTxHashes:
				pBlock, ok := result.(*Pre_confirmedBlockTxHashes)
				require.Truef(t, ok, "should return *Pre_confirmedBlockTxHashes, instead: %T\n", result)

				if test.ExpectedPre_confirmedBlockWithTxHashes == nil {
					validatePre_confirmedBlockHeader(t, &pBlock.Pre_confirmedBlockHeader)
				} else {
					assert.Exactly(t, test.ExpectedPre_confirmedBlockWithTxHashes, pBlock)
				}
			default:
				t.Fatalf("unexpected block type, found: %T\n", resultType)
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

	type testSetType struct {
		BlockID                    BlockID
		ExpectedBlock              *Block
		ExpectedPre_confirmedBlock *Pre_confirmedBlock
		InvokeV0Index              int // TODO: implement mainnet testcases as Sepolia doesn't contains V0 transactions
		InvokeV1Index              int
		InvokeV3Index              int
		DeclareV0Index             int // TODO: implement mainnet testcases as Sepolia doesn't contains V0 transactions
		DeclareV1Index             int
		DeclareV2Index             int
		DeclareV3Index             int // TODO: implement testcase
		DeployAccountV1Index       int
		DeployAccountV3Index       int // TODO: implement testcase
		L1HandlerV0Index           int
		DeployV0Index              int // TODO: implement testcase
	}

	fullBlockSepolia65083 := *internalUtils.TestUnmarshalJSONFileToType[Block](t, "./testData/block/sepoliaBlockTxs65083.json", "result")
	fullBlockSepolia122476 := *internalUtils.TestUnmarshalJSONFileToType[Block](t, "./testData/block/sepoliaBlockTxs122476.json", "result")
	fullBlockIntegration1300000 := *internalUtils.TestUnmarshalJSONFileToType[Block](t, "./testData/block/integration1_300_000.json", "result")

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPre_confirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
				ExpectedPre_confirmedBlock: &Pre_confirmedBlock{
					Pre_confirmedBlockHeader{
						Number:           1234,
						Timestamp:        1234,
						SequencerAddress: internalUtils.RANDOM_FELT,
						L1GasPrice: ResourcePrice{
							PriceInFRI: internalUtils.RANDOM_FELT,
							PriceInWei: internalUtils.RANDOM_FELT,
						},
						L2GasPrice: ResourcePrice{
							PriceInFRI: internalUtils.RANDOM_FELT,
							PriceInWei: internalUtils.RANDOM_FELT,
						},
						L1DataGasPrice: ResourcePrice{
							PriceInFRI: internalUtils.RANDOM_FELT,
							PriceInWei: internalUtils.RANDOM_FELT,
						},
						L1DAMode:        L1DAModeBlob,
						StarknetVersion: "0.14.0",
					},
					[]BlockTransaction{},
				},
			},
			{
				BlockID:       WithBlockNumber(65083),
				ExpectedBlock: &fullBlockSepolia65083,
				InvokeV1Index: 1,
			},
			{
				BlockID:       WithBlockHash(internalUtils.TestHexToFelt(t, "0x549770b5b74df90276277ff7a11af881c998dffa452f4156f14446db6005174")),
				ExpectedBlock: &fullBlockSepolia65083,
				InvokeV1Index: 1,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPre_confirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID:              WithBlockNumber(65083),
				ExpectedBlock:        &fullBlockSepolia65083,
				InvokeV1Index:        1,
				InvokeV3Index:        3,
				DeclareV1Index:       26,
				DeclareV2Index:       25,
				DeployAccountV1Index: 42,
			},
			{
				BlockID:          WithBlockHash(internalUtils.TestHexToFelt(t, "0x56a71e0443d2fbfaa91b1000e830b516ca0d4a424abb9c970d23801957dbfa3")),
				ExpectedBlock:    &fullBlockSepolia122476,
				L1HandlerV0Index: 4,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPre_confirmed),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID:       WithBlockNumber(1300000),
				ExpectedBlock: &fullBlockIntegration1300000,
				InvokeV3Index: 1,
			},
		},
	}[tests.TEST_ENV]

	// TODO: refactor test to check the marshal result against the expected json file
	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			blockWithTxsInterface, err := testConfig.Provider.BlockWithTxs(
				context.Background(),
				test.BlockID,
			)
			require.NoError(t, err, "Unable to fetch the given block.")

			switch block := blockWithTxsInterface.(type) {
			case *Pre_confirmedBlock:
				if test.ExpectedPre_confirmedBlock == nil {
					validatePre_confirmedBlockHeader(t, &block.Pre_confirmedBlockHeader)
				} else {
					assert.Exactly(t, test.ExpectedPre_confirmedBlock, block)
				}
			case *Block:
				if test.ExpectedBlock == nil {
					assert.Equal(t, block.Hash.String()[:2], "0x", "Block Hash should start with \"0x\".")
				} else {
					assert.Exactly(t, test.ExpectedBlock, block)

					// validates an BlockInvokeV1 transaction
					if test.InvokeV1Index > 0 {
						invokeV1Expected, ok := test.ExpectedBlock.Transactions[test.InvokeV1Index].Transaction.(InvokeTxnV1)
						assert.True(t, ok, "Expected invoke v1 transaction.")
						invokeV1Block, ok := block.Transactions[test.InvokeV1Index].Transaction.(InvokeTxnV1)
						assert.True(t, ok, "Expected invoke v1 transaction.")

						assert.Exactly(t, invokeV1Expected, invokeV1Block)
					}

					// validates an BlockInvokeV3 transaction
					if test.InvokeV3Index > 0 {
						invokeV3Expected, ok := test.ExpectedBlock.Transactions[test.InvokeV3Index].Transaction.(InvokeTxnV3)
						assert.True(t, ok, "Expected invoke v3 transaction.")
						invokeV3Block, ok := block.Transactions[test.InvokeV3Index].Transaction.(InvokeTxnV3)
						assert.True(t, ok, "Expected invoke v3 transaction.")

						assert.Exactly(t, invokeV3Expected, invokeV3Block)
					}

					// validates an BlockDeclareV1 transaction
					if test.DeclareV1Index > 0 {
						declareV1Expected, ok := test.ExpectedBlock.Transactions[test.DeclareV1Index].Transaction.(DeclareTxnV1)
						assert.True(t, ok, "Expected declare v1 transaction.")
						declareV1Block, ok := block.Transactions[test.DeclareV1Index].Transaction.(DeclareTxnV1)
						assert.True(t, ok, "Expected declare v1 transaction.")

						assert.Exactly(t, declareV1Expected, declareV1Block)
					}

					// validates an BlockDeclareV2 transaction
					if test.DeclareV2Index > 0 {
						declareV2Expected, ok := test.ExpectedBlock.Transactions[test.DeclareV2Index].Transaction.(DeclareTxnV2)
						assert.True(t, ok, "Expected declare v2 transaction.")
						declareV2Block, ok := block.Transactions[test.DeclareV2Index].Transaction.(DeclareTxnV2)
						assert.True(t, ok, "Expected declare v2 transaction.")

						assert.Exactly(t, declareV2Expected, declareV2Block)
					}

					// validates an BlockDeployAccountV1 transaction
					if test.DeployAccountV1Index > 0 {
						deployAccountV1Expected, ok := test.ExpectedBlock.Transactions[test.DeployAccountV1Index].Transaction.(DeployAccountTxnV1)
						assert.True(t, ok, "Expected deploy account v1 transaction.")
						deployAccountV1Block, ok := block.Transactions[test.DeployAccountV1Index].Transaction.(DeployAccountTxnV1)
						assert.True(t, ok, "Expected deploy account v1 transaction.")

						assert.Exactly(t, deployAccountV1Expected, deployAccountV1Block)
					}

					// validates an BlockL1HandlerV0 transaction
					if test.L1HandlerV0Index > 0 {
						l1HandlerV0Expected, ok := test.ExpectedBlock.Transactions[test.L1HandlerV0Index].Transaction.(L1HandlerTxn)
						assert.True(t, ok, "Expected L1 handler transaction.")
						l1HandlerV0Block, ok := block.Transactions[test.L1HandlerV0Index].Transaction.(L1HandlerTxn)
						assert.True(t, ok, "Expected L1 handler transaction.")

						assert.Exactly(t, l1HandlerV0Expected, l1HandlerV0Block)
					}
				}
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
				BlockID:       WithBlockTag(BlockTagPre_confirmed),
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
				BlockID:       WithBlockTag(BlockTagPre_confirmed),
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
				BlockID: WithBlockTag(BlockTagPre_confirmed),
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
				BlockID: WithBlockTag(BlockTagPre_confirmed),
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

func validatePre_confirmedBlockHeader(t *testing.T, pBlock *Pre_confirmedBlockHeader) {
	assert.NotZero(t, pBlock.Number)
	assert.NotZero(t, pBlock.Timestamp)
	assert.NotZero(t, pBlock.SequencerAddress)
	assert.NotZero(t, pBlock.L1GasPrice)
	assert.NotZero(t, pBlock.L2GasPrice)
	assert.NotZero(t, pBlock.L1DataGasPrice)
	assert.NotZero(t, pBlock.StarknetVersion)
}
