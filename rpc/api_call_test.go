package rpc

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

// TestCall tests Call
func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress    string
		EntrypointSelector string
		ExpectedResult     string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress:    "0xdeadbeef",
				EntrypointSelector: "decimals",
				ExpectedResult:     "0x12",
			},
		},
		"testnet": {
			{
				ContractAddress:    "0x029260ce936efafa6d0042bc59757a653e3f992b97960c1c4f8ccd63b7a90136",
				EntrypointSelector: "decimals",
				ExpectedResult:     "0x12",
			},
		},
		"mainnet": {
			{
				ContractAddress:    "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75",
				EntrypointSelector: "decimals",
				ExpectedResult:     "0x12",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		function := types.FunctionCall{
			ContractAddress:    test.ContractAddress,
			EntryPointSelector: test.EntrypointSelector,
		}
		output, err := testConfig.client.Call(context.Background(), function, "latest")
		if err != nil {
			t.Fatal(err)
		}
		if len(output) == 0 {
			t.Fatal("should return an output")
		}
		if output[0] != test.ExpectedResult {
			t.Fatalf("1st output expecting %s,git %s", test.ExpectedResult, output[0])
		}
	}
}
