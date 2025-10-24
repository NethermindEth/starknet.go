package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/typeddata"
	"github.com/NethermindEth/starknet.go/utils"
)

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.

func main() {
	// Setup the account
	accnt := localSetup()
	fmt.Println("Account address:", accnt.Address)

	// This is how you can initialise a typed data from a JSON file
	ttd, err := utils.UnmarshalJSONFileToType[typeddata.TypedData]("./baseExample.json", "")
	if err != nil {
		panic(fmt.Errorf("fail to unmarshal TypedData: %w", err))
	}

	// get the message hash linked to your account address
	messageHash, err := ttd.GetMessageHash(accnt.Address.String())
	if err != nil {
		panic(fmt.Errorf("fail to get message hash: %w", err))
	}
	fmt.Println("Message hash:", messageHash)

	// sign the message hash
	signature, err := accnt.Sign(context.Background(), messageHash)
	if err != nil {
		panic(fmt.Errorf("fail to sign message: %w", err))
	}
	fmt.Println("Signature:", signature)

	// verify the signature
	isValid, err := accnt.Verify(messageHash, signature)
	if err != nil {
		panic(fmt.Errorf("fail to verify signature: %w", err))
	}

	fmt.Println("Verification result:", isValid)
}

func localSetup() *account.Account {
	// Load variables from '.env' file
	rpcProviderURL := setup.GetRPCProviderURL()
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()

	// Initialise connection to RPC provider
	client, err := rpc.NewProvider(context.Background(), rpcProviderURL)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the RPC provider: %s", err))
	}

	// Initialise the account memkeyStore (set public and private keys)
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
	// Initialise the account
	accnt, err := account.NewAccount(
		client,
		accountAddressInFelt,
		publicKey,
		ks,
		accountCairoVersion,
	)
	if err != nil {
		panic(err)
	}

	return accnt
}
