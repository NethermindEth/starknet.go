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

			rawExpectedClassHash := testConfig.RPCSpy.LastResponse()
			rawClassHash, err := json.Marshal(classhash)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedClassHash), string(rawClassHash))
		})
	}
}
