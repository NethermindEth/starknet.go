package rpc

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestRPCError tests the RPCError type, checking if a given error is unwrapped correctly.
func TestRPCError(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	// invalid msg - empty payload
	msgFromL1 := MsgFromL1{
		FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
		ToAddress: internalUtils.TestHexToFelt(
			t,
			"0x04c5772d1914fe6ce891b64eb35bf3522aeae1315647314aac58b01137607f3f",
		),
		Selector: internalUtils.TestHexToFelt(
			t,
			"0x1b64b1b3b690b43b9b514fb81377518f4039cd3e4f4914d8a6bdf01d679fb19",
		),
		Payload: []*felt.Felt{},
	}

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_estimateMessageFee",
				msgFromL1,
				WithBlockNumber(523066),
			).
			DoAndReturn(func(_, result, _ any, args ...any) error {
				rpcErr := internalUtils.TestUnmarshalJSONFileToType[RPCError](
					t,
					"./testData/errors/contractError.json",
					"error",
				)
				rawData, ok := rpcErr.Data.(StringErrData)
				require.True(t, ok)
				var contractErrData ContractErrData
				err := json.Unmarshal([]byte(rawData), &contractErrData)
				require.NoError(t, err)
				rpcErr.Data = &contractErrData

				return rpcErr
			}).
			Times(1)
	}

	_, err := testConfig.Provider.EstimateMessageFee(
		t.Context(),
		msgFromL1,
		WithBlockNumber(523066),
	)
	require.Error(t, err)
	rpcErr := err.(*RPCError)

	// check if the error code, message, and data are correct
	assert.Equal(t, ErrContractError.Code, rpcErr.Code)
	assert.Equal(t, ErrContractError.Message, rpcErr.Message)
	assert.IsType(t, ErrContractError.Data, rpcErr.Data)
	assert.NotEmpty(t, rpcErr.Data)

	// check if the error message contains the error code, message, and data
	assert.ErrorContains(t, err, strconv.Itoa(rpcErr.Code))
	assert.ErrorContains(t, err, rpcErr.Message)
	assert.ErrorContains(t, err, rpcErr.Data.ErrorMessage())
}
