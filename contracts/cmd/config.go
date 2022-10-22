package main

import (
	"flag"
	"log"
	"strings"
)

const SECRET_FILE_NAME = ".starknet-account.json"

const (
	DEVNET_ENV       = "devnet"
	TESTNET_ENV      = "testnet"
	MAINNET_ENV      = "mainnet"
	PROVIDER_GATEWAY = "gateway"
	PROVIDER_RPCV01  = "rpcv01"
	ACCOUNT_VERSION0 = "v0"
	ACCOUNT_VERSION1 = "v1"
)

type defaultURLS map[string]string

var providerDefaultURLS = map[string]defaultURLS{
	PROVIDER_GATEWAY: {
		DEVNET_ENV:  "http://localhost:5050",
		TESTNET_ENV: "https://alpha4.starknet.io",
		MAINNET_ENV: "https://alpha4-mainnet.starknet.io",
	},
	PROVIDER_RPCV01: {
		DEVNET_ENV:  "http://localhost:5050/rpc",
		TESTNET_ENV: "https://localhost:9545/rpc/v0.1",
		MAINNET_ENV: "https://localhost:9545/rpc/v0.1",
	},
}

type config struct {
	command        string
	accountVersion string
	withPlugin     bool
	skipCharge     bool
	withProxy      bool
	debug          bool
	provider       string
	baseURL        string
	env            string
	filename       string
}

func parse(args []string) (*config, error) {
	command := "install"
	skipCharge := false
	withPlugin := false
	withProxy := false
	accountVersion := ACCOUNT_VERSION0
	env := DEVNET_ENV
	provider := PROVIDER_RPCV01
	baseURL := ""
	debug := false

	flagset := flag.NewFlagSet("", flag.ExitOnError)
	flagset.StringVar(&command, "command", "install", "defines the operation to execute")
	flagset.BoolVar(&skipCharge, "skip-charge", false, "do not charge the account (on devnet only)")
	flagset.BoolVar(&withPlugin, "with-plugin", false, "use a plugin/session-key account")
	flagset.BoolVar(&withProxy, "with-proxy", false, "use a proxy account")
	flagset.BoolVar(&debug, "debug", false, "run the command in debug mode")
	flagset.StringVar(&accountVersion, "account-version", "v0", "choose v0 or v1 account")
	flagset.StringVar(&provider, "provider", "rpcv01", "choose rpc01 or gateway provider")
	flagset.StringVar(&env, "env", "devnet", "change the environment between devnet/testnet/mainnet")
	flagset.StringVar(&baseURL, "base-url", "", "baseURL depends on the app")
	err := flagset.Parse(args)

	if strings.Contains(strings.ToLower(env), "devnet") {
		env = DEVNET_ENV
	}
	if strings.Contains(strings.ToLower(env), "goerli") || strings.Contains(strings.ToLower(env), "testnet") {
		env = TESTNET_ENV
	}
	if strings.Contains(strings.ToLower(env), "mainnet") {
		env = MAINNET_ENV
	}
	if env == TESTNET_ENV || env == MAINNET_ENV {
		skipCharge = true
	}
	if baseURL == "" {
		baseURL = providerDefaultURLS[provider][env]
	}
	if accountVersion != ACCOUNT_VERSION0 && accountVersion != ACCOUNT_VERSION1 {
		log.Fatal("account-version only supports v0 and v1")
	}
	if err != nil {
		return nil, err
	}
	return &config{
		command:        command,
		accountVersion: accountVersion,
		withPlugin:     withPlugin,
		withProxy:      withProxy,
		skipCharge:     skipCharge,
		provider:       provider,
		debug:          debug,
		env:            env,
		baseURL:        baseURL,
		filename:       SECRET_FILE_NAME,
	}, nil
}

// func (c *config) incrementWithSessionKey() error {
// 	accountWithPlugin := &contracts.AccountManager{}
// 	accountWithPlugin.Read(c.filename)
// 	if accountWithPlugin.PluginClassHash == "" {
// 		log.Println("account not installed with plugin, stop!")
// 		return fmt.Errorf("account not installed with plugin")
// 	}
// 	ctx := context.Background()
// 	baseURL := "http://localhost:5050/rpc"
// 	if c.baseURL != "" {
// 		baseURL = c.baseURL
// 	}
// 	client, err := ethrpc.DialContext(ctx, baseURL)
// 	if err != nil {
// 		log.Fatalf("error connecting to devnet, %v\n", err)
// 		return err
// 	}
// 	provider := rpcv01.NewProvider(client)
// 	counterAddress, err := accountWithPlugin.InstallCounterWithRPCv01(ctx, provider)
// 	if err != nil {
// 		log.Fatalf("could not deploy the counter contract, %v\n", err)
// 		return err
// 	}
// 	tx, err := accountWithPlugin.ExecuteWithSessionKey(counterAddress, "increment", provider)
// 	if err != nil {
// 		log.Fatalf("count not execute transaction, %v\n", err)
// 		return err
// 	}
// 	log.Println("transaction executed with success", tx)
// 	return nil
// }

// func (c *config) incrementWithRPCv01() error {
// 	accountWithPlugin := &contracts.AccountPlugin{}
// 	accountWithPlugin.Read(c.filename)
// 	ctx := context.Background()
// 	client, err := ethrpc.DialContext(ctx, c.baseURL)
// 	if err != nil {
// 		log.Fatalf("error connecting to devnet, %v\n", err)
// 		return fmt.Errorf("error connecting to devnet, %v\n", err)
// 	}
// 	provider := rpcv01.NewProvider(client)
// 	counterAddress, err := accountWithPlugin.InstallCounterWithRPCv01(ctx, provider)
// 	if err != nil {
// 		log.Fatalf("could not deploy the counter contract, %v\n", err)
// 		return err
// 	}
// 	tx, err := accountWithPlugin.ExecuteWithRPCv01(counterAddress, "increment", provider)
// 	if err != nil {
// 		log.Fatalf("count not execute transaction, %v\n", err)
// 		return err
// 	}
// 	log.Println("transaction executed with success", tx)
// 	return nil
// }

