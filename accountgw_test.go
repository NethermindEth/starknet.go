package caigo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo/artifacts"
	"github.com/dontpanicdao/caigo/gateway"
	devtest "github.com/dontpanicdao/caigo/test"
	"github.com/dontpanicdao/caigo/types"
)

const (
	SEED             int    = 100000000
	CONTRACT_ADDRESS string = "0x02b795d8c5e38c45da3b89c91174c66a3c77845bbeb87a36038f19c521dbe87e"
)

type TestAccountType struct {
	PrivateKey   string               `json:"private_key"`
	PublicKey    string               `json:"public_key"`
	Address      string               `json:"address"`
	Transactions []types.FunctionCall `json:"transactions,omitempty"`
}

func InstallCounterContract(t *testing.T, provider *gateway.GatewayProvider) string {
	class := types.ContractClass{}
	if err := json.Unmarshal(artifacts.CounterCompiled, &class); err != nil {
		t.Fatal("error Unmashaling counter contract", err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	tx, err := provider.Deploy(context.Background(), class, types.DeployRequest{
		ContractAddressSalt: "0x0",
		ConstructorCalldata: []string{},
	})
	if err != nil {
		t.Fatal("error deploying contract", err)
	}
	fmt.Println("deploy counter txHash", tx.TransactionHash)
	_, receipt, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 3, 20)
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		t.Fatal("contract installation did not succeed", err)
	}
	if !receipt.Status.IsTransactionFinal() ||
		receipt.Status == types.TransactionRejected {
		t.Fatal("contract installation status:", receipt.Status)
	}
	return tx.ContractAddress
}
func TestGatewayAccount_ExtimateAndExecute(t *testing.T) {
	testConfig := beforeGatewayEach(t)

	if testEnv != "devnet" {
		t.Skip("this test is only available on devnet")
	}
	counterAddress := InstallCounterContract(t, testConfig.client)
	type testSetType struct {
		Calls []types.FunctionCall `json:"transactions,omitempty"`
	}

	testSet := map[string][]testSetType{
		"devnet": {{
			Calls: []types.FunctionCall{{
				EntryPointSelector: "increment",
				ContractAddress:    types.HexToHash(counterAddress),
			}},
		}},
	}[testEnv]

	for _, test := range testSet {
		accounts, err := devtest.NewDevNet().Accounts()
		if err != nil {
			t.Fatal("should access the existing accounts", err)
		}
		account, err := NewGatewayAccount(
			accounts[0].PrivateKey,
			accounts[0].Address,
			testConfig.client,
			AccountVersion1)
		if err != nil {
			t.Fatal("should access the existing accounts", err)
		}
		if err != nil {
			t.Fatal("should access the existing accounts", err)
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		estimateFee, err := account.EstimateFee(ctx, test.Calls, types.ExecuteDetails{})
		if err != nil {
			t.Fatal("should succeed with EstimateFee, instead:", err)
		}
		fmt.Printf("estimate fee is %+v\n", estimateFee.OverallFee)
		tx, err := account.Execute(ctx, test.Calls, types.ExecuteDetails{})
		if err != nil {
			t.Fatal("should succeed with Execute, instead:", err)
		}
		fmt.Printf("Execute txHash: %v\n", tx.TransactionHash)
		_, state, err := testConfig.client.WaitForTransaction(ctx, tx.TransactionHash, 3, 10)
		if err != nil {
			t.Fatal("should succeed with Execute, instead:", err)
		}
		if state.Status != types.TransactionAcceptedOnL1 && state.Status != types.TransactionAcceptedOnL2 {
			t.Fatal("should be final, instead:", state.Status)
		}
		result, err := account.Call(ctx, types.FunctionCall{
			ContractAddress:    types.HexToHash(counterAddress),
			EntryPointSelector: "get_count",
		})
		if err != nil {
			t.Fatal("should succeed with Call, instead:", err)
		}
		if len(result) == 0 {
			t.Fatal("should return data, instead 0")
		}
		fmt.Println("count is now: ", result[0])
	}
}
