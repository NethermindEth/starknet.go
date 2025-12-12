package rpc

import (
	"encoding/json"
	"fmt"
	"reflect"
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
							class = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"testData/class/0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9.json",
								"result",
							)
						}
						if contractAddress.String() == "0x456" {
							// sierra class
							class = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
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

			rawExpectedResp := testConfig.Spy.LastResponse()

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
				t.Fatalf("Received unknown response type: %v", reflect.TypeOf(resp))
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

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description     string
		Block           BlockID
		ContractAddress *felt.Felt
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "normal call",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
			},
			{
				Description:     "invalid contract",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.DeadBeef,
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.DevnetEnv: {
			{
				Description:     "normal call",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x41A78E741E5AF2FEC34B695679BC6891742439F7AFB8484ECD7766661AD02BF"),
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "normal call",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
			},
			{
				Description:     "invalid contract",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.DeadBeef,
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "normal call",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "normal call",
				Block:           WithBlockTag(BlockTagLatest),
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

						*rawResp = json.RawMessage("\"0xdeadbeef\"")

						return nil
					}).
					Times(1)
			}

			classhash, err := testConfig.Provider.ClassHashAt(
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

			rawExpectedClassHash := testConfig.Spy.LastResponse()
			rawClassHash, err := json.Marshal(classhash)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedClassHash), string(rawClassHash))
		})
	}
}

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
							class = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"testData/class/0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9.json",
								"result",
							)
						}
						if classHash.String() == "0x456" {
							// sierra class
							class = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
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

			rawExpectedClass := testConfig.Spy.LastResponse()

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
				t.Fatalf("Received unknown response type: %v", reflect.TypeOf(resp))
			}
		})
	}
}

// TestStorageAt tests the StorageAt function.
func TestStorageAt(t *testing.T) {
	tests.RunTestOn(
		t,
		tests.DevnetEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description     string
		ContractAddress *felt.Felt
		StorageKey      string
		Block           BlockID
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				StorageKey:      "_signer",
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				StorageKey:      "_signer",
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				StorageKey:      "_signer",
				Block:           WithBlockTag(BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.DevnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				StorageKey:      "ERC20_name",
				Block:           WithBlockTag(BlockTagLatest),
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				StorageKey:      "_signer",
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				StorageKey:      "_signer",
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				StorageKey:      "_signer",
				Block:           WithBlockTag(BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				StorageKey:      "ERC20_decimals",
				Block:           WithBlockTag(BlockTagLatest),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x8d17e6a3B92a2b5Fa21B8e7B5a3A794B05e06C5FD6C6451C6F2695Ba77101"),
				StorageKey:      "_signer",
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
						"starknet_getStorageAt",
						test.ContractAddress,
						// the StorateAt function is not compliant with the spec
						fmt.Sprintf("0x%x", internalUtils.GetSelectorFromName(test.StorageKey)),
						test.Block,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						contractAddress := args[0].(*felt.Felt)
						blockID := args[2].(BlockID)

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

			value, err := testConfig.Provider.StorageAt(
				t.Context(),
				test.ContractAddress,
				test.StorageKey,
				test.Block,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedValue := testConfig.Spy.LastResponse()
			rawValue, err := json.Marshal(value)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedValue), string(rawValue))
		})
	}
}

// TestNonce tests the Nonce function.
func TestNonce(t *testing.T) {
	tests.RunTestOn(
		t,
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
				Description:     "normal call",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
			},
			{
				Description:     "invalid contract",
				Block:           WithBlockTag(BlockTagLatest),
				ContractAddress: internalUtils.DeadBeef,
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				Block:           WithBlockHash(internalUtils.DeadBeef),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				Block:           WithBlockTag(BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
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
		tests.IntegrationEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0567f76279d525c7d02057465dd492526b291f864484f3e9c1371c0f770acf0c"),
				Block:           WithBlockTag(BlockTagLatest),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x00bE9AeF00Ec751Ba252A595A473315FBB8DA629850e13b8dB83d0fACC44E4f2"),
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
						"starknet_getNonce",
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

						*rawResp = json.RawMessage("\"0xdeadbeef\"")

						return nil
					}).
					Times(1)
			}

			nonce, err := testConfig.Provider.Nonce(
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

			rawExpectedNonce := testConfig.Spy.LastResponse()
			rawNonce, err := json.Marshal(nonce)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedNonce), string(rawNonce))
		})
	}
}

