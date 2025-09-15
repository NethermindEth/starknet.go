package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			Input:         `"priority"`,
			Expected:      FeeModePriority,
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

func TestBuildTransaction(t *testing.T) {
	t.Parallel()
	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)

		t.Run("'deploy' transaction type", func(t *testing.T) {
			t.Parallel()
			// *** setup account data
			_, pubK, _ := account.GetRandomKeys()
			// OZ account class hash
			classHash := internalUtils.TestHexToFelt(t, "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")
			salt := internalUtils.RANDOM_FELT
			// precomputedAddress := account.PrecomputeAccountAddress(salt, classHash, []*felt.Felt{pubK})

			// *** build request
			reqBody := BuildTransactionRequest{
				Transaction: &UserTransaction{
					Type: UserTxnDeploy,
					Deployment: &AccDeploymentData{
						Address:             internalUtils.TestHexToFelt(t, "0x736b7c3fac1586518b55cccac1f675ca1bd0570d7354e2f2d23a0975a31f220"),
						ClassHash:           classHash,
						Salt:                salt,
						ConstructorCalldata: []*felt.Felt{pubK},
						SignatureData:       []*felt.Felt{internalUtils.RANDOM_FELT},
						Version:             2,
					},
				},
				Parameters: nil,
			}

			t.Run("'sponsored' fee mode", func(t *testing.T) {
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

			t.Run("'default' fee mode", func(t *testing.T) {
				t.Parallel()
				pm, _ := SetupPaymaster(t)

				request := reqBody
				request.Parameters = &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode:      FeeModeDefault,
						GasToken:  internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
						TipInStrk: internalUtils.TestHexToFelt(t, "0xfff"),
					},
				}

				_, err := pm.BuildTransaction(context.Background(), &request)
				require.Error(t, err) // it seems that the default fee mode is not supported for the 'deploy' transaction type
			})
		})
	})
}
