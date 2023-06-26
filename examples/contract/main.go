package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/caigo/contracts"
	"github.com/smartcontractkit/caigo/gateway"

	"github.com/smartcontractkit/caigo/types"
)

// Start Devnet:
//   - starknet-devnet
var (
	name         string = "local"
	maxPoll      int    = 5
	pollInterval int    = 5
)

func main() {
	// init starknet gateway client
	gw := gateway.NewClient(gateway.WithChain(name))

	counterClass := types.ContractClass{}
	err := json.Unmarshal(contracts.CounterCompiled, &counterClass)
	if err != nil {
		panic(err.Error())
	}

	// will fail w/o new seed
	deployResponse, err := gw.Deploy(context.Background(), counterClass, types.DeployRequest{
		ContractAddressSalt: fmt.Sprintf("0x0%x", time.Now().UnixNano()),
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deployment Response: \n\t%+v\n\n", deployResponse)

	// poll until the desired transaction status
	n, receipt, err := gw.WaitForTransaction(context.Background(), deployResponse.TransactionHash, pollInterval, maxPoll)
	if err != nil {
		fmt.Println("Transaction Failure: ", receipt.Status)
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) status: %s\n\n", n*pollInterval, n, deployResponse.TransactionHash, receipt.Status)

	// fetch transaction data
	tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: deployResponse.TransactionHash})
	if err != nil {
		panic(err.Error())
	}

	// call StarkNet contract
	callResp, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    types.StrToFelt(tx.Transaction.ContractAddress),
		EntryPointSelector: "get_count",
	}, "")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
}
