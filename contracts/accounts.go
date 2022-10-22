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

func guessProviderType(p interface{}) (Provider, error) {
	switch v := p.(type) {
	case *rpcv01.Provider:
		provider := RPCv01Provider(*v)
		return &provider, nil
	}
	return nil, errors.New("unsupported type")
}

// InstallAndWaitForAccount installs an account with a DEPLOY command.
//
// Deprecated: this function should be replaced by InstallAndWaitForAccount
// that will use the DEPLOY_ACCOUNT syscall.
func InstallAndWaitForAccountNoWallet[V *rpcv01.Provider | *gateway.GatewayProvider](ctx context.Context, provider V, privateKey *big.Int, compiledPlugin, compiledAccount, compiledProxy []byte) (*AccountManager, error) {
	if len(compiledAccount) == 0 {
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
	calldata := []string{}
	if len(compiledProxy) != 0 {
		output, err := p.declareAndWaitNoWallet(ctx, compiledAccount)
		if err != nil {
			return nil, err
		}
		initialize := types.GetSelectorFromName("initialize")
		calldata = append(calldata, output.classHash)
		calldata = append(calldata, fmt.Sprintf("0x%s", initialize.Text(16)))
		paramLen := "0x1"
		if len(compiledPlugin) != 0 {
			paramLen = "0x2"
		}
		calldata = append(calldata, paramLen)
		accountClassHash = output.classHash
	}
	calldata = append(calldata, publicKeyString)
	pluginClassHash := ""
	if len(compiledPlugin) != 0 {
		output, err := p.declareAndWaitNoWallet(ctx, compiledPlugin)
		if err != nil {
			return nil, err
		}
		pluginClassHash = output.classHash
		calldata = append(calldata, pluginClassHash)
	}
	compiledDeployed := compiledAccount
	if len(compiledProxy) != 0 {
		compiledDeployed = compiledProxy
	}
	deployedOutput, err := p.deployAndWaitNoWallet(ctx, compiledDeployed, publicKeyString, calldata)
	if err != nil {
		return nil, err
	}
	proxyClassHash := ""
	switch len(compiledProxy) {
	case 0:
		accountClassHash = deployedOutput.classHash
	default:
		proxyClassHash = deployedOutput.classHash
	}
	return &AccountManager{
		AccountAddress:   deployedOutput.contractAddress,
		AccountClassHash: accountClassHash,
		PluginClassHash:  pluginClassHash,
		PrivateKey:       privateKeyString,
		ProxyClassHash:   proxyClassHash,
		PublicKey:        publicKeyString,
		TransactionHash:  deployedOutput.transactionHash,
		Version:          "v0",
	}, nil

}
