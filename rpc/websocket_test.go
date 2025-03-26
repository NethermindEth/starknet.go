package rpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestSubscribeNewHeads(t *testing.T) {
	t.Parallel()

	testConfig := beforeEach(t, true)

	type testSetType struct {
		headers         chan *BlockHeader
		subBlockID      *SubscriptionBlockID
		counter         int
		isErrorExpected bool
		description     string
	}

	provider := testConfig.provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1} // for the case the latest block number is updated

	testSet, ok := map[string][]testSetType{
		"testnet": {
			{
				headers:         make(chan *BlockHeader),
				isErrorExpected: false,
				description:     "normal call, without subBlockID",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      &SubscriptionBlockID{Tag: "latest"},
				isErrorExpected: false,
				description:     "with tag latest",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      &SubscriptionBlockID{Number: blockNumber - 100},
				counter:         100,
				isErrorExpected: false,
				description:     "with block number within the range of 1024 blocks",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      &SubscriptionBlockID{Number: blockNumber - 1025},
				isErrorExpected: true,
				description:     "invalid, with block number out of the range of 1024 blocks",
			},
		},
	}[testEnv]

	if !ok {
		t.Skip("test environment not supported")
	}

	for _, test := range testSet {
		t.Run(fmt.Sprintf("test: %s", test.description), func(t *testing.T) {
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
							require.Contains(t, latestBlockNumbers, resp.BlockNumber+1)
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

func TestSubscribeEvents(t *testing.T) {
	t.Parallel()

	testConfig := beforeEach(t, true)

	type testSetType struct {
		// Example values for the test
		fromAddressExample *felt.Felt
		keyExample         *felt.Felt
	}

	testSet, ok := map[string]testSetType{
		"testnet": {
			// sepolia StarkGate: ETH Token
			fromAddressExample: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			// random key from StarkGate: ETH Token
			keyExample: internalUtils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
		},
	}[testEnv]

	if !ok {
		t.Skip("test environment not supported")
	}

	provider := testConfig.provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	// 'blockNumber + 1' for the case the latest block number is updated
	// '0' for the case of events from pending blocks
	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1, 0}

	t.Run("normal call, with empty args", func(t *testing.T) {
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
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("normal call, fromAddress only", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			FromAddress: testSet.fromAddressExample,
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
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)

				// Subscription with only fromAddress should return events from the specified address from the latest block onwards.
				require.Equal(t, testSet.fromAddressExample, resp.FromAddress)
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("normal call, keys only", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			Keys: [][]*felt.Felt{{testSet.keyExample}},
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
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)

				// Subscription with only keys should return events with the specified keys from the latest block onwards.
				require.Equal(t, testSet.keyExample, resp.Keys[0])
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("normal call, blockID only", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID: SubscriptionBlockID{Number: blockNumber - 100},
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

				// check if there are at least 3 different addresses and keys in the received events
				if len(uniqueAddresses) >= 3 && len(uniqueKeys) >= 3 {
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("normal call, with all arguments, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.wsProvider

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			BlockID:     SubscriptionBlockID{Number: blockNumber - 100},
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
				// 'fromAddress' is the address of the sepolia StarkGate: ETH Token, which is very likely to have events,
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
					BlockID: SubscriptionBlockID{Number: blockNumber - 1025},
				},
				expectedError: ErrTooManyBlocksBack,
			},
			{
				input: EventSubscriptionInput{
					BlockID: SubscriptionBlockID{Number: blockNumber + 100},
				},
				expectedError: ErrBlockNotFound,
			},
		}

		for _, test := range testSet {
			t.Logf("test: %+v", test.expectedError.Error())
			events := make(chan *EmittedEvent)
			defer close(events)
			sub, err := wsProvider.SubscribeEvents(context.Background(), events, &test.input)
			if sub != nil {
				defer sub.Unsubscribe()
			}
			require.Nil(t, sub)
			require.EqualError(t, err, test.expectedError.Error())
		}
	})
}

func TestSubscribeTransactionStatus(t *testing.T) {
	t.Parallel()

	testConfig := beforeEach(t, true)

	testSet := map[string]bool{
		"testnet": true,
	}[testEnv]

	if !testSet {
		t.Skip("test environment not supported")
	}

	provider := testConfig.provider
	blockInterface, err := provider.BlockWithTxHashes(context.Background(), WithBlockTag("latest"))
	require.NoError(t, err)
	block := blockInterface.(*BlockTxHashes)

	txHash := new(felt.Felt)
	for _, tx := range block.Transactions {
		status, err := provider.GetTransactionStatus(context.Background(), tx)
		require.NoError(t, err)
		if status.FinalityStatus == TxnStatus_Accepted_On_L2 {
			txHash = tx
			break
		}
	}

	t.Run("normal call", func(t *testing.T) {
		wsProvider := testConfig.wsProvider

		events := make(chan *NewTxnStatusResp)
		sub, err := wsProvider.SubscribeTransactionStatus(context.Background(), events, txHash)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &NewTxnStatusResp{}, resp)
				require.Equal(t, txHash, resp.TransactionHash)
				require.Equal(t, TxnStatus_Accepted_On_L2, resp.Status.FinalityStatus)
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})
}

func TestSubscribePendingTransactions(t *testing.T) {
	t.Parallel()

	testConfig := beforeEach(t, true)

	type testSetType struct {
		pendingTxns   chan *SubPendingTxns
		options       *SubPendingTxnsInput
		expectedError error
		description   string
	}

	addresses := make([]*felt.Felt, 1025)
	for i := range 1025 {
		addresses[i] = internalUtils.TestHexToFelt(t, "0x1")
	}

	testSet, ok := map[string][]testSetType{
		"testnet": {
			{
				pendingTxns: make(chan *SubPendingTxns),
				options:     nil,
				description: "nil input",
			},
			{
				pendingTxns: make(chan *SubPendingTxns),
				options:     &SubPendingTxnsInput{},
				description: "empty input",
			},
			{
				pendingTxns: make(chan *SubPendingTxns),
				options:     &SubPendingTxnsInput{TransactionDetails: true},
				description: "with transanctionDetails true",
			},
			{
				pendingTxns:   make(chan *SubPendingTxns),
				options:       &SubPendingTxnsInput{SenderAddress: addresses},
				expectedError: ErrTooManyAddressesInFilter,
				description:   "error: too many addresses",
			},
		},
	}[testEnv]

	if !ok {
		t.Skip("test environment not supported")
	}

	for _, test := range testSet {
		t.Run(fmt.Sprintf("test: %s", test.description), func(t *testing.T) {
			t.Parallel()

			wsProvider := testConfig.wsProvider

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
					require.IsType(t, &SubPendingTxns{}, resp)

					if test.options == nil || !test.options.TransactionDetails {
						require.NotEmpty(t, resp.TransactionHash)
						require.Empty(t, resp.Transaction)
					} else {
						require.NotEmpty(t, resp.TransactionHash)
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
	t.Parallel()

	testConfig := beforeEach(t, true)

	testSet := map[string]bool{
		"testnet": true,
	}[testEnv]

	if !testSet {
		t.Skip("test environment not supported")
	}

	wsProvider := testConfig.wsProvider

	events := make(chan *EmittedEvent)
	sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{})
	if sub != nil {
		defer sub.Unsubscribe()
	}
	require.NoError(t, err)
	require.NotNil(t, sub)

	go func(t *testing.T) {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		sub.Unsubscribe()
	}(t)

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
