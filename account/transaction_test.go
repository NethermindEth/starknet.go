package account_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/mocks"
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
	testSet := map[string]bool{
		"testnet": true,
		"devnet":  false, // TODO:change to true once devnet supports full v3 transaction type, and adapt the code to use it
	}[testEnv]

	if !testSet {
		t.Skip("test environment not supported")
	}

	provider, err := rpc.NewProvider(tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Build and send invoke txn
	resp, err := acc.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{
		{
			// same ERC20 contract as in examples/simpleInvoke
			ContractAddress: internalUtils.TestHexToFelt(t, "0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54"),
			FunctionName:    "mint",
			CallData:        []*felt.Felt{new(felt.Felt).SetUint64(10000), &felt.Zero},
		},
	}, 1.5, false)
	require.NoError(t, err, "Error building and sending invoke txn")

	// check the transaction hash
	require.NotNil(t, resp.Hash)
	t.Logf("Invoke transaction hash: %s", resp.Hash)

	txReceipt, err := acc.WaitForTransactionReceipt(context.Background(), resp.Hash, 1*time.Second)
	require.NoError(t, err, "Error waiting for invoke transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	assert.Equal(t, rpc.TxnFinalityStatusAcceptedOnL2, txReceipt.FinalityStatus)
}

// TestBuildAndSendDeclareTxn is a test function that tests the BuildAndSendDeclareTxn method.
//
// This function tests the BuildAndSendDeclareTxn method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
func TestBuildAndSendDeclareTxn(t *testing.T) {
	testSet := map[string]bool{
		"testnet": true,
		"devnet":  false, // TODO:change to true once devnet supports full v3 transaction type, and adapt the code to use it
	}[testEnv]

	if !testSet {
		t.Skip("test environment not supported")
	}

	provider, err := rpc.NewProvider(tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Class
	class := *internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](t, "./tests/contracts_v2_HelloStarknet.sierra.json", "")

	// Casm Class
	casmClass := *internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](t, "./tests/contracts_v2_HelloStarknet.casm.json", "")

	// Build and send declare txn
	resp, err := acc.BuildAndSendDeclareTxn(context.Background(), &casmClass, &class, 1.5, false)
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

	txReceipt, err := acc.WaitForTransactionReceipt(context.Background(), resp.Hash, 1*time.Second)
	require.NoError(t, err, "Error waiting for declare transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	assert.Equal(t, rpc.TxnFinalityStatusAcceptedOnL2, txReceipt.FinalityStatus)
}

// BuildAndEstimateDeployAccountTxn is a test function that tests the BuildAndSendDeployAccount method.
//
// This function tests the BuildAndSendDeployAccount method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
func TestBuildAndEstimateDeployAccountTxn(t *testing.T) {
	testSet := map[string]bool{
		"testnet": true,
		"devnet":  false, // TODO:change to true once devnet supports full v3 transaction type, and adapt the code to use it
	}[testEnv]

	if !testSet {
		t.Skip("test environment not supported")
	}

	provider, err := rpc.NewProvider(tConfig.providerURL)
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
	classHash := internalUtils.TestHexToFelt(t, "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")

	// Build, estimate the fee and precompute the address of the new account
	deployAccTxn, precomputedAddress, err := tempAcc.BuildAndEstimateDeployAccountTxn(
		context.Background(),
		new(felt.Felt).SetUint64(uint64(time.Now().UnixNano())), // random salt
		classHash,
		[]*felt.Felt{pub},
		1.5,
		false,
	)
	require.NoError(t, err, "Error building and estimating deploy account txn")
	require.NotNil(t, deployAccTxn)
	require.NotNil(t, precomputedAddress)
	t.Logf("Precomputed address: %s", precomputedAddress)

	overallFee, err := utils.ResBoundsMapToOverallFee(deployAccTxn.ResourceBounds, 1)
	require.NoError(t, err, "Error converting resource bounds to overall fee")

	// Fund the new account with STRK tokens
	transferSTRKAndWaitConfirmation(t, acc, overallFee, precomputedAddress)

	// Deploy the new account
	resp, err := provider.AddDeployAccountTransaction(context.Background(), deployAccTxn)
	require.NoError(t, err, "Error deploying new account")

	require.NotNil(t, resp.Hash)
	t.Logf("Deploy account transaction hash: %s", resp.Hash)
	require.NotNil(t, resp.ContractAddress)

	txReceipt, err := acc.WaitForTransactionReceipt(context.Background(), resp.Hash, 1*time.Second)
	require.NoError(t, err, "Error waiting for deploy account transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	assert.Equal(t, rpc.TxnFinalityStatusAcceptedOnL2, txReceipt.FinalityStatus)
}

// a helper function that transfers STRK tokens to a given address and waits for confirmation,
// used to fund the new account with STRK tokens in the TestBuildAndEstimateDeployAccountTxn test
func transferSTRKAndWaitConfirmation(t *testing.T, acc *account.Account, amount, recipient *felt.Felt) {
	t.Helper()
	// Build and send invoke txn
	u256Amount, err := internalUtils.HexToU256Felt(amount.String())
	require.NoError(t, err, "Error converting amount to u256")
	resp, err := acc.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{
		{
			// STRK contract address in Sepolia
			ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
			FunctionName:    "transfer",
			CallData:        append([]*felt.Felt{recipient}, u256Amount...),
		},
	}, 1.5, false)
	require.NoError(t, err, "Error transferring STRK tokens")

	// check the transaction hash
	require.NotNil(t, resp.Hash)
	t.Logf("Transfer transaction hash: %s", resp.Hash)

	txReceipt, err := acc.WaitForTransactionReceipt(context.Background(), resp.Hash, 1*time.Second)
	require.NoError(t, err, "Error waiting for transfer transaction receipt")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txReceipt.ExecutionStatus)
	assert.Equal(t, rpc.TxnFinalityStatusAcceptedOnL2, txReceipt.FinalityStatus)
}

