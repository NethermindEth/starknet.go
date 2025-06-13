package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/starknet.go/internal"
	"github.com/stretchr/testify/require"
)

// TestSpecVersion tests starknet_specVersion
func TestSpecVersion(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		ExpectedResp string
	}
	testSet := map[internal.TestEnv][]testSetType{
		internal.DevnetEnv:  {},
		internal.MainnetEnv: {},
		internal.MockEnv:    {},
		internal.TestnetEnv: {{
			ExpectedResp: rpcVersion,
		}},
	}[internal.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.provider.SpecVersion(context.Background())
		require.NoError(t, err)
		require.Equal(t, test.ExpectedResp, resp)
	}
}
