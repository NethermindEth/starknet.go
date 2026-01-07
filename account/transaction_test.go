package account_test

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/tests/mocks/rpcv10mock"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestBuildAndSendInvokeTxn is a test function that tests the BuildAndSendInvokeTxn method.
//
// This function tests the BuildAndSendInvokeTxn method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
func TestBuildAndSendInvokeTxn(t *testing.T) {
	// TODO: implement devnet support
	tests.RunTestOn(t, tests.TestnetEnv)

	provider, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Build and send invoke txn
	resp, err := acc.BuildAndSendInvokeTxn(t.Context(), []rpc.InvokeFunctionCall{
		{
			// same ERC20 contract as in examples/simpleInvoke
			ContractAddress: internalUtils.TestHexToFelt(
				t,
				"0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54",
			),
			FunctionName: "mint",
			CallData:     []*felt.Felt{new(felt.Felt).SetUint64(10000), &felt.Zero},
		},
	}, nil)
	require.NoError(t, err, "Error building and sending invoke txn")

	// check the transaction hash
	require.NotNil(t, resp.Hash)
	t.Logf("Invoke transaction hash: %s", resp.Hash)

	txReceipt, err := acc.WaitForTransactionReceipt(
		t.Context(),
		resp.Hash,
		500*time.Millisecond,
	)
	require.NoError(t, err, "Error waiting for invoke transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)

	// testing the default tip estimation feature
	txn, err := acc.Provider.TransactionByHash(t.Context(), resp.Hash)
	require.NoError(t, err, "Error getting transaction by hash")
	require.NotNil(t, txn)
	assert.NotEqual(t, "0x0", txn.Transaction.(rpc.InvokeTxnV3).Tip)
}

// TestBuildAndSendDeclareTxn is a test function that tests the BuildAndSendDeclareTxn method.
//
// This function tests the BuildAndSendDeclareTxn method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
func TestBuildAndSendDeclareTxn(t *testing.T) {
	// TODO: implement devnet support
	tests.RunTestOn(t, tests.TestnetEnv)

	provider, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Class
	class := internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](
		t,
		"./testData/contracts_v2_HelloStarknet.sierra.json",
	)

	// Casm Class
	casmClass := internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](
		t,
		"./testData/contracts_v2_HelloStarknet.casm.json",
	)

	// Build and send declare txn
	resp, err := acc.BuildAndSendDeclareTxn(
		t.Context(),
		&casmClass,
		&class,
		nil,
	)
	if err != nil {
		require.EqualError(
			t,
			err,
			"41 Transaction execution error: Class with hash 0x0224518978adb773cfd4862a894e9d333192fbd24bc83841dc7d4167c09b89c5 is already declared.",
		)
		t.Log("declare txn not sent: class already declared")

		return
	}

	// check the transaction and class hash
	require.NotNil(t, resp.Hash)
	require.NotNil(t, resp.ClassHash)
	t.Logf("Declare transaction hash: %s", resp.Hash)
	t.Logf("Class hash: %s", resp.ClassHash)

	txReceipt, err := acc.WaitForTransactionReceipt(
		t.Context(),
		resp.Hash,
		500*time.Millisecond,
	)
	require.NoError(t, err, "Error waiting for declare transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)

	// testing the default tip estimation feature
	txn, err := acc.Provider.TransactionByHash(t.Context(), resp.Hash)
	require.NoError(t, err, "Error getting transaction by hash")
	require.NotNil(t, txn)
	assert.NotEqual(t, "0x0", txn.Transaction.(rpc.DeclareTxnV3).Tip)
}

