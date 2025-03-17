package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

// TODO: improve this example

// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.

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

	casmClass, err := utils.UnmarshallJSONFileToType[contracts.CasmClass]("./contracts_v2_HelloStarknet.casm.json", "")
	if err != nil {
		panic(err)
	}

	contractClass, err := utils.UnmarshallJSONFileToType[rpc.ContractClass]("./contracts_v2_HelloStarknet.sierra.json", "")
	if err != nil {
		panic(err)
	}

	// Building and sending the Broadcast Invoke Txn.
	//
	// note: in Starknet, you can execute multiple function calls in the same transaction, even if they are from different contracts.
	// To do this in Starknet.go, just group all the 'InvokeFunctionCall' in the same slice and pass it to BuildInvokeTxn.
	resp, err := accnt.BuildAndSendDeclareTxn(context.Background(), *casmClass, contractClass, 1.5)
	if err != nil {
		setup.PanicRPC(err)
	}

	fmt.Println("Waiting for the transaction status...")
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := client.GetTransactionStatus(context.Background(), resp.TransactionHash)
	if err != nil {
		setup.PanicRPC(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Transaction hash response: %v\n", resp.TransactionHash)
	fmt.Printf("Transaction execution status: %s\n", txStatus.ExecutionStatus)
	fmt.Printf("Transaction status: %s\n", txStatus.FinalityStatus)

}
