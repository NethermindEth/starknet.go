package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/require"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
//
// The function tests the MarshalJSON method of the BlockID struct by providing
// different scenarios and verifying the output against the expected values.
// The scenarios include testing the serialization of the "latest" and
// "pending" tags, testing an invalid tag, testing the serialization of a
// block number, and testing the serialization of a block hash.
// The function uses the testing.T parameter to report any errors that occur
// during the execution of the test cases.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
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
			h, _ := new(felt.Felt).SetString("0xdead")
			return BlockID{
				Hash: h,
			}
		}(),
		want: `{"block_hash":"0xdead"}`,
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

// TestBlockStatus is a unit test for the BlockStatus function.
//
// The test checks the behavior of the BlockStatus function by iterating through a list of test cases.
//
// Parameters:
// - t: A testing.T object used for reporting test failures and logging.
// Returns:
//
//	none
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

//go:embed tests/block/block.json
var rawBlock []byte

// TestBlock_Unmarshal tests the Unmarshal function of the Block type.
//
// This test case unmarshals raw block data into a Block instance and verifies
// that there are no errors during the process. If any error occurs, the test
// fails with a fatal error message.
//
// Parameters:
// - t: the testing object for running the test
// Returns:
//
//	none
func TestBlock_Unmarshal(t *testing.T) {
	b := Block{}
	if err := json.Unmarshal(rawBlock, &b); err != nil {
		t.Fatalf("Unmarshalling block: %v", err)
	}
}

func TestBlockWithReceipts(t *testing.T) {
	provider := &Provider{c: &rpcMock{}}

	ctx := context.Background()

	type testSetType struct {
		BlockID                   BlockID
		ExpectedBlockWithReceipts BlockWithReceipts
		ExpectedErr               *RPCError
	}

	var expectedBlockWithReceipts struct {
		Result BlockWithReceipts `json:"result"`
	}
	read, err := os.ReadFile("tests/blockWithReceipts/integration332275.json")
	require.Nil(t, err)
	require.Nil(t, json.Unmarshal(read, &expectedBlockWithReceipts))

	testSet := map[string][]testSetType{
		"mock": {testSetType{
			BlockID:                   BlockID{Tag: "tests/blockWithReceipts/integration332275.json"},
			ExpectedBlockWithReceipts: expectedBlockWithReceipts.Result,
			ExpectedErr:               nil,
		},
		},
	}[testEnv]

	for _, test := range testSet {
		t.Run("BlockWithReceipts - block", func(t *testing.T) {

			block, err := provider.BlockWithReceipts(ctx, test.BlockID)
			require.Nil(t, err)
			blockCasted := block.(*BlockWithReceipts)
			require.Equal(t, test.ExpectedBlockWithReceipts, *blockCasted)

		})
	}
}
