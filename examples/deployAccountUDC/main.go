package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.
var (
	name                  string = "testnet"                                                            // env."name"
	account_addr          string = "0x06f36e8a0fc06518125bbb1c63553e8a7d8597d437f9d56d891b8c7d3c977716" // Replace it with your account address
	account_cairo_version        = 0                                                                    // Replace  with the cairo version of your account
	privateKey            string = "0x0687bf84896ee63f52d69e6de1b41492abeadc0dc3cb7bd351d0a52116915937" // Replace it with your account private key
	public_key            string = "0x58b0824ee8480133cad03533c8930eda6888b3c5170db2f6e4f51b519141963"  // Replace it with your account public key
	someContract          string = "0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf" // UDC contract address
	contractMethod        string = "deployContract"                                                     // UDC method to deploy account (from pre-declared contract)
)

// Example succesful transaction created from this example on Goerli
// https://goerli.voyager.online/tx/0x9576bad061e1790ea1785cb3a950a5724390ea3d0bbb65fc09cc300d801b22

func main() {
	fmt.Println("Starting deployAccountUDC example")

	// Loading the env
	godotenv.Load(fmt.Sprintf(".env.%s", name))
	url := os.Getenv("INTEGRATION_BASE") //please modify the .env.testnet and replace the INTEGRATION_BASE with a starknet goerli RPC.

	// Initialize connection to RPC provider
	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	// Here we are converting the account address to felt
	account_address, err := utils.HexToFelt(account_addr)
	if err != nil {
		panic(err.Error())
	}
	// Initialize the account memkeyStore (set public and private keys)
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic(err.Error())
	}
	ks.Put(public_key, fakePrivKeyBI)

	fmt.Println("Established connection with the client")

	// Set maxFee
	maxfee, err := utils.HexToFelt("0x9184e72a000")
	if err != nil {
		panic(err.Error())
	}

	// Initialize the account
	accnt, err := account.NewAccount(clientv02, account_address, public_key, ks, account_cairo_version)
	if err != nil {
		panic(err.Error())
	}

	// Get the accounts nonce
	nonce, err := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	if err != nil {
		panic(err.Error())
	}

	// Build the InvokeTx struct
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        maxfee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: accnt.AccountAddress,
	}

	// Convert the contractAddress from hex to felt
	contractAddress, err := utils.HexToFelt(someContract)
	if err != nil {
		panic(err.Error())
	}

	// Build the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,                               //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod), //this is the function that we want to call
		Calldata:           getUDCCalldata(),
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		panic(err.Error())
	}

	// Sign the transaction
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

// getUDCCalldata is a simple helper to set the call data required by the UDCs deployContract function. Update as needed.
func getUDCCalldata() []*felt.Felt {

	classHash, err := new(felt.Felt).SetString("0x32f352d58c0a96d594de0ab19c24b9e6ed1e6310f805e61369ff156310827a")
	if err != nil {
		panic(err.Error())
	}

	randomInt := rand.Uint64()
	salt := new(felt.Felt).SetUint64(randomInt) // to prevent address clashes

	unique, err := new(felt.Felt).SetString("0x0")
	if err != nil {
		panic(err.Error())
	}

	calldataLen, err := new(felt.Felt).SetString("0x5")
	if err != nil {
		panic(err.Error())
	}

	calldata, err := utils.HexArrToFelt([]string{
		"0x477261696c7320455243343034",
		"0x475241494c53",
		"0x2710",
		"0x00",
		"0x07820b89733f802708f8eb768b59615f986205adc6eb6917c38b7771f7801caa"})
	if err != nil {
		panic(err.Error())
	}
	return append([]*felt.Felt{classHash, salt, unique, calldataLen}, calldata...)
}
