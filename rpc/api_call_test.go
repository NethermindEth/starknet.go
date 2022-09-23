package rpc

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestCall tests Call
func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FunctionCall   types.FunctionCall
		BlockID        types.BlockID
		ExpectedResult string
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				FunctionCall: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x035a55a64238b776664d7723de1f6b50350116a1ab1ca1fe154320a0eba53d3a"),
					EntryPointSelector: "get_count",
					CallData:           []string{},
				},
				BlockID:        WithBlockTag("latest"),
				ExpectedResult: "0x01",
			},
		},
		"mock": {
			{
				FunctionCall: types.FunctionCall{
					ContractAddress:    types.HexToHash("0xdeadbeef"),
					EntryPointSelector: "decimals",
					CallData:           []string{},
				},
				BlockID:        WithBlockTag("latest"),
				ExpectedResult: "0x12",
			},
		},
		"testnet": {
			{
				FunctionCall: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x029260ce936efafa6d0042bc59757a653e3f992b97960c1c4f8ccd63b7a90136"),
					EntryPointSelector: "decimals",
					CallData:           []string{},
				},
				BlockID:        WithBlockTag("latest"),
				ExpectedResult: "0x12",
			},
			{
				FunctionCall: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					EntryPointSelector: "balanceOf",
					CallData:           []string{"0x0207aCC15dc241e7d167E67e30E769719A727d3E0fa47f9E187707289885Dfde"},
				},
				BlockID:        WithBlockNumber(310000),
				ExpectedResult: "0x2f0e64b37383fa",
			},
		},
		"mainnet": {
			{
				FunctionCall: types.FunctionCall{
					ContractAddress:    types.HexToHash("0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75"),
					EntryPointSelector: "decimals",
					CallData:           []string{},
				},
				BlockID:        WithBlockTag("latest"),
				ExpectedResult: "0x12",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		function := test.FunctionCall
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		output, err := testConfig.client.Call(context.Background(), function, test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		if diff, err := spy.Compare(output, false); err != nil || diff != "FullMatch" {
			spy.Compare(output, true)
			t.Fatal("expecting to match", err)
		}
		if len(output) == 0 {
			t.Fatal("should return an output")
		}
		if output[0] != test.ExpectedResult {
			t.Fatalf("1st output expecting %s, got: %s", test.ExpectedResult, output[0])
		}
	}
}
