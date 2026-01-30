package rpcv10

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestTransactionByBlockIdAndIndex tests the TransactionByBlockIdAndIndex function.
func TestTransactionByBlockIdAndIndex(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		blockID       types.BlockID
		Index         uint64
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				blockID: types.WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:   3,
			},
			{
				blockID:       types.WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:         99999999999999999,
				ExpectedError: ErrInvalidTxnIndex,
			},
			{
				blockID:       types.WithBlockHash(internalUtils.DeadBeef),
				Index:         3,
				ExpectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				blockID: types.WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:   3,
			},
			{
				blockID: types.WithBlockTag(types.BlockTagPreConfirmed),
				Index:   0,
			},
			{
				blockID: types.WithBlockTag(types.BlockTagL1Accepted),
				Index:   0,
			},
			{
				blockID: types.WithBlockTag(types.BlockTagLatest),
				Index:   0,
			},
		},
		tests.IntegrationEnv: {
			{
				blockID: types.WithBlockNumber(1_300_000),
				Index:   0,
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		t.Run(fmt.Sprintf("Index: %d, BlockID: %v", test.Index, test.blockID), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getTransactionByBlockIdAndIndex",
						test.blockID,
						test.Index,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(types.BlockID)

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

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/txnWithHash/sepoliaTxn.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			tx, err := TransactionByBlockIDAndIndex(
				t.Context(),
				testConfig.Provider,
				test.blockID,
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

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawTx, err := json.Marshal(tx)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawTx))
		})
	}
}

// TestTransactionByHash tests the TransactionByHash function.
func TestTransactionByHash(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := internal.BeforeEach(t, false)

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

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/txnWithHash/sepoliaTxn.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			tx, err := TransactionByHash(
				t.Context(),
				testConfig.Provider,
				test.TxHash,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)
			require.NotNil(t, tx)

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawTx, err := json.Marshal(tx)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawTx))
		})
	}
}

// TestTransactionReceipt tests the TransactionReceipt function.
func TestTransactionReceipt(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		TxnHash       *felt.Felt
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TxnHash: internalUtils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
			},
			{
				TxnHash:       internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				TxnHash: internalUtils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
			},
			{
				TxnHash:       internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				TxnHash: internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
			},
			{
				TxnHash:       internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.TxnHash.String(), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getTransactionReceipt",
						test.TxnHash,
					).
					DoAndReturn(func(_, result, _ any, _ ...any) error {
						rawResp := result.(*json.RawMessage)

						if test.TxnHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    29,
								Message: "Transaction hash not found",
							}
						}

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/receipt/sepoliaReceipt.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			txReceiptWithBlockInfo, err := GetTransactionReceipt(
				t.Context(),
				testConfig.Provider,
				test.TxnHash,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawReceipt, err := json.Marshal(txReceiptWithBlockInfo)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawReceipt))
		})
	}
}

// TestTransactionStatus tests the TransactionStatus function
func TestTransactionStatus(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		Description   string
		TxnHash       *felt.Felt
		ExpectedError error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "only with FinalityStatus and ExecutionStatus",
				TxnHash:     internalUtils.TestHexToFelt(t, "0x1"),
			},
			{
				Description: "with FailureReason",
				TxnHash:     internalUtils.TestHexToFelt(t, "0x2"),
			},
			{
				Description:   "error - hash not found",
				TxnHash:       internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "only with FinalityStatus and ExecutionStatus",
				TxnHash:     internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
			},
			{
				Description: "with FailureReason",
				TxnHash:     internalUtils.TestHexToFelt(t, "0x5adf825a4b7fc4d2d99e65be934bd85c83ca2b9383f2ff28fc2a4bc2e6382fc"),
			},
			{
				Description:   "error - hash not found",
				TxnHash:       internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				Description: "only with FinalityStatus and ExecutionStatus",
				TxnHash:     internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
			},
			{
				Description:   "error - hash not found",
				TxnHash:       internalUtils.DeadBeef,
				ExpectedError: ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getTransactionStatus",
						test.TxnHash,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						txnHash := args[0].(*felt.Felt)

						if txnHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    29,
								Message: "Transaction hash not found",
							}
						}

						if txnHash.String() == "0x1" {
							*rawResp = json.RawMessage(`
								{
									"finality_status": "ACCEPTED_ON_L2",
									"execution_status": "SUCCEEDED"
								}
							`)

							return nil
						}
						if txnHash.String() == "0x2" {
							*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"./testData/txnStatus/sepoliaStatus.json",
								"result",
							)

							return nil
						}

						return nil
					}).
					Times(1)
			}

			resp, err := TransactionStatus(
				t.Context(),
				testConfig.Provider,
				test.TxnHash,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}
