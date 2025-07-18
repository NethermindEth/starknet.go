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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
//
// The function tests the MarshalJSON method of the BlockID struct by providing
// different scenarios and verifying the output against the expected values.
// The scenarios include testing the serialisation of the "latest" and
// "pre_confirmed" tags, testing an invalid tag, testing the serialisation of a
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
			Tag: "pre_confirmed",
		},
		want: `"pre_confirmed"`,
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
	}{
		{
			status: `"PRE_CONFIRMED"`,
			want:   BlockStatus_Pre_confirmed,
		},
		{
			status: `"ACCEPTED_ON_L2"`,
			want:   BlockStatus_AcceptedOnL2,
		},
		{
			status: `"ACCEPTED_ON_L1"`,
			want:   BlockStatus_AcceptedOnL1,
		},
	} {
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
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		BlockID                                BlockID
		ExpectedBlockWithReceipts              *BlockWithReceipts
		ExpectedPre_confirmedBlockWithReceipts *Pre_confirmedBlockWithReceipts
	}

	var blockWithReceipt BlockWithReceipts

	switch tests.TEST_ENV {
	case tests.TestnetEnv:
		blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/sepoliaBlockReceipts64159.json", "result")
	case tests.MainnetEnv:
		blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/mainnetBlockReceipts588763.json", "result")
	case tests.IntegrationEnv:
		blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/integration1300000.json", "result")
	}

	blockMock123 := BlockWithReceipts{
		BlockHeader{
			Hash: internalUtils.RANDOM_FELT,
		},
		"ACCEPTED_ON_L1",
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						Hash: internalUtils.RANDOM_FELT,
						Transaction: InvokeTxnV1{
							Type:          "INVOKE",
							Version:       TransactionV1,
							SenderAddress: internalUtils.RANDOM_FELT,
						},
					},
					Receipt: TransactionReceipt{
						Type:            "INVOKE",
						Hash:            internalUtils.RANDOM_FELT,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
						ActualFee: FeePayment{
							Amount: internalUtils.RANDOM_FELT,
							Unit:   UnitFri,
						},
					},
				},
			},
		},
	}

	pre_confirmedBlockMock123 := Pre_confirmedBlockWithReceipts{
		Pre_confirmedBlockHeader{
			Number: 1234,
		},
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						Hash: internalUtils.RANDOM_FELT,
						Transaction: InvokeTxnV1{
							Type:          "INVOKE",
							Version:       TransactionV1,
							SenderAddress: internalUtils.RANDOM_FELT,
						},
					},
					Receipt: TransactionReceipt{
						Type:            "INVOKE",
						Hash:            internalUtils.RANDOM_FELT,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
						ActualFee: FeePayment{
							Amount: internalUtils.RANDOM_FELT,
							Unit:   UnitFri,
						},
					},
				},
			},
		},
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:                                WithBlockTag("latest"),
				ExpectedBlockWithReceipts:              &blockMock123,
				ExpectedPre_confirmedBlockWithReceipts: nil,
			},
			{
				BlockID:                                WithBlockTag("pre_confirmed"),
				ExpectedBlockWithReceipts:              nil,
				ExpectedPre_confirmedBlockWithReceipts: &pre_confirmedBlockMock123,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockTag("pre_confirmed"),
			},
			{
				BlockID:                   WithBlockNumber(64159),
				ExpectedBlockWithReceipts: &blockWithReceipt,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID: WithBlockTag("pre_confirmed"),
			},
			{
				BlockID:                   WithBlockNumber(1_300_000),
				ExpectedBlockWithReceipts: &blockWithReceipt,
			},
		},
		tests.MainnetEnv: {
			{
				BlockID: WithBlockTag("pre_confirmed"),
			},
			{
				BlockID:                   WithBlockNumber(588763),
				ExpectedBlockWithReceipts: &blockWithReceipt,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(string(blockID), func(t *testing.T) {
			result, err := testConfig.Provider.BlockWithReceipts(context.Background(), test.BlockID)
			require.NoError(t, err, "Error in BlockWithReceipts")

			switch resultType := result.(type) {
			case *BlockWithReceipts:
				block, ok := result.(*BlockWithReceipts)
				require.True(t, ok, fmt.Sprintf("should return *BlockWithReceipts, instead: %T\n", result))
				assert.True(t, strings.HasPrefix(block.Hash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.Hash)
				assert.NotEmpty(t, block.Transactions, "the number of transactions should not be 0")

				assert.Exactly(t, block, test.ExpectedBlockWithReceipts)
			case *Pre_confirmedBlockWithReceipts:
				pBlock, ok := result.(*Pre_confirmedBlockWithReceipts)
				require.True(t, ok, fmt.Sprintf("should return *Pre_confirmedBlockWithReceipts, instead: %T\n", result))

				if tests.TEST_ENV == tests.MockEnv {
					assert.Exactly(t, pBlock, test.ExpectedPre_confirmedBlockWithReceipts)
				} else {
					assert.NotEmpty(t, pBlock.Number, "Error in Pre_confirmedBlockWithReceipts ParentHash")
					assert.NotEmpty(t, pBlock.SequencerAddress, "Error in Pre_confirmedBlockWithReceipts SequencerAddress")
					assert.NotEmpty(t, pBlock.Timestamp, "Error in Pre_confirmedBlockWithReceipts Timestamp")
				}

			default:
				t.Fatalf("unexpected block type, found: %T\n", resultType)
			}
		})
	}
}
