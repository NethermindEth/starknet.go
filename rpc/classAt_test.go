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

// TestClassAt tests the ClassAt function.
func TestClassAt(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.TestnetEnv,
		tests.MainnetEnv,
		tests.IntegrationEnv,
	)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description     string
		ContractAddress *felt.Felt
		Block           BlockID
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "deprecated class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x456"),
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x789"),
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				Block:           WithBlockTag(BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "deprecated class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x073ad76dCF68168cBF68EA3EC0382a3605F3dEAf24dc076C355e275769b3c561"),
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "invalid contract",
				ContractAddress: internalUtils.DeadBeef,
				Block:           WithBlockTag(BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				Block:           WithBlockTag(BlockTagLatest),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "sierra class",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x004b3d247e79c58e77c93e2c52025d0bb1727957cc9c33b33f7216f369c77be5"),
				Block:           WithBlockTag(BlockTagLatest),
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
						blockID := args[0].(BlockID)
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

			resp, err := testConfig.Provider.ClassAt(
				t.Context(),
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
