package types

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestTransaction(t *testing.T) {
	hash := HexToHash("0xdead")
	th := TransactionHash{hash}
	b, err := json.Marshal(th)
	if err != nil {
		t.Fatalf("marshalling transaction hash: %v", err)
	}

	marshalled, err := hash.MarshalText()
	if err != nil {
		t.Fatalf("marshalling transaction hash: %v", err)
	}

	if string(b) != strconv.Quote(string(marshalled)) {
		t.Fatalf("Marshalled hash mismatch, want: %s, got: %s", marshalled, b)
	}
}
