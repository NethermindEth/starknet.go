package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestSpecVersion tests the SpecVersion function.
func TestSpecVersion(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_specVersion",
			).
			DoAndReturn(func(_, result, _ any, _ ...any) error {
				rawResp := result.(*json.RawMessage)
				*rawResp = json.RawMessage("\"0.10.0\"")

				return nil
			}).
			Times(1)
	}

	resp, err := testConfig.Provider.SpecVersion(t.Context())
	require.NoError(t, err)

	rawExpectedResp := testConfig.RPCSpy.LastResponse()
	rawResp, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.Equal(t, string(rawExpectedResp), string(rawResp))
}
