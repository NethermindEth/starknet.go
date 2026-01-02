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

// TestTransactionStatus tests the TransactionStatus function
func TestTransactionStatus(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

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

			resp, err := testConfig.Provider.TransactionStatus(t.Context(), test.TxnHash)
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
