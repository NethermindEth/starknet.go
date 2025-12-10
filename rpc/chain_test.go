package rpc

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
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

// TestSyncing is a test function that tests the syncing functionality of the provider.
//
// It checks the synchronisation status and verifies the values returned by the provider.
// The test is performed for different test environments, such as devnet, mainnet, mock, and testnet.
// For each test environment, it retrieves the synchronisation status from the provider and performs the necessary assertions.
// If the test environment is "mock", it verifies that the returned values match the expected values.
// Otherwise, it checks that the synchronisation status is false and verifies the values returned by the provider.
// The function uses the testing.T type for assertions and the context.Background() function for the context.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestSyncing(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	sync, err := testConfig.Provider.Syncing(context.Background())
	require.NoError(t, err)

	if tests.TEST_ENV == tests.MockEnv {
		value := SyncStatus{
			IsSyncing:         true,
			StartingBlockHash: internalUtils.DeadBeef,
			StartingBlockNum:  1234,
			CurrentBlockHash:  internalUtils.DeadBeef,
			CurrentBlockNum:   1234,
			HighestBlockHash:  internalUtils.DeadBeef,
			HighestBlockNum:   1234,
		}
		assert.Exactly(t, value, sync)

		return
	}

	if sync.IsSyncing {
		require.True(
			t,
			strings.HasPrefix(sync.CurrentBlockHash.String(), "0x"),
			"current block hash should return a string starting with 0x",
		)
		assert.NotZero(t, sync.StartingBlockHash)
		assert.NotZero(t, sync.StartingBlockNum)
		assert.NotZero(t, sync.CurrentBlockHash)
		assert.NotZero(t, sync.CurrentBlockNum)
		assert.NotZero(t, sync.HighestBlockHash)
		assert.NotZero(t, sync.HighestBlockNum)
	} else {
		assert.False(t, sync.IsSyncing)
		assert.Zero(t, sync.StartingBlockHash)
		assert.Zero(t, sync.StartingBlockNum)
		assert.Zero(t, sync.CurrentBlockHash)
		assert.Zero(t, sync.CurrentBlockNum)
		assert.Zero(t, sync.HighestBlockHash)
		assert.Zero(t, sync.HighestBlockNum)
	}
}
