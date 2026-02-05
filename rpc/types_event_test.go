package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddressList_MarshalJSON tests the MarshalJSON method of the AddressList type.
func TestAddressList_MarshalJSON(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	for _, test := range []struct {
		name string
		al   AddressList
		want string
	}{
		{
			name: "single address marshals as string",
			al:   AddressList{internalUtils.TestHexToFelt(t, "0x1234")},
			want: `"0x1234"`,
		},
		{
			name: "multiple addresses marshal as array",
			al: AddressList{
				internalUtils.TestHexToFelt(t, "0x1234"),
				internalUtils.TestHexToFelt(t, "0x5678"),
			},
			want: `["0x1234","0x5678"]`,
		},
		{
			name: "empty list marshals as empty array",
			al:   AddressList{},
			want: `[]`,
		},
		{
			name: "nil list marshals as null",
			al:   nil,
			want: `null`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.al)
			require.NoError(t, err)
			assert.JSONEq(t, test.want, string(b))
		})
	}
}

// TestAddressList_UnmarshalJSON tests the UnmarshalJSON method of the AddressList type.
func TestAddressList_UnmarshalJSON(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	for _, test := range []struct {
		name    string
		input   string
		want    AddressList
		wantErr bool
	}{
		{
			name:  "single address as string",
			input: `"0x1234"`,
			want:  AddressList{internalUtils.TestHexToFelt(t, "0x1234")},
		},
		{
			name:  "multiple addresses as array",
			input: `["0x1234","0x5678"]`,
			want: AddressList{
				internalUtils.TestHexToFelt(t, "0x1234"),
				internalUtils.TestHexToFelt(t, "0x5678"),
			},
		},
		{
			name:  "empty array",
			input: `[]`,
			want:  AddressList{},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var al AddressList
			err := json.Unmarshal([]byte(test.input), &al)
			if test.wantErr {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)
			require.Len(t, al, len(test.want))
			for i := range al {
				assert.Equal(t, test.want[i], al[i])
			}
		})
	}
}

// TestAddressList_RoundTrip tests that marshalling and unmarshalling produces the same result.
func TestAddressList_RoundTrip(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	for _, test := range []struct {
		name string
		al   AddressList
	}{
		{
			name: "single address round trip",
			al:   AddressList{internalUtils.TestHexToFelt(t, "0xdeadbeef")},
		},
		{
			name: "multiple addresses round trip",
			al: AddressList{
				internalUtils.TestHexToFelt(t, "0x1"),
				internalUtils.TestHexToFelt(t, "0x2"),
				internalUtils.TestHexToFelt(t, "0x3"),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Marshal
			b, err := json.Marshal(test.al)
			require.NoError(t, err)

			// Unmarshal
			var result AddressList
			err = json.Unmarshal(b, &result)
			require.NoError(t, err)

			// Compare
			require.Len(t, result, len(test.al))
			for i := range result {
				assert.Equal(t, test.al[i], result[i])
			}
		})
	}
}
