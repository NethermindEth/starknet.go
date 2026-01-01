package rpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
