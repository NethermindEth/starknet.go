package types

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	for _, tc := range []struct {
		hash string
	}{{
		hash: "0x0315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313",
	}, {
		hash: "0x008d6955e1bc0d5ba78b04630375f962ce6e5c91115173bc5f6e7843c6ee1269",
	}, {
		hash: "0x00680b0616e65633dfaf06d5a5ee5f1d1d4b641396009f00a67c0d18dc0f9638",
	}} {
		var th Hash
		if err := json.Unmarshal([]byte(strconv.Quote(tc.hash)), &th); err != nil {
			t.Fatalf("Unmarshalling text: %v", err)
		}
		h := th
		h2 := HexToHash(tc.hash)

		if h != h2 {
			t.Fatalf("Hashes not equal: %s %s", h, h2)
		}

		m, err := h.MarshalText()
		if err != nil {
			t.Fatalf("Marshalling text: %v", err)
		}

		m2, err := json.Marshal(h)
		if err != nil {
			t.Fatalf("Marshalling json: %v", err)
		}

		if tc.hash != string(m) {
			t.Errorf("Hash mismatch, want: %s, got: %s", tc.hash, m)
		}

		if strconv.Quote(tc.hash) != string(m2) {
			t.Errorf("Hash mismatch, want: %s, got: %s", tc.hash, m2)
		}
	}
}
