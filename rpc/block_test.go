package rpc

import (
	"encoding/json"
	"fmt"
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

	rawExpectedResp := testConfig.RPCSpy.LastResponse()
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

	rawExpectedResp := testConfig.RPCSpy.LastResponse()
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
				blockSepolia3100000 := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithHashes/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithHashes/sepoliaPreConfirmed.json",
					"result",
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
			require.NotNil(t, result)
			rawExpectedBlock := testConfig.RPCSpy.LastResponse()

			switch {
			case result.Block != nil:
				rawBlock, err := json.Marshal(result.Block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
				assert.Nil(t, result.PreConfirmed)
			case result.PreConfirmed != nil:
				rawBlock, err := json.Marshal(result.PreConfirmed)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
				assert.Nil(t, result.Block)
			default:
				t.Fatalf("unexpected block result: both block fields are nil")
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
				blockSepolia3100000 := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithTxns/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithTxns/sepoliaPreConfirmed.json",
					"result",
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

			blockWithTxsOutput, err := provider.BlockWithTxs(
				t.Context(),
				test.BlockID,
			)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)
			require.NotNil(t, blockWithTxsOutput)
			rawExpectedBlock := testConfig.RPCSpy.LastResponse()

			switch {
			case blockWithTxsOutput.Block != nil:
				rawBlock, err := json.Marshal(blockWithTxsOutput.Block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
				assert.Nil(t, blockWithTxsOutput.PreConfirmed)
			case blockWithTxsOutput.PreConfirmed != nil:
				rawBlock, err := json.Marshal(blockWithTxsOutput.PreConfirmed)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
				assert.Nil(t, blockWithTxsOutput.Block)
			default:
				t.Fatalf("unexpected block result: both block fields are nil")
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

			rawExpectedCount := testConfig.RPCSpy.LastResponse()

			rawCount, err := json.Marshal(count)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedCount), string(rawCount))
		})
	}
}

// TestStateUpdate is a test function for the StateUpdate method.
func TestStateUpdate(t *testing.T) {
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
				BlockID: WithBlockTag(BlockTagPreConfirmed),
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
				blockSepolia3100000 := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/stateUpdate/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/stateUpdate/sepoliaPreConfirmed.json",
					"result",
				)

				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getStateUpdate",
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

			stateUpdate, err := provider.StateUpdate(t.Context(), test.BlockID)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedStateUpdate := testConfig.RPCSpy.LastResponse()

			rawStateUpdate, err := json.Marshal(stateUpdate)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedStateUpdate), string(rawStateUpdate))
		})
	}
}

// TestBlockWithReceipts tests the BlockWithReceipts function.
func TestBlockWithReceipts(t *testing.T) {
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
				BlockID:     WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID:     WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     WithBlockHash(internalUtils.DeadBeef),
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
		t.Run(string(blockID), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				blockSepolia3100000 := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithReceipts/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithReceipts/sepoliaPreConfirmed.json",
					"result",
				)

				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getBlockWithReceipts",
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

							if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
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
			result, err := provider.BlockWithReceipts(t.Context(), test.BlockID)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedBlock := testConfig.RPCSpy.LastResponse()

			switch block := result.(type) {
			case *BlockWithReceipts:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			case *PreConfirmedBlockWithReceipts:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			default:
				t.Fatalf("unexpected block type, found: %T\n", block)
			}
		})
	}
}
