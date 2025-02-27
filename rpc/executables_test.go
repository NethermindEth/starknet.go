package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompiledCasm(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		Description        string
		ClassHash          *felt.Felt
		ExpectedResultPath string
		ExpectedError      *RPCError
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				Description:        "success - get compiled CASM",
				ClassHash:          utils.RANDOM_FELT,
				ExpectedResultPath: "./tests/compiledCasm.json",
			},
			{
				Description:   "error - class hash not found",
				ClassHash:     utils.TestHexToFelt(t, "0xdadadadada"),
				ExpectedError: ErrClassHashNotFound,
			},
			{
				Description:   "error - compilation error",
				ClassHash:     utils.TestHexToFelt(t, "0xbad"),
				ExpectedError: ErrCompilationError,
			},
		},
		"devnet": {},
		"testnet": {
			{
				Description:        "normal call, with field class_hash",
				ClassHash:          utils.TestHexToFelt(t, "0x00d764f235da1c654c4ca14c47bfc2a54ccd4c0c56b3f4570cd241bd638db448"),
				ExpectedResultPath: "./tests/compiledCasm.json",
			},
			{
				Description:   "error call, inexistent class_hash",
				ClassHash:     utils.TestHexToFelt(t, "0xdedededededede"),
				ExpectedError: ErrClassHashNotFound,
			},
			// TODO: add test for compilation error. Need to find the invalid class hash.
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		var expectedResult *CasmCompiledContractClass
		var rawFile []byte
		var err error

		if test.ExpectedResultPath != "" {
			// getting the expected result from the file
			rawFile, err = os.ReadFile(test.ExpectedResultPath)
			require.NoError(t, err)
			var rpcResponse struct {
				Result json.RawMessage `json:"result"`
			}
			err = json.Unmarshal(rawFile, &rpcResponse)
			require.NoError(t, err)
			rawFile = rpcResponse.Result
			err = json.Unmarshal(rawFile, &expectedResult)
			require.NoError(t, err)
		}

		// getting the result from the provider and asserting equality
		result, err := testConfig.provider.CompiledCasm(context.Background(), test.ClassHash)
		if test.ExpectedError != nil {
			rpcErr, ok := err.(*RPCError)
			require.True(t, ok)
			assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		// asserting equality of the json results
		jsonResult, err := json.Marshal(result)
		assert.NoError(t, err)

		assert.JSONEq(t, string(rawFile), string(jsonResult))
	}
}
