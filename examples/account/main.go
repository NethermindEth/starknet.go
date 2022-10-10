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
	counterContract string = "0x002bb8640f7875a6c73469149b261c99897f33aea38f272b5d8644c6c67bd278"
	address         string = "0x3ce251e4c470648a913346c218bd5f7925560d1811a2e49179958f13d03ffa6"
	privakeKey      string = "0x2d363b753bdfe5dfc77f4818bad6e25b"
	feeMargin       uint64 = 115
	maxPoll         int    = 5
	pollInterval    int    = 150
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
