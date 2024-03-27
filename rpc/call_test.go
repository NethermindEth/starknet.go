package rpc

import (
	"context"
	"log"
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
// The function uses a spy to monitor the calls made to the provider, and compares
// the output against the expected output. It also checks that the output is not empty
// and that it matches the expected pattern result. If any of the assertions fail,
// the test fails with an error.
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
					ContractAddress:    utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("name"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: utils.TestHexToFelt(t, "0x6574686572"),
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
					ContractAddress:    utils.TestHexToFelt(t, "0x029260ce936efafa6d0042bc59757a653e3f992b97960c1c4f8ccd63b7a90136"),
					EntryPointSelector: utils.TestHexToFelt(t, "0x004c4fb1ab068f6039d5780c68dd0fa2f8742cceb3426d19667778ca7f3518a9"),
					Calldata:           []*felt.Felt{},
				},
				BlockID:               WithBlockTag("latest"),
				ExpectedPatternResult: utils.TestHexToFelt(t, "0x12"),
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
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		output, err := testConfig.provider.Call(context.Background(), FunctionCall(test.FunctionCall), test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		if diff, err := spy.Compare(output, false); err != nil || diff != "FullMatch" {
			if _, err := spy.Compare(output, true); err != nil {
				log.Fatal(err)
			}
			t.Fatal("expecting to match", err)
		}
		if len(output) == 0 {
			t.Fatal("should return an output")
		}
		require.Equal(t, test.ExpectedPatternResult, output[0])
	}
}
