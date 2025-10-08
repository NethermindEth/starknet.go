package rpc

import (
	"context"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChainID is a function that tests the ChainID function in the Go test file.
func TestChainID(t *testing.T) {
	tests.RunTestOn(
		t,
		tests.MockEnv,
		tests.TestnetEnv,
		tests.MainnetEnv,
		tests.DevnetEnv,
		tests.IntegrationEnv,
	)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		ChainID string
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.DevnetEnv:      {{ChainID: "SN_SEPOLIA"}},
		tests.MainnetEnv:     {{ChainID: "SN_MAIN"}},
		tests.MockEnv:        {{ChainID: "SN_SEPOLIA"}},
		tests.TestnetEnv:     {{ChainID: "SN_SEPOLIA"}},
		tests.IntegrationEnv: {{ChainID: "SN_INTEGRATION_SEPOLIA"}},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		chain, err := testConfig.Provider.ChainID(context.Background())
		require.NoError(t, err)
		require.Equal(t, test.ChainID, chain)
	}
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
			StartingBlockHash: internalUtils.RANDOM_FELT,
			StartingBlockNum:  1234,
			CurrentBlockHash:  internalUtils.RANDOM_FELT,
			CurrentBlockNum:   1234,
			HighestBlockHash:  internalUtils.RANDOM_FELT,
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
