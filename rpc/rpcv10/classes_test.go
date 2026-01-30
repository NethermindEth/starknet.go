package rpcv10

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestClass tests the Class function.
func TestClass(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.IntegrationEnv)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description   string
		blockID       types.BlockID
		ClassHash     *felt.Felt
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "deprecated class",
				ClassHash:   internalUtils.TestHexToFelt(t, "0x123"),
				blockID:     types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description: "sierra class",
				ClassHash:   internalUtils.TestHexToFelt(t, "0x456"),
				blockID:     types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:   "invalid block",
				ClassHash:     internalUtils.TestHexToFelt(t, "0x789"),
				blockID:       types.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description:   "invalid class hash",
				ClassHash:     internalUtils.DeadBeef,
				blockID:       types.WithBlockTag(types.BlockTagLatest),
				ExpectedError: ErrClassHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "deprecated class",
				blockID:     types.WithBlockTag(types.BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
			},
			{
				Description: "sierra class",
				blockID:     types.WithBlockTag(types.BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
			},
			{
				Description:   "invalid block",
				ClassHash:     internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
				blockID:       types.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description:   "invalid class hash",
				ClassHash:     internalUtils.DeadBeef,
				blockID:       types.WithBlockTag(types.BlockTagLatest),
				ExpectedError: ErrClassHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description: "sierra class",
				blockID:     types.WithBlockTag(types.BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x941a2dc3ab607819fdc929bea95831a2e0c1aab2f2f34b3a23c55cebc8a040"),
			},
		},
		tests.MainnetEnv: {
			{
				Description: "sierra class",
				blockID:     types.WithBlockTag(types.BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x029927c8af6bccf3f6fda035981e765a7bdbf18a2dc0d630494f8758aa908e2b"),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getClass",
						test.blockID,
						test.ClassHash,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(types.BlockID)
						classHash := args[1].(*felt.Felt)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if classHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    28,
								Message: "Class hash not found",
							}
						}

						var class json.RawMessage
						if classHash.String() == "0x123" {
							// deprecated class
							class = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"testData/class/0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9.json",
								"result",
							)
						}
						if classHash.String() == "0x456" {
							// sierra class
							class = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"testData/class/0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855.json",
								"result",
							)
						}

						*rawResp = class

						return nil
					}).
					Times(1)
			}

			resp, err := Class(
				t.Context(),
				testConfig.Provider,
				test.blockID,
				test.ClassHash,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedClass := testConfig.RPCSpy.LastResponse()

			switch class := resp.(type) {
			case *contracts.DeprecatedContractClass:
				rawClass, err := json.Marshal(class)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedClass), string(rawClass))
			case *contracts.ContractClass:
				rawClass, err := json.Marshal(class)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedClass), string(rawClass))
			default:
				t.Fatalf("Received unknown response type: %T", resp)
			}
		})
	}
}

// TestClassAt tests the ClassAt function.
func TestClassAt(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.TestnetEnv,
		tests.MainnetEnv,
		tests.IntegrationEnv,
	)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description     string
		ContractAddress *felt.Felt
		Block           types.BlockID
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "deprecated class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x456"),
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x789"),
				Block:           types.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "deprecated class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x073ad76dCF68168cBF68EA3EC0382a3605F3dEAf24dc076C355e275769b3c561"),
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:     "invalid contract",
				ContractAddress: internalUtils.DeadBeef,
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				Block:           types.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x004b3d247e79c58e77c93e2c52025d0bb1727957cc9c33b33f7216f369c77be5"),
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getClassAt",
						test.Block,
						test.ContractAddress,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(types.BlockID)
						contractAddress := args[1].(*felt.Felt)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if contractAddress == internalUtils.DeadBeef {
							return RPCError{
								Code:    20,
								Message: "Contract not found",
							}
						}

						var class json.RawMessage
						if contractAddress.String() == "0x123" {
							// deprecated class
							class = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"testData/class/0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9.json",
								"result",
							)
						}
						if contractAddress.String() == "0x456" {
							// sierra class
							class = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"testData/class/0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855.json",
								"result",
							)
						}

						*rawResp = class

						return nil
					}).
					Times(1)
			}

			resp, err := ClassAt(
				t.Context(),
				testConfig.Provider,
				test.Block,
				test.ContractAddress,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.RPCSpy.LastResponse()

			switch class := resp.(type) {
			case *contracts.DeprecatedContractClass:
				rawClass, err := json.Marshal(class)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedResp), string(rawClass))
			case *contracts.ContractClass:
				rawClass, err := json.Marshal(class)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedResp), string(rawClass))
			default:
				t.Fatalf("Received unknown response type: %T", resp)
			}
		},
		)
	}
}

// TestClassHashAt tests the ClassHashAt function.
func TestClassHashAt(t *testing.T) {
	tests.RunTestOn(
		t,
		tests.MockEnv,
		tests.DevnetEnv,
		tests.TestnetEnv,
		tests.MainnetEnv,
		tests.IntegrationEnv,
	)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description     string
		Block           types.BlockID
		ContractAddress *felt.Felt
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "normal call",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
			},
			{
				Description:     "invalid contract",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.DeadBeef,
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				Block:           types.WithBlockHash(internalUtils.DeadBeef),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.DevnetEnv: {
			{
				Description:     "normal call",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x41A78E741E5AF2FEC34B695679BC6891742439F7AFB8484ECD7766661AD02BF"),
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "normal call",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
			},
			{
				Description:     "invalid contract",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.DeadBeef,
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				Block:           types.WithBlockHash(internalUtils.DeadBeef),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "normal call",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "normal call",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa"),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getClassHashAt",
						test.Block,
						test.ContractAddress,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(types.BlockID)
						contractAddress := args[1].(*felt.Felt)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if contractAddress == internalUtils.DeadBeef {
							return RPCError{
								Code:    20,
								Message: "Contract not found",
							}
						}

						*rawResp = json.RawMessage("\"0xdeadbeef\"")

						return nil
					}).
					Times(1)
			}

			classhash, err := ClassHashAt(
				t.Context(),
				testConfig.Provider,
				test.Block,
				test.ContractAddress,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedClassHash := testConfig.RPCSpy.LastResponse()
			rawClassHash, err := json.Marshal(classhash)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedClassHash), string(rawClassHash))
		})
	}
}

// TestCompiledCasm tests the CompiledCasm function.
func TestCompiledCasm(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)
	t.Parallel()

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description   string
		ClassHash     *felt.Felt
		ExpectedError *RPCError
	}

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
			result, err := CompiledCasm(
				t.Context(),
				testConfig.Provider,
				test.ClassHash,
			)
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
	}
}