// TestEstimateMessageFee tests the EstimateMessageFee function.
func TestEstimateMessageFee(t *testing.T) {
	// TODO: add integration testcase
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description string
		MsgFromL1
		BlockID
		ExpectedError *RPCError
	}

	// https://sepolia.voyager.online/message/0x273f4e20fc522098a60099e5872ab3deeb7fb8321a03dadbd866ac90b7268361
	l1Handler := MsgFromL1{
		FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
		ToAddress: internalUtils.TestHexToFelt(
			t,
			"0x04c5772d1914fe6ce891b64eb35bf3522aeae1315647314aac58b01137607f3f",
		),
		Selector: internalUtils.TestHexToFelt(
			t,
			"0x1b64b1b3b690b43b9b514fb81377518f4039cd3e4f4914d8a6bdf01d679fb19",
		),
		Payload: internalUtils.TestHexArrToFelt(t, []string{
			"0x455448",
			"0x2f14d277fc49e0e2d2967d019aea8d6bd9cb3998",
			"0x02000e6213e24b84012b1f4b1cbd2d7a723fb06950aeab37bedb6f098c7e051a",
			"0x01a055690d9db80000",
			"0x00",
		}),
	}

	l1HandlerInvalidSelector := l1Handler
	l1HandlerInvalidSelector.Selector = internalUtils.DeadBeef

	l1HandlerInvalidToAddress := l1Handler
	l1HandlerInvalidToAddress.ToAddress = internalUtils.DeadBeef

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "normal call",
				MsgFromL1:   l1Handler,
				BlockID:     WithBlockTag(BlockTagLatest),
			},
			{
				Description:   "contract error",
				MsgFromL1:     l1HandlerInvalidSelector,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractError,
			},
			{
				Description:   "contract not found",
				MsgFromL1:     l1HandlerInvalidToAddress,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractNotFound,
			},
			{
				Description:   "invalid block",
				MsgFromL1:     l1Handler,
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call",
				MsgFromL1:   l1Handler,
				BlockID:     WithBlockTag(BlockTagLatest),
			},
			{
				Description:   "contract error",
				MsgFromL1:     l1HandlerInvalidSelector,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractError,
			},
			{
				Description:   "contract not found",
				MsgFromL1:     l1HandlerInvalidToAddress,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractNotFound,
			},
			{
				Description:   "invalid block",
				MsgFromL1:     l1Handler,
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
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
						"starknet_estimateMessageFee",
						test.MsgFromL1,
						test.BlockID,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						msgFromL1 := args[0].(MsgFromL1)
						blockID := args[1].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if msgFromL1.ToAddress == internalUtils.DeadBeef {
							return RPCError{
								Code:    20,
								Message: "Contract not found",
							}
						}

						if msgFromL1.Selector == internalUtils.DeadBeef {
							return RPCError{
								Code:    40,
								Message: "Contract error",
								Data:    &ContractErrData{},
							}
						}

						*rawResp = json.RawMessage(`
							{
								"l1_gas_consumed": "0x4ed3",
								"l1_gas_price": "0x7e15d2b5",
								"l2_gas_consumed": "0x0",
								"l2_gas_price": "0x1",
								"l1_data_gas_consumed": "0x80",
								"l1_data_gas_price": "0x1",
								"overall_fee": "0x26d2922fd1af",
								"unit": "WEI"
							}
						`)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.EstimateMessageFee(
				t.Context(),
				test.MsgFromL1,
				test.BlockID,
			)
			if test.ExpectedError != nil {
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

				return
			}
			require.NoError(t, err)
			rawExpectedFeeEst := testConfig.Spy.LastResponse()

			rawFeeEst, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedFeeEst), string(rawFeeEst))
		},
		)
	}
}

// TestEstimateFee tests the EstimateFee function.
func TestEstimateFee(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		description   string
		txs           []BroadcastTxn
		simFlags      []SimulationFlag
		blockID       BlockID
		expectedError *RPCError
	}

	sepoliaInvokeV3 := *internalUtils.TestUnmarshalJSONFileToType[BroadcastInvokeTxnV3](
		t,
		"./testData/transactions/sepoliaInvokeV3_0x6035477af07a1b0a0186bec85287a6f629791b2f34b6e90eec9815c7a964f64.json",
		"",
	)
	invalidSepoliaInvokeV3 := sepoliaInvokeV3
	invalidSepoliaInvokeV3.Calldata = []*felt.Felt{internalUtils.DeadBeef}

	integrationInvokeV3 := *internalUtils.TestUnmarshalJSONFileToType[BroadcastInvokeTxnV3](
		t,
		"./testData/transactions/integrationInvokeV3_0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019.json",
		"",
	)
	invalidIntegrationInvokeV3 := integrationInvokeV3
	invalidIntegrationInvokeV3.Calldata = []*felt.Felt{internalUtils.DeadBeef}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				description: "without flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags: []SimulationFlag{},
				blockID:  WithBlockTag(BlockTagLatest),
			},
			{
				description: "with flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags: []SimulationFlag{SkipValidate},
				blockID:  WithBlockTag(BlockTagLatest),
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					invalidSepoliaInvokeV3,
				},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				blockID:       WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				description: "normal call - without flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
			},
			{
				description: "normal call - two transactions",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
					sepoliaInvokeV3,
				},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
			},
			{
				description: "normal call - with skip validate flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags:      []SimulationFlag{SkipValidate},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					invalidSepoliaInvokeV3,
				},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				blockID:       WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
			// the contract_not_found error will not be tested since it's still not clear
			// when it should be returned (Pathfinder and Juno behave differently)
		},
		tests.IntegrationEnv: {
			{
				description: "without flag",
				txs: []BroadcastTxn{
					integrationInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(1_300_000),
				expectedError: nil,
			},
			{
				description: "with flag",
				txs: []BroadcastTxn{
					integrationInvokeV3,
				},
				simFlags:      []SimulationFlag{SkipValidate},
				blockID:       WithBlockNumber(1_300_000),
				expectedError: nil,
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					invalidIntegrationInvokeV3,
				},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					integrationInvokeV3,
				},
				blockID:       WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_estimateFee",
						test.txs,
						test.simFlags,
						test.blockID,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						txs := args[0].([]BroadcastTxn)
						blockID := args[2].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if txs[0].(BroadcastInvokeTxnV3).Calldata[0] == internalUtils.DeadBeef {
							return RPCError{
								Code:    41,
								Message: "Transaction execution error",
								Data:    &TransactionExecErrData{},
							}
						}

						*rawResp = json.RawMessage(`
							[
								{
									"l1_data_gas_consumed": "0x80",
									"l1_data_gas_price": "0x75da",
									"l1_gas_consumed": "0x0",
									"l1_gas_price": "0x1b709d15a1c6",
									"l2_gas_consumed": "0xc25b1",
									"l2_gas_price": "0xb2d05e00",
									"overall_fee": "0x87c1827e1eb00",
									"unit": "FRI"
								}
							]
						`)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.EstimateFee(
				t.Context(),
				test.txs,
				test.simFlags,
				test.blockID,
			)
			if test.expectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.expectedError.Code, rpcErr.Code)
				assert.Equal(t, test.expectedError.Message, rpcErr.Message)
				assert.IsType(t, rpcErr.Data, rpcErr.Data)

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}

