package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.
var (
	someContract   string = "0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54" // Replace it with the contract that you want to invoke. In this case, an ERC20
	contractMethod string = "mint"                                                               // Replace it with the function name that you want to invoke
)

// main is the main function that will be executed when the program is run.
// It will load the variables from the '.env' file, initialise the connection to the RPC provider,
// initialise the account, and then call the simpleInvoke and verboseInvoke functions passing the account,
// the contract address, the contract method and the amount to be sent.
func main() {
	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()

	// Initialise connection to RPC provider
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the RPC provider: %s", err))
	}

	// Initialise the account memkeyStore (set public and private keys)
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Failed to convert privKey to bigInt")
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

	fmt.Println("Established connection with the client")

	// Converting the contractAddress from hex to felt
	contractAddress, err := utils.HexToFelt(someContract)
	if err != nil {
		panic(err)
	}

	amount, err := utils.HexToFelt("0xffffffff")
	if err != nil {
		panic(err)
	}

	// Here we have two examples of how to send an invoke transaction, one is simple and the other one is verbose.
	// The simple example is more user-friendly and easier to use, while the verbose example is more detailed and informative.
	// You can choose one of them to run, or both! Each one will send a different transaction, but with almost the same parameters.
	simpleInvoke(accnt, contractAddress, contractMethod, amount)
	fmt.Println("--------------------------------")
	verboseInvoke(accnt, contractAddress, contractMethod, amount)
}
