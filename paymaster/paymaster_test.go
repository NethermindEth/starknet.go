package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
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

		pm := SetupPaymaster(t)
		available, err := pm.IsAvailable(context.Background())
		require.NoError(t, err)
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

		pm := SetupPaymaster(t)
		tokens, err := pm.GetSupportedTokens(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, tokens)

		for _, token := range tokens {
			assert.NotNil(t, token.TokenAddress)
			assert.NotZero(t, token.Decimals)
			assert.NotZero(t, token.PriceInStrk)
		}
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

func TestOutsideExecutionTypedData(t *testing.T) {
	t.Run("GetOutsideExecutionTypedDataV1", func(t *testing.T) {
		message := OutsideExecutionMessageV1{
			Caller:        new(felt.Felt).SetUint64(1),
			Nonce:         new(felt.Felt).SetUint64(2),
			ExecuteAfter:  new(felt.Felt).SetUint64(3),
			ExecuteBefore: new(felt.Felt).SetUint64(4),
			CallsLen:      new(felt.Felt).SetUint64(1),
			Calls: []*OutsideCallV1{
				{
					To:          new(felt.Felt).SetUint64(5),
					Selector:    new(felt.Felt).SetUint64(6),
					CalldataLen: new(felt.Felt).SetUint64(1),
					Calldata:    []*felt.Felt{new(felt.Felt).SetUint64(7)},
				},
			},
		}
		typedData := GetOutsideExecutionTypedDataV1(message)

		assert.Equal(t, "OutsideExecution", typedData.PrimaryType)
		assert.Equal(t, "Account.execute_from_outside", typedData.Domain.Name)
		assert.Equal(t, "1", typedData.Domain.Version)
		assert.Equal(t, "0x534e5f4d41494e", typedData.Domain.ChainID)
		assert.Equal(t, message, typedData.Message)

		assert.Contains(t, typedData.Types, "StarkNetDomain")
		assert.Contains(t, typedData.Types, "OutsideExecution")
		assert.Contains(t, typedData.Types, "OutsideCall")
	})

	t.Run("GetOutsideExecutionTypedDataV2", func(t *testing.T) {
		message := OutsideExecutionMessageV2{
			Caller:        new(felt.Felt).SetUint64(1),
			Nonce:         new(felt.Felt).SetUint64(2),
			ExecuteAfter:  "3",
			ExecuteBefore: "4",
			Calls:         []Call{},
		}
		typedData := GetOutsideExecutionTypedDataV2(message)

		assert.Equal(t, "OutsideExecution", typedData.PrimaryType)
		assert.Equal(t, "Account.execute_from_outside", typedData.Domain.Name)
		assert.Equal(t, "2", typedData.Domain.Version)
		assert.Equal(t, "0x534e5f4d41494e", typedData.Domain.ChainID)
		assert.Equal(t, message, typedData.Message)

		assert.Contains(t, typedData.Types, "StarknetDomain")
		assert.Contains(t, typedData.Types, "OutsideExecution")
		assert.Contains(t, typedData.Types, "Call")
	})

	t.Run("GetOutsideExecutionTypedDataV3RC", func(t *testing.T) {
		message := OutsideExecutionMessageV3{
			Caller:        new(felt.Felt).SetUint64(1),
			Nonce:         new(felt.Felt).SetUint64(2),
			ExecuteAfter:  "3",
			ExecuteBefore: "4",
			Calls:         []Call{},
			Fee:           map[string]interface{}{"No Fee": "test"},
		}
		typedData := GetOutsideExecutionTypedDataV3RC(&message)

		assert.Equal(t, "OutsideExecution", typedData.PrimaryType)
		assert.Equal(t, "Account.execute_from_outside", typedData.Domain.Name)
		assert.Equal(t, "3", typedData.Domain.Version)
		assert.Equal(t, "0x534e5f4d41494e", typedData.Domain.ChainID)
		assert.Equal(t, message, typedData.Message)

		assert.Contains(t, typedData.Types, "StarknetDomain")
		assert.Contains(t, typedData.Types, "OutsideExecution")
		assert.Contains(t, typedData.Types, "Call")
		assert.Contains(t, typedData.Types, "Fee Mode")
		assert.Contains(t, typedData.Types, "Fee Transfer")
	})
}
