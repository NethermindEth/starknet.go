package rpc

import "testing"

func TestTxnStatus(t *testing.T) {
	for _, tc := range []struct {
		status string
		want   TxnStatus
	}{{
		status: "PENDING",
		want:   TxnStatus_Pending,
	}, {
		status: "ACCEPTED_ON_L2",
		want:   TxnStatus_AcceptedOnL2,
	}, {
		status: "ACCEPTED_ON_L1",
		want:   TxnStatus_AcceptedOnL1,
	}, {
		status: "REJECTED",
		want:   TxnStatus_Rejected,
	}} {
		tx := new(TxnStatus)
		if err := tx.UnmarshalJSON([]byte(tc.status)); err != nil {
			t.Errorf("unmarshalling status want: %s", err)
		}
	}
}