// TestBuildAndSendMethodsWithQueryBit is a test function that tests the BuildAndSendDeclareTxn, BuildAndSendInvokeTxn
// and BuildAndEstimateDeployAccountTxn methods with a query bit version.
//
// The tests will test the methods when called with the 'hasQueryBitVersion' parameter set to true.
// It'll check if the txn version when estimating has the query bit, and if the txn version when sending does NOT have it.
func TestBuildAndSendMethodsWithQueryBit(t *testing.T) {
	// Class
	class := *internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](t, "./tests/contracts_v2_HelloStarknet.sierra.json", "")

	// Casm Class
	casmClass := *internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](t, "./tests/contracts_v2_HelloStarknet.casm.json", "")

	t.Run("on mock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRpcProvider := mocks.NewMockRpcProvider(ctrl)

		mockRpcProvider.EXPECT().Nonce(gomock.Any(), gomock.Any(), gomock.Any()).Return(new(felt.Felt).SetUint64(1), nil).Times(2)

		ks, pub, _ := account.GetRandomKeys()
		// called when instantiating the account
		mockRpcProvider.EXPECT().ClassHashAt(gomock.Any(), gomock.Any(), gomock.Any()).Return(internalUtils.RANDOM_FELT, nil).Times(1)
		mockRpcProvider.EXPECT().ChainID(gomock.Any()).Return("SN_SEPOLIA", nil).Times(1)
		acnt, err := account.NewAccount(mockRpcProvider, internalUtils.RANDOM_FELT, pub.String(), ks, account.CairoV2)
		require.NoError(t, err)

		// setting the expected behaviour for each call to EstimateFee,
		// asserting if the passed txn has the query bit version
		mockRpcProvider.EXPECT().EstimateFee(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_, request, _, _ any) ([]rpc.FeeEstimation, error) {
				reqArr, ok := request.([]rpc.BroadcastTxn)
				require.True(t, ok)
				txn, ok := reqArr[0].(rpc.Transaction)
				require.True(t, ok)

				// assert that the transaction being estimated has the query bit version
				assert.Equal(t, txn.GetVersion(), rpc.TransactionV3WithQueryBit)

				return []rpc.FeeEstimation{
					{
						L1GasPrice:        new(felt.Felt).SetUint64(10),
						L1GasConsumed:     new(felt.Felt).SetUint64(100),
						L1DataGasPrice:    new(felt.Felt).SetUint64(5),
						L1DataGasConsumed: new(felt.Felt).SetUint64(50),
						L2GasPrice:        new(felt.Felt).SetUint64(3),
						L2GasConsumed:     new(felt.Felt).SetUint64(200),
					},
				}, nil
			},
		).Times(3)

		t.Run("BuildAndSendInvokeTxn", func(t *testing.T) {
			mockRpcProvider.EXPECT().AddInvokeTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_, txn any) (*rpc.AddInvokeTransactionResponse, error) {
					bcTxn, ok := txn.(*rpc.BroadcastInvokeTxnV3)
					require.True(t, ok)

					// assert that the transaction being added does NOT have the query bit version
					assert.Equal(t, bcTxn.GetVersion(), rpc.TransactionV3)

					return &rpc.AddInvokeTransactionResponse{}, nil
				},
			)

			_, err = acnt.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{
				{
					ContractAddress: internalUtils.RANDOM_FELT,
					FunctionName:    "transfer",
				},
			}, 1.5, true)
			require.NoError(t, err)
		})

		t.Run("BuildAndSendDeclareTxn", func(t *testing.T) {
			mockRpcProvider.EXPECT().AddDeclareTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_, txn any) (*rpc.AddDeclareTransactionResponse, error) {
					bcTxn, ok := txn.(*rpc.BroadcastDeclareTxnV3)
					require.True(t, ok)

					// assert that the transaction being added does NOT have the query bit version
					assert.Equal(t, bcTxn.GetVersion(), rpc.TransactionV3)

					return &rpc.AddDeclareTransactionResponse{}, nil
				},
			)

			_, err = acnt.BuildAndSendDeclareTxn(context.Background(), &casmClass, &class, 1.5, true)
			require.NoError(t, err)
		})

		t.Run("TestBuildAndEstimateDeployAccountTxn", func(t *testing.T) {
			txn, _, err := acnt.BuildAndEstimateDeployAccountTxn(
				context.Background(),
				pub,
				internalUtils.RANDOM_FELT,
				[]*felt.Felt{pub},
				1.5,
				true,
			)
			require.NoError(t, err)

			// assert the returned transaction does NOT have the query bit version
			assert.Equal(t, txn.Version, rpc.TransactionV3)
		})
	})

	t.Run("on devnet", func(t *testing.T) {
		if testEnv != "devnet" {
			t.Skip("Skipping test as it requires a devnet environment")
		}
		client, err := rpc.NewProvider(tConfig.providerURL)
		require.NoError(t, err, "Error in rpc.NewProvider")

		_, acnts, err := newDevnet(t, tConfig.providerURL)
		require.NoError(t, err, "Error setting up Devnet")

		acnt := newDevnetAccount(t, client, acnts[0], account.CairoV2)

		t.Run("BuildAndSendDeclareTxn", func(t *testing.T) {
			resp, err := acnt.BuildAndSendDeclareTxn(context.Background(), &casmClass, &class, 1.5, true)
			require.NoError(t, err)

			txn, err := client.TransactionByHash(context.Background(), resp.Hash)
			require.NoError(t, err)

			// assert the returned transaction does NOT have the query bit version
			assert.Equal(t, txn.GetVersion(), rpc.TransactionV3)
		})

		t.Run("BuildAndSendInvokeTxn", func(t *testing.T) {
			u256Amount, err := internalUtils.HexToU256Felt("0x10000")
			acntaddr2 := internalUtils.TestHexToFelt(t, acnts[1].Address)

			require.NoError(t, err, "Error converting amount to u256")

			resp, err := acnt.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{
				{
					// STRK contract address in Sepolia
					ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
					FunctionName:    "transfer",
					CallData:        append([]*felt.Felt{acntaddr2}, u256Amount...),
				},
			}, 1.5, true)
			require.NoError(t, err)

			txn, err := client.TransactionByHash(context.Background(), resp.Hash)
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
				"0x02b31e19e45c06f29234e06e2ee98a9966479ba3067f8785ed972794fdb0065c",
			) // preDeployed OZ account classhash in devnet
			// Build and send deploy account txn
			txn, _, err := tempAcc.BuildAndEstimateDeployAccountTxn(
				context.Background(),
				pub,
				classHash,
				[]*felt.Felt{pub},
				1.5,
				true,
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
	type testSetType struct {
		ExpectedErr          error
		CairoContractVersion account.CairoVersion
		SetKS                bool
		AccountAddress       *felt.Felt
		PubKey               *felt.Felt
		PrivKey              *felt.Felt
		InvokeTx             rpc.BroadcastInvokeTxnV3
	}
	testSet := map[string][]testSetType{
		"mock":   {},
		"devnet": {},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x7aac4792c8fd7578dd01b20ff04565f2e2ce6ea3c792c5e609a088704c1dd87
				ExpectedErr:          rpc.ErrDuplicateTx,
				CairoContractVersion: account.CairoV2,
				AccountAddress:       internalUtils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				SetKS:                true,
				PubKey:               internalUtils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e"),
				PrivKey:              internalUtils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"), //
				InvokeTx: rpc.BroadcastInvokeTxnV3{
					Nonce:   internalUtils.TestHexToFelt(t, "0xd"),
					Type:    rpc.TransactionType_Invoke,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x7bff07f1c2f6dc0eeaa9e622a0ee35f6e2e9855b39ed757236970a71b7c9e2e"),
						internalUtils.TestHexToFelt(t, "0x588b821ccb9f61ca217bfb0a580f889886742c2fd63526009eb401a9cf951e3"),
					},
					ResourceBounds: &rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x4305031628668",
						},
						L1DataGas: rpc.ResourceBounds{
							MaxAmount:       "0x210",
							MaxPricePerUnit: "0x948",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x15cde0",
							MaxPricePerUnit: "0x18955dc56",
						},
					},
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					SenderAddress:         internalUtils.TestHexToFelt(t, "0x1ae6fe02fcd9f61a3a8c30d68a8a7c470b0d7dd6f0ee685d5bbfa0d79406ff9"),
					Calldata: internalUtils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54",
						"0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
						"0x2",
						"0xffffffff",
						"0x0",
					}),
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		client, err := rpc.NewProvider(tConfig.providerURL)
		require.NoError(t, err, "Error in rpc.NewProvider")

		// Set up ks
		ks := account.NewMemKeystore()
		if test.SetKS {
			fakePrivKeyBI, ok := new(big.Int).SetString(test.PrivKey.String(), 0)
			require.True(t, ok)
			ks.Put(test.PubKey.String(), fakePrivKeyBI)
		}

		acnt, err := account.NewAccount(client, test.AccountAddress, test.PubKey.String(), ks, account.CairoV2)
		require.NoError(t, err)

		err = acnt.SignInvokeTransaction(context.Background(), &test.InvokeTx)
		require.NoError(t, err)

		resp, err := acnt.SendTransaction(context.Background(), test.InvokeTx)
		if err != nil {
			require.Equal(t, test.ExpectedErr.Error(), err.Error(), "AddInvokeTransaction returned an unexpected error")
			require.Nil(t, resp)
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
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}
	expectedTxHash := internalUtils.TestHexToFelt(t, "0x1c3df33f06f0da7f5df72bbc02fb8caf33e91bdd2433305dd007c6cd6acc6d0")
	expectedClassHash := internalUtils.TestHexToFelt(t, "0x06ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17")

	AccountAddress := internalUtils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9")
	PubKey := internalUtils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e")
	PrivKey := internalUtils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9")

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(PrivKey.String(), 0)
	require.True(t, ok)
	ks.Put(PubKey.String(), fakePrivKeyBI)

	client, err := rpc.NewProvider(tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acnt, err := account.NewAccount(client, AccountAddress, PubKey.String(), ks, account.CairoV0)
	require.NoError(t, err)

	// Class
	class := *internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](t, "./tests/contracts_v2_HelloStarknet.sierra.json", "")

	// Compiled Class Hash
	casmClass := *internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](t, "./tests/contracts_v2_HelloStarknet.casm.json", "")
	compClassHash, err := hash.CompiledClassHash(&casmClass)
	require.NoError(t, err)

	broadcastTx := rpc.BroadcastDeclareTxnV3{
		Type:              rpc.TransactionType_Declare,
		SenderAddress:     AccountAddress,
		CompiledClassHash: compClassHash,
		Version:           rpc.TransactionV3,
		Signature: []*felt.Felt{
			internalUtils.TestHexToFelt(t, "0x74a20e84469ecf7bfaa7eb82a803621357b695af5ac6f857c0615c7e9fa94e3"),
			internalUtils.TestHexToFelt(t, "0x3a79c411c05fc60fe6da68bd4a1cc57745a7e1e6cfa95dd7c3466fae384cfc3"),
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

	err = acnt.SignDeclareTransaction(context.Background(), &broadcastTx)
	require.NoError(t, err)

	resp, err := acnt.SendTransaction(context.Background(), broadcastTx)

	if err != nil {
		require.Equal(t, rpc.ErrDuplicateTx.Error(), err.Error(), "AddDeclareTransaction error not what expected")
	} else {
		require.Equal(t, expectedTxHash.String(), resp.Hash.String(), "AddDeclareTransaction TxHash not what expected")
		require.Equal(t, expectedClassHash.String(), resp.ClassHash.String(), "AddDeclareTransaction ClassHash not what expected")
	}
}

// TestAddDeployAccountDevnet tests the functionality of adding a deploy account in the devnet environment.
//
// The test checks if the environment is set to "devnet" and skips the test if not. It then initialises a new RPC client
// and provider using the base URL. After that, it sets up a devnet environment and creates a fake user account. The
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
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}
	client, err := rpc.NewProvider(tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	devnetClient, acnts, err := newDevnet(t, tConfig.providerURL)
	require.NoError(t, err, "Error setting up Devnet")

	fakeUser := acnts[0]
	fakeUserPub := internalUtils.TestHexToFelt(t, fakeUser.PublicKey)
	acnt := newDevnetAccount(t, client, fakeUser, account.CairoV2)

	classHash := internalUtils.TestHexToFelt(
		t,
		"0x02b31e19e45c06f29234e06e2ee98a9966479ba3067f8785ed972794fdb0065c",
	) // preDeployed classhash
	require.NoError(t, err)

	tx := rpc.DeployAccountTxnV3{
		Type:                rpc.TransactionType_DeployAccount,
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

	precomputedAddress := account.PrecomputeAccountAddress(fakeUserPub, classHash, tx.ConstructorCalldata)
	require.NoError(t, acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress))

	_, err = devnetClient.Mint(precomputedAddress, new(big.Int).SetUint64(10000000000000000000))
	require.NoError(t, err)

	resp, err := acnt.SendTransaction(context.Background(), tx)
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
	if testEnv != "mock" {
		t.Skip("Skipping test as it requires a mock environment")
	}
	mockCtrl := gomock.NewController(t)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)
	// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
	mockRpcProvider.EXPECT().ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).Return(internalUtils.RANDOM_FELT, nil)
	acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), account.CairoV0)
	require.NoError(t, err, "error returned from account.NewAccount()")

	type testSetType struct {
		Timeout                      time.Duration
		ShouldCallTransactionReceipt bool
		Hash                         *felt.Felt
		ExpectedErr                  error
		ExpectedReceipt              *rpc.TransactionReceiptWithBlockInfo
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				Timeout:                      time.Duration(1000),
				ShouldCallTransactionReceipt: true,
				Hash:                         new(felt.Felt).SetUint64(1),
				ExpectedReceipt:              nil,
				ExpectedErr:                  rpc.Err(rpc.InternalError, rpc.StringErrData("UnExpectedErr")),
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
				ExpectedErr:                  rpc.Err(rpc.InternalError, rpc.StringErrData(context.DeadlineExceeded.Error())),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), test.Timeout*time.Second)
			defer cancel()
			if test.ShouldCallTransactionReceipt {
				mockRpcProvider.EXPECT().TransactionReceipt(ctx, test.Hash).Return(test.ExpectedReceipt, test.ExpectedErr)
			}
			resp, err := acnt.WaitForTransactionReceipt(ctx, test.Hash, 2*time.Second)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
			} else {
				// check
				require.Equal(t, test.ExpectedReceipt.ExecutionStatus, (resp.TransactionReceipt).ExecutionStatus)
			}
		}()
	}
}

