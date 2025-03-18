package account_test

import (
	"context"
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
	"github.com/NethermindEth/starknet.go/internal"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	// the environment for the test, default: mock
	testEnv = ""
	// the base url for the test
	base = ""
	// the test account data
	privKey        = ""
	pubKey         = ""
	accountAddress = ""
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
//
// It sets up the test environment by parsing command line flags and loading environment variables.
// The test environment can be set using the "env" flag.
// It then sets the base path for integration tests by reading the value from the "HTTP_PROVIDER_URL" environment variable.
// If the base path is not set and the test environment is not "mock", it panics.
// Finally, it exits with the return value of the test suite
//
// Parameters:
// - m: is the test main
// Returns:
//
//	none
func TestMain(m *testing.M) {
	testEnv = internal.LoadEnv()

	if testEnv == "mock" {
		os.Exit(m.Run())
	}
	base = os.Getenv("HTTP_PROVIDER_URL")
	if base == "" {
		panic("Failed to load HTTP_PROVIDER_URL, empty string")
	}

	// load the test account data, only required for some tests
	privKey = os.Getenv("STARKNET_PRIVATE_KEY")
	pubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	accountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	os.Exit(m.Run())
}

func setupAcc(t *testing.T, provider rpc.RpcProvider) (*account.Account, error) {
	t.Helper()

	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privKey, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert privKey to big.Int")
	}
	ks.Put(pubKey, privKeyBI)

	accAddress, err := internalUtils.HexToFelt(accountAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert accountAddress to felt: %w", err)
	}

	acc, err := account.NewAccount(provider, accAddress, pubKey, ks, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return acc, nil
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
				ExpectedHash:   internalUtils.TestHexToFelt(t, "0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406"),
				SetKS:          true,
				AccountAddress: internalUtils.TestHexToFelt(t, "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d"),
				PrivKey:        internalUtils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				PubKey:         "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d",
				ChainID:        "SN_SEPOLIA",
				FnCall: rpc.FunctionCall{
					Calldata: internalUtils.TestHexArrToFelt(t, []string{
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
					Nonce:   internalUtils.TestHexToFelt(t, "0x3cf"),
					MaxFee:  internalUtils.TestHexToFelt(t, "0x1a6f9d0dc5952"),
					Version: rpc.TransactionV1,
				},
			},
		},
		"devnet": {},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406
				ExpectedHash:   internalUtils.TestHexToFelt(t, "0x5d307ad21a407ab6e93754b2fca71dd2d3b28313f6e844a7f3ecc404263a406"),
				SetKS:          true,
				AccountAddress: internalUtils.TestHexToFelt(t, "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d"),
				PrivKey:        internalUtils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				PubKey:         "0x06fb2806bc2564827796e0796144f8104581acdcbcd7721615ad376f70baf87d",
				ChainID:        "SN_SEPOLIA",
				FnCall: rpc.FunctionCall{
					Calldata: internalUtils.TestHexArrToFelt(t, []string{
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
					Nonce:   internalUtils.TestHexToFelt(t, "0x3cf"),
					MaxFee:  internalUtils.TestHexToFelt(t, "0x1a6f9d0dc5952"),
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
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x04daadb9d30c887e1ab2cf7d78dfe444a77aab5a49c3353d6d9977e7ed669902"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name_set"),
					Calldata: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x617279616e5f676f64617261"),
					},
				},
				ExpectedCallData: internalUtils.TestHexArrToFelt(t, []string{
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
					ContractAddress:    internalUtils.TestHexToFelt(t, "0x017cE9DffA7C87a03EB496c96e04ac36c4902085030763A83a35788d475e15CA"),
					EntryPointSelector: internalUtils.GetSelectorFromNameFelt("name_set"),
					Calldata: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x737461726b6e6574"),
					},
				},
				ExpectedCallData: internalUtils.TestHexArrToFelt(t, []string{
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
				Address:    internalUtils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				PrivKey:    internalUtils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"),
				ChainId:    "SN_SEPOLIA",
				FeltToSign: internalUtils.TestHexToFelt(t, "0x4b2e6743b03a0412f8450dd1d337f37a0e946603c3e6fbf4ba2469703c1705b"),
				ExpectedSig: []*felt.Felt{
					internalUtils.TestHexToFelt(t, "0xfa671736285eb70057579532f0efb6fde09ecefe323755ffd126537234e9c5"),
					internalUtils.TestHexToFelt(t, "0x27bf55daa78a3ccfb7a4ee6576a13adfc44af707c28588be8292b8476bb27ef"),
				},
			},
		},
		"devnet":  {},
		"testnet": {},
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
		InvokeTx             rpc.BroadcastInvokev3Txn
	}
	testSet := map[string][]testSetType{
		"mock":   {},
		"devnet": {},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x7aac4792c8fd7578dd01b20ff04565f2e2ce6ea3c792c5e609a088704c1dd87
				ExpectedErr:          rpc.ErrDuplicateTx,
				CairoContractVersion: 2,
				AccountAddress:       internalUtils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9"),
				SetKS:                true,
				PubKey:               internalUtils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e"),
				PrivKey:              internalUtils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9"), //
				InvokeTx: rpc.BroadcastInvokev3Txn{
					InvokeTxnV3: rpc.InvokeTxnV3{
						Nonce:   internalUtils.TestHexToFelt(t, "0xd"),
						Type:    rpc.TransactionType_Invoke,
						Version: rpc.TransactionV3,
						Signature: []*felt.Felt{
							internalUtils.TestHexToFelt(t, "0x7bff07f1c2f6dc0eeaa9e622a0ee35f6e2e9855b39ed757236970a71b7c9e2e"),
							internalUtils.TestHexToFelt(t, "0x588b821ccb9f61ca217bfb0a580f889886742c2fd63526009eb401a9cf951e3"),
						},
						ResourceBounds: rpc.ResourceBoundsMapping{
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

		err = acnt.SignInvokeTransaction(context.Background(), &test.InvokeTx.InvokeTxnV3)
		require.NoError(t, err)

		resp, err := acnt.SendTransaction(context.Background(), test.InvokeTx)
		if err != nil {
			require.Equal(t, test.ExpectedErr.Error(), err.Error(), "AddInvokeTransaction returned an unexpected error")
			require.Nil(t, resp)
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
	expectedTxHash := internalUtils.TestHexToFelt(t, "0x1c3df33f06f0da7f5df72bbc02fb8caf33e91bdd2433305dd007c6cd6acc6d0")
	expectedClassHash := internalUtils.TestHexToFelt(t, "0x06ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17")

	AccountAddress := internalUtils.TestHexToFelt(t, "0x01AE6Fe02FcD9f61A3A8c30D68a8a7c470B0d7dD6F0ee685d5BBFa0d79406ff9")
	PubKey := internalUtils.TestHexToFelt(t, "0x022288424ec8116c73d2e2ed3b0663c5030d328d9c0fb44c2b54055db467f31e")
	PrivKey := internalUtils.TestHexToFelt(t, "0x04818374f8071c3b4c3070ff7ce766e7b9352628df7b815ea4de26e0fadb5cc9")

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(PrivKey.String(), 0)
	require.True(t, ok)
	ks.Put(PubKey.String(), fakePrivKeyBI)

	client, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

	acnt, err := account.NewAccount(client, AccountAddress, PubKey.String(), ks, 0)
	require.NoError(t, err)

	// Class
	class := *internalUtils.TestUnmarshalJSONFileToType[contracts.ContractClass](t, "./tests/contracts_v2_HelloStarknet.sierra.json", "")

	// Compiled Class Hash
	casmClass := *internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](t, "./tests/contracts_v2_HelloStarknet.casm.json", "")
	compClassHash, err := hash.CompiledClassHash(casmClass)
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
		ResourceBounds: rpc.ResourceBoundsMapping{
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
		require.Equal(t, expectedTxHash.String(), resp.TransactionHash.String(), "AddDeclareTransaction TxHash not what expected")
		require.Equal(t, expectedClassHash.String(), resp.ClassHash.String(), "AddDeclareTransaction ClassHash not what expected")
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
	fakeUserAddr := internalUtils.TestHexToFelt(t, fakeUser.Address)
	fakeUserPub := internalUtils.TestHexToFelt(t, fakeUser.PublicKey)

	// Set up ks
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(fakeUser.PrivateKey, 0)
	require.True(t, ok)
	ks.Put(fakeUser.PublicKey, fakePrivKeyBI)

	acnt, err := account.NewAccount(client, fakeUserAddr, fakeUser.PublicKey, ks, 0)
	require.NoError(t, err)

	classHash := internalUtils.TestHexToFelt(t, "0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f") // preDeployed classhash
	require.NoError(t, err)

	tx := rpc.DeployAccountTxnV3{
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV3,
		Signature:           []*felt.Felt{},
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		ContractAddressSalt: fakeUserPub,
		ConstructorCalldata: []*felt.Felt{fakeUserPub},
		ClassHash:           classHash,
		ResourceBounds: rpc.ResourceBoundsMapping{
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
		Tip:           "0x0",
		PayMasterData: []*felt.Felt{},
		NonceDataMode: rpc.DAModeL1,
		FeeMode:       rpc.DAModeL1,
	}

	precomputedAddress := account.PrecomputeAccountAddress(fakeUserPub, classHash, tx.ConstructorCalldata)
	require.NoError(t, acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress))

	_, err = devnet.Mint(precomputedAddress, new(big.Int).SetUint64(10000000000000000000))
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
					Nonce:   internalUtils.TestHexToFelt(t, "0x1"),
					Type:    rpc.TransactionType_Declare,
					Version: rpc.TransactionV2,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x713765e220325edfaf5e033ad77b1ba4eceabe66333893b89845c2ddc744d34"),
						internalUtils.TestHexToFelt(t, "0x4f28b1c15379c0ceb1855c09ed793e7583f875a802cbf310a8c0c971835c5cf")},
					SenderAddress:     internalUtils.TestHexToFelt(t, "0x0019bd7ebd72368deb5f160f784e21aa46cd09e06a61dc15212456b5597f47b8"),
					CompiledClassHash: internalUtils.TestHexToFelt(t, "0x017f655f7a639a49ea1d8d56172e99cff8b51f4123b733f0378dfd6378a2cd37"),
					ClassHash:         internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
					MaxFee:            internalUtils.TestHexToFelt(t, "0x177e06ff6cab2"),
				},
				ExpectedHash: internalUtils.TestHexToFelt(t, "0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3"),
				ExpectedErr:  nil,
			},
			{
				// https://sepolia.voyager.online/tx/0x30c852c522274765e1d681bc8a84ce7c41118370ef2ba7d18a427ed29f5b155
				Txn: rpc.DeclareTxnV3{
					Nonce:   internalUtils.TestHexToFelt(t, "0x2b"),
					Type:    rpc.TransactionType_Declare,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x5c6a94302ef4b6d80a4c6a3eaf5ad30e11fa13aa78f7397a4f69901ceb12b7"),
						internalUtils.TestHexToFelt(t, "0x25bf97f481061f8abf5eb93e67eaebe6bb74dda34d7378a506f5ee2ff1daef1"),
					},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x10968159929e",
						},
						L1DataGas: rpc.ResourceBounds{
							MaxAmount:       "0x120",
							MaxPricePerUnit: "0x99f",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x1ff3ec0",
							MaxPricePerUnit: "0x197aa1ce3",
						},
					},
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					SenderAddress:         internalUtils.TestHexToFelt(t, "0x36d67ab362562a97f9fba8a1051cf8e37ff1a1449530fb9f1f0e32ac2da7d06"),
					ClassHash:             internalUtils.TestHexToFelt(t, "0x224518978adb773cfd4862a894e9d333192fbd24bc83841dc7d4167c09b89c5"),
					CompiledClassHash:     internalUtils.TestHexToFelt(t, "0x6ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17"),
					AccountDeploymentData: []*felt.Felt{},
					NonceDataMode:         rpc.DAModeL1,
					FeeMode:               rpc.DAModeL1,
				},
				ExpectedHash: internalUtils.TestHexToFelt(t, "0x30c852c522274765e1d681bc8a84ce7c41118370ef2ba7d18a427ed29f5b155"),
				ExpectedErr:  nil,
			},
		},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3
				Txn: rpc.DeclareTxnV2{
					Nonce:   internalUtils.TestHexToFelt(t, "0x1"),
					Type:    rpc.TransactionType_Declare,
					Version: rpc.TransactionV2,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x713765e220325edfaf5e033ad77b1ba4eceabe66333893b89845c2ddc744d34"),
						internalUtils.TestHexToFelt(t, "0x4f28b1c15379c0ceb1855c09ed793e7583f875a802cbf310a8c0c971835c5cf")},
					SenderAddress:     internalUtils.TestHexToFelt(t, "0x0019bd7ebd72368deb5f160f784e21aa46cd09e06a61dc15212456b5597f47b8"),
					CompiledClassHash: internalUtils.TestHexToFelt(t, "0x017f655f7a639a49ea1d8d56172e99cff8b51f4123b733f0378dfd6378a2cd37"),
					ClassHash:         internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
					MaxFee:            internalUtils.TestHexToFelt(t, "0x177e06ff6cab2"),
				},
				ExpectedHash: internalUtils.TestHexToFelt(t, "0x28e430cc73715bd1052e8db4f17b053c53dd8174341cba4b1a337b9fecfa8c3"),
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
				// https://sepolia.voyager.online/tx/0x76b52e17bc09064bd986ead34263e6305ef3cecfb3ae9e19b86bf4f1a1a20ea
				Txn: rpc.InvokeTxnV3{
					Nonce:   internalUtils.TestHexToFelt(t, "0x9803"),
					Type:    rpc.TransactionType_Invoke,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x17bacc700df6c82682139e8e550078a5daa75dfe356577f78f7e57fd7c56245"),
						internalUtils.TestHexToFelt(t, "0x4eb8734727eb9412b79ba6d14ff1c9a6beb0dc0b811e3f97168c747f8d427b3")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x2d79883d20000",
						},
						L1DataGas: rpc.ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x2d79883d20000",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x5f5e100",
							MaxPricePerUnit: "0xba43b7400",
						},
					},
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					SenderAddress:         internalUtils.TestHexToFelt(t, "0x745d525a3582e91299d8d7c71730ffc4b1f191f5b219d800334bc0edad0983b"),
					Calldata: internalUtils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x4138fd51f90d171df37e9d4419c8cdb67d525840c58f8a5c347be93a1c5277d",
						"0x2468d193cd15b621b24c2a602b8dbcfa5eaa14f88416c40c09d7fd12592cb4b",
						"0x0",
					}),
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
				},
				ExpectedHash: internalUtils.TestHexToFelt(t, "0x76b52e17bc09064bd986ead34263e6305ef3cecfb3ae9e19b86bf4f1a1a20ea"),
				ExpectedErr:  nil,
			},
		},
		"testnet": {
			{
				// https://sepolia.voyager.online/tx/0x76b52e17bc09064bd986ead34263e6305ef3cecfb3ae9e19b86bf4f1a1a20ea
				Txn: rpc.InvokeTxnV3{
					Nonce:   internalUtils.TestHexToFelt(t, "0x9803"),
					Type:    rpc.TransactionType_Invoke,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x17bacc700df6c82682139e8e550078a5daa75dfe356577f78f7e57fd7c56245"),
						internalUtils.TestHexToFelt(t, "0x4eb8734727eb9412b79ba6d14ff1c9a6beb0dc0b811e3f97168c747f8d427b3")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x2d79883d20000",
						},
						L1DataGas: rpc.ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x2d79883d20000",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0x5f5e100",
							MaxPricePerUnit: "0xba43b7400",
						},
					},
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					SenderAddress:         internalUtils.TestHexToFelt(t, "0x745d525a3582e91299d8d7c71730ffc4b1f191f5b219d800334bc0edad0983b"),
					Calldata: internalUtils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x4138fd51f90d171df37e9d4419c8cdb67d525840c58f8a5c347be93a1c5277d",
						"0x2468d193cd15b621b24c2a602b8dbcfa5eaa14f88416c40c09d7fd12592cb4b",
						"0x0",
					}),
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
				},
				ExpectedHash: internalUtils.TestHexToFelt(t, "0x76b52e17bc09064bd986ead34263e6305ef3cecfb3ae9e19b86bf4f1a1a20ea"),
				ExpectedErr:  nil,
			},
		},
	}[testEnv]
	for _, test := range testSet {
		hash, err := acnt.TransactionHashInvoke(test.Txn)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hash.String(), "TransactionHashInvoke not what expected")
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
					Nonce:   internalUtils.TestHexToFelt(t, "0x0"),
					Type:    rpc.TransactionType_DeployAccount,
					MaxFee:  internalUtils.TestHexToFelt(t, "0x1d2109b99cf94"),
					Version: rpc.TransactionV1,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x427df9a1a4a0b7b9011a758524b8a6c2595aac9140608fe24c66efe04b340d7"),
						internalUtils.TestHexToFelt(t, "0x4edc73cd97dab7458a08fec6d7c0e1638c3f1111646fc8a91508b4f94b36310"),
					},
					ClassHash:           internalUtils.TestHexToFelt(t, "0x1e60c8722677cfb7dd8dbea5be86c09265db02cdfe77113e77da7d44c017388"),
					ContractAddressSalt: internalUtils.TestHexToFelt(t, "0x15d621f9515c6197d3117eb1a25c7a4a669317be8f49831e03fcc00d855352e"),
					ConstructorCalldata: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x960532cfba33384bbec41aa669727a9c51e995c87e101c86706aaf244f7e4e"),
					},
				},
				SenderAddress: internalUtils.TestHexToFelt(t, "0x05dd5faeddd4a9e01231f3bb9b95ec93426d08977b721c222e45fd98c5f353ff"),
				ExpectedHash:  internalUtils.TestHexToFelt(t, "0x66d1d9d50d308a9eb16efedbad208b0672769a545a0b828d357757f444e9188"),
				ExpectedErr:   nil,
			},
			{
				// https://sepolia.voyager.online/tx/0x32413f8cee053089d6d7026a72e4108262ca3cfe868dd9159bc1dd160aec975
				Txn: rpc.DeployAccountTxnV3{
					Nonce:   internalUtils.TestHexToFelt(t, "0x0"),
					Type:    rpc.TransactionType_DeployAccount,
					Version: rpc.TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x3ef7f047c95592a04d4d754888dd8f125480a48dee23ee86c115d5da2a86573"),
						internalUtils.TestHexToFelt(t, "0x65e8661ab1526b4f8ea50b76fea1a0e82543de1eb3885e415790d7e1b5a93c7")},
					ResourceBounds: rpc.ResourceBoundsMapping{
						L1Gas: rpc.ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x1597b3274d88",
						},
						L1DataGas: rpc.ResourceBounds{
							MaxAmount:       "0x210",
							MaxPricePerUnit: "0x97c",
						},
						L2Gas: rpc.ResourceBounds{
							MaxAmount:       "0xe6fa0",
							MaxPricePerUnit: "0x1920d1317",
						},
					},
					Tip:           "0x0",
					PayMasterData: []*felt.Felt{},
					NonceDataMode: rpc.DAModeL1,
					FeeMode:       rpc.DAModeL1,
					ClassHash:     internalUtils.TestHexToFelt(t, "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"),
					ConstructorCalldata: internalUtils.TestHexArrToFelt(t, []string{
						"0x2e94ba2293dfa45f86dfcf9952d7a33dc50ce2b00b932999fbe0844772604f3",
					}),
					ContractAddressSalt: internalUtils.TestHexToFelt(t, "0x2e94ba2293dfa45f86dfcf9952d7a33dc50ce2b00b932999fbe0844772604f3"),
				},
				SenderAddress: internalUtils.TestHexToFelt(t, "0x48419d3cc27f158917b45255d5376c06a9524484e19a1102279cbdc715c5522"),
				ExpectedHash:  internalUtils.TestHexToFelt(t, "0x32413f8cee053089d6d7026a72e4108262ca3cfe868dd9159bc1dd160aec975"),
				ExpectedErr:   nil,
			},
		},
	}[testEnv]
	for _, test := range testSet {
		hash, err := acnt.TransactionHashDeployAccount(test.Txn, test.SenderAddress)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hash.String(), "TransactionHashDeployAccount not what expected")
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(test.Timeout)*time.Second)
		defer cancel()

		resp, err := acnt.WaitForTransactionReceipt(ctx, test.Hash, 1*time.Second)
		if test.ExpectedErr != nil {
			rpcErr, ok := err.(*rpc.RPCError)
			require.True(t, ok)
			require.Equal(t, test.ExpectedErr.Code, rpcErr.Code)
			require.Contains(t, rpcErr.Data.ErrorMessage(), test.ExpectedErr.Data.ErrorMessage()) // sometimes the error message starts with "Post \"http://localhost:5050\":..."
		} else {
			require.Equal(t, test.ExpectedReceipt.ExecutionStatus, (*resp).ExecutionStatus)
		}
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

	provider, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

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
	}, 1.5)
	require.NoError(t, err, "Error building and sending invoke txn")

	// check the transaction hash
	require.NotNil(t, resp.TransactionHash)
	t.Logf("Invoke transaction hash: %s", resp.TransactionHash)

	// Waiting for the transaction status (TODO: update this for use WaitForTransactionReceipt when merged with PR 677 that fixed it)
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := provider.GetTransactionStatus(context.Background(), resp.TransactionHash)
	require.NoError(t, err, "Error in provider.GetTransactionStatus")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txStatus.ExecutionStatus)
	assert.Equal(t, rpc.TxnStatus_Accepted_On_L2, txStatus.FinalityStatus)
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

	provider, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Class
	class := *internalUtils.TestUnmarshalJSONFileToType[rpc.ContractClass](t, "./tests/contracts_v2_HelloStarknet.sierra.json", "")

	// Casm Class
	casmClass := *internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](t, "./tests/contracts_v2_HelloStarknet.casm.json", "")

	// Build and send declare txn
	resp, err := acc.BuildAndSendDeclareTxn(context.Background(), casmClass, &class, 1.5)
	if err != nil {
		require.EqualError(t, err, "Transaction execution error: Class with hash 0x0224518978adb773cfd4862a894e9d333192fbd24bc83841dc7d4167c09b89c5 is already declared.")
		t.Log("declare txn not sent: class already declared")
		return
	}

	// check the transaction and class hash
	require.NotNil(t, resp.TransactionHash)
	require.NotNil(t, resp.ClassHash)

	// Waiting for the transaction status (TODO: update this for use WaitForTransactionReceipt when merged with PR 677 that fixed it)
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := provider.GetTransactionStatus(context.Background(), resp.TransactionHash)
	require.NoError(t, err, "Error getting declare transaction status")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txStatus.ExecutionStatus)
	assert.Equal(t, rpc.TxnStatus_Accepted_On_L2, txStatus.FinalityStatus)
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

	provider, err := rpc.NewProvider(base)
	require.NoError(t, err, "Error in rpc.NewClient")

	// we need this account to fund the new account with STRK tokens, in order to deploy it
	acc, err := setupAcc(t, provider)
	require.NoError(t, err, "Error in setupAcc")

	// Get random keys to create the new account
	ks, pub, _ := account.GetRandomKeys()

	// Set up the account passing random values to 'accountAddress' and 'cairoVersion' variables,
	// as for this case we only need the 'ks' to sign the deploy transaction.
	tempAcc, err := account.NewAccount(provider, pub, pub.String(), ks, 2)
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
		1.5)
	require.NoError(t, err, "Error building and estimating deploy account txn")
	require.NotNil(t, deployAccTxn)
	require.NotNil(t, precomputedAddress)
	t.Logf("Precomputed address: %s", precomputedAddress)

	overallFee, err := utils.ResBoundsMapToOverallFee(deployAccTxn.ResourceBounds, 1)
	require.NoError(t, err, "Error converting resource bounds to overall fee")

	// Fund the new account with STRK tokens
	transferSTRKAndWaitConfirmation(t, provider, acc, overallFee, precomputedAddress)

	// Deploy the new account
	resp, err := provider.AddDeployAccountTransaction(context.Background(), deployAccTxn)
	require.NoError(t, err, "Error deploying new account")

	require.NotNil(t, resp.TransactionHash)
	t.Logf("Deploy account transaction hash: %s", resp.TransactionHash)
	require.NotNil(t, resp.ContractAddress)

	// Waiting for the transaction status (TODO: update this for use WaitForTransactionReceipt when merged with PR 677 that fixed it)
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := provider.GetTransactionStatus(context.Background(), resp.TransactionHash)
	require.NoError(t, err, "Error getting deploy account transaction status")

	assert.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txStatus.ExecutionStatus)
	assert.Equal(t, rpc.TxnStatus_Accepted_On_L2, txStatus.FinalityStatus)
}

