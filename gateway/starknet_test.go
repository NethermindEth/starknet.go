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
	accountCompiled string = "contracts/account_class.json"
)

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
		declareTx, err := gw.Declare(context.Background(), accountCompiled, DeclareRequest{})
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
		"devnet": {{
			testnetAccounts: []TestAccountType{
				{
					PrivateKey: "0x28a778906e0b5f4d240ad25c5993422e06769eb799483ae602cc3830e3f538",
					PublicKey:  "0x63f0f116c78146e1e4e193923fe3cad5f236c0ed61c2dc04487a733031359b8",
					Address:    "0x0254cfb85c43dee6f410867b9795b5309beb4a2640211c8f5b2c7681a47e5f3c",
					Transactions: []types.FunctionCall{
						{
							ContractAddress:    types.HexToHash("0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62"),
							EntryPointSelector: "update_single_store",
							Calldata:           []string{"3"},
						},
					},
				},
				{
					PrivateKey: "0x879d7dad7f9df54e1474ccf572266bba36d40e3202c799d6c477506647c126",
					PublicKey:  "0xb95246e1caeaf34672906d7b74bd6968231a2130f41e85aebb62d43b88068",
					Address:    "0x0126dd900b82c7fc95e8851f9c64d0600992e82657388a48d3c466553d4d9246",
					Transactions: []types.FunctionCall{
						{
							ContractAddress:    types.HexToHash("0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62"),
							EntryPointSelector: "update_multi_store",
							Calldata:           []string{"4", "7"},
						},
						{
							ContractAddress:    types.HexToHash("0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62"),
							EntryPointSelector: "update_struct_store",
							Calldata:           []string{"435921360636", "15000000000000000000", "0"},
						},
					},
				},
			},
		}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			testnetAccounts: []TestAccountType{
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
			},
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

