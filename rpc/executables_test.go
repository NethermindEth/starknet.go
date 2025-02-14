package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompiledCasm(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Description        string
		ClassHash          *felt.Felt
		ExpectedResultPath string
		ExpectedError      error
	}
	testSet := map[string][]testSetType{
		"mock":   {},
		"devnet": {},
		"testnet": {
			{
				Description:        "normal call, only required field class_hash",
				ClassHash:          utils.TestHexToFelt(t, "0x00d764f235da1c654c4ca14c47bfc2a54ccd4c0c56b3f4570cd241bd638db448"),
				ExpectedResultPath: "./tests/compiledCasm.json",
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		// getting the expected result from the file
		rawFile, err := os.ReadFile(test.ExpectedResultPath)
		require.NoError(t, err)
		var rpcResponse struct {
			Result json.RawMessage `json:"result"`
		}
		err = json.Unmarshal(rawFile, &rpcResponse)
		require.NoError(t, err)
		rawFile = rpcResponse.Result
		var expectedResult CasmCompiledContractClass
		err = json.Unmarshal(rawFile, &expectedResult)
		require.NoError(t, err)

		// getting the result from the provider and asserting equality
		result, err := testConfig.provider.CompiledCasm(context.Background(), test.ClassHash)
		assert.NoError(t, err)
		assert.Equal(t, &expectedResult, result)

		// asserting equality of the json results
		jsonResult, err := json.Marshal(result)
		assert.NoError(t, err)

		assert.JSONEq(t, string(rawFile), string(jsonResult))
	}
}
