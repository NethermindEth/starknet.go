package account_test

import (
	"context"
	"flag"
	"math/big"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/require"
)

var (
	// set the environment for the test, default: mock
	testEnv = "mock"
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "devnet", "set the test environment")
	flag.Parse()
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
		account, err := account.NewAccount(mockRpcProvider, 1, address, starknetgo.NewMemKeystore(), false)
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
		account, err := account.NewAccount(mockRpcProvider, 1, accountAddress, starknetgo.NewMemKeystore(), false)
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

	t.Run("ChainId mainnet", func(t *testing.T) {
		mainnetID := utils.TestHexToFelt(t, "0x534e5f4d41494e")
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_MAIN", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, starknetgo.NewMemKeystore(), false)
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), mainnetID.String())
	})

	t.Run("ChainId testnet", func(t *testing.T) {
		testnetID := utils.TestHexToFelt(t, "0x534e5f474f45524c49")
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, starknetgo.NewMemKeystore(), false)
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), testnetID.String())
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
		account, err := account.NewAccount(mockRpcProvider, 1, address, ks, true)
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
		devNetURL := "http://0.0.0.0:5050"
		accounts, err := NewDevnet(t, devNetURL)
		require.NoError(t, err, "Error in NewDevnet")

		client, err := rpc.NewClient(devNetURL)
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		devAccount := accounts[0]
		priv, ok := new(big.Int).SetString(devAccount.PrivateKey, 0)
		require.True(t, ok)
		ks := starknetgo.SetNewMemKeystore(devAccount.PublicKey, priv)

		account, err := account.NewAccount(provider, 1, utils.TestHexToFelt(t, devAccount.Address), ks, false)
		require.NoError(t, err)

		invokeTx := rpc.BroadcastedInvokeV1Transaction{}
		_, err = account.AddInvokeTransaction(context.Background(), &invokeTx)
		require.NoError(t, err)

	})
}

func NewDevnet(t *testing.T, url string) ([]test.TestAccount, error) {
	// url := SetupLocalStarknetNode(t)
	devnet := test.NewDevNet(url)
	acnts, err := devnet.Accounts()
	return acnts, err
}

// // SetupLocalStarknetNode sets up a local starknet node via cli, and returns the url
// func SetupLocalStarknetNode(t *testing.T) string {
// 	url := "http://127.0.0.1:" + "5050"
// 	cmd := exec.Command("starknet-devnet",
// 		"--seed", "0", // use same seed for testing
// 		"--port", "5050",
// 		"--lite-mode",
// 	)
// 	var stdErr bytes.Buffer
// 	cmd.Stderr = &stdErr
// 	require.NoError(t, cmd.Start())
// 	t.Cleanup(func() {
// 		assert.NoError(t, cmd.Process.Kill())
// 		if err2 := cmd.Wait(); assert.Error(t, err2) {
// 			if !assert.Contains(t, err2.Error(), "signal: killed", cmd.ProcessState.String()) {
// 				t.Log("starknet-devnet stderr:", stdErr.String())
// 			}
// 		}
// 		t.Log("starknet-devnet server closed")
// 	})

// 	// Wait for api server to boot
// 	var ready bool
// 	for i := 0; i < 30; i++ {
// 		time.Sleep(time.Second)
// 		res, err := http.Get(url + "/is_alive")
// 		if err != nil || res.StatusCode != 200 {
// 			t.Logf("API server not ready yet (attempt %d)\n", i+1)
// 			continue
// 		}
// 		ready = true
// 		t.Logf("API server ready at %s\n", url)
// 		break
// 	}
// 	require.True(t, ready)
// 	return url
// }
