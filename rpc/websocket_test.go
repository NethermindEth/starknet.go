package rpc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribeNewHeads(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)

	type testSetType struct {
		headers         chan *BlockHeader
		subBlockID      SubscriptionBlockID
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
				subBlockID:      new(SubscriptionBlockID).WithLatestTag(),
				isErrorExpected: false,
				description:     "with tag latest",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      new(SubscriptionBlockID).WithBlockNumber(blockNumber - 100),
				counter:         100,
				isErrorExpected: false,
				description:     "with block number within the range of 1024 blocks",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1025),
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
				subBlockID:      new(SubscriptionBlockID).WithLatestTag(),
				isErrorExpected: false,
				description:     "with tag latest",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      new(SubscriptionBlockID).WithBlockNumber(blockNumber - 100),
				counter:         100,
				isErrorExpected: false,
				description:     "with block number within the range of 1024 blocks",
			},
			{
				headers:         make(chan *BlockHeader),
				subBlockID:      new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1025),
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

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				require.Contains(t, latestBlockNumbers, resp.BlockNumber)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("blockID only - 1000 blocks back", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			SubBlockID: new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
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
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
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
			case <-time.After(10 * time.Second):
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

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			SubBlockID: new(SubscriptionBlockID).WithLatestTag(),
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
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
			case <-time.After(10 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("finalityStatus only", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		t.Run("with finality status ACCEPTED_ON_L2", func(t *testing.T) {
			t.Parallel()

			events := make(chan *EmittedEventWithFinalityStatus)
			sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
				FinalityStatus: TxnFinalityStatusAcceptedOnL2,
			})
			if sub != nil {
				defer sub.Unsubscribe()
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			var blockNum uint64

			for {
				select {
				case resp := <-events:
					require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)

					if blockNum == 0 {
						// first event
						blockNum = resp.BlockNumber
					} else if resp.BlockNumber > blockNum {
						// that means we received events from the next block, and no PRE_CONFIRMED was received. Success!
						return
					}

					assert.Equal(t, TxnFinalityStatusAcceptedOnL2, resp.FinalityStatus)
				case err := <-sub.Err():
					require.NoError(t, err)
				case <-time.After(10 * time.Second):
					t.Fatal("timeout waiting for events")
				}
			}
		})
		t.Run("with finality status PRE_CONFIRMED", func(t *testing.T) {
			t.Parallel()

			events := make(chan *EmittedEventWithFinalityStatus)
			sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
				FinalityStatus: TxnFinalityStatusPre_confirmed,
			})
			if sub != nil {
				defer sub.Unsubscribe()
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			var preConfirmedEventFound bool
			var acceptedOnL2EventFound bool

			for {
				select {
				case resp := <-events:
					require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)

					if preConfirmedEventFound && acceptedOnL2EventFound {
						// subscribing with PRE_CONFIRMED should return both PRE_CONFIRMED and ACCEPTED_ON_L2 events
						return
					}

					switch resp.FinalityStatus {
					case TxnFinalityStatusPre_confirmed:
						preConfirmedEventFound = true
					case TxnFinalityStatusAcceptedOnL2:
						acceptedOnL2EventFound = true
					default:
						t.Fatalf("unexpected finality status: %s", resp.FinalityStatus)
					}
				case err := <-sub.Err():
					require.NoError(t, err)
				case <-time.After(10 * time.Second):
					t.Fatal("timeout waiting for events")
				}
			}
		})
	})

	t.Run("fromAddress + blockID, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			FromAddress: testSet.fromAddressExample,
			SubBlockID:  new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				assert.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				assert.Less(t, resp.BlockNumber, blockNumber)

				assert.Equal(t, testSet.fromAddressExample, resp.FromAddress)

				if resp.BlockNumber >= blockNumber-100 {
					// we searched more than 900 blocks back, it's fine
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("keys + blockID, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			Keys:       [][]*felt.Felt{{testSet.keyExample}},
			SubBlockID: new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
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
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
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

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(context.Background(), events, &EventSubscriptionInput{
			SubBlockID:     new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
			FromAddress:    testSet.fromAddressExample,
			Keys:           [][]*felt.Felt{{testSet.keyExample}},
			FinalityStatus: TxnFinalityStatusAcceptedOnL2,
		})
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		for {
			select {
			case resp := <-events:
				assert.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				assert.Less(t, resp.BlockNumber, blockNumber)
				// 'fromAddressExample' is the address of the sepolia StarkGate: ETH Token, which is very likely to have events,
				// so we can use it to verify the events are returned correctly.
				assert.Equal(t, testSet.fromAddressExample, resp.FromAddress)
				assert.Equal(t, testSet.keyExample, resp.Keys[0])
				assert.Equal(t, TxnFinalityStatusAcceptedOnL2, resp.FinalityStatus)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
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
					SubBlockID: new(SubscriptionBlockID).WithBlockNumber(blockNumber - 2000),
				},
				expectedError: ErrTooManyBlocksBack,
			},
			{
				input: EventSubscriptionInput{
					SubBlockID: new(SubscriptionBlockID).WithBlockNumber(blockNumber + 10000),
				},
				expectedError: ErrBlockNotFound,
			},
		}

		for _, test := range testSet {
			t.Run(test.expectedError.Error(), func(t *testing.T) {
				t.Parallel()

				events := make(chan *EmittedEventWithFinalityStatus)
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

//nolint:dupl,gocyclo
func TestSubscribeNewTransactionReceipts(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)
	wsProvider := testConfig.WsProvider

	t.Run("general cases", func(t *testing.T) {
		t.Parallel()

		type testSetType struct {
			txnReceipts   chan *TransactionReceiptWithBlockInfo
			options       *SubNewTxnReceiptsInput
			expectedError error
			description   string
		}

		addresses := make([]*felt.Felt, 1025)
		for i := range 1025 {
			addresses[i] = internalUtils.TestHexToFelt(t, "0x1")
		}

		testSet := []testSetType{
			{
				txnReceipts: make(chan *TransactionReceiptWithBlockInfo),
				options:     nil,
				description: "nil input",
			},
			{
				txnReceipts: make(chan *TransactionReceiptWithBlockInfo),
				options:     &SubNewTxnReceiptsInput{},
				description: "empty input",
			},
			{
				txnReceipts:   make(chan *TransactionReceiptWithBlockInfo),
				options:       &SubNewTxnReceiptsInput{SenderAddress: addresses},
				expectedError: ErrTooManyAddressesInFilter,
				description:   "error: too many addresses",
			},
		}

		for _, test := range testSet {
			t.Run("test: "+test.description, func(t *testing.T) {
				t.Parallel()

				sub, err := wsProvider.SubscribeNewTransactionReceipts(context.Background(), test.txnReceipts, test.options)
				if test.expectedError != nil {
					require.EqualError(t, err, test.expectedError.Error())

					return
				}
				defer sub.Unsubscribe()

				require.NoError(t, err)
				require.NotNil(t, sub)

				for {
					select {
					case resp := <-test.txnReceipts:
						assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
						assert.Equal(t, TxnFinalityStatusAcceptedOnL2, resp.FinalityStatus) // default finality status is ACCEPTED_ON_L2

						return
					case err := <-sub.Err():
						require.NoError(t, err)
					}
				}
			})
		}
	})

	t.Run("with finality status ACCEPTED_ON_L2", func(t *testing.T) {
		t.Parallel()

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			FinalityStatus: []TxnFinalityStatus{TxnFinalityStatusAcceptedOnL2},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(context.Background(), txnReceipts, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		counter := 0
		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
				assert.Equal(t, TxnFinalityStatusAcceptedOnL2, resp.FinalityStatus)
				assert.NotEmpty(t, resp.BlockHash)
				assert.NotEmpty(t, resp.BlockNumber)
				assert.NotEmpty(t, resp.TransactionReceipt)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with finality status PRE_CONFIRMED", func(t *testing.T) {
		t.Parallel()

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			FinalityStatus: []TxnFinalityStatus{TxnFinalityStatusPre_confirmed},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(context.Background(), txnReceipts, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		counter := 0
		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
				assert.Equal(t, TxnFinalityStatusPre_confirmed, resp.FinalityStatus)
				assert.Empty(t, resp.BlockHash)
				assert.NotEmpty(t, resp.BlockNumber)
				assert.NotEmpty(t, resp.TransactionReceipt)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with both PRE_CONFIRMED and ACCEPTED_ON_L2 finality statuses", func(t *testing.T) {
		t.Parallel()

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			FinalityStatus: []TxnFinalityStatus{TxnFinalityStatusPre_confirmed, TxnFinalityStatusAcceptedOnL2},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(context.Background(), txnReceipts, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		preConfirmedReceived := false
		acceptedOnL2Received := false

		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
				assert.NotEmpty(t, resp.BlockNumber)
				assert.NotEmpty(t, resp.TransactionReceipt)

				if resp.FinalityStatus == TxnFinalityStatusPre_confirmed {
					preConfirmedReceived = true
				}

				if resp.FinalityStatus == TxnFinalityStatusAcceptedOnL2 {
					acceptedOnL2Received = true
				}

				if preConfirmedReceived && acceptedOnL2Received {
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			}
		}
	})

	t.Run("with sender address filter", func(t *testing.T) {
		t.Parallel()

		// and address currently sending a lot of transactions in Sepolia
		randAddress := internalUtils.TestHexToFelt(t, "0x0352057331d5ad77465315d30b98135ddb815b86aa485d659dfeef59a904f88d")
		provider := testConfig.Provider
		tempStruct := struct {
			SenderAddress *felt.Felt `json:"sender_address"`
		}{}

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			SenderAddress: []*felt.Felt{randAddress},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(context.Background(), txnReceipts, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		counter := 0
		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)

				txn, err := provider.TransactionByHash(context.Background(), resp.Hash)
				require.NoError(t, err)

				raw, err := json.Marshal(txn)
				require.NoError(t, err)

				err = json.Unmarshal(raw, &tempStruct)
				require.NoError(t, err)

				assert.Equal(t, randAddress, tempStruct.SenderAddress)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.Greater(t, counter, 0, "no receipts received")

				return
			}
		}
	})
}

//nolint:dupl,gocyclo
func TestSubscribeNewTransactions(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)
	wsProvider := testConfig.WsProvider

	t.Run("general cases", func(t *testing.T) {
		t.Parallel()

		type testSetType struct {
			newTxns       chan *TxnWithHashAndStatus
			options       *SubNewTxnsInput
			expectedError error
			description   string
		}

		addresses := make([]*felt.Felt, 1025)
		for i := range 1025 {
			addresses[i] = internalUtils.TestHexToFelt(t, "0x1")
		}

		testSet := []testSetType{
			{
				newTxns:     make(chan *TxnWithHashAndStatus),
				options:     nil,
				description: "nil input",
			},
			{
				newTxns:     make(chan *TxnWithHashAndStatus),
				options:     &SubNewTxnsInput{},
				description: "empty input",
			},
			{
				newTxns:       make(chan *TxnWithHashAndStatus),
				options:       &SubNewTxnsInput{SenderAddress: addresses},
				expectedError: ErrTooManyAddressesInFilter,
				description:   "error: too many addresses",
			},
		}

		for _, test := range testSet {
			t.Run("test: "+test.description, func(t *testing.T) {
				t.Parallel()

				sub, err := wsProvider.SubscribeNewTransactions(context.Background(), test.newTxns, test.options)
				if test.expectedError != nil {
					require.EqualError(t, err, test.expectedError.Error())

					return
				}
				defer sub.Unsubscribe()

				require.NoError(t, err)
				require.NotNil(t, sub)

				for {
					select {
					case resp := <-test.newTxns:
						assert.IsType(t, &TxnWithHashAndStatus{}, resp)
						assert.Equal(t, TxnStatus_Accepted_On_L2, resp.FinalityStatus) // default finality status is ACCEPTED_ON_L2

						return
					case <-time.After(10 * time.Second):
						assert.Fail(t, "no txns received within timeout")

						return
					case err := <-sub.Err():
						require.NoError(t, err)
					}
				}
			})
		}
	})

	t.Run("with finality status ACCEPTED_ON_L2", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{TxnStatus_Accepted_On_L2},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		counter := 0
		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.Equal(t, TxnStatus_Accepted_On_L2, resp.FinalityStatus)
				assert.NotEmpty(t, resp.Transaction)
				assert.NotEmpty(t, resp.Hash)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with finality status PRE_CONFIRMED", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{TxnStatus_Pre_confirmed},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		counter := 0
		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.Equal(t, TxnStatus_Pre_confirmed, resp.FinalityStatus)
				assert.NotEmpty(t, resp.Hash)
				assert.NotEmpty(t, resp.Transaction)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with both PRE_CONFIRMED and ACCEPTED_ON_L2 finality statuses", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{TxnStatus_Pre_confirmed, TxnStatus_Accepted_On_L2},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		preConfirmedReceived := false
		acceptedOnL2Received := false

		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.NotEmpty(t, resp.Hash)
				assert.NotEmpty(t, resp.Transaction)

				if resp.FinalityStatus == TxnStatus_Pre_confirmed {
					preConfirmedReceived = true
				}

				if resp.FinalityStatus == TxnStatus_Accepted_On_L2 {
					acceptedOnL2Received = true
				}

				if preConfirmedReceived && acceptedOnL2Received {
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.True(t, (preConfirmedReceived && acceptedOnL2Received), "no txns received from both finality statuses")

				return
			}
		}
	})

	t.Run("with all finality statuses, except ACCEPTED_ON_L1", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{TxnStatus_Received, TxnStatus_Candidate, TxnStatus_Pre_confirmed, TxnStatus_Accepted_On_L2},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		receivedReceived := false
		candidateReceived := false
		preConfirmedReceived := false
		acceptedOnL2Received := false

		timeout := time.After(10 * time.Second)

		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.NotEmpty(t, resp.Hash)
				assert.NotEmpty(t, resp.Transaction)

				switch resp.FinalityStatus {
				case TxnStatus_Received:
					t.Log("RECEIVED txn received")
					receivedReceived = true
				case TxnStatus_Candidate:
					t.Log("CANDIDATE txn received")
					candidateReceived = true
				case TxnStatus_Pre_confirmed:
					t.Log("PRE_CONFIRMED txn received")
					preConfirmedReceived = true
				case TxnStatus_Accepted_On_L2:
					t.Log("ACCEPTED_ON_L2 txn received")
					acceptedOnL2Received = true
				}

			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				assert.True(
					t,
					(receivedReceived || candidateReceived || preConfirmedReceived || acceptedOnL2Received),
					"no txns received",
				)

				return
			}
		}
	})

	t.Run("with sender address filter", func(t *testing.T) {
		t.Parallel()

		// and address currently sending a lot of transactions in Sepolia
		randAddress := internalUtils.TestHexToFelt(t, "0x0352057331d5ad77465315d30b98135ddb815b86aa485d659dfeef59a904f88d")
		provider := testConfig.Provider
		tempStruct := struct {
			SenderAddress *felt.Felt `json:"sender_address"`
		}{}

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			SenderAddress: []*felt.Felt{randAddress},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		counter := 0
		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)

				txn, err := provider.TransactionByHash(context.Background(), resp.Hash)
				require.NoError(t, err)

				raw, err := json.Marshal(txn)
				require.NoError(t, err)

				err = json.Unmarshal(raw, &tempStruct)
				require.NoError(t, err)

				assert.Equal(t, randAddress, tempStruct.SenderAddress)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				assert.Greater(t, counter, 0, "no receipts received")

				return
			}
		}
	})
}

func TestUnsubscribe(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv, tests.IntegrationEnv)

	t.Parallel()

	testConfig := BeforeEach(t, true)

	wsProvider := testConfig.WsProvider

	events := make(chan *EmittedEventWithFinalityStatus)
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
			require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
		case err := <-sub.Err():
			// when unsubscribing, the error channel should return nil
			require.Nil(t, err)

			break loop
		case <-time.After(10 * time.Second):
			t.Fatal("timeout waiting for unsubscription")
		}
	}
}

// A simple test was made to make sure the reorg events are received. Ref:
// https://github.com/NethermindEth/starknet.go/pull/651#discussion_r1927356194
func TestReorgEvents(t *testing.T) {
	t.Skip("TODO: implement reorg test")
}
