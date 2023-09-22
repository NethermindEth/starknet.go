package gateway_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/rpc"
	devtest "github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	counterCompiled = artifacts.CounterCompiled
	counterAddress  = "0x0"
	accountCompiled = artifacts.AccountCompiled
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
		accountClass := rpc.DeprecatedContractClass{}
		err := json.Unmarshal(accountCompiled, &accountClass)
		if err != nil {
			t.Fatalf("could not parse contract: %v\n", err)
		}
		declareTx, err := gw.Declare(context.Background(), accountClass, gateway.DeclareRequest{})
		if err != nil {
			t.Errorf("%s: could not 'DECLARE' contract: %v\n", env, err)
			return
		}

		tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: declareTx.TransactionHash})
		if err != nil {
			t.Errorf("%s: could not get 'DECLARE' transaction: %v\n", env, err)
		}
		if tx.Transaction.Type != gateway.DECLARE {
			t.Errorf("%s: incorrect declare transaction: %v\n", env, tx)
		}
	}
}

func TestDeployCounterContract(t *testing.T) {
	t.Skip() // TODO: use account
	testConfig := beforeEach(t)

	type testSetType struct{}
	testSet := map[string][]testSetType{
		"devnet":  {{}},
		"mainnet": {},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for range testSet {

		gw := testConfig.client

		counterClass := rpc.DeprecatedContractClass{}
		err := json.Unmarshal(counterCompiled, &counterClass)
		if err != nil {
			t.Fatalf("could not parse contract: %v\n", err)
		}
		tx, err := gw.Deploy(context.Background(), counterClass, rpc.DeployAccountTxn{
			ContractAddressSalt: utils.TestHexToFelt(t, "0x01"),
			ConstructorCalldata: []*felt.Felt{},
		})
		if err != nil {
			t.Fatalf("testnet: could not deploy contract: %v\n", err)
		}
		fmt.Println("txHash: ", tx.TransactionHash)
	}
}

func TestDeployAccountContract(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Salt                    string
		PrivateKey              string
		PublicKey               string
		ExpectedTxForClass      string
		ExpectedTxForContract   string
		ExpectedClassHash       string
		ExpectedContractAddress string
	}

	// This is an example of the DEPLOY_ACCOUNT payload to the gateway
	// {
	// 	"signature": [
	// 		"2150375523489051434313910706914087130364875815064017773574340106797124994336",
	// 		"1209118725260616106547286820901693710951319632876052767730937249518036922842"
	// 	],
	// 	"class_hash": "0x1fac3074c9d5282f0acc5c69a4781a1c711efea5e73c550c5d9fb253cf7fd3d",
	// 	"constructor_calldata": [
	// 		"1817738695255614008829969791198936129031107774334328524814838251584466428552"
	// 	],
	// 	"version": "0x1",
	// 	"max_fee": "0x1fa865ed09401",
	// 	"contract_address_salt": "0x3ae263cecc600734abd80d26629a09578bcfbef768354085595800ff1b05efb",
	// 	"nonce": "0x0",
	// 	"type": "DEPLOY_ACCOUNT"
	// }

	// To load the account, run:
	// curl -X POST localhost:5050/mint
	// -d'{
	//   "address": "0x0655f6d28c50e8187929a14639b73797cb9174e03aefe3908047b57a28688b34",
	// "amount": 1000000000000000000000 }' \
	// -H 'Content-Type: application/json'

	// The response just looks like that:
	// {
	//   "address": "0x0655f6d28c50e8187929a14639b73797cb9174e03aefe3908047b57a28688b34",
	//   "code": "TRANSACTION_RECEIVED",
	//   "transaction_hash": "0xb48337f8d256e3b44c635c0160de2d9fba0ef5655d9b4bce864d6c711ef343"
	// }

	testSet := map[string][]testSetType{
		"devnet": {{
			Salt:                    "0x0",
			PrivateKey:              "0x1",
			PublicKey:               "0x1ef15c18599971b7beced415a40f0c7deacfd9b0d1819e03d723d8bc943cfca",
			ExpectedClassHash:       "0x4529debda559ad97e5fb30c50a6cd12f498a562f5c31b636c6bd157abdb3dfb",
			ExpectedContractAddress: "0x02f10e67bc68db8d0dac316512e390721df3e397fdc147594944da58b17ec9b7",
		}},
	}[testEnv]

	for _, test := range testSet {
		gw := testConfig.client
		// Step 1: deploy the a class
		accountClass := rpc.DeprecatedContractClass{}
		err := json.Unmarshal(accountCompiled, &accountClass)
		if err != nil {
			t.Fatalf("could not parse account: %v\n", err)
		}
		tx, err := gw.Declare(context.Background(), accountClass, gateway.DeclareRequest{})
		if err != nil {
			t.Fatalf("could not declare contract: %v\n", err)
		}
		fmt.Println("txHash: ", tx.TransactionHash)
		_, receipt, err := gw.WaitForTransaction(context.Background(), tx.TransactionHash, 3, 10)
		if err != nil {
			t.Fatalf("could not declare contract: %v\n", err)
		}
		if receipt.Status != types.TransactionAcceptedOnL1 && receipt.Status != types.TransactionAcceptedOnL2 {
			t.Fatalf("unexpected status: %s\n", receipt.Status)
		}
		if test.ExpectedClassHash != tx.ClassHash {
			t.Fatalf("unexpected class hash: %s, instead %s\n", test.ExpectedClassHash, tx.ClassHash)
		}
		// Step 4: send some eth
		mint, err := devtest.NewDevNet().Mint(utils.TestHexToFelt(t, test.ExpectedContractAddress), big.NewInt(int64(1000000000000000000)))
		if err != nil {
			t.Fatalf("could not declare contract: %v\n", err)
		}
		fmt.Println(mint.NewBalance)
		// Step 5: deploy the account
		resp, err := gw.DeployAccount(context.Background(), types.DeployAccountRequest{
			MaxFee:              big.NewInt(10000000000000000),
			Version:             big.NewInt(1),
			ContractAddressSalt: "0x0",
			ConstructorCalldata: []string{test.PublicKey},
			ClassHash:           tx.ClassHash,
		})
		if err != nil {
			t.Fatalf("could not declare contract: %v\n", err)
		}
		// Step 6: make sure the account is deployed
		fmt.Println("txHash: ", resp.TransactionHash)
		if resp.ContractAddress != test.ExpectedContractAddress {
			t.Fatalf("expecting contract %s, instead %s\n", test.ExpectedContractAddress, resp.ContractAddress)
		}
		_, receipt, err = gw.WaitForTransaction(context.Background(), resp.TransactionHash, 3, 10)
		if err != nil {
			t.Fatalf("could not declare contract: %v\n", err)
		}
		if receipt.Status != types.TransactionAcceptedOnL1 && receipt.Status != types.TransactionAcceptedOnL2 {
			t.Fatalf("unexpected status: %s\n", receipt.Status)
		}
		if test.ExpectedClassHash != tx.ClassHash {
			t.Fatalf("unexpected class hash: %s, instead %s\n", test.ExpectedClassHash, tx.ClassHash)
		}
		// Step 7: Get the block from the TX
	}
}

func TestCall(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Call rpc.FunctionCall
	}
	testSet := map[string][]testSetType{
		"devnet": {
			{
				Call: rpc.FunctionCall{
					ContractAddress:    utils.TestHexToFelt(t, counterAddress),
					EntryPointSelector: types.GetSelectorFromNameFelt("get_count"),
					Calldata:           []*felt.Felt{},
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
