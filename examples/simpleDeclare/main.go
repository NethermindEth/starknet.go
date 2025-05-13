package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

const (
	sierraContractFilePath = "./contract.sierra.json"
	casmContractFilePath   = "./contract.casm.json"
)

// This example demonstrates how to declare a contract on Starknet.
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
	accnt, err := account.NewAccount(client, accountAddressInFelt, publicKey, ks, accountCairoVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("Established connection with the client")

	// Unmarshalling the casm contract class from a JSON file.
	casmClass, err := utils.UnmarshalJSONFileToType[contracts.CasmClass](casmContractFilePath, "")
	if err != nil {
		panic(err)
	}

	// Unmarshalling the sierra contract class from a JSON file.
	contractClass, err := utils.UnmarshalJSONFileToType[contracts.ContractClass](sierraContractFilePath, "")
	if err != nil {
		panic(err)
	}

	// Building and sending the Broadcast Invoke Txn.
	//
	// note: in Starknet, you can execute multiple function calls in the same transaction, even if they are from different contracts.
	// To do this in Starknet.go, just group all the 'InvokeFunctionCall' in the same slice and pass it to BuildInvokeTxn.
	resp, err := accnt.BuildAndSendDeclareTxn(
		context.Background(),
		casmClass,
		contractClass,
		&utils.TransactionOptions{Multiplier: 1.5, WithQueryBitVersion: false, Tip: "0x0"},
	)
	if err != nil {
		if strings.Contains(err.Error(), "is already declared") {
			fmt.Println("")
			fmt.Println("Error: ooops, this contract class was already declared.")
			fmt.Println("You need to: ")
			fmt.Println("- create a different Cairo contract,")
			fmt.Println("- compile it,")
			fmt.Println("- paste the new casm and sierra json files in this 'examples/simpleDeclare' folder,")
			fmt.Println("- change the 'casmContractFilePath' and 'sierraContractFilePath' variables to the new files names,")
			fmt.Println("and then, run the example again. You can use Scarb for it: https://docs.swmansion.com/scarb/")

			return
		}
		panic(err)
	}

	fmt.Println("Waiting for the transaction status...")

	txReceipt, err := accnt.WaitForTransactionReceipt(context.Background(), resp.Hash, time.Second)
	if err != nil {
		panic(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Transaction hash response: %v\n", resp.Hash)
	fmt.Printf("Transaction execution status: %s\n", txReceipt.ExecutionStatus)
	fmt.Printf("Transaction status: %s\n", txReceipt.FinalityStatus)
	fmt.Printf("Class hash: %s\n", resp.ClassHash)
}
