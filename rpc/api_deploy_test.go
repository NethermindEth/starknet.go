package rpc

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

// TestAddDeployTransaction tests AddDeployTransaction
func TestAddDeployTransaction(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

	type testSetType struct {
	}
	testSet := map[string][]testSetType{
		"mock": {
			{},
		},
		"testnet": {
			{},
		},
		"mainnet": {
			{},
		},
	}[testEnv]

	for range testSet {
		output, err := testConfig.client.AddDeployTransaction(context.Background(), "0x0", []string{}, types.ContractClass{})
		if err != nil {
			t.Fatal(err)
		}
		if output.TransactionHash == "" {
			t.Fatal("should return a transaction")
		}
		if output.ContractAddress == "" {
			t.Fatal("should return a contract address")
		}
	}
}
