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

// TestTransactionReceipt tests the TransactionReceipt function.
func TestTransactionReceipt(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

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

			txReceiptWithBlockInfo, err := testConfig.Provider.TransactionReceipt(
				t.Context(),
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
