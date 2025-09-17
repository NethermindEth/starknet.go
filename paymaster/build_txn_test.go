package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var STRKContractAddress, _ = internalUtils.HexToFelt("0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D")

// Test the UserTxnType type
//
//nolint:dupl
func TestUserTxnType(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	type testCase struct {
		Input         string
		Expected      UserTxnType
		ErrorExpected bool
	}

	testCases := []testCase{
		{
			Input:         `"deploy"`,
			Expected:      UserTxnDeploy,
			ErrorExpected: false,
		},
		{
			Input:         `"invoke"`,
			Expected:      UserTxnInvoke,
			ErrorExpected: false,
		},
		{
			Input:         `"deploy_and_invoke"`,
			Expected:      UserTxnDeployAndInvoke,
			ErrorExpected: false,
		},
		{
			Input:         `"unknown"`,
			ErrorExpected: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.Input, func(t *testing.T) {
			t.Parallel()
			CompareEnumsHelper(t, test.Input, test.Expected, test.ErrorExpected)
		})
	}
}

// Test the FeeModeType type
//
//nolint:dupl
func TestFeeModeType(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	type testCase struct {
		Input         string
		Expected      FeeModeType
		ErrorExpected bool
	}

	testCases := []testCase{
		{
			Input:         `"default"`,
			Expected:      FeeModeDefault,
			ErrorExpected: false,
		},
		{
			Input:         `"sponsored"`,
			Expected:      FeeModeSponsored,
			ErrorExpected: false,
		},
		{
			Input:         `"unknown"`,
			ErrorExpected: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.Input, func(t *testing.T) {
			t.Parallel()
			CompareEnumsHelper(t, test.Input, test.Expected, test.ErrorExpected)
		})
	}
}

// Test the UserParamVersion type
//

func TestUserParamVersion(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	type testCase struct {
		Input         string
		Expected      UserParamVersion
		ErrorExpected bool
	}

	testCases := []testCase{
		{
			Input:         `"0x1"`,
			Expected:      UserParamV1,
			ErrorExpected: false,
		},
		{
			Input:         `"0x2"`,
			ErrorExpected: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.Input, func(t *testing.T) {
			t.Parallel()
			CompareEnumsHelper(t, test.Input, test.Expected, test.ErrorExpected)
		})
	}
}

// CompareEnumsHelper compares an enum type with the expected value and error expected.
func CompareEnumsHelper[T any](t *testing.T, input string, expected T, errorExpected bool) {
	t.Helper()

	var actual T
	err := json.Unmarshal([]byte(input), &actual)
	if errorExpected {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		marshalled, err := json.Marshal(actual)
		assert.NoError(t, err)
		assert.Equal(t, input, string(marshalled))
	}
}

