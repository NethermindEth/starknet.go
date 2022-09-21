package rpc

import (
	"context"
	"os"
	"testing"
)

// TestCall tests Call
func TestAccountNonce(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
			Provider *Client
			Address  string
			PrivateKeyEnvVar  string
	}

	testSet := map[string][]testSetType{
		"mock": {
			
		},
		"testnet": {
			{
				Address: "0x19e63006d7df131737f5222283da28de2d9e2f0ee92fdc4c4c712d1659826b0",
				PrivateKeyEnvVar: "TESTNET_ACCOUNT_PRIVATE_KEY",
			},
			
		},
		"mainnet": {
			
		},
	}[testEnv]

	for _, test := range testSet {

		account, err := testConfig.client.NewAccount(os.Getenv(test.PrivateKeyEnvVar),test.Address)
		 
		if err != nil {
			t.Fatal(err)
		}

		nonce, err := account.Nonce(context.Background())

		if err != nil {
			t.Fatal(err)
		}
		if nonce.Uint64() <= 1 {
			t.Fatal("nonce should be > 1", nonce.Uint64())
		}
	}
}