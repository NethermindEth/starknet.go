package gateway

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func Test_Block_Devnet(t *testing.T) {
	gw := NewClient(WithChain("devnet"))

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

		block, _ := gw.Block(context.Background(), &BlockOptions{BlockNumber: test.BlockNumber})
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
		opts      *BlockOptions
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			BlockHash: "0x57f5102f7c61826926a4d76e544d2272cad091aa4e4b12e8e3e2120a220bd11",
			opts:      &BlockOptions{BlockNumber: 159179}}},

		"mainnet": {{
			BlockHash: "0x3bb30a6d1a3b6dcbc935b18c976126ab8d1fea60ef055be3c78530624824d50",
			opts:      &BlockOptions{BlockNumber: 5879},
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

func Test_BlockHashByID(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash string
		opts      *BlockOptions
	}

	testSet := map[string][]testSetType{
		"mainnet": {{
			BlockHash: "0x032f952924a746346868fecd72066df6092b416836c89ae00082b8f54c8e3331",
			opts:      &BlockOptions{BlockNumber: 6319},
		}},
		"testnet": {{
			BlockHash: "0x052af1130ada6c9d735e8cb4d513f00d2fc488dd27739550e384c712d73b8e06",
			opts:      &BlockOptions{BlockNumber: 380445},
		}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.BlockHashByID(context.Background(), test.opts.BlockNumber)

		if err != nil {
			t.Fatal(err)
		}
		if block != test.BlockHash {
			t.Fatalf("expecting %s, instead: %s", block, test.BlockHash)
		} else {
			fmt.Println(block)
		}
	}
}

func Test_BlockIDByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockNumber uint64
		opts        *BlockOptions
	}

	testSet := map[string][]testSetType{
		"mainnet": {{
			BlockNumber: 6319,
			opts:        &BlockOptions{BlockHash: "0x032f952924a746346868fecd72066df6092b416836c89ae00082b8f54c8e3331"},
		}},
		"testnet": {{
			BlockNumber: 380445,
			opts:        &BlockOptions{BlockHash: "0x052af1130ada6c9d735e8cb4d513f00d2fc488dd27739550e384c712d73b8e06"},
		}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.BlockIDByHash(context.Background(), test.opts.BlockHash)

		if err != nil {
			t.Fatal(err)
		}
		if block != test.BlockNumber {
			t.Fatalf("expecting %v, instead: %v", block, test.BlockNumber)
		} else {
			fmt.Println(block)
		}
	}
}
