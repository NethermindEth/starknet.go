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
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
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

//go:embed tests/block/block.json
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

	type testSetType struct {
		BlockID                          BlockID
		ExpectedBlockWithReceipts        *BlockWithReceipts
		ExpectedPendingBlockWithReceipts *PendingBlockWithReceipts
	}

	var blockSepolia64159 = BlockWithReceipts{
		BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6df565874b2ea6a02d346a23f9efb0b26abbf5708b51bb12587f88a49052964"),
			ParentHash:       utils.TestHexToFelt(t, "0x1406ec9385293905d6c20e9c5aa0bbf9f63f87d39cf12fcdfef3ed0d056c0f5"),
			BlockNumber:      64159,
			NewRoot:          utils.TestHexToFelt(t, "0x310be818a18de0d6f6c1391f467d0dbd1a2753e6dde876449448465f8e617f0"),
			Timestamp:        1714901729,
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			L1GasPrice: ResourcePrice{
				PriceInFRI: utils.TestHexToFelt(t, "0xdf0413d3c777"),
				PriceInWei: utils.TestHexToFelt(t, "0x185f2d3eb5"),
			},
			L1DataGasPrice: ResourcePrice{
				PriceInFRI: utils.TestHexToFelt(t, "0xa41c1219f8849"),
				PriceInWei: utils.TestHexToFelt(t, "0x11ef315a9ab"),
			},
			L1DAMode:        L1DAModeBlob,
			StarknetVersion: "0.13.1.1",
		},
		"ACCEPTED_ON_L1",
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				// a lot of previous transactions ...
				{ // this is the last transaction of this block
					Transaction: BlockTransaction{
						BlockInvokeTxnV1{
							TransactionHash: utils.TestHexToFelt(t, "0x5d41f4dec3678156d3888d6b890648c3baa02d866820689d5f8b3e20735521b"),
							InvokeTxnV1: InvokeTxnV1{
								Type:          "INVOKE",
								Version:       TransactionV1,
								Nonce:         utils.TestHexToFelt(t, "0x3a"),
								MaxFee:        utils.TestHexToFelt(t, "0x1bad55a98e1c1"),
								SenderAddress: utils.TestHexToFelt(t, "0x3543d2f0290e39a08cfdf2245f14aec7dca60672b7c7458375f3cb3834e1067"),
								Signature: []*felt.Felt{
									utils.TestHexToFelt(t, "0x1"),
									utils.TestHexToFelt(t, "0x7f1d7eb58b1467c75dad461953ba6d931c7b46b51ee2055996ebaa595583985"),
									utils.TestHexToFelt(t, "0x34fb6291f60e2a7ec088a92075bb977d89eb779e8f51fa638799d27a999d9fe"),
								},
								Calldata: []*felt.Felt{
									utils.TestHexToFelt(t, "0x1"),
									utils.TestHexToFelt(t, "0x517567ac7026ce129c950e6e113e437aa3c83716cd61481c6bb8c5057e6923e"),
									utils.TestHexToFelt(t, "0xcaffbd1bd76bd7f24a3fa1d69d1b2588a86d1f9d2359b13f6a84b7e1cbd126"),
									utils.TestHexToFelt(t, "0x6"),
									utils.TestHexToFelt(t, "0x5265706f73736573734275696c64696e67"),
									utils.TestHexToFelt(t, "0x4"),
									utils.TestHexToFelt(t, "0x5"),
									utils.TestHexToFelt(t, "0xa25"),
									utils.TestHexToFelt(t, "0x1"),
									utils.TestHexToFelt(t, "0xe52"),
								},
							},
						},
					},
					Receipt: UnknownTransactionReceipt{
						TransactionReceipt: InvokeTransactionReceipt{
							Type:            "INVOKE",
							TransactionHash: utils.TestHexToFelt(t, "0x5d41f4dec3678156d3888d6b890648c3baa02d866820689d5f8b3e20735521b"),
							ActualFee: FeePayment{
								Amount: utils.TestHexToFelt(t, "0x1276c2c6477ed"),
								Unit:   UnitWei,
							},
							ExecutionStatus: TxnExecutionStatusSUCCEEDED,
							FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
							MessagesSent:    []MsgToL1{},
							Events: []Event{
								{
									FromAddress: utils.TestHexToFelt(t, "0x517567ac7026ce129c950e6e113e437aa3c83716cd61481c6bb8c5057e6923e"),
									Keys: []*felt.Felt{
										utils.TestHexToFelt(t, "0x297be67eb977068ccd2304c6440368d4a6114929aeb860c98b6a7e91f96e2ef"),
										utils.TestHexToFelt(t, "0x436f6e74726f6c"),
									},
									Data: []*felt.Felt{
										utils.TestHexToFelt(t, "0x1"),
										utils.TestHexToFelt(t, "0xa250005"),
										utils.TestHexToFelt(t, "0x1"),
										utils.TestHexToFelt(t, "0xe52"),
									},
								},
								{
									FromAddress: utils.TestHexToFelt(t, "0x517567ac7026ce129c950e6e113e437aa3c83716cd61481c6bb8c5057e6923e"),
									Keys: []*felt.Felt{
										utils.TestHexToFelt(t, "0x1085a37d58e6a75db0dadc9bb9e6707ed9c5630aec61fdcdcd832decec751c0"),
									},
									Data: []*felt.Felt{
										utils.TestHexToFelt(t, "0x5"),
										utils.TestHexToFelt(t, "0xa25"),
										utils.TestHexToFelt(t, "0x1"),
										utils.TestHexToFelt(t, "0xe52"),
										utils.TestHexToFelt(t, "0x3543d2f0290e39a08cfdf2245f14aec7dca60672b7c7458375f3cb3834e1067"),
									},
								},
								{
									FromAddress: utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
									Keys: []*felt.Felt{
										utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
									},
									Data: []*felt.Felt{
										utils.TestHexToFelt(t, "0x3543d2f0290e39a08cfdf2245f14aec7dca60672b7c7458375f3cb3834e1067"),
										utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
										utils.TestHexToFelt(t, "0x1276c2c6477ed"),
										utils.TestHexToFelt(t, "0x0"),
									},
								},
							},
							ExecutionResources: ExecutionResources{
								ComputationResources{
									Steps:          34408,
									PedersenApps:   27,
									RangeCheckApps: 1077,
									ECOPApps:       3,
									PoseidonApps:   65,
								},
								DataAvailability{
									L1Gas:     0,
									L1DataGas: 256,
								},
							},
						},
					},
				},
			},
		},
	}

	var blockMainnet655660 = BlockWithReceipts{
		BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x7e53153d9737c3b60f917ae6df26b10bed5771ca2fce980c1cea9973e97ee7e"),
			ParentHash:       utils.TestHexToFelt(t, "0x2147e34aba1742219cb6a702476b55cd959bb70e44550ca9b9ce545125bac42"),
			BlockNumber:      655660,
			NewRoot:          utils.TestHexToFelt(t, "0x5c49cad0c4eac00da62a9c17f4b043973bbb58af5e94c59da7a215626559154"),
			Timestamp:        1720209856,
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			L1GasPrice: ResourcePrice{
				PriceInFRI: utils.TestHexToFelt(t, "0x151a2612eeb3"),
				PriceInWei: utils.TestHexToFelt(t, "0xedd8f5ff"),
			},
			L1DataGasPrice: ResourcePrice{
				PriceInFRI: utils.TestHexToFelt(t, "0xa7dc58"),
				PriceInWei: utils.TestHexToFelt(t, "0x764"),
			},
			L1DAMode:        L1DAModeBlob,
			StarknetVersion: "0.13.1.1",
		},
		"ACCEPTED_ON_L1",
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				// a lot of previous transactions ...
				{ // this is the last transaction of this block
					Transaction: BlockTransaction{
						BlockInvokeTxnV1{
							TransactionHash: utils.TestHexToFelt(t, "0x7dd8facd75bdebed2a76eb29dfa49172efea4913eea0abbc3e90a1af3d2c6ed"),
							InvokeTxnV1: InvokeTxnV1{
								Type:          "INVOKE",
								Version:       TransactionV1,
								Nonce:         utils.TestHexToFelt(t, "0x42"),
								MaxFee:        utils.TestHexToFelt(t, "0x286a5dccd4"),
								SenderAddress: utils.TestHexToFelt(t, "0x4dcbde5783e7131bd21c4114e22723ee0db79a5060933235fabea35049b766e"),
								Signature: []*felt.Felt{
									utils.TestHexToFelt(t, "0x1"),
									utils.TestHexToFelt(t, "0x2f9c7fa1b56357a343cf982b8f0beb0fddc8dcbac0e2bb65f4b98edb255ebcc"),
									utils.TestHexToFelt(t, "0x48b6a3060e267da25fc6d6ba4496f2380f11ece6e4231be2048a22ca273bee9"),
								},
								Calldata: []*felt.Felt{
									utils.TestHexToFelt(t, "0x1"),
									utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
									utils.TestHexToFelt(t, "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
									utils.TestHexToFelt(t, "0x3"),
									utils.TestHexToFelt(t, "0x73b7f95acad70fc8b2746062d9cf87c3e1600e0add99f469945b9d03d35637a"),
									utils.TestHexToFelt(t, "0x7180a5ffdd4fe"),
									utils.TestHexToFelt(t, "0x0"),
								},
							},
						},
					},
					Receipt: UnknownTransactionReceipt{
						TransactionReceipt: InvokeTransactionReceipt{
							Type:            "INVOKE",
							TransactionHash: utils.TestHexToFelt(t, "0x7dd8facd75bdebed2a76eb29dfa49172efea4913eea0abbc3e90a1af3d2c6ed"),
							ActualFee: FeePayment{
								Amount: utils.TestHexToFelt(t, "0x1bdf725ee2"),
								Unit:   UnitWei,
							},
							ExecutionStatus: TxnExecutionStatusSUCCEEDED,
							FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
							MessagesSent:    []MsgToL1{},
							Events: []Event{
								{
									FromAddress: utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
									Keys: []*felt.Felt{
										utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
									},
									Data: []*felt.Felt{
										utils.TestHexToFelt(t, "0x4dcbde5783e7131bd21c4114e22723ee0db79a5060933235fabea35049b766e"),
										utils.TestHexToFelt(t, "0x73b7f95acad70fc8b2746062d9cf87c3e1600e0add99f469945b9d03d35637a"),
										utils.TestHexToFelt(t, "0x7180a5ffdd4fe"),
										utils.TestHexToFelt(t, "0x0"),
									},
								},
								{
									FromAddress: utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
									Keys: []*felt.Felt{
										utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
									},
									Data: []*felt.Felt{
										utils.TestHexToFelt(t, "0x4dcbde5783e7131bd21c4114e22723ee0db79a5060933235fabea35049b766e"),
										utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
										utils.TestHexToFelt(t, "0x1bdf725ee2"),
										utils.TestHexToFelt(t, "0x0"),
									},
								},
							},
							ExecutionResources: ExecutionResources{
								ComputationResources{
									Steps:          11453,
									PedersenApps:   25,
									RangeCheckApps: 224,
									ECOPApps:       3,
								},
								DataAvailability{
									L1Gas:     0,
									L1DataGas: 192,
								},
							},
						},
					},
				},
			},
		},
	}

	var blockMock123 = BlockWithReceipts{
		BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "deadbeef"),
			ParentHash:       new(felt.Felt).SetUint64(1),
			BlockNumber:      1,
			NewRoot:          new(felt.Felt).SetUint64(1),
			Timestamp:        123,
			SequencerAddress: new(felt.Felt).SetUint64(1),
			L1GasPrice: ResourcePrice{
				PriceInFRI: new(felt.Felt).SetUint64(1),
				PriceInWei: new(felt.Felt).SetUint64(1),
			},
			L1DataGasPrice: ResourcePrice{
				PriceInFRI: new(felt.Felt).SetUint64(1),
				PriceInWei: new(felt.Felt).SetUint64(1),
			},
			L1DAMode:        L1DAModeBlob,
			StarknetVersion: "0.13",
		},
		"ACCEPTED_ON_L1",
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						BlockInvokeTxnV1{
							TransactionHash: utils.TestHexToFelt(t, "deadbeef"),
							InvokeTxnV1: InvokeTxnV1{
								Type:          "INVOKE",
								Version:       TransactionV1,
								Nonce:         new(felt.Felt).SetUint64(1),
								MaxFee:        new(felt.Felt).SetUint64(1),
								SenderAddress: utils.TestHexToFelt(t, "deadbeef"),
								Signature: []*felt.Felt{
									utils.TestHexToFelt(t, "deadbeef"),
								},
								Calldata: []*felt.Felt{
									new(felt.Felt).SetUint64(1),
								},
							},
						},
					},
					Receipt: UnknownTransactionReceipt{
						TransactionReceipt: InvokeTransactionReceipt{
							Type:            "INVOKE",
							TransactionHash: utils.TestHexToFelt(t, "deadbeef"),
							ActualFee: FeePayment{
								Amount: new(felt.Felt).SetUint64(1),
								Unit:   UnitWei,
							},
							ExecutionStatus: TxnExecutionStatusSUCCEEDED,
							FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
							MessagesSent:    []MsgToL1{},
							Events:          []Event{},
						},
					},
				},
			},
		},
	}

	var pendingBlockMock123 = PendingBlockWithReceipts{
		PendingBlockHeader{
			ParentHash:       new(felt.Felt).SetUint64(1),
			Timestamp:        123,
			SequencerAddress: new(felt.Felt).SetUint64(1),
			L1GasPrice: ResourcePrice{
				PriceInFRI: new(felt.Felt).SetUint64(1),
				PriceInWei: new(felt.Felt).SetUint64(1),
			},
			L1DataGasPrice: ResourcePrice{
				PriceInFRI: new(felt.Felt).SetUint64(1),
				PriceInWei: new(felt.Felt).SetUint64(1),
			},
			L1DAMode:        L1DAModeBlob,
			StarknetVersion: "0.13",
		},
		BlockBodyWithReceipts{
			Transactions: []TransactionWithReceipt{
				{
					Transaction: BlockTransaction{
						BlockInvokeTxnV1{
							TransactionHash: utils.TestHexToFelt(t, "deadbeef"),
							InvokeTxnV1: InvokeTxnV1{
								Type:          "INVOKE",
								Version:       TransactionV1,
								Nonce:         new(felt.Felt).SetUint64(1),
								MaxFee:        new(felt.Felt).SetUint64(1),
								SenderAddress: utils.TestHexToFelt(t, "deadbeef"),
								Signature: []*felt.Felt{
									utils.TestHexToFelt(t, "deadbeef"),
								},
								Calldata: []*felt.Felt{
									new(felt.Felt).SetUint64(1),
								},
							},
						},
					},
					Receipt: UnknownTransactionReceipt{
						TransactionReceipt: InvokeTransactionReceipt{
							Type:            "INVOKE",
							TransactionHash: utils.TestHexToFelt(t, "deadbeef"),
							ActualFee: FeePayment{
								Amount: new(felt.Felt).SetUint64(1),
								Unit:   UnitWei,
							},
							ExecutionStatus: TxnExecutionStatusSUCCEEDED,
							FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
							MessagesSent:    []MsgToL1{},
							Events:          []Event{},
						},
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
				ExpectedBlockWithReceipts: &blockSepolia64159,
			},
		},
		"mainnet": {
			{
				BlockID: WithBlockTag("pending"),
			},
			{
				BlockID:                   WithBlockNumber(655660),
				ExpectedBlockWithReceipts: &blockMainnet655660,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)

		result, err := testConfig.provider.BlockWithReceipts(context.Background(), test.BlockID)
		require.NoError(err, "Error in BlockWithReceipts")
		switch resultType := result.(type) {
		case *BlockWithReceipts:
			block, ok := result.(*BlockWithReceipts)
			require.True(ok, fmt.Sprintf("should return *BlockWithReceipts, instead: %T\n", result))
			require.True(strings.HasPrefix(block.BlockHash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.BlockHash)
			require.NotEmpty(block.Transactions, "the number of transactions should not be 0")

			if test.ExpectedBlockWithReceipts != nil {
				require.Equal(block.BlockHeader.BlockHash, test.ExpectedBlockWithReceipts.BlockHeader.BlockHash, "Error in BlockTxHash BlockHash")
				require.Equal(block.BlockHeader.ParentHash, test.ExpectedBlockWithReceipts.BlockHeader.ParentHash, "Error in BlockTxHash ParentHash")
				require.Equal(block.BlockHeader.Timestamp, test.ExpectedBlockWithReceipts.BlockHeader.Timestamp, "Error in BlockTxHash Timestamp")
				require.Equal(block.BlockHeader.SequencerAddress, test.ExpectedBlockWithReceipts.BlockHeader.SequencerAddress, "Error in BlockTxHash SequencerAddress")
				require.Equal(block.Status, test.ExpectedBlockWithReceipts.Status, "Error in BlockTxHash Status")
				require.Equal(block.Transactions[len(block.Transactions)-1], test.ExpectedBlockWithReceipts.Transactions[0], "Error in BlockTxHash Transactions")
			}
		case *PendingBlockWithReceipts:
			pBlock, ok := result.(*PendingBlockWithReceipts)
			require.True(ok, fmt.Sprintf("should return *PendingBlockWithReceipts, instead: %T\n", result))

			if testEnv == "mock" {
				require.Equal(pBlock.ParentHash, test.ExpectedPendingBlockWithReceipts.ParentHash, "Error in PendingBlockWithReceipts ParentHash")
				require.Equal(pBlock.SequencerAddress, test.ExpectedPendingBlockWithReceipts.SequencerAddress, "Error in PendingBlockWithReceipts SequencerAddress")
				require.Equal(pBlock.Timestamp, test.ExpectedPendingBlockWithReceipts.Timestamp, "Error in PendingBlockWithReceipts Timestamp")
				require.Equal(pBlock.Transactions, test.ExpectedPendingBlockWithReceipts.Transactions, "Error in PendingBlockWithReceipts Transactions")
			} else {
				require.NotEmpty(pBlock.ParentHash, "Error in PendingBlockWithReceipts ParentHash")
				require.NotEmpty(pBlock.SequencerAddress, "Error in PendingBlockWithReceipts SequencerAddress")
				require.NotEmpty(pBlock.Timestamp, "Error in PendingBlockWithReceipts Timestamp")
				require.NotEmpty(pBlock.Transactions, "Error in PendingBlockWithReceipts Transactions")
			}

		default:
			t.Fatalf("unexpected block type, found: %T\n", resultType)
		}
	}
}
