package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/require"
)

// TestEvents is a test function for testing the Events function.
//
// It creates a test configuration and defines a test set with different scenarios for the Events function.
// The test set includes a "mock" scenario and a "mainnet" scenario.
// In the "mock" scenario, it sets up the event filter, result page request, and expected response.
// In the "mainnet" scenario, it sets up the event filter, result page request, and expected response.
// It then iterates through the test set and performs the following steps for each test:
//   - Creates a spy object.
//   - Sets the provider's context to the spy object.
//   - Sets up the event input with the event filter and result page request.
//   - Calls the Events function with the event input.
//   - Checks if there is an error and fails the test if there is.
//   - Compares the events' block hash, block number, and transaction hash with the expected response.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestEvents(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv, tests.MainnetEnv, tests.IntegrationEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		eventFilter  EventFilter
		resPageReq   ResultPageRequest
		expectedResp EventChunk
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {{
			eventFilter: EventFilter{},
			resPageReq: ResultPageRequest{
				ChunkSize: 1000,
			},
			expectedResp: EventChunk{
				Events: []EmittedEvent{
					{
						BlockHash:       internalUtils.TestHexToFelt(t, "0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb"),
						BlockNumber:     1472,
						TransactionHash: internalUtils.TestHexToFelt(t, "0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d"),
					},
				},
			},
		}},
		tests.MainnetEnv: {
			{
				eventFilter: EventFilter{
					FromBlock: WithBlockNumber(1471),
					ToBlock:   WithBlockNumber(1473),
					Address:   internalUtils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					Keys: [][]*felt.Felt{{
						internalUtils.TestHexToFelt(t, "0x3774b0545aabb37c45c1eddc6a7dae57de498aae6d5e3589e362d4b4323a533"),
						internalUtils.TestHexToFelt(t, "0x714ae72367a39c17df987cf00f7cbb69c8cdcfa611e69e3511b5d16a23e2ec5"),
					}},
				},
				resPageReq: ResultPageRequest{
					ChunkSize: 1000,
				},
				expectedResp: EventChunk{
					Events: []EmittedEvent{
						{
							Event: Event{
								FromAddress: internalUtils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
								EventContent: EventContent{
									Keys: []*felt.Felt{
										internalUtils.TestHexToFelt(t, "0x3774b0545aabb37c45c1eddc6a7dae57de498aae6d5e3589e362d4b4323a533"),
									},
									Data: []*felt.Felt{
										internalUtils.TestHexToFelt(t, "0x0714ae72367a39c17df987cf00f7cbb69c8cdcfa611e69e3511b5d16a23e2ec5"),
										internalUtils.TestHexToFelt(t, "0x0714ae72367a39c17df987cf00f7cbb69c8cdcfa611e69e3511b5d16a23e2ec5"),
									},
								},
							},
							BlockHash:       internalUtils.TestHexToFelt(t, "0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb"),
							BlockNumber:     1472,
							TransactionHash: internalUtils.TestHexToFelt(t, "0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d"),
						},
					},
				},
			},
		},
		tests.TestnetEnv: {
			{
				eventFilter: EventFilter{
					FromBlock: WithBlockNumber(144932),
					ToBlock:   WithBlockNumber(144933),
					Address:   internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
					Keys: [][]*felt.Felt{{
						internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
						internalUtils.TestHexToFelt(t, "0x0143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6"),
					}},
				},
				resPageReq: ResultPageRequest{
					ChunkSize: 1000,
				},
				expectedResp: EventChunk{
					Events: []EmittedEvent{
						{
							Event: Event{
								FromAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
								EventContent: EventContent{
									Keys: []*felt.Felt{
										internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
									},
									Data: internalUtils.TestHexArrToFelt(t, []string{
										"0x0143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6",
										"0x01176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8",
										"0x04cffc7cb4ec5bd0",
										"0x00",
									}),
								},
							},
							BlockHash:       internalUtils.TestHexToFelt(t, "0x18a38934263d8b27a15d2e170d90030daa3a66e589b545908f376a8fdc971c8"),
							BlockNumber:     144933,
							TransactionHash: internalUtils.TestHexToFelt(t, "0x622817859a37dedf36cfb1417247f93dcc5840845bb8969df47491ef33e088e"),
						},
					},
				},
			},
		},
		tests.IntegrationEnv: {
			{
				eventFilter: EventFilter{
					FromBlock: WithBlockNumber(144932),
					ToBlock:   WithBlockNumber(144933),
					Address:   internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
					Keys: [][]*felt.Felt{{
						internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
						internalUtils.TestHexToFelt(t, "0x0143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6"),
					}},
				},
				resPageReq: ResultPageRequest{
					ChunkSize: 1000,
				},
				expectedResp: EventChunk{
					Events: []EmittedEvent{
						{
							Event: Event{
								FromAddress: internalUtils.TestHexToFelt(t, "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
								EventContent: EventContent{
									Keys: []*felt.Felt{
										internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
										internalUtils.TestHexToFelt(t, "0x352057331d5ad77465315d30b98135ddb815b86aa485d659dfeef59a904f88d"),
										internalUtils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
									},
									Data: internalUtils.TestHexArrToFelt(t, []string{
										"0x1e03a73820",
										"0x0",
									}),
								},
							},
							BlockHash:       internalUtils.TestHexToFelt(t, "0x2bc26a46f2bcf0f163514d0936040030fc52a3896133bfede267e49a1d5906f"),
							BlockNumber:     144932,
							TransactionHash: internalUtils.TestHexToFelt(t, "0x15be3c702ac6982b37e112962b1e94f896d8dcb90fc1655e1ea1571047e2a64"),
						},
					},
				},
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		eventInput := EventsInput{
			EventFilter:       test.eventFilter,
			ResultPageRequest: test.resPageReq,
		}
		events, err := testConfig.Provider.Events(context.Background(), eventInput)
		require.NoError(t, err, "Events failed")
		require.Exactly(t, test.expectedResp.Events[0], events.Events[0], "Events mismatch")
	}
}

func TestEventWith(t *testing.T) {
	key := "0xabc"
	feltKey := internalUtils.TestHexToFelt(t, key)

	events := []Event{
		{
			EventContent: EventContent{
				Keys: []*felt.Felt{feltKey},
			},
		},
	}

	found := EventWith(events, key)
	require.NotNil(t, found, "Expected to find event")
	require.True(t, found.Keys[0].Equal(feltKey), "Expected matching key")
}

func TestTransactionReceiptWithBlockInfo_EventWith(t *testing.T) {
	key := "0xdead"
	feltKey := internalUtils.TestHexToFelt(t, key)

	receipt := &TransactionReceiptWithBlockInfo{
		TransactionReceipt: TransactionReceipt{
			Events: []Event{
				{
					EventContent: EventContent{
						Keys: []*felt.Felt{feltKey},
					},
				},
			},
		},
	}

	found := receipt.EventWith(key)
	require.NotNil(t, found, "Expected to find event from receipt method")
	require.True(t, found.Keys[0].Equal(feltKey), "Expected matching key in receipt method")
}
