package rpc

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

// TestBlockWithTxHashes tests TestBlockWithTxHashes
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockIDOption             BlockIDOption
		ExpectedError             error
		ExpectedBlockWithTxHashes BlockWithTxHashes
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockIDOption: WithBlockIDTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockIDOption: WithBlockIDTag("error"),
				ExpectedError: errInvalidBlockTag,
			},
			{
				BlockIDOption:             WithBlockIDHash(BlockHash("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72")),
				ExpectedError:             nil,
				ExpectedBlockWithTxHashes: blockGoerli310370,
			},
			{
				BlockIDOption:             WithBlockIDNumber(BlockNumber(310370)),
				ExpectedError:             nil,
				ExpectedBlockWithTxHashes: blockGoerli310370,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		block := blockID{}
		_ = test.BlockIDOption(&block)
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		blockWithTxHashesInterface, err := testConfig.client.BlockWithTxHashes(context.Background(), test.BlockIDOption)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxHashesInterface == nil {
			continue
		}
		blockWithTxHashes, ok := blockWithTxHashesInterface.(*BlockWithTxHashes)
		if !ok {
			t.Fatalf("expecting BlockWithTxHashes, instead %T", blockWithTxHashesInterface)
		}
		if diff, err := spy.Compare(blockWithTxHashes); err != nil || diff != "FullMatch" {
			t.Fatal("expecting to match", err)
		}
		if !strings.HasPrefix(string(blockWithTxHashes.BlockHash), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxHashes.BlockHash)
		}

		if len(blockWithTxHashes.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedBlockWithTxHashes.BlockHash == "" {
			continue
		}
		if !cmp.Equal(test.ExpectedBlockWithTxHashes, *blockWithTxHashes) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxHashes, blockWithTxHashes))
		}
	}
}

// TestBlockWithTxs tests BlockWithTxs
func TestBlockWithTxs(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockIDOption        BlockIDOption
		ExpectedError        error
		ExpectedTxNumber     int
		ExpectedBlockWithTxs BlockWithTxs
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			// {
			// 	BlockIDOption: WithBlockIDTag("latest"),
			// 	ExpectedError: nil,
			// },
			// {
			// 	BlockIDOption: WithBlockIDTag("error"),
			// 	ExpectedError: errInvalidBlockTag,
			// },
			{
				BlockIDOption:        WithBlockIDHash(BlockHash("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72")),
				ExpectedError:        nil,
				ExpectedBlockWithTxs: fullBlockGoerli310370,
			},
			{
				BlockIDOption:        WithBlockIDNumber(BlockNumber(310370)),
				ExpectedError:        nil,
				ExpectedBlockWithTxs: fullBlockGoerli310370,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		block := blockID{}
		_ = test.BlockIDOption(&block)
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		blockWithTxsInterface, err := testConfig.client.BlockWithTxs(context.Background(), test.BlockIDOption)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*BlockWithTxs)
		if !ok {
			t.Fatalf("expecting BlockWithTxs, instead %T", blockWithTxsInterface)
		}
		diff, err := spy.Compare(blockWithTxs)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if !strings.HasPrefix(string(blockWithTxs.BlockHash), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}

		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedBlockWithTxs.BlockHash == "" {
			continue
		}
		if !cmp.Equal(test.ExpectedBlockWithTxs, *blockWithTxs) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxs, blockWithTxs))
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
