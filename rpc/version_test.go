package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/require"
)

// TestSpecVersion tests starknet_specVersion
func TestSpecVersion(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		ExpectedResp string
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {{
			ExpectedResp: rpcVersion.String(),
		}},
		tests.IntegrationEnv: {{
			ExpectedResp: rpcVersion.String(),
		}},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.Provider.SpecVersion(context.Background())
		require.NoError(t, err)
		require.Equal(t, test.ExpectedResp, resp)
	}
}
