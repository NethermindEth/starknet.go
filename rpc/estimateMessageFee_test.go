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

// TestEstimateMessageFee tests the EstimateMessageFee function.
func TestEstimateMessageFee(t *testing.T) {
	// TODO: add integration testcase
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description string
		MsgFromL1
		BlockID
		ExpectedError *RPCError
	}

	// https://sepolia.voyager.online/message/0x273f4e20fc522098a60099e5872ab3deeb7fb8321a03dadbd866ac90b7268361
	l1Handler := MsgFromL1{
		FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
		ToAddress: internalUtils.TestHexToFelt(
			t,
			"0x04c5772d1914fe6ce891b64eb35bf3522aeae1315647314aac58b01137607f3f",
		),
		Selector: internalUtils.TestHexToFelt(
			t,
			"0x1b64b1b3b690b43b9b514fb81377518f4039cd3e4f4914d8a6bdf01d679fb19",
		),
		Payload: internalUtils.TestHexArrToFelt(t, []string{
			"0x455448",
			"0x2f14d277fc49e0e2d2967d019aea8d6bd9cb3998",
			"0x02000e6213e24b84012b1f4b1cbd2d7a723fb06950aeab37bedb6f098c7e051a",
			"0x01a055690d9db80000",
			"0x00",
		}),
	}

	l1HandlerInvalidSelector := l1Handler
	l1HandlerInvalidSelector.Selector = internalUtils.DeadBeef

	l1HandlerInvalidToAddress := l1Handler
	l1HandlerInvalidToAddress.ToAddress = internalUtils.DeadBeef

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "normal call",
				MsgFromL1:   l1Handler,
				BlockID:     WithBlockTag(BlockTagLatest),
			},
			{
				Description:   "contract error",
				MsgFromL1:     l1HandlerInvalidSelector,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractError,
			},
			{
				Description:   "contract not found",
				MsgFromL1:     l1HandlerInvalidToAddress,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractNotFound,
			},
			{
				Description:   "invalid block",
				MsgFromL1:     l1Handler,
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call",
				MsgFromL1:   l1Handler,
				BlockID:     WithBlockTag(BlockTagLatest),
			},
			{
				Description:   "contract error",
				MsgFromL1:     l1HandlerInvalidSelector,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractError,
			},
			{
				Description:   "contract not found",
				MsgFromL1:     l1HandlerInvalidToAddress,
				BlockID:       WithBlockTag(BlockTagLatest),
				ExpectedError: ErrContractNotFound,
			},
			{
				Description:   "invalid block",
				MsgFromL1:     l1Handler,
				BlockID:       WithBlockHash(internalUtils.DeadBeef),
				ExpectedError: ErrBlockNotFound,
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
						"starknet_estimateMessageFee",
						test.MsgFromL1,
						test.BlockID,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						msgFromL1 := args[0].(MsgFromL1)
						blockID := args[1].(BlockID)

						if blockID.Hash != nil && blockID.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if msgFromL1.ToAddress == internalUtils.DeadBeef {
							return RPCError{
								Code:    20,
								Message: "Contract not found",
							}
						}

						if msgFromL1.Selector == internalUtils.DeadBeef {
							return RPCError{
								Code:    40,
								Message: "Contract error",
								Data:    &ContractErrData{},
							}
						}

						*rawResp = json.RawMessage(`
							{
								"l1_gas_consumed": "0x4ed3",
								"l1_gas_price": "0x7e15d2b5",
								"l2_gas_consumed": "0x0",
								"l2_gas_price": "0x1",
								"l1_data_gas_consumed": "0x80",
								"l1_data_gas_price": "0x1",
								"overall_fee": "0x26d2922fd1af",
								"unit": "WEI"
							}
						`)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.EstimateMessageFee(
				t.Context(),
				test.MsgFromL1,
				test.BlockID,
			)
			if test.ExpectedError != nil {
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

				return
			}
			require.NoError(t, err)
			rawExpectedFeeEst := testConfig.RPCSpy.LastResponse()

			rawFeeEst, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedFeeEst), string(rawFeeEst))
		},
		)
	}
}
