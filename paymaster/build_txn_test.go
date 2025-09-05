package paymaster

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the UserTxnType type
func TestUserTxnType(t *testing.T) {
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
			CompareEnumsHelper(t, test.Input, test.Expected, test.ErrorExpected)
		})
	}
}

// Test the FeeModeType type
func TestFeeModeType(t *testing.T) {
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
