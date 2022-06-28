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
			{MissingMethod: "starknet_getNonce"},
			{MissingMethod: "starknet_protocolVersion"},
			{MissingMethod: "starknet_pendingTransactions"},
			{MissingMethod: "starknet_getStateUpdateByHash"},
		},
		"mock": {},
		"testnet": {
			{MissingMethod: "starknet_traceTransaction"},
			{MissingMethod: "starknet_traceBlockTransactions"},
			{MissingMethod: "starknet_getNonce"},
			{MissingMethod: "starknet_protocolVersion"},
			{MissingMethod: "starknet_pendingTransactions"},
			{MissingMethod: "starknet_getStateUpdateByHash"},
		},
	}[testEnv]

	for _, test := range testSet {
		var out string
		err := testConfig.client.do(context.Background(), test.MissingMethod, &out)

		if err == nil || err.Error() != "Method not found" {
			t.Fatalf("Method is now available, got %v\n", err)
		}
	}
}
