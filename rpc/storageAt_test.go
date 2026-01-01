package rpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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

			rawExpectedValue := testConfig.RPCSpy.LastResponse()
			rawValue, err := json.Marshal(value)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedValue), string(rawValue))
		})
	}
}
