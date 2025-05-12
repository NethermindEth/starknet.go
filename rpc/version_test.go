package rpc

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSpecVersion tests starknet_specVersion
func TestSpecVersion(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		ExpectedResp string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			ExpectedResp: "0.8.1",
		}},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.SpecVersion(context.Background())
		require.NoError(t, err)
		require.Equal(t, test.ExpectedResp, resp)
	}
}

// TestVersionCompatibility tests that the provider correctly handles version compatibility
func TestVersionCompatibility(t *testing.T) {
	if testEnv == "mock" {
		t.Skip("Skipping integration test in mock environment")
	}

	// Create a new provider
	provider, err := NewProvider(os.Getenv("HTTP_PROVIDER_URL"))
	require.NoError(t, err)
	require.NotNil(t, provider)

	// Get the version from the provider
	version, err := provider.SpecVersion(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, version)

	// Check that the version contains our expected version
	require.Contains(t, version, RPCVersion)
}
