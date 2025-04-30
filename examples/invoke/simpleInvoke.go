package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// simpleInvoke is a function that shows how to easily send an invoke transaction.
func simpleInvoke(accnt *account.Account, contractAddress *felt.Felt, contractMethod string, amount *felt.Felt) {
	u256Amount, err := utils.HexToU256Felt(amount.String())
	if err != nil {
		panic(err)
	}
	// Building the functionCall struct, where :
	FnCall := rpc.InvokeFunctionCall{
		ContractAddress: contractAddress, //contractAddress is the contract that we want to call
		FunctionName:    contractMethod,  //this is the function that we want to call
		CallData:        u256Amount,      //the calldata necessary to call the function. Here we are passing the "amount" value (a u256 cairo variable) for the "mint" function
	}

	// Building and sending the Broadcast Invoke Txn.
	//
	// note: in Starknet, you can execute multiple function calls in the same transaction, even if they are from different contracts.
	// To do this in Starknet.go, just group all the 'InvokeFunctionCall' in the same slice and pass it to BuildInvokeTxn.
	resp, err := accnt.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{FnCall}, 1.5, false)
	if err != nil {
		panic(err)
	}

	fmt.Println("Simple Invoke : Waiting for the transaction receipt...")

	txReceipt, err := accnt.WaitForTransactionReceipt(context.Background(), resp.Hash, time.Second)
	if err != nil {
		panic(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Simple Invoke : Transaction hash response: %v\n", resp.Hash)
	fmt.Printf("Simple Invoke : Transaction execution status: %s\n", txReceipt.ExecutionStatus)
	fmt.Printf("Simple Invoke : Transaction status: %s\n", txReceipt.FinalityStatus)
	fmt.Printf("Simple Invoke : Block number: %d\n", txReceipt.BlockNumber)
}
