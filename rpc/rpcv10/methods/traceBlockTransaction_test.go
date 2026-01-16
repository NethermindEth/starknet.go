package methods

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	. "github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestTraceBlockTransactions tests the TraceBlockTransactions function.
func TestTraceBlockTransactions(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.MockEnv)

	testConfig := internal.BeforeEach(t, false)

	type testSetType struct {
		BlockID     rpcv10.BlockID
		ExpectedErr error
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
			},
			{
				BlockID:     rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagPreConfirmed),
				// not the exact error, but it should contain it due to the checkForPreConfirmed() function
				ExpectedErr: rpcv10.ErrInvalidBlockID,
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: rpcv10.WithBlockNumber(99433),
			},
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
			},
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagL1Accepted),
			},
			{
				BlockID:     rpcv10.WithBlockHash(internalUtils.DeadBeef),
				ExpectedErr: ErrBlockNotFound,
			},
			{
				BlockID: rpcv10.WithBlockTag(rpcv10.BlockTagPreConfirmed),
				// not the exact error, but it should contain it due to the checkForPreConfirmed() function
				ExpectedErr: rpcv10.ErrInvalidBlockID,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(fmt.Sprintf("blockID: %v", test.BlockID), func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv && test.BlockID.Tag != rpcv10.BlockTagPreConfirmed {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_traceBlockTransactions",
						test.BlockID,
					).
					DoAndReturn(func(_, result, _ any, args ...any) error {
						rawResp := result.(*json.RawMessage)
						blockID := args[0].(rpcv10.BlockID)

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

			resp, err := TraceBlockTransactions(
				t.Context(),
				testConfig.Provider,
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
