package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
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
