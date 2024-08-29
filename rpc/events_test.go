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
					FromBlock: WithBlockNumber(1471),
					ToBlock:   WithBlockNumber(1473),
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
							Event: Event{
								FromAddress: utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
								Keys: []*felt.Felt{
									utils.TestHexToFelt(t, "0x3774b0545aabb37c45c1eddc6a7dae57de498aae6d5e3589e362d4b4323a533"),
								},
								Data: []*felt.Felt{
									utils.TestHexToFelt(t, "0x0714ae72367a39c17df987cf00f7cbb69c8cdcfa611e69e3511b5d16a23e2ec5"),
									utils.TestHexToFelt(t, "0x0714ae72367a39c17df987cf00f7cbb69c8cdcfa611e69e3511b5d16a23e2ec5"),
								},
							},
							BlockHash:       utils.TestHexToFelt(t, "0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb"),
							BlockNumber:     1472,
							TransactionHash: utils.TestHexToFelt(t, "0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d"),
						},
					},
				},
			},
		},
		"testnet": {
			{
				eventFilter: EventFilter{
					FromBlock: WithBlockNumber(144932),
					ToBlock:   WithBlockNumber(144933),
					Address:   utils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
					Keys: [][]*felt.Felt{{
						utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
						utils.TestHexToFelt(t, "0x0143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6"),
					}},
				},
				resPageReq: ResultPageRequest{
					ChunkSize: 1000,
				},
				expectedResp: EventChunk{
					Events: []EmittedEvent{
						{
							Event: Event{
								FromAddress: utils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
								Keys: []*felt.Felt{
									utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
								},
								Data: utils.TestHexArrToFelt(t, []string{
									"0x0143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6",
									"0x01176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8",
									"0x04cffc7cb4ec5bd0",
									"0x00",
								}),
							},
							BlockHash:       utils.TestHexToFelt(t, "0x18a38934263d8b27a15d2e170d90030daa3a66e589b545908f376a8fdc971c8"),
							BlockNumber:     144933,
							TransactionHash: utils.TestHexToFelt(t, "0x622817859a37dedf36cfb1417247f93dcc5840845bb8969df47491ef33e088e"),
						},
					},
				},
			},
		},
	}[testEnv]

	for _, test := range testSet {
		eventInput := EventsInput{
			EventFilter:       test.eventFilter,
			ResultPageRequest: test.resPageReq,
		}
		events, err := testConfig.provider.Events(context.Background(), eventInput)
		require.NoError(t, err, "Events failed")
		require.Exactly(t, test.expectedResp.Events[0], events.Events[0], "Events mismatch")
	}
}
