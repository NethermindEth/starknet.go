package account_test

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
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
		account, err := account.NewAccount(mockRpcProvider, 1, address, starknetgo.NewMemKeystore())
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

// func TestExecute(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	t.Cleanup(mockCtrl.Finish)
// 	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

// 	t.Run("Test Execute", func(t *testing.T) {
// 		expectedHash, _ := new(felt.Felt).SetString("0x0")
// 		address := &felt.Zero

// 		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
// 		account, err := account.NewAccount(mockRpcProvider, 1, address, starknetgo.NewMemKeystore())
// 		require.NoError(t, err, "error returned from account.NewAccount()")

// 		call := rpc.FunctionCall{ContractAddress: &felt.Zero, EntryPointSelector: &felt.Zero, Calldata: []*felt.Felt{&felt.Zero}}
// 		txDetails := rpc.TxDetails{Nonce: &felt.Zero, MaxFee: &felt.Zero}

// 		mockRpcProvider.EXPECT().AddInvokeTransaction(context.Background()).Return("SN_GOERLI", nil)

// 		resp, err := account.Execute(context.Background(), call, txDetails)
// 		require.NoError(t, err, "error returned from account.Execute()")
// 		require.Equal(t, resp.TransactionHash, expectedHash.String(), "transaction hash does not match expected")
// 	})
// }
