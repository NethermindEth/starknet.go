package rpc

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
func TestBlockID_Marshal(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	blockNumber := uint64(420)
	for _, test := range []struct {
		id      BlockID
		want    string
		wantErr error
	}{
		{
			id: BlockID{
				Tag: "latest",
			},
			want: `"latest"`,
		},
		{
			id: BlockID{
				Tag: "pre_confirmed",
			},
			want: `"pre_confirmed"`,
		},
		{
			id: BlockID{
				Tag: "l1_accepted",
			},
			want: `"l1_accepted"`,
		},
		{
			id: BlockID{
				Tag: "bad tag",
			},
			wantErr: ErrInvalidBlockID,
		},
		{
			id: BlockID{
				Number: &blockNumber,
			},
			want: `{"block_number":420}`,
		},
		{
			id: BlockID{
				Hash: internalUtils.TestHexToFelt(t, "0xdead"),
			},
			want: `{"block_hash":"0xdead"}`,
		},
	} {
		b, err := test.id.MarshalJSON()
		if test.wantErr != nil {
			require.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())

			return
		}
		require.NoError(t, err)

		assert.JSONEq(t, string(b), test.want)
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
				blockSepolia3100000 := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithReceipts/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/blockWithReceipts/sepoliaPreConfirmed.json", "result",
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

			rawExpectedBlock := testConfig.Spy.LastResponse()

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
