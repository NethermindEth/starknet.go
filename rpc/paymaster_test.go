package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			Invoke: UserInvoke{
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
				Invoke: UserInvoke{
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
		typedData := GetOutsideExecutionTypedDataV3RC(message)

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

func TestPaymasterClient(t *testing.T) {
	testConfig := beforeEach(t, false)
	client := &PaymasterClient{c: testConfig.provider.c}

	t.Run("IsAvailable", func(t *testing.T) {
		spy := NewSpy(client.c)
		client.c = spy
		_, err := client.IsAvailable(context.Background())
		require.NoError(t, err)
	})

	t.Run("GetSupportedTokens", func(t *testing.T) {
		spy := NewSpy(client.c)
		client.c = spy
		_, err := client.GetSupportedTokens(context.Background())
		require.NoError(t, err)
	})

}

func TestPaymasterConstants(t *testing.T) {
	assert.Equal(t, "OUTSIDE_EXECUTION_TYPED_DATA_V1", OUTSIDE_EXECUTION_TYPED_DATA_V1)
	assert.Equal(t, "OUTSIDE_EXECUTION_TYPED_DATA_V2", OUTSIDE_EXECUTION_TYPED_DATA_V2)
	assert.Equal(t, "OUTSIDE_EXECUTION_TYPED_DATA_V3_RC", OUTSIDE_EXECUTION_TYPED_DATA_V3_RC)
}
