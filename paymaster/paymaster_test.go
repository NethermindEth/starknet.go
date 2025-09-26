package paymaster

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Test the 'paymaster_isAvailable' method
func TestIsAvailable(t *testing.T) {
	t.Parallel()
	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)
		t.Parallel()

		pm, spy := SetupPaymaster(t)
		available, err := pm.IsAvailable(context.Background())
		require.NoError(t, err)

		assert.Equal(t, string(spy.LastResponse()), strconv.FormatBool(available))
		assert.True(t, available)
	})

	t.Run("mock", func(t *testing.T) {
		tests.RunTestOn(t, tests.MockEnv)
		t.Parallel()

		pm := SetupMockPaymaster(t)
		pm.c.EXPECT().
			CallContextWithSliceArgs(context.Background(), gomock.AssignableToTypeOf(new(bool)), "paymaster_isAvailable").
			SetArg(1, true).
			Return(nil)
		available, err := pm.IsAvailable(context.Background())
		assert.NoError(t, err)
		assert.True(t, available)
	})
}

// Test the 'paymaster_getSupportedTokens' method
func TestGetSupportedTokens(t *testing.T) {
	t.Parallel()
	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)
		t.Parallel()

		pm, spy := SetupPaymaster(t)
		tokens, err := pm.GetSupportedTokens(context.Background())
		require.NoError(t, err)

		rawResult, err := json.Marshal(tokens)
		require.NoError(t, err)
		assert.EqualValues(t, spy.LastResponse(), rawResult)
	})

	t.Run("mock", func(t *testing.T) {
		tests.RunTestOn(t, tests.MockEnv)
		t.Parallel()

		pm := SetupMockPaymaster(t)

		expectedRawResult := `[
			{
				"token_address": "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
				"decimals": 18,
				"price_in_strk": "0x288aa92ed8c5539ae80"
			},
			{
				"token_address": "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d",
				"decimals": 18,
				"price_in_strk": "0xde0b6b3a7640000"
			},
			{
				"token_address": "0x53b40a647cedfca6ca84f542a0fe36736031905a9639a7f19a3c1e66bfd5080",
				"decimals": 6,
				"price_in_strk": "0x48e1ecdbbe883b08"
			},
			{
				"token_address": "0x30058f19ed447208015f6430f0102e8ab82d6c291566d7e73fe8e613c3d2ed",
				"decimals": 6,
				"price_in_strk": "0x2c3460a7992f8a"
			}
		]`

		var expectedResult []TokenData
		err := json.Unmarshal([]byte(expectedRawResult), &expectedResult)
		require.NoError(t, err)

		pm.c.EXPECT().
			CallContextWithSliceArgs(context.Background(), gomock.AssignableToTypeOf(new([]TokenData)), "paymaster_getSupportedTokens").
			SetArg(1, expectedResult).
			Return(nil)
		result, err := pm.GetSupportedTokens(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		rawResult, err := json.Marshal(result)
		require.NoError(t, err)
		assert.JSONEq(t, expectedRawResult, string(rawResult))
	})
}
