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

// verboseInvoke is a function that shows how to send an invoke transaction step by step, using only
// a few helper functions.
func verboseInvoke(accnt *account.Account, contractAddress *felt.Felt, contractMethod string, amount *felt.Felt) {
	// Getting the nonce from the account
	nonce, err := accnt.Nonce(context.Background())
	if err != nil {
		panic(err)
	}

	u256Amount, err := utils.HexToU256Felt(amount.String())
	if err != nil {
		panic(err)
	}
	// Building the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,                               //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod), //this is the function that we want to call
		Calldata:           u256Amount,                                    //the calldata necessary to call the function. Here we are passing the "amount" value (a u256 cairo variable) for the "mint" function
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	//
	// note: in Starknet, you can execute multiple function calls in the same transaction, even if they are from different contracts.
	// To do this in Starknet.go, just group all the function calls in the same slice and pass it to FmtCalldata
	// e.g. : InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{funcCall, anotherFuncCall, yetAnotherFuncCallFromDifferentContract})
	calldata, err := accnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		panic(err)
	}

	// Using the BuildInvokeTxn helper to build the BroadInvokeTx
	InvokeTx := utils.BuildInvokeTxn(accnt.Address, nonce, calldata, rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L1DataGas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L2Gas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
	})

	// We need to sign the transaction to be able to estimate the fee
	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV3)
	if err != nil {
		panic(err)
	}

	// Estimate the transaction fee
	feeRes, err := accnt.Provider.EstimateFee(context.Background(), []rpc.BroadcastTxn{InvokeTx}, []rpc.SimulationFlag{}, rpc.WithBlockTag("pending"))
	if err != nil {
		panic(err)
	}

	// assign the estimated fee to the transaction, multiplying the estimated fee by 1.5 for a better chance of success
	InvokeTx.InvokeTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(feeRes[0], 1.5)

	// As we changed the resource bounds, we need to sign the transaction again, since the resource bounds are part of the signature
	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV3)
	if err != nil {
		panic(err)
	}

	// After signing, we finally send the transaction in order to invoke the contract function
	resp, err := accnt.SendTransaction(context.Background(), InvokeTx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Verbose Invoke : Waiting for the transaction receipt...")

	txReceipt, err := accnt.WaitForTransactionReceipt(context.Background(), resp.Hash, time.Second)
	if err != nil {
		panic(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Verbose Invoke : Transaction hash response: %v\n", resp.Hash)
	fmt.Printf("Verbose Invoke : Transaction execution status: %s\n", txReceipt.ExecutionStatus)
	fmt.Printf("Verbose Invoke : Transaction status: %s\n", txReceipt.FinalityStatus)
	fmt.Printf("Verbose Invoke : Block number: %d\n", txReceipt.BlockNumber)
}
