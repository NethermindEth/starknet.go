package contracts

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "embed"

	"github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/plugins/xsessions"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/NethermindEth/juno/core/felt"
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

// func (ap *AccountManager) ExecuteWithSessionKey(counterAddress, selector string, provider *rpcv02.Provider) (string, error) {
// 	sessionPrivateKey, _ := starknet.go.Curve.GetRandomPrivateKey()
// 	sessionPublicKey, _, _ := starknet.go.Curve.PrivateToPoint(sessionPrivateKey)

// 	signedSessionKey, err := signSessionKey(ap.PrivateKey, ap.AccountAddress, counterAddress, "increment", types.BigToHex(sessionPublicKey))
// 	if err != nil {
// 		return "", err
// 	}
// 	plugin := xsessions.WithSessionKeyPlugin(
// 		ap.PluginClassHash,
// 		signedSessionKey,
// 	)
// 	v := starknet.go.AccountVersion0
// 	if ap.Version == "v1" {
// 		v = starknet.go.AccountVersion1
// 	}
// 	account, err := starknet.go.NewRPCAccount(
// 		types.BigToHex(sessionPrivateKey),
// 		ap.AccountAddress,
// 		provider,
// 		plugin,
// 		v,
// 	)
// 	if err != nil {
// 		return "", err
// 	}
// 	calls := []types.FunctionCall{
// 		{
// 			ContractAddress:    types.StrToFelt(counterAddress),
// 			EntryPointSelector: "increment",
// 			Calldata:           []string{},
// 		},
// 	}
// 	ctx := context.Background()
// 	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
// 	if err != nil {
// 		log.Printf("could not execute transaction %v\n", err)
// 		return "", err
// 	}
// 	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
// 	status, err := provider.WaitForTransaction(ctx, types.StrToFelt(tx.TransactionHash), 8*time.Second)
// 	if err != nil {
// 		log.Printf("could not execute transaction %v\n", err)
// 		return tx.TransactionHash, err
// 	}
// 	if status != types.TransactionAcceptedOnL2 {
// 		log.Printf("transaction has failed with %s", status)
// 		return tx.TransactionHash, fmt.Errorf("unexpected status: %s", status)
// 	}
// 	return tx.TransactionHash, nil
// }

func (ap *AccountManager) ExecuteWithGateway(counterAddress *felt.Felt, selector string, provider *gateway.GatewayProvider) (string, error) {
	v := starknet.go.AccountVersion0
	if ap.Version == "v1" {
		v = starknet.go.AccountVersion1
	}
	// shim in  the keystore. while weird and awkward, it's functionally ok because
	// 1. account manager doesn't seem to be used any where
	// 2. the account that is created below is scoped to this func
	ks := starknet.go.NewMemKeystore()
	fakeSenderAddress := ap.PrivateKey
	k := types.SNValToBN(ap.PrivateKey)
	ks.Put(fakeSenderAddress, k)
	fakeSenderAdd, err := utils.HexToFelt(fakeSenderAddress)
	if err != nil {
		return "", err
	}
	apAcntAdd, err := utils.HexToFelt(ap.AccountAddress)
	if err != nil {
		return "", err
	}
	account, err := starknet.go.NewGatewayAccount(
		fakeSenderAdd,
		apAcntAdd,
		ks,
		provider,
		v,
	)
	if err != nil {
		return "", err
	}
	calls := []types.FunctionCall{
		{
			ContractAddress:    counterAddress,
			EntryPointSelector: types.GetSelectorFromNameFelt("increment"),
			Calldata:           []*felt.Felt{},
		},
	}
	ctx := context.Background()
	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return "", err
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
	_, receipt, err := provider.WaitForTransaction(ctx, tx.TransactionHash.String(), 3, 10)
	if err != nil {
		log.Printf("could not execute transaction %v\n", err)
		return tx.TransactionHash.String(), err
	}
	if receipt.Status != types.TransactionAcceptedOnL2 {
		log.Printf("transaction has failed with %s", receipt.Status)
		return tx.TransactionHash.String(), fmt.Errorf("unexpected status: %s", receipt.Status)
	}
	return tx.TransactionHash.String(), nil
}

func (ap *AccountManager) CallWithGateway(call types.FunctionCall, provider *gateway.GatewayProvider) ([]string, error) {
	//  shim in  the keystore. while weird and awkward, it's functionally ok because
	// 1. account manager doesn't seem to be used any where
	// 2. the account that is created below is scoped to this func
	ks := starknet.go.NewMemKeystore()
	fakeSenderAddress := ap.PrivateKey
	k := types.SNValToBN(ap.PrivateKey)
	ks.Put(fakeSenderAddress, k)
	fakeSenderAdd, err := utils.HexToFelt(fakeSenderAddress)
	if err != nil {
		return nil, err
	}
	apAcntAdd, err := utils.HexToFelt(ap.AccountAddress)
	if err != nil {
		return nil, err
	}
	account, err := starknet.go.NewGatewayAccount(
		fakeSenderAdd,
		apAcntAdd,
		ks,
		provider,
	)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return account.Call(ctx, call)
}
