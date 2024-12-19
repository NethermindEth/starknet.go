package account_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/devnet"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	// set the environment for the test, default: mock
	testEnv = "mock"
	base    = ""
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
//
// It sets up the test environment by parsing command line flags and loading environment variables.
// The test environment can be set using the "env" flag.
// It then sets the base path for integration tests by reading the value from the "INTEGRATION_BASE" environment variable.
// If the base path is not set and the test environment is not "mock", it panics.
// Finally, it exits with the return value of the test suite
//
// Parameters:
// - m: is the test main
// Returns:
//
//	none
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()
	if testEnv == "mock" {
		return
	}
	base = os.Getenv("INTEGRATION_BASE")
	if base == "" {
		if err := godotenv.Load(fmt.Sprintf(".env.%s", testEnv)); err != nil {
			panic(fmt.Sprintf("Failed to load .env.%s, err: %s", testEnv, err))
		}
		base = os.Getenv("INTEGRATION_BASE")
	}
	os.Exit(m.Run())
}

// TestTransactionHashInvoke tests the TransactionHashInvoke function.
//
// This function tests the TransactionHashInvoke method of the Account struct.
// It generates a set of test cases and iterates over them to verify the correctness
// of the transaction hash. Each test case consists of the expected hash, a flag
// indicating whether the KeyStore should be set, account address, public key,
// private key, chain ID, function call, and transaction details.
//
// Parameters:
//   - t: The testing.T object for running the test
//
// Returns:
//
//	none
func TestTransactionHashInvoke(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	type testSetType struct {
		ExpectedHash   *felt.Felt
		SetKS          bool
		AccountAddress *felt.Felt
		PubKey         string
		PrivKey        *felt.Felt
		ChainID        string
		FnCall         rpc.FunctionCall
		TxDetails      rpc.TxDetails
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				// https://sepolia.voyager.online/tx/0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406
				ExpectedHash:   utils.TestHexToFelt(t, "0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406"),
				SetKS:          true,
				AccountAddress: utils.TestHexToFelt(t, "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d"),
				PrivKey:        utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				PubKey:         "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d",
				ChainID:        "SN_SEPOLIA",
				FnCall: rpc.FunctionCall{
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x517567ac7026ce129c950e6e113e437aa3c83716cd61481c6bb8c5057e6923e",
						"0xcaffbd1bd76bd7f24a3fa1d69d1b2588a86d1f9d2359b13f6a84b7e1cbd126",
						"0x7",
						"0x457874726163745265736f7572636546696e697368",
						"0x5",
						"0x5",
						"0xb82",
						"0x1",
						"0x1",
						"0x35c",
					}),
				},
				TxDetails: rpc.TxDetails{
					Nonce:   utils.TestHexToFelt(t, "0x3cf"),
					MaxFee:  utils.TestHexToFelt(t, "0x1a6f9d0dc5952"),
					Version: rpc.TransactionV1,
				},
			},
		},
		"devnet": {},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406
				ExpectedHash:   utils.TestHexToFelt(t, "0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406"),
				SetKS:          true,
				AccountAddress: utils.TestHexToFelt(t, "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d"),
				PrivKey:        utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				PubKey:         "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d",
				ChainID:        "SN_SEPOLIA",
				FnCall: rpc.FunctionCall{
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x517567ac7026ce129c950e6e113e437aa3c83716cd61481c6bb8c5057e6923e",
						"0xcaffbd1bd76bd7f24a3fa1d69d1b2588a86d1f9d2359b13f6a84b7e1cbd126",
						"0x7",
						"0x457874726163745265736f7572636546696e697368",
						"0x5",
						"0x5",
						"0xb82",
						"0x1",
						"0x1",
						"0x35c",
					}),
				},
				TxDetails: rpc.TxDetails{
					Nonce:   utils.TestHexToFelt(t, "0x3cf"),
					MaxFee:  utils.TestHexToFelt(t, "0x1a6f9d0dc5952"),
					Version: rpc.TransactionV1,
				},
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {

		t.Run("Transaction hash", func(t *testing.T) {
			ks := account.NewMemKeystore()
			if test.SetKS {
				privKeyBI, ok := new(big.Int).SetString(test.PrivKey.String(), 0)
				require.True(t, ok)
				ks.Put(test.PubKey, privKeyBI)
			}

			mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
			account, err := account.NewAccount(mockRpcProvider, test.AccountAddress, test.PubKey, ks, 0)
			require.NoError(t, err, "error returned from account.NewAccount()")
			invokeTxn := rpc.BroadcastInvokev1Txn{
				InvokeTxnV1: rpc.InvokeTxnV1{
					Calldata:      test.FnCall.Calldata,
					Nonce:         test.TxDetails.Nonce,
					MaxFee:        test.TxDetails.MaxFee,
					SenderAddress: account.AccountAddress,
					Version:       test.TxDetails.Version,
				}}
			hash, err := account.TransactionHashInvoke(invokeTxn.InvokeTxnV1)
			require.NoError(t, err, "error returned from account.TransactionHash()")
			require.Equal(t, test.ExpectedHash.String(), hash.String(), "transaction hash does not match expected")
		})
	}

}

