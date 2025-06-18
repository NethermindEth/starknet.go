package account_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
		CairoVersion     account.CairoVersion
		ChainID          string
		FnCall           rpc.FunctionCall
		ExpectedCallData []*felt.Felt
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mock":   {},
		"testnet": {
			{
				CairoVersion: account.CairoV2,
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
				CairoVersion: account.CairoV2,
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
			client, err = rpc.NewProvider(tConfig.providerURL)
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
		acc, err := account.NewAccount(mockRpcProvider, &felt.Zero, "pubkey", account.NewMemKeystore(), account.CairoV0)
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
		client, err := rpc.NewProvider(tConfig.providerURL)
		require.NoError(t, err, "Error in rpc.NewClient")

		acc, err := account.NewAccount(client, &felt.Zero, "pubkey", account.NewMemKeystore(), account.CairoV0)
		require.NoError(t, err)
		require.Equal(t, acc.ChainId.String(), test.ExpectedID)
	}
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
			_, err = account.NewAccount(mockRpcProvider, internalUtils.RANDOM_FELT, "pubkey", account.NewMemKeystore(), account.CairoV2)
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
