package methods

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	. "github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestBlockWithReceipts tests the BlockWithReceipts function.
func TestBlockWithReceipts(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := internal.BeforeEach(t, false)
	provider := testConfig.Provider

	type testSetType struct {
		BlockID     rpcv10.BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagPreConfirmed),
			},
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
			},
			{
				BlockID:     rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID:     rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID:     rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID:     rpcv10.WithBlockHash(internalUtils.DeadBeef),
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
							blockID := args[0].(rpcv10.BlockID)

							switch blockID.Tag {
							case rpcv10.BlockTagPreConfirmed:
								*rawResp = blockSepoliaPreConfirmed
							case rpcv10.BlockTagLatest:
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
			result, err := GetBlockWithReceipts(t.Context(), testConfig.Provider, test.BlockID)
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
