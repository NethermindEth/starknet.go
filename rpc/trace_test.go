package rpc

import (
	"context"
	"testing"
)

// TestNotImplemented checks the method not yet implemented
func TestNotImplemented(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		MissingMethod string
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mainnet": {
			{MissingMethod: "starknet_traceTransaction"},
			{MissingMethod: "starknet_traceBlockTransactions"},
		},
		"mock": {},
		"testnet": {
			{MissingMethod: "starknet_traceTransaction"},
			{MissingMethod: "starknet_traceBlockTransactions"},
		},
	}[testEnv]

	for _, test := range testSet {
		var out string
		err := do(context.Background(), testConfig.provider.c, test.MissingMethod, &out)

		if err == nil || err.Error() != "Method Not Found" {
			t.Fatalf("Method %s is now available, got %v\n", test.MissingMethod, err)
		}
	}
}
