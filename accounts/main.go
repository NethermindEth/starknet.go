package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/test"
	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

func parse(args []string) (string, bool) {
	command := "install"
	skipCharge := true
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	flagset.StringVar(&command, "command", "install", "defines the operation to execute")
	flagset.BoolVar(&skipCharge, "skip-charge", false, "do not charge the account")
	flagset.Parse(args)
	return command, skipCharge
}

func main() {
	command, skipCharge := parse(os.Args[1:])
	switch command {
	case "install":
		accountWithPlugin := newAccount()
		ctx := context.Background()
		client, err := ethrpc.DialContext(ctx, "http://localhost:5050/rpc")
		if err != nil {
			log.Fatalf("error connecting to devnet, %v\n", err)
		}
		provider := rpcv01.NewProvider(client)
		err = accountWithPlugin.installAccount(ctx, *provider)
		if err != nil {
			log.Fatalf("error deploying account, %v\n", err)
		}
		if !skipCharge {
			d := test.NewDevNet()
			_, err := d.Mint(types.HexToHash(accountWithPlugin.ProxyAccountAddress), 1000000000000000000)
			if err != nil {
				log.Fatalf("error loading ETH, %v\n", err)
			}
		}
		log.Println("account installed with success", accountWithPlugin.ProxyAccountAddress)
	case "sessionkey":
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
		tx, err := accountWithPlugin.simulateSessionKey(counterAddress, "increment", provider)
		if err != nil {
			log.Fatalf("count not execute transaction, %v\n", err)
		}
		log.Println("transaction executed with success", tx)
	default:
		log.Fatalf("unknown command: %s\n", command)
	}
}
