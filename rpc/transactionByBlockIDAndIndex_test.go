package rpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
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

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawTx, err := json.Marshal(tx)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawTx))
		})
	}
}
