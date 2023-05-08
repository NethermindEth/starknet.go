package gateway_test

import (
	"context"
	"log"
	"testing"

	"github.com/smartcontractkit/caigo/gateway"
)

func Test_Block_Devnet(t *testing.T) {
	gw := gateway.NewClient(gateway.WithChain("devnet"))

	type testSetType struct {
		BlockNumber uint64
	}

	testSet := map[string][]testSetType{
		"devnet": {{BlockNumber: 1}},
	}[testEnv]

	for _, test := range testSet {

		err := setupDevnet(context.Background())

		if err != nil {
			log.Fatal("error starting test", err)
		}
		blockNumber := test.BlockNumber
		block, _ := gw.Block(context.Background(), &gateway.BlockOptions{BlockNumber: &blockNumber})
		if block.BlockNumber != int(test.BlockNumber) {
			t.Fatal("block number should be 1")
		}
		if len(block.Transactions) == 0 {
			log.Fatal("should have atleast 1 tx")
		}
	}
}

func Test_Block(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash string
		opts      *gateway.BlockOptions
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			// TODO/3: instead of testing just the blockHash we should
			// (1) Marshal the block back into a []byte and (2) compare it with the original message
			// check "github.com/nsf/jsondiff" for ideas how to do that.
			BlockHash: "0x57f5102f7c61826926a4d76e544d2272cad091aa4e4b12e8e3e2120a220bd11",
			opts:      &gateway.BlockOptions{BlockNumber: func() *uint64 { var v uint64 = 159179; return &v }()}}},

		"mainnet": {{
			BlockHash: "0x3bb30a6d1a3b6dcbc935b18c976126ab8d1fea60ef055be3c78530624824d50",
			opts:      &gateway.BlockOptions{BlockNumber: func() *uint64 { var v uint64 = 5879; return &v }()},
		}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.Block(context.Background(), test.opts)

		if err != nil {
			t.Fatal(err)
		}
		if block.BlockHash != test.BlockHash {
			t.Fatalf("expecting %s, instead: %s", block.BlockHash, test.BlockHash)
		}
	}
}
