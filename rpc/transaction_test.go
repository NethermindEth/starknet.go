package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestTransactionByHash tests the TransactionByHash function.
func TestTransactionByHash(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TxHash        *felt.Felt
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TxHash: internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
			},
			{
				TxHash:        internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				TxHash: internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
			},
			{
				TxHash:        internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				TxHash: internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
			},
			{
				TxHash:        internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		t.Run(test.TxHash.String(), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getTransactionByHash",
						test.TxHash,
					).
					DoAndReturn(func(_, result, _ any, _ ...any) error {
						rawResp := result.(*json.RawMessage)

						if test.TxHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    29,
								Message: "Transaction hash not found",
							}
						}

						*rawResp = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/txnWithHash/sepoliaTxn.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			tx, err := testConfig.Provider.TransactionByHash(t.Context(), test.TxHash)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
		require.NoError(t, err)
		require.NotNil(t, tx)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawTx, err := json.Marshal(tx)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawTx))
		})
	}
}

// TestTransactionByBlockIdAndIndex tests the TransactionByBlockIdAndIndex function.
func TestTransactionByBlockIdAndIndex(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		BlockID       BlockID
		Index         uint64
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:   3,
			},
			{
				BlockID:       WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:         99999999999999999,
				ExpectedError: ErrInvalidTxnIndex,
			},
			{
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				Index:         3,
				ExpectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:   3,
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
				Index:   0,
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
				Index:   0,
			},
			{
				BlockID: WithBlockTag(BlockTagLatest),
				Index:   0,
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID: WithBlockNumber(1_300_000),
				Index:   0,
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		t.Run(fmt.Sprintf("Index: %d, BlockID: %v", test.Index, test.BlockID), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getTransactionByBlockIdAndIndex",
						test.BlockID,
						test.Index,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(BlockID)

						if blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if test.Index == 99999999999999999 {
							return RPCError{
								Code:    27,
								Message: "Invalid transaction index in a block",
							}
						}

						*rawResp = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/txnWithHash/sepoliaTxn.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			tx, err := testConfig.Provider.TransactionByBlockIDAndIndex(
				t.Context(),
				test.BlockID,
				test.Index,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			if err != nil {
				// in case the block has no transactions
				assert.EqualError(t, err, ErrInvalidTxnIndex.Error())

				return
			}

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawTx, err := json.Marshal(tx)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawTx))
		})
	}
}

func TestTransactionReceipt(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TransactionReceiptWithBlockInfo
	}

	receiptTxn52767_16 := *internalUtils.TestUnmarshalJSONFileToType[TransactionReceiptWithBlockInfo](t, "./testData/receipt/sepoliaRec_0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b.json", "")

	integrationInvokeV3Example := *internalUtils.TestUnmarshalJSONFileToType[TransactionReceiptWithBlockInfo](t, "./testData/receipt/integration_0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019.json", "result")

	// https://voyager.online/tx/0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519
	receiptL1Handler := *internalUtils.TestUnmarshalJSONFileToType[TransactionReceiptWithBlockInfo](t, "./testData/receipt/mainnetRc_0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519.json", "")

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: receiptTxn52767_16,
			},
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519"),
				ExpectedResp: receiptL1Handler,
			},
		},
		tests.TestnetEnv: {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: receiptTxn52767_16,
			},
		},
		tests.IntegrationEnv: {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
				ExpectedResp: integrationInvokeV3Example,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		txReceiptWithBlockInfo, err := testConfig.Provider.TransactionReceipt(
			context.Background(),
			test.TxnHash,
		)
		require.Nil(t, err)
		require.Equal(t, test.ExpectedResp, *txReceiptWithBlockInfo)
	}
}

// TestGetTransactionStatus tests starknet_getTransactionStatus in the GetTransactionStatus function
func TestGetTransactionStatus(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TxnStatusResult
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedResp: TxnStatusResult{FinalityStatus: TxnStatusAcceptedOnL1, ExecutionStatus: TxnExecutionStatusSUCCEEDED},
			},
			{
				TxnHash: internalUtils.TestHexToFelt(t, "0x5adf825a4b7fc4d2d99e65be934bd85c83ca2b9383f2ff28fc2a4bc2e6382fc"),
				ExpectedResp: TxnStatusResult{
					FinalityStatus:  TxnStatusAcceptedOnL1,
					ExecutionStatus: TxnExecutionStatusREVERTED,
					FailureReason:   "Transaction execution has failed:\n0: Error in the called contract (contract address: 0x036d67ab362562a97f9fba8a1051cf8e37ff1a1449530fb9f1f0e32ac2da7d06, class hash: 0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f, selector: 0x015d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad):\nError at pc=0:4835:\nCairo traceback (most recent call last):\nUnknown location (pc=0:67)\nUnknown location (pc=0:1835)\nUnknown location (pc=0:2554)\nUnknown location (pc=0:3436)\nUnknown location (pc=0:4040)\n\n1: Error in the called contract (contract address: 0x00000000000000000000000000000000000000000000000000000000ffffffff, class hash: 0x0000000000000000000000000000000000000000000000000000000000000000, selector: 0x02f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354):\nRequested contract address 0x00000000000000000000000000000000000000000000000000000000ffffffff is not deployed.\n",
				},
			},
		},
		tests.IntegrationEnv: {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
				ExpectedResp: TxnStatusResult{FinalityStatus: TxnStatusAcceptedOnL2, ExecutionStatus: TxnExecutionStatusSUCCEEDED},
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.Provider.TransactionStatus(context.Background(), test.TxnHash)
		require.Nil(t, err)
		require.Equal(t, resp.FinalityStatus, test.ExpectedResp.FinalityStatus)
		require.Equal(t, resp.ExecutionStatus, test.ExpectedResp.ExecutionStatus)
		require.Equal(t, resp.FailureReason, test.ExpectedResp.FailureReason)
	}
}

// TestGetMessagesStatus tests starknet_getMessagesStatus in the GetMessagesStatus function
func TestGetMessagesStatus(t *testing.T) {
	// TODO: add integration testcases
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TxHash       NumAsHex
		ExpectedResp []MessageStatus
		ExpectedErr  error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TxHash: "0x123",
				ExpectedResp: []MessageStatus{
					{
						Hash:            internalUtils.DeadBeef,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL2,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
					},
					{
						Hash:            internalUtils.DeadBeef,
						FinalityStatus:  TxnFinalityStatusAcceptedOnL2,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
					},
				},
			},
			{
				TxHash:      "0xdededededededededededededededededededededededededededededededede",
				ExpectedErr: ErrHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				TxHash: "0x06c5ca541e3d6ce35134e1de3ed01dbf106eaa770d92744432b497f59fddbc00",
				ExpectedResp: []MessageStatus{
					{
						Hash:            internalUtils.TestHexToFelt(t, "0x71660e0442b35d307fc07fa6007cf2ae4418d29fd73833303e7d3cfe1157157"),
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
					},
					{
						Hash:            internalUtils.TestHexToFelt(t, "0x28a3d1f30922ab86bb240f7ce0f5e8cbbf936e5d2fcfe52b8ffbe71e341640"),
						FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
						ExecutionStatus: TxnExecutionStatusSUCCEEDED,
					},
				},
			},
			{
				TxHash:      "0xdededededededededededededededededededededededededededededededede",
				ExpectedErr: ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(string(test.TxHash), func(t *testing.T) {
			resp, err := testConfig.Provider.MessagesStatus(context.Background(), test.TxHash)
			if test.ExpectedErr != nil {
				require.EqualError(t, err, test.ExpectedErr.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, test.ExpectedResp, resp)
			}
		})
	}
}
