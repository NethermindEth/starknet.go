package rpc

import (
	"context"
	"fmt"
	"math/big"
	"testing"
)

// TestBlockNumber tests BlockNumber and check the returned value is strictly positive
func TestBlockNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}

	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {{}},
		"mainnet": {{}},
		"devnet":  {},
	}[testEnv]

	for range testSet {
		blockNumber, err := testConfig.client.BlockNumber(context.Background())

		if err != nil {
			t.Fatal(err)
		}
		if blockNumber == nil || blockNumber.Int64() <= 0 {
			t.Fatal("current block number should be higher or equal to 1")
		}
	}
}

// TestBlockByNumber tests BlockByNumber
func TestBlockByNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockNumber       *big.Int
		BlockScope        string
		ExpectedBlockHash string
		ExpectedStatus    string
		ExpectedTx0Hash   string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockNumber:       big.NewInt(1000),
				BlockScope:        "FULL_TXN_AND_RECEIPTS",
				ExpectedBlockHash: "0xdeadbeef",
				ExpectedStatus:    "ACCEPTED_ON_L1",
				ExpectedTx0Hash:   "0xdeadbeef",
			},
			{
				BlockNumber:       big.NewInt(1000),
				BlockScope:        "FULL_TXNS",
				ExpectedBlockHash: "0xdeadbeef",
				ExpectedStatus:    "",
				ExpectedTx0Hash:   "0xdeadbeef",
			},
		},
		"testnet": {
			{
				BlockNumber:       big.NewInt(242060),
				BlockScope:        "FULL_TXN_AND_RECEIPTS",
				ExpectedBlockHash: "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044",
				ExpectedStatus:    "ACCEPTED_ON_L1",
				ExpectedTx0Hash:   "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0",
			},
			{
				BlockNumber:       big.NewInt(242060),
				BlockScope:        "FULL_TXNS",
				ExpectedBlockHash: "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044",
				ExpectedStatus:    "",
				ExpectedTx0Hash:   "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0",
			},
		},
		"mainnet": {{
			BlockNumber:       big.NewInt(1500),
			BlockScope:        "FULL_TXN_AND_RECEIPTS",
			ExpectedBlockHash: "0x6f8e6413281c43bfcb9f96e315a08c57c619c9da4b10e2cb7d33369f3fb75a0",
			ExpectedStatus:    "ACCEPTED_ON_L1",
			ExpectedTx0Hash:   "0x5f904b9185d4ed442846ac7e26bc4c60249a2a7f0bb85376c0bc7459665bae6",
		}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.BlockByNumber(context.Background(), test.BlockNumber, test.BlockScope)
		if err != nil {
			t.Fatal(err)
		}
		if block.BlockHash != test.ExpectedBlockHash {
			t.Fatalf("blockhash mismatch, expect %s, got %s :",
				test.ExpectedBlockHash,
				block.BlockHash)
		}
		if block.Transactions[0].TransactionHash != test.ExpectedTx0Hash {
			t.Fatalf("tx[0] mismatch, expect %s, got %s :",
				test.ExpectedTx0Hash,
				block.Transactions[0].TransactionHash)
		}
		if block.Transactions[0].TransactionReceipt.Status != test.ExpectedStatus {
			t.Fatalf("tx receipt mismatch, expect %s, got %s :",
				test.ExpectedStatus,
				block.Transactions[0].TransactionReceipt.Status)
		}
	}
}