func TestBuildAndSendDeclareTxnMock(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	// Class
	class := internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](
		t,
		"./testData/contracts_v2_HelloStarknet.sierra.json",
	)
	// Casm Class
	casmClass := internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](
		t,
		"./testData/contracts_v2_HelloStarknet.casm.json",
	)

	t.Run("compiled class hash", func(t *testing.T) {
		testcases := []struct {
			name                     string
			txnOptions               *account.TxnOptions
			starknetVersion          string
			expectedCompileClassHash string
		}{
			{
				name:                     "before 0.14.1",
				txnOptions:               nil,
				starknetVersion:          "0.14.0",
				expectedCompileClassHash: "0x6ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17",
			},
			{
				name: "before 0.14.1 + UseBlake2sHash true",
				txnOptions: &account.TxnOptions{
					UseBlake2sHash: &[]bool{true}[0],
				},
				starknetVersion:          "0.14.0",
				expectedCompileClassHash: "0x23c2091df2547f77185ba592b06ee2e897b0c2a70f968521a6a24fc5bfc1b1e",
			},
			{
				name: "before 0.14.1 + UseBlake2sHash false",
				txnOptions: &account.TxnOptions{
					UseBlake2sHash: &[]bool{false}[0],
				},
				starknetVersion:          "0.14.0",
				expectedCompileClassHash: "0x6ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17",
			},
			{
				name:                     "after 0.14.1",
				txnOptions:               nil,
				starknetVersion:          "0.14.1",
				expectedCompileClassHash: "0x23c2091df2547f77185ba592b06ee2e897b0c2a70f968521a6a24fc5bfc1b1e",
			},
			{
				name: "after 0.14.1 + UseBlake2sHash true",
				txnOptions: &account.TxnOptions{
					UseBlake2sHash: &[]bool{true}[0],
				},
				starknetVersion:          "0.14.1",
				expectedCompileClassHash: "0x23c2091df2547f77185ba592b06ee2e897b0c2a70f968521a6a24fc5bfc1b1e",
			},
			{
				name: "after 0.14.1 + UseBlake2sHash false",
				txnOptions: &account.TxnOptions{
					UseBlake2sHash: &[]bool{false}[0],
				},
				starknetVersion:          "0.14.1",
				expectedCompileClassHash: "0x6ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17",
			},
		}

		for _, test := range testcases {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockRPCProvider := rpcv10mock.NewMockRPCProvider(ctrl)

				ks, pub, _ := account.GetRandomKeys()
				// called when instantiating the account
				mockRPCProvider.EXPECT().ChainID(gomock.Any()).Return("SN_SEPOLIA", nil).Times(1)
				acnt, err := account.NewAccount(
					mockRPCProvider,
					internalUtils.DeadBeef,
					pub.String(),
					ks,
					account.CairoV2,
				)
				require.NoError(t, err)

				// called in the BuildAndSendDeclareTxn method
				mockRPCProvider.EXPECT().
					Nonce(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(new(felt.Felt).SetUint64(1), nil).
					Times(1)
				mockRPCProvider.EXPECT().
					BlockWithTxs(t.Context(), rpc.WithBlockTag(rpc.BlockTagLatest)).
					Return(&rpc.BlockWithTxsOutput{Block: &rpc.Block{}}, nil).Times(1)
					mockRPCProvider.EXPECT().
					EstimateFee(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]rpc.FeeEstimation{
						{
							FeeEstimationCommon: rpc.FeeEstimationCommon{
								L1GasPrice:        new(felt.Felt).SetUint64(10),
								L1GasConsumed:     new(felt.Felt).SetUint64(100),
								L1DataGasPrice:    new(felt.Felt).SetUint64(5),
								L1DataGasConsumed: new(felt.Felt).SetUint64(50),
								L2GasPrice:        new(felt.Felt).SetUint64(3),
								L2GasConsumed:     new(felt.Felt).SetUint64(200),
							},
						},
					}, nil).
					Times(1)

				var compiledClassHash *felt.Felt

				if test.txnOptions == nil {
					// if txnOptions is nil, the code should call the BlockWithTxHashes method to get the
					// Starknet version and decide whether to use the Blake2s hash function
					mockRPCProvider.EXPECT().
						BlockWithTxHashes(gomock.Any(), rpc.WithBlockTag(rpc.BlockTagLatest)).
						Return(&rpc.BlockWithTxHashesOutput{
							Block: &rpc.BlockTxHashes{
								BlockHeader: rpc.BlockHeader{StarknetVersion: test.starknetVersion},
							},
						}, nil).Times(1)
				}

				mockRPCProvider.EXPECT().
					AddDeclareTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(
						func(_, txn any) (rpc.AddDeclareTransactionResponse, error) {
							declareTxn, ok := txn.(*rpc.BroadcastDeclareTxnV3)
							require.True(t, ok)

							compiledClassHash = declareTxn.CompiledClassHash

							return rpc.AddDeclareTransactionResponse{}, nil
						},
					).Times(1)

				_, err = acnt.BuildAndSendDeclareTxn(
					t.Context(),
					&casmClass,
					&class,
					test.txnOptions,
				)
				require.NoError(t, err)

				assert.Equal(t, test.expectedCompileClassHash, compiledClassHash.String())
			})
		}
	})
}

