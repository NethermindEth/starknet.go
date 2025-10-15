package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Test the TxnStatus enum type
//
//nolint:dupl // The enum tests are similar, but with different enum values
func TestTxnStatusType(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	type testCase struct {
		Input         string
		Expected      TxnStatus
		ErrorExpected bool
	}

	testCases := []testCase{
		{
			Input:         `"active"`,
			Expected:      TxnActive,
			ErrorExpected: false,
		},
		{
			Input:         `"accepted"`,
			Expected:      TxnAccepted,
			ErrorExpected: false,
		},
		{
			Input:         `"dropped"`,
			Expected:      TxnDropped,
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

// Test the 'paymaster_trackingIdToLatestHash' method
func TestTrackingIdToLatestHash(t *testing.T) {
	// The AVNU paymaster does not support this method yet, so we can't have integration tests
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	expectedRawResp := `{
		"transaction_hash": "0xdeadbeef",
		"status": "active"
	}`

	var expectedResp TrackingIDResponse
	err := json.Unmarshal([]byte(expectedRawResp), &expectedResp)
	require.NoError(t, err)

	trackingID := internalUtils.DeadBeef

	pm := SetupMockPaymaster(t)
	pm.c.EXPECT().
		CallContextWithSliceArgs(
			context.Background(),
			gomock.AssignableToTypeOf(new(TrackingIDResponse)),
			"paymaster_trackingIdToLatestHash",
			trackingID,
		).
		SetArg(1, expectedResp).
		Return(nil)

	response, err := pm.TrackingIDToLatestHash(context.Background(), trackingID)
	require.NoError(t, err)
	assert.Equal(t, TxnActive, response.Status)
	assert.Equal(t, expectedResp.TransactionHash, response.TransactionHash)

	rawResp, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, expectedRawResp, string(rawResp))
}
