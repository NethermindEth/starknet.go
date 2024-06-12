package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
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
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FunctionCall          FunctionCall
		BlockID               BlockID
		ExpectedPatternResult *felt.Felt
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				FunctionCall: FunctionCall{
					// ContractAddress of predeployed devnet Feetoken
					ContractAddress:    utils.TestHexToFelt(t, DevNetETHAddress),
					EntryPointSelector: utils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: utils.TestHexToFelt(t, "0x12"),
			},
		},
		"mock": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0xdeadbeef"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: utils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		"testnet": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x025633c6142D9CA4126e3fD1D522Faa6e9f745144aba728c0B3FEE38170DF9e7"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: utils.TestHexToFelt(t, "0x506f736974696f6e"),
			},
		},
		"mainnet": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("decimals"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: utils.TestHexToFelt(t, "0x12"),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		output, err := testConfig.provider.Call(context.Background(), FunctionCall(test.FunctionCall), test.BlockID)
		require.NoError(err)
		require.NotEmpty(output, "should return an output")
		require.Equal(test.ExpectedPatternResult, output[0])
	}
}
