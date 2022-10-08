package caigo

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo/gateway"
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

var (
	devnetAccounts  []TestAccountType
	testnetAccounts = []TestAccountType{
		{
			PrivateKey: "0x2294a8695b61f3a7ae8ddcb2cdfa72f1973dbeb22955aa43286a57685aa0e91",
			PublicKey:  "0x4672cbb8f57ff12043861effdb7abc21eb81b8a1473868d91bb0681c7e4f269",
			Address:    "0x1343858d3b9315df9155106c29103102e893252ded58884098be03060da347f",
			Transactions: []types.FunctionCall{
				{
					ContractAddress:    types.HexToHash(CONTRACT_ADDRESS),
					EntryPointSelector: "increase_balance",
					Calldata: []string{
						"1",
					},
				},
			},
		},
	}
)

func TestExecuteGoerli(t *testing.T) {
	testConfig := beforeGatewayEach(t)
	for _, testAccount := range testnetAccounts {
		account, err := NewGatewayAccount(testAccount.PrivateKey, types.HexToHash(testAccount.Address), testConfig.client)
		if err != nil {
			t.Errorf("testnet: could not create account: %v\n", err)
		}

		feeEstimate, err := account.EstimateFee(context.Background(), testAccount.Transactions, types.ExecuteDetails{})
		if err != nil {
			t.Errorf("testnet: could not estimate fee for transaction: %v\n", err)
		}

		fee, _ := big.NewInt(0).SetString(string(feeEstimate.OverallFee), 0)
		expandedFee := big.NewInt(0).Mul(fee, big.NewInt(int64(FEE_MARGIN)))
		max := big.NewInt(0).Div(expandedFee, big.NewInt(100))

		_, err = account.Execute(context.Background(), testAccount.Transactions,
			types.ExecuteDetails{
				MaxFee: max,
			})
		if err != nil {
			t.Errorf("Could not execute test transaction: %v\n", err)
		}
	}
}

func TestE2EDevnet(t *testing.T) {
	testConfig := beforeGatewayEach(t)

	type testSetType struct{}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for _, env := range testSet {
		gw := testConfig.client

		deployTx, err := gw.Deploy(context.Background(), "../rpc/tests/counter.json", types.DeployRequest{})
		if err != nil {
			t.Errorf("%s: could not deploy devnet counter: %v", env, err)
		}

		_, _, err = gw.PollTx(context.Background(), deployTx.TransactionHash, gateway.ACCEPTED_ON_L2, 1, 10)
		if err != nil {
			t.Errorf("%s: could not deploy devnet counter: %v", env, err)
		}

		txDetails, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: deployTx.TransactionHash})
		if err != nil {
			t.Errorf("%s: fetching transaction: %v", env, err)
		}

		for i := 0; i < 3; i++ {
			rand := fmt.Sprintf("0x%x", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(SEED))

			tx := []types.FunctionCall{
				{
					ContractAddress:    types.HexToHash(txDetails.Transaction.ContractAddress),
					EntryPointSelector: "set_rand",
					Calldata:           []string{rand},
				},
			}

			account, err := NewGatewayAccount(devnetAccounts[i].PrivateKey, types.HexToHash(devnetAccounts[i].Address), gw)
			if err != nil {
				t.Errorf("testnet: could not create account: %v\n", err)
			}

			feeEstimate, err := account.EstimateFee(context.Background(), tx, types.ExecuteDetails{})
			if err != nil {
				t.Errorf("testnet: could not estimate fee for transaction: %v\n", err)
			}
			fee, _ := big.NewInt(0).SetString(string(feeEstimate.OverallFee), 0)
			expandedFee := big.NewInt(0).Mul(fee, big.NewInt(int64(FEE_MARGIN)))
			max := big.NewInt(0).Div(expandedFee, big.NewInt(100))

			nonce, err := gw.AccountNonce(context.Background(), account.Address)
			if err != nil {
				t.Errorf("testnet: could not get account nonce: %v", err)
			}

			execResp, err := account.Execute(context.Background(), tx,
				types.ExecuteDetails{
					MaxFee: max,
					Nonce:  nonce,
				})
			if err != nil {
				t.Errorf("Could not execute test transaction: %v\n", err)
			}

			_, _, err = gw.PollTx(context.Background(), execResp.TransactionHash, gateway.ACCEPTED_ON_L2, 1, 10)
			if err != nil {
				t.Errorf("could not deploy devnet counter: %v\n", err)
			}

			call := types.FunctionCall{
				ContractAddress:    types.HexToHash(txDetails.Transaction.ContractAddress),
				EntryPointSelector: "get_rand",
			}
			callResp, err := gw.Call(context.Background(), call, "")
			if err != nil {
				t.Errorf("could not call counter contract: %v\n", err)
			}

			if rand != callResp[0] {
				t.Errorf("could not set value on counter contract: %v\n", err)
			}
		}
	}
}
