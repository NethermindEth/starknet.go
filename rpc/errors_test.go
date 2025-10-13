package rpc

import (
	"context"
	"strconv"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCError(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	if tests.TEST_ENV == tests.MockEnv {
		testConfig := BeforeEach(t, false)
		_, err := testConfig.Provider.ChainID(context.Background())
		require.NoError(t, err)

		_, err = testConfig.Provider.Events(
			context.Background(),
			EventsInput{ResultPageRequest: ResultPageRequest{ChunkSize: 0}},
		)
		require.Error(t, err)
		rpcErr := err.(*RPCError)
		assert.Equal(t, rpcErr.Code, rpcerr.InternalError)
		assert.NotNil(t, rpcErr.Message, "Internal Error")
		assert.NotNil(t, rpcErr.Data, "-ChuckSize error message-")

		assert.ErrorContains(t, err, strconv.Itoa(rpcErr.Code))
		assert.ErrorContains(t, err, rpcErr.Message)
		assert.ErrorContains(t, err, rpcErr.Data.ErrorMessage())

		return
	}

	testConfig := BeforeEach(t, false)

	// invalid msg
	msgFromL1 := MsgFromL1{
		FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
		ToAddress:   internalUtils.DeadBeef,
		Selector:    internalUtils.DeadBeef,
		Payload:     []*felt.Felt{},
	}

	_, err := testConfig.Provider.EstimateMessageFee(
		context.Background(),
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
