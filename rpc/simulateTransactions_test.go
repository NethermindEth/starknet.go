package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
