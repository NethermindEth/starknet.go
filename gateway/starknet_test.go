package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

const (
	FEE_MARGIN         uint64 = 115
	SEED               int    = 100000000
	ACCOUNT_CLASS_HASH string = "0x3e327de1c40540b98d05cbcb13552008e36f0ec8d61d46956d2f9752c294328"
	CONTRACT_ADDRESS   string = "0x02b795d8c5e38c45da3b89c91174c66a3c77845bbeb87a36038f19c521dbe87e"
)

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
	accountCompiled string = "contracts/account_class.json"
)

type TestAccountType struct {
	PrivateKey   string               `json:"private_key"`
	PublicKey    string               `json:"public_key"`
	Address      string               `json:"address"`
	Transactions []types.FunctionCall `json:"transactions,omitempty"`
}

func TestDeclare(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}
	testSet := map[string][]testSetType{
		"devnet":  {{}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{}},
	}[testEnv]

	for _, env := range testSet {
		gw := testConfig.client
		declareTx, err := gw.Declare(context.Background(), accountCompiled, types.DeclareRequest{})
		if err != nil {
			t.Errorf("%s: could not 'DECLARE' contract: %v\n", env, err)
			return
		}

		tx, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: declareTx.TransactionHash})
		if err != nil {
			t.Errorf("%s: could not get 'DECLARE' transaction: %v\n", env, err)
		}
		if tx.Transaction.Type != DECLARE {
			t.Errorf("%s: incorrect delcare transaction: %v\n", env, tx)
		}
	}
}

func TestDeployCounterContract(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}
	testSet := map[string][]testSetType{
		"devnet":  {{}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{}},
	}[testEnv]

	for range testSet {

		gw := testConfig.client

		tx, err := gw.Deploy(context.Background(), "contracts/counter.json", types.DeployRequest{
			ContractAddressSalt: "0x1",
			ConstructorCalldata: []string{},
		})
		if err != nil {
			t.Errorf("testnet: could not deploy contract: %v\n", err)
		}
		fmt.Println("txHash: ", tx.TransactionHash)
	}

}

func TestCallGoerli(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		testnetAccounts []TestAccountType
	}
	testSet := map[string][]testSetType{
		"devnet":  {{}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			testnetAccounts: testnetAccounts,
		}},
	}[testEnv]

	for _, env := range testSet {
		gw := testConfig.client
		for _, testAccount := range env.testnetAccounts {
			call := types.FunctionCall{
				ContractAddress:    types.HexToHash(testAccount.Address),
				EntryPointSelector: "get_public_key",
			}

			resp, err := gw.Call(context.Background(), call, "")
			if err != nil {
				t.Errorf("testnet: could 'Call' deployed contract: %v\n", err)
			}
			if len(resp) == 0 {
				t.Errorf("testnet: could get signing key for account: %v\n", err)
			}

			if resp[0] != testAccount.PublicKey {
				t.Errorf("testnet: signing key is incorrect: \n%s %v\n", resp[0], testAccount.PublicKey)
			}
		}
	}
}

func DevnetAccounts() ([]TestAccountType, error) {
	req, err := http.NewRequest("GET", "http://localhost:5050/predeployed_accounts", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accounts []TestAccountType
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	return accounts, err
}
