package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestMessagesStatus tests the MessagesStatus function
func TestMessagesStatus(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TxHash      NumAsHex
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TxHash: "0x123",
			},
			{
				TxHash:      "0xdeadbeef",
				ExpectedErr: ErrHashNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				TxHash: "0x06c5ca541e3d6ce35134e1de3ed01dbf106eaa770d92744432b497f59fddbc00",
			},
			{
				TxHash:      "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				ExpectedErr: ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(string(test.TxHash), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getMessagesStatus",
						test.TxHash,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						txnHash := args[0].(NumAsHex)

						if txnHash == "0xdeadbeef" {
							return RPCError{
								Code:    29,
								Message: "Transaction hash not found",
							}
						}

						*rawResp = json.RawMessage(`
							[
								{
									"transaction_hash": "0x71660e0442b35d307fc07fa6007cf2ae4418d29fd73833303e7d3cfe1157157",
									"finality_status": "ACCEPTED_ON_L1",
									"execution_status": "SUCCEEDED"
								},
								{
									"transaction_hash": "0x28a3d1f30922ab86bb240f7ce0f5e8cbbf936e5d2fcfe52b8ffbe71e341640",
									"finality_status": "ACCEPTED_ON_L1",
									"execution_status": "SUCCEEDED"
								}
							]
						`)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.MessagesStatus(t.Context(), test.TxHash)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedErr.Error())

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
