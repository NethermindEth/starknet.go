package rpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TODO: we need to add logic to check the raw request sent to the RPC server.

const testDuration = 10 * time.Second

// an address currently sending a lot of transactions in Sepolia
var randAddress, _ = internalUtils.HexToFelt(
	"0x04f4e29add19afa12c868ba1f4439099f225403ff9a71fe667eebb50e13518d3",
)

// TestSubscribeEvents tests the SubscribeEvents function.
func TestSubscribeEvents(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.TestnetEnv,
	)

	t.Parallel()
	testConfig := BeforeEach(t, true)

	// STRK Token
	fromAddress := internalUtils.TestHexToFelt(
		t,
		"0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d",
	)
	// "Transfer" event key
	key := internalUtils.TestHexToFelt(
		t,
		"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9",
	)

	tooManyKeys := make([][]*felt.Felt, 10000)
	for i := range 10000 {
		tooManyKeys[i] = []*felt.Felt{new(felt.Felt).SetUint64(uint64(i))}
	}

	if tests.TEST_ENV == tests.MockEnv {
		testConfig.MockClient.EXPECT().
			CallContextWithSliceArgs(
				t.Context(),
				gomock.Any(),
				"starknet_blockNumber",
			).
			DoAndReturn(
				func(_, result, _ any, _ ...any) error {
					rawResp := result.(*json.RawMessage)
					rawBlockNumber := json.RawMessage("1234567890")
					*rawResp = rawBlockNumber

					return nil
				},
			)
	}
	blockNumber, err := testConfig.Provider.BlockNumber(t.Context())
	require.NoError(t, err)

	type testSetType struct {
		description   string
		input         *EventSubscriptionInput
		expectedError error
	}

	template := []testSetType{
		{
			description: "from address only",
			input: &EventSubscriptionInput{
				FromAddress: fromAddress,
			},
		},
		{
			description: "keys only",
			input: &EventSubscriptionInput{
				Keys: [][]*felt.Felt{{key}},
			},
		},
		{
			description: "with block ID only",
			input: &EventSubscriptionInput{
				SubBlockID: SubscriptionBlockID{
					Tag: BlockTagLatest,
				},
			},
		},
		{
			description: "with finality status PRE_CONFIRMED",
			input: &EventSubscriptionInput{
				FinalityStatus: TxnFinalityStatusPreConfirmed,
			},
		},
		{
			description: "with finality status ACCEPTED_ON_L2",
			input: &EventSubscriptionInput{
				FinalityStatus: TxnFinalityStatusAcceptedOnL2,
			},
		},
		{
			description: "all filters",
			input: &EventSubscriptionInput{
				FromAddress:    fromAddress,
				Keys:           [][]*felt.Felt{{key}},
				SubBlockID:     new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
				FinalityStatus: TxnFinalityStatusAcceptedOnL2,
			},
		},
		{
			description: "error: too many keys",
			input: &EventSubscriptionInput{
				Keys: tooManyKeys,
			},
			expectedError: ErrTooManyKeysInFilter,
		},
		{
			description: "error: too many blocks back",
			input: &EventSubscriptionInput{
				SubBlockID: new(SubscriptionBlockID).WithBlockNumber(3_000_000),
			},
			expectedError: ErrTooManyBlocksBack,
		},
		{
			description: "error: block not found",
			input: &EventSubscriptionInput{
				SubBlockID: new(SubscriptionBlockID).WithBlockHash(internalUtils.DeadBeef),
			},
			expectedError: ErrBlockNotFound,
		},
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv:        template,
		tests.IntegrationEnv: template,
		tests.MainnetEnv:     template,
		tests.TestnetEnv:     template,
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			tsetup := BeforeEach(t, true)

			if tests.TEST_ENV == tests.MockEnv {
				tsetup.MockClient.EXPECT().
					Subscribe(
						t.Context(),
						"starknet",
						"_subscribeEvents",
						gomock.Any(),
						test.input,
					).
					DoAndReturn(func(_, _, _, channel any, arg any) (*client.ClientSubscription, error) {
						ch := channel.(chan json.RawMessage)
						input := arg.(*EventSubscriptionInput)

						if input.SubBlockID.Number != nil && *input.SubBlockID.Number == 3_000_000 {
							return nil, RPCError{
								Code:    68,
								Message: "Cannot go back more than 1024 blocks",
							}
						}

						if input.SubBlockID.Hash != nil &&
							input.SubBlockID.Hash == internalUtils.DeadBeef {
							return nil, RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						if len(input.Keys) > 1000 {
							return nil, RPCError{
								Code:    34,
								Message: "Too many keys provided in a filter",
							}
						}

						msg := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/ws/sepoliaEvents.json",
							"params", "result",
						)

						go func() {
							for {
								select {
								case <-time.Tick(2 * time.Second):
									ch <- msg
								case <-t.Context().Done():
									return
								}
							}
						}()

						return &client.ClientSubscription{}, nil
					})
			}

			events := make(chan *EmittedEventWithFinalityStatus)
			sub, err := tsetup.WsProvider.SubscribeEvents(
				t.Context(),
				events,
				test.input,
			)
			if test.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError.Error())

				return
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			if tests.TEST_ENV != tests.MockEnv {
				// this would block mock tests since the ClientSubscription is empty
				defer sub.Unsubscribe()
			}

			stopTest := time.After(testDuration)
			for {
				select {
				case resp := <-events:
					require.NotNil(t, resp)

					rawExpectedMsg := <-tsetup.WSSpy.SpyChannel()
					rawMsg, err := json.Marshal(resp)
					require.NoError(t, err)
					assert.JSONEq(t, string(rawExpectedMsg), string(rawMsg))

					// stop test after a few seconds
					select {
					case <-stopTest:
						return
					default:
					}
				case err := <-sub.Err():
					require.NoError(t, err)
				case <-time.After(testDuration * 2):
					// Since we are setting some filters, it could be the case that no events match the filters
					// at the time. So we skip the test instead of failing it.
					t.Skip("no events received")

					return
				}
			}
		})
	}

	t.Run("with default options - nil input", func(t *testing.T) {
		t.Parallel()
		tsetup := BeforeEach(t, true)

		if tests.TEST_ENV == tests.MockEnv {
			tsetup.MockClient.EXPECT().
				Subscribe(
					t.Context(),
					"starknet",
					"_subscribeEvents",
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _, _, channel any, _ any) (*client.ClientSubscription, error) {
					ch := channel.(chan json.RawMessage)

					msg := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
						t,
						"./testData/ws/sepoliaEvents.json",
						"params", "result",
					)

					go func() {
						for {
							select {
							case <-time.Tick(2 * time.Second):
								ch <- msg
							case <-t.Context().Done():
								return
							}
						}
					}()

					return &client.ClientSubscription{}, nil
				})
		}

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := tsetup.WsProvider.SubscribeEvents(
			t.Context(),
			events,
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, sub)

		if tests.TEST_ENV != tests.MockEnv {
			// this would block mock tests since the ClientSubscription is empty
			defer sub.Unsubscribe()
		}

		stopTest := time.After(testDuration)
		for {
			select {
			case resp := <-events:
				require.NotNil(t, resp)

				rawExpectedMsg := <-tsetup.WSSpy.SpyChannel()
				rawMsg, err := json.Marshal(resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedMsg), string(rawMsg))

				// stop test after a few seconds
				select {
				case <-stopTest:
					return
				default:
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(testDuration * 2):
				t.Fatal("no events received")

				return
			}
		}
	})
}

// TestUnsubscribe tests the Unsubscribe method.
func TestUnsubscribe(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)
	t.Parallel()

	testConfig := BeforeEach(t, true)
	wsProvider := testConfig.WsProvider

	events := make(chan *EmittedEventWithFinalityStatus)
	sub, err := wsProvider.SubscribeEvents(t.Context(), events, nil)
	require.NoError(t, err)
	require.NotNil(t, sub)

	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		sub.Unsubscribe()
	}()

	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-events:
		case err := <-sub.Err():
			// when unsubscribing, the error channel should return nil
			assert.Nil(t, err)

			return
		case <-timeout:
			t.Fatal("timeout waiting for unsubscription")
		}
	}
}

// A simple test was made to make sure the reorg events are received. Ref:
// https://github.com/NethermindEth/starknet.go/pull/651#discussion_r1927356194
// Also here: https://github.com/NethermindEth/starknet.go/pull/781
func TestReorgEvents(t *testing.T) {
	t.Skip("TODO: implement reorg test")
}