// BuildAndEstimateDeployAccountTxn is a test function that tests the BuildAndSendDeployAccount method.
//
// This function tests the BuildAndSendDeployAccount method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
func TestBuildAndEstimateDeployAccountTxn(t *testing.T) {
	// TODO: implement devnet support
	tests.RunTestOn(t, tests.TestnetEnv)

	provider, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	// we need this account to fund the new account with STRK tokens, in order to deploy it
	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Get random keys to create the new account
	ks, pub, _ := account.GetRandomKeys()

	// Set up the account passing random values to 'accountAddress' and 'cairoVersion' variables,
	// as for this case we only need the 'ks' to sign the deploy transaction.
	tempAcc, err := account.NewAccount(provider, pub, pub.String(), ks, account.CairoV2)
	if err != nil {
		panic(err)
	}

	// OpenZeppelin Account Class Hash in Sepolia
	classHash := internalUtils.TestHexToFelt(
		t,
		"0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f",
	)

	// Build, estimate the fee and precompute the address of the new account
	deployAccTxn, precomputedAddress, err := tempAcc.BuildAndEstimateDeployAccountTxn(
		t.Context(),
		new(felt.Felt).SetUint64(uint64(time.Now().UnixNano())), // random salt
		classHash,
		[]*felt.Felt{pub},
		nil,
	)
	require.NoError(t, err, "Error building and estimating deploy account txn")
	require.NotNil(t, deployAccTxn)
	require.NotNil(t, precomputedAddress)
	t.Logf("Precomputed address: %s", precomputedAddress)

	// multiplier is 1, since the BuildAndEstimateDeployAccountTxn method already
	// multiplies the fee by 1.5
	overallFee, err := utils.ResBoundsMapToOverallFee(
		deployAccTxn.ResourceBounds,
		1,
		deployAccTxn.Tip,
	)
	require.NoError(t, err, "Error converting resource bounds to overall fee")

	// Fund the new account with STRK tokens
	transferSTRKAndWaitConfirmation(t, acc, overallFee, precomputedAddress)

	// There's something like a bug, a delay in the sequencer in fetching the balance, so let's wait for 5 seconds
	// to be sure that the balance is fetched
	time.Sleep(5 * time.Second)

	// Deploy the new account
	resp, err := provider.AddDeployAccountTransaction(t.Context(), deployAccTxn)
	require.NoError(t, err, "Error deploying new account")

	require.NotNil(t, resp.Hash)
	t.Logf("Deploy account transaction hash: %s", resp.Hash)
	require.NotNil(t, resp.ContractAddress)

	txReceipt, err := acc.WaitForTransactionReceipt(
		t.Context(),
		resp.Hash,
		500*time.Millisecond,
	)
	require.NoError(t, err, "Error waiting for deploy account transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)

	// testing the default tip estimation feature
	txn, err := acc.Provider.TransactionByHash(t.Context(), resp.Hash)
	require.NoError(t, err, "Error getting transaction by hash")
	require.NotNil(t, txn)
	assert.NotEqual(t, "0x0", txn.Transaction.(rpc.DeployAccountTxnV3).Tip)
}

// a helper function that transfers STRK tokens to a given address and waits for confirmation,
// used to fund the new account with STRK tokens in the TestBuildAndEstimateDeployAccountTxn test
func transferSTRKAndWaitConfirmation(
	t *testing.T,
	acc *account.Account,
	amount, recipient *felt.Felt,
) {
	t.Helper()
	// Build and send invoke txn
	u256Amount, err := internalUtils.HexToU256Felt(amount.String())
	require.NoError(t, err, "Error converting amount to u256")
	resp, err := acc.BuildAndSendInvokeTxn(t.Context(), []rpc.InvokeFunctionCall{
		{
			// STRK contract address in Sepolia
			ContractAddress: internalUtils.TestHexToFelt(
				t,
				"0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D",
			),
			FunctionName: "transfer",
			CallData:     append([]*felt.Felt{recipient}, u256Amount...),
		},
	}, nil)
	require.NoError(t, err, "Error transferring STRK tokens")

	// check the transaction hash
	require.NotNil(t, resp.Hash)
	t.Logf("Transfer transaction hash: %s", resp.Hash)

	txReceipt, err := acc.WaitForTransactionReceipt(
		t.Context(),
		resp.Hash,
		500*time.Millisecond,
	)
	require.NoError(t, err, "Error waiting for transfer transaction receipt")

	err = waitForTransactionStatus(
		t.Context(),
		acc.Provider,
		resp.Hash,
		rpc.TxnStatusAcceptedOnL2,
		500*time.Millisecond,
	)
	require.NoError(t, err, "Error waiting for transfer transaction status")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
}

// TODO: make it an exported utility function
func waitForTransactionStatus(
	ctx context.Context,
	provider rpc.RPCProvider,
	transactionHash *felt.Felt,
	txnStatus rpc.TxnStatus,
	pollInterval time.Duration,
) error {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return rpcerr.Err(rpcerr.InternalError, rpc.StringErrData(ctx.Err().Error()))
		case <-t.C:
			returnedTxnStatus, err := provider.TransactionStatus(ctx, transactionHash)
			if err != nil {
				rpcErr := err.(*rpc.RPCError)
				if rpcErr.Code == rpc.ErrHashNotFound.Code &&
					rpcErr.Message == rpc.ErrHashNotFound.Message {
					continue
				} else {
					return err
				}
			}

			if returnedTxnStatus.FinalityStatus == txnStatus {
				return nil
			}
		}
	}
}

