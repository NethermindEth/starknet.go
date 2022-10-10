package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dontpanicdao/caigo/accounts"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/test"
	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

const SECRET_FILE_NAME = ".starknet-account.json"

func (c *config) incrementWithSessionKey() error {
	accountWithPlugin := &accounts.AccountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	if accountWithPlugin.PluginClassHash == "" {
		log.Println("account not installed with plugin, stop!")
		return fmt.Errorf("account not installed with plugin")
	}
	ctx := context.Background()
	baseURL := "http://localhost:5050/rpc"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	client, err := ethrpc.DialContext(ctx, baseURL)
	if err != nil {
		log.Fatalf("error connecting to devnet, %v\n", err)
		return err
	}
	provider := rpcv01.NewProvider(client)
	counterAddress, err := accountWithPlugin.InstallCounterWithRPCv01(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return err
	}
	tx, err := accountWithPlugin.ExecuteWithSessionKey(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	log.Println("transaction executed with success", tx)
	return nil
}

func (c *config) incrementWithRPCv01() error {
	accountWithPlugin := &accounts.AccountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	client, err := ethrpc.DialContext(ctx, c.baseURL)
	if err != nil {
		log.Fatalf("error connecting to devnet, %v\n", err)
		return fmt.Errorf("error connecting to devnet, %v\n", err)
	}
	provider := rpcv01.NewProvider(client)
	counterAddress, err := accountWithPlugin.InstallCounterWithRPCv01(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return err
	}
	tx, err := accountWithPlugin.ExecuteWithRPCv01(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	log.Println("transaction executed with success", tx)
	return nil
}

func (c *config) incrementWithGateway() error {
	accountWithPlugin := &accounts.AccountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	provider := gateway.NewClient(gateway.WithBaseURL(c.baseURL))
	counterAddress, err := accountWithPlugin.InstallCounterWithGateway(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return err
	}
	tx, err := accountWithPlugin.ExecuteWithGateway(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	log.Println("transaction executed with success", tx)
	return nil
}

func (c *config) sumWithGateway() error {
	accountWithPlugin := &accounts.AccountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	provider := gateway.NewClient(gateway.WithBaseURL(c.baseURL))
	counterAddress, err := accountWithPlugin.InstallCounterWithGateway(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return fmt.Errorf("could not deploy the counter contract, %v\n", err)
	}
	res, err := accountWithPlugin.CallWithGateway(types.FunctionCall{
		ContractAddress:    types.HexToHash(counterAddress),
		EntryPointSelector: "sum",
		Calldata:           []string{"0x1", "0x2"},
	}, provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	if len(res) != 1 || res[0] != "0x3" {
		log.Fatalf("sum should be 3, instead %+v", res)
		return fmt.Errorf("wrong results")
	}
	log.Printf("sum(1+2)=%s", res[0])
	return nil
}

func (c *config) installAccountWithRPCv01() {
	account := accounts.NewAccount(".starknet")
	ctx := context.Background()
	baseURL := "http://localhost:5050/rpc"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	client, err := ethrpc.DialContext(ctx, baseURL)
	if err != nil {
		log.Fatalf("error connecting to devnet, %v\n", err)
	}
	if c.accountVersion == "v1" && c.withPlugin {
		log.Fatalf("account v1 with plugin is not supported yet")
	}
	compiledAccount := accounts.AccountContent[c.accountVersion][c.withPlugin]
	plugin := []byte{}
	if c.withPlugin {
		plugin = accounts.PluginCompiled
	}
	proxy := []byte{}
	if c.withProxy {
		proxy = accounts.ProxyCompiled
	}
	provider := rpcv01.NewProvider(client)
	account.Version = c.accountVersion
	account.Plugin = c.withPlugin
	err = account.InstallAccountWithRPCv01(ctx, *provider, plugin, compiledAccount, proxy)
	if err != nil {
		log.Fatalf("error installing account to devnet/rpcv01, %v\n", err)
	}
	if !c.skipCharge {
		d := test.NewDevNet()
		_, err := d.Mint(types.HexToHash(account.AccountAddress), 1000000000000000000)
		if err != nil {
			log.Fatalf("error loading ETH, %v\n", err)
		}
	}
	log.Println("account installed with success", account.AccountAddress)
}

func (c *config) installAccountWithGateway() {
	account := accounts.NewAccount(SECRET_FILE_NAME)
	ctx := context.Background()
	baseURL := "http://localhost:5050"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	if c.accountVersion == "v1" && c.withPlugin {
		log.Fatalf("account v1 with plugin is not supported yet")
	}
	compiledAccount := accounts.AccountContent[c.accountVersion][c.withPlugin]
	plugin := []byte{}
	if c.withPlugin {
		plugin = accounts.PluginCompiled
	}
	proxy := []byte{}
	if c.withProxy {
		proxy = accounts.ProxyCompiled
	}
	provider := gateway.NewClient(gateway.WithBaseURL(baseURL))
	account.Version = c.accountVersion
	account.Plugin = c.withPlugin
	err := account.InstallAccountWithGateway(ctx, *provider, plugin, compiledAccount, proxy)
	if err != nil {
		log.Fatalf("error installing account to devnet/gateway, %v\n", err)
	}
	if !c.skipCharge {
		d := test.NewDevNet()
		_, err := d.Mint(types.HexToHash(account.AccountAddress), 1000000000000000000)
		if err != nil {
			log.Fatalf("error loading ETH, %v\n", err)
		}
	}
	log.Println("account installed with success", account.AccountAddress)
}
