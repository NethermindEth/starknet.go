package rpc

import (
	_ "embed"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
func TestBlockID_Marshal(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	blockNumber := uint64(420)
	for _, test := range []struct {
		id      BlockID
		want    string
		wantErr error
	}{
		{
			id: BlockID{
				Tag: "latest",
			},
			want: `"latest"`,
		},
		{
			id: BlockID{
				Tag: "pre_confirmed",
			},
			want: `"pre_confirmed"`,
		},
		{
			id: BlockID{
				Tag: "l1_accepted",
			},
			want: `"l1_accepted"`,
		},
		{
			id: BlockID{
				Tag: "bad tag",
			},
			wantErr: ErrInvalidBlockID,
		},
		{
			id: BlockID{
				Number: &blockNumber,
			},
			want: `{"block_number":420}`,
		},
		{
			id: BlockID{
				Hash: internalUtils.TestHexToFelt(t, "0xdead"),
			},
			want: `{"block_hash":"0xdead"}`,
		},
	} {
		b, err := test.id.MarshalJSON()
		if test.wantErr != nil {
			require.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())

			return
		}
		require.NoError(t, err)

		assert.JSONEq(t, string(b), test.want)
	}
}
