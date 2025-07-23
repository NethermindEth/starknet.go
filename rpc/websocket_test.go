package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestSubscribeNewHeads(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)

	type testSetType struct {
		headers         chan *BlockHeader
		subBlockID      BlockID
		counter         int
		isErrorExpected bool
		description     string
	}

	provider := testConfig.Provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1} // for the case the latest block number is updated

	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {
			{
				headers:         make(chan *BlockHeader),
				isErrorExpected: false,
				description:     "normal call, without subBlockID",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      WithBlockTag(BlockTagLatest),
				isErrorExpected: false,
				description:     "with tag latest",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      WithBlockNumber(blockNumber - 100),
				counter:         100,
				isErrorExpected: false,
				description:     "with block number within the range of 1024 blocks",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      WithBlockNumber(blockNumber - 1025),
				isErrorExpected: true,
				description:     "invalid, with block number out of the range of 1024 blocks",
			},
		},
		tests.IntegrationEnv: {
			{
				headers:         make(chan *BlockHeader),
				isErrorExpected: false,
				description:     "normal call, without subBlockID",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      WithBlockTag(BlockTagLatest),
				isErrorExpected: false,
				description:     "with tag latest",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      WithBlockNumber(blockNumber - 100),
				counter:         100,
				isErrorExpected: false,
				description:     "with block number within the range of 1024 blocks",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      WithBlockNumber(blockNumber - 1025),
				isErrorExpected: true,
				description:     "invalid, with block number out of the range of 1024 blocks",
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run("test: "+test.description, func(t *testing.T) {
			t.Parallel()

			wsProvider := testConfig.WsProvider

			var sub *client.ClientSubscription
			sub, err = wsProvider.SubscribeNewHeads(context.Background(), test.headers, test.subBlockID)
			if sub != nil {
				defer sub.Unsubscribe()
			}

			if test.isErrorExpected {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			for {
				select {
				case resp := <-test.headers:
					require.IsType(t, &BlockHeader{}, resp)

					if test.counter != 0 {
						if test.counter == 1 {
							require.Contains(t, latestBlockNumbers, resp.Number+1)

							return
						} else {
							test.counter--
						}
					} else {
						return
					}
				case err := <-sub.Err():
					require.NoError(t, err)
				}
			}
		})
	}
}

//nolint:gocyclo
func TestSubscribeEvents(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)

	type testSetType struct {
		// Example values for the test
		fromAddressExample *felt.Felt
		keyExample         *felt.Felt
	}

	testSet := map[tests.TestEnv]testSetType{
		tests.TestnetEnv: {
			// sepolia StarkGate: STRK Token
			fromAddressExample: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
			// "Transfer" event key, used by StarkGate ETH Token and STRK Token contracts
			keyExample: internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
		},
		tests.IntegrationEnv: {
			// a contract with a lot of txns in integration network
			fromAddressExample: internalUtils.TestHexToFelt(t, "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
			// "Transfer" event key
			keyExample: internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
		},
	}[tests.TEST_ENV]

	provider := testConfig.Provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	// 'blockNumber + 1' for the case the latest block number is updated
	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1}

	t.Run("with empty args", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("blockID only - 1000 blocks back", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID: WithBlockNumber(blockNumber - 1000),
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		uniqueAddresses := make(map[string]bool)
		uniqueKeys := make(map[string]bool)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)
				// Subscription with only blockID should return events from all addresses and keys from the specified block onwards.
				// As none filters are applied, the events should be from all addresses and keys.

				uniqueAddresses[resp.FromAddress.String()] = true
				uniqueKeys[resp.Keys[0].String()] = true

				if tests.TEST_ENV == tests.IntegrationEnv {
					// in integration network, there less unique addresses and keys
					if len(uniqueAddresses) >= 2 && len(uniqueKeys) >= 2 {
						return
					}
				} else {
					// check if there are at least 3 different addresses and keys in the received events
					if len(uniqueAddresses) >= 3 && len(uniqueKeys) >= 3 {
						return
					}
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("blockID only - with tag latest", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider
		provider := testConfig.Provider
		rawBlock, err := provider.BlockWithTxHashes(context.Background(), WithBlockTag(BlockTagLatest))
		require.NoError(t, err)
		expectedBlock, ok := rawBlock.(*BlockTxHashes)
		require.True(t, ok)

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID: WithBlockTag(BlockTagLatest),
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				if len(expectedBlock.Transactions) > 0 {
					// since we are subscribing to the latest block, the event block number should be the same as the latest block number,
					// but if the latest block is empty, the subscription will return events from later blocks.
					// Also, we can have race condition here, in case the latest block is updated between the `BlockWithTxHashes`
					// request and the subscription.
					require.Equal(t, expectedBlock.Number, resp.BlockNumber)
				}

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("fromAddress only, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			FromAddress: testSet.fromAddressExample,
			BlockID:     WithBlockNumber(blockNumber - 1000),
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		uniqueKeys := make(map[string]bool)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				// Subscription with fromAddress should only return events from the specified address.
				// 'fromAddressExample' is the address of the sepolia StarkGate: ETH Token, which is very likely to have events,
				// so we can use it to verify the events are returned correctly.
				require.Equal(t, testSet.fromAddressExample, resp.FromAddress)

				if tests.TEST_ENV == tests.IntegrationEnv {
					// integration network is not very used by external users, so let's skip the keys verification
					return
				}

				uniqueKeys[resp.Keys[0].String()] = true

				// check if there are at least 2 different keys in the received events
				if len(uniqueKeys) >= 2 {
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("keys only, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			Keys:    [][]*felt.Felt{{testSet.keyExample}},
			BlockID: WithBlockNumber(blockNumber - 1000),
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		uniqueAddresses := make(map[string]bool)
		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				// Subscription with keys should only return events with the specified keys.
				require.Equal(t, testSet.keyExample, resp.Keys[0])

				uniqueAddresses[resp.FromAddress.String()] = true

				// check if there are at least 2 different addresses in the received events
				if len(uniqueAddresses) >= 2 {
					return
				}

				if tests.TEST_ENV == tests.IntegrationEnv {
					// integration network is not very used by external users, so let's skip the keys verification
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("with all arguments, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID:     WithBlockNumber(blockNumber - 1000),
			FromAddress: testSet.fromAddressExample,
			Keys:        [][]*felt.Felt{{testSet.keyExample}},
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)
				// 'fromAddressExample' is the address of the sepolia StarkGate: ETH Token, which is very likely to have events,
				// so we can use it to verify the events are returned correctly.
				require.Equal(t, testSet.fromAddressExample, resp.FromAddress)
				require.Equal(t, testSet.keyExample, resp.Keys[0])

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("error calls", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		type testSetType struct {
			input         EventSubscriptionInput
			expectedError error
		}

		keys := make([][]*felt.Felt, 2000)
		for i := range 2000 {
			keys[i] = []*felt.Felt{internalUtils.TestHexToFelt(t, "0x1")}
		}

		testSet := []testSetType{
			{
				input: EventSubscriptionInput{
					Keys: keys,
				},
				expectedError: ErrTooManyKeysInFilter,
			},
			{
				input: EventSubscriptionInput{
					BlockID: WithBlockNumber(blockNumber - 2000),
				},
				expectedError: ErrTooManyBlocksBack,
			},
			{
				input: EventSubscriptionInput{
					BlockID: WithBlockNumber(blockNumber + 10000),
				},
				expectedError: ErrBlockNotFound,
			},
		}

		for _, test := range testSet {
			t.Run(test.expectedError.Error(), func(t *testing.T) {
				t.Parallel()

				events := make(chan *EmittedEvent)
				defer close(events)
				sub, err := wsProvider.SubscribeEvents(context.Background(), events, &test.input)
				if sub != nil {
					defer sub.Unsubscribe()
				}
				require.Nil(t, sub)
				require.EqualError(t, err, test.expectedError.Error())
			})
		}
	})
}

func TestSubscribePendingTransactions(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)

	type testSetType struct {
		pendingTxns   chan *PendingTxn
		options       *SubPendingTxnsInput
		expectedError error
		description   string
	}

	addresses := make([]*felt.Felt, 1025)
	for i := range 1025 {
		addresses[i] = internalUtils.TestHexToFelt(t, "0x1")
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {
			{
				pendingTxns: make(chan *PendingTxn),
				options:     nil,
				description: "nil input",
			},
			{
				pendingTxns: make(chan *PendingTxn),
				options:     &SubPendingTxnsInput{},
				description: "empty input",
			},
			{
				pendingTxns: make(chan *PendingTxn),
				options:     &SubPendingTxnsInput{TransactionDetails: true},
				description: "with transanctionDetails true",
			},
			{
				pendingTxns:   make(chan *PendingTxn),
				options:       &SubPendingTxnsInput{SenderAddress: addresses},
				expectedError: ErrTooManyAddressesInFilter,
				description:   "error: too many addresses",
			},
		},
		tests.IntegrationEnv: {
			{
				pendingTxns: make(chan *PendingTxn),
				options:     nil,
				description: "nil input",
			},
			{
				pendingTxns: make(chan *PendingTxn),
				options:     &SubPendingTxnsInput{},
				description: "empty input",
			},
			{
				pendingTxns: make(chan *PendingTxn),
				options:     &SubPendingTxnsInput{TransactionDetails: true},
				description: "with transanctionDetails true",
			},
			{
				pendingTxns:   make(chan *PendingTxn),
				options:       &SubPendingTxnsInput{SenderAddress: addresses},
				expectedError: ErrTooManyAddressesInFilter,
				description:   "error: too many addresses",
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run("test: "+test.description, func(t *testing.T) {
			t.Parallel()

			wsProvider := testConfig.WsProvider

			sub, err := wsProvider.SubscribePendingTransactions(context.Background(), test.pendingTxns, test.options)
			if sub != nil {
				defer sub.Unsubscribe()
			}

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())

				return
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			for {
				select {
				case resp := <-test.pendingTxns:
					require.IsType(t, &PendingTxn{}, resp)

					if test.options == nil || !test.options.TransactionDetails {
						require.NotEmpty(t, resp.Hash)
						require.Empty(t, resp.Transaction)
					} else {
						require.NotEmpty(t, resp.Hash)
						require.NotEmpty(t, resp.Transaction)
					}

					return
				case err := <-sub.Err():
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)

	wsProvider := testConfig.WsProvider

	events := make(chan *EmittedEvent)
	sub, err := wsProvider.SubscribeEvents(context.Background(), events, nil)
	if sub != nil {
		defer sub.Unsubscribe()
	}
	require.NoError(t, err)
	require.NotNil(t, sub)

	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		sub.Unsubscribe()
	}()

loop:
	for {
		select {
		case resp := <-events:
			require.IsType(t, &EmittedEvent{}, resp)
		case err := <-sub.Err():
			// when unsubscribing, the error channel should return nil
			require.Nil(t, err)

			break loop
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for unsubscription")
		}
	}
}

// A simple test was made to make sure the reorg events are received. Ref:
// https://github.com/NethermindEth/starknet.go/pull/651#discussion_r1927356194
func TestReorgEvents(t *testing.T) {
	t.Skip("TODO: implement reorg test")
}
