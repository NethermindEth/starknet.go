package contracts

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/smartcontractkit/caigo"
	"github.com/smartcontractkit/caigo/artifacts"
	devtest "github.com/smartcontractkit/caigo/test"
	"github.com/smartcontractkit/caigo/types"
)

func TestGateway_InstallCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		providerType  caigo.ProviderType
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  caigo.ProviderGateway,
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
		case caigo.ProviderGateway:
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

func TestRPCv01_InstallCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		providerType  caigo.ProviderType
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  caigo.ProviderRPCv01,
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
		case caigo.ProviderRPCv01:
			provider := RPCv01Provider(*testConfiguration.rpcv01)
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
		providerType  caigo.ProviderType
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  caigo.ProviderRPCv02,
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
		case caigo.ProviderRPCv02:
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
		providerType    caigo.ProviderType
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0x01",
				providerType:    caigo.ProviderGateway,
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
		var account *caigo.Account
		// shim a keystore into existing tests.
		// use string representation of the PK as a fake sender address for the keystore
		ks := caigo.NewMemKeystore()

		fakeSenderAddress := test.privateKey
		k := types.SNValToBN(test.privateKey)
		ks.Put(fakeSenderAddress, k)
		switch test.providerType {
		case caigo.ProviderGateway:
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
			mint, err := devtest.NewDevNet().Mint(types.HexToHash(accountManager.AccountAddress), 1000000000000000000)
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
			account, err = caigo.NewGatewayAccount(types.StrToFelt(fakeSenderAddress), types.StrToFelt(accountManager.AccountAddress), ks, testConfiguration.gateway, caigo.AccountVersion1)
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		tx, err := account.Execute(ctx, []types.FunctionCall{{ContractAddress: types.HexToHash(counterTransaction.ContractAddress), EntryPointSelector: "increment", Calldata: []string{}}}, types.ExecuteDetails{})
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
		providerType    caigo.ProviderType
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0xe3e70682c2094cac629f6fbed82c07cd",
				providerType:    caigo.ProviderRPCv02,
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
		var account *caigo.Account
		ks := caigo.NewMemKeystore()

		fakeSenderAddress := test.privateKey
		k := types.SNValToBN(test.privateKey)
		ks.Put(fakeSenderAddress, k)
		switch test.providerType {
		case caigo.ProviderRPCv01:
			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
			fmt.Println("befor")
			accountManager := &AccountManager{}
			accountManager, err := InstallAndWaitForAccount(
				ctx,
				testConfiguration.rpcv01,
				pk,
				test.accountContract,
			)
			if err != nil {
				t.Fatal("error deploying account", err)
			}
			fmt.Println("after")
			mint, err := devtest.NewDevNet().Mint(types.HexToHash(accountManager.AccountAddress), 1000000000000000000)
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
			account, err = caigo.NewRPCAccount(types.StrToFelt(fakeSenderAddress), types.StrToFelt(accountManager.AccountAddress), ks, testConfiguration.rpcv02, caigo.AccountVersion1)
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		tx, err := account.Execute(ctx, []types.FunctionCall{{ContractAddress: types.HexToHash(counterTransaction.ContractAddress), EntryPointSelector: "increment", Calldata: []string{}}}, types.ExecuteDetails{})
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("increment transaction", tx.TransactionHash)
	}
}
