package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

var (
	name            string = "testnet"
	counterContract string = "0x0331034cbde9af8aef62929b5886b096ae3d11e33e6ca23122669e928d406500"
	address         string = "0x126dd900b82c7fc95e8851f9c64d0600992e82657388a48d3c466553d4d9246"
	privakeKey      string = "0x879d7dad7f9df54e1474ccf572266bba36d40e3202c799d6c477506647c126"
	feeMargin       uint64 = 115
	maxPoll         int    = 5
	pollInterval    int    = 150
)

func main() {
	// init starknet gateway client
	gw := gateway.NewProvider(gateway.WithChain(name))

	// get count before tx
	callResp, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    counterContract,
		EntryPointSelector: "get_count",
	}, "")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])

	// init account handler
	account, err := caigo.NewAccount(privakeKey, address, gw)
	if err != nil {
		panic(err.Error())
	}

	increment := []types.Transaction{
		{
			ContractAddress:    counterContract,
			EntryPointSelector: "increment",
		},
	}

	// estimate fee for executing transaction
	feeEstimate, err := account.EstimateFee(context.Background(), increment, caigo.ExecuteDetails{})
	if err != nil {
		panic(err.Error())
	}
	fee := types.Felt{
		Int: new(big.Int).SetUint64((feeEstimate.OverallFee * feeMargin) / 100),
	}
	fmt.Printf("Fee:\n\tEstimate\t\t%v wei\n\tEstimate+Margin\t\t%v wei\n\n", feeEstimate.OverallFee, fee)

	// execute transaction
	execResp, err := account.Execute(context.Background(), increment, caigo.ExecuteDetails{MaxFee: &fee})
	if err != nil {
		panic(err.Error())
	}

	n, receipt, err := gw.PollTx(context.Background(), execResp.TransactionHash, types.ACCEPTED_ON_L2, pollInterval, maxPoll)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) receipt: %s\n\n", n*pollInterval, n, execResp.TransactionHash, receipt.Status)

	// get count after tx
	callResp, err = gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    counterContract,
		EntryPointSelector: "get_count",
	}, "")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
}
