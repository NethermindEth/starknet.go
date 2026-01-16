package methods

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	. "github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestNonce tests the Nonce function.
func TestNonce(t *testing.T) {
	tests.RunTestOn(
		t,
		tests.MockEnv,
		tests.TestnetEnv,
		tests.MainnetEnv,
		tests.IntegrationEnv,
	)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description     string
		ContractAddress *felt.Felt
		Block           rpcv10.BlockID
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "normal call",
				Block:           rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
			},
			{
				Description:     "invalid contract",
				Block:           rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ContractAddress: internalUtils.DeadBeef,
				ExpectedError:   ErrContractNotFound,
			},
			{
				Description:     "invalid block",
				Block:           rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ContractAddress: internalUtils.TestHexToFelt(t, "0x123"),
				ExpectedError:   ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				Block:           rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
			},
			{
				Description:     "invalid block",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				Block:           rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "invalid contract address",
				ContractAddress: internalUtils.DeadBeef,
				Block:           rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedError:   ErrContractNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0567f76279d525c7d02057465dd492526b291f864484f3e9c1371c0f770acf0c"),
				Block:           rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
			},
		},
		tests.MainnetEnv: {
			{
				Description:     "normal call",
				ContractAddress: internalUtils.TestHexToFelt(t, "0x00bE9AeF00Ec751Ba252A595A473315FBB8DA629850e13b8dB83d0fACC44E4f2"),
				Block:           rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
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
						blockID := args[0].(rpcv10.BlockID)
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

			nonce, err := Nonce(
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

			rawExpectedNonce := testConfig.RPCSpy.LastResponse()
			rawNonce, err := json.Marshal(nonce)
			require.NoError(t, err)
			assert.Equal(t, string(rawExpectedNonce), string(rawNonce))
		})
	}
}
