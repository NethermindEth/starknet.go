package contracts

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/NethermindEth/caigo"
	"github.com/NethermindEth/caigo/artifacts"
	"github.com/joho/godotenv"
)

func TestGateway_InstallAccounts(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		privateKey       string
		CompiledContract artifacts.CompiledContract
		providerType     caigo.ProviderType
	}

	devnet := []TestCase{}
	for _, provider := range []caigo.ProviderType{caigo.ProviderGateway} {
		for _, version := range []string{"v1"} {
			for _, proxy := range []bool{false, true} {
				for _, plugin := range []bool{false, true} {
					devnet = append(devnet, TestCase{
						privateKey:       "0x1",
						CompiledContract: artifacts.AccountContracts[version][proxy][plugin],
						providerType:     provider,
					})
				}
			}
		}
	}
	TestCases := map[string][]TestCase{
		"devnet": devnet,
	}[testEnv]
	for _, test := range TestCases {
		privateKey, _ := big.NewInt(0).SetString(test.privateKey, 0)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		var accountManager *AccountManager
		var err error
		switch test.providerType {
		case caigo.ProviderGateway:
			accountManager, err = InstallAndWaitForAccount(
				ctx,
				testConfiguration.gateway,
				privateKey,
				test.CompiledContract,
			)
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("deployment transaction", accountManager.TransactionHash)
	}
}

func TestRPCv02_InstallAccounts(t *testing.T) {
	godotenv.Load()
	testConfiguration := beforeEach(t)

	type TestCase struct {
		privateKey       string
		CompiledContract artifacts.CompiledContract
		providerType     caigo.ProviderType
	}

	devnet := []TestCase{}
	for _, provider := range []caigo.ProviderType{caigo.ProviderRPCv02} {
		for _, version := range []string{"v1"} {
			for _, proxy := range []bool{false, true} {
				for _, plugin := range []bool{false, true} {
					devnet = append(devnet, TestCase{
						privateKey:       "0x1",
						CompiledContract: artifacts.AccountContracts[version][proxy][plugin],
						providerType:     provider,
					})
				}
			}
		}
	}
	TestCases := map[string][]TestCase{
		"devnet": devnet,
	}[testEnv]
	for _, test := range TestCases {
		privateKey, _ := big.NewInt(0).SetString(test.privateKey, 0)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		var accountManager *AccountManager
		var err error
		switch test.providerType {
		case caigo.ProviderRPCv02:
			accountManager, err = InstallAndWaitForAccount(
				ctx,
				testConfiguration.rpcv02,
				privateKey,
				test.CompiledContract,
			)
		case caigo.ProviderGateway:
			accountManager, err = InstallAndWaitForAccount(
				ctx,
				testConfiguration.gateway,
				privateKey,
				test.CompiledContract,
			)
		default:
			t.Fatal("unsupported client type", test.providerType)
		}
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("deployment transaction", accountManager.TransactionHash)
	}
}