// TestBuildAndSendMethodsWithQueryBit is a test function that tests the BuildAndSendDeclareTxn, BuildAndSendInvokeTxn
// and BuildAndEstimateDeployAccountTxn methods with a query bit version.
//
// The tests will test the methods when called with the 'hasQueryBitVersion' parameter set to true.
// It'll check if the txn version when estimating has the query bit, and if the txn version when sending does NOT have it.
func TestBuildAndSendMethodsWithQueryBit(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.DevnetEnv)

	// Class
	class := internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](
		t,
		"./testData/contracts_v2_HelloStarknet.sierra.json",
	)

	// Casm Class
	casmClass := internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](
		t,
		"./testData/contracts_v2_HelloStarknet.casm.json",
	)

	t.Run("on mock", func(t *testing.T) {
		tests.RunTestOn(t, tests.MockEnv)

		ctrl := gomock.NewController(t)
		mockRPCProvider := rpcv10mock.NewMockRPCProvider(ctrl)

		mockRPCProvider.EXPECT().
			Nonce(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(new(felt.Felt).SetUint64(1), nil).
			Times(2)

		ks, pub, _ := account.GetRandomKeys()

		// called when instantiating the account
		mockRPCProvider.EXPECT().ChainID(gomock.Any()).Return("SN_SEPOLIA", nil).Times(1)

		acnt, err := account.NewAccount(
			mockRPCProvider,
			internalUtils.DeadBeef,
			pub.String(),
			ks,
			account.CairoV2,
		)
		require.NoError(t, err)

		// setting the expected behaviour for each call to EstimateFee,
		// asserting if the passed txn has the query bit version
		mockRPCProvider.EXPECT().
			EstimateFee(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(_, request, _, _ any) ([]rpc.FeeEstimation, error) {
					reqArr, ok := request.([]rpc.BroadcastTxn)
					require.True(t, ok)
					txn, ok := reqArr[0].(rpc.Transaction)
					require.True(t, ok)

					// assert that the transaction being estimated has the query bit version
					assert.Equal(t, txn.GetVersion(), rpc.TransactionV3WithQueryBit)

					return []rpc.FeeEstimation{
						{
							FeeEstimationCommon: rpc.FeeEstimationCommon{
								L1GasPrice:        new(felt.Felt).SetUint64(10),
								L1GasConsumed:     new(felt.Felt).SetUint64(100),
								L1DataGasPrice:    new(felt.Felt).SetUint64(5),
								L1DataGasConsumed: new(felt.Felt).SetUint64(50),
								L2GasPrice:        new(felt.Felt).SetUint64(3),
								L2GasConsumed:     new(felt.Felt).SetUint64(200),
							},
						},
					}, nil
				},
			).
			Times(3)

		// modified txn with a 10000 tip
		fakeTxn := internalUtils.TestUnmarshalJSONFileToType[rpc.InvokeTxnV3](
			t,
			"./testData/fakeInvokeTxn.json",
		)
		// called when estimating the tip
		mockRPCProvider.EXPECT().
			BlockWithTxs(t.Context(), rpc.WithBlockTag(rpc.BlockTagLatest)).
			Return(&rpc.BlockWithTxsOutput{
				Block: &rpc.Block{
					BlockHeader: rpc.BlockHeader{},
					Status:      rpc.BlockStatusAcceptedOnL2,
					Transactions: []rpc.BlockTransaction{
						{
							Hash:        internalUtils.DeadBeef,
							Transaction: fakeTxn,
						},
						{
							Hash:        internalUtils.DeadBeef,
							Transaction: fakeTxn,
						},
					},
				},
			}, nil).Times(3)

		t.Run("BuildAndSendInvokeTxn", func(t *testing.T) {
			mockRPCProvider.EXPECT().AddInvokeTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_, txn any) (rpc.AddInvokeTransactionResponse, error) {
					bcTxn, ok := txn.(*rpc.BroadcastInvokeTxnV3)
					require.True(t, ok)

					// assert that the transaction being added does NOT have the query bit version
					assert.Equal(t, bcTxn.GetVersion(), rpc.TransactionV3)

					return rpc.AddInvokeTransactionResponse{}, nil
				},
			)

			_, err = acnt.BuildAndSendInvokeTxn(t.Context(), []rpc.InvokeFunctionCall{
				{
					ContractAddress: internalUtils.DeadBeef,
					FunctionName:    "transfer",
				},
			}, &account.TxnOptions{
				UseQueryBit: true,
			})
			require.NoError(t, err)
		})

		t.Run("BuildAndSendDeclareTxn", func(t *testing.T) {
			mockRPCProvider.EXPECT().
				BlockWithTxHashes(gomock.Any(), rpc.WithBlockTag(rpc.BlockTagLatest)).
				Return(&rpc.BlockWithTxHashesOutput{
					Block: &rpc.BlockTxHashes{
						BlockHeader: rpc.BlockHeader{StarknetVersion: "0.14.1"},
					},
				}, nil).Times(1)
			mockRPCProvider.EXPECT().AddDeclareTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_, txn any) (rpc.AddDeclareTransactionResponse, error) {
					bcTxn, ok := txn.(*rpc.BroadcastDeclareTxnV3)
					require.True(t, ok)

					// assert that the transaction being added does NOT have the query bit version
					assert.Equal(t, bcTxn.GetVersion(), rpc.TransactionV3)

					return rpc.AddDeclareTransactionResponse{}, nil
				},
			)

			_, err = acnt.BuildAndSendDeclareTxn(
				t.Context(),
				&casmClass,
				&class,
				&account.TxnOptions{
					UseQueryBit: true,
				},
			)
			require.NoError(t, err)
		})

		t.Run("TestBuildAndEstimateDeployAccountTxn", func(t *testing.T) {
			txn, _, err := acnt.BuildAndEstimateDeployAccountTxn(
				t.Context(),
				pub,
				internalUtils.DeadBeef,
				[]*felt.Felt{pub},
				&account.TxnOptions{
					UseQueryBit: true,
				},
			)
			require.NoError(t, err)

			// assert the returned transaction does NOT have the query bit version
			assert.Equal(t, txn.Version, rpc.TransactionV3)
		})
	})

	t.Run("on devnet", func(t *testing.T) {
		tests.RunTestOn(t, tests.DevnetEnv)

		client, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
		require.NoError(t, err, "Error in rpc.NewProvider")

		_, acnts, err := newDevnet(t, tConfig.providerURL)
		require.NoError(t, err, "Error setting up Devnet")

		acnt := newDevnetAccount(t, client, acnts[0], account.CairoV2)

		t.Run("BuildAndSendDeclareTxn", func(t *testing.T) {
			resp, err := acnt.BuildAndSendDeclareTxn(
				t.Context(),
				&casmClass,
				&class,
				&account.TxnOptions{
					UseQueryBit: true,
				},
			)
			require.NoError(t, err)

			txn, err := client.TransactionByHash(t.Context(), resp.Hash)
			require.NoError(t, err)

			// assert the returned transaction does NOT have the query bit version
			assert.Equal(t, txn.GetVersion(), rpc.TransactionV3)
		})

		t.Run("BuildAndSendInvokeTxn", func(t *testing.T) {
			u256Amount, err := internalUtils.HexToU256Felt("0x10000")
			acntaddr2 := internalUtils.TestHexToFelt(t, acnts[1].Address)

			require.NoError(t, err, "Error converting amount to u256")

			resp, err := acnt.BuildAndSendInvokeTxn(t.Context(), []rpc.InvokeFunctionCall{
				{
					// STRK contract address in Sepolia
					ContractAddress: internalUtils.TestHexToFelt(
						t,
						"0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D",
					),
					FunctionName: "transfer",
					CallData:     append([]*felt.Felt{acntaddr2}, u256Amount...),
				},
			}, &account.TxnOptions{
				UseQueryBit: true,
			})
			require.NoError(t, err)

			txn, err := client.TransactionByHash(t.Context(), resp.Hash)
			require.NoError(t, err)

			// assert the returned transaction does NOT have the query bit version
			assert.Equal(t, txn.GetVersion(), rpc.TransactionV3)
		})

		t.Run("BuildAndEstimateDeployAccountTxn", func(t *testing.T) {
			// Get random keys to create the new account
			ks, pub, _ := account.GetRandomKeys()
			tempAcc, err := account.NewAccount(client, pub, pub.String(), ks, account.CairoV2)
			require.NoError(t, err)

			classHash := internalUtils.TestHexToFelt(
				t,
				"0x05b4b537eaa2399e3aa99c4e2e0208ebd6c71bc1467938cd52c798c601e43564",
			) // preDeployed OZ account classhash in devnet
			// Build and send deploy account txn
			txn, _, err := tempAcc.BuildAndEstimateDeployAccountTxn(
				t.Context(),
				pub,
				classHash,
				[]*felt.Felt{pub},
				&account.TxnOptions{
					UseQueryBit: true,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, txn)

			// assert the returned transaction does NOT have the query bit version
			assert.Equal(t, txn.Version, rpc.TransactionV3)
		})
	})
}

