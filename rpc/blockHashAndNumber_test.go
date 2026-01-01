package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestBlockHashAndNumber tests the BlockHashAndNumber function.
func TestBlockHashAndNumber(t *testing.T) {
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
				"starknet_blockHashAndNumber",
			).
			DoAndReturn(
				func(_, result, _ any, _ ...any) error {
					rawResp := result.(*json.RawMessage)
					rawBlockHashAndNumber := json.RawMessage(
						`{
							"block_hash": "0x7fcc97a2e4e4a328582326254baca628cad2a82b17a711e7a8e5c9edd8022e6",
							"block_number": 3640605
						}`,
					)
					*rawResp = rawBlockHashAndNumber

					return nil
				},
			).
			Times(1)
	}

	blockHashAndNumber, err := provider.BlockHashAndNumber(t.Context())
	require.NoError(t, err, "BlockHashAndNumber should not return an error")

	rawExpectedResp := testConfig.RPCSpy.LastResponse()
	rawActualResp, err := json.Marshal(blockHashAndNumber)
	require.NoError(t, err)
	assert.JSONEq(t, string(rawExpectedResp), string(rawActualResp))
}
