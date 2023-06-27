package rpcv02

import (
	"context"
	"strings"
	"testing"

	"github.com/NethermindEth/caigo/types"
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
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockNumber, err := testConfig.provider.BlockNumber(context.Background())
		if err != nil {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if diff, err := spy.Compare(blockNumber, false); err != nil || diff != "FullMatch" {
			t.Fatal("expecting to match", err)
		}
		if blockNumber <= 3000 {
			t.Fatal("Block number should be > 3000, instead: ", blockNumber)
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
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockHashAndNumber, err := testConfig.provider.BlockHashAndNumber(context.Background())
		if err != nil {
			t.Fatal("BlockHashAndNumber match the expected error:", err)
		}
		if diff, err := spy.Compare(blockHashAndNumber, false); err != nil || diff != "FullMatch" {
			t.Fatal("expecting to match", err)
		}
		if blockHashAndNumber.BlockNumber < 3000 {
			t.Fatal("Block number should be > 3000, instead: ", blockHashAndNumber.BlockNumber)
		}
		if !strings.HasPrefix(blockHashAndNumber.BlockHash, "0x") {
			t.Fatal("current block hash should return a string starting with 0x")
		}
	}
}

// TestBlockWithTxHashes tests BlockWithTxHashes
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                   BlockID
		ExpectedError             error
		ExpectedBlockWithTxHashes Block
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:                   WithBlockHash(types.StrToFelt("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72")),
				ExpectedError:             nil,
				ExpectedBlockWithTxHashes: blockGoerli310370,
			},
			{
				BlockID:                   WithBlockNumber(310370),
				ExpectedError:             nil,
				ExpectedBlockWithTxHashes: blockGoerli310370,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		result, err := testConfig.provider.BlockWithTxHashes(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		block, ok := result.(*Block)
		if !ok {
			t.Fatalf("should return *Block, instead: %T\n", result)
		}
		if test.ExpectedError != nil {
			continue
		}
		if !strings.HasPrefix(block.BlockHash.String(), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", block.BlockHash)
		}

		if len(block.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedBlockWithTxHashes.BlockHash == types.StrToFelt("0x0") {
			continue
		}

		if !cmp.Equal(test.ExpectedBlockWithTxHashes, *block) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxHashes, block))
		}
	}
}

// TestBlockWithTxsAndInvokeTXNV0 tests block with Invoke TXN V0
func TestBlockWithTxsAndInvokeTXNV0(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		want                        Block
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:       WithBlockHash(types.StrToFelt("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72")),
				ExpectedError: nil,
				want:          fullBlockGoerli310370,
			},
			{
				BlockID:       WithBlockNumber(310370),
				ExpectedError: nil,
				want:          fullBlockGoerli310370,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		if !ok {
			t.Fatalf("expecting *rpv02.Block, instead %T", blockWithTxsInterface)
		}
		_, err = spy.Compare(blockWithTxs, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if !strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}

		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.want.BlockHash == types.StrToFelt("0x0") {
			continue
		}
		if !cmp.Equal(test.want.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.want.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]))
		}
	}
}

// TestBlockWithTxsAndDeployOrDeclare tests BlockWithTxs with Deploy or Declare TXN
func TestBlockWithTxsAndDeployOrDeclare(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		ExpectedBlockWithTxs        Block
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:                     WithBlockHash(types.StrToFelt("0x424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9")),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 14,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        fullBlockGoerli310843,
			},
			{
				BlockID:                     WithBlockNumber(310843),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 14,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        fullBlockGoerli310843,
			},
			{
				BlockID:                     WithBlockNumber(300114),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 3,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        fullBlockGoerli300114,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		if !ok {
			t.Fatalf("expecting *rpcv02.Block, instead %T", blockWithTxsInterface)
		}
		diff, err := spy.Compare(blockWithTxs, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(blockWithTxs, false)
		}
		if !strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}

		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedBlockWithTxs.BlockHash == types.StrToFelt("0x0") {
			continue
		}
		if !cmp.Equal(test.ExpectedBlockWithTxs.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxs.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]))
		}
	}
}

