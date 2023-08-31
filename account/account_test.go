package account_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/require"
)

func TestTransactionHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	t.Run("Transaction hash", func(t *testing.T) {
		expectedHash, _ := new(felt.Felt).SetString("0x135c34f53f8b7f59efd450eb689fccd9dd4cfe7f9d9dc4d09954c5653138698")
		address := &felt.Zero

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, starknetgo.NewMemKeystore(), "")
		require.NoError(t, err, "error returned from account.NewAccount()")

		call := rpc.FunctionCall{
			ContractAddress:    &felt.Zero,
			EntryPointSelector: &felt.Zero,
			Calldata:           []*felt.Felt{&felt.Zero},
		}
		txDetails := rpc.TxDetails{
			Nonce:  &felt.Zero,
			MaxFee: &felt.Zero,
		}
		hash, err := account.TransactionHash(call, txDetails)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, hash.String(), expectedHash.String(), "transaction hash does not match expected")
	})
}

func TestSign(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	t.Run("Sign", func(t *testing.T) {
		expectedS1, _ := new(felt.Felt).SetString("0x0")
		expectedS2, _ := new(felt.Felt).SetString("0x0")

		ks := starknetgo.NewMemKeystore()
		fakeSenderAddress, _ := new(felt.Felt).SetString("0x135c34f53f8b7f59efd450eb689fccd9dd4cfe7f9d9dc4d09954c5653138698")
		fakeSenderAddressStr := fakeSenderAddress.String()
		k := types.SNValToBN(fakeSenderAddressStr)
		ks.Put(fakeSenderAddressStr, k)
		qwe, err := ks.Get(fakeSenderAddressStr)
		fmt.Println("qwe", qwe)
		require.NoError(t, err)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, ks, fakeSenderAddressStr)
		require.NoError(t, err, "error returned from account.NewAccount()")

		msg := new(felt.Felt).SetUint64(1)
		sig, err := account.Sign(context.Background(), msg)

		require.NoError(t, err, "error returned from account.Sign()")
		require.Equal(t, sig[0], expectedS1.String(), "s1 does not match expected")
		require.Equal(t, sig[1], expectedS2.String(), "s2 does not match expected")
	})
}

func TestExecute(t *testing.T) {
	// mockCtrl := gomock.NewController(t)
	// t.Cleanup(mockCtrl.Finish)
	// mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	t.Run("Test Execute", func(t *testing.T) {
		panic("Tests for Execute need implemented")
	})
}
