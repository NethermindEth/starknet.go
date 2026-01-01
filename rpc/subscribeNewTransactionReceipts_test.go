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

// TestSubscribeNewTransactionReceipts tests the SubscribeNewTransactionReceipts function.
func TestSubscribeNewTransactionReceipts(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.TestnetEnv,
	)
	t.Parallel()

	type testSetType struct {
		description   string
		input         *SubNewTxnReceiptsInput
		expectedError error
	}

	tooManyAddresses := make([]*felt.Felt, 10000)
	for i := range 10000 {
		tooManyAddresses[i] = new(felt.Felt).SetUint64(uint64(i))
	}

	template := []testSetType{
		{
			description: "from address only",
			input: &SubNewTxnReceiptsInput{
				SenderAddress: []*felt.Felt{randAddress},
			},
		},
		{
			description: "with finality status PRE_CONFIRMED",
			input: &SubNewTxnReceiptsInput{
				FinalityStatus: []TxnFinalityStatus{TxnFinalityStatusPreConfirmed},
			},
		},
		{
			description: "with finality status ACCEPTED_ON_L2",
			input: &SubNewTxnReceiptsInput{
				FinalityStatus: []TxnFinalityStatus{TxnFinalityStatusAcceptedOnL2},
			},
		},
		{
			description: "all filters",
			input: &SubNewTxnReceiptsInput{
				SenderAddress: []*felt.Felt{randAddress},
				FinalityStatus: []TxnFinalityStatus{
					TxnFinalityStatusAcceptedOnL2,
					TxnFinalityStatusPreConfirmed,
				},
			},
		},
		{
			description: "error: too many addresses",
			input: &SubNewTxnReceiptsInput{
				SenderAddress: tooManyAddresses,
			},
			expectedError: ErrTooManyAddressesInFilter,
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
						"_subscribeNewTransactionReceipts",
						gomock.Any(),
						test.input,
					).
					DoAndReturn(func(_, _, _, channel any, arg any) (*client.ClientSubscription, error) {
						ch := channel.(chan json.RawMessage)
						input := arg.(*SubNewTxnReceiptsInput)

						if len(input.SenderAddress) > 1000 {
							return nil, RPCError{
								Code:    67,
								Message: "Too many addresses in filter sender_address filter",
							}
						}

						msg := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
							t,
							"./testData/ws/sepoliaNewTxnReceipts.json",
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

			receipts := make(chan *TransactionReceiptWithBlockInfo)
			sub, err := tsetup.WsProvider.SubscribeNewTransactionReceipts(
				t.Context(),
				receipts,
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
				case resp := <-receipts:
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
					"_subscribeNewTransactionReceipts",
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _, _, channel any, _ any) (*client.ClientSubscription, error) {
					ch := channel.(chan json.RawMessage)

					msg := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
						t,
						"./testData/ws/sepoliaNewTxnReceipts.json",
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

		receipts := make(chan *TransactionReceiptWithBlockInfo)
		sub, err := tsetup.WsProvider.SubscribeNewTransactionReceipts(
			t.Context(),
			receipts,
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
			case resp := <-receipts:
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
