package contracts

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
	"github.com/dontpanicdao/caigo/types"
)

//go:embed artifacts/proxyv0.json
var ProxyV0Compiled []byte

//go:embed artifacts/accountv0.json
var AccountV0 []byte

//go:embed artifacts/proxyv0.json
var ProxyV0WithPluginCompiled []byte

//go:embed artifacts/accountv0_plugin.json
var AccountV0WithPluginCompiled []byte

//go:embed artifacts/pluginv0.json
var PluginV0Compiled []byte

//go:embed artifacts/accountv0_plugin.json
var AccountV0WithPlugin []byte

//go:embed artifacts/proxy.json
var ProxyCompiled []byte

//go:embed artifacts/account.json
var Account []byte

//go:embed artifacts/plugin.json
var PluginCompiled []byte

//go:embed artifacts/account_plugin.json
var AccountWithPlugin []byte

var AccountContent = map[string]map[bool][]byte{
	"v0": {
		true:  AccountV0WithPlugin,
		false: AccountV0,
	},
	"v1": {
		false: AccountV1,
	},
}

func NewAccount(filename string) *AccountPlugin {
	if _, err := os.Stat(filename); err == nil {
		log.Fatalf("file .starknet-account.json exists! exit...")
	}
	privateKey, _ := caigo.Curve.GetRandomPrivateKey()
	publicKey, _, _ := caigo.Curve.PrivateToPoint(privateKey)
	return &AccountPlugin{
		PrivateKey: fmt.Sprintf("0x%s", privateKey.Text(16)),
		PublicKey:  fmt.Sprintf("0x%s", publicKey.Text(16)),
		filename:   filename,
	}
}

type RPCProvider rpcv01.Provider

func (p *RPCProvider) DeclareClass(ctx context.Context, compiledClass []byte) (string, error) {
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

func (ap *AccountPlugin) InstallAccountWithRPCv01(ctx context.Context, provider rpcv01.Provider, plugin, account, proxy []byte) error {
	p := RPCProvider(provider)
	inputs := []string{}
	if len(proxy) != 0 {
		accountClassHash, err := (&p).DeclareClass(ctx, account)
		if err != nil {
			return err
		}
		ap.AccountClassHash = accountClassHash
		inputs = append(inputs, accountClassHash)
	}
	inputs = append(inputs, ap.PublicKey)

	if len(plugin) != 0 {
		pluginClassHash, err := (&p).DeclareClass(ctx, plugin)
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
	err = ap.Write(ap.filename)
	return err
}

func (ap *AccountPlugin) InstallAccountWithGateway(ctx context.Context, provider gateway.Gateway, plugin, account, proxy []byte) error {
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
	err = ap.Write(ap.filename)
	return err
}
