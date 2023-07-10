package starknet.go

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/NethermindEth/starknet.go/types"
)

type TestAccountType struct {
	PrivateKey   string               `json:"private_key"`
	PublicKey    string               `json:"public_key"`
	Address      string               `json:"address"`
	Transactions []types.FunctionCall `json:"transactions,omitempty"`
}

func TestGatewayAccount_EstimateAndExecute(t *testing.T) {
	testConfig := beforeGatewayEach(t)
	type testSetType struct {
		ExecuteCalls []types.FunctionCall
		QueryCall    types.FunctionCall
	}

	testSet := map[string][]testSetType{
		"devnet": {{
			ExecuteCalls: []types.FunctionCall{{
				EntryPointSelector: "increment",
				ContractAddress:    types.StrToFelt(testConfig.CounterAddress),
			}},
			QueryCall: types.FunctionCall{
				EntryPointSelector: "get_count",
				ContractAddress:    types.StrToFelt(testConfig.CounterAddress),
			},
		}},
		"testnet": {{
			ExecuteCalls: []types.FunctionCall{{
				EntryPointSelector: "increment",
				ContractAddress:    types.StrToFelt(testConfig.CounterAddress),
			}},
			QueryCall: types.FunctionCall{
				EntryPointSelector: "get_count",
				ContractAddress:    types.StrToFelt(testConfig.CounterAddress),
			},
		}},
	}[testEnv]

	for _, test := range testSet {
		// shim a keystore into existing tests.
		ks := NewMemKeystore()
		fakeSenderAddress := testConfig.AccountPrivateKey
		k := types.SNValToBN(testConfig.AccountPrivateKey)
		ks.Put(fakeSenderAddress, k)
		account, err := NewGatewayAccount(
			types.StrToFelt(fakeSenderAddress),
			types.StrToFelt(testConfig.AccountAddress),
			ks,
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
		fmt.Printf("estimate fee is %+v\n", estimateFee)
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
