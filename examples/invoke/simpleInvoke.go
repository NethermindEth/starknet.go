package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

func simpleInvoke(accnt *account.Account, contractAddress *felt.Felt, contractMethod string, amount *felt.Felt) {
	// Building the functionCall struct, where :
	FnCall := rpc.InvokeFunctionCall{
		ContractAddress: contractAddress,                  //contractAddress is the contract that we want to call
		FunctionName:    contractMethod,                   //this is the function that we want to call
		CallData:        []*felt.Felt{amount, &felt.Zero}, //the calldata necessary to call the function. Here we are passing the "amount" value for the "mint" function
	}

	// Building and sending the Broadcast Invoke Txn.
	//
	// note: in Starknet, you can execute multiple function calls in the same transaction, even if they are from different contracts.
	// To do this in Starknet.go, just group all the 'InvokeFunctionCall' in the same slice and pass it to BuildInvokeTxn.
	resp, err := accnt.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{FnCall}, 1.5)
	if err != nil {
		setup.PanicRPC(err)
	}

	fmt.Println("Waiting for the transaction status...")
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := accnt.Provider().GetTransactionStatus(context.Background(), resp.TransactionHash)
	if err != nil {
		setup.PanicRPC(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Transaction hash response: %v\n", resp.TransactionHash)
	fmt.Printf("Transaction execution status: %s\n", txStatus.ExecutionStatus)
	fmt.Printf("Transaction status: %s\n", txStatus.FinalityStatus)

}