// TestBlockTransactionCount tests BlockTransactionCount
func TestBlockTransactionCount(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID       BlockID
		ExpectedCount uint64
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:       WithBlockNumber(300000),
				ExpectedCount: 10,
			},
		},
		"testnet": {
			{
				BlockID:       WithBlockNumber(300000),
				ExpectedCount: 23,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		count, err := testConfig.provider.BlockTransactionCount(context.Background(), test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(count, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(count, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if count != test.ExpectedCount {
			t.Fatalf("structure expecting %d, instead: %d", test.ExpectedCount, count)
		}
	}
}

func TestCaptureUnsupportedBlockTxn(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		StartBlock uint64
		EndBlock   uint64
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				StartBlock: 375919,
				EndBlock:   376000,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		for i := test.StartBlock; i < test.EndBlock; i++ {
			blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), WithBlockNumber(i))
			if err != nil {
				t.Fatal("BlockWithTxHashes match the expected error:", err)
			}
			blockWithTxs, ok := blockWithTxsInterface.(*Block)
			if !ok {
				t.Fatalf("expecting *rpcv02.Block, instead %T", blockWithTxsInterface)
			}
			for k, v := range blockWithTxs.Transactions {
				_, okv1 := v.(InvokeTxnV1)
				_, okl1 := v.(L1HandlerTxn)
				_, okdec := v.(DeclareTxn)
				_, okdep := v.(DeployTxn)
				_, okdepac := v.(DeployAccountTxn)
				if !okv1 && !okl1 && !okdec && !okdep && !okdepac {
					t.Fatalf("New Type Detected %T at Block(%d)/Txn(%d)", v, i, k)
				}
			}
		}
	}
}

// TODO: Find a block with such a Txn
// TestBlockWithTxsAndInvokeTXNV1 tests BlockWithTxs with Invoke V1
func TestBlockWithTxsAndInvokeTXNV1(t *testing.T) {
	_ = beforeEach(t)

	type testSetType struct {
		check bool
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				check: false,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		if test.check {
			t.Fatalf("error running test: %v", ErrNotImplemented)
		}
	}
}

// TestStateUpdate tests StateUpdateByHash
// TODO: this is not implemented yet with pathfinder as you can see from the
// [code](https://github.com/eqlabs/pathfinder/blob/927183552dad6dcdfebac16c8c1d2baf019127b1/crates/pathfinder/rpc_examples.sh#L37)
// check when it is and test when it is the case.
func TestStateUpdate(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                   BlockID
		ExpectedStateUpdateOutput StateUpdateOutput
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID: WithBlockNumber(300000),
				ExpectedStateUpdateOutput: StateUpdateOutput{
					BlockHash: types.StrToFelt("0x4f1cee281edb6cb31b9ba5a8530694b5527cf05c5ac6502decf3acb1d0cec4"),
					NewRoot:   "0x70677cda9269d47da3ff63bc87cf1c87d0ce167b05da295dc7fc68242b250b",
					OldRoot:   "0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af",
					StateDiff: StateDiff{
						StorageDiffs: []ContractStorageDiffItem{{
							Address: "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4",
							Key:     "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e",
							Value:   "0x630b4197",
						}},
					},
				},
			},
		},
		"testnet": {
			{
				BlockID: WithBlockNumber(300000),
				ExpectedStateUpdateOutput: StateUpdateOutput{
					BlockHash: types.StrToFelt("0x03b6d94b246815960f38b7dffc53cda192e7d1dcfff61e1bc042fb57e95f8349"),
					NewRoot:   "0x70677cda9269d47da3ff63bc87cf1c87d0ce167b05da295dc7fc68242b250b",
					OldRoot:   "0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af",
					StateDiff: StateDiff{
						StorageDiffs: []ContractStorageDiffItem{{
							Address: "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4",
							Key:     "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e",
							Value:   "0x630b4197",
						}},
					},
				},
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		stateUpdate, err := testConfig.provider.StateUpdate(context.Background(), test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		if stateUpdate.BlockHash.String() != test.ExpectedStateUpdateOutput.BlockHash.String() {
			t.Fatalf("structure expecting %s, instead: %s", test.ExpectedStateUpdateOutput.BlockHash.String(), stateUpdate.BlockHash.String())
		}
	}
}
