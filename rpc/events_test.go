package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestEvents is a test function for testing the Events function.
func TestEvents(t *testing.T) {
	tests.RunTestOn(t,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.MockEnv,
		tests.TestnetEnv,
	)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		description   string
		eventFilter   EventFilter
		resPageReq    ResultPageRequest
		expectedError error
	}

	// for this method, it seems the same data works for all the networks,
	// so we can use a single test set
	evFilter := EventFilter{
		FromBlock: WithBlockNumber(2000000),
		ToBlock:   WithBlockNumber(2000100),
		Address: internalUtils.TestHexToFelt(
			t,
			"0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d",
		),
		Keys: [][]*felt.Felt{{
			internalUtils.TestHexToFelt(
				t,
				"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9",
			),
			internalUtils.TestHexToFelt(
				t,
				"0x0143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6",
			),
		}},
	}

	tooManyKeys := make([][]*felt.Felt, 10000)
	for i := range 10000 {
		tooManyKeys[i] = []*felt.Felt{internalUtils.DeadBeef}
	}
	tooManyKeysFilter := evFilter
	tooManyKeysFilter.Keys = tooManyKeys

	invalidBlockFilter := evFilter
	invalidBlockFilter.FromBlock = WithBlockHash(internalUtils.DeadBeef)

	testSets := []testSetType{
		{
			description: "normal call",
			eventFilter: evFilter,
			resPageReq:  ResultPageRequest{ChunkSize: 10},
		},
		{
			description:   "invalid chunk size",
			eventFilter:   evFilter,
			resPageReq:    ResultPageRequest{ChunkSize: 10000000000},
			expectedError: ErrPageSizeTooBig,
		},
		{
			description: "invalid continuation token",
			eventFilter: evFilter,
			resPageReq: ResultPageRequest{
				ChunkSize:         10,
				ContinuationToken: "deadbeef",
			},
			expectedError: ErrInvalidContinuationToken,
		},
		{
			description:   "too many keys in filter",
			eventFilter:   tooManyKeysFilter,
			resPageReq:    ResultPageRequest{ChunkSize: 10},
			expectedError: ErrTooManyKeysInFilter,
		},
		{
			description:   "invalid block",
			eventFilter:   invalidBlockFilter,
			resPageReq:    ResultPageRequest{ChunkSize: 10},
			expectedError: ErrBlockNotFound,
		},
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv:        testSets,
		tests.IntegrationEnv: testSets,
		tests.MainnetEnv:     testSets,
		tests.TestnetEnv:     testSets,
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.description, func(t *testing.T) {
			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_getEvents",
						EventsInput{
							EventFilter:       test.eventFilter,
							ResultPageRequest: test.resPageReq,
						},
					).
					DoAndReturn(func(_, result, _ any, _ ...any) error {
						rawResp := result.(*json.RawMessage)

						if test.resPageReq.ChunkSize > 1000 {
							return RPCError{
								Code:    31,
								Message: "Requested page size is too big",
							}
						}

						if len(test.eventFilter.Keys) > 1000 {
							return RPCError{
								Code:    34,
								Message: "Too many keys provided in a filter",
							}
						}

						if test.resPageReq.ContinuationToken == "deadbeef" {
							return RPCError{
								Code:    33,
								Message: "The supplied continuation token is invalid or unknown",
							}
						}

						if test.eventFilter.FromBlock.Hash != nil &&
							test.eventFilter.FromBlock.Hash == internalUtils.DeadBeef {
							return RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						*rawResp = *internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/events/sepoliaEvents.json",
							"result",
						)

						return nil
					}).
					Times(1)
			}

			events, err := testConfig.Provider.Events(
				t.Context(),
				EventsInput{
					EventFilter:       test.eventFilter,
					ResultPageRequest: test.resPageReq,
				},
			)
			if test.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawEvents, err := json.Marshal(events)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawEvents))
		})
	}
}
