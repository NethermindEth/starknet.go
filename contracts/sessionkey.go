package contracts

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "embed"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/plugins/xsessions"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
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

// func (ap *AccountManager) ExecuteWithSessionKey(counterAddress, selector string, provider *rpc.Provider) (string, error) {
// 	sessionPrivateKey, _ := starknetgo.Curve.GetRandomPrivateKey()
// 	sessionPublicKey, _, _ := starknetgo.Curve.PrivateToPoint(sessionPrivateKey)

// 	signedSessionKey, err := signSessionKey(ap.PrivateKey, ap.AccountAddress, counterAddress, "increment", types.BigToHex(sessionPublicKey))
// 	if err != nil {
// 		return "", err
// 	}
// 	plugin := xsessions.WithSessionKeyPlugin(
// 		ap.PluginClassHash,
// 		signedSessionKey,
// 	)
// 	v := starknetgo.AccountVersion0
// 	if ap.Version == "v1" {
// 		v = starknetgo.AccountVersion1
// 	}
// 	account, err := starknetgo.NewRPCAccount(
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
