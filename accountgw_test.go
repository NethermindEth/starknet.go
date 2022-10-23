package caigo

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestGatewayAccount_ExtimateAndExecute(t *testing.T) {
	testConfig := beforeGatewayEach(t)

	if testEnv != "devnet" {
		t.Skip("this test is only available on devnet")
	}
	type testSetType struct {
		ExecuteCalls []types.FunctionCall
		QueryCall    types.FunctionCall
	}

	testSet := map[string][]testSetType{
		"devnet": {{
			ExecuteCalls: []types.FunctionCall{{
				EntryPointSelector: "increment",
				ContractAddress:    types.HexToHash(testConfig.CounterAddress),
			}},
			QueryCall: types.FunctionCall{
				EntryPointSelector: "get_count",
				ContractAddress:    types.HexToHash(testConfig.CounterAddress),
			},
		}},
	}[testEnv]

	for _, test := range testSet {
		account, err := NewGatewayAccount(
			testConfig.AccountPrivateKey,
			testConfig.AccountAddress,
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
		estimateFee, err := account.EstimateFee(ctx, test.ExecuteCalls, types.ExecuteDetails{})
		if err != nil {
			t.Fatal("should succeed with EstimateFee, instead:", err)
		}
		fmt.Printf("estimate fee is %+v\n", estimateFee.OverallFee)
		tx, err := account.Execute(ctx, test.ExecuteCalls, types.ExecuteDetails{})
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
		result, err := account.Call(ctx, test.QueryCall)
		if err != nil {
			t.Fatal("should succeed with Call, instead:", err)
		}
		if len(result) == 0 {
			t.Fatal("should return data, instead 0")
		}
		fmt.Println("count is now: ", result[0])
	}
}
