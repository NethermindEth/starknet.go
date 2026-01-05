package methods

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestCall tests the Call function.
func TestCall(t *testing.T) {
	tests.RunTestOn(
		t,
		tests.MockEnv,
		tests.TestnetEnv,
		tests.MainnetEnv,
		tests.DevnetEnv,
		tests.IntegrationEnv,
	)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		name                  string
		FunctionCall          rpcv10.FunctionCall
		BlockID               rpcv10.BlockID
		ExpectedPatternResult *felt.Felt
		ExpectedError         *RPCError
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.DevnetEnv: {
			{
				name: "Ok",
				FunctionCall: rpcv10.FunctionCall{
					// ContractAddress of predeployed devnet Feetoken
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x12"),
			},
		},
		tests.MockEnv: {
			{
				name: "Ok",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.DeadBeef,
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedPatternResult: internalUtils.DeadBeef,
			},
		},
		tests.TestnetEnv: {
			{
				name: "Ok - latest block tag",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x506f736974696f6e"),
			},
			{
				name: "Ok - pre_confirmed block tag",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagPreConfirmed),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x506f736974696f6e"),
			},
			{
				name: "Ok - l1_accepted block tag",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagL1Accepted),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x506f736974696f6e"),
			},
			{
				name: "ContractError",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{&felt.Zero},
				},
				BlockID:       rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedError: ErrContractError,
			},
			{
				name: "EntrypointNotFound",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("RANDOM_STRINGGG"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:       rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedError: ErrEntrypointNotFound,
			},
			{
				name: "BlockNotFound",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:       rpcv10.WithBlockNumber(999999999999999),
				ExpectedError: ErrBlockNotFound,
			},
			{
				name: "ContractNotFound",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.DeadBeef,
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:       rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedError: ErrContractNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				name: "Ok",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x12"),
			},
		},
		tests.MainnetEnv: {
			{
				name: "Ok",
				FunctionCall: rpcv10.FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x12"),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run("Test: "+test.name, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_call",
						test.FunctionCall,
						test.BlockID,
					).
					DoAndReturn(
						func(_, result, _ any, _ ...any) error {
							rawResp := result.(*json.RawMessage)
							*rawResp = json.RawMessage("[\"0xdeadbeef\"]")

							return nil
						},
					).
					Times(1)
			}
			output, err := Call(
				t.Context(),
				testConfig.Provider.c,
				test.FunctionCall,
				test.BlockID,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.ExpectedError.Message)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, output, "should return an output")
				assert.Equal(t, test.ExpectedPatternResult, output[0])
			}
		})
	}
}
