package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockPaymaster struct {
	*Paymaster
	// this should be a pointer to the mock client used in the Paymaster struct.
	// This is intended to have an easy access to the mock client, without having to
	// type cast it from the `callCloser` interface every time.
	c *mocks.MockClient
}

// Creates a real Sepolia paymaster client.
func SetupPaymaster(t *testing.T) *Paymaster {
	t.Helper()
	pm, err := NewPaymasterClient("https://sepolia.paymaster.avnu.fi")
	require.NoError(t, err, "failed to create paymaster client")

	return pm
}

// Creates a mock paymaster client.
func SetupMockPaymaster(t *testing.T) *MockPaymaster {
	t.Helper()

	client := mocks.NewMockClient(gomock.NewController(t))
	mpm := &MockPaymaster{
		Paymaster: &Paymaster{c: client},
		c:         client,
	}

	return mpm
}

// Test the 'paymaster_isAvailable' method
func TestIsAvailable(t *testing.T) {
	t.Parallel()
	t.Run("integration", func(t *testing.T) {
		t.Parallel()
		pm := SetupPaymaster(t)
		available, err := pm.IsAvailable(context.Background())
		require.NoError(t, err)
		assert.True(t, available)
	})

	t.Run("mock", func(t *testing.T) {
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

func TestPaymasterTypes(t *testing.T) {
	t.Run("Call", func(t *testing.T) {
		call := Call{
			To:       &felt.Felt{},
			Selector: &felt.Felt{},
			Calldata: []*felt.Felt{},
		}

		assert.NotNil(t, call)
		assert.NotNil(t, call.To)
		assert.NotNil(t, call.Selector)
		assert.NotNil(t, call.Calldata)
	})

	t.Run("UserInvoke", func(t *testing.T) {
		invoke := UserInvoke{
			UserAddress: &felt.Felt{},
			Calls:       []Call{},
		}

		assert.NotNil(t, invoke)
		assert.NotNil(t, invoke.UserAddress)
		assert.NotNil(t, invoke.Calls)
	})

	t.Run("UserTransaction", func(t *testing.T) {
		transaction := UserTransaction{
			Type: "invoke",
			Invoke: &UserInvoke{
				UserAddress: &felt.Felt{},
				Calls:       []Call{},
			},
		}

		assert.NotNil(t, transaction)
		assert.Equal(t, "invoke", transaction.Type)
		assert.NotNil(t, transaction.Invoke)
	})

	t.Run("FeeMode", func(t *testing.T) {
		feeMode := FeeMode{
			Mode: "sponsored",
		}

		assert.NotNil(t, feeMode)
		assert.Equal(t, "sponsored", feeMode.Mode)
	})

	t.Run("UserParameters", func(t *testing.T) {
		params := UserParameters{
			Version: "0x1",
			FeeMode: FeeMode{
				Mode: "sponsored",
			},
		}

		assert.NotNil(t, params)
		assert.Equal(t, "0x1", params.Version)
		assert.NotNil(t, params.FeeMode)
	})

	t.Run("BuildTransactionRequest", func(t *testing.T) {
		request := BuildTransactionRequest{
			Transaction: UserTransaction{
				Type: "invoke",
				Invoke: &UserInvoke{
					UserAddress: &felt.Felt{},
					Calls:       []Call{},
				},
			},
			Parameters: UserParameters{
				Version: "0x1",
				FeeMode: FeeMode{
					Mode: "sponsored",
				},
			},
		}

		assert.NotNil(t, request)
		assert.NotNil(t, request.Transaction)
		assert.NotNil(t, request.Parameters)
	})

	t.Run("FeeEstimateResponse", func(t *testing.T) {
		estimate := FeeEstimateResponse{
			GasTokenPriceInStrk:       &felt.Felt{},
			EstimatedFeeInStrk:        &felt.Felt{},
			EstimatedFeeInGasToken:    &felt.Felt{},
			SuggestedMaxFeeInStrk:     &felt.Felt{},
			SuggestedMaxFeeInGasToken: &felt.Felt{},
		}

		assert.NotNil(t, estimate)
		assert.NotNil(t, estimate.GasTokenPriceInStrk)
		assert.NotNil(t, estimate.EstimatedFeeInStrk)
		assert.NotNil(t, estimate.EstimatedFeeInGasToken)
		assert.NotNil(t, estimate.SuggestedMaxFeeInStrk)
		assert.NotNil(t, estimate.SuggestedMaxFeeInGasToken)
	})

	t.Run("ExecutableUserInvoke", func(t *testing.T) {
		invoke := ExecutableUserInvoke{
			UserAddress: &felt.Felt{},
			TypedData:   map[string]interface{}{},
			Signature:   []*felt.Felt{},
		}

		assert.NotNil(t, invoke)
		assert.NotNil(t, invoke.UserAddress)
		assert.NotNil(t, invoke.TypedData)
		assert.NotNil(t, invoke.Signature)
	})

	t.Run("TokenData", func(t *testing.T) {
		token := TokenData{
			TokenAddress: &felt.Felt{},
			Decimals:     18,
			PriceInStrk:  "0x1234567890abcdef",
		}

		assert.NotNil(t, token)
		assert.NotNil(t, token.TokenAddress)
		assert.Equal(t, 18, token.Decimals)
		assert.Equal(t, "0x1234567890abcdef", token.PriceInStrk)
	})

	t.Run("TrackingIdResponse", func(t *testing.T) {
		response := TrackingIdResponse{
			TransactionHash: &felt.Felt{},
			Status:          "active",
		}

		assert.NotNil(t, response)
		assert.NotNil(t, response.TransactionHash)
		assert.Equal(t, "active", response.Status)
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

// func TestPaymasterClient(t *testing.T) {
// 	testConfig := rpc.BeforeEach(t, false)
// 	client := &Paymaster{c: testConfig.Provider.c}

// 	t.Run("IsAvailable", func(t *testing.T) {
// 		_, err := client.IsAvailable(context.Background())
// 		require.NoError(t, err)
// 	})

// 	t.Run("GetSupportedTokens", func(t *testing.T) {
// 		_, err := client.GetSupportedTokens(context.Background())
// 		require.NoError(t, err)
// 	})

// }

func TestPaymasterConstants(t *testing.T) {
	assert.Equal(t, "OUTSIDE_EXECUTION_TYPED_DATA_V1", OUTSIDE_EXECUTION_TYPED_DATA_V1)
	assert.Equal(t, "OUTSIDE_EXECUTION_TYPED_DATA_V2", OUTSIDE_EXECUTION_TYPED_DATA_V2)
	assert.Equal(t, "OUTSIDE_EXECUTION_TYPED_DATA_V3_RC", OUTSIDE_EXECUTION_TYPED_DATA_V3_RC)
}
