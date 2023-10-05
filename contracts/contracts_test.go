package contracts

import (
	"context"
	"fmt"

	// "math/big"
	"testing"
	"time"

	"github.com/joho/godotenv"
	// "github.com/NethermindEth/juno/core/felt"
	// starknetgo "github.com/NethermindEth/starknet.go"
	// "github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/artifacts"
	// devtest "github.com/NethermindEth/starknet.go/test"
	// "github.com/NethermindEth/starknet.go/utils"
)

func TestRPC_InstallCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		providerType  ProviderType
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  ProviderRPC,
				CompiledClass: artifacts.CounterCompiled,
				Salt:          "0x01",
				Inputs:        []string{},
			},
		},
	}[testEnv]
	for _, test := range TestCases {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		var err error
		var tx *DeployOutput
		switch test.providerType {
		case ProviderRPC:
			provider := RPCProvider(*testConfiguration.rpc)
			tx, err = provider.deployAndWaitWithWallet(ctx, test.CompiledClass, test.Salt, test.Inputs)
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("deployment transaction", tx.TransactionHash)
	}
}

// TODO: update test with new account implementation
// func TestRPC_LoadAndExecuteCounter(t *testing.T) {
// 	godotenv.Load()
// 	testConfiguration := beforeEach(t)

// 	type TestCase struct {
// 		privateKey      string
// 		providerType    ProviderType
// 		accountContract artifacts.CompiledContract
// 	}

// 	TestCases := map[string][]TestCase{
// 		"devnet": {
// 			{
// 				privateKey:      "0xe3e70682c2094cac629f6fbed82c07cd",
// 				providerType:    ProviderRPC,
// 				accountContract: artifacts.AccountContracts[ACCOUNT_VERSION1][false][false],
// 			},
// 		},
// 	}[testEnv]
// 	for _, test := range TestCases {
// 		ctx := context.Background()
// 		ctx, cancel := context.WithTimeout(ctx, time.Second*120)
// 		defer cancel()
// 		var err error
// 		var counterTransaction *DeployOutput
// 		var acc *account.Account
// 		ks := starknetgo.NewMemKeystore()

// 		fakeSenderAddress := test.privateKey
// 		k := utils.SNValToBN(test.privateKey)
// 		ks.Put(fakeSenderAddress, k)
// 		switch test.providerType {
// 		case ProviderRPC:
// 			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
// 			accountManager := &AccountManager{}
// 			accountManager, err := InstallAndWaitForAccount(
// 				ctx,
// 				testConfiguration.rpc,
// 				pk,
// 				test.accountContract,
// 			)
// 			if err != nil {
// 				t.Fatal("error deploying account", err)
// 			}
// 			mint, err := devtest.NewDevNet().Mint(utils.TestHexToFelt(t, accountManager.AccountAddress), big.NewInt(int64(1000000000000000000)))
// 			if err != nil {
// 				t.Fatal("error deploying account", err)
// 			}
// 			fmt.Printf("current balance is %d\n", mint.NewBalance)
// 			provider := RPCProvider(*testConfiguration.rpc)
// 			counterTransaction, err = provider.deployAndWaitWithWallet(ctx, artifacts.CounterCompiled, "0x0", []string{})
// 			if err != nil {
// 				t.Fatal("should succeed, instead", err)
// 			}
// 			fmt.Println("deployment transaction", counterTransaction.TransactionHash)
// 			acc, err = account.NewAccount(testConfiguration.rpc, utils.TestHexToFelt(t, accountManager.AccountAddress), accountManager.PublicKey, ks)
// 			if err != nil {
// 				t.Fatal("should succeed, instead", err)
// 			}
// 		default:
// 			t.Fatal("unsupported client type", test.providerType)
// 		}
// 		tx, err := acc.Execute(ctx, []utils.FunctionCall{{ContractAddress: utils.TestHexToFelt(t, counterTransaction.ContractAddress), EntryPointSelector: utils.GetSelectorFromNameFelt("increment"), Calldata: []*felt.Felt{}}}, utils.ExecuteDetails{})
// 		if err != nil {
// 			t.Fatal("should succeed, instead", err)
// 		}
// 		fmt.Println("increment transaction", tx.TransactionHash)
// 	}
// }
