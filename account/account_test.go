package account_test

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/test-go/testify/require"
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
//  none
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	base = os.Getenv("INTEGRATION_BASE")
	if base == "" && testEnv != "mock" {
		panic(fmt.Sprint("Failed to set INTEGRATION_BASE for ", testEnv))
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
// Returns:
//   none
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
				// https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8
				ExpectedHash:   utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8"),
				SetKS:          true,
				AccountAddress: utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
				PrivKey:        utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				PubKey:         "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e",
				ChainID:        "SN_GOERLI",
				FnCall: rpc.FunctionCall{
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
						"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
						"0x0",
						"0x3",
						"0x3",
						"0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
						"0x1",
						"0x0",
					}),
				},
				TxDetails: rpc.TxDetails{
					Nonce:   utils.TestHexToFelt(t, "0x2"),
					MaxFee:  utils.TestHexToFelt(t, "0x574fbde6000"),
					Version: rpc.TransactionV1,
				},
			},
			{
				ExpectedHash:   utils.TestHexToFelt(t, "0x135c34f53f8b7f59efd450eb689fccd9dd4cfe7f9d9dc4d09954c5653138698"),
				SetKS:          false,
				AccountAddress: &felt.Zero,
				ChainID:        "SN_GOERLI",
				FnCall: rpc.FunctionCall{
					ContractAddress:    &felt.Zero,
					EntryPointSelector: &felt.Zero,
					Calldata:           []*felt.Felt{&felt.Zero},
				},
				TxDetails: rpc.TxDetails{
					Nonce:   &felt.Zero,
					MaxFee:  &felt.Zero,
					Version: rpc.TransactionV1,
				},
			},
			{
				ExpectedHash:   utils.TestHexToFelt(t, "0x3476c76a81522fe52616c41e95d062f5c3ea4eeb6c652904ad389fcd9ff4637"),
				SetKS:          false,
				AccountAddress: utils.TestHexToFelt(t, "0x59cd166e363be0a921e42dd5cfca0049aedcf2093a707ef90b5c6e46d4555a8"),
				ChainID:        "SN_MAIN",
				FnCall: rpc.FunctionCall{
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x1",
						"0x5dbdedc203e92749e2e746e2d40a768d966bd243df04a6b712e222bc040a9af",
						"0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
						"0x0",
						"0x1",
						"0x1",
						"0x52884ee3f",
					}),
				},
				TxDetails: rpc.TxDetails{
					Nonce:   utils.TestHexToFelt(t, "0x1"),
					MaxFee:  utils.TestHexToFelt(t, "0x2a173cd36e400"),
					Version: rpc.TransactionV1,
				},
			},
		},
		"devnet":  {},
		"testnet": {},
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
			account, err := account.NewAccount(mockRpcProvider, test.AccountAddress, test.PubKey, ks)
			require.NoError(t, err, "error returned from account.NewAccount()")
			invokeTxn := rpc.InvokeTxnV1{
				Calldata:      test.FnCall.Calldata,
				Nonce:         test.TxDetails.Nonce,
				MaxFee:        test.TxDetails.MaxFee,
				SenderAddress: account.AccountAddress,
				Version:       test.TxDetails.Version,
			}
			hash, err := account.TransactionHashInvoke(invokeTxn)
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
//   none
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
		"mock": {
			{
				CairoVersion: 0,
				ChainID:      "SN_GOERLI",
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("transfer"),
					Calldata: utils.TestHexArrToFelt(t, []string{
						"0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
						"0x1"}),
				},
				ExpectedCallData: utils.TestHexArrToFelt(t, []string{
					"0x1",
					"0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
					"0x0",
					"0x3",
					"0x3",
					"0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					"0x1",
					"0x0",
				}),
			},
		},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
		acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore())
		require.NoError(t, err)

		fmtCallData, err := acnt.FmtCalldata([]rpc.FunctionCall{test.FnCall}, test.CairoVersion)
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
//   none
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
				ChainID:    "SN_GOERLI",
				ExpectedID: "0x534e5f474f45524c49",
			},
		},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
		account, err := account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore())
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), test.ExpectedID)
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
//   none
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
				ChainID:    "SN_GOERLI",
				ExpectedID: "0x534e5f474f45524c49",
			},
		},
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		client, err := rpc.NewClient(base + "/rpc")
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		account, err := account.NewAccount(provider, &felt.Zero, "pubkey", account.NewMemKeystore())
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
//  none
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
			// Accepted on testnet https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8
			{
				Address:    utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
				PrivKey:    utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				ChainId:    "SN_GOERLI",
				FeltToSign: utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8"),
				ExpectedSig: []*felt.Felt{
					utils.TestHexToFelt(t, "0x10d405427040655f118bc8b897e2f2f8147858bbcb0e3d6bc6dfbc6d0205e8"),
					utils.TestHexToFelt(t, "0x5cdfe4a3d5b63002e9011ec0ba59ae2b75a43cb2a3bc1699b35aa64cb9ca3cf"),
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
		account, err := account.NewAccount(mockRpcProvider, test.Address, test.Address.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		msg := utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8")
		sig, err := account.Sign(context.Background(), msg)

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
//   none
func TestAddInvoke(t *testing.T) {

	type testSetType struct {
		ExpectedError        *rpc.RPCError
		CairoContractVersion int
		SetKS                bool
		AccountAddress       *felt.Felt
		PubKey               *felt.Felt
		PrivKey              *felt.Felt
		InvokeTx             rpc.InvokeTxnV1
		FnCall               rpc.FunctionCall
		TxDetails            rpc.TxDetails
	}
	testSet := map[string][]testSetType{
		"mock":   {},
		"devnet": {},
		"testnet": {
			{
				// https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8#overview
				ExpectedError:        rpc.ErrDuplicateTx,
				CairoContractVersion: 0,
				AccountAddress:       utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x049f060d2dffd3bf6f2c103b710baf519530df44529045f92c3903097e8d861f"),
				PrivKey:              utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				InvokeTx: rpc.InvokeTxnV1{
					Nonce:         new(felt.Felt).SetUint64(2),
					MaxFee:        utils.TestHexToFelt(t, "0x574fbde6000"),
					Version:       rpc.TransactionV1,
					Type:          rpc.TransactionType_Invoke,
					SenderAddress: utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
				},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("transfer"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
						utils.TestHexToFelt(t, "0x1"),
					},
				},
			},
			{
				// https://goerli.voyager.online/tx/0x171537c58b16db45aeec3d3f493617cd3dd571561b856c115dc425b85212c86#overview
				ExpectedError:        rpc.ErrDuplicateTx,
				CairoContractVersion: 0,
				AccountAddress:       utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x049f060d2dffd3bf6f2c103b710baf519530df44529045f92c3903097e8d861f"),
				PrivKey:              utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa"),
				InvokeTx: rpc.InvokeTxnV1{
					Nonce:         new(felt.Felt).SetUint64(6),
					MaxFee:        utils.TestHexToFelt(t, "0x9184e72a000"),
					Version:       rpc.TransactionV1,
					Type:          rpc.TransactionType_Invoke,
					SenderAddress: utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
				},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x03E85bFbb8E2A42B7BeaD9E88e9A1B19dbCcf661471061807292120462396ec9"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("burn"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e"),
						utils.TestHexToFelt(t, "0x1"),
					},
				},
			},
			{
				// https://goerli.voyager.online/tx/0x1bc0f8c04584735ea9e4485f927c25a6e025bda3117beb508cd1bb5e41f08d9
				ExpectedError:        rpc.ErrDuplicateTx,
				CairoContractVersion: 2,
				AccountAddress:       utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883"),
				PrivKey:              utils.TestHexToFelt(t, "0x07514c4f0de1f800b0b0c7377ef39294ce218a7abd9a1c9b6aa574779f7cdc6a"),
				InvokeTx: rpc.InvokeTxnV1{
					Nonce:         new(felt.Felt).SetUint64(6),
					MaxFee:        utils.TestHexToFelt(t, "0x9184e72a000"),
					Version:       rpc.TransactionV1,
					Type:          rpc.TransactionType_Invoke,
					SenderAddress: utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1"),
				},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x05044dfb70b9475663e3ddddb11bbbeccc71614b8db86fc3dc0c16b2b9d3151d"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("increase_value_8"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x1234"),
					},
				},
			},
			{
				// https://goerli.voyager.online/tx/0xe8cdb03ddc6b65c2c268eb8084bef41ef63009c10a38f8d1e167652a721588
				ExpectedError:        rpc.ErrDuplicateTx,
				CairoContractVersion: 2,
				AccountAddress:       utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883"),
				PrivKey:              utils.TestHexToFelt(t, "0x07514c4f0de1f800b0b0c7377ef39294ce218a7abd9a1c9b6aa574779f7cdc6a"),
				InvokeTx: rpc.InvokeTxnV1{
					Nonce:         new(felt.Felt).SetUint64(7),
					MaxFee:        utils.TestHexToFelt(t, "0x9184e72a000"),
					Version:       rpc.TransactionV1,
					Type:          rpc.TransactionType_Invoke,
					SenderAddress: utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1"),
				},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x05044dfb70b9475663e3ddddb11bbbeccc71614b8db86fc3dc0c16b2b9d3151d"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("increase_value"),
					Calldata:           []*felt.Felt{},
				},
			},
			{
				// https://goerli.voyager.online/tx/0xdcec9fdd48440243fa8fdb8bf87cc40d5ef91181d5a4a0304140df5701c238
				ExpectedError:        rpc.ErrDuplicateTx,
				CairoContractVersion: 2,
				AccountAddress:       utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1"),
				SetKS:                true,
				PubKey:               utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883"),
				PrivKey:              utils.TestHexToFelt(t, "0x07514c4f0de1f800b0b0c7377ef39294ce218a7abd9a1c9b6aa574779f7cdc6a"),
				InvokeTx: rpc.InvokeTxnV1{
					Nonce:         new(felt.Felt).SetUint64(18),
					MaxFee:        utils.TestHexToFelt(t, "0x9184e72a000"),
					Version:       rpc.TransactionV1,
					Type:          rpc.TransactionType_Invoke,
					SenderAddress: utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1"),
				},
				FnCall: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, "0x05044dfb70b9475663e3ddddb11bbbeccc71614b8db86fc3dc0c16b2b9d3151d"),
					EntryPointSelector: utils.GetSelectorFromNameFelt("increase_value_8"),
					Calldata:           []*felt.Felt{utils.TestHexToFelt(t, "0xaC25b2B9F4ca06179fA0D2522F47Bc86A9DF9314")},
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		client, err := rpc.NewClient(base)
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		// Set up ks
		ks := account.NewMemKeystore()
		if test.SetKS {
			fakePrivKeyBI, ok := new(big.Int).SetString(test.PrivKey.String(), 0)
			require.True(t, ok)
			ks.Put(test.PubKey.String(), fakePrivKeyBI)
		}

		acnt, err := account.NewAccount(provider, test.AccountAddress, test.PubKey.String(), ks)
		require.NoError(t, err)

		test.InvokeTx.Calldata, err = acnt.FmtCalldata([]rpc.FunctionCall{test.FnCall}, test.CairoContractVersion)
		require.NoError(t, err)

		err = acnt.SignInvokeTransaction(context.Background(), &test.InvokeTx)
		require.NoError(t, err)

		resp, err := acnt.AddInvokeTransaction(context.Background(), test.InvokeTx)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError.Error())
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
//  - t: is the testing framework
// Returns:
//  none
func TestAddDeployAccountDevnet(t *testing.T) {
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}
	client, err := rpc.NewClient(base + "/rpc")
	require.NoError(t, err, "Error in rpc.NewClient")
	provider := rpc.NewProvider(client)

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

	acnt, err := account.NewAccount(provider, fakeUserAddr, fakeUser.PublicKey, ks)
	require.NoError(t, err)

	classHash := utils.TestHexToFelt(t, "0x7b3e05f48f0c69e4a65ce5e076a66271a527aff2c34ce1083ec6e1526997a69") // preDeployed classhash
	require.NoError(t, err)

	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		MaxFee:              new(felt.Felt).SetUint64(4724395326064),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: fakeUserPub,
		ConstructorCalldata: []*felt.Felt{fakeUserPub},
	}

	precomputedAddress, err := acnt.PrecomputeAddress(&felt.Zero, fakeUserPub, classHash, tx.ConstructorCalldata)
	require.NoError(t, acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress))

	_, err = devnet.Mint(precomputedAddress, new(big.Int).SetUint64(10000000000000000000))
	require.NoError(t, err)

	resp, err := acnt.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
	require.NoError(t, err, "AddDeployAccountTransaction gave an Error")
	require.NotNil(t, resp, "AddDeployAccountTransaction resp not nil")
}

