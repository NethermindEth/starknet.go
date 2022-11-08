package contracts

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/artifacts"
	devtest "github.com/dontpanicdao/caigo/test"
	"github.com/dontpanicdao/caigo/types"
	"github.com/joho/godotenv"
)

func TestGateway_InstallCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		providerType  string
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  PROVIDER_GATEWAY,
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
		case "gateway":
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
		providerType  string
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  PROVIDER_RPCV01,
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
		case "rpcv01":
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
		providerType  string
		CompiledClass []byte
		Salt          string
		Inputs        []string
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				providerType:  PROVIDER_RPCV01,
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
		case "rpcv01":
			provider := RPCv01Provider(*testConfiguration.rpcv01)
			tx, err = provider.deployAndWaitNoWallet(ctx, test.CompiledClass, test.Salt, test.Inputs)
		case "gateway":
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

func TestGateway_LoadAndExecuteCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		privateKey      string
		providerType    string
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0x1",
				providerType:    PROVIDER_GATEWAY,
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
		switch test.providerType {
		case "gateway":
			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
			accountManager, err := InstallAndWaitForAccountNoWallet(
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
			account, err = caigo.NewGatewayAccount(test.privateKey, accountManager.AccountAddress, testConfiguration.gateway, caigo.AccountVersion1)
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

func TestRPCv01_LoadAndExecuteCounter(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		privateKey      string
		providerType    string
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0x1",
				providerType:    PROVIDER_RPCV01,
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
		switch test.providerType {
		case "rpcv01":
			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
			accountManager, err := InstallAndWaitForAccountNoWallet(
				ctx,
				testConfiguration.rpcv01,
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
			provider := RPCv01Provider(*testConfiguration.rpcv01)
			counterTransaction, err = provider.deployAndWaitNoWallet(ctx, artifacts.CounterCompiled, "0x0", []string{})
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
			fmt.Println("deployment transaction", counterTransaction.TransactionHash)
			account, err = caigo.NewRPCAccount(test.privateKey, accountManager.AccountAddress, testConfiguration.rpcv01, caigo.AccountVersion1)
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
		providerType    string
		accountContract artifacts.CompiledContract
	}

	TestCases := map[string][]TestCase{
		"devnet": {
			{
				privateKey:      "0x1",
				providerType:    PROVIDER_RPCV01,
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
		switch test.providerType {
		case "rpcv01":
			pk, _ := big.NewInt(0).SetString(test.privateKey, 0)
			accountManager, err := InstallAndWaitForAccountNoWallet(
				ctx,
				testConfiguration.rpcv01,
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
			provider := RPCv01Provider(*testConfiguration.rpcv01)
			counterTransaction, err = provider.deployAndWaitNoWallet(ctx, artifacts.CounterCompiled, "0x0", []string{})
			if err != nil {
				t.Fatal("should succeed, instead", err)
			}
			fmt.Println("deployment transaction", counterTransaction.TransactionHash)
			account, err = caigo.NewRPCAccount(test.privateKey, accountManager.AccountAddress, testConfiguration.rpcv01, caigo.AccountVersion1)
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
