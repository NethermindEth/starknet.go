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

// TestEstimateFee tests EstimateFee
func TestEstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		call                types.FunctionInvoke
		BlockHashOrTag      string
		ExpectedOverallFee  string
		ExpectedGasPrice    string
		ExpectedGasConsumed string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				call: types.FunctionInvoke{
					FunctionCall: types.FunctionCall{
						ContractAddress: "0x0019fcae2482de8fb3afaf8d4b219449bec93a5928f02f58eef645cc071767f4",
						Calldata: []string{
							"0x0000000000000000000000000000000000000000000000000000000000000001",
							"0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
							"0x0083afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
							"0x0000000000000000000000000000000000000000000000000000000000000000",
							"0x0000000000000000000000000000000000000000000000000000000000000003",
							"0x0000000000000000000000000000000000000000000000000000000000000003",
							"0x04681402a7ab16c41f7e5d091f32fe9b78de096e0bd5962ce5bd7aaa4a441f64",
							"0x000000000000000000000000000000000000000000000000001d41f6331e6800",
							"0x0000000000000000000000000000000000000000000000000000000000000000",
							"0x0000000000000000000000000000000000000000000000000000000000000001",
						},
						EntryPointSelector: "0x015d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
					},
				},
				BlockHashOrTag:      "0x0147c4b0f702079384e26d9d34a15e7758881e32b219fc68c076b09d0be13f8c",
				ExpectedOverallFee:  "0x7134",
				ExpectedGasPrice:    "0x45",
				ExpectedGasConsumed: "0x1a4",
			},
		},
		"testnet": {},
		"mainnet": {
			{
				call: types.FunctionInvoke{
					FunctionCall: types.FunctionCall{
						ContractAddress: "0x0019fcae2482de8fb3afaf8d4b219449bec93a5928f02f58eef645cc071767f4",
						Calldata: []string{
							"0x0000000000000000000000000000000000000000000000000000000000000001",
							"0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
							"0x0083afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
							"0x0000000000000000000000000000000000000000000000000000000000000000",
							"0x0000000000000000000000000000000000000000000000000000000000000003",
							"0x0000000000000000000000000000000000000000000000000000000000000003",
							"0x04681402a7ab16c41f7e5d091f32fe9b78de096e0bd5962ce5bd7aaa4a441f64",
							"0x000000000000000000000000000000000000000000000000001d41f6331e6800",
							"0x0000000000000000000000000000000000000000000000000000000000000000",
							"0x0000000000000000000000000000000000000000000000000000000000000001",
						},
						EntryPointSelector: "0x015d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
					},
				},
				BlockHashOrTag:      "0x0147c4b0f702079384e26d9d34a15e7758881e32b219fc68c076b09d0be13f8c",
				ExpectedOverallFee:  "0xc84c599f51bd",
				ExpectedGasPrice:    "0x5df32828e",
				ExpectedGasConsumed: "0x2178",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		call := test.call
		output, err := testConfig.client.EstimateFee(context.Background(), call, test.BlockHashOrTag)
		if err != nil || output == nil {
			t.Fatalf("output is nil, go err %v", err)
		}
		if fmt.Sprintf("0x%x", output.OverallFee) != test.ExpectedOverallFee {
			t.Fatalf("expected %s, got %s", test.ExpectedOverallFee, fmt.Sprintf("0x%x", output.OverallFee))
		}
		if fmt.Sprintf("0x%x", output.GasConsumed) != test.ExpectedGasConsumed {
			t.Fatalf("expected %s, got %s", test.ExpectedGasConsumed, fmt.Sprintf("0x%x", output.GasConsumed))
		}
		if fmt.Sprintf("0x%x", output.GasPrice) != test.ExpectedGasPrice {
			t.Fatalf("expected %s, got %s", test.ExpectedGasPrice, fmt.Sprintf("0x%x", output.GasPrice))
		}
	}
}
