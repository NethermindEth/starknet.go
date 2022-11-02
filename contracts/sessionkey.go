package contracts

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "embed"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/plugins/xsessions"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/rpcv02"
	"github.com/dontpanicdao/caigo/types"
)

func signSessionKey(privateKey, accountAddress, counterAddress, selector, sessionPublicKey string) (*xsessions.SessionKeyToken, error) {
	return xsessions.SignToken(
		privateKey,
		"0x"+types.UTF8StrToBig("SN_GOERLI").Text(16),
		sessionPublicKey,
		accountAddress,
		time.Hour*2,
		[]xsessions.Policy{{
			ContractAddress: counterAddress,
			Selector:        selector,
		}},
	)
}

// TODO: why we do have *rpcv02.Provider
func (ap *AccountManager) ExecuteWithSessionKey(counterAddress, selector string, provider *rpcv02.Provider) (string, error) {
	sessionPrivateKey, _ := caigo.Curve.GetRandomPrivateKey()
	sessionPublicKey, _, _ := caigo.Curve.PrivateToPoint(sessionPrivateKey)

	signedSessionKey, err := signSessionKey(ap.PrivateKey, ap.AccountAddress, counterAddress, "increment", types.BigToHex(sessionPublicKey))
	if err != nil {
		return "", err
	}
	plugin := xsessions.WithSessionKeyPlugin(
		ap.PluginClassHash,
		signedSessionKey,
	)
	v := caigo.AccountVersion0
	if ap.Version == "v1" {
		v = caigo.AccountVersion1
	}
	account, err := caigo.NewRPCAccount(
		types.BigToHex(sessionPrivateKey),
		ap.AccountAddress,
		provider,
		plugin,
		v,
	)
	if err != nil {
		return "", err
	}
	calls := []types.FunctionCall{
		{
			ContractAddress:    types.HexToHash(counterAddress),
			EntryPointSelector: "increment",
			Calldata:           []string{},
		},
	}
	ctx := context.Background()
	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return "", err
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return tx.TransactionHash, err
	}
	if status != types.TransactionAcceptedOnL2 {
		log.Printf("transaction has failed with %s", status)
		return tx.TransactionHash, fmt.Errorf("unexpected status: %s", status)
	}
	return tx.TransactionHash, nil
}

func (ap *AccountManager) ExecuteWithRPCv01(counterAddress, selector string, provider *rpcv01.Provider) (string, error) {
	v := caigo.AccountVersion0
	if ap.Version == "v1" {
		v = caigo.AccountVersion1
	}
	account, err := caigo.NewRPCAccount(
		ap.PrivateKey,
		ap.AccountAddress,
		provider,
		v,
	)
	if err != nil {
		return "", err
	}
	calls := []types.FunctionCall{
		{
			ContractAddress:    types.HexToHash(counterAddress),
			EntryPointSelector: "increment",
			Calldata:           []string{},
		},
	}
	ctx := context.Background()
	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return "", err
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return tx.TransactionHash, err
	}
	if status != types.TransactionAcceptedOnL2 {
		log.Printf("transaction has failed with %s", status)
		return tx.TransactionHash, fmt.Errorf("unexpected status: %s", status)
	}
	return tx.TransactionHash, nil
}

func (ap *AccountManager) ExecuteWithRPCv02(counterAddress, selector string, provider *rpcv02.Provider) (string, error) {
	v := caigo.AccountVersion0
	if ap.Version == "v1" {
		v = caigo.AccountVersion1
	}
	account, err := caigo.NewRPCAccount(
		ap.PrivateKey,
		ap.AccountAddress,
		provider,
		v,
	)
	if err != nil {
		return "", err
	}
	calls := []types.FunctionCall{
		{
			ContractAddress:    types.HexToHash(counterAddress),
			EntryPointSelector: "increment",
			Calldata:           []string{},
		},
	}
	ctx := context.Background()
	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return "", err
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return tx.TransactionHash, err
	}
	if status != types.TransactionAcceptedOnL2 {
		log.Printf("transaction has failed with %s", status)
		return tx.TransactionHash, fmt.Errorf("unexpected status: %s", status)
	}
	return tx.TransactionHash, nil
}

func (ap *AccountManager) ExecuteWithGateway(counterAddress, selector string, provider *gateway.GatewayProvider) (string, error) {
	v := caigo.AccountVersion0
	if ap.Version == "v1" {
		v = caigo.AccountVersion1
	}
	account, err := caigo.NewGatewayAccount(
		ap.PrivateKey,
		ap.AccountAddress,
		provider,
		v,
	)
	if err != nil {
		return "", err
	}
	calls := []types.FunctionCall{
		{
			ContractAddress:    types.HexToHash(counterAddress),
			EntryPointSelector: "increment",
			Calldata:           []string{},
		},
	}
	ctx := context.Background()
	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return "", err
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
	_, receipt, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 3, 10)
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return tx.TransactionHash, err
	}
	if receipt.Status != types.TransactionAcceptedOnL2 {
		log.Printf("transaction has failed with %s", receipt.Status)
		return tx.TransactionHash, fmt.Errorf("unexpected status: %s", receipt.Status)
	}
	return tx.TransactionHash, nil
}

func (ap *AccountManager) CallWithGateway(call types.FunctionCall, provider *gateway.GatewayProvider) ([]string, error) {
	account, err := caigo.NewGatewayAccount(
		ap.PrivateKey,
		ap.AccountAddress,
		provider,
	)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return account.Call(ctx, call)
}
