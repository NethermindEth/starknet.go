package types

import (
	"encoding/json"
	"testing"
)

func TestTransactionStatus(t *testing.T) {
	for _, tc := range []struct {
		status string
		want   TransactionState
	}{{
		status: `"PENDING"`,
		want:   TransactionPending,
	}, {
		status: `"ACCEPTED_ON_L2"`,
		want:   TransactionAcceptedOnL2,
	}, {
		status: `"ACCEPTED_ON_L1"`,
		want:   TransactionAcceptedOnL1,
	}, {
		status: `"REJECTED"`,
		want:   TransactionRejected,
	}} {
		var tx TransactionState
		if err := json.Unmarshal([]byte(tc.status), &tx); err != nil {
			t.Fatalf("unmarshalling status want: %s", err)
		}
		if tx != tc.want {
			t.Fatalf("status error want %s, got %s", tc.want, tx)
		}
	}
}
