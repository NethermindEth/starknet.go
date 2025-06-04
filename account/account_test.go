package account_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/devnet"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/internal"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
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
//   - m: is the test main
//
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
		return nil, errors.New("failed to convert privKey to big.Int")
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
	// TODO: improve test cases to include invoke txns v0 and v3
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

			var acc *account.Account
			var err error
			if testEnv == "testnet" {
				var client *rpc.Provider
				client, err = rpc.NewProvider(base)
				require.NoError(t, err, "Error in rpc.NewClient")
				acc, err = account.NewAccount(client, test.AccountAddress, test.PubKey, ks, 0)
				require.NoError(t, err, "error returned from account.NewAccount()")
			}
			if testEnv == "mock" {
				mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
				// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
				mockRpcProvider.EXPECT().
					ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).
					Return(internalUtils.RANDOM_FELT, nil)
				acc, err = account.NewAccount(mockRpcProvider, test.AccountAddress, test.PubKey, ks, 0)
				require.NoError(t, err, "error returned from account.NewAccount()")
			}
			invokeTxn := rpc.InvokeTxnV1{
				Calldata:      test.FnCall.Calldata,
				Nonce:         test.TxDetails.Nonce,
				MaxFee:        test.TxDetails.MaxFee,
				SenderAddress: acc.Address,
				Version:       test.TxDetails.Version,
			}
			hashResp, err := acc.TransactionHashInvoke(invokeTxn)
			require.NoError(t, err, "error returned from account.TransactionHash()")
			require.Equal(t, test.ExpectedHash.String(), hashResp.String(), "transaction hash does not match expected")

			hash2, err := hash.TransactionHashInvokeV1(&invokeTxn, acc.ChainId)
			require.NoError(t, err)
			assert.Equal(t, hashResp, hash2)
		})
	}
}

// TestFmtCallData tests the FmtCallData function.
//
// It tests the FmtCallData function by providing different test sets
// and comparing the output with the expected call data.
//
// Parameters:
//   - t: The testing.T instance for running the test
//
// Return:
//
//	none
func TestFmtCallData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
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
		var acc *account.Account
		var err error
		if testEnv == "testnet" {
			var client *rpc.Provider
			client, err = rpc.NewProvider(base)
			require.NoError(t, err, "Error in rpc.NewClient")
			acc, err = account.NewAccount(client, &felt.Zero, "pubkey", account.NewMemKeystore(), test.CairoVersion)
			require.NoError(t, err)
		}
		if testEnv == "mock" {
			mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
			acc, err = account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore(), test.CairoVersion)
			require.NoError(t, err)
		}

		fmtCallData, err := acc.FmtCalldata([]rpc.FunctionCall{test.FnCall})
		require.NoError(t, err)
		require.Equal(t, fmtCallData, test.ExpectedCallData)
	}
}

// TestChainIdMOCK is a test function that tests the behaviour of the ChainId function.
//
// It creates a mock controller and a mock RpcProvider. It defines a test set
// consisting of different ChainID and ExpectedID pairs. It then iterates over
// the test set and sets the expected behaviour for the ChainID method of the
// mockRpcProvider. It creates a new account using the mockRpcProvider,
// Zero value, "pubkey", and a new in-memory keystore. It asserts that the
// account's ChainId matches the expected ID for each test case in the test set.
//
// Parameters:
//   - t: The testing.T instance for running the test
//
// Return:
//
//	none
func TestChainIdMOCK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
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
		// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
		mockRpcProvider.EXPECT().ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).Return(internalUtils.RANDOM_FELT, nil)
		acc, err := account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore(), 0)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedID, acc.ChainId.String())
	}
}