// TestGetStorageProof tests the GetStorageProof function.
func TestGetStorageProof(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description       string
		StorageProofInput StorageProofInput
		ExpectedError     error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "block_id + class_hashes + contract_addresses + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
						internalUtils.TestHexToFelt(t, "0x009524a94b41c4440a16fd96d7c1ef6ad6f44c1c013e96662734502cd4ee9b1f"),
					},
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
					},
				},
			},
			{
				Description: "error: using pre_confirmed tag in block_id",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagPreConfirmed),
				},
				ExpectedError: ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockHash(internalUtils.DeadBeef),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockNumber(123456),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call, only required field block_id with 'latest' tag",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
				},
			},
			{
				Description: "block_id + class_hashes parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
					},
				},
			},
			{
				Description: "block_id + contract_addresses parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
				},
			},
			{
				Description: "block_id + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
							},
						},
					},
				},
			},
			{
				Description: "block_id + class_hashes + contract_addresses + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
						internalUtils.TestHexToFelt(t, "0x009524a94b41c4440a16fd96d7c1ef6ad6f44c1c013e96662734502cd4ee9b1f"),
					},
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
					},
				},
			},
			{
				Description: "error: using pre_confirmed tag in block_id",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagPreConfirmed),
				},
				ExpectedError: ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockHash(internalUtils.DeadBeef),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockNumber(123456),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
		tests.IntegrationEnv: {
			{
				Description: "normal call, only required field block_id with 'latest' tag",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
				},
			},
			{
				Description: "block_id + class_hashes parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
					},
				},
			},
			{
				Description: "block_id + contract_addresses parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
				},
			},
			{
				Description: "block_id + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
							},
						},
					},
				},
			},
			{
				Description: "block_id + class_hashes + contract_addresses + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
						internalUtils.TestHexToFelt(t, "0x009524a94b41c4440a16fd96d7c1ef6ad6f44c1c013e96662734502cd4ee9b1f"),
					},
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
					},
				},
			},
			{
				Description: "error: using pre_confirmed tag in block_id",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockTag(BlockTagPreConfirmed),
				},
				ExpectedError: ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockHash(internalUtils.DeadBeef),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					BlockID: WithBlockNumber(123456),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv &&
				test.StorageProofInput.BlockID.Tag != BlockTagPreConfirmed {
				testConfig.MockClient.EXPECT().
					CallContext(
						t.Context(),
						gomock.Any(),
						"starknet_getStorageProof",
						test.StorageProofInput,
					).
					DoAndReturn(func(_, result, _ any, arg any) error {
						rawResp := result.(*json.RawMessage)
						storageProofInput := arg.(StorageProofInput)
						blockID := storageProofInput.BlockID

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if blockID.Number != nil && *blockID.Number < 3000000 {
							return RPCError{
								Code:    42,
								Message: "the node doesn't support storage proofs for blocks that are too far in the past",
							}
						}

						*rawResp = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/storageProof/sepoliaLatestFullResp.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			result, err := testConfig.Provider.StorageProof(t.Context(), test.StorageProofInput)
			if test.ExpectedError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			// verify JSON equality
			rawResult := testConfig.Spy.LastResponse()
			marshalledResult, err := json.Marshal(result)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawResult), string(marshalledResult))
		})
	}
}