// Test the BuildTransaction method with different transaction types and fee modes.
func TestBuildTransaction(t *testing.T) {
	t.Parallel()

	// *** setup for deploy type transactions
	classHash := internalUtils.TestHexToFelt(
		t,
		"0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f", // OZ account class hash
	)

	deploymentData := &AccDeploymentData{
		Address:             internalUtils.TestHexToFelt(t, "0x736b7c3fac1586518b55cccac1f675ca1bd0570d7354e2f2d23a0975a31f220"),
		ClassHash:           classHash,
		Salt:                internalUtils.RANDOM_FELT,
		ConstructorCalldata: []*felt.Felt{internalUtils.RANDOM_FELT},
		SignatureData:       []*felt.Felt{internalUtils.RANDOM_FELT},
		Version:             2,
	}

	// *** setup for invoke type transactions
	accountAddress := internalUtils.TestHexToFelt(t, "0x5c74db20fa8f151bfd3a7a462cf2e8d4578a88aa4bd7a1746955201c48d8e5e")
	transferAmount, _ := internalUtils.HexToU256Felt("0xfff")

	invokeData := &UserInvoke{
		UserAddress: accountAddress,
		Calls: []Call{
			{
				To:       STRKContractAddress,
				Selector: internalUtils.GetSelectorFromNameFelt("transfer"),
				Calldata: append([]*felt.Felt{accountAddress}, transferAmount...),
			},
			{
				// same ERC20 contract as in examples/simpleInvoke
				To:       internalUtils.TestHexToFelt(t, "0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54"),
				Selector: internalUtils.GetSelectorFromNameFelt("mint"),
				Calldata: []*felt.Felt{new(felt.Felt).SetUint64(10000), &felt.Zero},
			},
		},
	}

	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)

		t.Run("deploy transaction type", func(t *testing.T) {
			t.Parallel()
			// *** build request
			reqBody := BuildTransactionRequest{
				Transaction: &UserTransaction{
					Type:       UserTxnDeploy,
					Deployment: deploymentData,
				},
				Parameters: nil,
			}

			t.Run("sponsored fee mode", func(t *testing.T) {
				t.Parallel()
				pm, spy := SetupPaymaster(t)

				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
					},
				}

				resp, err := pm.BuildTransaction(context.Background(), &request)
				require.NoError(t, err)

				rawResp, err := json.Marshal(resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
			})

			t.Run("default fee mode", func(t *testing.T) {
				t.Parallel()
				pm, _ := SetupPaymaster(t)

				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode:     FeeModeDefault,
						GasToken: STRKContractAddress,
						Tip:      nil,
					},
				}

				_, err := pm.BuildTransaction(context.Background(), &request)
				require.Error(t, err) // it seems that the default fee mode is not supported for the 'deploy' transaction type
			})
		})

		t.Run("invoke transaction type", func(t *testing.T) {
			t.Parallel()

			reqBody := BuildTransactionRequest{
				Transaction: &UserTransaction{
					Type:   UserTxnInvoke,
					Invoke: invokeData,
				},
				Parameters: nil,
			}

			t.Run("sponsored fee mode - with nil tip", func(t *testing.T) {
				t.Parallel()
				pm, spy := SetupPaymaster(t)

				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
					},
				}

				resp, err := pm.BuildTransaction(context.Background(), &request)
				require.NoError(t, err)
				// The default tip priority is normal
				assert.Equal(t, TipPriorityNormal, resp.Parameters.FeeMode.Tip.Priority)

				rawResp, err := json.Marshal(resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
			})

			t.Run("default fee mode - with custom tip", func(t *testing.T) {
				t.Parallel()
				pm, spy := SetupPaymaster(t)

				customTip := uint64(1000)
				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode:     FeeModeDefault,
						GasToken: STRKContractAddress,
						Tip: &TipPriority{
							Custom: &customTip,
						},
					},
				}

				resp, err := pm.BuildTransaction(context.Background(), &request)
				require.NoError(t, err)

				assert.Equal(t, customTip, *resp.Parameters.FeeMode.Tip.Custom)

				rawResp, err := json.Marshal(resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
			})
		})

		t.Run("deploy-and-invoke transaction type", func(t *testing.T) {
			t.Parallel()

			reqBody := BuildTransactionRequest{
				Transaction: &UserTransaction{
					Type:   UserTxnDeployAndInvoke,
					Invoke: invokeData,
				},
				Parameters: nil,
			}
			reqBody.Transaction.Deployment = deploymentData
			reqBody.Transaction.Invoke = invokeData

			t.Run("sponsored fee mode - with slow tip", func(t *testing.T) {
				t.Parallel()
				pm, spy := SetupPaymaster(t)

				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
						Tip: &TipPriority{
							Priority: TipPrioritySlow,
						},
					},
				}

				resp, err := pm.BuildTransaction(context.Background(), &request)
				require.NoError(t, err)

				assert.Equal(t, TipPrioritySlow, resp.Parameters.FeeMode.Tip.Priority)

				rawResp, err := json.Marshal(resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
			})

			t.Run("default fee mode - with fast tip", func(t *testing.T) {
				t.Parallel()
				pm, spy := SetupPaymaster(t)

				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode:     FeeModeDefault,
						GasToken: STRKContractAddress,
						Tip: &TipPriority{
							Priority: TipPriorityFast,
						},
					},
				}

				resp, err := pm.BuildTransaction(context.Background(), &request)
				require.NoError(t, err)

				assert.Equal(t, TipPriorityFast, resp.Parameters.FeeMode.Tip.Priority)

				rawResp, err := json.Marshal(resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
			})
		})
	})

	t.Run("mock", func(t *testing.T) {
		tests.RunTestOn(t, tests.MockEnv)

		t.Run("deploy transaction type - sponsored fee mode", func(t *testing.T) {
			t.Parallel()
			// *** build request
			request := BuildTransactionRequest{
				Transaction: &UserTransaction{
					Type:       UserTxnDeploy,
					Deployment: deploymentData,
				},
				Parameters: &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
					},
				},
			}

			// *** assert the request marshalled is equal to the expected request
			expectedReqs := *internalUtils.TestUnmarshalJSONFileToType[[]json.RawMessage](t, "testdata/build_txn/deploy-request.json", "params")
			expectedReq := expectedReqs[0]

			rawReq, err := json.Marshal(request)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedReq), string(rawReq))

			// *** assert the response marshalled is equal to the expected response
			expectedResp := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](t, "testdata/build_txn/deploy-response.json", "result")

			var response BuildTransactionResponse
			err = json.Unmarshal(expectedResp, &response)
			require.NoError(t, err)

			pm := SetupMockPaymaster(t)
			pm.c.EXPECT().CallContextWithSliceArgs(
				context.Background(),
				gomock.AssignableToTypeOf(new(BuildTransactionResponse)),
				"paymaster_buildTransaction",
				&request,
			).Return(nil).
				SetArg(1, response)

			resp, err := pm.BuildTransaction(context.Background(), &request)
			require.NoError(t, err)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedResp), string(rawResp))
		})

		t.Run("invoke transaction type - default fee mode with custom tip", func(t *testing.T) {
			t.Parallel()
			// *** build request
			customTip := uint64(1000)

			request := BuildTransactionRequest{
				Transaction: &UserTransaction{
					Type:   UserTxnInvoke,
					Invoke: invokeData,
				},
				Parameters: &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode:     FeeModeDefault,
						GasToken: STRKContractAddress,
						Tip: &TipPriority{
							Custom: &customTip,
						},
					},
				},
			}

			// *** assert the request marshalled is equal to the expected request
			expectedReqs := *internalUtils.TestUnmarshalJSONFileToType[[]json.RawMessage](t, "testdata/build_txn/invoke-request.json", "params")
			expectedReq := expectedReqs[0]

			rawReq, err := json.Marshal(request)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedReq), string(rawReq))

			// *** assert the response marshalled is equal to the expected response
			expectedResp := *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](t, "testdata/build_txn/invoke-response.json", "result")

			var response BuildTransactionResponse
			err = json.Unmarshal(expectedResp, &response)
			require.NoError(t, err)

			pm := SetupMockPaymaster(t)
			pm.c.EXPECT().CallContextWithSliceArgs(
				context.Background(),
				gomock.AssignableToTypeOf(new(BuildTransactionResponse)),
				"paymaster_buildTransaction",
				&request,
			).Return(nil).
				SetArg(1, response)

			resp, err := pm.BuildTransaction(context.Background(), &request)
			require.NoError(t, err)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedResp), string(rawResp))
		})
	})
}