// TestAddInvoke is a test function that verifies the behaviour of the AddInvokeTransaction method.
//
// This function tests the AddInvokeTransaction method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
//
// Parameters:
//   - t: The testing.T instance for running the test
//
// Returns:
//
//	none
func TestSendInvokeTxn(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	type testSetType struct {
		ExpectedErr          *rpc.RPCError
		CairoContractVersion account.CairoVersion
		SetKS                bool
		AccountAddress       *felt.Felt
		PubKey               *felt.Felt
		PrivKey              *felt.Felt
		InvokeTx             rpc.BroadcastInvokeTxnV3
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.TestnetEnv: {
			{
				// https://sepolia.voyager.online/tx/0x51d224a96a8a07e07e31754e20e713c9cccdcfd7f61105700b1a72b2715ed9f
				ExpectedErr:          rpc.ErrInvalidTransactionNonce,
				CairoContractVersion: account.CairoV2,
				AccountAddress:       internalUtils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				SetKS:                true,
				PubKey:               internalUtils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e"),
				PrivKey:              internalUtils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"), //
				InvokeTx: rpc.BroadcastInvokeTxnV3{
					Nonce:   internalUtils.TestHexToFelt(t, "0x196bfe"),
					Type:    rpc.TransactionTypeInvoke,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x7d975abd8cb41ad812a57f509b5ce8c696a56dd1d133baeb8a8c6804e1f24ac"),
						internalUtils.TestHexToFelt(t, "0x82285115ef2f99fa7e69ee04d11054c94b0686e62873bdb4b3377efc0830b4"),
					},
					ResourceBounds: &rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x11170",
							MaxPricePerUnit: "0x8d79883d20000",
						},
						L1DataGas: rpc.ResourceBounds{
							MaxAmount:       "0x2710",
							MaxPricePerUnit: "0x62448724953354",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x5f5e100",
							MaxPricePerUnit: "0xba43b7400",
						},
					},
					Tip:                   "0x5f5e100",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					SenderAddress:         internalUtils.TestHexToFelt(t, "0x4f4e29add19afa12c868ba1f4439099f225403ff9a71fe667eebb50e13518d3"),
					Calldata: internalUtils.TestHexArrToFelt(t, []string{
						"0x2",
						"0x3eaf27245e5a10286542e75c216d17432dd077984c86d37944ba7f5002d10d3",
						"0x382be990ca34815134e64a9ac28f41a907c62e5ad10547f97174362ab94dc89",
						"0x0",
						"0x2a730fc5366a8932645ada40338487d5c272294d70a43dc2d53f03534f418ea",
						"0x3d3da80997f8be5d16e9ae7ee6a4b5f7191d60765a1a6c219ab74269c85cf97",
						"0x0",
					}),
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
				},
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		client, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
		require.NoError(t, err, "Error in rpc.NewProvider")

		// Set up ks
		ks := account.NewMemKeystore()
		if test.SetKS {
			fakePrivKeyBI, ok := new(big.Int).SetString(test.PrivKey.String(), 0)
			require.True(t, ok)
			ks.Put(test.PubKey.String(), fakePrivKeyBI)
		}

		acnt, err := account.NewAccount(
			client,
			test.AccountAddress,
			test.PubKey.String(),
			ks,
			account.CairoV2,
		)
		require.NoError(t, err)

		err = acnt.SignInvokeTransaction(t.Context(), &test.InvokeTx)
		require.NoError(t, err)

		resp, err := acnt.SendTransaction(t.Context(), test.InvokeTx)
		if err != nil {
			rpcErr := err.(*rpc.RPCError)
			require.Equal(
				t,
				test.ExpectedErr.Code,
				rpcErr.Code,
				"AddInvokeTransaction returned an unexpected error",
			)
			require.Zero(t, resp)
		}
	}
}

