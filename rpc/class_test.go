package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestClass tests the Class function.
func TestClass(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description   string
		BlockID       BlockID
		ClassHash     *felt.Felt
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "deprecated class",
				ClassHash:   internalUtils.TestHexToFelt(t, "0x123"),
				BlockID:     WithBlockTag(BlockTagLatest),
			},
			{
				Description: "sierra class",
				ClassHash:   internalUtils.TestHexToFelt(t, "0x456"),
				BlockID:     WithBlockTag(BlockTagLatest),
			},
			{
				Description:   "invalid block",
				ClassHash:     internalUtils.TestHexToFelt(t, "0x789"),
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description:   "invalid class hash",
				ClassHash:     internalUtils.DeadBeef,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrClassHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "deprecated class",
				BlockID:     WithBlockTag(BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
			},
			{
				Description: "sierra class",
				BlockID:     WithBlockTag(BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
			},
			{
				Description:   "invalid block",
				ClassHash:     internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description:   "invalid class hash",
				ClassHash:     internalUtils.DeadBeef,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrClassHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description: "sierra class",
				BlockID:     WithBlockTag(BlockTagLatest),
				ClassHash:   internalUtils.TestHexToFelt(t, "0x941a2dc3ab607819fdc929bea95831a2e0c1aab2f2f34b3a23c55cebc8a040"),
			},
		},
		tests.MainnetEnv: {
			{
				Description: "sierra class",
				BlockID:     WithBlockTag(BlockTagLatest),
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
						test.BlockID,
						test.ClassHash,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(BlockID)
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

			resp, err := testConfig.Provider.Class(
				t.Context(),
				test.BlockID,
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
