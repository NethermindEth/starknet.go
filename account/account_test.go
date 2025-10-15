package account_test

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/internal/tests"
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
	tests.RunTestOn(t, tests.MockEnv)

	mockCtrl := gomock.NewController(t)
	mockRPCProvider := mocks.NewMockRPCProvider(mockCtrl)

	type testSetType struct {
		CairoVersion     account.CairoVersion
		ChainID          string
		FnCall           rpc.FunctionCall
		ExpectedCallData []*felt.Felt
	}
	testSet := []testSetType{
		{
			CairoVersion: account.CairoV2,
			ChainID:      "SN_SEPOLIA",
			FnCall: rpc.FunctionCall{
				ContractAddress: internalUtils.TestHexToFelt(
					t,
					"0x04daadb9d30c887e1ab2cf7d78dfe444a77aab5a49c3353d6d9977e7ed669902",
				),
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
				ContractAddress: internalUtils.TestHexToFelt(
					t,
					"0x017cE9DffA7C87a03EB496c96e04ac36c4902085030763A83a35788d475e15CA",
				),
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
	}

	for _, test := range testSet {
		mockRPCProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
		acc, err := account.NewAccount(
			mockRPCProvider,
			&felt.Zero,
			"pubkey",
			account.NewMemKeystore(),
			test.CairoVersion,
		)
		require.NoError(t, err)

		fmtCallData, err := acc.FmtCalldata([]rpc.FunctionCall{test.FnCall})
		require.NoError(t, err)
		assert.Equal(t, fmtCallData, test.ExpectedCallData)
	}
}

// TestChainIdMOCK is a test function that tests the behaviour of the ChainId function.
//
// It creates a mock controller and a mock RpcProvider. It defines a test set
// consisting of different ChainID and ExpectedID pairs. It then iterates over
// the test set and sets the expected behaviour for the ChainID method of the
// mockRPCProvider. It creates a new account using the mockRPCProvider,
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
	tests.RunTestOn(t, tests.MockEnv)

	mockCtrl := gomock.NewController(t)
	mockRPCProvider := mocks.NewMockRPCProvider(mockCtrl)

	type testSetType struct {
		ChainID    string
		ExpectedID string
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				ChainID:    "SN_MAIN",
				ExpectedID: "0x534e5f4d41494e",
			},
			{
				ChainID:    "SN_SEPOLIA",
				ExpectedID: "0x534e5f5345504f4c4941",
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		mockRPCProvider.EXPECT().ChainID(context.Background()).Return(test.ChainID, nil)
		acc, err := account.NewAccount(
			mockRPCProvider,
			&felt.Zero,
			"pubkey",
			account.NewMemKeystore(),
			account.CairoV0,
		)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedID, acc.ChainID.String())
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
	tests.RunTestOn(t, tests.DevnetEnv)

	type testSetType struct {
		ChainID    string
		ExpectedID string
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.DevnetEnv: {
			{
				ChainID:    "SN_SEPOLIA",
				ExpectedID: "0x534e5f5345504f4c4941",
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		client, err := rpc.NewProvider(tConfig.providerURL)
		require.NoError(t, err, "Error in rpc.NewClient")

		acc, err := account.NewAccount(
			client,
			&felt.Zero,
			"pubkey",
			account.NewMemKeystore(),
			account.CairoV0,
		)
		require.NoError(t, err)
		require.Equal(t, acc.ChainID.String(), test.ExpectedID)
	}
}
