package rpc

import (
	"context"
	"testing"

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
