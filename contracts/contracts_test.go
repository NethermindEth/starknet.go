package contracts

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/artifacts"
	devtest "github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/joho/godotenv"
)

func TestGateway_InstallCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		providerType  starknet.go.ProviderType
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  starknet.go.ProviderGateway,
				CompiledClass: artifacts.CounterCompiled,
				Salt:          "0x0",
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
		case starknet.go.ProviderGateway:
			provider := GatewayProvider(*testConfiguration.gateway)
			tx, err = provider.deployAndWaitNoWallet(ctx, test.CompiledClass, test.Salt, test.Inputs)
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("deployment transaction", tx.TransactionHash)
	}
}

func TestRPCv02_InstallCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		providerType  starknet.go.ProviderType
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  starknet.go.ProviderRPCv02,
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
		case starknet.go.ProviderRPCv02:
			provider := RPCv02Provider(*testConfiguration.rpcv02)
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

func TestGateway_LoadAndExecuteCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		privateKey      string
		providerType    starknet.go.ProviderType
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0x01",
				providerType:    starknet.go.ProviderGateway,
				accountContract: artifacts.AccountContracts[ACCOUNT_VERSION1][false][false],
			},
		},
	}[testEnv]
	for _, test := range TestCases {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*120)
		defer cancel()
		var err error
		var counterTransaction *DeployOutput
		var account *starknet.go.Account
		// shim a keystore into existing tests.
		// use string representation of the PK as a fake sender address for the keystore
		ks := starknet.go.NewMemKeystore()

		fakeSenderAddress := test.privateKey
		k := types.SNValToBN(test.privateKey)
		ks.Put(fakeSenderAddress, k)
		switch test.providerType {
		case starknet.go.ProviderGateway:
			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
			accountManager, err := InstallAndWaitForAccount(
				ctx,
				testConfiguration.gateway,
				pk,
				test.accountContract,
			)
			if err != nil {
				t.Fatal("error deploying account", err)
			}
			mint, err := devtest.NewDevNet().Mint(types.StrToFelt(accountManager.AccountAddress), big.NewInt(int64(1000000000000000000)))
			if err != nil {
				t.Fatal("error deploying account", err)
			}
			fmt.Printf("current balance is %d\n", mint.NewBalance)
			provider := GatewayProvider(*testConfiguration.gateway)
			counterTransaction, err = provider.deployAndWaitNoWallet(ctx, artifacts.CounterCompiled, "0x0", []string{})
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
			fmt.Println("deployment transaction", counterTransaction.TransactionHash)
			account, err = starknet.go.NewGatewayAccount(types.StrToFelt(fakeSenderAddress), types.StrToFelt(accountManager.AccountAddress), ks, testConfiguration.gateway, starknet.go.AccountVersion1)
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		tx, err := account.Execute(ctx, []types.FunctionCall{{ContractAddress: types.StrToFelt(counterTransaction.ContractAddress), EntryPointSelector: "increment", Calldata: []string{}}}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("increment transaction", tx.TransactionHash)
	}
}

func TestRPCv02_LoadAndExecuteCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		privateKey      string
		providerType    starknet.go.ProviderType
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0xe3e70682c2094cac629f6fbed82c07cd",
				providerType:    starknet.go.ProviderRPCv02,
				accountContract: artifacts.AccountContracts[ACCOUNT_VERSION1][false][false],
			},
		},
	}[testEnv]
	for _, test := range TestCases {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*120)
		defer cancel()
		var err error
		var counterTransaction *DeployOutput
		var account *starknet.go.Account
		ks := starknet.go.NewMemKeystore()

		fakeSenderAddress := test.privateKey
		k := types.SNValToBN(test.privateKey)
		ks.Put(fakeSenderAddress, k)
		switch test.providerType {
		case starknet.go.ProviderRPCv02:
			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
			fmt.Println("befor")
			accountManager := &AccountManager{}
			accountManager, err := InstallAndWaitForAccount(
				ctx,
				testConfiguration.rpcv02,
				pk,
				test.accountContract,
			)
			if err != nil {
				t.Fatal("error deploying account", err)
			}
			fmt.Println("after")
			mint, err := devtest.NewDevNet().Mint(types.StrToFelt(accountManager.AccountAddress), big.NewInt(int64(1000000000000000000)))
			if err != nil {
				t.Fatal("error deploying account", err)
			}
			fmt.Printf("current balance is %d\n", mint.NewBalance)
			provider := RPCv02Provider(*testConfiguration.rpcv02)
			counterTransaction, err = provider.deployAndWaitWithWallet(ctx, artifacts.CounterCompiled, "0x0", []string{})
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
			fmt.Println("deployment transaction", counterTransaction.TransactionHash)
			account, err = starknet.go.NewRPCAccount(types.StrToFelt(fakeSenderAddress), types.StrToFelt(accountManager.AccountAddress), ks, testConfiguration.rpcv02, starknet.go.AccountVersion1)
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		tx, err := account.Execute(ctx, []types.FunctionCall{{ContractAddress: types.StrToFelt(counterTransaction.ContractAddress), EntryPointSelector: "increment", Calldata: []string{}}}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("increment transaction", tx.TransactionHash)
	}
}
