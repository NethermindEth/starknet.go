package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/test"
	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

//go:embed artifacts/proxy.json
var proxyCompiled []byte

//go:embed artifacts/plugin.json
var pluginCompiled []byte

//go:embed artifacts/accountv0.json
var accountV0 []byte

//go:embed artifacts/accountv0_plugin.json
var accountV0WithPlugin []byte

//go:embed artifacts/accountv1.json
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

func newAccount() *accountPlugin {
	if _, err := os.Stat(SECRET_FILE_NAME); err == nil {
		log.Fatalf("file .starknet-account.json exists! exit...")
	}
	privateKey, _ := caigo.Curve.GetRandomPrivateKey()
	publicKey, _, _ := caigo.Curve.PrivateToPoint(privateKey)
	return &accountPlugin{
		PrivateKey: fmt.Sprintf("0x%s", privateKey.Text(16)),
		PublicKey:  fmt.Sprintf("0x%s", publicKey.Text(16)),
	}
}

type RPCProvider rpcv01.Provider

func (p *RPCProvider) declareClass(ctx context.Context, compiledClass []byte) (string, error) {
	provider := rpcv01.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return "", err
	}
	tx, err := provider.AddDeclareTransaction(ctx, class, "0x0")
	if err != nil {
		return "", err
	}
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("class Hash: %s\n", tx.ClassHash)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return "", err
	}
	if status == types.TransactionRejected {
		log.Printf("class Hash: %s\n", tx.ClassHash)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return "", errors.New("declare rejected")
	}
	return tx.ClassHash, nil
}

func (p *RPCProvider) deployContract(ctx context.Context, compiledClass []byte, salt string, inputs []string) (string, error) {
	provider := rpcv01.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return "", err
	}
	tx, err := provider.AddDeployTransaction(ctx, salt, inputs, class)
	if err != nil {
		return "", err
	}
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return "", err
	}
	if status == types.TransactionRejected {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return "", errors.New("deploy rejected")
	}
	return tx.ContractAddress, nil
}

type GatewayProvider gateway.Gateway

func (p *GatewayProvider) declareClass(ctx context.Context, compiledClass []byte) (string, error) {
	provider := gateway.Gateway(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return "", err
	}
	tx, err := provider.Declare(ctx, class, gateway.DeclareRequest{})
	if err != nil {
		return "", err
	}
	_, receipt, err := (&provider).WaitForTransaction(ctx, tx.TransactionHash, 3, 10)
	if err != nil {
		return "", err
	}
	if !receipt.Status.IsTransactionFinal() ||
		receipt.Status == types.TransactionRejected {
		return "", fmt.Errorf("wrong status: %s", receipt.Status)
	}
	return tx.ClassHash, nil
}

func (p *GatewayProvider) deployContract(ctx context.Context, compiledClass []byte, salt string, inputs []string) (string, error) {
	provider := gateway.Gateway(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return "", err
	}
	tx, err := provider.Deploy(ctx, class, types.DeployRequest{
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
	})
	if err != nil {
		return "", err
	}
	_, receipt, err := (&provider).WaitForTransaction(ctx, tx.TransactionHash, 3, 10)
	if err != nil {
		return "", err
	}
	if !receipt.Status.IsTransactionFinal() ||
		receipt.Status == types.TransactionRejected {
		return "", fmt.Errorf("wrong status: %s", receipt.Status)
	}
	return tx.ContractAddress, nil
}

func (ap *accountPlugin) installAccountWithRPCv01(ctx context.Context, provider rpcv01.Provider, plugin, account, proxy []byte) error {
	p := RPCProvider(provider)
	inputs := []string{}
	if len(proxy) != 0 {
		accountClassHash, err := (&p).declareClass(ctx, account)
		if err != nil {
			return err
		}
		ap.AccountClassHash = accountClassHash
		inputs = append(inputs, accountClassHash)
	}
	inputs = append(inputs, ap.PublicKey)

	if len(plugin) != 0 {
		pluginClassHash, err := (&p).declareClass(ctx, plugin)
		if err != nil {
			return err
		}
		ap.PluginClassHash = pluginClassHash
		inputs = append(inputs, pluginClassHash)
	}
	if len(proxy) == 0 {
		proxy = account
	}
	accountAddress, err := (&p).deployContract(ctx, proxy, ap.PublicKey, inputs)
	if err != nil {
		return err
	}
	ap.AccountAddress = accountAddress
	err = ap.Write(SECRET_FILE_NAME)
	return err
}

func (ap *accountPlugin) installAccountWithGateway(ctx context.Context, provider gateway.Gateway, plugin, account, proxy []byte) error {
	p := GatewayProvider(provider)
	inputs := []string{}
	if len(proxy) != 0 {
		accountClassHash, err := (&p).declareClass(ctx, account)
		if err != nil {
			return err
		}
		ap.AccountClassHash = accountClassHash
		inputs = append(inputs, accountClassHash)
	}
	inputs = append(inputs, ap.PublicKey)

	if len(plugin) != 0 {
		pluginClassHash, err := (&p).declareClass(ctx, plugin)
		if err != nil {
			return err
		}
		ap.PluginClassHash = pluginClassHash
		inputs = append(inputs, pluginClassHash)
	}
	if len(proxy) == 0 {
		proxy = account
	}
	accountAddress, err := (&p).deployContract(ctx, proxy, ap.PublicKey, inputs)
	if err != nil {
		return err
	}
	ap.AccountAddress = accountAddress
	err = ap.Write(SECRET_FILE_NAME)
	return err
}

func (c *config) installAccountWithRPCv01() {
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
		log.Fatalf("account v1 with plugin is not supported yet")
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
	err = account.installAccountWithRPCv01(ctx, *provider, plugin, compiledAccount, proxy)
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
	account := newAccount()
	ctx := context.Background()
	baseURL := "http://localhost:5050"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	if c.accountVersion == "v1" && c.withPlugin {
		log.Fatalf("account v1 with plugin is not supported yet")
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
	provider := gateway.NewClient(gateway.WithBaseURL(baseURL))
	account.Version = c.accountVersion
	account.Plugin = c.withPlugin
	err := account.installAccountWithGateway(ctx, *provider, plugin, compiledAccount, proxy)
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