// TestAddDeclareTxn is a test function that verifies the behaviour of the AddDeclareTransaction method.
//
// This function tests the AddDeclareTransaction method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
//
// Parameters:
//   - t: The testing.T instance for running the test
//
// Returns:
//
//	none
func TestSendDeclareTxn(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	expectedTxHash := internalUtils.TestHexToFelt(
		t,
		"0x1c3df33f06f0da7f5df72bbc02fb8caf33e91bdd2433305dd007c6cd6acc6d0",
	)
	expectedClassHash := internalUtils.TestHexToFelt(
		t,
		"0x06ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17",
	)

	AccountAddress := internalUtils.TestHexToFelt(
		t,
		"0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9",
	)
	PubKey := internalUtils.TestHexToFelt(
		t,
		"0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e",
	)
	PrivKey := internalUtils.TestHexToFelt(
		t,
		"0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9",
	)

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(PrivKey.String(), 0)
	require.True(t, ok)
	ks.Put(PubKey.String(), fakePrivKeyBI)

	client, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acnt, err := account.NewAccount(client, AccountAddress, PubKey.String(), ks, account.CairoV0)
	require.NoError(t, err)

	// Class
	class := internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](
		t,
		"./testData/contracts_v2_HelloStarknet.sierra.json",
	)

	// Compiled Class Hash
	casmClass := internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](
		t,
		"./testData/contracts_v2_HelloStarknet.casm.json",
	)
	compClassHash, err := hash.CompiledClassHashV2(&casmClass)
	require.NoError(t, err)

	broadcastTx := rpc.BroadcastDeclareTxnV3{
		Type:              rpc.TransactionTypeDeclare,
		SenderAddress:     AccountAddress,
		CompiledClassHash: compClassHash,
		Version:           rpc.TransactionV3,
		Signature: []*felt.Felt{
			internalUtils.TestHexToFelt(
				t,
				"0x74a20e84469ecf7bfaa7eb82a803621357b695af5ac6f857c0615c7e9fa94e3",
			),
			internalUtils.TestHexToFelt(
				t,
				"0x3a79c411c05fc60fe6da68bd4a1cc57745a7e1e6cfa95dd7c3466fae384cfc3",
			),
		},
		Nonce:         internalUtils.TestHexToFelt(t, "0xe"),
		ContractClass: &class,
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x1597b3274d88",
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x210",
				MaxPricePerUnit: "0x997c",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x1115cde0",
				MaxPricePerUnit: "0x11920d1317",
			},
		},
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	err = acnt.SignDeclareTransaction(t.Context(), &broadcastTx)
	require.NoError(t, err)

	resp, err := acnt.SendTransaction(t.Context(), broadcastTx)

	if err != nil {
		rpcErr := err.(*rpc.RPCError)
		require.Equal(
			t,
			rpc.ErrInvalidTransactionNonce.Code,
			rpcErr.Code,
			"AddDeclareTransaction error not what expected",
		)
	} else {
		require.Equal(t, expectedTxHash.String(), resp.Hash.String(), "AddDeclareTransaction TxHash not what expected")
		require.Equal(t, expectedClassHash.String(), resp.ClassHash.String(), "AddDeclareTransaction ClassHash not what expected")
	}
}

// TestAddDeployAccountDevnet tests the functionality of adding a deploy account in the devnet environment.
//
// The test checks if the environment is set to "devnet" and skips the test if not. It then initialises a new RPC client
// and provider using the tConfig.base URL. After that, it sets up a devnet environment and creates a fake user account. The
// fake user's address and public key are converted to the appropriate format. The test also sets up a memory keystore
// and puts the fake user's public key and private key in it. Then, it creates a new account using the provider, fake
// user's address, public key, and keystore. Next, it converts a class hash to the appropriate format. The test
// constructs a deploy account transaction and precomputes the address. It then signs the transaction and mints coins to
// the precomputed address. Finally, it adds the deploy account transaction and verifies that no errors occurred and the
// response is not nil.
//
// Parameters:
//   - t: is the testing framework
//
// Returns:
//
//	none
func TestSendDeployAccountDevnet(t *testing.T) {
	tests.RunTestOn(t, tests.DevnetEnv)

	client, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	devnetClient, acnts, err := newDevnet(t, tConfig.providerURL)
	require.NoError(t, err, "Error setting up Devnet")

	fakeUser := acnts[0]
	fakeUserPub := internalUtils.TestHexToFelt(t, fakeUser.PublicKey)
	acnt := newDevnetAccount(t, client, fakeUser, account.CairoV2)

	classHash := internalUtils.TestHexToFelt(
		t,
		"0x05b4b537eaa2399e3aa99c4e2e0208ebd6c71bc1467938cd52c798c601e43564",
	) // preDeployed classhash
	require.NoError(t, err)

	tx := rpc.DeployAccountTxnV3{
		Type:                rpc.TransactionTypeDeployAccount,
		Version:             rpc.TransactionV3,
		Signature:           []*felt.Felt{},
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		ContractAddressSalt: fakeUserPub,
		ConstructorCalldata: []*felt.Felt{fakeUserPub},
		ClassHash:           classHash,
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x997c",
				MaxPricePerUnit: "0x1597b3274d88",
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x2230",
				MaxPricePerUnit: "0x9924327c",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x15cde0",
				MaxPricePerUnit: "0x11920d1317",
			},
		},
		Tip:           "0x0",
		PayMasterData: []*felt.Felt{},
		NonceDataMode: rpc.DAModeL1,
		FeeMode:       rpc.DAModeL1,
	}

	precomputedAddress := account.PrecomputeAccountAddress(
		fakeUserPub,
		classHash,
		tx.ConstructorCalldata,
	)
	require.NoError(
		t,
		acnt.SignDeployAccountTransaction(t.Context(), &tx, precomputedAddress),
	)

	_, err = devnetClient.Mint(precomputedAddress, new(big.Int).SetUint64(10000000000000000000))
	require.NoError(t, err)

	resp, err := acnt.SendTransaction(t.Context(), tx)
	if err != nil {
		// TODO: remove this once devnet supports full v3 transaction type
		require.ErrorContains(t, err, "unsupported transaction type")

		return
	}
	require.Nil(t, err, "AddDeployAccountTransaction gave an Error")
	require.NotNil(t, resp, "AddDeployAccountTransaction resp not nil")
}