// a helper function that transfers STRK tokens to a given address and waits for confirmation,
// used to fund the new account with STRK tokens in the TestBuildAndEstimateDeployAccountTxn test
func transferSTRKAndWaitConfirmation(t *testing.T, provider rpc.RpcProvider, acc *account.Account, amount *felt.Felt, recipient *felt.Felt) {
	t.Helper()
	// Build and send invoke txn
	resp, err := acc.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{
		{
			// STRK contract address in Sepolia
			ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
			FunctionName:    "transfer",
			CallData:        []*felt.Felt{recipient, amount, &felt.Zero},
		},
	}, 1.5)
	require.NoError(t, err, "Error transferring STRK tokens")

	// check the transaction hash
	require.NotNil(t, resp.TransactionHash)
	t.Logf("Transfer transaction hash: %s", resp.TransactionHash)
	// Waiting for the transaction status (TODO: update this for use WaitForTransactionReceipt when merged with PR 677 that fixed it)
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := provider.GetTransactionStatus(context.Background(), resp.TransactionHash)
	require.NoError(t, err, "Error getting transfer transaction status")

	require.Equal(t, rpc.TxnExecutionStatusSUCCEEDED, txStatus.ExecutionStatus)
	require.Equal(t, rpc.TxnStatus_Accepted_On_L2, txStatus.FinalityStatus)
}
