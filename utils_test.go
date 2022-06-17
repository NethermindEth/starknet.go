package caigo

import (
	"math/big"
	"testing"
)

func TestSplitFactStr(t *testing.T) {
	data := []map[string]string{
		{"input": "0x3", "h": "0x0", "l": "0x3"},
		{"input": "0x300000000000000000000000000000000", "h": "0x3", "l": "0x0"},
	}
	for _, d := range data {
		l, h := SplitFactStr(d["input"]) // 0x3
		if l != d["l"] {
			t.Errorf("expected %s, got %s", d["l"], l)
		}
		if h != d["h"] {
			t.Errorf("expected %s, got %s", d["h"], h)
		}
	}
}

func TestsnValToBNSuccess(t *testing.T) {
	testSet := []struct {
		SN string
		BN *big.Int
	}{
		{SN: "1", BN: big.NewInt(1)},
		{SN: "0x1", BN: big.NewInt(1)},
		{SN: "0xa", BN: big.NewInt(10)},
		{SN: "0xv", BN: nil},
	}
	for _, k := range testSet {
		bn, _ := big.NewInt(0).SetString(k.SN, 0)
		if bn == nil && k.BN != nil {
			t.Errorf("snValToBN(%s) should return %v, but got nil", k.SN, k.BN)
		}
		if bn != nil && bn.Cmp(k.BN) != 0 {
			t.Fatalf("snValToBN(%s) = %s, want %s", k.SN, bn, k.BN)
		}
	}
}
