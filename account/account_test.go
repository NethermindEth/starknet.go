package account_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/joho/godotenv"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/require"
)

var (
	// set the environment for the test, default: mock
	testEnv = "mock"
	base    = ""
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "devnet", "set the test environment")
	flag.Parse()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	base = os.Getenv("INTEGRATION_BASE")
	if base == "" && testEnv != "mock" {
		panic("Failed to set INTEGRATION_BASE for non-mock test.")
	}
	os.Exit(m.Run())
}

func TestTransactionHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	t.Run("Transaction hash testnet", func(t *testing.T) {
		expectedHash := utils.TestHexToFelt(t, "0x135c34f53f8b7f59efd450eb689fccd9dd4cfe7f9d9dc4d09954c5653138698")
		address := &felt.Zero

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, "pubkey", starknetgo.NewMemKeystore())
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

	t.Run("Transaction hash mainnet", func(t *testing.T) {
		expectedHash := utils.TestHexToFelt(t, "0x3476c76a81522fe52616c41e95d062f5c3ea4eeb6c652904ad389fcd9ff4637")
		accountAddress := utils.TestHexToFelt(t, "0x59cd166e363be0a921e42dd5cfca0049aedcf2093a707ef90b5c6e46d4555a8")

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_MAIN", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, accountAddress, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err, "error returned from account.NewAccount()")

		call := rpc.FunctionCall{
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x5dbdedc203e92749e2e746e2d40a768d966bd243df04a6b712e222bc040a9af"),
				utils.TestHexToFelt(t, "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"),
				utils.TestHexToFelt(t, "0x0"),
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x52884ee3f"),
			},
		}
		txDetails := rpc.TxDetails{
			Nonce:   utils.TestHexToFelt(t, "0x1"),
			MaxFee:  utils.TestHexToFelt(t, "0x2a173cd36e400"),
			Version: rpc.TransactionV1,
		}
		hash, err := account.TransactionHash(call, txDetails)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, expectedHash.String(), hash.String(), "transaction hash does not match expected")
	})
}

func TestChainId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	t.Run("ChainId mainnet - mock", func(t *testing.T) {
		mainnetID := utils.TestHexToFelt(t, "0x534e5f4d41494e")
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_MAIN", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), mainnetID.String())
	})

	t.Run("ChainId testnet - mock", func(t *testing.T) {
		testnetID := utils.TestHexToFelt(t, "0x534e5f474f45524c49")
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), testnetID.String())
	})

	t.Run("ChainId devnet", func(t *testing.T) {
		if testEnv != "devnet" {
			t.Skip("Skipping test as it requires a devnet environment")
		}
		devNetURL := "http://0.0.0.0:5050/rpc"

		fmt.Println("devNetURL", devNetURL)
		client, err := rpc.NewClient(devNetURL)
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		_, err = account.NewAccount(provider, 1, &felt.Zero, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err)
	})
}