// TestChainId tests the ChainId function.
//
// This function tests the ChainId function by setting up a mock controller, defining a test set,
// and running a series of assertions on the expected results.
// It checks if the ChainId function returns the correct ChainID and ExpectedID values
// for different test environments.
// Parameters:
//   - t: The testing.T instance for running the test
//
// Return:
//
//	none
func TestChainId(t *testing.T) {
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

		acc, err := account.NewAccount(client, &felt.Zero, "pubkey", account.NewMemKeystore(), 0)
		require.NoError(t, err)
		require.Equal(t, acc.ChainId.String(), test.ExpectedID)
	}
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
//   - t: reference to the testing.T object
//
// Returns:
//
//	none
func TestTransactionHashDeclare(t *testing.T) {
	var acnt *account.Account
	var err error
	if testEnv == "mock" {
		mockCtrl := gomock.NewController(t)

		mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)
		// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
		mockRpcProvider.EXPECT().ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).Return(internalUtils.RANDOM_FELT, nil)
		acnt, err = account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
		require.NoError(t, err)
	}
	if testEnv == "testnet" {
		client, err := rpc.NewProvider(base)
		require.NoError(t, err, "Error in rpc.NewClient")
		acnt, err = account.NewAccount(client, &felt.Zero, "", account.NewMemKeystore(), 0)
		require.NoError(t, err)
	}

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
						internalUtils.TestHexToFelt(t, "0x4f28b1c15379c0ceb1855c09ed793e7583f875a802cbf310a8c0c971835c5cf"),
					},
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
					ResourceBounds: &rpc.ResourceBoundsMapping{
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
						internalUtils.TestHexToFelt(t, "0x4f28b1c15379c0ceb1855c09ed793e7583f875a802cbf310a8c0c971835c5cf"),
					},
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
		hashResp, err := acnt.TransactionHashDeclare(test.Txn)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hashResp.String(), "TransactionHashDeclare not what expected")

		var hash2 *felt.Felt
		switch txn := test.Txn.(type) {
		case rpc.DeclareTxnV2:
			hash2, err = hash.TransactionHashDeclareV2(&txn, acnt.ChainId)
		case rpc.DeclareTxnV3:
			hash2, err = hash.TransactionHashDeclareV3(&txn, acnt.ChainId)
		}
		require.NoError(t, err)
		assert.Equal(t, hashResp, hash2)
	}
}

func TestTransactionHashInvokeV3(t *testing.T) {
	var acnt *account.Account
	var err error
	if testEnv == "mock" {
		mockCtrl := gomock.NewController(t)

		mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)
		// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
		mockRpcProvider.EXPECT().ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).Return(internalUtils.RANDOM_FELT, nil)
		acnt, err = account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
		require.NoError(t, err)
	}
	if testEnv == "testnet" {
		client, err := rpc.NewProvider(base)
		require.NoError(t, err, "Error in rpc.NewClient")
		acnt, err = account.NewAccount(client, &felt.Zero, "", account.NewMemKeystore(), 0)
		require.NoError(t, err)
	}

	type testSetType struct {
		Txn          rpc.InvokeTxnV3
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
						internalUtils.TestHexToFelt(t, "0x4eb8734727eb9412b79ba6d14ff1c9a6beb0dc0b811e3f97168c747f8d427b3"),
					},
					ResourceBounds: &rpc.ResourceBoundsMapping{
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
		hashResp, err := acnt.TransactionHashInvoke(test.Txn)
		require.Equal(t, test.ExpectedErr, err)
		require.Equal(t, test.ExpectedHash.String(), hashResp.String(), "TransactionHashInvoke not what expected")

		hash2, err := hash.TransactionHashInvokeV3(&test.Txn, acnt.ChainId)
		require.NoError(t, err)
		assert.Equal(t, hashResp, hash2)
	}
}

func TestTransactionHashdeployAccount(t *testing.T) {
	var acnt *account.Account
	var err error
	if testEnv == "mock" {
		mockCtrl := gomock.NewController(t)

		mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)
		// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
		mockRpcProvider.EXPECT().ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).Return(internalUtils.RANDOM_FELT, nil)

		acnt, err = account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore(), 0)
		require.NoError(t, err)
	}
	if testEnv == "testnet" {
		client, err := rpc.NewProvider(base)
		require.NoError(t, err, "Error in rpc.NewClient")
		acnt, err = account.NewAccount(client, &felt.Zero, "", account.NewMemKeystore(), 0)
		require.NoError(t, err)
	}
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
				Txn: rpc.DeployAccountTxnV1{
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
						internalUtils.TestHexToFelt(t, "0x65e8661ab1526b4f8ea50b76fea1a0e82543de1eb3885e415790d7e1b5a93c7"),
					},
					ResourceBounds: &rpc.ResourceBoundsMapping{
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
		hashResp, err := acnt.TransactionHashDeployAccount(test.Txn, test.SenderAddress)
		require.Equal(t, test.ExpectedErr, err)
		assert.Equal(t, test.ExpectedHash.String(), hashResp.String(), "TransactionHashDeployAccount not what expected")

		var hash2 *felt.Felt
		switch txn := test.Txn.(type) {
		case rpc.DeployAccountTxnV1:
			hash2, err = hash.TransactionHashDeployAccountV1(&txn, test.SenderAddress, acnt.ChainId)
		case rpc.DeployAccountTxnV3:
			hash2, err = hash.TransactionHashDeployAccountV3(&txn, test.SenderAddress, acnt.ChainId)
		}
		require.NoError(t, err)
		assert.Equal(t, hashResp, hash2)
	}
}