// TestFmtCallData tests the FmtCallData function.
//
// It tests the FmtCallData function by providing different test sets
// and comparing the output with the expected call data.
//
// Parameters:
// - t: The testing.T instance for running the test
// Return:
//
//	none
func TestFmtCallData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	type testSetType struct {
		CairoVersion     int
		ChainID          string
		FnCall           rpc.FunctionCall
		ExpectedCallData []*felt.Felt
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mock":   {},
		"testnet": {
			{
				CairoVersion: 2,
				ChainID:      "SN_SEPOLIA",
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x04daadb9d30c887e1ab2cf7d78dfe444a77aab5a49c3353d6d9977e7ed669902"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("name_set"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x617279616e5f676f64617261"),
					},
				},
				ExpectedCallData: utils.TestHexArrToFelt(t, []string{
					"0x01",
					"0x04daadb9d30c887e1ab2cf7d78dfe444a77aab5a49c3353d6d9977e7ed669902",
					"0x0166d775d0cf161f1ce9b90698485f0c7a0e249af1c4b38126bddb37859737ac",
					"0x01",
					"0x617279616e5f676f64617261",
				}),
			},
			{
				CairoVersion: 2,
				ChainID:      "SN_SEPOLIA",
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x017cE9DffA7C87a03EB496c96e04ac36c4902085030763A83a35788d475e15CA"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("name_set"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x737461726b6e6574"),
					},
				},
				ExpectedCallData: utils.TestHexArrToFelt(t, []string{
					"0x01",
					"0x017ce9dffa7c87a03eb496c96e04ac36c4902085030763a83a35788d475e15ca",
					"0x0166d775d0cf161f1ce9b90698485f0c7a0e249af1c4b38126bddb37859737ac",
					"0x01",
					"0x737461726b6e6574",
				}),
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
		acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore(), test.CairoVersion)
		require.NoError(t, err)

		fmtCallData, err := acnt.FmtCalldata([]rpc.FunctionCall{test.FnCall})
		require.NoError(t, err)
		require.Equal(t, fmtCallData, test.ExpectedCallData)
	}
}

// TestChainIdMOCK is a test function that tests the behavior of the ChainId function.
//
// It creates a mock controller and a mock RpcProvider. It defines a test set
// consisting of different ChainID and ExpectedID pairs. It then iterates over
// the test set and sets the expected behavior for the ChainID method of the
// mockRpcProvider. It creates a new account using the mockRpcProvider,
// Zero value, "pubkey", and a new in-memory keystore. It asserts that the
// account's ChainId matches the expected ID for each test case in the test set.
//
// Parameters:
// - t: The testing.T instance for running the test
// Return:
//
//	none
func TestChainIdMOCK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	type testSetType struct {
		ChainID    string
		ExpectedID string
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mock": {
			{
				ChainID:    "SN_MAIN",
				ExpectedID: "0x534e5f4d41494e",
			},
			{
				ChainID:    "SN_SEPOLIA",
				ExpectedID: "0x534e5f5345504f4c4941",
			},
		},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
		account, err := account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore(), 0)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedID, account.ChainId.String())
	}
}

// TestChainId tests the ChainId function.
//
// This function tests the ChainId function by setting up a mock controller, defining a test set,
// and running a series of assertions on the expected results.
// It checks if the ChainId function returns the correct ChainID and ExpectedID values
// for different test environments.
// Parameters:
// - t: The testing.T instance for running the test
// Return:
//
//	none
func TestChainId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	type testSetType struct {
		ChainID    string
		ExpectedID string
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				ChainID:    "SN_SEPOLIA",
				ExpectedID: "0x534e5f5345504f4c4941",
			},
		},
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		client, err := rpc.NewProvider(base)
		require.NoError(t, err, "Error in rpc.NewClient")

		account, err := account.NewAccount(client, &felt.Zero, "pubkey", account.NewMemKeystore(), 0)
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), test.ExpectedID)
	}

}