func TestSign(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	// Accepted on testnet https://goerli.voyager.online/tx/0x2a7eec54aab835323a810e893354368a496f1a217e8b6ef295476568ef08f0d
	t.Run("Sign testnet", func(t *testing.T) {
		expectedS1 := utils.TestHexToFelt(t, "0x6bf7980d98fa300ed9565b8cd5efcf5582133daa961b5e1d9477bf1bd750727")
		expectedS2 := utils.TestHexToFelt(t, "0x5886b8236b7dc3665c0014876a644ddd0800a167ff1036fb82af1a6f4134c91")

		address := utils.TestHexToFelt(t, "0x476466998f22e0b0177ddc76afcf8e3b5d30164f3eb33031aae7a9cb63c831")
		// pubKey := utils.TestHexToFelt(t, "0xc1e5fc3f93e04dac29878e4efccd81a96547884628d83b40fc9b758b5a349")
		privKey := utils.TestHexToFelt(t, "0x15d0b81e6140f4cce02b47609879a723f9f5b7b9f3ffca346018c73fe81847e")
		privKeyBI, ok := new(big.Int).SetString(privKey.String(), 0)
		require.True(t, ok)
		ks := starknetgo.NewMemKeystore()
		ks.Put(address.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, address.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		msg := utils.TestHexToFelt(t, "0x2a7eec54aab835323a810e893354368a496f1a217e8b6ef295476568ef08f0d")
		sig, err := account.Sign(context.Background(), msg)

		require.NoError(t, err, "error returned from account.Sign()")
		require.Equal(t, expectedS1.String(), sig[0].String(), "s1 does not match expected")
		require.Equal(t, expectedS2.String(), sig[1].String(), "s2 does not match expected")
	})
}

func TestAddInvoke(t *testing.T) {

	t.Run("Test AddInvoke devnet", func(t *testing.T) {
		if testEnv != "devnet" {
			t.Skip("Skipping test as it requires a devnet environment")
		}
		// New devnet
		devNetURL := "http://0.0.0.0:5050"
		accounts, err := newDevnet(t, devNetURL)
		require.NoError(t, err, "Error in NewDevnet")

		// New Client
		client, err := rpc.NewClient(devNetURL + "/rpc")
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		// Get devnet account + set up ks
		devAccount := accounts[0]
		priv, ok := new(big.Int).SetString(devAccount.PrivateKey, 0)
		require.True(t, ok)
		ks := starknetgo.SetNewMemKeystore(devAccount.PublicKey, priv)

		// Get account
		account, err := account.NewAccount(provider, 1, utils.TestHexToFelt(t, devAccount.Address), devAccount.PublicKey, ks)
		require.NoError(t, err)

		// Now build the trasaction
		invokeTx := rpc.BroadcastedInvokeV1Transaction{
			BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
				Nonce:   &felt.Zero,
				MaxFee:  new(felt.Felt).SetUint64(100),
				Version: rpc.TransactionV1,
				Type:    rpc.TransactionType_Invoke,
			},
			SenderAddress: account.AccountAddress,
		}
		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x49D36570D4E46F48E99674BD3FCC84644DDD6B96F7C741B1562B82F9E004DC7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("balanceOf"),
			Calldata:           []*felt.Felt{account.AccountAddress},
		}

		err = account.BuildInvokeTx(context.Background(), &invokeTx, &[]rpc.FunctionCall{fnCall})
		require.NoError(t, err, "Error in BuildInvokeTx")

		// Send tx
		resp, err := account.AddInvokeTransaction(context.Background(), &invokeTx)
		fmt.Println("resp", resp, err)
		require.NoError(t, err)

	})

	t.Run("Test AddInvoke mainnet", func(t *testing.T) {
		if testEnv != "mainnet" {
			t.Skip("Skipping test as it requires a mainnet environment")
		}

		// New Client
		fmt.Println("base", base)
		client, err := rpc.NewClient(base + "/rpc")
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		// Set up ks
		ks := starknetgo.NewMemKeystore()
		fakePrivKey := new(big.Int).SetUint64(123)
		fakePubKey := new(felt.Felt).SetUint64(456)
		ks.Put(fakePubKey.String(), fakePrivKey)

		// Get account
		account, err := account.NewAccount(provider, 1, fakePubKey, fakePubKey.String(), ks)
		require.NoError(t, err)

		// Now build the trasaction
		invokeTx := rpc.BroadcastedInvokeV1Transaction{
			BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
				Nonce:   &felt.Zero,
				MaxFee:  new(felt.Felt).SetUint64(1),
				Version: rpc.TransactionV1,
				Type:    rpc.TransactionType_Invoke,
			},
			SenderAddress: account.AccountAddress,
		}
		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("name"),
			Calldata:           []*felt.Felt{account.AccountAddress},
		}

		err = account.BuildInvokeTx(context.Background(), &invokeTx, &[]rpc.FunctionCall{fnCall})
		require.NoError(t, err, "Error in BuildInvokeTx")

		// Send tx
		resp, err := account.AddInvokeTransaction(context.Background(), &invokeTx)
		fmt.Println("resp", resp, err)
		require.NoError(t, err)
	})

	t.Run("Test AddInvoke testnet", func(t *testing.T) {
		if testEnv != "testnet" {
			t.Skip("Skipping test as it requires a testnet environment")
		}

		// New Client
		fmt.Println("base", base)
		client, err := rpc.NewClient(base + "/rpc")
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		// Set up ks
		ks := starknetgo.NewMemKeystore()
		fakePrivKey := new(big.Int).SetUint64(123)
		fakePubKey := new(felt.Felt).SetUint64(456)

		ks.Put(fakePubKey.String(), fakePrivKey)

		// Get account
		accountAddress := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")
		account, err := account.NewAccount(provider, 1, accountAddress, fakePubKey.String(), ks)
		require.NoError(t, err)

		// Now build the trasaction
		invokeTx := rpc.BroadcastedInvokeV1Transaction{
			BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
				Nonce:   &felt.Zero,
				MaxFee:  new(felt.Felt).SetUint64(1),
				Version: rpc.TransactionV1,
				Type:    rpc.TransactionType_Invoke,
			},
			SenderAddress: account.AccountAddress,
		}
		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("balanceOf"),
			Calldata:           []*felt.Felt{account.AccountAddress},
		}
		fmt.Println(fnCall.EntryPointSelector)
		err = account.BuildInvokeTx(context.Background(), &invokeTx, &[]rpc.FunctionCall{fnCall})
		require.NoError(t, err, "Error in BuildInvokeTx")

		qwe, _ := json.MarshalIndent(invokeTx, "", "")
		fmt.Println(string(qwe))
		// Send tx
		resp, err := account.AddInvokeTransaction(context.Background(), &invokeTx)
		fmt.Println("resp", resp, err)
		require.NoError(t, err)
	})
}

func newDevnet(t *testing.T, url string) ([]test.TestAccount, error) {
	// url := SetupLocalStarknetNode(t)
	devnet := test.NewDevNet(url)
	acnts, err := devnet.Accounts()
	return acnts, err
}