// TestBlockByHash tests BlockByHash
func TestBlockByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash           string
		BlockScope          string
		ExpectedBlockNumber int
		ExpectedTx0Hash     string
		ExpectedStatus      string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockHash:           "0xdeadbeef",
				BlockScope:          "FULL_TXN_AND_RECEIPTS",
				ExpectedBlockNumber: 1000,
				ExpectedTx0Hash:     "0xdeadbeef",
				ExpectedStatus:      "ACCEPTED_ON_L1",
			},
			{
				BlockHash:           "0xdeadbeef",
				BlockScope:          "FULL_TXNS",
				ExpectedBlockNumber: 1000,
				ExpectedTx0Hash:     "0xdeadbeef",
				ExpectedStatus:      "",
			},
		},
		"testnet": {
			{
				BlockHash:           "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044",
				BlockScope:          "FULL_TXN_AND_RECEIPTS",
				ExpectedBlockNumber: 242060,
				ExpectedTx0Hash:     "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0",
				ExpectedStatus:      "ACCEPTED_ON_L1",
			},
			{
				BlockHash:           "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044",
				BlockScope:          "FULL_TXNS",
				ExpectedBlockNumber: 242060,
				ExpectedTx0Hash:     "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0",
				ExpectedStatus:      "",
			},
		},
		"mainnet": {{
			BlockHash:           "0x6f8e6413281c43bfcb9f96e315a08c57c619c9da4b10e2cb7d33369f3fb75a0",
			BlockScope:          "FULL_TXN_AND_RECEIPTS",
			ExpectedBlockNumber: 1500,
			ExpectedTx0Hash:     "0x5f904b9185d4ed442846ac7e26bc4c60249a2a7f0bb85376c0bc7459665bae6",
			ExpectedStatus:      "ACCEPTED_ON_L1",
		}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.BlockByHash(context.Background(), test.BlockHash, test.BlockScope)
		if err != nil {
			t.Fatal(err)
		}
		if block.BlockNumber != test.ExpectedBlockNumber {
			t.Fatalf("blockNumber mismatch, expect %d, got %d :",
				test.ExpectedBlockNumber,
				block.BlockNumber)
		}
		if block.Transactions[0].TransactionHash != test.ExpectedTx0Hash {
			t.Fatalf("tx[0] mismatch, expect %s, got %s :",
				test.ExpectedTx0Hash,
				block.Transactions[0].TransactionHash)
		}
		if block.Transactions[0].TransactionReceipt.Status != test.ExpectedStatus {
			t.Fatalf("tx receipt mismatch, expect %s, got %s :",
				test.ExpectedStatus,
				block.Transactions[0].TransactionReceipt.Status)
		}
	}
}

// TestStateUpdateByHash tests StateUpdateByHash
// TODO: this is not implemented yet with pathfinder as you can see from the
// [code](https://github.com/eqlabs/pathfinder/blob/927183552dad6dcdfebac16c8c1d2baf019127b1/crates/pathfinder/rpc_examples.sh#L37)
// check when it is and test when it is the case.
func TestStateUpdateByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHashOrTag string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockHashOrTag: "0xdeadbeef",
			},
		},
	}[testEnv]

	if len(testSet) == 0 {
		t.Skip(fmt.Sprintf("not implemented on %s", testEnv))
	}
	for _, test := range testSet {
		output, err := testConfig.client.StateUpdateByHash(context.Background(), test.BlockHashOrTag)
		if err != nil {
			t.Fatal(err)
		}
		if output.BlockHash != test.BlockHashOrTag {
			t.Fatalf("expecting block %s, got %s", test.BlockHashOrTag, output.BlockHash)
		}
	}
}

// TestBlockTransactionCountByHash tests BlockTransactionCountByHash
func TestBlockTransactionCountByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash       string
		ExpectedTxCount int
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockHash:       "0xdeadbeef",
				ExpectedTxCount: 7,
			},
		},
		"testnet": {
			{
				BlockHash:       "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044",
				ExpectedTxCount: 31,
			},
		},
		"mainnet": {{
			BlockHash:       "0x6f8e6413281c43bfcb9f96e315a08c57c619c9da4b10e2cb7d33369f3fb75a0",
			ExpectedTxCount: 31,
		}},
	}[testEnv]

	for _, test := range testSet {
		txCount, err := testConfig.client.BlockTransactionCountByHash(context.Background(), test.BlockHash)
		if err != nil {
			t.Fatal(err)
		}
		if txCount != test.ExpectedTxCount {
			t.Fatalf("txCount mismatch, expect %d, got %d :",
				test.ExpectedTxCount,
				txCount)
		}
	}
}

// TestBlockTransactionCountByHash tests BlockTransactionCountByHash
func TestBlockTransactionCountByNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockNumberOrTag interface{}
		ExpectedTxCount  int
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockNumberOrTag: 666,
				ExpectedTxCount:  7,
			},
		},
		"testnet": {
			{
				BlockNumberOrTag: 242060,
				ExpectedTxCount:  31,
			},
		},
		"mainnet": {{
			BlockNumberOrTag: 1500,
			ExpectedTxCount:  31,
		}},
	}[testEnv]

	for _, test := range testSet {
		txCount, err := testConfig.client.BlockTransactionCountByNumber(context.Background(), test.BlockNumberOrTag)
		if err != nil {
			t.Fatal(err)
		}
		if txCount != test.ExpectedTxCount {
			t.Fatalf("txCount mismatch, expect %d, got %d :",
				test.ExpectedTxCount,
				txCount)
		}
	}
}
