package types

import (
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
