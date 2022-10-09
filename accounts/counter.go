package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

//go:embed artifacts/counter.json
var counterCompiled []byte

func (ap *accountPlugin) installCounterWithRPCv01(ctx context.Context, provider *rpcv01.Provider) (string, error) {
	p := RPCProvider(*provider)
	return (&p).deployContract(ctx, counterCompiled, ap.PublicKey, []string{})
}

func (ap *accountPlugin) installCounterWithGateway(ctx context.Context, provider *gateway.Gateway) (string, error) {
	p := GatewayProvider(*provider)
	return (&p).deployContract(ctx, counterCompiled, ap.PublicKey, []string{})
}

func (c *config) incrementWithSessionKey() error {
	accountWithPlugin := &accountPlugin{}
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
	counterAddress, err := accountWithPlugin.installCounterWithRPCv01(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return err
	}
	tx, err := accountWithPlugin.executeWithSessionKey(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	log.Println("transaction executed with success", tx)
	return nil
}

func (c *config) incrementWithRPCv01() error {
	accountWithPlugin := &accountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	client, err := ethrpc.DialContext(ctx, c.baseURL)
	if err != nil {
		log.Fatalf("error connecting to devnet, %v\n", err)
		return fmt.Errorf("error connecting to devnet, %v\n", err)
	}
	provider := rpcv01.NewProvider(client)
	counterAddress, err := accountWithPlugin.installCounterWithRPCv01(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return err
	}
	tx, err := accountWithPlugin.executeWithRPCv01(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	log.Println("transaction executed with success", tx)
	return nil
}

func (c *config) incrementWithGateway() error {
	accountWithPlugin := &accountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	provider := gateway.NewClient(gateway.WithBaseURL(c.baseURL))
	counterAddress, err := accountWithPlugin.installCounterWithGateway(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return err
	}
	tx, err := accountWithPlugin.executeWithGateway(counterAddress, "increment", provider)
	if err != nil {
		log.Fatalf("count not execute transaction, %v\n", err)
		return err
	}
	log.Println("transaction executed with success", tx)
	return nil
}

func (c *config) sumWithGateway() error {
	accountWithPlugin := &accountPlugin{}
	accountWithPlugin.Read(SECRET_FILE_NAME)
	ctx := context.Background()
	provider := gateway.NewClient(gateway.WithBaseURL(c.baseURL))
	counterAddress, err := accountWithPlugin.installCounterWithGateway(ctx, provider)
	if err != nil {
		log.Fatalf("could not deploy the counter contract, %v\n", err)
		return fmt.Errorf("could not deploy the counter contract, %v\n", err)
	}
	res, err := accountWithPlugin.callWithGateway(types.FunctionCall{
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
