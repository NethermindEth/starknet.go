package utils

import (
	"testing"
)

func TestByteArrToFelt(t *testing.T) {
	var tests = []struct {
		in  string
		out []string
	}{
		{
			in:  "hello",
			out: []string{"0x0", "0x68656c6c6f", "0x5"},
		},
		{
			in:  "Long string, more than 31 characters.",
			out: []string{"0x1", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x63746572732e", "0x6"},
		},
		{
			in:  "Long string, more than 31 chara",
			out: []string{"0x1", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261"},
		},
		{
			in:  "Long string, more than 31 charaLong string, more than 31 chara",
			out: []string{"0x2", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261"},
		},
	}

	for _, tc := range tests {
		res, err := StringToByteArrFelt(tc.in)
		if err != nil {
			t.Fatalf("error in byte array conversion, err: %v", err)
		}

		if len(res) != len(tc.out) {
			t.Fatalf("error in byte array conversion, invalid length")
		}

		out, _ := HexArrToFelt(tc.out)
		for i, cmp := range out {
			if !cmp.Equal(res[i]) {
				t.Fatalf("invalid conversion, arr: %v", res)
			}
		}
	}
}
