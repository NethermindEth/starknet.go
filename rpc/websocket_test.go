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
	tests.RunTestOn(t, tests.TestnetEnv)

	t.Parallel()

	testConfig := beforeEach(t, true)

	type testSetType struct {
		headers         chan *BlockHeader
		subBlockID      BlockID
		counter         int
		isErrorExpected bool
		description     string
	}

	provider := testConfig.provider
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
				subBlockID:      WithBlockTag("latest"),
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

			wsProvider := testConfig.wsProvider

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
	tests.RunTestOn(t, tests.TestnetEnv)

	t.Parallel()

	testConfig := beforeEach(t, true)

	type testSetType struct {
		// Example values for the test
		fromAddressExample *felt.Felt
		keyExample         *felt.Felt
	}

	testSet := map[tests.TestEnv]testSetType{
		tests.TestnetEnv: {
			// sepolia StarkGate: ETH Token
			fromAddressExample: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			// "Transfer" event key, used by StarkGate ETH Token and STRK Token contracts
			keyExample: internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
		},
	}[tests.TEST_ENV]

	provider := testConfig.provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	t.Run("with empty args", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

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

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("blockID only", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID: WithBlockNumber(blockNumber - 1000),
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

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("fromAddress only, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

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

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				// Subscription with fromAddress should only return events from the specified address.
				// 'fromAddressExample' is the address of the sepolia StarkGate: ETH Token, which is very likely to have events,
				// so we can use it to verify the events are returned correctly.
				require.Equal(t, testSet.fromAddressExample, resp.FromAddress)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Skip("timeout waiting for events")
			}
		}
	})

	t.Run("keys only, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

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

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				// Subscription with keys should only return events with the specified keys.
				require.Equal(t, testSet.keyExample, resp.Keys[0])

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Skip("timeout waiting for events")
			}
		}
	})

	t.Run("with all arguments, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID:     WithBlockNumber(blockNumber - 100),
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
			case <-time.After(20 * time.Second):
				t.Skip("timeout waiting for events")
			}
		}
	})

	t.Run("error calls", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		type testSetType struct {
			input         EventSubscriptionInput
			expectedError error
		}

		keys := make([][]*felt.Felt, 1025)
		for i := range 1025 {
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
					BlockID: WithBlockNumber(blockNumber - 1025),
				},
				expectedError: ErrTooManyBlocksBack,
			},
			{
				input: EventSubscriptionInput{
					BlockID: WithBlockNumber(blockNumber + 100),
				},
				expectedError: ErrBlockNotFound,
			},
		}

		for _, test := range testSet {
			t.Run("test: "+test.expectedError.Error(), func(t *testing.T) {
				t.Logf("test: %+v", test.expectedError.Error())
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

func TestSubscribeTransactionStatus(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	t.Parallel()

	testConfig := beforeEach(t, true)

	provider := testConfig.provider
	var txHash *felt.Felt
	isBlockEmpty := true

	for isBlockEmpty {
		blockInterface, err := provider.BlockWithTxHashes(context.Background(), WithBlockTag(BlockTagLatest))
		require.NoError(t, err)
		block := blockInterface.(*BlockTxHashes)
		if len(block.Transactions) > 0 {
			isBlockEmpty = false
			txHash = block.Transactions[0]
		}
	}

	t.Run("normal call", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *NewTxnStatus)
		sub, err := wsProvider.SubscribeTransactionStatus(context.Background(), events, txHash)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &NewTxnStatus{}, resp)
				require.Equal(t, txHash, resp.TransactionHash)
				require.Equal(t, TxnStatus_Accepted_On_L2, resp.Status.FinalityStatus)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})
}

func TestUnsubscribe(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	t.Parallel()

	testConfig := beforeEach(t, true)

	wsProvider := testConfig.wsProvider

	events := make(chan *EmittedEvent)
	sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{})
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
