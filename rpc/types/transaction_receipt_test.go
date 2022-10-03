package types

import (
	"encoding/json"
	"testing"
)

func TestTransactionStatus(t *testing.T) {
	for _, tc := range []struct {
		status string
		want   TransactionStatus
	}{{
		status: `"PENDING"`,
		want:   TransactionStatus_Pending,
	}, {
		status: `"ACCEPTED_ON_L2"`,
		want:   TransactionStatus_AcceptedOnL2,
	}, {
		status: `"ACCEPTED_ON_L1"`,
		want:   TransactionStatus_AcceptedOnL1,
	}, {
		status: `"REJECTED"`,
		want:   TransactionStatus_Rejected,
	}} {
		tx := new(TransactionStatus)
		if err := json.Unmarshal([]byte(tc.status), tx); err != nil {
			t.Errorf("unmarshalling status want: %s", err)
		}
	}
}
