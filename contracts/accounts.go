package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"encoding/json"
	"os"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/artifacts"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv02"
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
	declareAndWaitWithWallet(context context.Context, contractClass []byte) (*DeclareOutput, error)
	deployAccountAndWaitNoWallet(ctx context.Context, classHash types.Felt, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error)
}

const (
	ACCOUNT_VERSION0 = "v0"
	ACCOUNT_VERSION1 = "v1"
)

func guessProviderType(p interface{}) (Provider, error) {
	switch v := p.(type) {
	case *rpcv02.Provider:
		provider := RPCv02Provider(*v)
		return &provider, nil
	case *gateway.GatewayProvider:
		provider := GatewayProvider(*v)
		return &provider, nil
	}
	return nil, errors.New("unsupported type")
}

// InstallAndWaitForAccount installs an account with a DEPLOY command.
//
// Deprecated: this function should be replaced by InstallAndWaitForAccount
// that will use the DEPLOY_ACCOUNT syscall.
func InstallAndWaitForAccount[V *rpcv02.Provider | *gateway.GatewayProvider](ctx context.Context, provider V, privateKey *big.Int, compiledContracts artifacts.CompiledContract) (*AccountManager, error) {
	if len(compiledContracts.AccountCompiled) == 0 {
		return nil, errors.New("empty account")
	}
	privateKeyString := fmt.Sprintf("0x%s", privateKey.Text(16))
	publicKey, _, err := caigo.Curve.PrivateToPoint(privateKey)
	if err != nil {
		return nil, err
	}
	publicKeyString := fmt.Sprintf("0x0%s", publicKey.Text(16))
	fmt.Println("z")
	p, err := guessProviderType(provider)
	if err != nil {
		return nil, err
	}
	accountClassHash := ""
	// if len(compiledContracts.ProxyCompiled) != 0 {
	output, err := p.declareAndWaitWithWallet(ctx, compiledContracts.AccountCompiled)
	if err != nil {
		return nil, err
	}
	accountClassHash = output.classHash
	// }
	pluginClassHash := ""
	if len(compiledContracts.PluginCompiled) != 0 {
		output, err := p.declareAndWaitWithWallet(ctx, compiledContracts.PluginCompiled)
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
	fmt.Println("d")
	// TODO: compiledDeploed could be proxy
	deployedOutput, err := p.deployAccountAndWaitNoWallet(ctx, types.StrToFelt(accountClassHash), compiledDeployed, publicKeyString, calldata)
	if err != nil {
		return nil, err
	}
	fmt.Println("e")
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
