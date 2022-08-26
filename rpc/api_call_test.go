package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

// TestCall tests Call
func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FunctionCall   FunctionCall
		BlockIDOption  BlockIDOption
		ExpectedResult string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    "0xdeadbeef",
					EntryPointSelector: "decimals",
					CallData:           []string{},
				},
				BlockIDOption:  WithBlockIDTag("latest"),
				ExpectedResult: "0x12",
			},
		},
		"testnet": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    "0x029260ce936efafa6d0042bc59757a653e3f992b97960c1c4f8ccd63b7a90136",
					EntryPointSelector: "decimals",
					CallData:           []string{},
				},
				BlockIDOption:  WithBlockIDTag("latest"),
				ExpectedResult: "0x12",
			},
			{
				FunctionCall: FunctionCall{
					ContractAddress:    "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					EntryPointSelector: "balanceOf",
					CallData:           []string{"0x0207aCC15dc241e7d167E67e30E769719A727d3E0fa47f9E187707289885Dfde"},
				},
				BlockIDOption:  WithBlockIDNumber(310000),
				ExpectedResult: "0x2f0e64b37383fa",
			},
		},
		"mainnet": {
			{
				FunctionCall: FunctionCall{
					ContractAddress:    "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75",
					EntryPointSelector: "decimals",
					CallData:           []string{},
				},
				BlockIDOption:  WithBlockIDTag("latest"),
				ExpectedResult: "0x12",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		function := test.FunctionCall
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		output, err := testConfig.client.Call(context.Background(), function, test.BlockIDOption)
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
			t.Fatalf("1st output expecting %s,git %s", test.ExpectedResult, output[0])
		}
	}
}

// TestEstimateFee tests EstimateFee
func TestEstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		call                types.FunctionInvoke
		BlockIDOption       BlockIDOption
		ExpectedOverallFee  string
		ExpectedGasPrice    string
		ExpectedGasConsumed string
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		call := test.call
		output, err := testConfig.client.EstimateFee(context.Background(), call, test.BlockIDOption)
		if err != nil || output == nil {
			t.Fatalf("output is nil, go err %v", err)
		}
		if string(output.OverallFee) != test.ExpectedOverallFee {
			t.Fatalf("expected %s, got %s", test.ExpectedOverallFee, output.OverallFee)
		}
		if fmt.Sprintf("0x%x", output.GasConsumed) != test.ExpectedGasConsumed {
			t.Fatalf("expected %s, got %s", test.ExpectedGasConsumed, fmt.Sprintf("0x%x", output.GasConsumed))
		}
		if fmt.Sprintf("0x%x", output.GasPrice) != test.ExpectedGasPrice {
			t.Fatalf("expected %s, got %s", test.ExpectedGasPrice, fmt.Sprintf("0x%x", output.GasPrice))
		}
	}
}
