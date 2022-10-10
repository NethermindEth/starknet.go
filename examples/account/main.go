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
	name            string = "devnet"
	counterContract string = "0x002bb8640f7875a6c73469149b261c99897f33aea38f272b5d8644c6c67bd278"
	address         string = "0x036f1f6a2c7f7b5180164a7f51a7e46d63bf5b898bd5144375858d4453725309"
	privakeKey      string = "0x2623ddff7889c38d0504dbe90d655402376601a73f7fea9864d8d302e8ba34d"
	feeMargin       uint64 = 115
	maxPoll         int    = 5
	pollInterval    int    = 3
)

func main() {
	// init starknet gateway client
	gw := gateway.NewClient(gateway.WithChain(name), gateway.WithBaseURL("http://localhost:8080"))

	// get count before tx
	// callResp, err := gw.Call(context.Background(), types.FunctionCall{
	// 	ContractAddress:    types.HexToHash(counterContract),
	// 	EntryPointSelector: "get_count",
	// }, "")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("Counter is currently at: ", callResp[0])

	// init account handler
	account, err := caigo.NewGatewayAccount(privakeKey, address, gw)
	if err != nil {
		panic(err.Error())
	}

	increment := []types.FunctionCall{
		{
			ContractAddress:    types.HexToHash(counterContract),
			EntryPointSelector: "increment",
		},
	}

	// estimate fee for executing transaction
	feeEstimate, err := account.EstimateFee(context.Background(), increment, types.ExecuteDetails{})
	if err != nil {
		panic(err.Error())
	}
	fee, _ := big.NewInt(0).SetString(string(feeEstimate.OverallFee), 0)
	expandedFee := big.NewInt(0).Mul(fee, big.NewInt(int64(feeMargin)))
	max := big.NewInt(0).Div(expandedFee, big.NewInt(100))
	fmt.Printf("Fee:\n\tEstimate\t\t%v wei\n\tEstimate+Margin\t\t%v wei\n\n", feeEstimate.OverallFee, max)

	// execute transaction
	execResp, err := account.Execute(context.Background(), increment, types.ExecuteDetails{MaxFee: max})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", execResp)
	n, receipt, err := gw.WaitForTransaction(context.Background(), execResp.TransactionHash, pollInterval, maxPoll)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) receipt: %s\n\n", n*pollInterval, n, execResp.TransactionHash, receipt.Status)

	// get count after tx
	callResp, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    types.HexToHash(counterContract),
		EntryPointSelector: "get_count",
	}, "")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
}
