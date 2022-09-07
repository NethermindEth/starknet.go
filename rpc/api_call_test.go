package rpc

import (
	"context"
	"fmt"
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
			t.Fatalf("1st output expecting %s,git %s", test.ExpectedResult, output[0])
		}
	}
}

// TestEstimateFee tests EstimateFee
func TestEstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Call                types.BroadcastedInvokeTxn
		BlockID             types.BlockID
		ExpectedOverallFee  string
		ExpectedGasPrice    string
		ExpectedGasConsumed string
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {
			// 	TxnHash: "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99",
			// {
			// 	Call: BroadcastedInvokeTxnV0{
			// 		BroadcastedCommonTxnProperties: BroadcastedCommonTxnProperties{
			// 			// Type:    "INVOKE",
			// 			MaxFee:  "0xde0b6b3a7640000",
			// 			Version: "0x0",
			// 			Signature: []string{
			// 				"0x7bc0a22005a54ec6a005c1e89ab0201cbd0819621edd9fe4d5ef177a4ff33dd",
			// 				"0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea",
			// 			},
			// 			// Nonce: "0x0",
			// 		},
			// 		InvokeV0: InvokeV0(FunctionCall{
			// 			ContractAddress:    "0x2e28403d7ee5e337b7d456327433f003aa875c29631906908900058c83d8cb6",
			// 			EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
			// 			CallData: []string{
			// 				"0x1",
			// 				"0x33830ce413e4c096eef81b5e6ffa9b9f5d963f57b8cd63c9ae4c839c383c1a6",
			// 				"0x2db698626ed7f60212e1ce6e99afb796b6b423d239c3f0ecef23e840685e866",
			// 				"0x0",
			// 				"0x2",
			// 				"0x2",
			// 				"0x61c6e7484657e5dc8b21677ffa33e4406c0600bba06d12cf1048fdaa55bdbc3",
			// 				"0x6307b990",
			// 				"0x2b81",
			// 				"0x0",
			// 			},
			// 		}),
			// 	},
			// 	BlockID: WithBlockTag("latest"),
			// WithBlockHash("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
			// },
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		call := test.Call
		output, err := testConfig.client.EstimateFee(context.Background(), call, test.BlockID)
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
