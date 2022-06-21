package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	// "math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	// "time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

const (
	FEE_MARGIN float64 = 1.15
	PEDERSON_JSON string = "pedersen_params.json"
)

var (
	snTest         StarknetTest
	snTransactions []string
	testProxy	*RawContractDefinition
	testImplementation	*RawContractDefinition

	_, b, _, _   = runtime.Caller(0)
	projectRoot  = strings.TrimRight(filepath.Dir(b), "gateway")
)

type StarknetTest struct {
	Environments  []TestEnvironment `json:"environments"`
}

type TestEnvironment struct {
	Chain      string `json:"chain"`
	Accounts []struct {
		Address      string              `json:"address,omitempty"`
		Public       string              `json:"public_key,omitempty"`
		Private      string              `json:"private_key,omitempty"`
		Transactions []types.Transaction `json:"transactions"`
	} `json:"accounts"`
	ContractAddresses types.ContractAddresses `json:"contractAddresses"`
}

// requires starknet-devnet to be running and accessible on port 5000
// and seed for accounts to be specified to 0
// ex: starknet-devnet --port 5000 --seed 0
// (ref: https://github.com/Shard-Labs/starknet-devnet)
func init() {
	testFile, err := os.Open(projectRoot + "/gateway/starknet_test.json")
	if err != nil {
		panic(err.Error())
	}

	defer testFile.Close()

	raw, _ := ioutil.ReadAll(testFile)
	json.Unmarshal(raw, &snTest)

	gw := NewClient(WithChain("main"))
	testProxy, err = gw.FullContract(context.Background(), snTest.Environments[2].Accounts[0].Address)
	if err != nil {
		panic(err.Error())
	}

	implResp, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress: snTest.Environments[2].Accounts[0].Address,
		EntryPointSelector: "get_implementation",
	}, "")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("IMPL: ", implResp[0])
	testImplementation, err = gw.FullContract(context.Background(), implResp[0])
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("testCont, err: ", err, testImplementation)
}

func TestContractAddresses(t *testing.T) {
	for _, env := range snTest.Environments {
		if env.Chain != "devnet" {
			gw := NewClient(WithChain(env.Chain))
			addresses, err := gw.ContractAddresses(context.Background())
			if err != nil {
				t.Errorf("%s: could not get starknet addresses - \n%v\n", env.Chain, err)
			}
	
			if strings.ToLower(addresses.Starknet) != env.ContractAddresses.Starknet {
				t.Errorf("%s: fetched incorrect addresses - \n%s %s\n", env.Chain, strings.ToLower(addresses.Starknet), env.ContractAddresses.Starknet)
			}
		}
	}
}

func TestCall(t *testing.T) {
	for _, env := range snTest.Environments {
		gw := NewClient(WithChain(env.Chain))
		for _, testAccount := range env.Accounts {
			call := types.FunctionCall{
				ContractAddress:    testAccount.Address,
			}
			if env.Chain == "devnet" {
				call.EntryPointSelector = "get_public_key"
			} else {
				call.EntryPointSelector = "get_signer"
			}
			resp, err := gw.Call(context.Background(), call, "")
			if err != nil {
				t.Errorf("%s: could 'Call' deployed contract: %v\n", env.Chain, err)
			}
			if len(resp) == 0 {
				t.Errorf("%s: could get signing key for account: %v\n", env.Chain, err)
			}
			
			if resp[0] != testAccount.Public {
				t.Errorf("%s: signing key is incorrect: \n%s %v\n", env.Chain, resp[0], testAccount.Public)
			}
		}
	}
}

func TestExecute(t *testing.T) {
	for _, env := range snTest.Environments {
		// signature scheme for devnet is not curently compatible
		if env.Chain != "devnet" {
			for _, testAccount := range env.Accounts {
				if testAccount.Private != "" {
					curve, err := caigo.SC(caigo.WithConstants(projectRoot + PEDERSON_JSON))
					if err != nil {
						t.Errorf("Could not init with constant points: %v\n", err)
					}
					
					account, err := caigo.NewAccount(&curve, testAccount.Private, testAccount.Address, NewProvider(WithChain(env.Chain)))
					if err != nil {
						t.Errorf("Could not create account: %v\n", err)
					}
					
					feeEstimate, err := account.EstimateFee(context.Background(), testAccount.Transactions)
					if err != nil {
						t.Errorf("Could not estimate fee for transaction: %v\n", err)
					}
					fee := &types.Felt{
						Int: big.NewInt(int64(float64(feeEstimate.Amount) * FEE_MARGIN)),
					}
					
					// txResp, err := account.Execute(context.Background(), fee, testAccount.Transactions)
					// if err != nil {
					// 	t.Errorf("Could not execute test transaction: %v\n", err)
					// }
					fmt.Println("RESP: ", fee)
					
					// snTransactions = append(snTransactions, txResp.TransactionHash)
				}
			}
		}
	}
}

func TestDeploy(t *testing.T) {
	// for _, env := range snTest.Environments {

	// }
	// ctx := context.Background()
	// gw := NewClient(WithChain("local"))

	// salt := rand.New(rand.NewSource(time.Now().UnixNano()))
	// deployTx, err := gw.Deploy(ctx, testContract, types.DeployRequest{
	// 	ContractAddressSalt: fmt.Sprintf("0x%x", salt.Intn(1000000)),
	// 	ConstructorCalldata: []string{},
	// })
	// if err != nil {
	// 	t.Errorf("Could not deploy contract: %v\n", err)
	// }

	// tx, err := gw.Transaction(ctx, TransactionOptions{TransactionHash: deployTx.TransactionHash})
	// fmt.Println("TX: ", tx)
	// if err != nil {
	// 	t.Errorf("Could not get tx: %v\n", err)
	// }
	// if tx.Transaction.Type != DEPLOY || tx.Status != types.ACCEPTED_ON_L2.String() {
	// 	t.Errorf("Incorrect deployment transaction: %+v\n", tx)
	// }
}

// func TestDevnetDeclare(t *testing.T) {
// 	ctx := context.Background()
// 	gw := NewClient(WithChain("local"))

// 	declareTx, err := gw.Declare(ctx, testProxy, types.DeclareRequest{})
// 	if err != nil {
// 		t.Errorf("Could not 'DECLARE' contract: %v\n", err)
// 	}

// 	tx, err := gw.Transaction(ctx, TransactionOptions{TransactionHash: declareTx.TransactionHash})
// 	if err != nil {
// 		t.Errorf("Could not get 'DECLARE' transaction: %v\n", err)
// 	}
// 	if tx.Transaction.Type != DECLARE || tx.Status != types.ACCEPTED_ON_L2.String() {
// 		t.Errorf("Incorrect delcare transaction: %v\n", tx)
// 	}
// }
