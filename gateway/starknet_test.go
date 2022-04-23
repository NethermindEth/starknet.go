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
)

func TestExecuteGoerli(t *testing.T) {
	curve, err := caigo.SC(caigo.WithConstants("./pedersen_params.json"))
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	priv, _ := new(big.Int).SetString("3904bd288b88a1dcd73e648b10642d63cb9b2ffd86526deee9d073f0690139e", 16)
	x, y, _ := curve.PrivateToPoint(priv)

	signer, err := curve.NewSigner(priv, x, y, NewClient())
	if err != nil {
		t.Errorf("Could not create signer: %v\n", err)
	}

	calls := []caigo.Transaction{
		{
			ContractAddress:    "0x07394cbe418daa16e42b87ba67372d4ab4a5df0b05c6e554d158458ce245bc10",
			EntryPointSelector: "mint",
			Calldata: []string{
				"3139631220741201955103162951941433790693684583007823827564831473435921360636",
				"1500000000000000000000",
				"0",
			},
		},
		{
			ContractAddress:    "0x07394cbe418daa16e42b87ba67372d4ab4a5df0b05c6e554d158458ce245bc10",
			EntryPointSelector: "transfer",
			Calldata: []string{
				"0x02e1b1ae589af66432469af22a38e84a6ac17202c55e3af2d40f8e18d3395398",
				"1500000000000000000000",
				"0",
			},
		},
	}

	_, err = signer.Execute(context.Background(), "0x6f0f7e2594028a454bed6bd856cc566763a6bef3d9965d79bd888ccea7426fc", calls)
	if err != nil {
		t.Errorf("Could not execute multicall with account: %v\n", err)
	}
}

func TestInvokeContract(t *testing.T) {
	gw := NewClient()

	req := caigo.Transaction{
		ContractAddress:    "0x077fd9aee87891eb334448c26e01020c8cffec0bf62a959bd373490542bdd812",
		EntryPointSelector: "increment",
	}

	_, err := gw.Invoke(context.Background(), req)
	if err != nil {
		t.Errorf("Could not add tx: %v\n", err)
	}
}

func TestLocalStarkNet(t *testing.T) {
	ctx := context.Background()
	setupTestEnvironment()

	curve, _ := caigo.SC()

	gw := NewClient(WithChain("local"))

	pr, _ := curve.GetRandomPrivateKey()

	deployRequest := caigo.DeployRequest{
		ContractAddressSalt: caigo.BigToHex(pr),
		ConstructorCalldata: []string{},
	}

	resp, err := gw.Deploy(ctx, "tmp/counter_compiled.json", deployRequest)
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

	_, err = gw.StorageAt(ctx, tx.Transaction.ContractAddress, 0, &StorageAtOptions{BlockNumber: 0})
	if err != nil {
		t.Errorf("Could not get storage: %v\n", err)
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
