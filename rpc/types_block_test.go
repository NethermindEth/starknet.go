package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBlockID_Marshal tests the MarshalJSON method of the BlockID struct.
//
// The function tests the MarshalJSON method of the BlockID struct by providing
// different scenarios and verifying the output against the expected values.
// The scenarios include testing the serialisation of the "latest" and
// "pre_confirmed" tags, testing an invalid tag, testing the serialisation of a
// block number, and testing the serialisation of a block hash.
// The function uses the testing.T parameter to report any errors that occur
// during the execution of the test cases.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestBlockID_Marshal(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

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
			Tag: "pre_confirmed",
		},
		want: `"pre_confirmed"`,
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
// The test checks the behaviour of the BlockStatus function by iterating through a list of test cases.
//
// Parameters:
//   - t: A testing.T object used for reporting test failures and logging.
//
// Returns:
//
//	none
func TestBlockStatus(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	for _, tc := range []struct {
		status string
		want   BlockStatus
	}{
		{
			status: `"PRE_CONFIRMED"`,
			want:   BlockStatusPreConfirmed,
		},
		{
			status: `"ACCEPTED_ON_L2"`,
			want:   BlockStatusAcceptedOnL2,
		},
		{
			status: `"ACCEPTED_ON_L1"`,
			want:   BlockStatusAcceptedOnL1,
		},
	} {
		tx := new(BlockStatus)
		if err := json.Unmarshal([]byte(tc.status), tx); err != nil {
			t.Errorf("unmarshalling status want: %s", err)
		}
	}
}

func TestBlockWithReceipts(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)
	provider := testConfig.Provider
	spy := tests.NewJSONRPCSpy(provider.c)
	provider.c = spy

	type testSetType struct {
		BlockID BlockID
	}

	// TODO: use these blocks for mock tests
	// var blockWithReceipt BlockWithReceipts
	// switch tests.TEST_ENV {
	// case tests.TestnetEnv:
	// 	blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/sepoliaBlockReceipts64159.json", "result")
	// case tests.MainnetEnv:
	// 	blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/mainnetBlockReceipts588763.json", "result")
	// case tests.IntegrationEnv:
	// 	blockWithReceipt = *internalUtils.TestUnmarshalJSONFileToType[BlockWithReceipts](t, "./testData/blockWithReceipts/integration1300000.json", "result")
	// }

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
		},
		tests.TestnetEnv: {
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
			{
				BlockID: WithBlockNumber(64159),
			},
		},
		tests.IntegrationEnv: {
			{
				BlockID: WithBlockTag(BlockTagL1Accepted),
			},
			{
				BlockID: WithBlockTag(BlockTagLatest),
			},
			{
				BlockID: WithBlockTag(BlockTagPreConfirmed),
			},
			{
				BlockID: WithBlockNumber(1_300_000),
			},
		},
		tests.MainnetEnv: {
			{
				BlockID: WithBlockTag("pre_confirmed"),
			},
			{
				BlockID: WithBlockNumber(588763),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		blockID, _ := test.BlockID.MarshalJSON()
		t.Run(string(blockID), func(t *testing.T) {
			result, err := provider.BlockWithReceipts(context.Background(), test.BlockID)
			require.NoError(t, err, "Error in BlockWithReceipts")

			rawExpectedBlock := spy.LastResponse()

			switch block := result.(type) {
			case *BlockWithReceipts:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			case *PreConfirmedBlockWithReceipts:
				rawBlock, err := json.Marshal(block)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedBlock), string(rawBlock))
			default:
				t.Fatalf("unexpected block type, found: %T\n", block)
			}
		})
	}
}
