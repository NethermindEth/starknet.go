package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

// TestEvents is a test function for testing the Events function.
//
// It creates a test configuration and defines a test set with different scenarios for the Events function.
// The test set includes a "mock" scenario and a "mainnet" scenario.
// In the "mock" scenario, it sets up the event filter, result page request, and expected response.
// In the "mainnet" scenario, it sets up the event filter, result page request, and expected response.
// It then iterates through the test set and performs the following steps for each test:
// - Creates a spy object.
// - Sets the provider's context to the spy object.
// - Sets up the event input with the event filter and result page request.
// - Calls the Events function with the event input.
// - Checks if there is an error and fails the test if there is.
// - Compares the events' block hash, block number, and transaction hash with the expected response.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestEvents(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		eventFilter  EventFilter
		resPageReq   ResultPageRequest
		expectedResp EventChunk
	}
	var fromNum uint64 = 1471
	var toNum uint64 = 1473

	testSet := map[string][]testSetType{
		"mock": {{
			eventFilter: EventFilter{},
			resPageReq: ResultPageRequest{
				ChunkSize: 1000,
			},
			expectedResp: EventChunk{
				Events: []EmittedEvent{
					{
						BlockHash:       utils.TestHexToFelt(t, "0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb"),
						BlockNumber:     1472,
						TransactionHash: utils.TestHexToFelt(t, "0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d"),
					},
				},
			},
		}},
		"mainnet": {
			{
				eventFilter: EventFilter{
					FromBlock: BlockID{Number: &fromNum},
					ToBlock:   BlockID{Number: &toNum},
					Address:   utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					Keys: [][]*felt.Felt{{
						utils.TestHexToFelt(t, "0x3774b0545aabb37c45c1eddc6a7dae57de498aae6d5e3589e362d4b4323a533"),
						utils.TestHexToFelt(t, "0x714ae72367a39c17df987cf00f7cbb69c8cdcfa611e69e3511b5d16a23e2ec5"),
					}},
				},
				resPageReq: ResultPageRequest{
					ChunkSize: 1000,
				},
				expectedResp: EventChunk{
					Events: []EmittedEvent{
						{
							BlockHash:       utils.TestHexToFelt(t, "0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb"),
							BlockNumber:     1472,
							TransactionHash: utils.TestHexToFelt(t, "0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d"),
						},
					},
				},
			},
		},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		eventInput := EventsInput{
			EventFilter:       test.eventFilter,
			ResultPageRequest: test.resPageReq,
		}
		events, err := testConfig.provider.Events(context.Background(), eventInput)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, events.Events[0].BlockHash, test.expectedResp.Events[0].BlockHash, "BlockHash mismatch")
		require.Equal(t, events.Events[0].BlockNumber, test.expectedResp.Events[0].BlockNumber, "BlockNumber mismatch")
		require.Equal(t, events.Events[0].TransactionHash, test.expectedResp.Events[0].TransactionHash, "TransactionHash mismatch")
	}
}
