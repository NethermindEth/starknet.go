package rpc

import (
	"context"
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
			ExpectedResp: rpcVersion,
		}},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.SpecVersion(context.Background())
		require.NoError(t, err)
		require.Equal(t, test.ExpectedResp, resp)
	}
}
