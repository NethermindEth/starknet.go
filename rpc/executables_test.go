package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompiledCasm(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description        string
		ClassHash          *felt.Felt
		ExpectedResultPath string
		ExpectedError      *RPCError
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:        "success - get compiled CASM",
				ClassHash:          internalUtils.DeadBeef,
				ExpectedResultPath: "./testData/compiledCasm/sepolia.json",
			},
			{
				Description:   "error - class hash not found",
				ClassHash:     internalUtils.TestHexToFelt(t, "0xdadadadada"),
				ExpectedError: ErrClassHashNotFound,
			},
			{
				Description:   "error - compilation error",
				ClassHash:     internalUtils.TestHexToFelt(t, "0xbad"),
				ExpectedError: ErrCompilationError,
			},
		},
		tests.TestnetEnv: {
			{
				Description:        "normal call, with field class_hash",
				ClassHash:          internalUtils.TestHexToFelt(t, "0x00d764f235da1c654c4ca14c47bfc2a54ccd4c0c56b3f4570cd241bd638db448"),
				ExpectedResultPath: "./testData/compiledCasm/sepolia.json",
			},
			{
				Description:   "error call, inexistent class_hash",
				ClassHash:     internalUtils.TestHexToFelt(t, "0xdedededededede"),
				ExpectedError: ErrClassHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:        "normal call, with field class_hash",
				ClassHash:          internalUtils.TestHexToFelt(t, "0x941a2dc3ab607819fdc929bea95831a2e0c1aab2f2f34b3a23c55cebc8a040"),
				ExpectedResultPath: "./testData/compiledCasm/integration.json",
			},
			{
				Description:   "error call, inexistent class_hash",
				ClassHash:     internalUtils.TestHexToFelt(t, "0xdedededededede"),
				ExpectedError: ErrClassHashNotFound,
			},
			// TODO: add test for compilation error when Juno implements it (maybe the class hash from block 1 could be a valid input)
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		expectedResult, err := internalUtils.UnmarshalJSONFileToType[contracts.CasmClass](
			test.ExpectedResultPath,
			"result",
		)
		if test.ExpectedResultPath != "" {
			require.NoError(t, err)
		}

		// getting the result from the provider and asserting equality
		result, err := testConfig.Provider.CompiledCasm(context.Background(), test.ClassHash)
		if test.ExpectedError != nil {
			rpcErr, ok := err.(*RPCError)
			require.True(t, ok)
			assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
			assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

			continue
		}
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		// asserting equality of the json results
		resultJSON, err := json.Marshal(result)
		require.NoError(t, err)

		rawFile, err := os.ReadFile(test.ExpectedResultPath)
		require.NoError(t, err)
		var rpcResponse struct {
			Result json.RawMessage `json:"result"`
		}
		err = json.Unmarshal(rawFile, &rpcResponse)
		require.NoError(t, err)
		expectedResultJSON, err := json.Marshal(rpcResponse.Result)
		require.NoError(t, err)

		assert.JSONEq(t, string(expectedResultJSON), string(resultJSON))
	}
}
