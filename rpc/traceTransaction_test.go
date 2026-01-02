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

// TestTraceTransaction tests the TraceTransaction function.
func TestTransactionTrace(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)
	testConfig := BeforeEach(t, false)

	type testSetType struct {
		TransactionHash *felt.Felt
		ExpectedError   error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				TransactionHash: internalUtils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
			},
			{
				TransactionHash: internalUtils.DeadBeef,
				ExpectedError:   ErrHashNotFound,
			},
			{
				TransactionHash: &felt.Zero,
				ExpectedError: &RPCError{
					Code:    10,
					Message: "No trace available for transaction",
					Data: &TraceStatusErrData{
						Status: TraceStatusReceived,
					},
				},
			},
		},
		tests.TestnetEnv: {
			{ // with 5 out of 6 fields (without state diff)
				TransactionHash: internalUtils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
			},
			{ // with 6 out of 6 fields
				TransactionHash: internalUtils.TestHexToFelt(t, "0x49d98a0328fee1de19d43d950cbaeb973d080d0c74c652523371e034cc0bbb2"),
			},
			{
				TransactionHash: internalUtils.DeadBeef,
				ExpectedError:   ErrHashNotFound,
			},
		},
		tests.IntegrationEnv: {
			{
				TransactionHash: internalUtils.TestHexToFelt(t, "0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019"),
			},
			{
				TransactionHash: internalUtils.DeadBeef,
				ExpectedError:   ErrHashNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.TransactionHash.String(), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_traceTransaction",
						test.TransactionHash,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						transactionHash := args[0].(*felt.Felt)

						if transactionHash == internalUtils.DeadBeef {
							return RPCError{
								Code:    29,
								Message: "Transaction hash not found",
							}
						}

						if *transactionHash == felt.Zero {
							return RPCError{
								Code:    10,
								Message: "No trace available for transaction",
								Data: &TraceStatusErrData{
									Status: TraceStatusReceived,
								},
							}
						}

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/trace/sepoliaInvokeTrace_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.TraceTransaction(
				t.Context(),
				test.TransactionHash,
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