// TestSignMOCK is a test function that tests the Sign method of the Account struct using mock objects.
//
// It sets up a mock controller and a mock RPC provider, and defines a test set containing different scenarios.
// Each scenario includes an address, private key, chain ID, a felt to sign, and the expected signatures.
// The function iterates over the test set and performs the following steps for each test case:
// - Converts the private key to a big.Int object and stores it in a memory keystore.
// - Mocks the ChainID method of the RPC provider to return the specified chain ID.
// - Creates an account using the mock RPC provider, the test address, the address string, and the keystore.
// - Converts the felt to sign to a big.Int object.
// - Calls the Sign method of the account with the felt to sign and retrieves the signature.
// - Verifies that the obtained signature matches the expected signature.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//
//	none
func TestSignMOCK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	type testSetType struct {
		Address     *felt.Felt
		PrivKey     *felt.Felt
		ChainId     string
		FeltToSign  *felt.Felt
		ExpectedSig []*felt.Felt
	}
	testSet := map[string][]testSetType{
		"mock": {
			// Accepted on testnet https://sepolia.voyager.online/tx/0x4b2e6743b03a0412f8450dd1d337f37a0e946603c3e6fbf4ba2469703c1705b
			{
				Address:    utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				PrivKey:    utils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"),
				ChainId:    "SN_SEPOLIA",
				FeltToSign: utils.TestHexToFelt(t, "0x4b2e6743b03a0412f8450dd1d337f37a0e946603c3e6fbf4ba2469703c1705b"),
				ExpectedSig: []*felt.Felt{
					utils.TestHexToFelt(t, "0xfa671736285eb70057579532f0efb6fde09ecefe323755ffd126537234e9c5"),
					utils.TestHexToFelt(t, "0x27bf55daa78a3ccfb7a4ee6576a13adfc44af707c28588be8292b8476bb27ef"),
				},
			},
		},
		"devnet": {},
		"testnet": {
			// Accepted on testnet https://sepolia.voyager.online/tx/0x4b2e6743b03a0412f8450dd1d337f37a0e946603c3e6fbf4ba2469703c1705b
			{
				Address:    utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				PrivKey:    utils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"),
				ChainId:    "SN_SEPOLIA",
				FeltToSign: utils.TestHexToFelt(t, "0x4b2e6743b03a0412f8450dd1d337f37a0e946603c3e6fbf4ba2469703c1705b"),
				ExpectedSig: []*felt.Felt{
					utils.TestHexToFelt(t, "0xfa671736285eb70057579532f0efb6fde09ecefe323755ffd126537234e9c5"),
					utils.TestHexToFelt(t, "0x27bf55daa78a3ccfb7a4ee6576a13adfc44af707c28588be8292b8476bb27ef"),
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		privKeyBI, ok := new(big.Int).SetString(test.PrivKey.String(), 0)
		require.True(t, ok)
		ks := account.NewMemKeystore()
		ks.Put(test.Address.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainId, nil)
		account, err := account.NewAccount(mockRpcProvider, test.Address, test.Address.String(), ks, 0)
		require.NoError(t, err, "error returned from account.NewAccount()")

		sig, err := account.Sign(context.Background(), test.FeltToSign)

		require.NoError(t, err, "error returned from account.Sign()")
		require.Equal(t, test.ExpectedSig[0].String(), sig[0].String(), "s1 does not match expected")
		require.Equal(t, test.ExpectedSig[1].String(), sig[1].String(), "s2 does not match expected")
	}

}

// TestAddInvoke is a test function that verifies the behavior of the AddInvokeTransaction method.
//
// This function tests the AddInvokeTransaction method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//
//	none
func TestSendInvokeTxn(t *testing.T) {

	type testSetType struct {
		ExpectedErr          error
		CairoContractVersion int
		SetKS                bool
		AccountAddress       *felt.Felt
		PubKey               *felt.Felt
		PrivKey              *felt.Felt
		InvokeTx             rpc.BroadcastInvokev1Txn
		FnCall               rpc.FunctionCall
		TxDetails            rpc.TxDetails
	}
	testSet := map[string][]testSetType{
		"mock":   {},
		"devnet": {},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x04b2e6743b03a0412f8450dd1d337f37a0e946603c3e6fbf4ba2469703c1705b
				ExpectedErr:          rpc.ErrDuplicateTx,
				CairoContractVersion: 2,
				AccountAddress:       utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e"),
				PrivKey:              utils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"), //
				InvokeTx: rpc.BroadcastInvokev1Txn{
					InvokeTxnV1: rpc.InvokeTxnV1{
						Nonce:         new(felt.Felt).SetUint64(5),
						MaxFee:        utils.TestHexToFelt(t, "0x26112A960026"),
						Version:       rpc.TransactionV1,
						Type:          rpc.TransactionType_Invoke,
						SenderAddress: utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
					}},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x04daadb9d30c887e1ab2cf7d78dfe444a77aab5a49c3353d6d9977e7ed669902"),
					EntryPointSelector: utils.TestHexToFelt(t, "0x166d775d0cf161f1ce9b90698485f0c7a0e249af1c4b38126bddb37859737ac"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x737461726b6e6574"),
					},
				},
			},
			{
				// https://sepolia.voyager.online/tx/0x32b46053f669fc198c2647bdc150c6a83d4a44a00e7d85fd10afca52706e6fa
				ExpectedErr:          rpc.ErrDuplicateTx,
				CairoContractVersion: 2,
				AccountAddress:       utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e"),
				PrivKey:              utils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"),
				InvokeTx: rpc.BroadcastInvokev1Txn{
					InvokeTxnV1: rpc.InvokeTxnV1{
						Nonce:         new(felt.Felt).SetUint64(8),
						MaxFee:        utils.TestHexToFelt(t, "0x1f6410500832"),
						Version:       rpc.TransactionV1,
						Type:          rpc.TransactionType_Invoke,
						SenderAddress: utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
					}},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x04daadb9d30c887e1ab2cf7d78dfe444a77aab5a49c3353d6d9977e7ed669902"),
					EntryPointSelector: utils.TestHexToFelt(t, "0x166d775d0cf161f1ce9b90698485f0c7a0e249af1c4b38126bddb37859737ac"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x617279616e5f676f64617261"),
					},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		client, err := rpc.NewProvider(base)
		require.NoError(t, err, "Error in rpc.NewClient")

		// Set up ks
		ks := account.NewMemKeystore()
		if test.SetKS {
			fakePrivKeyBI, ok := new(big.Int).SetString(test.PrivKey.String(), 0)
			require.True(t, ok)
			ks.Put(test.PubKey.String(), fakePrivKeyBI)
		}

		acnt, err := account.NewAccount(client, test.AccountAddress, test.PubKey.String(), ks, 2)
		require.NoError(t, err)

		test.InvokeTx.Calldata, err = acnt.FmtCalldata([]rpc.FunctionCall{test.FnCall})
		require.NoError(t, err)

		err = acnt.SignInvokeTransaction(context.Background(), &test.InvokeTx.InvokeTxnV1)
		require.NoError(t, err)

		resp, err := acnt.SendTransaction(context.Background(), test.InvokeTx)
		if err != nil {
			require.Equal(t, test.ExpectedErr.Error(), err.Error(), "AddInvokeTransaction returned an unexpected error")
			require.Nil(t, resp)
		}

	}
}

// TestAddDeployAccountDevnet tests the functionality of adding a deploy account in the devnet environment.
//
// The test checks if the environment is set to "devnet" and skips the test if not. It then initializes a new RPC client
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
	client, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

	devnet, acnts, err := newDevnet(t, base)
	require.NoError(t, err, "Error setting up Devnet")
	fakeUser := acnts[0]
	fakeUserAddr := utils.TestHexToFelt(t, fakeUser.Address)
	fakeUserPub := utils.TestHexToFelt(t, fakeUser.PublicKey)

	// Set up ks
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(fakeUser.PrivateKey, 0)
	require.True(t, ok)
	ks.Put(fakeUser.PublicKey, fakePrivKeyBI)

	acnt, err := account.NewAccount(client, fakeUserAddr, fakeUser.PublicKey, ks, 0)
	require.NoError(t, err)

	classHash := utils.TestHexToFelt(t, "0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f") // preDeployed classhash
	require.NoError(t, err)

	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		MaxFee:              utils.TestHexToFelt(t, "0xc5cb22092551"),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: fakeUserPub,
		ConstructorCalldata: []*felt.Felt{fakeUserPub},
	}

	precomputedAddress, err := acnt.PrecomputeAccountAddress(fakeUserPub, classHash, tx.ConstructorCalldata)
	require.Nil(t, err)
	require.NoError(t, acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress))

	_, err = devnet.Mint(precomputedAddress, new(big.Int).SetUint64(10000000000000000000))
	require.NoError(t, err)

	resp, err := acnt.SendTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
	require.Nil(t, err, "AddDeployAccountTransaction gave an Error")
	require.NotNil(t, resp, "AddDeployAccountTransaction resp not nil")
}

