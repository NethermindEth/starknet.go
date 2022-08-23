package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

// TestBlockHashAndNumber tests BlockHashAndNumber and check the returned value is strictly positive
func TestBlockHashAndNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}

	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {{}},
		"mainnet": {{}},
		"devnet":  {},
	}[testEnv]

	for range testSet {
		blockHashAndNumber, err := testConfig.client.BlockHashAndNumber(context.Background())

		if err != nil {
			t.Fatal(err)
		}
		if blockHashAndNumber.BlockNumber == 0 {
			t.Fatal("current block number should be higher or equal to 1")
		}
		if !strings.HasPrefix(blockHashAndNumber.BlockHash, "0x") {
			t.Fatal("current block hash should return a string starting with 0x")
		}
	}
}

// TestPendingBlockWithTxHashes tests TestPendingBlockWithTxHashes
func TestPendingBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{},
		},
		"mainnet": {{}},
	}[testEnv]

	for range testSet {
		_, err := testConfig.client.PendingBlockWithTxHashes(context.Background())
		if err == nil || !strings.Contains(err.Error(), "Pending data not supported in this configuration") {
			t.Fatal("PendingBlockWithTxHashes should not yet be supported")
		}
	}
}

// TestPendingBlockWithTxHashes tests TestPendingBlockWithTxHashes
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		RequestBlockHash         *BlockHash
		RequestBlockNumber       *BlockNumber
		RequestBlockTag          *string
		ExpectedError            error
		ExpectedFirstTransaction TxnHash
	}
	latestTag := "latest"
	errorTag := "error"
	testnetBlockHash := BlockHash("0x631127f10ab881f17c2cb1a3375e1c71352777b9ab0c1a2a7fe8fa9e201456e")
	testnetBlockNumber := BlockNumber(307417)
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				RequestBlockHash:         nil,
				RequestBlockNumber:       nil,
				RequestBlockTag:          &latestTag,
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash(""),
			},
			{
				RequestBlockHash:         nil,
				RequestBlockNumber:       nil,
				RequestBlockTag:          &errorTag,
				ExpectedError:            errBadRequest,
				ExpectedFirstTransaction: TxnHash(""),
			},
			{
				RequestBlockHash:         &testnetBlockHash,
				RequestBlockNumber:       nil,
				RequestBlockTag:          nil,
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash("0x32be2ddc447a19466760ef64a1c92e0683a7e1bcc68a677138020a65a81763d"),
			},
			{
				RequestBlockHash:         nil,
				RequestBlockNumber:       &testnetBlockNumber,
				RequestBlockTag:          nil,
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash("0x32be2ddc447a19466760ef64a1c92e0683a7e1bcc68a677138020a65a81763d"),
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		blockId := BlockID{
			BlockHash:   test.RequestBlockHash,
			BlockNumber: test.RequestBlockNumber,
			BlockTag:    test.RequestBlockTag,
		}
		block, err := testConfig.client.BlockWithTxHashes(context.Background(), blockId)
		if err != test.ExpectedError {
			t.Fatal("PendingBlockWithTxHashes match the expected error", err)
		}
		if test.ExpectedError != nil && block == nil {
			continue
		}
		if !strings.HasPrefix(string(block.BlockHash), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", block.BlockHash)
		}
		if block.Status == "" {
			t.Fatal("Status not be empty")
		}
		if len(block.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedFirstTransaction != "" && block.Transactions[0] != test.ExpectedFirstTransaction {
			t.Fatalf("the expected transaction 0 is %s, instead %s", test.ExpectedFirstTransaction, block.Transactions[0])
		}
	}
}

// TestDemonstrateMultipleEmbedding shows how you can guess what a type is and apply it
func TestDemonstrateMultipleEmbedding(t *testing.T) {
	type V1 struct {
		Label1 string
	}

	type V2 struct {
		Label2 string
	}

	type V interface{}

	type MyType struct {
		Data string
		Tx   []V
	}

	var my MyType

	jsonContent := `{"data": "data", "tx": [{"label2": "yes"}, {"label1": "no"}]}`

	err := json.Unmarshal([]byte(jsonContent), &my)
	if err != nil {
		t.Fatal("should succeed, instead", err)
	}
	for key, value := range my.Tx {
		switch local := value.(type) {
		case map[string]interface{}:
			if _, ok := local["label1"]; ok {
				labelOutput, err := json.Marshal(local)
				if err != nil {
					t.Fatal("label1Output should succeed, instead", err)
				}
				v := V1{}
				err = json.Unmarshal(labelOutput, &v)
				if err != nil {
					t.Fatal("V1 should succeed, instead", err)
				}
				my.Tx[key] = v
				continue
			}
			if _, ok := local["label2"]; ok {
				labelOutput, err := json.Marshal(local)
				if err != nil {
					t.Fatal("label1Output should succeed, instead", err)
				}
				v := V2{}
				err = json.Unmarshal(labelOutput, &v)
				if err != nil {
					t.Fatal("V1 should succeed, instead", err)
				}
				my.Tx[key] = v
				continue
			}
			fmt.Printf("we should not get here \n")
		default:
			fmt.Printf("%T", value)
		}
	}
	if _, ok := my.Tx[0].(V2); !ok {
		t.Fatalf("Tx[0] should be a V2, instead, %T", my.Tx[0])
	}
	if _, ok := my.Tx[1].(V1); !ok {
		t.Fatalf("Tx[0] should be a V1, instead, %T", my.Tx[1])
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
		t.Skipf("not implemented on %s", testEnv)
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