// TestWaitForTransactionReceipt is a test function that tests the WaitForTransactionReceipt method.
//
// It checks if the test environment is "devnet" and skips the test if it's not.
// It creates a new RPC client using the base URL and "/rpc" endpoint.
// It creates a new RPC provider using the client.
// It creates a new account using the provider, a zero-value Felt object, the "pubkey" string, and a new memory keystore.
// It defines a testSet variable that contains an array of testSetType structs.
// Each testSetType struct contains a Timeout integer, a Hash object, an ExpectedErr error, and an ExpectedReceipt TransactionReceipt object.
// It retrieves the testSet based on the testEnv variable.
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
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}
	client, err := rpc.NewProvider(tConfig.providerURL)
	require.NoError(t, err, "Error in rpc.NewProvider")

	acnt, err := account.NewAccount(client, &felt.Zero, "pubkey", account.NewMemKeystore(), account.CairoV0)
	require.NoError(t, err, "error returned from account.NewAccount()")

	type testSetType struct {
		Timeout         int
		Hash            *felt.Felt
		ExpectedErr     *rpc.RPCError
		ExpectedReceipt rpc.TransactionReceipt
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Timeout:         3, // Should poll 3 times
				Hash:            new(felt.Felt).SetUint64(100),
				ExpectedReceipt: rpc.TransactionReceipt{},
				ExpectedErr:     rpc.Err(rpc.InternalError, rpc.StringErrData("context deadline exceeded")),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(test.Timeout)*time.Second)
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
		}()
	}
}
