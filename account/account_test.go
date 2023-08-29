package account_test

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/test-go/testify/require"
)

func TestTransactionHash(t *testing.T) {
	c, err := ethrpc.DialContext(context.Background(), "http://0.0.0.0:5050/rpc")
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	provider := rpc.NewProvider(c)

	t.Run("Transaction hash", func(t *testing.T) {
		expectedHash, _ := new(felt.Felt).SetString("0x135c34f53f8b7f59efd450eb689fccd9dd4cfe7f9d9dc4d09954c5653138698")

		address := &felt.Zero
		account, err := account.NewAccount(provider, 1, address, starknetgo.NewMemKeystore())
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
