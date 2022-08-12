package rpc

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

// TestCall tests Call
func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress    *types.Felt
		EntrypointSelector *types.Felt
		Calldata           []*types.Felt
		ExpectedResult     string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress:    types.StrToFelt("0xdeadbeef"),
				EntrypointSelector: types.StrToFelt("decimals"),
				Calldata:           []*types.Felt{types.StrToFelt("1234"), types.StrToFelt("5678")},
				ExpectedResult:     "0x12",
			},
		},
		"testnet": {
			{
				ContractAddress:    types.StrToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				EntrypointSelector: types.StrToFelt("decimals"),
				Calldata:           []*types.Felt{types.StrToFelt("78910"), types.StrToFelt("111213")},
				ExpectedResult:     "0x12",
			},
		},
		"mainnet": {
			{
				ContractAddress:    types.StrToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				EntrypointSelector: types.StrToFelt("decimals"),
				Calldata:           []*types.Felt{types.StrToFelt("141516"), types.StrToFelt("17181920")},
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
		call               types.FunctionInvoke
		BlockHashOrTag     string
		ExpectedOverallFee string
		ExpectedGasPrice   string
		ExpectedGasUsage   string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				call: types.FunctionInvoke{
					FunctionCall: types.FunctionCall{
						ContractAddress: &types.Felt{Int: big.NewInt(10000)},
						Calldata: []*types.Felt{
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000001"),
							types.StrToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							types.StrToFelt("0x0083afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000000"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000003"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000003"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000001"),
						},
						EntryPointSelector: types.StrToFelt("0x015d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad"),
					},
					Signature: []*types.Felt{
						types.StrToFelt("0x010e400d046147777c2ac5645024e1ee81c86d90b52d76ab8a8125e5f49612f9"),
						types.StrToFelt("0xadb92739205b4626fefb533b38d0071eb018e6ff096c98c17a6826b536817b"),
					},
					MaxFee:  types.StrToFelt("0x012c72866efa9b"),
					Version: 0,
				},
				BlockHashOrTag:     "0x0147c4b0f702079384e26d9d34a15e7758881e32b219fc68c076b09d0be13f8c",
				ExpectedOverallFee: "0x7134",
				ExpectedGasPrice:   "0x45",
				ExpectedGasUsage:   "0x1a4",
			},
		},
		"testnet": {},
		"mainnet": {
			{
				call: types.FunctionInvoke{
					FunctionCall: types.FunctionCall{
						ContractAddress: types.StrToFelt("0xdeadbeef"),
						Calldata: []*types.Felt{
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000001"),
							types.StrToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							types.StrToFelt("0x0083afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000000"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000003"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000003"),
							types.StrToFelt("0x04681402a7ab16c41f7e5d091f32fe9b78de096e0bd5962ce5bd7aaa4a441f64"),
							types.StrToFelt("0x000000000000000000000000000000000000000000000000001d41f6331e6800"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000000"),
							types.StrToFelt("0x0000000000000000000000000000000000000000000000000000000000000001"),
						},
						EntryPointSelector: types.StrToFelt("0x015d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad"),
					},
					Signature: []*types.Felt{
						types.StrToFelt("0x010e400d046147777c2ac5645024e1ee81c86d90b52d76ab8a8125e5f49612f9"),
						types.StrToFelt("0xadb92739205b4626fefb533b38d0071eb018e6ff096c98c17a6826b536817b"),
					},
					MaxFee:  types.StrToFelt("0x012c72866efa9b"),
					Version: 0,
				},
				BlockHashOrTag:     "0x0147c4b0f702079384e26d9d34a15e7758881e32b219fc68c076b09d0be13f8c",
				ExpectedOverallFee: "0xc84c599f51bd",
				ExpectedGasPrice:   "0x5df32828e",
				ExpectedGasUsage:   "0x221c",
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
		if fmt.Sprintf("0x%x", output.GasUsage) != test.ExpectedGasUsage {
			t.Fatalf("expected %s, got %s", test.ExpectedGasUsage, fmt.Sprintf("0x%x", output.GasUsage))
		}
		if fmt.Sprintf("0x%x", output.GasPrice) != test.ExpectedGasPrice {
			t.Fatalf("expected %s, got %s", test.ExpectedGasPrice, fmt.Sprintf("0x%x", output.GasPrice))
		}
	}
}
