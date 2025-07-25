package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/require"
)

// TestCall tests the Call function.
//
// It sets up different test scenarios and asserts the result of each scenario.
// The test scenarios include different contract addresses, entry point selectors,
// and expected results for different environments (devnet, mock, testnet, mainnet).
// The function calls the provider, and compares the output against the expected output.
// It also checks that the output is not empty and that it matches the expected pattern result.
// If any of the assertions fail, the test fails with an error.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestCall(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.DevnetEnv)

	testConfig := beforeEach(t, false)

	type testSetType struct {
		name                  string
		FunctionCall          FunctionCall
		BlockID               BlockID
		ExpectedPatternResult *felt.Felt
		ExpectedError         *RPCError
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.DevnetEnv: {
			{
				name: "Ok",
				FunctionCall: FunctionCall{
					// ContractAddress of predeployed devnet Feetoken
					ContractAddress:    internalUtils.TestHexToFelt(t, DevNetETHAddress),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x12"),
			},
		},
		tests.MockEnv: {
			{
				name: "Ok",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0xdeadbeef"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		tests.TestnetEnv: {
			{
				name: "Ok",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x506f736974696f6e"),
			},
			{
				name: "ContractError",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{&felt.Zero},
				},
				BlockID:       WithBlockTag("latest"),
				ExpectedError: ErrContractError,
			},
			{
				name: "EntrypointNotFound",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("RANDOM_STRINGGG"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:       WithBlockTag("latest"),
				ExpectedError: ErrEntrypointNotFound,
			},
			{
				name: "BlockNotFound",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:       WithBlockNumber(9999999999999999999),
				ExpectedError: ErrBlockNotFound,
			},
			{
				name: "ContractNotFound",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.RANDOM_FELT,
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:       WithBlockTag("latest"),
				ExpectedError: ErrContractNotFound,
			},
		},
		tests.MainnetEnv: {
			{
				name: "Ok",
				FunctionCall: FunctionCall{
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: internalUtils.TestHexToFelt(t, "0x12"),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(fmt.Sprintf("Network: %s, Test: %s", tests.TEST_ENV, test.name), func(t *testing.T) {
			output, err := testConfig.provider.Call(context.Background(), test.FunctionCall, test.BlockID)
			if err != nil {
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				require.ErrorContains(t, test.ExpectedError, rpcErr.Message)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, output, "should return an output")
				require.Equal(t, test.ExpectedPatternResult, output[0])
			}
		})
	}
}
