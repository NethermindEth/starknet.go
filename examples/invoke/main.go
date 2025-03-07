package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.
var (
	someContract   string = "0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54" //Replace it with the contract that you want to invoke. In this case, an ERC20
	contractMethod string = "mint"                                                               //Replace it with the function name that you want to invoke
)

func main() {
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

	simpleInvoke(accnt, contractAddress, contractMethod, amount)
	verboseInvoke(accnt, contractAddress, contractMethod, amount)
}