// TestTransactionHashDeployAccountTestnet tests the TransactionHashDeployAccount function when using the testnet environment.
//
// It creates a client and provider, initializes the required addresses and keys, and sets up the transaction parameters.
// It then precomputes the address and calculates the hash of the transaction.
// Finally, it verifies that the calculated hash matches the expected hash.
//
// Parameters:
//  - t: is the testing framework
// Returns:
//  none
func TestTransactionHashDeployAccountTestnet(t *testing.T) {

	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}

	client, err := rpc.NewClient(base)
	require.NoError(t, err, "Error in rpc.NewClient")
	provider := rpc.NewProvider(client)

	AccountAddress := utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1")
	PubKey := utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883")
	PrivKey := utils.TestHexToFelt(t, "0x07514c4f0de1f800b0b0c7377ef39294ce218a7abd9a1c9b6aa574779f7cdc6a")

	ExpectedHash := utils.TestHexToFelt(t, "0x5b6b5927cd70ad7a80efdbe898244525871875c76540b239f6730118598b9cb")
	ExpectedPrecomputeAddr := utils.TestHexToFelt(t, "0x88d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1")
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(PrivKey.String(), 0)
	require.True(t, ok)
	ks.Put(PubKey.String(), fakePrivKeyBI)

	acnt, err := account.NewAccount(provider, AccountAddress, PubKey.String(), ks)
	require.NoError(t, err)

	classHash := utils.TestHexToFelt(t, "0x3131fa018d520a037686ce3efddeab8f28895662f019ca3ca18a626650f7d1e")

	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero,
		MaxFee:              utils.TestHexToFelt(t, "0x105ef39b2000"),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883"),
		ConstructorCalldata: []*felt.Felt{
			utils.TestHexToFelt(t, "0x5aa23d5bb71ddaa783da7ea79d405315bafa7cf0387a74f4593578c3e9e6570"),
			utils.TestHexToFelt(t, "0x2dd76e7ad84dbed81c314ffe5e7a7cacfb8f4836f01af4e913f275f89a3de1a"),
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883"),
		},
	}
	precomputedAddress, err := acnt.PrecomputeAddress(&felt.Zero, tx.ContractAddressSalt, classHash, tx.ConstructorCalldata)
	require.Equal(t, ExpectedPrecomputeAddr.String(), precomputedAddress.String(), "Error with calulcating PrecomputeAddress")

	hash, err := acnt.TransactionHashDeployAccount(tx, precomputedAddress)
	require.NoError(t, err, "TransactionHashDeployAccount gave an Error")
	require.Equal(t, hash.String(), ExpectedHash.String(), "Error with calulcating TransactionHashDeployAccount")
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
//  none
func TestTransactionHashDeclare(t *testing.T) {
	// https://goerli.voyager.online/tx/0x4e0519272438a3ae0d0fca776136e2bb6fcd5d3b2af47e53575c5874ccfce92
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}
	expectedHash := utils.TestHexToFelt(t, "0x4e0519272438a3ae0d0fca776136e2bb6fcd5d3b2af47e53575c5874ccfce92")

	client, err := rpc.NewClient(base)
	require.NoError(t, err, "Error in rpc.NewClient")
	provider := rpc.NewProvider(client)

	acnt, err := account.NewAccount(provider, &felt.Zero, "", account.NewMemKeystore())
	require.NoError(t, err)

	tx := rpc.DeclareTxnV2{
		Nonce:             utils.TestHexToFelt(t, "0xb"),
		MaxFee:            utils.TestHexToFelt(t, "0x50c8f3053db"),
		Type:              rpc.TransactionType_Declare,
		Version:           rpc.TransactionV2,
		Signature:         []*felt.Felt{},
		SenderAddress:     utils.TestHexToFelt(t, "0x36437dffa1b0bf630f04690a3b302adbabb942deb488ea430660c895ff25acf"),
		CompiledClassHash: utils.TestHexToFelt(t, "0x615a5260d3d47d79fba87898da95cb5394b181c7d5097bc8ced4ed06ac24ac5"),
		ClassHash:         utils.TestHexToFelt(t, "0x639cdc0c42c8c4d3d805e56294fa0e6bf5a584ad0fcd538b843cc294913b982"),
	}

	hash, err := acnt.TransactionHashDeclare(tx)
	require.NoError(t, err)
	require.Equal(t, expectedHash.String(), hash.String(), "TransactionHashDeclare not what expected")
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
//  none
func TestWaitForTransactionReceiptMOCK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
	acnt, err := account.NewAccount(mockRpcProvider, &felt.Zero, "", account.NewMemKeystore())
	require.NoError(t, err, "error returned from account.NewAccount()")

	type testSetType struct {
		Timeout                      time.Duration
		ShouldCallTransactionReceipt bool
		Hash                         *felt.Felt
		ExpectedErr                  error
		ExpectedReceipt              rpc.TransactionReceipt
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				Timeout:                      time.Duration(1000),
				ShouldCallTransactionReceipt: true,
				Hash:                         new(felt.Felt).SetUint64(1),
				ExpectedReceipt:              nil,
				ExpectedErr:                  errors.New("UnExpectedErr"),
			},
			{
				Timeout:                      time.Duration(1000),
				Hash:                         new(felt.Felt).SetUint64(2),
				ShouldCallTransactionReceipt: true,
				ExpectedReceipt: rpc.InvokeTransactionReceipt{
					TransactionHash: new(felt.Felt).SetUint64(2),
					ExecutionStatus: rpc.TxnExecutionStatusSUCCEEDED,
				},
				ExpectedErr: nil,
			},
			{
				Timeout:                      time.Duration(1),
				Hash:                         new(felt.Felt).SetUint64(3),
				ShouldCallTransactionReceipt: false,
				ExpectedReceipt:              nil,
				ExpectedErr:                  context.DeadlineExceeded,
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
			require.Equal(t, test.ExpectedReceipt.GetExecutionStatus(), (*resp).GetExecutionStatus())
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
//  none
func TestWaitForTransactionReceipt(t *testing.T) {
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}
	client, err := rpc.NewClient(base + "/rpc")
	require.NoError(t, err, "Error in rpc.NewClient")
	provider := rpc.NewProvider(client)

	acnt, err := account.NewAccount(provider, &felt.Zero, "pubkey", account.NewMemKeystore())
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
				ExpectedReceipt: nil,
				ExpectedErr:     errors.New("Post \"http://0.0.0.0:5050/rpc\": context deadline exceeded"),
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
			require.Equal(t, test.ExpectedReceipt.GetExecutionStatus(), (*resp).GetExecutionStatus())
		}

	}
}

