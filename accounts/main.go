package main

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"os"

	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/test"
	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type config struct {
	command        string
	accountVersion string
	withPlugin     bool
	skipCharge     bool
	withProxy      bool
	provider       string
	baseURL        string
}

func (c *config) installAccount() {
	account := newAccount()
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
		log.Fatalf("acount v1 with plugin is not supported yet")
	}
	compiledAccount := accountContent[c.accountVersion][c.withPlugin]
	plugin := []byte{}
	if c.withPlugin {
		plugin = pluginCompiled
	}
	proxy := []byte{}
	if c.withProxy {
		proxy = proxyCompiled
	}
	provider := rpcv01.NewProvider(client)
	account.Version = c.accountVersion
	account.Plugin = c.withPlugin
	err = account.installAccount(ctx, *provider, plugin, compiledAccount, proxy)
	if err != nil {
		log.Fatalf("error installing account to devnet, %v\n", err)
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

func (c *config) incrementWithSessionKey() {
	accountWithPlugin := &accountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	if accountWithPlugin.PluginClassHash == "" {
		log.Println("account not installed with plugin, stop!")
	}
	ctx := context.Background()
	baseURL := "http://localhost:5050/rpc"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	client, err := ethrpc.DialContext(ctx, baseURL)
	if err != nil {
		log.Fatalf("error connecting to devnet, %v\n", err)
	}
	provider := rpcv01.NewProvider(client)
	counterAddress, err := accountWithPlugin.installCounter(ctx, *provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
	}
	tx, err := accountWithPlugin.executeWithSessionKey(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
	}
	log.Println("transaction executed with success", tx)
}

func (c *config) increment() {
	accountWithPlugin := &accountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	client, err := ethrpc.DialContext(ctx, "http://localhost:5050/rpc")
	if err != nil {
		log.Fatalf("error connecting to devnet, %v\n", err)
	}
	provider := rpcv01.NewProvider(client)
	counterAddress, err := accountWithPlugin.installCounter(ctx, *provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
	}
	tx, err := accountWithPlugin.execute(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
	}
	log.Println("transaction executed with success", tx)
}

//go:embed artifacts/account_v0.3.2.json
var accountV0 []byte

//go:embed artifacts/account_plugin.json
var accountV0WithPlugin []byte

//go:embed artifacts/account_v0.4.0b.json
var accountV1 []byte

var accountContent = map[string]map[bool][]byte{
	"v0": {
		true:  accountV0WithPlugin,
		false: accountV0,
	},
	"v1": {
		false: accountV1,
	},
}

func parse(args []string) (*config, error) {
	command := "install"
	skipCharge := false
	withPlugin := false
	withProxy := false
	accountVersion := "v0"
	provider := "rpcv01"
	baseURL := "http://localhost:5050/rpc"
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	flagset.StringVar(&command, "command", "install", "defines the operation to execute")
	flagset.BoolVar(&skipCharge, "skip-charge", false, "do not charge the account")
	flagset.BoolVar(&withPlugin, "with-plugin", false, "use a plugin/session-key account")
	flagset.BoolVar(&withProxy, "with-proxy", false, "use a proxy account")
	flagset.StringVar(&accountVersion, "account-version", "v0", "choose v0 or v1 account")
	flagset.StringVar(&provider, "provider", "rpcv01", "choose rpc01 or gateway provider")
	flagset.StringVar(&baseURL, "base-url", "http://localhost:5050/rpc", "change the default baseURL")
	err := flagset.Parse(args)

	if provider != "rpcv01" {
		log.Fatal("provider provider only supports rpcv01")
	}
	if accountVersion != "v0" && accountVersion != "v1" {
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
		baseURL:        baseURL,
	}, nil
}

func main() {
	c, err := parse(os.Args[1:])
	if err != nil {
		log.Fatalf("could not run the command...")
	}
	switch c.command {
	case "install":
		c.installAccount()
	case "execute":
		if c.withPlugin {
			c.incrementWithSessionKey()
			return
		}
		c.increment()
	default:
		log.Fatalf("unknown command: %s\n", c.command)
	}
}