// TestTransactionHashDeclare tests the TransactionHashDeclare function.
//
// This function verifies that the TransactionHashDeclare function returns the
// expected hash value for a given transaction.
// The function requires a testnet environment to run.
// It creates a new client using the provided base URL and verifies that no
// error occurs.
// It then creates a new account using the provider and verifies that no error
// occurs.
// It constructs a DeclareTxnV2 struct with test hex values for the nonce,
// max fee, signature, sender address, compiled class hash, and class hash.
// Finally, it calls the TransactionHashDeclare function and compares the
// returned hash with the expected hash, ensuring they match.
//
// Parameters:
// - t: reference to the testing.T object
// Returns:
//
//	none
func TestTransactionHashDeclare(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
	mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)

	acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
	require.NoError(t, err)

	type testSetType struct {
		Txn          rpc.DeclareTxnType
		ExpectedHash *felt.Felt
		ExpectedErr  error
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				// https://sepolia.voyager.online/tx/0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3
				Txn: rpc.DeclareTxnV2{
					Nonce:   utils.TestHexToFelt(t, "0x1"),
					Type:    rpc.TransactionType_Declare,
					Version: rpc.TransactionV2,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x713765e220325edfaf5e033ad77b1ba4eceabe66333893b89845c2ddc744d34"),
						utils.TestHexToFelt(t, "0x4f28b1c15379c0ceb1855c09ed793e7583f875a802cbf310a8c0c971835c5cf")},
					SenderAddress:     utils.TestHexToFelt(t, "0x0019bd7ebd72368deb5f160f784e21aa46cd09e06a61dc15212456b5597f47b8"),
					CompiledClassHash: utils.TestHexToFelt(t, "0x017f655f7a639a49ea1d8d56172e99cff8b51f4123b733f0378dfd6378a2cd37"),
					ClassHash:         utils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
					MaxFee:            utils.TestHexToFelt(t, "0x177e06ff6cab2"),
				},
				ExpectedHash: utils.TestHexToFelt(t, "0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3"),
				ExpectedErr:  nil,
			},
		},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3
				Txn: rpc.DeclareTxnV2{
					Nonce:   utils.TestHexToFelt(t, "0x1"),
					Type:    rpc.TransactionType_Declare,
					Version: rpc.TransactionV2,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x713765e220325edfaf5e033ad77b1ba4eceabe66333893b89845c2ddc744d34"),
						utils.TestHexToFelt(t, "0x4f28b1c15379c0ceb1855c09ed793e7583f875a802cbf310a8c0c971835c5cf")},
					SenderAddress:     utils.TestHexToFelt(t, "0x0019bd7ebd72368deb5f160f784e21aa46cd09e06a61dc15212456b5597f47b8"),
					CompiledClassHash: utils.TestHexToFelt(t, "0x017f655f7a639a49ea1d8d56172e99cff8b51f4123b733f0378dfd6378a2cd37"),
					ClassHash:         utils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
					MaxFee:            utils.TestHexToFelt(t, "0x177e06ff6cab2"),
				},
				ExpectedHash: utils.TestHexToFelt(t, "0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3"),
				ExpectedErr:  nil,
			},
		},
	}[testEnv]
	for _, test := range testSet {
		hash, err := acnt.TransactionHashDeclare(test.Txn)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hash.String(), "TransactionHashDeclare not what expected")
	}
}

