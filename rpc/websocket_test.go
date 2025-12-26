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

			timeout := time.After(10 * time.Second)
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
					case <-timeout:
						return
					default:
					}
				case err := <-sub.Err():
					require.NoError(t, err)
				}
			}
		})
	}
}

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

	// TODO for all the cases: add logic to marshal get request type and compare with the raw request sent to the RPC server.
	// Maybe a websocket spy could help here.

	t.Run("with empty args", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(
			context.Background(),
			events,
			&EventSubscriptionInput{},
		)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		// outside the loop, to avoid it being resetted
		timeout := time.After(10 * time.Second)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("blockID only - 1000 blocks back", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(
			context.Background(),
			events,
			&EventSubscriptionInput{
				SubBlockID: new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
			},
		)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		// outside the loop, to avoid it being resetted
		timeout := time.After(10 * time.Second)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Fatal("timeout waiting for events")
			}
		}
	})

	t.Run("blockID only - with tag latest", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(
			context.Background(),
			events,
			&EventSubscriptionInput{
				SubBlockID: new(SubscriptionBlockID).WithLatestTag(),
			},
		)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		// outside the loop, to avoid it being resetted
		timeout := time.After(10 * time.Second)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
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
			sub, err := wsProvider.SubscribeEvents(
				context.Background(),
				events,
				&EventSubscriptionInput{
					FinalityStatus: TxnFinalityStatusAcceptedOnL2,
				},
			)
			if sub != nil {
				defer sub.Unsubscribe()
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			// outside the loop, to avoid it being resetted
			timeout := time.After(10 * time.Second)

			for {
				select {
				case resp := <-events:
					require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
					assert.Equal(t, TxnFinalityStatusAcceptedOnL2, resp.FinalityStatus)

					return
				case err := <-sub.Err():
					require.NoError(t, err)
				case <-timeout:
					t.Fatal("timeout waiting for events")
				}
			}
		})

		t.Run("with finality status PRE_CONFIRMED", func(t *testing.T) {
			t.Parallel()

			events := make(chan *EmittedEventWithFinalityStatus)
			sub, err := wsProvider.SubscribeEvents(
				context.Background(),
				events,
				&EventSubscriptionInput{
					FinalityStatus: TxnFinalityStatusPreConfirmed,
				},
			)
			if sub != nil {
				defer sub.Unsubscribe()
			}
			require.NoError(t, err)
			require.NotNil(t, sub)

			// outside the loop, to avoid it being resetted
			timeout := time.After(10 * time.Second)

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
					case TxnFinalityStatusPreConfirmed:
						preConfirmedEventFound = true
					case TxnFinalityStatusAcceptedOnL2:
						acceptedOnL2EventFound = true
					default:
						t.Fatalf("unexpected finality status: %s", resp.FinalityStatus)
					}
				case err := <-sub.Err():
					require.NoError(t, err)
				case <-timeout:
					t.Fatal("timeout waiting for events")
				}
			}
		})
	})

	t.Run("fromAddress + blockID, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(
			context.Background(),
			events,
			&EventSubscriptionInput{
				FromAddress: testSet.fromAddressExample,
				SubBlockID:  new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
			},
		)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		// outside the loop, to avoid it being resetted
		timeout := time.After(10 * time.Second)

		for {
			select {
			case resp := <-events:
				assert.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				assert.Less(t, resp.BlockNumber, blockNumber)

				assert.Equal(t, testSet.fromAddressExample, resp.FromAddress)

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Skip("timeout reached, no events received")
			}
		}
	})

	t.Run("keys + blockID, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(
			context.Background(),
			events,
			&EventSubscriptionInput{
				Keys:       [][]*felt.Felt{{testSet.keyExample}},
				SubBlockID: new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
			},
		)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		// outside the loop, to avoid it being resetted
		timeout := time.After(20 * time.Second)

		for {
			select {
			case resp := <-events:
				require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				require.Less(t, resp.BlockNumber, blockNumber)

				// Subscription with keys should only return events with the specified keys.
				require.Equal(t, testSet.keyExample, resp.Keys[0])

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Skip("timeout reached, no events received")
			}
		}
	})

	t.Run("with all arguments, within the range of 1024 blocks", func(t *testing.T) {
		t.Parallel()

		wsProvider := testConfig.WsProvider

		events := make(chan *EmittedEventWithFinalityStatus)
		sub, err := wsProvider.SubscribeEvents(
			context.Background(),
			events,
			&EventSubscriptionInput{
				SubBlockID:     new(SubscriptionBlockID).WithBlockNumber(blockNumber - 1000),
				FromAddress:    testSet.fromAddressExample,
				Keys:           [][]*felt.Felt{{testSet.keyExample}},
				FinalityStatus: TxnFinalityStatusAcceptedOnL2,
			},
		)
		if sub != nil {
			defer sub.Unsubscribe()
		}
		require.NoError(t, err)
		require.NotNil(t, sub)

		// outside the loop, to avoid it being resetted
		timeout := time.After(20 * time.Second)

		for {
			select {
			case resp := <-events:
				assert.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
				assert.Less(t, resp.BlockNumber, blockNumber)
				// 'fromAddressExample' is the address of the sepolia StarkGate: ETH Token, which is very likely to have events,
				// so we can use it to verify the events are returned correctly.
				assert.Equal(t, testSet.fromAddressExample, resp.FromAddress)
				assert.Equal(t, testSet.keyExample, resp.Keys[0])

				return
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Skip("timeout reached, no events received")
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

				sub, err := wsProvider.SubscribeNewTransactionReceipts(
					context.Background(),
					test.txnReceipts,
					test.options,
				)
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
						assert.Equal(
							t,
							TxnFinalityStatusAcceptedOnL2,
							resp.FinalityStatus,
						) // default finality status is ACCEPTED_ON_L2

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

		sub, err := wsProvider.SubscribeNewTransactionReceipts(
			context.Background(),
			txnReceipts,
			options,
		)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		timeout := time.After(20 * time.Second)

		counter := 0
		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
				assert.Equal(t, TxnFinalityStatusAcceptedOnL2, resp.FinalityStatus)
				assert.NotEmpty(t, resp.BlockNumber)
				assert.NotEmpty(t, resp.TransactionReceipt)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with finality status PRE_CONFIRMED", func(t *testing.T) {
		t.Parallel()

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			FinalityStatus: []TxnFinalityStatus{TxnFinalityStatusPreConfirmed},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(
			context.Background(),
			txnReceipts,
			options,
		)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		timeout := time.After(10 * time.Second)

		counter := 0
		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
				assert.Equal(t, TxnFinalityStatusPreConfirmed, resp.FinalityStatus)
				assert.Empty(t, resp.BlockHash)
				assert.NotEmpty(t, resp.BlockNumber)
				assert.NotEmpty(t, resp.TransactionReceipt)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with both PRE_CONFIRMED and ACCEPTED_ON_L2 finality statuses", func(t *testing.T) {
		t.Parallel()

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			FinalityStatus: []TxnFinalityStatus{
				TxnFinalityStatusPreConfirmed,
				TxnFinalityStatusAcceptedOnL2,
			},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(
			context.Background(),
			txnReceipts,
			options,
		)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		preConfirmedReceived := false
		acceptedOnL2Received := false

		timeout := time.After(20 * time.Second)

		for {
			select {
			case resp := <-txnReceipts:
				assert.IsType(t, &TransactionReceiptWithBlockInfo{}, resp)
				assert.NotEmpty(t, resp.BlockNumber)
				assert.NotEmpty(t, resp.TransactionReceipt)

				if resp.FinalityStatus == TxnFinalityStatusPreConfirmed {
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

			case <-timeout:
				assert.True(
					t,
					(preConfirmedReceived && acceptedOnL2Received),
					"no txns received from both finality statuses",
				)

				return
			}
		}
	})

	t.Run("with sender address filter", func(t *testing.T) {
		t.Parallel()

		// and address currently sending a lot of transactions in Sepolia
		randAddress := internalUtils.TestHexToFelt(
			t,
			"0x00395a96a5b6343fc0f543692fd36e7034b54c2a276cd1a021e8c0b02aee1f43",
		)
		provider := testConfig.Provider
		tempStruct := struct {
			SenderAddress *felt.Felt `json:"sender_address"`
		}{}

		txnReceipts := make(chan *TransactionReceiptWithBlockInfo)
		options := &SubNewTxnReceiptsInput{
			SenderAddress: []*felt.Felt{randAddress},
		}

		sub, err := wsProvider.SubscribeNewTransactionReceipts(
			context.Background(),
			txnReceipts,
			options,
		)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		timeout := time.After(10 * time.Second)

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
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Skip("no txns received")
			}
		}
	})
}

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

				sub, err := wsProvider.SubscribeNewTransactions(
					context.Background(),
					test.newTxns,
					test.options,
				)
				if test.expectedError != nil {
					require.EqualError(t, err, test.expectedError.Error())

					return
				}
				defer sub.Unsubscribe()

				require.NoError(t, err)
				require.NotNil(t, sub)

				timeout := time.After(20 * time.Second)

				for {
					select {
					case resp := <-test.newTxns:
						assert.IsType(t, &TxnWithHashAndStatus{}, resp)
						assert.Equal(
							t,
							TxnStatusAcceptedOnL2,
							resp.FinalityStatus,
						) // default finality status is ACCEPTED_ON_L2

						return
					case <-timeout:
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
			FinalityStatus: []TxnStatus{TxnStatusAcceptedOnL2},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		// outside the loop, to avoid it being resetted
		timeout := time.After(20 * time.Second)

		counter := 0
		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.Equal(t, TxnStatusAcceptedOnL2, resp.FinalityStatus)
				assert.NotEmpty(t, resp.Transaction)
				assert.NotEmpty(t, resp.Hash)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with finality status PRE_CONFIRMED", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{TxnStatusPreConfirmed},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		// outside the loop, to avoid it being resetted
		timeout := time.After(20 * time.Second)

		counter := 0
		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.Equal(t, TxnStatusPreConfirmed, resp.FinalityStatus)
				assert.NotEmpty(t, resp.Hash)
				assert.NotEmpty(t, resp.Transaction)

				counter++
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				assert.Greater(t, counter, 0, "no txns received")

				return
			}
		}
	})

	t.Run("with both PRE_CONFIRMED and ACCEPTED_ON_L2 finality statuses", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{TxnStatusPreConfirmed, TxnStatusAcceptedOnL2},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		preConfirmedReceived := false
		acceptedOnL2Received := false

		// outside the loop, to avoid it being resetted
		timeout := time.After(20 * time.Second)

		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.NotEmpty(t, resp.Hash)
				assert.NotEmpty(t, resp.Transaction)

				if resp.FinalityStatus == TxnStatusPreConfirmed {
					preConfirmedReceived = true
				}

				if resp.FinalityStatus == TxnStatusAcceptedOnL2 {
					acceptedOnL2Received = true
				}

				if preConfirmedReceived && acceptedOnL2Received {
					return
				}
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				assert.True(
					t,
					(preConfirmedReceived && acceptedOnL2Received),
					"no txns received from both finality statuses",
				)

				return
			}
		}
	})

	t.Run("with all finality statuses, except ACCEPTED_ON_L1", func(t *testing.T) {
		t.Parallel()

		newTxns := make(chan *TxnWithHashAndStatus)
		options := &SubNewTxnsInput{
			FinalityStatus: []TxnStatus{
				TxnStatusReceived,
				TxnStatusCandidate,
				TxnStatusPreConfirmed,
				TxnStatusAcceptedOnL2,
			},
		}

		sub, err := wsProvider.SubscribeNewTransactions(context.Background(), newTxns, options)
		require.NoError(t, err)
		require.NotNil(t, sub)

		defer sub.Unsubscribe()

		receivedReceived := false
		candidateReceived := false
		preConfirmedReceived := false
		acceptedOnL2Received := false

		// outside the loop, to avoid it being resetted
		timeout := time.After(20 * time.Second)

		for {
			select {
			case resp := <-newTxns:
				assert.IsType(t, &TxnWithHashAndStatus{}, resp)
				assert.NotEmpty(t, resp.Hash)
				assert.NotEmpty(t, resp.Transaction)

				switch resp.FinalityStatus {
				case TxnStatusReceived:
					t.Log("RECEIVED txn received")
					receivedReceived = true
				case TxnStatusCandidate:
					t.Log("CANDIDATE txn received")
					candidateReceived = true
				case TxnStatusPreConfirmed:
					t.Log("PRE_CONFIRMED txn received")
					preConfirmedReceived = true
				case TxnStatusAcceptedOnL2:
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
		randAddress := internalUtils.TestHexToFelt(
			t,
			"0x00395a96a5b6343fc0f543692fd36e7034b54c2a276cd1a021e8c0b02aee1f43",
		)
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

		timeout := time.After(20 * time.Second)

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
			case err := <-sub.Err():
				require.NoError(t, err)
			case <-timeout:
				t.Skip("no txns received")
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

	timeout := time.After(10 * time.Second)

loop:

	for {
		select {
		case resp := <-events:
			require.IsType(t, &EmittedEventWithFinalityStatus{}, resp)
		case err := <-sub.Err():
			// when unsubscribing, the error channel should return nil
			require.Nil(t, err)

			break loop
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
