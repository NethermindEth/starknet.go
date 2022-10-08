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
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

//go:embed artifacts/proxy.json
var proxyCompiled []byte

//go:embed artifacts/plugin.json
var pluginCompiled []byte

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

func declareClass(ctx context.Context, provider rpcv01.Provider, compiledClass []byte) (string, error) {
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

func deployContract(ctx context.Context, provider rpcv01.Provider, compiledClass []byte, salt string, inputs []string) (string, error) {
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

func (ap *accountPlugin) installAccount(ctx context.Context, provider rpcv01.Provider, plugin, account, proxy []byte) error {
	inputs := []string{}
	if len(proxy) != 0 {
		accountClassHash, err := declareClass(ctx, provider, account)
		if err != nil {
			return err
		}
		ap.AccountClassHash = accountClassHash
		inputs = append(inputs, accountClassHash)
	}
	inputs = append(inputs, ap.PublicKey)

	if len(plugin) != 0 {
		pluginClassHash, err := declareClass(ctx, provider, plugin)
		if err != nil {
			return err
		}
		ap.PluginClassHash = pluginClassHash
		inputs = append(inputs, pluginClassHash)
	}
	if len(proxy) == 0 {
		proxy = account
	}
	accountAddress, err := deployContract(ctx, provider, proxy, ap.PublicKey, inputs)
	if err != nil {
		return err
	}
	ap.AccountAddress = accountAddress
	err = ap.Write(SECRET_FILE_NAME)
	return err
}
