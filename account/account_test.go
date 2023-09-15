package account_test

import (
	"context"
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
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	base = os.Getenv("INTEGRATION_BASE")
	if base == "" && testEnv != "mock" {
		panic(fmt.Sprint("Failed to set INTEGRATION_BASE for ", testEnv))
	}
	os.Exit(m.Run())
}

func TestTransactionHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	// https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8
	t.Run("Transaction hash mock", func(t *testing.T) {
		if testEnv != "mock" {
			t.Skip("Skipping test as it requires a mock environment")
		}
		expectedHash := utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8")
		acntAddress := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")
		privKey := utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
		privKeyBI, ok := new(big.Int).SetString(privKey.String(), 0)
		require.True(t, ok)
		ks := starknetgo.NewMemKeystore()
		ks.Put(acntAddress.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, acntAddress, acntAddress.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		call := rpc.FunctionCall{
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
				utils.TestHexToFelt(t, "0x0"),
				utils.TestHexToFelt(t, "0x3"),
				utils.TestHexToFelt(t, "0x3"),
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x0"),
			},
		}
		txDetails := rpc.TxDetails{
			Nonce:   utils.TestHexToFelt(t, "0x2"),
			MaxFee:  utils.TestHexToFelt(t, "0x574fbde6000"),
			Version: rpc.TransactionV1,
		}
		hash, err := account.TransactionHash2(call.Calldata, txDetails.Nonce, txDetails.MaxFee, account.AccountAddress)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, expectedHash.String(), hash.String(), "transaction hash does not match expected")
	})

	t.Run("Transaction hash testnet", func(t *testing.T) {
		if testEnv != "testnet" {
			t.Skip("Skipping test as it requires a testnet environment")
		}
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
		hash, err := account.TransactionHash2(call.Calldata, txDetails.Nonce, txDetails.MaxFee, account.AccountAddress)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, hash.String(), expectedHash.String(), "transaction hash does not match expected")
	})

	t.Run("Transaction hash mainnet", func(t *testing.T) {
		if testEnv != "mainnet" {
			t.Skip("Skipping test as it requires a mainnet environment")
		}
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
		hash, err := account.TransactionHash2(call.Calldata, txDetails.Nonce, txDetails.MaxFee, account.AccountAddress)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, expectedHash.String(), hash.String(), "transaction hash does not match expected")
	})
}

func TestFmtCallData(t *testing.T) {

	t.Run("ChainId mainnet - mock", func(t *testing.T) {

		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("transfer"),
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x1"),
			},
		}
		expectedCallData := []*felt.Felt{
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			utils.TestHexToFelt(t, "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
			utils.TestHexToFelt(t, "0x0"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x0"),
		}
		fmt.Println("fnCall.asd", fnCall.EntryPointSelector)
		fmtCallData := account.FmtCalldata([]rpc.FunctionCall{fnCall})
		fmt.Println("fmtCallData", fmtCallData)
		require.Equal(t, fmtCallData, expectedCallData)
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

	// Accepted on testnet https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8
	t.Run("Sign testnet - mock 2", func(t *testing.T) {
		expectedS1 := utils.TestHexToFelt(t, "0x10d405427040655f118bc8b897e2f2f8147858bbcb0e3d6bc6dfbc6d0205e8")
		expectedS2 := utils.TestHexToFelt(t, "0x5cdfe4a3d5b63002e9011ec0ba59ae2b75a43cb2a3bc1699b35aa64cb9ca3cf")

		address := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")
		privKey := utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
		privKeyBI, ok := new(big.Int).SetString(privKey.String(), 0)
		require.True(t, ok)
		ks := starknetgo.NewMemKeystore()
		ks.Put(address.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, address.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		msg := utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8")
		sig, err := account.Sign(context.Background(), msg)

		require.NoError(t, err, "error returned from account.Sign()")
		require.Equal(t, expectedS1.String(), sig[0].String(), "s1 does not match expected")
		require.Equal(t, expectedS2.String(), sig[1].String(), "s2 does not match expected")
	})
}

func TestAddInvoke(t *testing.T) {

	// https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8#overview
	t.Run("Test AddInvokeTransction testnet", func(t *testing.T) {
		if testEnv != "testnet" {
			t.Skip("Skipping test as it requires a testnet environment")
		}
		client, err := rpc.NewClient(base + "/rpc")
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		// account address
		accountAddress := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")

		// Set up ks
		ks := starknetgo.NewMemKeystore()
		fakePubKey, _ := new(felt.Felt).SetString("0x049f060d2dffd3bf6f2c103b710baf519530df44529045f92c3903097e8d861f")
		fakePrivKey, _ := new(big.Int).SetString("0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa", 0)
		fakePrivKeyFelt, _ := new(felt.Felt).SetString("0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
		ks.Put(fakePubKey.String(), fakePrivKey)

		// Get account
		acnt, err := account.NewAccount(provider, 1, accountAddress, fakePubKey.String(), ks)
		require.NoError(t, err)

		// Now build the trasaction
		maxFee, _ := new(felt.Felt).SetString("0x574fbde6000")
		invokeTx := rpc.BroadcastedInvokeV1Transaction{
			BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
				Nonce:   new(felt.Felt).SetUint64(3),
				MaxFee:  maxFee,
				Version: rpc.TransactionV1,
				Type:    rpc.TransactionType_Invoke,
			},
			SenderAddress: acnt.AccountAddress,
		}
		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("transfer"),
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x1"),
			},
		}
		invokeTx.Calldata = account.FmtCalldata([]rpc.FunctionCall{fnCall})

		/// NEED TO FORMAT THE CALL DATA
		txHash, err := acnt.TransactionHash2(invokeTx.Calldata, invokeTx.Nonce, invokeTx.MaxFee, acnt.AccountAddress)
		x, y, err := starknetgo.Curve.SignFelt(txHash, fakePrivKeyFelt)
		require.NoError(t, err)
		invokeTx.Signature = []*felt.Felt{x, y}
		resp, err := acnt.AddInvokeTransaction(context.Background(), &invokeTx)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

func newDevnet(t *testing.T, url string) ([]test.TestAccount, error) {
	devnet := test.NewDevNet(url)
	acnts, err := devnet.Accounts()
	return acnts, err
}
