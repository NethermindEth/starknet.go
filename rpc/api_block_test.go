package rpc

import (
	"context"
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
		if blockNumber == 0 {
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
		pending, err := testConfig.client.BlockWithTxHashes(context.Background(), WithBlockIDTag("pending"))
		if err == nil || !strings.Contains(err.Error(), "Pending data not supported in this configuration") {
			t.Fatal("PendingBlockWithTxHashes should not yet be supported")
		}
		if _, ok := pending.(PendingBlockWithTxHashes); !ok {
			t.Fatalf("expecting PendingBlockWithTxs, instead %T", pending)
		}
	}
}

// TestPendingBlockWithTxHashes tests TestPendingBlockWithTxHashes
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockIDOption            BlockIDOption
		ExpectedError            error
		ExpectedFirstTransaction TxnHash
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockIDOption:            WithBlockIDTag("latest"),
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash(""),
			},
			{
				BlockIDOption:            WithBlockIDTag("error"),
				ExpectedError:            errBadRequest,
				ExpectedFirstTransaction: TxnHash(""),
			},
			{
				BlockIDOption:            WithBlockIDHash(BlockHash("0x631127f10ab881f17c2cb1a3375e1c71352777b9ab0c1a2a7fe8fa9e201456e")),
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash("0x32be2ddc447a19466760ef64a1c92e0683a7e1bcc68a677138020a65a81763d"),
			},
			{
				BlockIDOption:            WithBlockIDNumber(BlockNumber(307417)),
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash("0x32be2ddc447a19466760ef64a1c92e0683a7e1bcc68a677138020a65a81763d"),
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		block := blockID{}
		_ = test.BlockIDOption(&block)
		blockWithTxHashesInterface, err := testConfig.client.BlockWithTxHashes(context.Background(), test.BlockIDOption)
		if err != test.ExpectedError {
			t.Fatal("PendingBlockWithTxHashes match the expected error", err)
		}
		if test.ExpectedError != nil && blockWithTxHashesInterface == nil {
			continue
		}
		blockWithTxHashes, ok := blockWithTxHashesInterface.(BlockWithTxHashes)
		if !ok {
			t.Fatalf("expecting BlockWithTxHashes, instead %T", blockWithTxHashesInterface)
		}
		if !strings.HasPrefix(string(blockWithTxHashes.BlockHash), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxHashes.BlockHash)
		}
		if blockWithTxHashes.Status == "" {
			t.Fatal("Status not be empty")
		}
		if len(blockWithTxHashes.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedFirstTransaction != "" && blockWithTxHashes.Transactions[0] != test.ExpectedFirstTransaction {
			t.Fatalf("the expected transaction 0 is %s, instead %s", test.ExpectedFirstTransaction, blockWithTxHashes.Transactions[0])
		}
	}
}

// TestBlockWithTxs tests TestPendingBlockWithTxHashes
func TestBlockWithTxs(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockIDOption            BlockIDOption
		ExpectedError            error
		ExpectedTxNumber         int
		ExpectedFirstTransaction TxnHash
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockIDOption:            WithBlockIDTag("latest"),
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash(""),
			},
			{
				BlockIDOption:            WithBlockIDTag("error"),
				ExpectedError:            errBadRequest,
				ExpectedFirstTransaction: TxnHash(""),
			},
			{
				BlockIDOption:            WithBlockIDHash(BlockHash("0x631127f10ab881f17c2cb1a3375e1c71352777b9ab0c1a2a7fe8fa9e201456e")),
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash("0x32be2ddc447a19466760ef64a1c92e0683a7e1bcc68a677138020a65a81763d"),
			},
			{
				BlockIDOption:            WithBlockIDNumber(BlockNumber(307417)),
				ExpectedError:            nil,
				ExpectedFirstTransaction: TxnHash("0x32be2ddc447a19466760ef64a1c92e0683a7e1bcc68a677138020a65a81763d"),
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		block := blockID{}
		_ = test.BlockIDOption(&block)
		blockWithTxsInterface, err := testConfig.client.BlockWithTxs(context.Background(), test.BlockIDOption)
		if err != test.ExpectedError {
			t.Fatal("PendingBlockWithTxs match the expected error", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(BlockWithTxs)
		if !ok {
			t.Fatalf("expecting BlockWithTxs, instead %T", blockWithTxsInterface)
		}
		if !strings.HasPrefix(string(blockWithTxs.BlockHash), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}
		if blockWithTxs.Status == "" {
			t.Fatal("Status not be empty")
		}
		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedFirstTransaction != "" && blockWithTxs.Transactions[0] != test.ExpectedFirstTransaction {
			t.Fatalf("the expected transaction 0 is %s, instead %s", test.ExpectedFirstTransaction, blockWithTxs.Transactions[0])
		}
	}
}

// TestStateUpdateByHash tests StateUpdateByHash
// TODO: this is not implemented yet with pathfinder as you can see from the
// [code](https://github.com/eqlabs/pathfinder/blob/927183552dad6dcdfebac16c8c1d2baf019127b1/crates/pathfinder/rpc_examples.sh#L37)
// check when it is and test when it is the case.
func TestStateUpdate(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockIDOption BlockIDOption
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockIDOption: WithBlockIDHash("0xdeadbeef"),
			},
		},
	}[testEnv]

	if len(testSet) == 0 {
		t.Skipf("not implemented on %s", testEnv)
	}
	for _, test := range testSet {
		output, err := testConfig.client.StateUpdate(context.Background(), test.BlockIDOption)
		if err != nil {
			t.Fatal(err)
		}
		blockID := &blockID{}
		test.BlockIDOption(blockID)
		if output.BlockHash != *blockID.BlockHash {
			t.Fatalf("expecting block %s, got %s", *blockID.BlockHash, output.BlockHash)
		}
	}
}
