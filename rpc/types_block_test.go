package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/require"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
//
// The function tests the MarshalJSON method of the BlockID struct by providing
// different scenarios and verifying the output against the expected values.
// The scenarios include testing the serialisation of the "latest" and
// "pending" tags, testing an invalid tag, testing the serialisation of a
// block number, and testing the serialisation of a block hash.
// The function uses the testing.T parameter to report any errors that occur
// during the execution of the test cases.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestBlockID_Marshal(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	blockNumber := uint64(420)
	for _, tc := range []struct {
		id      BlockID
		want    string
		wantErr error
	}{{
		id: BlockID{
			Tag: "latest",
		},
		want: `"latest"`,
	}, {
		id: BlockID{
			Tag: "pending",
		},
		want: `"pending"`,
	}, {
		id: BlockID{
			Tag: "bad tag",
		},
		wantErr: ErrInvalidBlockID,
	}, {
		id: BlockID{
			Number: &blockNumber,
		},
		want: `{"block_number":420}`,
	}, {
		id: func() BlockID {
			h, _ := new(felt.Felt).SetString("0xdead")

			return BlockID{
				Hash: h,
			}
		}(),
		want: `{"block_hash":"0xdead"}`,
	}} {
		b, err := tc.id.MarshalJSON()
		if err != nil && tc.wantErr == nil {
			t.Errorf("marshalling block id: %v", err)
		} else if err != nil && !errors.Is(err, tc.wantErr) {
			t.Errorf("block error mismatch, want: %v, got: %v", tc.wantErr, err)
		}

		if string(b) != tc.want {
			t.Errorf("block id mismatch, want: %s, got: %s", tc.want, b)
		}
	}
}

// TestBlockStatus is a unit test for the BlockStatus function.
//
// The test checks the behaviour of the BlockStatus function by iterating through a list of test cases.
//
// Parameters:
//   - t: A testing.T object used for reporting test failures and logging.
//
// Returns:
//
//	none
func TestBlockStatus(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	for _, tc := range []struct {
		status string
		want   BlockStatus
	}{{
		status: `"PENDING"`,
		want:   BlockStatus_Pending,
	}, {
		status: `"ACCEPTED_ON_L2"`,
		want:   BlockStatus_AcceptedOnL2,
	}, {
		status: `"ACCEPTED_ON_L1"`,
		want:   BlockStatus_AcceptedOnL1,
	}, {
		status: `"REJECTED"`,
		want:   BlockStatus_Rejected,
	}} {
		tx := new(BlockStatus)
		if err := json.Unmarshal([]byte(tc.status), tx); err != nil {
			t.Errorf("unmarshalling status want: %s", err)
		}
	}
}

//go:embed testData/block/sepoliaBlockTxs65083.json
var rawBlock []byte

// TestBlock_Unmarshal tests the Unmarshal function of the Block type.
//
// This test case unmarshals raw block data into a Block instance and verifies
// that there are no errors during the process. If any error occurs, the test
// fails with a fatal error message.
//
// Parameters:
//   - t: the testing object for running the test
//
// Returns:
//
//	none
func TestBlock_Unmarshal(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	b := Block{}
	if err := json.Unmarshal(rawBlock, &b); err != nil {
		t.Fatalf("Unmarshalling block: %v", err)
	}
}

func TestBlockWithReceipts(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv)

	testConfig := beforeEach(t, false)

	type testSetType struct {
		BlockID                          BlockID
		ExpectedBlockWithReceipts        *BlockWithReceipts
		ExpectedPendingBlockWithReceipts *PendingBlockWithReceipts
	}

	var blockWithReceipt BlockWithReceipts

	switch tests.TEST_ENV {
	case tests.TestnetEnv:
		blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/sepoliaBlockReceipts64159.json", "result")
	case tests.MainnetEnv:
		blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/mainnetBlockReceipts588763.json", "result")
	}

	deadBeef := internalUtils.TestHexToFelt(t, "0xdeadbeef")
	blockMock123 := BlockWithReceipts{
		BlockHeader{
			Hash: deadBeef,
		},
		"ACCEPTED_ON_L1",
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						Hash: deadBeef,
						Transaction: InvokeTxnV1{
							Type:          "INVOKE",
							Version:       TransactionV1,
							SenderAddress: deadBeef,
						},
					},
					Receipt: TransactionReceipt{
						Type:            "INVOKE",
						Hash:            deadBeef,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
					},
				},
			},
		},
	}

	pendingBlockMock123 := PendingBlockWithReceipts{
		PendingBlockHeader{
			ParentHash: deadBeef,
		},
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						Hash: deadBeef,
						Transaction: InvokeTxnV1{
							Type:          "INVOKE",
							Version:       TransactionV1,
							SenderAddress: deadBeef,
						},
					},
					Receipt: TransactionReceipt{
						Type:            "INVOKE",
						Hash:            deadBeef,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
					},
				},
			},
		},
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:                          WithBlockTag("latest"),
				ExpectedBlockWithReceipts:        &blockMock123,
				ExpectedPendingBlockWithReceipts: nil,
			},
			{
				BlockID:                          WithBlockTag("pending"),
				ExpectedBlockWithReceipts:        nil,
				ExpectedPendingBlockWithReceipts: &pendingBlockMock123,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockTag("pending"),
			},
			{
				BlockID:                   WithBlockNumber(64159),
				ExpectedBlockWithReceipts: &blockWithReceipt,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID: WithBlockTag("pending"),
			},
			{
				BlockID:                   WithBlockNumber(588763),
				ExpectedBlockWithReceipts: &blockWithReceipt,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		result, err := testConfig.provider.BlockWithReceipts(context.Background(), test.BlockID)
		require.NoError(t, err, "Error in BlockWithReceipts")

		switch resultType := result.(type) {
		case *BlockWithReceipts:
			block, ok := result.(*BlockWithReceipts)
			require.True(t, ok, fmt.Sprintf("should return *BlockWithReceipts, instead: %T\n", result))
			require.True(t, strings.HasPrefix(block.Hash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.Hash)
			require.NotEmpty(t, block.Transactions, "the number of transactions should not be 0")

			require.Exactly(t, block, test.ExpectedBlockWithReceipts)
		case *PendingBlockWithReceipts:
			pBlock, ok := result.(*PendingBlockWithReceipts)
			require.True(t, ok, fmt.Sprintf("should return *PendingBlockWithReceipts, instead: %T\n", result))

			if tests.TEST_ENV == tests.MockEnv {
				require.Exactly(t, pBlock, test.ExpectedPendingBlockWithReceipts)
			} else {
				require.NotEmpty(t, pBlock.ParentHash, "Error in PendingBlockWithReceipts ParentHash")
				require.NotEmpty(t, pBlock.SequencerAddress, "Error in PendingBlockWithReceipts SequencerAddress")
				require.NotEmpty(t, pBlock.Timestamp, "Error in PendingBlockWithReceipts Timestamp")
			}

		default:
			t.Fatalf("unexpected block type, found: %T\n", resultType)
		}
	}
}