// TestWaitForTransactionReceiptMOCK is a unit test for the WaitForTransactionReceipt function.
//
// It tests the functionality of WaitForTransactionReceipt by mocking the RpcProvider and simulating different test scenarios.
// It creates a test set with different parameters and expectations, and iterates over the test set to run the test cases.
// For each test case, it sets up the necessary mocks, creates a context with a timeout, and calls the WaitForTransactionReceipt function.
// It then asserts the expected result against the actual result.
// The function uses the testify package for assertions and the gomock package for creating mocks.
//
// Parameters:
//   - t: The testing.T object for test assertions and logging
//
// Returns:
//
//	none
func TestWaitForTransactionReceiptMOCK(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	mockCtrl := gomock.NewController(t)
	mockRPCProvider := rpcv10mock.NewMockRPCProvider(mockCtrl)

	mockRPCProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)

	acnt, err := account.NewAccount(
		mockRPCProvider,
		&felt.Zero,
		"",
		account.NewMemKeystore(),
		account.CairoV0,
	)
	require.NoError(t, err, "error returned from account.NewAccount()")

	type testSetType struct {
		Timeout                      time.Duration
		ShouldCallTransactionReceipt bool
		Hash                         *felt.Felt
		ExpectedErr                  error
		ExpectedReceipt              *rpc.TransactionReceiptWithBlockInfo
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Timeout:                      time.Duration(1000),
				ShouldCallTransactionReceipt: true,
				Hash:                         new(felt.Felt).SetUint64(1),
				ExpectedReceipt:              nil,
				ExpectedErr:                  rpcerr.Err(rpcerr.InternalError, rpc.StringErrData("UnExpectedErr")),
			},
			{
				Timeout:                      time.Duration(1000),
				Hash:                         new(felt.Felt).SetUint64(2),
				ShouldCallTransactionReceipt: true,
				ExpectedReceipt: &rpc.TransactionReceiptWithBlockInfo{
					TransactionReceipt: rpc.TransactionReceipt{},
					BlockHash:          new(felt.Felt).SetUint64(2),
					BlockNumber:        2,
				},

				ExpectedErr: nil,
			},
			{
				Timeout:                      time.Duration(1),
				Hash:                         new(felt.Felt).SetUint64(3),
				ShouldCallTransactionReceipt: false,
				ExpectedReceipt:              nil,
				ExpectedErr:                  rpcerr.Err(rpcerr.InternalError, rpc.StringErrData(context.DeadlineExceeded.Error())),
			},
		},
	}[tests.TEST_ENV]

	var wg sync.WaitGroup
	for _, test := range testSet {
		wg.Go(
			func() {
				ctx, cancel := context.WithTimeout(t.Context(), test.Timeout*time.Second)
				defer cancel()
				if test.ShouldCallTransactionReceipt {
					mockRPCProvider.EXPECT().
						TransactionReceipt(ctx, test.Hash).
						Return(test.ExpectedReceipt, test.ExpectedErr)
				}
				resp, err := acnt.WaitForTransactionReceipt(ctx, test.Hash, 2*time.Second)

				if test.ExpectedErr != nil {
					require.Equal(t, test.ExpectedErr, err)
				} else {
					// check
					require.Equal(t, test.ExpectedReceipt.ExecutionStatus, (resp.TransactionReceipt).ExecutionStatus)
				}
			})
	}
	wg.Wait()
}

// TestWaitForTransactionReceipt is a test function that tests the WaitForTransactionReceipt method.
//
// It checks if the test environment is "devnet" and skips the test if it's not.
// It creates a new RPC client using the tConfig.base URL and "/rpc" endpoint.
// It creates a new RPC provider using the client.
// It creates a new account using the provider, a zero-value Felt object, the "pubkey" string, and a new memory keystore.
// It defines a testSet variable that contains an array of testSetType structs.
// Each testSetType struct contains a Timeout integer, a Hash object, an ExpectedErr error, and an ExpectedReceipt TransactionReceipt object.
// It retrieves the testSet based on the tests.TEST_ENV variable.
// It iterates over each test in the testSet.
// For each test, it creates a new context with a timeout based on the test's Timeout value.
// It calls the WaitForTransactionReceipt method on the account object, passing the context, the test's Hash value, and a 1-second timeout.
// If the test's ExpectedErr is not nil, it asserts that the returned error matches the test's ExpectedErr error.
// Otherwise, it asserts that the ExecutionStatus of the returned receipt matches the ExecutionStatus of the test's ExpectedReceipt.
// It then cleans up the test environment.
//
// Parameters:
//   - t: The testing.T instance for running the test
//
// Returns:
//
//	none
func TestWaitForTransactionReceipt(t *testing.T) {
	tests.RunTestOn(t, tests.DevnetEnv)

	client, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acnt, err := account.NewAccount(
		client,
		&felt.Zero,
		"pubkey",
		account.NewMemKeystore(),
		account.CairoV0,
	)
	require.NoError(t, err, "error returned from account.NewAccount()")

	type testSetType struct {
		Timeout         int
		Hash            *felt.Felt
		ExpectedErr     *rpc.RPCError
		ExpectedReceipt rpc.TransactionReceipt
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.DevnetEnv: {
			{
				Timeout:         3, // Should poll 3 times
				Hash:            new(felt.Felt).SetUint64(100),
				ExpectedReceipt: rpc.TransactionReceipt{},
				ExpectedErr:     rpcerr.Err(rpcerr.InternalError, rpc.StringErrData("context deadline exceeded")),
			},
		},
	}[tests.TEST_ENV]

	var wg sync.WaitGroup
	for _, test := range testSet {
		wg.Go(
			func() {
				ctx, cancel := context.WithTimeout(
					t.Context(),
					time.Duration(test.Timeout)*time.Second,
				)
				defer cancel()

				resp, err := acnt.WaitForTransactionReceipt(ctx, test.Hash, 1*time.Second)
				if test.ExpectedErr != nil {
					rpcErr, ok := err.(*rpc.RPCError)
					require.True(t, ok)
					require.Equal(t, test.ExpectedErr.Code, rpcErr.Code)
					require.Contains(
						t,
						rpcErr.Data.ErrorMessage(),
						test.ExpectedErr.Data.ErrorMessage(),
					) // sometimes the error message starts with "Post \"http://localhost:5050\":..."
				} else {
					require.Equal(t, test.ExpectedReceipt.ExecutionStatus, resp.ExecutionStatus)
				}
			},
		)
	}
	wg.Wait()
}

