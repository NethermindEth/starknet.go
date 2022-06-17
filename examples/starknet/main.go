package main

import (
	"context"
	"fmt"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

func main() {
	// init the stark curve with constants
	// 'WithConstants()' will pull the StarkNet 'pedersen_params.json' file if you don't have it locally
	curve, err := caigo.SC(caigo.WithConstants())
	if err != nil {
		panic(err.Error())
	}

	// init starknet gateway client
	gw := gateway.NewClient() //defaults to goerli

	// get random value for salt
	priv, _ := curve.GetRandomPrivateKey()

	// example: https://github.com/starknet-edu/ultimate-env/blob/main/counter.cairo
	// starknet-compile counter.cairo --output counter_compiled.json --abi counter_abi.json
	deployResponse, err := gw.Deploy(context.Background(), "counter_compiled.json", types.DeployRequest{
		ContractAddressSalt: caigo.BigToHex(priv),
		ConstructorCalldata: []string{},
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deployment Response: \n\t%+v\n\n", deployResponse)

	// poll until the desired transaction status
	pollInterval := 5
	n, status, err := gw.PollTx(context.Background(), deployResponse.TransactionHash, types.ACCEPTED_ON_L2, pollInterval, 150)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) status: %s\n\n", n*pollInterval, n, deployResponse.TransactionHash, status)

	// fetch transaction details
	tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: deployResponse.TransactionHash})
	if err != nil {
		panic(err.Error())
	}

	// call StarkNet contract
	callResp, err := gw.Call(context.Background(), types.Transaction{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: "get_count",
	}, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])

	// invoke StarkNet contract external function
	invResp, err := gw.Invoke(context.Background(), types.Transaction{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: "increment",
	})
	if err != nil {
		panic(err.Error())
	}

	n, status, err = gw.PollTx(context.Background(), invResp.TransactionHash, types.ACCEPTED_ON_L2, 5, 150)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) status: %s\n\n", n*pollInterval, n, deployResponse.TransactionHash, status)

	callResp, err = gw.Call(context.Background(), types.Transaction{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: "get_count",
	}, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
}
