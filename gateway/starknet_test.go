package gateway

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

func TestContractAddresses(t *testing.T) {
	gw := NewClient()

	addresses, err := gw.ContractAddresses(context.Background())
	if err != nil {
		t.Errorf("Could not get starknet addresses: %v\n", err)
	}
	
	if addresses.Starknet != "0xde29d060D45901Fb19ED6C6e959EB22d8626708e" {
		t.Errorf("Fetched incorrect addresses: %v\n", err)
	}
}

func TestExecute(t *testing.T) {
	curve, err := caigo.SC(caigo.WithConstants("../pedersen_params.json"))
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	priv := "0x879d7dad7f9df54e1474ccf572266bba36d40e3202c799d6c477506647c126"
	addr := "0x126dd900b82c7fc95e8851f9c64d0600992e82657388a48d3c466553d4d9246"

	account, err := curve.NewAccount(priv, addr, NewProvider())
	if err != nil {
		t.Errorf("Could not create account: %v\n", err)
	}
	
	_, err = account.Execute(context.Background(), types.Transaction{
		ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
		EntryPointSelector: "update_struct_store",
		Calldata: []string{"435921360636", "1500000000000000000000", "0"},
	})
	if err != nil {
		t.Errorf("Could not execute multicall with account: %v\n", err)
	}
}

func TestExecuteMulti(t *testing.T) {
	curve, err := caigo.SC(caigo.WithConstants("../pedersen_params.json"))
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	priv := "0xefbd0ef595389f5b3d466f782cae093d99aaa5a312b1b99c4d6627318a3754"
	addr := "0x028105caf03e1c4eb96b1c18d39d9f03bd53e5d2affd0874792e5bf05f3e529f"

	account, err := curve.NewAccount(priv, addr, NewProvider())
	if err != nil {
		t.Errorf("Could not create account: %v\n", err)
	}
	
	calls := []types.Transaction{
		{
			ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
			EntryPointSelector: "update_single_store",
			Calldata: []string{"3"},
		},
		{
			ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
			EntryPointSelector: "update_multi_store",
			Calldata: []string{"4", "6"},
		},
	}

	_, err = account.ExecuteMultiCall(context.Background(), "0x1f2cefc44c8d", calls)
	if err != nil {
		t.Errorf("Could not execute multicall with account: %v\n", err)
	}
}

// requires starknet-devnet to be running and accessible on port 5000
// and seed for accounts to be specified to 0
// ex: starknet-devnet --port 5000 --seed 0
// (ref: https://github.com/Shard-Labs/starknet-devnet)
func TestLocalStarkNet(t *testing.T) {
	ctx := context.Background()
	setupTestEnvironment()

	curve, _ := caigo.SC()

	gw := NewClient(WithChain("local"))

	rand, _ := curve.GetRandomPrivateKey()
	deployRequest := types.DeployRequest{
		ContractAddressSalt: caigo.BigToHex(rand),
		ConstructorCalldata: []string{},
	}

	resp, err := gw.Deploy(ctx, "tmp/counter_compiled.json", deployRequest)
	if err != nil {
		t.Errorf("Could not deploy contract: %v\n", err)
	}
	
	depTx, err := gw.Transaction(ctx, TransactionOptions{TransactionHash: resp.TransactionHash})
	if err != nil || depTx.Status != "ACCEPTED_ON_L2" {
		t.Errorf("Could not get tx: %v\n", err)
	}

	// bug in starknet-devnet can only declare one class per devnet run
	resp, err = gw.Declare(ctx, "tmp/counter_compiled.json", types.DeclareRequest{})
	if err != nil {
		t.Errorf("Could not deploy contract: %v\n", err)
	}

	tx, err := gw.Transaction(ctx, TransactionOptions{TransactionHash: resp.TransactionHash})
	if err != nil || tx.Status != "ACCEPTED_ON_L2" {
		t.Errorf("Could not get tx: %v\n", err)
	}

	receipt, err := gw.TransactionReceipt(ctx, resp.TransactionHash)
	if err != nil || receipt.Status != "ACCEPTED_ON_L2" {
		t.Errorf("Could not get tx receipt: %v\n", err)
	}
	
	block, err := gw.Block(ctx, &BlockOptions{BlockHash: tx.BlockHash})
	if err != nil || block.Status != "ACCEPTED_ON_L2" {
		t.Errorf("Could not get block by hash: %v\n", err)
	}
}

func setupTestEnvironment() {
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		err := os.Mkdir("tmp", os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}

	url := "https://raw.githubusercontent.com/starknet-edu/ultimate-env/main/counter.cairo"
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = ioutil.WriteFile("tmp/counter.cairo", body, 0666)
	if err != nil {
		panic(err.Error())
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}

	err = exec.Command(fmt.Sprintf("%s/cairo_venv/bin/starknet-compile", dirname), "tmp/counter.cairo", "--output", "tmp/counter_compiled.json", "--abi", "tmp/counter_abi.json").Run()
	if err != nil {
		panic(err.Error())
	}
}
