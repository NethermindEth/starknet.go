package rpc

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

func TestSubscribeNewHeads(t *testing.T) {
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}

	testConfig := beforeEach(t)
	require.NotNil(t, testConfig.wsBase, "wsProvider base is not set")

	type testSetType struct {
		headers         chan *BlockHeader
		blockID         []BlockID
		counter         int
		isErrorExpected bool
	}

	provider := testConfig.provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1} // for the case the latest block number is updated

	testSet := map[string][]testSetType{
		"testnet": {
			{ // normal
				headers:         make(chan *BlockHeader),
				isErrorExpected: false,
			},
			{ // with tag latest
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockTag("latest")},
				isErrorExpected: false,
			},
			{ // with tag pending
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockTag("pending")},
				isErrorExpected: true,
			},
			{ // with block number within the range of 1024 blocks
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockNumber(blockNumber - 100)},
				counter:         100,
				isErrorExpected: false,
			},
			{ // invalid, with block number out of the range of 1024 blocks
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockNumber(blockNumber - 1025)},
				isErrorExpected: true,
			},
			{ // invalid, more than one blockID parameter
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockTag("latest"), WithBlockTag("latest")},
				isErrorExpected: true,
			},
		},
	}[testEnv]

	for index, test := range testSet {
		t.Run(fmt.Sprintf("test %d", index+1), func(t *testing.T) {

			wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
			require.NoError(t, err)
			defer wsProvider.Close()

			var sub *client.ClientSubscription
			if len(test.blockID) == 0 {
				sub, err = wsProvider.SubscribeNewHeads(context.Background(), test.headers)
			} else {
				sub, err = wsProvider.SubscribeNewHeads(context.Background(), test.headers, test.blockID...)
			}

			if test.isErrorExpected {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.NotNil(t, sub)
			defer sub.Unsubscribe()

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
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}

	testConfig := beforeEach(t)
	require.NotNil(t, testConfig.wsBase, "wsProvider base is not set")

	provider := testConfig.provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1}                                              // for the case the latest block number is updated
	fromAddress := utils.HexToFeltNoErr("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7") // sepolia StarkGate: ETH Token
	key := utils.HexToFeltNoErr("0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9")

	t.Run("normal call, with empty args", func(t *testing.T) {
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, EventSubscriptionInput{})
		require.NoError(t, err)
		require.NotNil(t, sub)
		defer sub.Unsubscribe()

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
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, EventSubscriptionInput{
			FromAddress: fromAddress,
		})
		require.NoError(t, err)
		require.NotNil(t, sub)
		defer sub.Unsubscribe()

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)

				// Subscription with only fromAddress should return events from the specified address from the latest block onwards.
				require.Equal(t, fromAddress, resp.FromAddress)
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("normal call, keys only", func(t *testing.T) {
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, EventSubscriptionInput{
			Keys: [][]*felt.Felt{{key}},
		})
		require.NoError(t, err)
		require.NotNil(t, sub)
		defer sub.Unsubscribe()

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)

				// Subscription with only keys should return events with the specified keys from the latest block onwards.
				require.Equal(t, key, resp.Keys[0])
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("normal call, blockID only", func(t *testing.T) {
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, EventSubscriptionInput{
			BlockID: WithBlockNumber(blockNumber - 100),
		})
		require.NoError(t, err)
		require.NotNil(t, sub)
		defer sub.Unsubscribe()

		differentFromAddressFound := false
		differentKeyFound := false

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				// Subscription with only blockID should return events from all addresses and keys from the specified block onwards.
				// Verify by checking for events with different addresses and keys than the test values.
				if !differentFromAddressFound {
					if resp.FromAddress != fromAddress {
						differentFromAddressFound = true
					}
				}

				if !differentKeyFound {
					if !slices.Contains(resp.Keys, key) {
						differentKeyFound = true
					}
				}

				if differentFromAddressFound && differentKeyFound {
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
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		events := make(chan *EmittedEvent)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, EventSubscriptionInput{
			BlockID:     WithBlockNumber(blockNumber - 100),
			FromAddress: fromAddress,
			Keys:        [][]*felt.Felt{{key}},
		})
		require.NoError(t, err)
		require.NotNil(t, sub)
		defer sub.Unsubscribe()

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEvent{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)
				require.Equal(t, fromAddress, resp.FromAddress)
				require.Equal(t, key, resp.Keys[0])
				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(4 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("error calls", func(t *testing.T) {
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		type testSetType struct {
			input         EventSubscriptionInput
			expectedError error
		}

		keys := make([][]*felt.Felt, 1025)
		for i := 0; i < 1025; i++ {
			keys[i] = []*felt.Felt{utils.HexToFeltNoErr("0x1")}
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
					BlockID: WithBlockNumber(blockNumber + 2),
				},
				expectedError: ErrBlockNotFound,
			},
			{
				input: EventSubscriptionInput{
					BlockID: WithBlockTag("pending"),
				},
				expectedError: ErrCallOnPending,
			},
		}

		for _, test := range testSet {
			events := make(chan *EmittedEvent)
			defer close(events)
			sub, err := wsProvider.SubscribeEvents(context.Background(), events, test.input)
			require.Nil(t, sub)
			require.EqualError(t, err, test.expectedError.Error())
		}
	})
}

func TestSubscribeTransactionStatus(t *testing.T) {
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}

	testConfig := beforeEach(t)
	require.NotNil(t, testConfig.wsBase, "wsProvider base is not set")

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
		wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
		require.NoError(t, err)
		defer wsProvider.Close()

		events := make(chan *NewTxnStatusResp)
		sub, err := wsProvider.SubscribeTransactionStatus(context.Background(), events, txHash)
		require.NoError(t, err)
		require.NotNil(t, sub)
		defer sub.Unsubscribe()

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
