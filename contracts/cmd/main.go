package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
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
		TESTNET_ENV: "https://localhost:9545/v0.1/rpc",
		MAINNET_ENV: "https://localhost:9545/v0.1/rpc",
	},
}

type config struct {
	command        string
	accountVersion string
	withPlugin     bool
	skipCharge     bool
	withProxy      bool
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

	flagset := flag.NewFlagSet("", flag.ExitOnError)
	flagset.StringVar(&command, "command", "install", "defines the operation to execute")
	flagset.BoolVar(&skipCharge, "skip-charge", false, "do not charge the account (on devnet only)")
	flagset.BoolVar(&withPlugin, "with-plugin", false, "use a plugin/session-key account")
	flagset.BoolVar(&withProxy, "with-proxy", false, "use a proxy account")
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

	if provider != PROVIDER_GATEWAY && provider != PROVIDER_RPCV01 {
		log.Fatal("provider provider only supports rpcv01 and gateway")
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
		env:            env,
		baseURL:        baseURL,
		filename:       SECRET_FILE_NAME,
	}, nil
}

func main() {
	c, err := parse(os.Args[1:])
	if err != nil {
		log.Fatalf("could not run the command...")
	}
	switch c.command {
	case "install":
		if c.provider == "rpcv01" {
			c.installAccountWithRPCv01()
			return
		}
		c.installAccountWithGateway()
	case "execute":
		if c.withPlugin {
			c.incrementWithSessionKey()
			return
		}
		if c.provider == "rpcv01" {
			c.incrementWithRPCv01()
			return
		}
		c.incrementWithGateway()
	case "sum":
		if c.provider == "rpcv01" {
			log.Fatalf("rpcv01 not yet implemented")
			return
		}
		c.sumWithGateway()
	default:
		log.Fatalf("unknown command: %s\n", c.command)
	}
}
