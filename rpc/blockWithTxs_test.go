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

			rawExpectedBlock := testConfig.RPCSpy.LastResponse()

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
