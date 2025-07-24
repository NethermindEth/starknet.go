package rpc_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// returns a new account type from the provided account data in the tConfig
func setupAcc(t *testing.T, tConfig *rpc.TestConfiguration) (*account.Account, error) {
	t.Helper()

	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(tConfig.PrivKey, 0)
	if !ok {
		return nil, errors.New("failed to convert privKey to big.Int")
	}
	ks.Put(tConfig.PubKey, privKeyBI)

	accAddress, err := internalUtils.HexToFelt(tConfig.AccountAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert accountAddress to felt: %w", err)
	}

	acc, err := account.NewAccount(tConfig.Provider, accAddress, tConfig.PubKey, ks, account.CairoV2)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return acc, nil
}

//nolint:tparallel
func TestSubscribeTransactionStatus(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	tConfig := rpc.BeforeEach(t, true)

	// ***** 1 - setup provider and account
	provider := tConfig.Provider
	wsProvider := tConfig.WsProvider

	acc, err := setupAcc(t, tConfig)
	require.NoError(t, err, "Error in setupAcc")

	nonce, err := acc.Nonce(context.Background())
	require.NoError(t, err, "Error getting nonce")

	// ***** 2 - build, sign, and get the hash of the txn

	calldata, err := acc.FmtCalldata([]rpc.FunctionCall{
		{
			// same ERC20 contract as in examples/simpleInvoke
			ContractAddress:    internalUtils.TestHexToFelt(t, "0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54"),
			EntryPointSelector: utils.GetSelectorFromNameFelt("mint"),
			Calldata:           []*felt.Felt{new(felt.Felt).SetUint64(10000), &felt.Zero},
		},
	})

	require.NoError(t, err, "Error building and sending invoke txn")

	invokeTx := utils.BuildInvokeTxn(
		acc.Address,
		nonce,
		calldata,
		&rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		nil,
	)

	err = acc.SignInvokeTransaction(context.Background(), invokeTx)
	require.NoError(t, err, "Error signing invoke txn")

	estimateFee, err := provider.EstimateFee(
		context.Background(),
		[]rpc.BroadcastTxn{invokeTx},
		[]rpc.SimulationFlag{},
		rpc.WithBlockTag(rpc.BlockTagLatest),
	)
	require.NoError(t, err, "Error estimating fee")

	txnFee := estimateFee[0]
	invokeTx.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, 1.5)

	// sign the txn
	err = acc.SignInvokeTransaction(context.Background(), invokeTx)
	require.NoError(t, err, "Error signing invoke txn")

	// get the hash of the txn
	txnHash, err := acc.TransactionHashInvoke(invokeTx)
	t.Logf("Txn hash: %s", txnHash)
	require.NoError(t, err, "Error getting txn hash")

	// ***** 3 - subscribe to txn status in parallel

	t.Run("subscribe to txn status", func(t *testing.T) {
		t.Parallel()

		txnStatus := make(chan *rpc.NewTxnStatus)
		sub, err := wsProvider.SubscribeTransactionStatus(context.Background(), txnStatus, txnHash) //nolint:govet
		require.NoError(t, err, "Error subscribing to txn status")
		defer sub.Unsubscribe()

		expectedStatus := rpc.TxnStatus_Received

		for {
			select {
			case txnStatus := <-txnStatus:
				// since Juno will only return the current status, skipping previous statuses
				// (e.g: it can go directly from RECEIVED to ACCEPTED_ON_L2),
				// we'll only check if the txn has been marked at least as received
				switch txnStatus.Status.FinalityStatus {
				case rpc.TxnStatus_Received:
					t.Logf("Txn status: %v", txnStatus.Status.FinalityStatus)
					assert.Equal(t, expectedStatus, rpc.TxnStatus_Received)

					expectedStatus = rpc.TxnStatus_Candidate
				case rpc.TxnStatus_Candidate:
					t.Logf("Txn status: %v", txnStatus.Status.FinalityStatus)
					assert.NotEqual(t, expectedStatus, rpc.TxnStatus_Received, "txn should have been marked as received first")

					expectedStatus = rpc.TxnStatus_Pre_confirmed
				case rpc.TxnStatus_Pre_confirmed:
					t.Logf("Txn status: %v", txnStatus.Status.FinalityStatus)
					assert.NotEqual(t, expectedStatus, rpc.TxnStatus_Received, "txn should have been marked as received first")

					expectedStatus = rpc.TxnStatus_Accepted_On_L2
				case rpc.TxnStatus_Accepted_On_L2:
					t.Logf("Txn status: %v", txnStatus.Status.FinalityStatus)
					assert.NotEqual(t, expectedStatus, rpc.TxnStatus_Received, "txn should have been marked as received first")

					return
				}
			case err = <-sub.Err():
				t.Fatal("error in subscription: ", err)
			}
		}
	})

	// wait for the subscription to start
	time.Sleep(3 * time.Second)

	// ***** 4 - send the txn

	_, err = provider.AddInvokeTransaction(context.Background(), invokeTx)
	require.NoError(t, err, "Error adding invoke txn")
}