// TestAddDeclareTxn is a test function that verifies the behavior of the AddDeclareTransaction method.
//
// This function tests the AddDeclareTransaction method by setting up test data and invoking the method with different test sets.
// It asserts that the expected hash and error values are returned for each test set.
//
// Parameters:
//  - t: The testing.T instance for running the test
// Returns:
//  none
func TestAddDeclareTxn(t *testing.T) {
	// https://goerli.voyager.online/tx/0x76af2faec46130ffad1ab2f615ad16b30afcf49cfbd09f655a26e545b03a21d
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}
	expectedTxHash := utils.TestHexToFelt(t, "0x76af2faec46130ffad1ab2f615ad16b30afcf49cfbd09f655a26e545b03a21d")
	expectedClassHash := utils.TestHexToFelt(t, "0x76af2faec46130ffad1ab2f615ad16b30afcf49cfbd09f655a26e545b03a21d")

	AccountAddress := utils.TestHexToFelt(t, "0x0088d0038623a89bf853c70ea68b1062ccf32b094d1d7e5f924cda8404dc73e1")
	PubKey := utils.TestHexToFelt(t, "0x7ed3c6482e12c3ef7351214d1195ee7406d814af04a305617599ff27be43883")
	PrivKey := utils.TestHexToFelt(t, "0x07514c4f0de1f800b0b0c7377ef39294ce218a7abd9a1c9b6aa574779f7cdc6a")

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(PrivKey.String(), 0)
	require.True(t, ok)
	ks.Put(PubKey.String(), fakePrivKeyBI)

	client, err := rpc.NewClient(base)
	require.NoError(t, err, "Error in rpc.NewClient")
	provider := rpc.NewProvider(client)

	acnt, err := account.NewAccount(provider, AccountAddress, PubKey.String(), ks)
	require.NoError(t, err)

	// Class Hash
	content, err := os.ReadFile("./tests/hello_starknet_compiled.sierra.json")
	require.NoError(t, err)

	var class rpc.ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
	classHash, err := hash.ClassHash(class)
	require.NoError(t, err)

	// Compiled Class Hash
	content2, err := os.ReadFile("./tests/hello_starknet_compiled.sierra.json")
	require.NoError(t, err)

	var casmClass contracts.CasmClass
	err = json.Unmarshal(content2, &casmClass)
	require.NoError(t, err)
	compClassHash := hash.CompiledClassHash(casmClass)

	nonce, err := acnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, acnt.AccountAddress)
	require.NoError(t, err)

	tx := rpc.DeclareTxnV2{
		Nonce:             nonce,
		MaxFee:            utils.TestHexToFelt(t, "0x50c8f3053db"),
		Type:              rpc.TransactionType_Declare,
		Version:           rpc.TransactionV2,
		Signature:         []*felt.Felt{},
		SenderAddress:     AccountAddress,
		CompiledClassHash: compClassHash,
		ClassHash:         classHash,
	}

	err = acnt.SignDeclareTransaction(context.Background(), &tx)
	require.NoError(t, err)

	resp, err := acnt.AddDeclareTransaction(context.Background(), tx)

	if err != nil {
		require.Equal(t, err.Error(), rpc.ErrDuplicateTx.Error())
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
	devnet := devnet.NewDevNet(url)
	acnts, err := devnet.Accounts()
	return devnet, acnts, err
}
