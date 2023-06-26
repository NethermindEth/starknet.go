package caigo

import (
	"testing"

	"github.com/NethermindEth/caigo/types"
)

func TestGeneral_SplitFactStr(t *testing.T) {
	data := []map[string]string{
		{"input": "0x3", "h": "0x0", "l": "0x3"},
		{"input": "0x300000000000000000000000000000000", "h": "0x3", "l": "0x0"},
	}
	for _, d := range data {
		l, h := types.SplitFactStr(d["input"]) // 0x3
		if l != d["l"] {
			t.Errorf("expected %s, got %s", d["l"], l)
		}
		if h != d["h"] {
			t.Errorf("expected %s, got %s", d["h"], h)
		}
	}
}
