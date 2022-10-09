package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
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

	if provider != "rpcv01" && provider != "gateway" {
		log.Fatal("provider provider only supports rpcv01 and gateway")
	}
	if provider == "gateway" && baseURL == "http://localhost:5050/rpc" {
		baseURL = "http://localhost:5050"
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
