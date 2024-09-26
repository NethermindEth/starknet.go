package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
//
// The function tests the MarshalJSON method of the BlockID struct by providing
// different scenarios and verifying the output against the expected values.
// The scenarios include testing the serialization of the "latest" and
// "pending" tags, testing an invalid tag, testing the serialization of a
// block number, and testing the serialization of a block hash.
// The function uses the testing.T parameter to report any errors that occur
// during the execution of the test cases.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockID_Marshal(t *testing.T) {
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
// The test checks the behavior of the BlockStatus function by iterating through a list of test cases.
//
// Parameters:
// - t: A testing.T object used for reporting test failures and logging.
// Returns:
//
//	none
func TestBlockStatus(t *testing.T) {
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

//go:embed tests/block/sepoliaBlockTxs65083.json
var rawBlock []byte

// TestBlock_Unmarshal tests the Unmarshal function of the Block type.
//
// This test case unmarshals raw block data into a Block instance and verifies
// that there are no errors during the process. If any error occurs, the test
// fails with a fatal error message.
//
// Parameters:
// - t: the testing object for running the test
// Returns:
//
//	none
func TestBlock_Unmarshal(t *testing.T) {
	b := Block{}
	if err := json.Unmarshal(rawBlock, &b); err != nil {
		t.Fatalf("Unmarshalling block: %v", err)
	}
}

func TestBlockWithReceipts(t *testing.T) {
	testConfig := beforeEach(t)
	require := require.New(t)

	type testSetType struct {
		BlockID                          BlockID
		ExpectedBlockWithReceipts        *BlockWithReceipts
		ExpectedPendingBlockWithReceipts *PendingBlockWithReceipts
	}

	var blockWithReceipt struct {
		Result BlockWithReceipts `json:"result"`
	}

	if testEnv == "testnet" {
		block, err := os.ReadFile("tests/blockWithReceipts/sepoliaBlockReceipts64159.json")
		require.NoError(err)
		require.NoError(json.Unmarshal(block, &blockWithReceipt))
	} else if testEnv == "mainnet" {
		block, err := os.ReadFile("tests/blockWithReceipts/mainnetBlockReceipts588763.json")
		require.NoError(err)
		require.NoError(json.Unmarshal(block, &blockWithReceipt))
	}

	deadBeef := utils.TestHexToFelt(t, "0xdeadbeef")
	var blockMock123 = BlockWithReceipts{
		BlockHeader{
			BlockHash: deadBeef,
		},
		"ACCEPTED_ON_L1",
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						BlockInvokeTxnV1{
							TransactionHash: deadBeef,
							InvokeTxnV1: InvokeTxnV1{
								Type:          "INVOKE",
								Version:       TransactionV1,
								SenderAddress: deadBeef,
							},
						},
					},
					Receipt: TransactionReceipt{
						TransactionHash: deadBeef,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
					},
				},
			},
		},
	}

	var pendingBlockMock123 = PendingBlockWithReceipts{
		PendingBlockHeader{
			ParentHash: deadBeef,
		},
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						BlockInvokeTxnV1{
							TransactionHash: deadBeef,
							InvokeTxnV1: InvokeTxnV1{
								Type:          "INVOKE",
								Version:       TransactionV1,
								SenderAddress: deadBeef,
							},
						},
					},
					Receipt: TransactionReceipt{
						TransactionHash: deadBeef,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
					},
				},
			},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:                          WithBlockTag("latest"),
				ExpectedBlockWithReceipts:        &blockMock123,
				ExpectedPendingBlockWithReceipts: nil,
			},
			{
				BlockID:                          WithBlockTag("latest"),
				ExpectedBlockWithReceipts:        nil,
				ExpectedPendingBlockWithReceipts: &pendingBlockMock123,
			},
		},
		"testnet": {
			{
				BlockID: WithBlockTag("pending"),
			},
			{
				BlockID:                   WithBlockNumber(64159),
				ExpectedBlockWithReceipts: &blockWithReceipt.Result,
			},
		},
		"mainnet": {
			{
				BlockID: WithBlockTag("pending"),
			},
			{
				BlockID:                   WithBlockNumber(588763),
				ExpectedBlockWithReceipts: &blockWithReceipt.Result,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		result, err := testConfig.provider.BlockWithReceipts(context.Background(), test.BlockID)
		require.NoError(err, "Error in BlockWithReceipts")

		switch resultType := result.(type) {
		case *BlockWithReceipts:
			block, ok := result.(*BlockWithReceipts)
			require.True(ok, fmt.Sprintf("should return *BlockWithReceipts, instead: %T\n", result))
			require.True(strings.HasPrefix(block.BlockHash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.BlockHash)
			require.NotEmpty(block.Transactions, "the number of transactions should not be 0")

			if test.ExpectedBlockWithReceipts != nil {
				require.Exactly(block, test.ExpectedBlockWithReceipts)
			}
		case *PendingBlockWithReceipts:
			pBlock, ok := result.(*PendingBlockWithReceipts)
			require.True(ok, fmt.Sprintf("should return *PendingBlockWithReceipts, instead: %T\n", result))

			if testEnv == "mock" {
				require.Exactly(pBlock, test.ExpectedPendingBlockWithReceipts)
			} else {
				require.NotEmpty(pBlock.ParentHash, "Error in PendingBlockWithReceipts ParentHash")
				require.NotEmpty(pBlock.SequencerAddress, "Error in PendingBlockWithReceipts SequencerAddress")
				require.NotEmpty(pBlock.Timestamp, "Error in PendingBlockWithReceipts Timestamp")
			}

		default:
			t.Fatalf("unexpected block type, found: %T\n", resultType)
		}
	}
}