// func (c *config) incrementWithGateway() error {
// 	accountWithPlugin := &contracts.AccountPlugin{}
// 	accountWithPlugin.Read(c.filename)
// 	ctx := context.Background()
// 	provider := gateway.NewClient(gateway.WithBaseURL(c.baseURL))
// 	counterAddress, err := accountWithPlugin.InstallCounterWithGateway(ctx, provider)
// 	if err != nil {
// 		log.Fatalf("could not deploy the counter contract, %v\n", err)
// 		return err
// 	}
// 	tx, err := accountWithPlugin.ExecuteWithGateway(counterAddress, "increment", provider)
// 	if err != nil {
// 		log.Fatalf("count not execute transaction, %v\n", err)
// 		return err
// 	}
// 	log.Println("transaction executed with success", tx)
// 	return nil
// }

// func (c *config) sumWithGateway() error {
// 	accountWithPlugin := &contracts.AccountPlugin{}
// 	accountWithPlugin.Read(c.filename)
// 	ctx := context.Background()
// 	provider := gateway.NewClient(gateway.WithBaseURL(c.baseURL))
// 	counterAddress, err := accountWithPlugin.InstallCounterWithGateway(ctx, provider)
// 	if err != nil {
// 		log.Fatalf("could not deploy the counter contract, %v\n", err)
// 		return fmt.Errorf("could not deploy the counter contract, %v\n", err)
// 	}
// 	res, err := accountWithPlugin.CallWithGateway(types.FunctionCall{
// 		ContractAddress:    types.HexToHash(counterAddress),
// 		EntryPointSelector: "sum",
// 		Calldata:           []string{"0x1", "0x2"},
// 	}, provider)
// 	if err != nil {
// 		log.Fatalf("count not execute transaction, %v\n", err)
// 		return err
// 	}
// 	if len(res) != 1 || res[0] != "0x3" {
// 		log.Fatalf("sum should be 3, instead %+v", res)
// 		return fmt.Errorf("wrong results")
// 	}
// 	log.Printf("sum(1+2)=%s", res[0])
// 	return nil
// }

// func (c *config) installAccountWithRPCv01() {
// 	account := contracts.NewAccount(c.filename)
// 	ctx := context.Background()
// 	baseURL := "http://localhost:5050/rpc"
// 	if c.baseURL != "" {
// 		baseURL = c.baseURL
// 	}
// 	client, err := ethrpc.DialContext(ctx, baseURL)
// 	if err != nil {
// 		log.Fatalf("error connecting to devnet, %v\n", err)
// 	}
// 	if c.accountVersion == "v1" && c.withPlugin {
// 		log.Fatalf("account v1 with plugin is not supported yet")
// 	}
// 	compiledAccount := contracts.CompiledAccounts[c.accountVersion][c.withPlugin]
// 	plugin := []byte{}
// 	if c.withPlugin {
// 		plugin = contracts.PluginCompiled
// 	}
// 	proxy := []byte{}
// 	if c.withProxy {
// 		proxy = contracts.ProxyCompiled
// 	}
// 	provider := rpcv01.NewProvider(client)
// 	account.Version = c.accountVersion
// 	account.Plugin = c.withPlugin
// 	err = account.InstallAccountWithRPCv01(ctx, *provider, plugin, compiledAccount, proxy)
// 	if err != nil {
// 		log.Fatalf("error installing account to devnet/rpcv01, %v\n", err)
// 	}
// 	if !c.skipCharge {
// 		d := test.NewDevNet()
// 		_, err := d.Mint(types.HexToHash(account.AccountAddress), 1000000000000000000)
// 		if err != nil {
// 			log.Fatalf("error loading ETH, %v\n", err)
// 		}
// 	}
// 	log.Println("account installed with success", account.AccountAddress)
// }

// func (c *config) installAccountWithGateway() {
// 	account := contracts.NewAccount(c.filename)
// 	ctx := context.Background()
// 	baseURL := "http://localhost:5050"
// 	if c.baseURL != "" {
// 		baseURL = c.baseURL
// 	}
// 	if c.accountVersion == "v1" && c.withPlugin {
// 		log.Fatalf("account v1 with plugin is not supported yet")
// 	}
// 	compiledAccount := contracts.CompiledAccounts[c.accountVersion][c.withPlugin]
// 	plugin := []byte{}
// 	if c.withPlugin {
// 		plugin = contracts.PluginCompiled
// 	}
// 	proxy := []byte{}
// 	if c.withProxy {
// 		proxy = contracts.ProxyCompiled
// 	}
// 	provider := gateway.NewClient(gateway.WithBaseURL(baseURL))
// 	account.Version = c.accountVersion
// 	account.Plugin = c.withPlugin
// 	err := account.InstallAccountWithGateway(ctx, *provider, plugin, compiledAccount, proxy)
// 	if err != nil {
// 		log.Fatalf("error installing account to devnet/gateway, %v\n", err)
// 	}
// 	if !c.skipCharge {
// 		d := test.NewDevNet()
// 		_, err := d.Mint(types.HexToHash(account.AccountAddress), 1000000000000000000)
// 		if err != nil {
// 			log.Fatalf("error loading ETH, %v\n", err)
// 		}
// 	}
// 	log.Println("account installed with success", account.AccountAddress)
// }
