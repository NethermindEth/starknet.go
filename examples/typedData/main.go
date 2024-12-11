package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/typedData"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.

func main() {
	// Setup the account
	accnt := localSetup()
	fmt.Println("Account address:", accnt.AccountAddress)

	// This is how you can initialize a typed data from a JSON file
	var ttd typedData.TypedData
	content, err := os.ReadFile("./baseExample.json")
	if err != nil {
		panic(fmt.Errorf("fail to read file: %w", err))
	}
	err = json.Unmarshal(content, &ttd)
	if err != nil {
		panic(fmt.Errorf("fail to unmarshal TypedData: %w", err))
	}

	// This is how you can get the message hash linked to your account address
	messageHash, err := ttd.GetMessageHash(accnt.AccountAddress.String())
	if err != nil {
		panic(fmt.Errorf("fail to get message hash: %w", err))
	}
	fmt.Println("Message hash:", messageHash)

	// This is how you can sign the message hash
	signature, err := accnt.Sign(context.Background(), messageHash)
	if err != nil {
		panic(fmt.Errorf("fail to sign message: %w", err))
	}
	fmt.Println("Signature:", signature)

	// This is how you can verify the signature
	privKeyBI, ok := new(big.Int).SetString(setup.GetPrivateKey(), 0)
	if !ok {
		panic("Fail to convert privKey to bitInt")
	}
	x, y, err := curve.Curve.PrivateToPoint(privKeyBI)
	if err != nil {
		panic(fmt.Errorf("fail to get point: %w", err))
	}
	isValid := curve.Curve.Verify(messageHash.BigInt(new(big.Int)), signature[0].BigInt(new(big.Int)), signature[1].BigInt(new(big.Int)), x, y)
	fmt.Println("Verification result:", isValid)
}

func localSetup() *account.Account {
	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()

	// Initialize connection to RPC provider
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	// Initialize the account memkeyStore (set public and private keys)
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Fail to convert privKey to bitInt")
	}
	ks.Put(publicKey, privKeyBI)

	// Here we are converting the account address to felt
	accountAddressInFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		fmt.Println("Failed to transform the account address, did you give the hex address?")
		panic(err)
	}
	// Initialize the account
	accnt, err := account.NewAccount(client, accountAddressInFelt, publicKey, ks, accountCairoVersion)
	if err != nil {
		panic(err)
	}

	return accnt
}
