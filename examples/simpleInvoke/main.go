package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.
var (
	name                  string = "testnet"                                                            //env."name"
	account_addr          string = "0x06f36e8a0fc06518125bbb1c63553e8a7d8597d437f9d56d891b8c7d3c977716" //Replace it with your account address
	account_cairo_version        = 0                                                                    //Replace  with the cairo version of your account
	privateKey            string = "0x0687bf84896ee63f52d69e6de1b41492abeadc0dc3cb7bd351d0a52116915937" //Replace it with your account private key
	public_key            string = "0x58b0824ee8480133cad03533c8930eda6888b3c5170db2f6e4f51b519141963"  //Replace it with your account public key
	someContract          string = "0x4c1337d55351eac9a0b74f3b8f0d3928e2bb781e5084686a892e66d49d510d"   //Replace it with the contract that you want to invoke
	contractMethod        string = "increase_value"                                                     //Replace it with the function name that you want to invoke
)

func main() {
	// Loading the env
	godotenv.Load(fmt.Sprintf(".env.%s", name))
	url := os.Getenv("INTEGRATION_BASE") //please modify the .env.testnet and replace the INTEGRATION_BASE with a starknet goerli RPC.
	fmt.Println("Starting simpleInvoke example")

	// Initialising the connection
	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	// Here we are converting the account address to felt
	account_address, err := utils.HexToFelt(account_addr)
	if err != nil {
		panic(err.Error())
	}
	// Initializing the account memkeyStore
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic(err.Error())
	}
	ks.Put(public_key, fakePrivKeyBI)

	fmt.Println("Established connection with the client")

	// Here we are setting the maxFee
	maxfee, err := utils.HexToFelt("0x9184e72a000")
	if err != nil {
		panic(err.Error())
	}

	// Initializing the account
	accnt, err := account.NewAccount(clientv02, account_address, public_key, ks, account_cairo_version)
	if err != nil {
		panic(err.Error())
	}

	// Getting the nonce from the account
	nonce, err := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	if err != nil {
		panic(err.Error())
	}

	// Building the InvokeTx struct
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        maxfee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: accnt.AccountAddress,
	}

	// Converting the contractAddress from hex to felt
	contractAddress, err := utils.HexToFelt(someContract)
	if err != nil {
		panic(err.Error())
	}

	// Building the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,                               //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod), //this is the function that we want to call
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		panic(err.Error())
	}

	// Signing of the transaction that is done by the account
	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx)
	if err != nil {
		panic(err.Error())
	}

	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
	resp, err := accnt.AddInvokeTransaction(context.Background(), InvokeTx)
	if err != nil {
		panic(err.Error())
	}
	// This returns us with the transaction hash
	fmt.Println("Transaction hash response : ", resp.TransactionHash)

}
