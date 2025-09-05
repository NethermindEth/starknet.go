package paymaster

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the UserTxnType type
func TestUserTxnType(t *testing.T) {
	t.Parallel()
	type testCase struct {
		Input         string
		Expected      UserTxnType
		ErrorExpected bool
	}

	tests := []testCase{
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

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			t.Parallel()
			CompareEnumsHelper(t, test.Input, test.Expected, test.ErrorExpected)
		})
	}
}

// Test the FeeModeType type
func TestFeeModeType(t *testing.T) {
	t.Parallel()
	type testCase struct {
		Input         string
		Expected      FeeModeType
		ErrorExpected bool
	}

	tests := []testCase{
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

	for _, test := range tests {
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
		t.Parallel()
		// pm := SetupPaymaster(t)

		_, pubK, _ := account.GetRandomKeys()

		classHash := new(felt.Felt).SetBytes([]byte("0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"))
		salt := new(felt.Felt).SetBytes([]byte("0xdeadbeef"))

		precomputedAddress := account.PrecomputeAccountAddress(salt, classHash, []*felt.Felt{pubK})

		request := &BuildTransactionRequest{
			Transaction: &UserTransaction{
				Type: UserTxnDeploy,
				Deployment: &AccDeploymentData{
					Address:             precomputedAddress,
					ClassHash:           classHash,
					Salt:                salt,
					ConstructorCalldata: []*felt.Felt{pubK},
					Version:             2,
				},
			},
			Parameters: &UserParameters{
				Version: "0x1",
				FeeMode: FeeMode{
					Mode: FeeModeSponsored,
				},
			},
		}
		raw, err := json.Marshal(request)
		require.NoError(t, err)
		t.Logf("Raw: %s", string(raw))

		// tokens, err := pm.BuildTransaction(context.Background(), request)
		// require.NoError(t, err)
		// assert.NotNil(t, tokens)

	})
}