// newDevnet creates a new devnet with the given URL.
//
// Parameters:
//   - t: The testing.T instance for running the test
//   - url: The URL of the devnet to be created
//
// Returns:
//   - *devnet.DevNet: a pointer to a devnet object
//   - []devnet.TestAccount: a slice of test accounts
//   - error: an error, if any
func newDevnet(t *testing.T, url string) (*devnet.DevNet, []devnet.TestAccount, error) {
	t.Helper()
	devnetInstance := devnet.NewDevNet(url)
	acnts, err := devnetInstance.Accounts()

	return devnetInstance, acnts, err
}

// newDevnetAccount creates a new devnet account from a test account.
//
// Parameters:
//   - t: The testing.T instance for running the test
//   - provider: The RPC provider
//   - accData: The test account data
//
// Returns:
//   - *account.Account: The new devnet account
//   - error: An error, if any
func newDevnetAccount(t *testing.T, provider *rpc.Provider, accData devnet.TestAccount, cairoVersion int) *account.Account {
	t.Helper()
	fakeUserAddr := internalUtils.TestHexToFelt(t, accData.Address)
	fakeUserPriv := internalUtils.TestHexToFelt(t, accData.PrivateKey)

	// Set up ks
	ks := account.NewMemKeystore()
	ks.Put(accData.PublicKey, fakeUserPriv.BigInt(new(big.Int)))

	acnt, err := account.NewAccount(provider, fakeUserAddr, accData.PublicKey, ks, cairoVersion)
	require.NoError(t, err)

	return acnt
}

func TestBraavosAccountWarning(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	type testSetType struct {
		ClassHash      *felt.Felt
		ExpectedOutput bool
	}

	// Known Braavos class hashes
	braavosClassHashes := []string{
		"0x2c8c7e6fbcfb3e8e15a46648e8914c6aa1fc506fc1e7fb3d1e19630716174bc",
		"0x816dd0297efc55dc1e7559020a3a825e81ef734b558f03c83325d4da7e6253",
		"0x41bf1e71792aecb9df3e9d04e1540091c5e13122a731e02bec588f71dc1a5c3",
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				ClassHash:      internalUtils.TestHexToFelt(t, braavosClassHashes[0]),
				ExpectedOutput: true,
			},
			{
				ClassHash:      internalUtils.TestHexToFelt(t, braavosClassHashes[1]),
				ExpectedOutput: true,
			},
			{
				ClassHash:      internalUtils.TestHexToFelt(t, braavosClassHashes[2]),
				ExpectedOutput: true,
			},
			{
				ClassHash:      internalUtils.RANDOM_FELT,
				ExpectedOutput: false,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		t.Run("ClassHash_"+test.ClassHash.String(), func(t *testing.T) {
			// Set up the mock to return the Braavos class hash
			mockRpcProvider.EXPECT().ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).Return(test.ClassHash, nil)
			mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_SEPOLIA", nil)

			// Create a buffer to capture stdout
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			require.NoError(t, err)
			os.Stdout = w

			// Create the account
			_, err = account.NewAccount(mockRpcProvider, internalUtils.RANDOM_FELT, "pubkey", account.NewMemKeystore(), 2)
			require.NoError(t, err)

			// Close the writer and restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the captured output
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			require.NoError(t, err)

			if test.ExpectedOutput {
				// Check if the warning message was printed
				assert.Contains(t, buf.String(), account.BRAAVOS_WARNING_MESSAGE)
			} else {
				// Check if the warning message was not printed
				assert.NotContains(t, buf.String(), account.BRAAVOS_WARNING_MESSAGE)
			}
		})
	}
}
