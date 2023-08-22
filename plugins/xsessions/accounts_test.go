package xsessions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/rpc"
	ctypes "github.com/NethermindEth/starknet.go/types"
)

const (
	privateKey        = "0x1"
	sessionPrivateKey = "0x2"
	counterAddress    = "0x07704fb2d72fcdae1e6f658ef8521415070a01a3bd3cc5788f7b082126922b7b"
	devnetEth         = "0x62230ea046a9a5fbc261ac77d03c8d41e5d442db2284587570ab46455fd2488"
)

type accountPlugin struct {
	PluginHash     string `json:"pluginHash"`
	AccountAddress string `json:"accountAddress"`
}

func (ap *accountPlugin) Read(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	json.Unmarshal(content, ap)
	return nil
}

func (ap *accountPlugin) Write(filename string) error {
	content, err := json.Marshal(ap)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0664)
}

var accountCompiled = artifacts.AccountV0WithPluginCompiled
var counterCompiled = artifacts.CounterCompiled

// TestCounter_DeployContract
func TestCounter_DeployContract(t *testing.T) {
	provider := beforeEachRPC(t)

	counterClass := rpc.DepcreatedContractClass{}
	inputs := []string{}

	if err := json.Unmarshal(counterCompiled, &counterClass); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tx, err := provider.AddDeployTransaction(ctx, "0xdeadbeef", inputs, counterClass)
	if err != nil {
		t.Fatal("deploy should succeed, instead:", err)
	}
	if tx.ContractAddress.String() != counterAddress {
		t.Fatal("deploy should return counter address, instead:", tx.ContractAddress)
	}
	status, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 8*time.Second)
	if err != nil {
		t.Fatal("declare should succeed, instead:", err)
	}
	if status != ctypes.TransactionAcceptedOnL2 {
		t.Log("unexpected status transaction status, check:", status)
		t.Log("...")
		t.Log("   verify transaction")
		t.Log("...")
		t.Log("export STARKNET_WALLET=starkware.starknet.wallets.open_zeppelin.OpenZeppelinAccount")
		t.Log("export STARKNET_NETWORK=alpha-goerli")
		t.Logf("export HASH=%s\n", tx.TransactionHash)
		t.Log("starknet get_transaction --hash $HASH --feeder_gateway http://localhost:5050/feeder_gateway")
		t.Log("...")
		t.Fail()
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
}