func TestTransactionHashInvokeV3(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
	mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)

	acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
	require.NoError(t, err)

	type testSetType struct {
		Txn          rpc.DeclareTxnType
		ExpectedHash *felt.Felt
		ExpectedErr  error
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				// https://sepolia.voyager.online/tx/0x8eb1104170ec42fd27c09ea78822dfb083ddd15324480f856bff01bc65e9d9
				Txn: rpc.InvokeTxnV3{
					Nonce:   utils.TestHexToFelt(t, "0x12eaa"),
					Type:    rpc.TransactionType_Invoke,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7121c34d7073fd21b73801000278883b332a6f8cdf90d7a84358748de811480"),
						utils.TestHexToFelt(t, "0x3df8c38724b89e9baa8dc0d0cb8fd14a8ca65308d7ca831793cc67394803b6c")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x2b",
							MaxPricePerUnit: "0x2eb31cc948ef",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					SenderAddress:         utils.TestHexToFelt(t, "0x1d091b30a2d20ca2509579f8beae26934bfdc3725c0b497f50b353b7a3c636f"),
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x132303a40ae2f271f4e1b707596a63f6f2921c4d400b38822548ed1bb0cbe0",
						"0xc844fd57777b0cd7e75c8ea68deec0adf964a6308da7a58de32364b7131cc8",
						"0x13",
						"0x46f1db039a8aa6edb473195b98421579517d79bbe026e74bbc3f172af0798",
						"0x1c4104",
						"0xde9f47f476fed1e72a6159aba4f458ac74cbbbf88a951758a1fc0276e27211",
						"0x66417c6a",
						"0x104030200000000000000000000000000000000000000000000000000000000",
						"0x4",
						"0x431d563dc0",
						"0x4329326f60",
						"0x432cdda6a2",
						"0x433403bcd6",
						"0xbc2ee78d0e41b9dd1",
						"0x1",
						"0x2",
						"0x75a8626edb90cc9983ae1dfca05c485c8ca6ca507f925ac8f28366aa8d7c211",
						"0x83b812e83b07feb3e898aa55db8552138580d63e4f827e28e2531bd308db29",
						"0x2e7dc996ebf724c1cf18d668fc3455df4245749ebc0724101cbc6c9cb13c962",
						"0x49e384b4c21fbb10318f461c7804432e068c7ff196647b4f3b470b4431c40e6",
						"0x2389f278922589f5f5d39b17339dc7ef80f13c8eb20b173c9eba52503c60874",
						"0x4225d1c8ee8e451a25e30c10689ef898e11ccf5c0f68d0fc7876c47b318e946",
					}),
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
				},
				ExpectedHash: utils.TestHexToFelt(t, "0x8eb1104170ec42fd27c09ea78822dfb083ddd15324480f856bff01bc65e9d9"),
				ExpectedErr:  nil,
			},
		},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x8eb1104170ec42fd27c09ea78822dfb083ddd15324480f856bff01bc65e9d9
				Txn: rpc.InvokeTxnV3{
					Nonce:   utils.TestHexToFelt(t, "0x12eaa"),
					Type:    rpc.TransactionType_Invoke,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7121c34d7073fd21b73801000278883b332a6f8cdf90d7a84358748de811480"),
						utils.TestHexToFelt(t, "0x3df8c38724b89e9baa8dc0d0cb8fd14a8ca65308d7ca831793cc67394803b6c")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x2b",
							MaxPricePerUnit: "0x2eb31cc948ef",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					SenderAddress:         utils.TestHexToFelt(t, "0x1d091b30a2d20ca2509579f8beae26934bfdc3725c0b497f50b353b7a3c636f"),
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x132303a40ae2f271f4e1b707596a63f6f2921c4d400b38822548ed1bb0cbe0",
						"0xc844fd57777b0cd7e75c8ea68deec0adf964a6308da7a58de32364b7131cc8",
						"0x13",
						"0x46f1db039a8aa6edb473195b98421579517d79bbe026e74bbc3f172af0798",
						"0x1c4104",
						"0xde9f47f476fed1e72a6159aba4f458ac74cbbbf88a951758a1fc0276e27211",
						"0x66417c6a",
						"0x104030200000000000000000000000000000000000000000000000000000000",
						"0x4",
						"0x431d563dc0",
						"0x4329326f60",
						"0x432cdda6a2",
						"0x433403bcd6",
						"0xbc2ee78d0e41b9dd1",
						"0x1",
						"0x2",
						"0x75a8626edb90cc9983ae1dfca05c485c8ca6ca507f925ac8f28366aa8d7c211",
						"0x83b812e83b07feb3e898aa55db8552138580d63e4f827e28e2531bd308db29",
						"0x2e7dc996ebf724c1cf18d668fc3455df4245749ebc0724101cbc6c9cb13c962",
						"0x49e384b4c21fbb10318f461c7804432e068c7ff196647b4f3b470b4431c40e6",
						"0x2389f278922589f5f5d39b17339dc7ef80f13c8eb20b173c9eba52503c60874",
						"0x4225d1c8ee8e451a25e30c10689ef898e11ccf5c0f68d0fc7876c47b318e946",
					}),
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
				},
				ExpectedHash: utils.TestHexToFelt(t, "0x8eb1104170ec42fd27c09ea78822dfb083ddd15324480f856bff01bc65e9d9"),
				ExpectedErr:  nil,
			},
		},
	}[testEnv]
	for _, test := range testSet {
		hash, err := acnt.TransactionHashInvoke(test.Txn)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hash.String(), "TransactionHashDeclare not what expected")
	}
}

