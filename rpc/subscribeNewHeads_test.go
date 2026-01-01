package rpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestSubscribeNewHeads tests the SubscribeNewHeads function.
func TestSubscribeNewHeads(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.TestnetEnv,
	)
	t.Parallel()

	type testSetType struct {
		description   string
		subBlockID    SubscriptionBlockID
		expectedError error
	}

	networkTestSet := []testSetType{
		{
			description: "normal call, zero subBlockID",
		},
		{
			description: "with tag latest",
			subBlockID:  new(SubscriptionBlockID).WithLatestTag(),
		},
		{
			description:   "error - too many blocks back",
			subBlockID:    new(SubscriptionBlockID).WithBlockNumber(3_000_000),
			expectedError: ErrTooManyBlocksBack,
		},
		{
			description:   "error - block not found",
			subBlockID:    new(SubscriptionBlockID).WithBlockHash(internalUtils.DeadBeef),
			expectedError: ErrBlockNotFound,
		},
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				description: "with tag latest",
				subBlockID:  new(SubscriptionBlockID).WithLatestTag(),
			},
			{
				description:   "error - too many blocks back",
				subBlockID:    new(SubscriptionBlockID).WithBlockNumber(3_000_000),
				expectedError: ErrTooManyBlocksBack,
			},
			{
				description:   "error - block not found",
				subBlockID:    new(SubscriptionBlockID).WithBlockHash(internalUtils.DeadBeef),
				expectedError: ErrBlockNotFound,
			},
		},
		tests.TestnetEnv:     networkTestSet,
		tests.IntegrationEnv: networkTestSet,
		tests.MainnetEnv:     networkTestSet,
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run("test: "+test.description, func(t *testing.T) {
			t.Parallel()
			tsetup := BeforeEach(t, true)

			if tests.TEST_ENV == tests.MockEnv {
				tsetup.MockClient.EXPECT().
					SubscribeWithSliceArgs(
						t.Context(),
						"starknet",
						"_subscribeNewHeads",
						gomock.Any(),
						test.subBlockID,
					).
					DoAndReturn(func(_, _, _, channel any, args ...any) (*client.ClientSubscription, error) {
						ch := channel.(chan json.RawMessage)
						subBlockID := args[0].(SubscriptionBlockID)

						if subBlockID.Number != nil && *subBlockID.Number == 3_000_000 {
							return nil, RPCError{
								Code:    68,
								Message: "Cannot go back more than 1024 blocks",
							}
						}

						if subBlockID.Hash != nil && subBlockID.Hash == internalUtils.DeadBeef {
							return nil, RPCError{
								Code:    24,
								Message: "Block not found",
							}
						}

						msg := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/ws/sepoliaNewHeads.json",
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

			headers := make(chan *BlockHeader)

			sub, err := tsetup.WsProvider.SubscribeNewHeads(
				t.Context(),
				headers,
				test.subBlockID,
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
				case resp := <-headers:
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
					t.Fatal("no new heads received")

					return
				}
			}
		})
	}
}
