package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringToByteArrFelt(t *testing.T) {
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
			out: []string{"0x1", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x0"},
		},
		{
			in:  "Long string, more than 31 charaLong string, more than 31 chara",
			out: []string{"0x2", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x4c6f6e6720737472696e672c206d6f7265207468616e203331206368617261", "0x0"},
		},
	}

	for _, tc := range tests {
		res, err := StringToByteArrFelt(tc.in)
		require.NoError(t, err, "error returned from StringToByteArrFelt")
		require.Len(t, res, len(tc.out), "invalid conversion: array sizes do not match")

		out, _ := HexArrToFelt(tc.out)
		require.Exactly(t, out, res, "invalid conversion: values do not match")
	}
}