func TestTransactionHashdeployAccount(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
	mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)

	acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
	require.NoError(t, err)

	type testSetType struct {
		Txn           rpc.DeployAccountType
		SenderAddress *felt.Felt
		ExpectedHash  *felt.Felt
		ExpectedErr   error
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				// https://sepolia.voyager.online/tx/0x66d1d9d50d308a9eb16efedbad208b0672769a545a0b828d357757f444e9188
				Txn: rpc.DeployAccountTxn{
					Nonce:   utils.TestHexToFelt(t, "0x0"),
					Type:    rpc.TransactionType_DeployAccount,
					MaxFee:  utils.TestHexToFelt(t, "0x1d2109b99cf94"),
					Version: rpc.TransactionV1,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x427df9a1a4a0b7b9011a758524b8a6c2595aac9140608fe24c66efe04b340d7"),
						utils.TestHexToFelt(t, "0x4edc73cd97dab7458a08fec6d7c0e1638c3f1111646fc8a91508b4f94b36310"),
					},
					ClassHash:           utils.TestHexToFelt(t, "0x1e60c8722677cfb7dd8dbea5be86c09265db02cdfe77113e77da7d44c017388"),
					ContractAddressSalt: utils.TestHexToFelt(t, "0x15d621f9515c6197d3117eb1a25c7a4a669317be8f49831e03fcc00d855352e"),
					ConstructorCalldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x960532cfba33384bbec41aa669727a9c51e995c87e101c86706aaf244f7e4e"),
					},
				},
				SenderAddress: utils.TestHexToFelt(t, "0x05dd5faeddd4a9e01231f3bb9b95ec93426d08977b721c222e45fd98c5f353ff"),
				ExpectedHash:  utils.TestHexToFelt(t, "0x66d1d9d50d308a9eb16efedbad208b0672769a545a0b828d357757f444e9188"),
				ExpectedErr:   nil,
			},
			{
				// https://sepolia.voyager.online/tx/0x4bf28fb0142063f1b9725ae490c6949e6f1842c79b49f7cc674b7e3f5ad4875
				Txn: rpc.DeployAccountTxnV3{
					Nonce:   utils.TestHexToFelt(t, "0x0"),
					Type:    rpc.TransactionType_DeployAccount,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0xaa580d6fd4bc056d6a9a49833e7fc966fe5f20cc283e05854e44a5d4516958"),
						utils.TestHexToFelt(t, "0x41a57fcb19908321f8e44c425ea419a1de272efd99888503ee0cdc0ddb6aee4")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x38",
							MaxPricePerUnit: "0x7cd9b6080b35",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:           "0x0",
					PayMasterData: []*felt.Felt{},
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
					ClassHash:     utils.TestHexToFelt(t, "0x29927c8af6bccf3f6fda035981e765a7bdbf18a2dc0d630494f8758aa908e2b"),
					ConstructorCalldata: utils.TestHexArrToFelt(t, []string{
						"0x1a09f0001cc46f82b1a805d07c13e235248a44ed13d87f170d7d925e3c86082",
						"0x0",
					}),
					ContractAddressSalt: utils.TestHexToFelt(t, "0x1a09f0001cc46f82b1a805d07c13e235248a44ed13d87f170d7d925e3c86082"),
				},
				SenderAddress: utils.TestHexToFelt(t, "0x0365633b6c2ca24b461747d2fe8e0c19a3637a954ee703a7ed0e5d1d9644ad1a"),
				ExpectedHash:  utils.TestHexToFelt(t, "0x4bf28fb0142063f1b9725ae490c6949e6f1842c79b49f7cc674b7e3f5ad4875"),
				ExpectedErr:   nil,
			},
		},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x66d1d9d50d308a9eb16efedbad208b0672769a545a0b828d357757f444e9188
				Txn: rpc.DeployAccountTxn{
					Nonce:   utils.TestHexToFelt(t, "0x0"),
					Type:    rpc.TransactionType_DeployAccount,
					MaxFee:  utils.TestHexToFelt(t, "0x1d2109b99cf94"),
					Version: rpc.TransactionV1,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x427df9a1a4a0b7b9011a758524b8a6c2595aac9140608fe24c66efe04b340d7"),
						utils.TestHexToFelt(t, "0x4edc73cd97dab7458a08fec6d7c0e1638c3f1111646fc8a91508b4f94b36310"),
					},
					ClassHash:           utils.TestHexToFelt(t, "0x1e60c8722677cfb7dd8dbea5be86c09265db02cdfe77113e77da7d44c017388"),
					ContractAddressSalt: utils.TestHexToFelt(t, "0x15d621f9515c6197d3117eb1a25c7a4a669317be8f49831e03fcc00d855352e"),
					ConstructorCalldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x960532cfba33384bbec41aa669727a9c51e995c87e101c86706aaf244f7e4e"),
					},
				},
				SenderAddress: utils.TestHexToFelt(t, "0x05dd5faeddd4a9e01231f3bb9b95ec93426d08977b721c222e45fd98c5f353ff"),
				ExpectedHash:  utils.TestHexToFelt(t, "0x66d1d9d50d308a9eb16efedbad208b0672769a545a0b828d357757f444e9188"),
				ExpectedErr:   nil,
			},
			{
				// https://sepolia.voyager.online/tx/0x4bf28fb0142063f1b9725ae490c6949e6f1842c79b49f7cc674b7e3f5ad4875
				Txn: rpc.DeployAccountTxnV3{
					Nonce:   utils.TestHexToFelt(t, "0x0"),
					Type:    rpc.TransactionType_DeployAccount,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0xaa580d6fd4bc056d6a9a49833e7fc966fe5f20cc283e05854e44a5d4516958"),
						utils.TestHexToFelt(t, "0x41a57fcb19908321f8e44c425ea419a1de272efd99888503ee0cdc0ddb6aee4")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x38",
							MaxPricePerUnit: "0x7cd9b6080b35",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:           "0x0",
					PayMasterData: []*felt.Felt{},
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
					ClassHash:     utils.TestHexToFelt(t, "0x29927c8af6bccf3f6fda035981e765a7bdbf18a2dc0d630494f8758aa908e2b"),
					ConstructorCalldata: utils.TestHexArrToFelt(t, []string{
						"0x1a09f0001cc46f82b1a805d07c13e235248a44ed13d87f170d7d925e3c86082",
						"0x0",
					}),
					ContractAddressSalt: utils.TestHexToFelt(t, "0x1a09f0001cc46f82b1a805d07c13e235248a44ed13d87f170d7d925e3c86082"),
				},
				SenderAddress: utils.TestHexToFelt(t, "0x0365633b6c2ca24b461747d2fe8e0c19a3637a954ee703a7ed0e5d1d9644ad1a"),
				ExpectedHash:  utils.TestHexToFelt(t, "0x4bf28fb0142063f1b9725ae490c6949e6f1842c79b49f7cc674b7e3f5ad4875"),
				ExpectedErr:   nil,
			},
		},
	}[testEnv]
	for _, test := range testSet {
		hash, err := acnt.TransactionHashDeployAccount(test.Txn, test.SenderAddress)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hash.String(), "TransactionHashDeclare not what expected")
	}
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
// - t: The testing.T object for test assertions and logging
// Returns:
//
//	none
func TestWaitForTransactionReceiptMOCK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)
	acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
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
				ExpectedErr:                  rpc.Err(rpc.InternalError, "UnExpectedErr"),
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
				ExpectedErr:                  rpc.Err(rpc.InternalError, context.DeadlineExceeded),
			},
		},
	}[testEnv]

	for _, test := range testSet {
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
// - t: The testing.T instance for running the test
// Returns:
//
//	none
func TestWaitForTransactionReceipt(t *testing.T) {
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}
	client, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

	acnt, err := account.NewAccount(client, &felt.Zero, "pubkey", account.NewMemKeystore(), 0)
	require.NoError(t, err, "error returned from account.NewAccount()")

	type testSetType struct {
		Timeout         int
		Hash            *felt.Felt
		ExpectedErr     error
		ExpectedReceipt rpc.TransactionReceipt
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Timeout:         3, // Should poll 3 times
				Hash:            new(felt.Felt).SetUint64(100),
				ExpectedReceipt: rpc.TransactionReceipt{},
				ExpectedErr:     rpc.Err(rpc.InternalError, "Post \"http://0.0.0.0:5050/\": context deadline exceeded"),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(test.Timeout)*time.Second)
		defer cancel()

		resp, err := acnt.WaitForTransactionReceipt(ctx, test.Hash, 1*time.Second)
		if test.ExpectedErr != nil {
			require.Equal(t, test.ExpectedErr.Error(), err.Error())
		} else {
			require.Equal(t, test.ExpectedReceipt.ExecutionStatus, (*resp).ExecutionStatus)
		}

	}
}

