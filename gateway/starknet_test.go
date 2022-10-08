package gateway

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
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
	//go:embed contracts/counter.json
	counterCompiled []byte

	counterAddress = "0x0"

	//go:embed contracts/account_class.json
	accountCompiled []byte
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
		accountClass := types.ContractClass{}
		err := json.Unmarshal(accountCompiled, &accountClass)
		if err != nil {
			t.Fatalf("could not parse contract: %v\n", err)
		}
		declareTx, err := gw.Declare(context.Background(), accountClass, DeclareRequest{})
		if err != nil {
			t.Errorf("%s: could not 'DECLARE' contract: %v\n", env, err)
			return
		}

		tx, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: declareTx.TransactionHash})
		if err != nil {
			t.Errorf("%s: could not get 'DECLARE' transaction: %v\n", env, err)
		}
		if tx.Transaction.Type != DECLARE {
			t.Errorf("%s: incorrect declare transaction: %v\n", env, tx)
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

		counterClass := types.ContractClass{}
		err := json.Unmarshal(counterCompiled, &counterClass)
		if err != nil {
			t.Fatalf("could not parse contract: %v\n", err)
		}
		tx, err := gw.Deploy(context.Background(), counterClass, types.DeployRequest{
			ContractAddressSalt: "0x1",
			ConstructorCalldata: []string{},
		})
		if err != nil {
			t.Fatalf("testnet: could not deploy contract: %v\n", err)
		}
		fmt.Println("txHash: ", tx.TransactionHash)
	}

}

func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Call types.FunctionCall
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Call: types.FunctionCall{
					ContractAddress:    types.HexToHash(counterAddress),
					EntryPointSelector: "get_count",
					Calldata:           []string{},
				},
			},
		},
		"mainnet": {},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		provider := testConfig.client
		ctx := context.Background()
		tx, err := provider.Call(ctx, test.Call, "latest")
		if err != nil {
			t.Fatalf("could not call contract: %v\n", err)
		}
		fmt.Printf("tx: %+v\n", tx)
		if len(tx) == 0 {
			t.Fatalf("error in tx: %+v\n", tx)
		}
	}
}
