package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestCompiledCasm tests the CompiledCasm function.
func TestCompiledCasm(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)
	t.Parallel()
	t.Parallel()

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider

	type testSetType struct {
		Description   string
		ClassHash     *felt.Felt
		ExpectedError *RPCError
		Description   string
		ClassHash     *felt.Felt
		ExpectedError *RPCError
	}

	// TODO: use the 'testData/compiledCasm' folder for mock tests

	// TODO: use the 'testData/compiledCasm' folder for mock tests
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "normal call",
				ClassHash:   internalUtils.TestHexToFelt(t, "0xdadadadada"),
			},
			{
				Description:   "class hash not found",
				ClassHash:     internalUtils.DeadBeef,
				ExpectedError: ErrClassHashNotFound,
			},
			{
				Description:   "compilation error",
				ClassHash:     internalUtils.TestHexToFelt(t, "0xbad"),
				ExpectedError: ErrCompilationError,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call",
				ClassHash:   internalUtils.TestHexToFelt(t, "0x00d764f235da1c654c4ca14c47bfc2a54ccd4c0c56b3f4570cd241bd638db448"),
			},
			{
				Description:   "class hash not found",
				ClassHash:     internalUtils.DeadBeef,
				ExpectedError: ErrClassHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description: "normal call",
				ClassHash:   internalUtils.TestHexToFelt(t, "0x941a2dc3ab607819fdc929bea95831a2e0c1aab2f2f34b3a23c55cebc8a040"),
			},
			{
				Description:   "class hash not found",
				ClassHash:     internalUtils.DeadBeef,
				ExpectedError: ErrClassHashNotFound,
			},
			// TODO: add test for compilation error when Juno implements it (maybe the class hash from block 1 could be a valid input)
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			t.Parallel()

			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getCompiledCasm",
						test.ClassHash,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						classHash := args[0].(*felt.Felt)

						if classHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    28,
								Message: "Class hash not found",
							}
						}

						if classHash.String() == "0xbad" {
							return RPCError{
								Code:    100,
								Message: "Failed to compile the contract",
								Data:    &CompilationErrData{},
							}
						}

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/compiledCasm/sepolia.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			// getting the result from the provider and asserting equality
			result, err := provider.CompiledCasm(t.Context(), test.ClassHash)
			if test.ExpectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

				return
			}
			require.NoError(t, err)
			rawExpectedResult := testConfig.RPCSpy.LastResponse()

			// asserting equality of the json results
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResult), string(resultJSON))
		})
			// asserting equality of the json results
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResult), string(resultJSON))
		})
	}
}
