package types

import (
	_ "embed"
	"encoding/json"
	"errors"
	"testing"

	ctypes "github.com/dontpanicdao/caigo/types"
)

func TestBlockID_Marshal(t *testing.T) {
	blockNumber := uint64(420)
	for _, tc := range []struct {
		id      BlockID
		want    string
		wantErr error
	}{{
		id: BlockID{
			Tag: "latest",
		},
		want: `"latest"`,
	}, {
		id: BlockID{
			Tag: "pending",
		},
		want: `"pending"`,
	}, {
		id: BlockID{
			Tag: "bad tag",
		},
		wantErr: ErrInvalidBlockID,
	}, {
		id: BlockID{
			Number: &blockNumber,
		},
		want: `{"block_number":420}`,
	}, {
		id: func() BlockID {
			h := ctypes.HexToHash("0xdead")
			return BlockID{
				Hash: &h,
			}
		}(),
		want: `{"block_hash":"0x000000000000000000000000000000000000000000000000000000000000dead"}`,
	}} {
		b, err := tc.id.MarshalJSON()
		if err != nil && tc.wantErr == nil {
			t.Errorf("marshalling block id: %v", err)
		} else if err != nil && !errors.Is(err, tc.wantErr) {
			t.Errorf("block error mismatch, want: %v, got: %v", tc.wantErr, err)
		}

		if string(b) != tc.want {
			t.Errorf("block id mismatch, want: %s, got: %s", tc.want, b)
		}
	}
}

func TestBlockStatus(t *testing.T) {
	for _, tc := range []struct {
		status string
		want   BlockStatus
	}{{
		status: `"PENDING"`,
		want:   BlockStatus_Pending,
	}, {
		status: `"ACCEPTED_ON_L2"`,
		want:   BlockStatus_AcceptedOnL2,
	}, {
		status: `"ACCEPTED_ON_L1"`,
		want:   BlockStatus_AcceptedOnL1,
	}, {
		status: `"REJECTED"`,
		want:   BlockStatus_Rejected,
	}} {
		tx := new(BlockStatus)
		if err := json.Unmarshal([]byte(tc.status), tx); err != nil {
			t.Errorf("unmarshalling status want: %s", err)
		}
	}
}

//go:embed testdata/block.json
var rawBlock []byte

func TestBlock_Unmarshal(t *testing.T) {
	b := Block{}
	if err := json.Unmarshal(rawBlock, &b); err != nil {
		t.Fatalf("Unmarshalling block: %v", err)
	}
}
