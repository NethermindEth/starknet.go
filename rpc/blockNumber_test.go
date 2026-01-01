package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestBlockNumber tests the BlockNumber function.
func TestBlockNumber(t *testing.T) {
	tests.RunTestOn(t,
		tests.DevnetEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_blockNumber",
			).
			DoAndReturn(
				func(_, result, _ any, _ ...any) error {
					rawResp := result.(*json.RawMessage)
					rawBlockNumber := json.RawMessage("1234")
					*rawResp = rawBlockNumber

					return nil
				},
			).
			Times(1)
	}

	blockNumber, err := provider.BlockNumber(t.Context())
	require.NoError(t, err)

	rawExpectedResp := testConfig.RPCSpy.LastResponse()
	rawActualResp, err := json.Marshal(blockNumber)
	require.NoError(t, err)
	assert.JSONEq(t, string(rawExpectedResp), string(rawActualResp))
}
