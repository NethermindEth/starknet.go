package rpc

import (
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

// TestSimulateTransaction tests the SimulateTransaction function.
func TestSimulateTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)
	testConfig := BeforeEach(t, false)

	type simulateTxnInput struct {
		BlockID         BlockID          `json:"block_id"`
		Txns            []BroadcastTxn   `json:"transactions"`
		SimulationFlags []SimulationFlag `json:"simulation_flags"`
	}
	input := internalUtils.TestUnmarshalJSONFileToType[simulateTxnInput](
		t, "./testData/trace/sepoliaSimulateInvokeTx.json", "params")

	type testSetType struct {
		Description     string
		BlockID         BlockID
		Txns            []BroadcastTxn
		SimulationFlags []SimulationFlag
		ExpectedError   *RPCError
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description:     "valid call, all flags",
				BlockID:         input.BlockID,
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{SkipValidate, SkipFeeCharge},
			},
			{
				Description:     "block not found",
				BlockID:         WithBlockHash(internalUtils.DeadBeef),
				Txns:            input.Txns,
				SimulationFlags: input.SimulationFlags,
				ExpectedError:   ErrBlockNotFound,
			},
			{
				Description:     "exec error, pre confirmed",
				BlockID:         WithBlockTag(BlockTagPreConfirmed),
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{},
				ExpectedError:   ErrTxnExec, // due to invalid nonce
			},
		},
		tests.TestnetEnv: {
			{
				Description:     "valid call, no flags",
				BlockID:         input.BlockID,
				Txns:            input.Txns,
				SimulationFlags: input.SimulationFlags,
			},
			{
				Description:     "valid call, all flags",
				BlockID:         input.BlockID,
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{SkipValidate, SkipFeeCharge},
			},
			{
				Description:     "exec error, pre confirmed",
				BlockID:         WithBlockTag(BlockTagPreConfirmed),
				Txns:            input.Txns,
				SimulationFlags: []SimulationFlag{},
				ExpectedError:   ErrTxnExec, // due to invalid nonce
			},
			{
				Description:     "block not found",
				BlockID:         WithBlockHash(internalUtils.DeadBeef),
				Txns:            input.Txns,
				SimulationFlags: input.SimulationFlags,
				ExpectedError:   ErrBlockNotFound,
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
						"starknet_simulateTransactions",
						test.BlockID,
						test.Txns,
						test.SimulationFlags,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if blockID.Tag == BlockTagPreConfirmed {
							return RPCError{
								Code:    41,
								Message: "Transaction execution error",
								Data:    &TransactionExecErrData{},
							}
						}
						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/trace/sepoliaSimulateInvokeTxResp.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.SimulateTransactions(
				t.Context(),
				test.BlockID,
				test.Txns,
				test.SimulationFlags)

			if test.ExpectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)

				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

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

// TestTraceBlockTransactions tests the TraceBlockTransactions function.
func TestTraceBlockTransactions(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.MockEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		BlockID     BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID:     WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
				// not the exact error, but it should contain it due to the checkForPreConfirmed() function
				ExpectedErr: ErrInvalidBlockID,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockNumber(99433),
			},
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID:     WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
				// not the exact error, but it should contain it due to the checkForPreConfirmed() function
				ExpectedErr: ErrInvalidBlockID,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(fmt.Sprintf("blockID: %v", test.BlockID), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv && test.BlockID.Tag != BlockTagPreConfirmed {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_traceBlockTransactions",
						test.BlockID,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						*rawResp = internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/trace/sepoliaBlockTrace_0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.TraceBlockTransactions(
				t.Context(),
				test.BlockID,
			)
			if test.ExpectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.ExpectedErr.Error())

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
