package rpc

import (
	"context"
	"strconv"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCError(t *testing.T) {
	if tests.TEST_ENV == tests.MockEnv {
		testConfig := beforeEach(t, false)
		_, err := testConfig.provider.ChainID(context.Background())
		require.NoError(t, err)

		_, err = testConfig.provider.Events(context.Background(), EventsInput{ResultPageRequest: ResultPageRequest{ChunkSize: 0}})
		require.Error(t, err)
		rpcErr := err.(*RPCError)
		assert.Equal(t, rpcErr.Code, InternalError)
		assert.NotNil(t, rpcErr.Message, "Internal Error")
		assert.NotNil(t, rpcErr.Data, "-ChuckSize error message-")

		assert.ErrorContains(t, err, strconv.Itoa(rpcErr.Code))
		assert.ErrorContains(t, err, rpcErr.Message)
		assert.ErrorContains(t, err, rpcErr.Data.ErrorMessage())
	}

	if tests.TEST_ENV == tests.TestnetEnv {
		testConfig := beforeEach(t, false)

		// invalid msg
		msgFromL1 := MsgFromL1{
			FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
			ToAddress:   internalUtils.RANDOM_FELT,
			Selector:    internalUtils.RANDOM_FELT,
			Payload:     []*felt.Felt{},
		}

		_, err := testConfig.provider.EstimateMessageFee(context.Background(), msgFromL1, WithBlockNumber(523066))
		require.Error(t, err)
		rpcErr := err.(*RPCError)

		// check if the error code, message, and data are correct
		assert.Equal(t, rpcErr.Code, ErrContractError.Code)
		assert.Equal(t, rpcErr.Message, ErrContractError.Message)
		assert.IsType(t, rpcErr.Data, ErrContractError.Data)
		assert.NotEmpty(t, rpcErr.Data)

		// check if the error message contains the error code, message, and data
		assert.ErrorContains(t, err, strconv.Itoa(rpcErr.Code))
		assert.ErrorContains(t, err, rpcErr.Message)
		assert.ErrorContains(t, err, rpcErr.Data.ErrorMessage())
	}
}