func TestDeployContractWithUDC(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	provider, err := rpc.NewProvider(t.Context(), tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	accnt, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	t.Run("UDCCairoV0, no constructor, udcOptions nil", func(t *testing.T) {
		classHash, _ := utils.HexToFelt(
			"0x0387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8",
		)

		resp, _, err := accnt.DeployContractWithUDC(t.Context(), classHash, nil, nil, nil)
		require.NoError(t, err, "DeployContractUDC failed")

		t.Logf("Transaction hash: %s", resp.Hash)

		txReceipt, err := accnt.WaitForTransactionReceipt(
			t.Context(),
			resp.Hash,
			500*time.Millisecond,
		)
		require.NoError(t, err, "Waiting for tx receipt failed")

		assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	})

	t.Run("error, UDCCairoV0, no constructor, all udcOptions set", func(t *testing.T) {
		classHash, _ := utils.HexToFelt(
			"0x0387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8",
		)

		_, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			nil,
			nil,
			&utils.UDCOptions{
				Salt:              internalUtils.DeadBeef,
				UDCVersion:        utils.UDCCairoV0,
				OriginIndependent: true,
			},
		)
		assert.ErrorContains(t, err, "contract already deployed")
	})

	t.Run("UDCCairoV2, no constructor, only UDCVersion set", func(t *testing.T) {
		classHash, _ := utils.HexToFelt(
			"0x0387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8",
		)

		resp, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			nil,
			nil,
			&utils.UDCOptions{
				UDCVersion: utils.UDCCairoV2,
			},
		)
		require.NoError(t, err, "DeployContractUDC failed")

		t.Logf("Transaction hash: %s", resp.Hash)

		txReceipt, err := accnt.WaitForTransactionReceipt(
			t.Context(),
			resp.Hash,
			500*time.Millisecond,
		)
		require.NoError(t, err, "Waiting for tx receipt failed")

		assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	})

	t.Run("error, UDCCairoV2, no constructor, all udcOptions set", func(t *testing.T) {
		classHash, _ := utils.HexToFelt(
			"0x0387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8",
		)

		_, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			nil,
			nil,
			&utils.UDCOptions{
				Salt:              internalUtils.DeadBeef,
				UDCVersion:        utils.UDCCairoV2,
				OriginIndependent: true,
			},
		)
		assert.ErrorContains(t, err, "contract already deployed")
	})

	// erc20 calldata:
	classHash, _ := utils.HexToFelt(
		"0x73d71c37e20c569186445d2c497d2195b4c0be9a255d72dbad86662fcc63ae6",
	)
	name, _ := utils.StringToByteArrFelt("My Test Token")
	symbol, _ := utils.StringToByteArrFelt("MTT")
	supply, _ := utils.HexToU256Felt("0x200000000000000000")
	recipient := accnt.Address
	owner := accnt.Address

	// https://docs.openzeppelin.com/contracts-cairo/1.0.0/api/erc20#ERC20Upgradeable-constructor
	constructorCalldata := make([]*felt.Felt, 0, 10)
	constructorCalldata = append(constructorCalldata, name...)
	constructorCalldata = append(constructorCalldata, symbol...)
	constructorCalldata = append(constructorCalldata, supply...)
	constructorCalldata = append(constructorCalldata, recipient, owner)

	t.Run("UDCCairoV0, with constructor - ERC20, udcOptions nil", func(t *testing.T) {
		resp, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			constructorCalldata,
			nil,
			nil,
		)
		require.NoError(t, err, "DeployContractUDC failed")

		t.Logf("Transaction hash: %s", resp.Hash)

		txReceipt, err := accnt.WaitForTransactionReceipt(
			t.Context(),
			resp.Hash,
			500*time.Millisecond,
		)
		require.NoError(t, err, "Waiting for tx receipt failed")

		assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	})

	t.Run("error, UDCCairoV0, with constructor - ERC20, all udcOptions set", func(t *testing.T) {
		_, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			constructorCalldata,
			nil,
			&utils.UDCOptions{
				Salt:              internalUtils.DeadBeef,
				UDCVersion:        utils.UDCCairoV0,
				OriginIndependent: true,
			},
		)
		assert.ErrorContains(t, err, "contract already deployed")
	})

	t.Run("UDCCairoV2, with constructor - ERC20, udcOptions nil", func(t *testing.T) {
		resp, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			constructorCalldata,
			nil,
			&utils.UDCOptions{
				UDCVersion: utils.UDCCairoV2,
			},
		)
		require.NoError(t, err, "DeployContractUDC failed")

		t.Logf("Transaction hash: %s", resp.Hash)

		txReceipt, err := accnt.WaitForTransactionReceipt(
			t.Context(),
			resp.Hash,
			500*time.Millisecond,
		)
		require.NoError(t, err, "Waiting for tx receipt failed")

		assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	})

	t.Run("error, UDCCairoV2, with constructor - ERC20, all udcOptions set", func(t *testing.T) {
		_, _, err := accnt.DeployContractWithUDC(
			t.Context(),
			classHash,
			constructorCalldata,
			nil,
			&utils.UDCOptions{
				Salt:              internalUtils.DeadBeef,
				UDCVersion:        utils.UDCCairoV2,
				OriginIndependent: true,
			},
		)
		assert.ErrorContains(t, err, "contract already deployed")
	})
}

// TODO: add more tests for the BuildAnd* functions, testing each of them with different TxnOption's