// TestAddDeclareTxn is a test function that verifies the behavior of the AddDeclareTransaction method.
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
	expectedTxHash := utils.TestHexToFelt(t, "0x0272ebd99f5d0a275b4bc26781f76c4c4e48050ce5f1c1ddafcdee48f0297255")
	expectedClassHash := utils.TestHexToFelt(t, "0x05e507b062836a3d73e71686ee62bca69026df94e72a657cbe0b954e6d3a0ce6")

	AccountAddress := utils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9")
	PubKey := utils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e")
	PrivKey := utils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9")

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(PrivKey.String(), 0)
	require.True(t, ok)
	ks.Put(PubKey.String(), fakePrivKeyBI)

	client, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

	acnt, err := account.NewAccount(client, AccountAddress, PubKey.String(), ks, 0)
	require.NoError(t, err)

	// Class Hash
	content, err := os.ReadFile("./tests/hello_world_compiled.sierra.json")
	require.NoError(t, err)

	var class rpc.ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
	classHash := hash.ClassHash(class)

	// Compiled Class Hash
	content2, err := os.ReadFile("./tests/hello_world_compiled.casm.json")
	require.NoError(t, err)

	var casmClass contracts.CasmClass
	err = json.Unmarshal(content2, &casmClass)
	require.NoError(t, err)
	compClassHash := hash.CompiledClassHash(casmClass)

	tx := rpc.DeclareTxnV2{
		Nonce:   utils.TestHexToFelt(t, "0xd"),
		MaxFee:  utils.TestHexToFelt(t, "0xc5cb22092551"),
		Type:    rpc.TransactionType_Declare,
		Version: rpc.TransactionV2,
		Signature: []*felt.Felt{
			utils.TestHexToFelt(t, "0x2975276c978f3cfbfa621b71085a910fe92ec32ba5995d8d70cfdd9c6db0ece"),
			utils.TestHexToFelt(t, "0x2f6eb4f42809ae38c8dfea82018451330ddcb276b63dde3ca8c64815e8f2fc0"),
		},
		SenderAddress:     AccountAddress,
		CompiledClassHash: compClassHash,
		ClassHash:         classHash,
	}

	err = acnt.SignDeclareTransaction(context.Background(), &tx)
	require.NoError(t, err)

	broadcastTx := rpc.BroadcastDeclareTxnV2{
		Nonce:             tx.Nonce,
		MaxFee:            tx.MaxFee,
		Type:              tx.Type,
		Version:           tx.Version,
		Signature:         tx.Signature,
		SenderAddress:     tx.SenderAddress,
		CompiledClassHash: tx.CompiledClassHash,
		ContractClass:     class,
	}

	resp, err := acnt.SendTransaction(context.Background(), broadcastTx)

	if err != nil {
		require.Equal(t, rpc.ErrDuplicateTx.Error(), err.Error(), "AddDeclareTransaction error not what expected")
	} else {
		require.Equal(t, expectedTxHash.String(), resp.TransactionHash.String(), "AddDeclareTransaction TxHash not what expected")
		require.Equal(t, expectedClassHash.String(), resp.ClassHash.String(), "AddDeclareTransaction ClassHash not what expected")
	}
}

// newDevnet creates a new devnet with the given URL.
//
// Parameters:
// - t: The testing.T instance for running the test
// - url: The URL of the devnet to be created
// Returns:
// - *devnet.DevNet: a pointer to a devnet object
// - []devnet.TestAccount: a slice of test accounts
// - error: an error, if any
func newDevnet(t *testing.T, url string) (*devnet.DevNet, []devnet.TestAccount, error) {
	t.Helper()
	devnet := devnet.NewDevNet(url)
	acnts, err := devnet.Accounts()
	return devnet, acnts, err
}
