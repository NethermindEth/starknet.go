package gateway

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

func TestExecuteGoerli(t *testing.T) {
	curve, err := caigo.SC(caigo.WithConstants("../pedersen_params.json"))
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	priv, _ := new(big.Int).SetString("879d7dad7f9df54e1474ccf572266bba36d40e3202c799d6c477506647c126", 16)
	addr := "0x126dd900b82c7fc95e8851f9c64d0600992e82657388a48d3c466553d4d9246"

	signer, err := curve.NewSigner(priv, addr, NewProvider())
	if err != nil {
		t.Errorf("Could not create signer: %v\n", err)
	}

	tx := types.Transaction{
		ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
		EntryPointSelector: "update_struct_store",
		Calldata: []string{"435921360636", "1500000000000000000000", "0"},
	}
	
	// fee, err := signer.Provider.EstimateFee(context.Background(), tx)

	// calls := []types.Transaction{
	// 	{
	// 		ContractAddress:    "0x07394cbe418daa16e42b87ba67372d4ab4a5df0b05c6e554d158458ce245bc10",
	// 		EntryPointSelector: "mint",
	// 		Calldata: []string{
	// 			"3139631220741201955103162951941433790693684583007823827564831473435921360636",
	// 			"1500000000000000000000",
	// 			"0",
	// 		},
	// 	},
	// 	{
	// 		ContractAddress:    "0x07394cbe418daa16e42b87ba67372d4ab4a5df0b05c6e554d158458ce245bc10",
	// 		EntryPointSelector: "transfer",
	// 		Calldata: []string{
	// 			"0x02e1b1ae589af66432469af22a38e84a6ac17202c55e3af2d40f8e18d3395398",
	// 			"1500000000000000000000",
	// 			"0",
	// 		},
	// 	},
	// }

	resp, err := signer.ExecuteSingle(context.Background(), tx)
	fmt.Println("RESP: ", resp)
	if err != nil {
		t.Errorf("Could not execute multicall with account: %v\n", err)
	}
}

// requires starknet-devnet to be running and accessible on port 5000
// (ref: https://github.com/Shard-Labs/starknet-devnet)
func TestLocalStarkNet(t *testing.T) {
	ctx := context.Background()
	setupTestEnvironment()

	curve, _ := caigo.SC()

	gw := NewClient(WithChain("local"))

	pr, _ := curve.GetRandomPrivateKey()

	deployRequest := types.DeployRequest{
		ContractAddressSalt: caigo.BigToHex(pr),
		ConstructorCalldata: []string{},
	}

	resp, err := gw.Deploy(ctx, "tmp/counter_compiled.json", deployRequest)
	if err != nil {
		t.Errorf("Could not deploy contract: %v\n", err)
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
