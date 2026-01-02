package rpc

import (
	"encoding/json"
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

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
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

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawTx, err := json.Marshal(tx)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawTx))
		})
	}
}
