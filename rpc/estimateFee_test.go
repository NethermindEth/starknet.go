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

// TestEstimateFee tests the EstimateFee function.
func TestEstimateFee(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		description   string
		txs           []BroadcastTxn
		simFlags      []SimulationFlag
		blockID       BlockID
		expectedError *RPCError
	}

	sepoliaInvokeV3 := internalUtils.TestUnmarshalJSONFileToType[BroadcastInvokeTxnV3](
		t,
		"./testData/transactions/sepoliaInvokeV3_0x6035477af07a1b0a0186bec85287a6f629791b2f34b6e90eec9815c7a964f64.json",
	)
	invalidSepoliaInvokeV3 := sepoliaInvokeV3
	invalidSepoliaInvokeV3.Calldata = []*felt.Felt{internalUtils.DeadBeef}

	integrationInvokeV3 := internalUtils.TestUnmarshalJSONFileToType[BroadcastInvokeTxnV3](
		t,
		"./testData/transactions/integrationInvokeV3_0x38f7c9972f2b6f6d92d474cf605a077d154d58de938125180e7c87f22c5b019.json",
	)
	invalidIntegrationInvokeV3 := integrationInvokeV3
	invalidIntegrationInvokeV3.Calldata = []*felt.Felt{internalUtils.DeadBeef}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				description: "without flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags: []SimulationFlag{},
				blockID:  WithBlockTag(BlockTagLatest),
			},
			{
				description: "with flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags: []SimulationFlag{SkipValidate},
				blockID:  WithBlockTag(BlockTagLatest),
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					invalidSepoliaInvokeV3,
				},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				blockID:       WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				description: "normal call - without flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
			},
			{
				description: "normal call - with skip validate flag",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				simFlags:      []SimulationFlag{SkipValidate},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					invalidSepoliaInvokeV3,
				},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					sepoliaInvokeV3,
				},
				blockID:       WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
			// the contract_not_found error will not be tested since it's still not clear
			// when it should be returned (Pathfinder and Juno behave differently)
		},
		tests.IntegrationEnv: {
			{
				description: "without flag",
				txs: []BroadcastTxn{
					integrationInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(1_300_000),
				expectedError: nil,
			},
			{
				description: "with flag",
				txs: []BroadcastTxn{
					integrationInvokeV3,
				},
				simFlags:      []SimulationFlag{SkipValidate},
				blockID:       WithBlockNumber(1_300_000),
				expectedError: nil,
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					invalidIntegrationInvokeV3,
				},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					integrationInvokeV3,
				},
				blockID:       WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_estimateFee",
						test.txs,
						test.simFlags,
						test.blockID,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						txs := args[0].([]BroadcastTxn)
						blockID := args[2].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if txs[0].(BroadcastInvokeTxnV3).Calldata[0] == internalUtils.DeadBeef {
							return RPCError{
								Code:    41,
								Message: "Transaction execution error",
								Data:    &TransactionExecErrData{},
							}
						}

						*rawResp = json.RawMessage(`
							[
								{
									"l1_data_gas_consumed": "0x80",
									"l1_data_gas_price": "0x75da",
									"l1_gas_consumed": "0x0",
									"l1_gas_price": "0x1b709d15a1c6",
									"l2_gas_consumed": "0xc25b1",
									"l2_gas_price": "0xb2d05e00",
									"overall_fee": "0x87c1827e1eb00",
									"unit": "FRI"
								}
							]
						`)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.EstimateFee(
				t.Context(),
				test.txs,
				test.simFlags,
				test.blockID,
			)
			if test.expectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.expectedError.Code, rpcErr.Code)
				assert.Equal(t, test.expectedError.Message, rpcErr.Message)
				assert.IsType(t, rpcErr.Data, rpcErr.Data)

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
