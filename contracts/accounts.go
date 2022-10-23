package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"encoding/json"
	"os"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

type AccountManager struct {
	AccountAddress   string `json:"accountAddress"`
	AccountClassHash string `json:"accountClassHash,omitempty"`
	filename         string
	PluginClassHash  string `json:"pluginClassHash,omitempty"`
	PrivateKey       string `json:"privateKey"`
	ProxyClassHash   string `json:"proxyClassHash,omitempty"`
	PublicKey        string `json:"publicKey"`
	TransactionHash  string `json:"transactionHash,omitempty"`
	Version          string `json:"accountVersion"`
}

func (ap *AccountManager) Read(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	json.Unmarshal(content, ap)
	return nil
}

func (ap *AccountManager) Write(filename string) error {
	content, err := json.Marshal(ap)
	if err != nil {
		return err
	}
	ap.filename = filename
	return os.WriteFile(filename, content, 0664)
}

type Provider interface {
	declareAndWaitNoWallet(context context.Context, contractClass []byte) (*DeclareOutput, error)
	deployAndWaitNoWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error)
}

const (
	PROVIDER_GATEWAY = "gateway"
	PROVIDER_RPCV01  = "rpcv01"
	ACCOUNT_VERSION0 = "v0"
	ACCOUNT_VERSION1 = "v1"
)

func guessProviderType(p interface{}) (Provider, error) {
	switch v := p.(type) {
	case *rpcv01.Provider:
		provider := RPCv01Provider(*v)
		return &provider, nil
	case *gateway.GatewayProvider:
		provider := GatewayProvider(*v)
		return &provider, nil
	}
	return nil, errors.New("unsupported type")
}

type calldataFormater func(accountHash, pluginHash, publicKey string) ([]string, error)

func AccountV0Formater(accountHash, pluginHash, publicKey string) ([]string, error) {
	calldata := []string{}
	if accountHash == "" {
		calldata = append(calldata, publicKey)
		if pluginHash != "" {
			calldata = append(calldata, pluginHash)
		}
		return calldata, nil
	}
	calldata = append(calldata, accountHash)
	initialize := fmt.Sprintf("0x%s", types.GetSelectorFromName("initialize").Text(16))
	calldata = append(calldata, initialize)
	paramLen := "0x1"
	if pluginHash != "" {
		paramLen = "0x2"
	}
	calldata = append(calldata, paramLen)
	calldata = append(calldata, publicKey)
	if pluginHash != "" {
		calldata = append(calldata, pluginHash)
	}
	return calldata, nil
}

func AccountFormater(accountHash, pluginHash, publicKey string) ([]string, error) {
	if pluginHash != "" {
		return nil, errors.New("plugin not supported")
	}
	calldata := []string{}
	if accountHash == "" {
		calldata = append(calldata, publicKey)
		if pluginHash != "" {
			calldata = append(calldata, pluginHash)
		}
		return calldata, nil
	}
	calldata = append(calldata, accountHash)
	initialize := fmt.Sprintf("0x%s", types.GetSelectorFromName("initialize").Text(16))
	calldata = append(calldata, initialize)
	calldata = append(calldata, "0x1")
	calldata = append(calldata, publicKey)
	return calldata, nil
}

func AccountPluginFormater(accountHash, pluginHash, publicKey string) ([]string, error) {
	if pluginHash == "" {
		return nil, errors.New("plugin is mandatory")
	}
	calldata := []string{}
	if accountHash != "" {
		calldata = append(calldata, accountHash)
		initialize := fmt.Sprintf("0x%s", types.GetSelectorFromName("initialize").Text(16))
		calldata = append(calldata, initialize)
		calldata = append(calldata, "0x4")
	}
	calldata = append(calldata, pluginHash)
	calldata = append(calldata, "0x2")
	calldata = append(calldata, "0x1")
	calldata = append(calldata, publicKey)
	return calldata, nil
}

// InstallAndWaitForAccount installs an account with a DEPLOY command.
//
// Deprecated: this function should be replaced by InstallAndWaitForAccount
// that will use the DEPLOY_ACCOUNT syscall.
func InstallAndWaitForAccountNoWallet[V *rpcv01.Provider | *gateway.GatewayProvider](ctx context.Context, provider V, privateKey *big.Int, compiledContracts CompiledContract) (*AccountManager, error) {
	if len(compiledContracts.AccountCompiled) == 0 {
		return nil, errors.New("empty account")
	}
	privateKeyString := fmt.Sprintf("0x%s", privateKey.Text(16))
	publicKey, _, err := caigo.Curve.PrivateToPoint(privateKey)
	if err != nil {
		return nil, err
	}
	publicKeyString := fmt.Sprintf("0x%s", publicKey.Text(16))
	p, err := guessProviderType(provider)
	if err != nil {
		return nil, err
	}
	accountClassHash := ""
	if len(compiledContracts.ProxyCompiled) != 0 {
		output, err := p.declareAndWaitNoWallet(ctx, compiledContracts.AccountCompiled)
		if err != nil {
			return nil, err
		}
		accountClassHash = output.classHash
	}
	pluginClassHash := ""
	if len(compiledContracts.PluginCompiled) != 0 {
		output, err := p.declareAndWaitNoWallet(ctx, compiledContracts.PluginCompiled)
		if err != nil {
			return nil, err
		}
		pluginClassHash = output.classHash
	}
	compiledDeployed := compiledContracts.AccountCompiled
	if len(compiledContracts.ProxyCompiled) != 0 {
		compiledDeployed = compiledContracts.ProxyCompiled
	}
	calldata, err := compiledContracts.Formatter(accountClassHash, pluginClassHash, publicKeyString)
	if err != nil {
		return nil, err
	}
	deployedOutput, err := p.deployAndWaitNoWallet(ctx, compiledDeployed, publicKeyString, calldata)
	if err != nil {
		return nil, err
	}
	proxyClassHash := ""
	switch len(compiledContracts.ProxyCompiled) {
	case 0:
		accountClassHash = deployedOutput.ClassHash
	default:
		proxyClassHash = deployedOutput.ClassHash
	}
	return &AccountManager{
		AccountAddress:   deployedOutput.ContractAddress,
		AccountClassHash: accountClassHash,
		PluginClassHash:  pluginClassHash,
		PrivateKey:       privateKeyString,
		ProxyClassHash:   proxyClassHash,
		PublicKey:        publicKeyString,
		TransactionHash:  deployedOutput.TransactionHash,
		Version:          "v0",
	}, nil

}
