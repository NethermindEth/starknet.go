package rpc

import (
	"context"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

// TestChainID is a function that tests the ChainID function in the Go test file.
//
// The function initializes a test configuration and defines a test set with different chain IDs for different environments.
// It then iterates over the test set and for each test, creates a new spy and sets the spy as the provider's client.
// The function calls the ChainID function and compares the returned chain ID with the expected chain ID.
// If there is a mismatch or an error occurs, the function logs a fatal error.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestChainID(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ChainID string
	}
	testSet := map[string][]testSetType{
		"devnet":  {{ChainID: "SN_SEPOLIA"}},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "SN_SEPOLIA"}},
		"testnet": {{ChainID: "SN_SEPOLIA"}},
	}[testEnv]

	for _, test := range testSet {
		chain, err := testConfig.provider.ChainID(context.Background())
		require.NoError(t, err)
		require.Equal(t, test.ChainID, chain)
	}
}

// TestSyncing is a test function that tests the syncing functionality of the provider.
//
// It checks the synchronization status and verifies the values returned by the provider.
// The test is performed for different test environments, such as devnet, mainnet, mock, and testnet.
// For each test environment, it retrieves the synchronization status from the provider and performs the necessary assertions.
// If the test environment is "mock", it verifies that the returned values match the expected values.
// Otherwise, it checks that the synchronization status is false and verifies the values returned by the provider.
// The function uses the testing.T type for assertions and the context.Background() function for the context.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestSyncing(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ChainID string
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_SEPOLIA"}},
	}[testEnv]

	for range testSet {
		sync, err := testConfig.provider.Syncing(context.Background())
		require.NoError(t, err)

		if testEnv == "mock" {
			value := SyncStatus{
				StartingBlockHash: utils.RANDOM_FELT,
				StartingBlockNum:  "0x4c602",
				CurrentBlockHash:  utils.RANDOM_FELT,
				CurrentBlockNum:   "0x4c727",
				HighestBlockHash:  utils.RANDOM_FELT,
				HighestBlockNum:   "0x4c727",
			}
			require.Exactly(t, &value, sync)

			continue
		}

		if sync.SyncStatus == nil {
			require.True(t, strings.HasPrefix(sync.CurrentBlockHash.String(), "0x"), "current block hash should return a string starting with 0x")
			require.NotZero(t, sync.StartingBlockHash)
			require.NotZero(t, sync.StartingBlockNum)
			require.NotZero(t, sync.CurrentBlockHash)
			require.NotZero(t, sync.CurrentBlockNum)
			require.NotZero(t, sync.HighestBlockHash)
			require.NotZero(t, sync.HighestBlockNum)
		} else {
			require.False(t, *sync.SyncStatus)
			require.Zero(t, sync.StartingBlockHash)
			require.Zero(t, sync.StartingBlockNum)
			require.Zero(t, sync.CurrentBlockHash)
			require.Zero(t, sync.CurrentBlockNum)
			require.Zero(t, sync.HighestBlockHash)
			require.Zero(t, sync.HighestBlockNum)
		}

	}
}
