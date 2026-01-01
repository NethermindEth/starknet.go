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

// TestSubscribeTransactionStatus tests the SubscribeTransactionStatus function.
func TestSubscribeTransactionStatus(t *testing.T) {
	tests.RunTestOn(t,
		tests.MockEnv,
		tests.IntegrationEnv,
		tests.MainnetEnv,
		tests.TestnetEnv,
	)
	t.Parallel()

	t.Run("network test", func(t *testing.T) {
		tests.RunTestOn(t,
			tests.IntegrationEnv,
			tests.MainnetEnv,
			tests.TestnetEnv,
		)
		t.Parallel()

		// getting a random new PRE_CONFIRMED transaction
		tempSetup := BeforeEach(t, true) // to avoid race condition (BeforeEach
		// must be called once per subscription)
		txns := make(chan *TxnWithHashAndStatus)
		sub, err := tempSetup.WsProvider.SubscribeNewTransactions(
			t.Context(),
			txns,
			&SubNewTxnsInput{
				FinalityStatus: []TxnStatus{TxnStatusPreConfirmed},
			},
		)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		txn := <-txns
		require.NotNil(t, txn)

		tsetup := BeforeEach(t, true)

		status := make(chan *NewTxnStatus)
		sub2, err := tsetup.WsProvider.SubscribeTransactionStatus(
			t.Context(),
			status,
			txn.Hash,
		)
		require.NoError(t, err)
		defer sub2.Unsubscribe()

		for {
			select {
			case status := <-status:
				rawExpectedMsg := <-tsetup.WSSpy.SpyChannel()
				rawMsg, err := json.Marshal(status)
				require.NoError(t, err)
				assert.JSONEq(t, string(rawExpectedMsg), string(rawMsg))

				return
			case err := <-sub2.Err():
				require.NoError(t, err)

				return
			case <-time.After(testDuration * 2):
				t.Fatal("timeout waiting for status")

				return
			}
		}
	})

	t.Run("mock tests", func(t *testing.T) {
		tests.RunTestOn(t, tests.MockEnv)
		t.Parallel()

		testSet := []*felt.Felt{
			new(felt.Felt).SetUint64(1), // RECEIVED
			new(felt.Felt).SetUint64(2), // CANDIDATE
			new(felt.Felt).SetUint64(3), // PRE_CONFIRMED
			new(felt.Felt).SetUint64(4), // ACCEPTED_ON_L2
			new(felt.Felt).SetUint64(5), // ACCEPTED_ON_L1
			new(felt.Felt).SetUint64(6), // REVERTED
		}

		for _, hash := range testSet {
			t.Run(hash.String(), func(t *testing.T) {
				t.Parallel()
				tsetup := BeforeEach(t, true)

				tsetup.MockClient.EXPECT().
					SubscribeWithSliceArgs(
						t.Context(),
						"starknet",
						"_subscribeTransactionStatus",
						gomock.Any(),
						gomock.Any(),
					).
					DoAndReturn(func(_, _, _, channel any, args ...any) (*client.ClientSubscription, error) {
						ch := channel.(chan json.RawMessage)
						hash := args[0].(*felt.Felt)
						var msg string

						switch hash.Uint64() {
						case 1:
							msg = `{
								"transaction_hash": "0x7ef307fc37cf8cfec75ba0ac88fa347bb6e4ff0b0937d45d64bde1cbe05a95",
								"status": {
									"finality_status": "RECEIVED",
									"execution_status": "SUCCEEDED"
								}}`
						case 2:
							msg = `{
								"transaction_hash": "0x7ef307fc37cf8cfec75ba0ac88fa347bb6e4ff0b0937d45d64bde1cbe05a95",
								"status": {
									"finality_status": "CANDIDATE",
									"execution_status": "SUCCEEDED"
								}}`
						case 3:
							msg = `{
								"transaction_hash": "0x7ef307fc37cf8cfec75ba0ac88fa347bb6e4ff0b0937d45d64bde1cbe05a95",
								"status": {
									"finality_status": "PRE_CONFIRMED",
									"execution_status": "SUCCEEDED"
								}}`
						case 4:
							msg = `{
								"transaction_hash": "0x7ef307fc37cf8cfec75ba0ac88fa347bb6e4ff0b0937d45d64bde1cbe05a95",
								"status": {
									"finality_status": "ACCEPTED_ON_L2",
									"execution_status": "SUCCEEDED"
								}}`
						case 5:
							msg = `{
								"transaction_hash": "0x7ef307fc37cf8cfec75ba0ac88fa347bb6e4ff0b0937d45d64bde1cbe05a95",
								"status": {
									"finality_status": "ACCEPTED_ON_L1",
									"execution_status": "SUCCEEDED"
								}}`
						case 6:
							rawMsg := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
								t,
								"./testData/ws/mainnetTxnStatus.json",
								"params", "result",
							)
							msg = string(rawMsg)
						}

						go func() {
							for {
								select {
								case <-time.Tick(2 * time.Second):
									ch <- json.RawMessage(msg)
								case <-t.Context().Done():
									return
								}
							}
						}()

						return &client.ClientSubscription{}, nil
					}).
					Times(1)

				status := make(chan *NewTxnStatus)
				sub, err := tsetup.WsProvider.SubscribeTransactionStatus(
					t.Context(),
					status,
					hash,
				)
				require.NoError(t, err)

				select {
				case status := <-status:
					rawExpectedMsg := <-tsetup.WSSpy.SpyChannel()
					rawMsg, err := json.Marshal(status)
					require.NoError(t, err)
					assert.JSONEq(t, string(rawExpectedMsg), string(rawMsg))

					return
				case err := <-sub.Err():
					require.NoError(t, err)

					return
				case <-time.After(testDuration * 2):
					t.Fatal("timeout waiting for status")

					return
				}
			})
		}
	})
}
