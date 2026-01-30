package rpcv10

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestStateUpdate is a test function for the StateUpdate method.
func TestStateUpdate(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := internal.BeforeEach(t, false)
	provider := testConfig.Provider

	type testSetType struct {
		blockID     types.BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				blockID: types.WithBlockTag(types.BlockTagLatest),
			},
			{
				blockID: types.WithBlockTag(types.BlockTagPreConfirmed),
			},
			{
				blockID:     types.WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				blockID:     types.WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				blockID:     types.WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				blockID:     types.WithBlockNumber(99999999999999999),
				ExpectedErr: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	if tests.TEST_ENV != tests.MockEnv {
		// add the common block IDs to the test set of network tests
		blockIDs := GetCommonBlockIDs(t, provider)
		for _, blockID := range blockIDs {
			testSet = append(testSet, testSetType{
				blockID: blockID,
			})
		}
	}

	for _, test := range testSet {
		blockID, _ := test.blockID.MarshalJSON()
		t.Run(fmt.Sprintf("BlockID: %v", string(blockID)), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				blockSepolia3100000 := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/stateUpdate/sepolia3100000.json", "result",
				)

				blockSepoliaPreConfirmed := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
					t,
					"./testData/stateUpdate/sepoliaPreConfirmed.json",
					"result",
				)

				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getStateUpdate",
						test.blockID,
					).
					DoAndReturn(
						func(_, result, _ any, args ...any) error {
							rawResp := result.(*json.RawMessage)
							blockID := args[0].(types.BlockID)

							switch blockID.Tag {
							case types.BlockTagPreConfirmed:
								*rawResp = blockSepoliaPreConfirmed
							case types.BlockTagLatest:
								*rawResp = blockSepolia3100000
							}

							if blockID.Number != nil && *blockID.Number == 99999999999999999 {
								return RPCError{
									Code:    24,
									Message: "Block not found",
								}
							}

							return nil
						},
					).
					Times(1)
			}

			stateUpdate, err := GetStateUpdate(
				t.Context(),
				provider,
				test.blockID,
			)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedStateUpdate := testConfig.RPCSpy.LastResponse()

			rawStateUpdate, err := json.Marshal(stateUpdate)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedStateUpdate), string(rawStateUpdate))
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

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description     string
		ContractAddress *felt.Felt
		StorageKey      string
		Block           types.BlockID
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				StorageKey:      "_signer",
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				StorageKey:      "_signer",
				Block:           types.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				StorageKey:      "_signer",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.DevnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				StorageKey:      "ERC20_name",
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				StorageKey:      "_signer",
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				StorageKey:      "_signer",
				Block:           types.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				StorageKey:      "_signer",
				Block:           types.WithBlockTag(types.BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				StorageKey:      "ERC20_decimals",
				Block:           types.WithBlockTag(types.BlockTagLatest),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x8d17e6a3B92a2b5Fa21B8e7B5a3A794B05e06C5FD6C6451C6F2695Ba77101"),
				StorageKey:      "_signer",
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
						"starknet_getStorageAt",
						test.ContractAddress,
						// the StorateAt function is not compliant with the spec
						fmt.Sprintf("0x%x", internalUtils.GetSelectorFromName(test.StorageKey)),
						test.Block,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						contractAddress := args[0].(*felt.Felt)
						blockID := args[2].(types.BlockID)

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

			value, err := StorageAt(
				t.Context(),
				testConfig.Provider,
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

			rawExpectedValue := testConfig.RPCSpy.LastResponse()
			rawValue, err := json.Marshal(value)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedValue), string(rawValue))
		})
	}
}

// TestStorageProof tests the StorageProof function.
func TestStorageProof(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := internal.BeforeEach(t, false)

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
					blockID: types.WithBlockTag(types.BlockTagLatest),
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
							StorageKeys: []types.StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []types.StorageKey{
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
					blockID: types.WithBlockTag(types.BlockTagPreConfirmed),
				},
				ExpectedError: types.ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockHash(internalUtils.DeadBeef),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockNumber(123456),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call, only required field block_id with 'latest' tag",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
				},
			},
			{
				Description: "block_id + class_hashes parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
					},
				},
			},
			{
				Description: "block_id + contract_addresses parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
				},
			},
			{
				Description: "block_id + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []types.StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
							},
						},
					},
				},
			},
			{
				Description: "block_id + class_hashes + contract_addresses + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
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
							StorageKeys: []types.StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []types.StorageKey{
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
					blockID: types.WithBlockTag(types.BlockTagPreConfirmed),
				},
				ExpectedError: types.ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockHash(internalUtils.DeadBeef),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockNumber(123456),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
		tests.IntegrationEnv: {
			{
				Description: "normal call, only required field block_id with 'latest' tag",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
				},
			},
			{
				Description: "block_id + class_hashes parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
					},
				},
			},
			{
				Description: "block_id + contract_addresses parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
				},
			},
			{
				Description: "block_id + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []types.StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
							},
						},
					},
				},
			},
			{
				Description: "block_id + class_hashes + contract_addresses + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockTag(types.BlockTagLatest),
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
							StorageKeys: []types.StorageKey{
								"0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
								"0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []types.StorageKey{
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
					blockID: types.WithBlockTag(types.BlockTagPreConfirmed),
				},
				ExpectedError: types.ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockHash(internalUtils.DeadBeef),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					blockID: types.WithBlockNumber(123456),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv &&
				test.StorageProofInput.blockID.Tag != types.BlockTagPreConfirmed {
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
						blockID := storageProofInput.blockID

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

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/storageProof/sepoliaLatestFullResp.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			result, err := StorageProof(
				t.Context(),
				testConfig.Provider,
				test.StorageProofInput,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			// verify JSON equality
			rawResult := testConfig.RPCSpy.LastResponse()
			marshalledResult, err := json.Marshal(result)
			require.NoError(t, err)

			assert.JSONEq(t, string(rawResult), string(marshalledResult))
		})
	}
}
