package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestChainID is a function that tests the ChainID function.
func TestChainID(t *testing.T) {
	tests.RunTestOn(
		t,
		tests.DevnetEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)

	testCase := map[tests.TestEnv]string{
		tests.DevnetEnv:      "SN_SEPOLIA",
		tests.IntegrationEnv: "SN_INTEGRATION_SEPOLIA",
		tests.MainnetEnv:     "SN_MAIN",
		tests.MockEnv:        "SN_SEPOLIA",
		tests.TestnetEnv:     "SN_SEPOLIA",
	}[tests.TEST_ENV]

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_chainId",
			).
			DoAndReturn(func(_, result, _ any, _ ...any) error {
				rawResp := result.(*json.RawMessage)
				*rawResp = json.RawMessage("\"0x534e5f5345504f4c4941\"") // "SN_SEPOLIA"

				return nil
			}).
			Times(1)
	}

	chain, err := testConfig.Provider.ChainID(t.Context())
	require.NoError(t, err)
	require.Equal(t, testCase, chain)
}

// TestSyncing tests the Syncing function.
func TestSyncing(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.MainnetEnv,
		tests.IntegrationEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_syncing",
			).
			DoAndReturn(func(_, result, _ any, _ ...any) error {
				rawResp := result.(*json.RawMessage)
				*rawResp = json.RawMessage(`
					{
						"starting_block_hash": "0x3be8e2915ffdb103cc0b999d7247e759b62c79d725634b943727bd4f7788c48",
						"starting_block_num": 3838962,
						"current_block_hash": "0x3be8e2915ffdb103cc0b999d7247e759b62c79d725634b943727bd4f7788c48",
						"current_block_num": 3838962,
						"highest_block_hash": "0x57d29828fad5bc1870c88b42685603fb9dab08c51118b9bf2e9474901b07501",
						"highest_block_num": 3839002
					}
				`)

				return nil
			}).
			Times(1)
	}

	sync, err := testConfig.Provider.Syncing(t.Context())
	require.NoError(t, err)

	if sync.IsSyncing {
		assert.NotZero(t, sync.StartingBlockHash)
		assert.NotZero(t, sync.StartingBlockNum)
		assert.NotZero(t, sync.CurrentBlockHash)
		assert.NotZero(t, sync.CurrentBlockNum)
		assert.NotZero(t, sync.HighestBlockHash)
		assert.NotZero(t, sync.HighestBlockNum)
	} else {
		assert.Zero(t, sync.StartingBlockHash)
		assert.Zero(t, sync.StartingBlockNum)
		assert.Zero(t, sync.CurrentBlockHash)
		assert.Zero(t, sync.CurrentBlockNum)
		assert.Zero(t, sync.HighestBlockHash)
		assert.Zero(t, sync.HighestBlockNum)
	}

	rawExpectedResp := testConfig.Spy.LastResponse()
	rawActualResp, err := json.Marshal(sync)
	require.NoError(t, err)
	assert.JSONEq(t, string(rawExpectedResp), string(rawActualResp))
}
