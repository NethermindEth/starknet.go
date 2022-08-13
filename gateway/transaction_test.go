package gateway

import (
	"testing"
	"context"
)

func TestGetTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionIndex int         
		BlockNumber      int         
		Transaction      Transaction 
		BlockHash        string      
		Status           string      
		opts    TransactionOptions
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"testnet": {{TransactionIndex: 3, 
			Status: "ACCEPTED_ON_L1", 
			BlockHash: "0x14a7a59e3e2d058d4c7c868e05907b2b49e324cc5b6af71182f008feb939e91",
			opts: TransactionOptions{TransactionHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5"}}},
	}[testEnv]

	for _, test := range testSet {
		tx, err := testConfig.client.Transaction(context.Background(), test.opts)

		if err != nil {
			t.Fatal(err)
		}
		if tx.TransactionIndex != test.TransactionIndex {
			t.Fatalf("expecting %d, instead: %d", test.TransactionIndex, tx.TransactionIndex)
		}
	}
}
